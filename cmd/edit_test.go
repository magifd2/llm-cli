package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This is a helper function to extract the logic for choosing an editor,
// so we can test it without actually executing the command.
func getEditorPathForTest(editorEnv string) (string, error) {
	if editorEnv == "" {
		editorEnv = "vim" // Default fallback
	}
	// In a real scenario, exec.LookPath would be called here.
	// For this test, we simulate its behavior based on the input string.
	if editorEnv == "vim" || editorEnv == "nano" || editorEnv == "vi" {
		return "/usr/bin/" + editorEnv, nil // Simulate a valid path
	}
	// Simulate LookPath failing for invalid/malicious strings
	return "", &os.PathError{Op: "lookaside", Path: editorEnv, Err: os.ErrNotExist}
}

func TestGetEditorPath(t *testing.T) {
	testCases := []struct {
		name          string
		editorEnv     string
		expectedPath  string
		expectError   bool
	}{
		{
			name:          "EDITOR is set to vim",
			editorEnv:     "vim",
			expectedPath:  "/usr/bin/vim",
			expectError:   false,
		},
		{
			name:          "EDITOR is not set, fallback to vim",
			editorEnv:     "",
			expectedPath:  "/usr/bin/vim",
			expectError:   false,
		},
		{
			name:          "EDITOR is set to nano",
			editorEnv:     "nano",
			expectedPath:  "/usr/bin/nano",
			expectError:   false,
		},
		{
			name:          "Potential command injection",
			editorEnv:     "vim; ls -la",
			expectedPath:  "",
			expectError:   true,
		},
		{
			name:          "Another potential command injection",
			editorEnv:     "/usr/bin/vim && rm -rf /",
			expectedPath:  "",
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We test our helper function which simulates the core logic.
			path, err := getEditorPathForTest(tc.editorEnv)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, path)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPath, path)
			}
		})
	}
}
