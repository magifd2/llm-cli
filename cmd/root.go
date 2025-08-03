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
	"github.com/spf13/cobra"
)



// rootCmd represents the base command for the llm-cli application.
// It defines the application's name, version, short description, and long description.
// rootCmdはllm-cliアプリケーションのベースコマンドです。
// アプリケーション名、バージョン、短い説明、長い説明を定義します。
var rootCmd = &cobra.Command{
	Use:   "llm-cli",
	Version: "v0.0.9",
	Short: "A CLI tool for interacting with various LLM providers.",
	Long: `llm-cli is a powerful command-line interface tool designed for seamless interaction with various Large Language Model (LLM) providers.
It supports providers such as Ollama, LM Studio (OpenAI compatible API), Amazon Bedrock, and Google Cloud Vertex AI.

Key Features:
- Interact with multiple LLM providers from your terminal.
- Manage different LLM configurations using profiles, allowing for easy switching between models and settings.
- Send prompts via command-line arguments, files, or standard input.
- Receive streaming responses for real-time interaction.

This tool simplifies the process of experimenting with and utilizing different LLMs directly from your command line.`,
	SilenceUsage: true, // Suppress usage message on error
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main() and should only be called once.
// Executeは全ての子コマンドをルートコマンドに追加し、フラグを適切に設定します。
// これはmain.main()から呼び出され、一度だけ呼び出されるべきです。
func Execute() error {
	return rootCmd.Execute()
}

// init function is called before main.
// It's used to define flags and configuration settings for the root command.
// init関数はmainの前に呼び出されます。
// ルートコマンドのフラグと設定を定義するために使用されます。
func init() {
}