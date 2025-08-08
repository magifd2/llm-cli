package cmd

import (
	"fmt"
	"strings"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/magifd2/llm-cli/internal/llm"
	"github.com/magifd2/llm-cli/internal/llm/bedrock"
	"github.com/magifd2/llm-cli/internal/llm/ollama"
	"github.com/magifd2/llm-cli/internal/llm/openai"
	"github.com/magifd2/llm-cli/internal/llm/openai2"
	"github.com/magifd2/llm-cli/internal/llm/vertexai"
	"github.com/magifd2/llm-cli/internal/llm/vertexai2"
)

// providerFactory defines the function signature for creating a new llm.Provider.
// It takes a profile and returns an implementation of the provider interface.
type providerFactory func(profile config.Profile) llm.Provider

// providerRegistry holds the mapping from a provider name (string) to its factory function.
// This is the central, explicit registry for all supported providers.
var providerRegistry = map[string]providerFactory{
	"ollama": func(p config.Profile) llm.Provider { return &ollama.Provider{Profile: p} },
	"openai": func(p config.Profile) llm.Provider { return &openai.Provider{Profile: p} },
	"openai2": func(p config.Profile) llm.Provider { return &openai2.Provider{Profile: p} },
	"bedrock": func(p config.Profile) llm.Provider {
		// The Bedrock provider has sub-types based on the model ID.
		if strings.HasPrefix(p.Model, "amazon.nova") {
			return &bedrock.NovaProvider{Profile: p}
		}
		// In the future, other Bedrock model families (e.g., Claude) can be added here.
		return nil // Return nil if no sub-type matches
	},
	"vertexai":  func(p config.Profile) llm.Provider { return &vertexai.Provider{Profile: p} },
	"vertexai2": func(p config.Profile) llm.Provider { return &vertexai2.Provider{Profile: p} },
}

// GetProvider retrieves a provider instance based on the provider name in the profile.
// It looks up the provider in the registry and uses the factory function to create it.
func GetProvider(profile config.Profile) (llm.Provider, error) {
	factory, ok := providerRegistry[profile.Provider]
	if !ok {
		return nil, fmt.Errorf("provider '%s' not recognized", profile.Provider)
	}

	provider := factory(profile)
	if provider == nil {
		return nil, fmt.Errorf("model '%s' is not supported by the '%s' provider yet", profile.Model, profile.Provider)
	}

	return provider, nil
}
