# CHANGELOG

## v0.0.5 - 2025-08-02

### ✨ Features
*   **Google Cloud Vertex AI プロバイダーのサポート**: Vertex AI との対話に対応しました。
*   **`profile add` コマンドの機能拡張**: `profile add` コマンドで、プロバイダー、モデル、エンドポイント、APIキー、AWS認証情報、GCPプロジェクトID、ロケーション、クレデンシャルファイルパスなどのパラメータを一括して指定できるようになりました。

### ♻️ Refactor
*   **Vertex AI SDK の移行**: `google.golang.org/genai` SDK の最新バージョンに移行しました。サービスアカウント認証の修正、`Client` オブジェクトの正しい利用、ストリーミングイテレータの適切な処理を含みます。
*   **クレデンシャルファイルパスの実行時展開**: `credentials_file` のパス展開を、設定時ではなく実行時に行うように変更しました。これにより、動的なホームディレクトリ環境でも柔軟に対応できるようになりました。

### 📝 Documentation
*   **開発経緯の更新**: `DEVELOPMENT_LOG.md` にVertex AI SDK移行の詳細な経緯とSDKの現状に関する記述を追加しました。
*   **関連ドキュメントの更新**: `README.ja.md` および `README.en.md` を、Vertex AI プロバイダーの追加と `profile add` コマンドの機能拡張、システムプロンプトの扱いについて更新しました。
*   **プロバイダー開発ガイドの修正**: `DEVELOPING_PROVIDERS.ja.md` および `DEVELOPING_PROVIDERS.en.md` から、特定のプロバイダー実装に関する記述を削除しました。
*   **変更履歴の更新**: `CHANGELOG.ja.md` および `CHANGELOG.en.md` を更新しました。

### ♻️ Refactor
*   **コード監査と品質改善**: コード全体の監査を実施し、潜在的な不具合や脆弱性を修正しました。`profile edit`のコマンドインジェクション対策、設定パス管理の一元化、エラーメッセージの改善など、堅牢性と保守性を向上させました。

## v0.0.4 - 2025-08-01

### 🐛 Bug Fixes
*   **APIエラーハンドリングの修正**: ストリーミングモードでAPIエラーが発生した際に、エラーを検知できずに正常終了してしまう問題を修正しました。非同期処理の競合状態を解消し、エラーハンドリングをより堅牢にしました。

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
