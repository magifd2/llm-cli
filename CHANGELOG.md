# CHANGELOG

## v0.0.1 - 2025-07-31

### ✨ Features

*   **LLMとの対話機能**:
    *   OllamaおよびLM Studio (OpenAI互換API) との対話に対応。
    *   コマンドライン引数、ファイル、標準入力からのプロンプト入力に対応。
    *   ストリーミング応答表示 (`--stream` フラグ) に対応。
*   **プロファイル管理機能**:
    *   複数のLLM設定をプロファイルとして管理 (`profile list`, `profile use`, `profile add`, `profile set`, `profile remove`, `profile edit`)。
*   **ビルドシステム**:
    *   Goモジュールの初期化とCobra CLIフレームワークの導入。
    *   `Makefile` によるビルド、テスト、クロスコンパイル (macOS Universal, Linux, Windows) に対応。
    *   ビルド成果物をプラットフォーム別ディレクトリに整理。

### 📝 Documentation

*   `README.md` の作成と機能説明の追加。
*   `BUILD.md` の作成とビルド手順の詳細化。
*   `README.md` にGeminiによる開発支援とMITライセンスの明記。
*   `BUILD.md` および `README.md` にコードレビューで指摘された改善点を反映。
