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
    *   現在利用しているOS・アーキテクチャ向けのバイナリを一つだけビルドします。ビルドされたバイナリは `bin/<OS>-<ARCH>/llm-cli` に配置されます。開発中に動作確認するのに便利です。

*   **`make cross-compile`**
    *   配布用に、複数のOS・アーキテクチャ向けのバイナリを一度にビルドし、圧縮ファイルを作成します。成果物は `bin/` ディレクトリ以下に生成されます。
        *   `bin/llm-cli-darwin-universal.tar.gz` (macOS Universal Binary)
        *   `bin/llm-cli-linux-amd64.tar.gz` (Linux amd64)
        *   `bin/llm-cli-windows-amd64.zip` (Windows amd64)

*   **`make all`**
    *   `make build` と `make cross-compile` の両方を実行します。現在のOS・アーキテクチャ向けのバイナリと、すべてのクロスコンパイル済みバイナリおよびアーカイブを作成します。
*   **`make test`**
    *   プロジェクトのテストを実行します。

*   **`make clean`**
    *   `bin/` ディレクトリとビルドキャッシュを削除します。