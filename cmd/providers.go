package cmd

import (
	"fmt"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/magifd2/llm-cli/internal/llm"
	"github.com/magifd2/llm-cli/internal/llm/bedrock"
	"github.com/magifd2/llm-cli/internal/llm/mock"
	"github.com/magifd2/llm-cli/internal/llm/ollama"
	"github.com/magifd2/llm-cli/internal/llm/openai"
	"github.com/magifd2/llm-cli/internal/llm/openai2"
	"github.com/magifd2/llm-cli/internal/llm/vertexai"
	"github.com/magifd2/llm-cli/internal/llm/vertexai2"
)

// providerFactory defines the function signature for creating a new llm.Provider.
// It takes a profile and returns an implementation of the provider interface and an error.
type providerFactory func(profile config.Profile) (llm.Provider, error)

// providerRegistry holds the mapping from a provider name (string) to its factory function.
// This is the central, explicit, and consistent registry for all supported providers.
var providerRegistry = map[string]providerFactory{
	"ollama":    ollama.NewProvider,
	"openai":    openai.NewProvider,
	"openai2":   openai2.NewProvider,
	"bedrock":   bedrock.NewProvider,
	"vertexai":  vertexai.NewProvider,
	"vertexai2": vertexai2.NewProvider,
	"mock":      mock.NewProvider,
}

// GetProvider retrieves a provider instance based on the provider name in the profile.
// It looks up the provider in the registry and uses the factory function to create it.
func GetProvider(profile config.Profile) (llm.Provider, error) {
	factory, ok := providerRegistry[profile.Provider]
	if !ok {
		return nil, fmt.Errorf("provider '%s' not recognized", profile.Provider)
	}

	return factory(profile)
}
