# CHANGELOG

## v0.0.3 - 2025-08-01

### 🚨 Breaking Changes
*   **コマンド名とフラグの変更**: `ask` コマンドが `prompt` に変更されました。また、`--prompt` フラグは `--user-prompt` にリネームされ、`--prompt-file` は `--user-prompt-file` にリネームされました。既存のスクリプトやワークフローを更新する必要があります。

### ✨ Features
*   **コマンド名とフラグのリファクタリング**: `ask` コマンドを `prompt` に変更し、`--prompt` フラグを `--user-prompt` にリネームしました。また、`--user-prompt` (`-p`), `--system-prompt` (`-P`), `--user-prompt-file` (`-f`), `--system-prompt-file` (`-F`) の省略形を追加しました。

### 📝 Documentation
*   **開発ログと開発計画の追加**: プロジェクトの経緯と今後の計画を記録するため、`DEVELOPMENT_LOG.md` と `DEVELOPMENT_PLAN.md` を追加しました。
*   **開発ルールの更新**: `GEMINI.md` に開発ルール（絶対パスの使用、根本原因の修正、ドキュメント更新、コミット前のコミット、コード品質とセキュリティ、コメントの原則）を追加しました。

## v0.0.2 - 2025-08-01

### ✨ Features

*   **Amazon Bedrock Novaモデルのサポート**: `amazon.nova-lite-v1:0` などのNovaモデルとの対話に対応しました。

### ♻️ Refactor

*   **Bedrockプロバイダーのリファクタリング**: Bedrockプロバイダーの内部実装を、NovaモデルのMessages API仕様に厳密に合わせて再設計しました。
    *   リクエスト/レスポンス構造体 (`novaMessageContent`, `novaMessage`, `novaSystemPrompt`, `novaMessagesAPIRequest`, `novaCombinedAPIResponse`, `novaMessagesAPIStreamChunk`) を更新。
    *   プロンプトおよび推論パラメータのハンドリングをNova APIに合わせて調整。
    *   ストリーミング応答のパースロジックを修正。
*   **プロバイダー選択ロジックの改善**: `cmd/ask.go` で、`bedrock` プロバイダーが選択された際に、モデルIDのプレフィックス (`amazon.nova`) に基づいて `NovaBedrockProvider` を動的に選択するように変更しました。

### 🐛 Bug Fixes

*   **プロンプトバリデーションの修正**: `cmd/ask.go` でプロンプトのバリデーションを一元的に行うように修正し、`internal/llm/bedrock_nova.go` から冗長なバリデーションを削除しました。これにより、`--prompt` または `--prompt-file`、あるいは標準入力からのプロンプトが必須となりました。

### 📝 Documentation

*   **ドキュメントの更新**: `README.md` および `BUILD.md` を、新しいBedrockのセットアップ手順とビルドプロセスの変更に合わせて更新しました。
*   **Makefileの改善**: `make all` コマンドがクロスコンパイルも実行するように修正し、ビルド出力ディレクトリの整理を行いました。

## v0.0.1 - 2025-07-31

### ✨ Features

*   **LLMとの対話機能**:\
    *   OllamaおよびLM Studio (OpenAI互換API) との対話に対応。\
    *   コマンドライン引数、ファイル、標準入力からのプロンプト入力に対応。\
    *   ストリーミング応答表示 (`--stream` フラグ) に対応。
*   **プロファイル管理機能**:\
    *   複数のLLM設定をプロファイルとして管理 (`profile list`, `profile use`, `profile add`, `profile set`, `profile remove`, `profile edit`)。
*   **ビルドシステム**:\
    *   Goモジュールの初期化とCobra CLIフレームワークの導入。\
    *   `Makefile` によるビルド、テスト、クロスコンパイル (macOS Universal, Linux, Windows) に対応。\
    *   ビルド成果物をプラットフォーム別ディレクトリに整理。

### 📝 Documentation

*   `README.md` の作成と機能説明の追加。
*   `BUILD.md` の作成とビルド手順の詳細化。
*   `README.md` にGeminiによる開発支援とMITライセンスの明記。
*   `BUILD.md` および `README.md` にコードレビューで指摘された改善点を反映。
