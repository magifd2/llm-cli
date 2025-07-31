# ビルド手順

このドキュメントは、`llm-cli` をソースコードからビルドする方法について説明します。

## 前提条件

*   [Go](https://go.dev/doc/install) (バージョン 1.21 以降を推奨)
*   [Git](https://git-scm.com/)
*   `make` コマンド
    *   macOS/Linuxでは標準で利用可能です。
    *   Windowsでは [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm) などをインストールしてください。

## ビルド

本プロジェクトでは `Makefile` を使用したビルドを推奨します。

### 1. リポジトリのクローン

```bash
git clone https://github.com/magifd2/llm-cli.git
cd llm-cli
```

### 2. ビルドコマンド

以下の `make` コマンドを利用できます。ビルドされたバイナリは `bin/` ディレクトリ以下に生成されます。

*   **`make build`**
    *   現在利用しているOS・アーキテクチャ向けのバイナリを一つだけビルドします。開発中に動作確認するのに便利です。

*   **`make cross-compile`**
    *   配布用に、複数のOS・アーキテクチャ向けのバイナリを一度にビルドします。成果物は `bin/` 以下のプラットフォーム別ディレクトリに生成されます。
        *   `bin/darwin-universal/llm-cli` (macOS Universal Binary)
            *   **注意**: このターゲットはmacOS環境でのみ実行可能です。
        *   `bin/linux-amd64/llm-cli` (Linux amd64)
        *   `bin/windows-amd64/llm-cli.exe` (Windows amd64)

*   **`make test`**
    *   プロジェクトのテストを実行します。

*   **`make clean`**
    *   `bin/` ディレクトリとビルドキャッシュを削除します。
