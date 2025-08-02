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

## インストール

`llm-cli` は、`Makefile` のターゲットを使用してインストールおよびアンインストールできます。

### `make install`

このターゲットは、`llm-cli` バイナリをビルドし、指定されたディレクトリにZshシェル補完スクリプトと共にインストールします。デフォルトのインストールパスは `/usr/local/bin` です。

*   **デフォルトインストール（システム全体）:**
    `llm-cli` を `/usr/local/bin` にインストールする場合（`sudo` が必要です）:
    ```bash
    sudo make install
    ```

*   **ユーザーローカルインストール:**
    `llm-cli` を `~/bin` にインストールする場合（ルート以外のユーザーに推奨。`~/bin` が `PATH` に含まれていることを確認してください）:
    ```bash
    make install PREFIX=~
    ```

*   **カスタムディレクトリへのインストール:**
    `llm-cli` をカスタムディレクトリ（例: `/opt/llm-cli`）にインストールする場合:
    ```bash
    sudo make install PREFIX=/opt/llm-cli
    ```

インストール後、Zshユーザーは補完スクリプトを有効にするために `compinit` を実行するか、シェルを再起動する必要がある場合があります。

### `make uninstall`

このターゲットは、インストールディレクトリから `llm-cli` バイナリとそれに関連する補完スクリプトを削除します。インストール時に使用した `PREFIX` と同じ値を指定することが重要です。

*   **デフォルトアンインストール:**
    ```bash
    sudo make uninstall
    ```

*   **ユーザーローカルアンインストール:**
    ```bash
    make uninstall PREFIX=~
    ```

*   **カスタムディレクトリからのアンインストール:**
    ```bash
    sudo make uninstall PREFIX=/opt/llm-cli
    ```

**注意:** アンインストールプロセスでは、`~/.config/llm-cli/config.json` にある設定ファイルは削除されません。これらのファイルにはLLMプロファイルが含まれており、インストール/アンインストールを跨いで保持されます。
