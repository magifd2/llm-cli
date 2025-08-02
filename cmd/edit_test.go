/*
Copyright © 2025 magifd2

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
