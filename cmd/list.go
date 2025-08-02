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
	"fmt"
	"os"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

// listCmd represents the 'profile list' command.
// This command lists all configured profiles and indicates the currently active one.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available profiles",
	Long:  `Lists all saved profiles and indicates which one is currently active.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Available profiles:")
		for name, p := range cfg.Profiles {
			activeMarker := " "
			if name == cfg.CurrentProfile {
				activeMarker = "*"
			}
			fmt.Printf("  %s %s (provider: %s, model: %s)\n", activeMarker, name, p.Provider, p.Model)
		}
	},

}

// init function registers the listCmd with the profileCmd.
func init() {
	profileCmd.AddCommand(listCmd)
}
