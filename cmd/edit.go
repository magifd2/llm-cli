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
	os "os"
	"os/exec"
	"path/filepath"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

// editCmd represents the 'profile edit' command.
// This command opens the configuration file in the user's default editor.
// editCmdは'profile edit'コマンドを表します。
// このコマンドは、ユーザーのデフォルトエディタで設定ファイルを開きます。
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the configuration file",
	Long:  `Opens the configuration file in the default editor ($EDITOR).`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runEditCommand(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// runEditCommand contains the core logic for the edit command.
// It determines the editor to use, finds its executable path for security, and opens the config file.
// runEditCommandはeditコマンドの主要なロジックを含みます。
// 使用するエディタを決定し、セキュリティのためにその実行可能パスを見つけ、設定ファイルを開きます。
func runEditCommand() error {
	// Determine the editor to use. Prioritize $EDITOR environment variable, fallback to vim, then nano.
	// 使用するエディタを決定します。$EDITOR環境変数を優先し、vim、nanoの順にフォールバックします。
	editorEnv := os.Getenv("EDITOR")
	if editorEnv == "" {
		editorEnv = "vim"
	}

	// Find the absolute path of the editor executable to prevent command injection.
	// This is a critical security measure.
	// コマンドインジェクションを防ぐために、エディタの実行可能ファイルの絶対パスを見つけます。
	// これは重要なセキュリティ対策です。
	editorPath, err := exec.LookPath(editorEnv)
	if err != nil {
		// If the primary editor is not found, try nano as a fallback.
		// プライマリのエディタが見つからない場合、フォールバックとしてnanoを試します。
		editorPath, err = exec.LookPath("nano")
		if err != nil {
			return fmt.Errorf("EDITOR environment variable not set, and vim/nano not found in PATH")
		}
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("getting config path: %w", err)
	}

	// Ensure the directory exists before trying to open the file
	// ファイルを開く前にディレクトリが存在することを確認します。
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	execCmd := exec.Command(editorPath, configPath)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	// Run the editor command.
	// エディタコマンドを実行します。
	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("opening editor: %w", err)
	}
	return nil
}

// init function registers the editCmd with the profileCmd.
// init関数はeditCmdをprofileCmdに登録します。
func init() {
	profileCmd.AddCommand(editCmd)
}
