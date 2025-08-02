# llm-cli

`llm-cli` は、ローカルおよびリモートのLLM（大規模言語モデル）と直接対話するためのコマンドラインインターフェースツールです。Ollama、LM Studio、Amazon Bedrockのような様々なプロバイダーに対して、統一された方法でプロンプトを送信したり、設定を管理したりする機能を提供します。

## 主な特徴

*   **マルチプロバイダー対応**: Ollama、LM Studio（およびその他のOpenAI互換API）、Amazon Bedrockとシームレスに連携します。
*   **プロファイル管理**: 複数のLLM設定（エンドポイント、モデル、APIキー）をプロファイルとして保存し、簡単に切り替えられます。
*   **柔軟な入力**: コマンドライン引数、ファイル、標準入力（パイプ）からプロンプトを渡せます。
*   **ストリーミング表示**: `--stream` フラグを使用することで、LLMからの応答をリアルタイムで表示します。
*   **シングルバイナリ**: 設定ファイルを除き、単一の実行ファイルで動作するため、配布や利用が簡単です。

## インストール

`llm-cli` は、提供されている `Makefile` を使用して簡単にインストールできます。

### `make install` を使用したインストール

この方法では、`llm-cli` バイナリをビルドし、指定されたディレクトリにZshシェル補完スクリプトと共にインストールします。

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
    `llm-cli` をカスタムディレクトリ（例: `/opt/llm-cli/bin`）にインストールする場合:
    ```bash
    sudo make install PREFIX=/opt/llm-cli
    ```

インストール後、Zshユーザーは補完スクリプトを有効にするために `compinit` を実行するか、シェルを再起動する必要がある場合があります。

### アンインストール

`llm-cli` とその補完スクリプトをアンインストールするには、インストール時に使用したのと同じ `PREFIX` を指定して `make uninstall` を使用します。

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

**注意:** アンインストールプロセスでは、`~/.config/llm-cli/config.json` にある設定ファイルは削除されません。

## クイックスタート

インストールと設定が完了すれば、すぐにLLMとの対話を開始できます。

```bash
# デフォルトのLLMに簡単なプロンプトを送信
llm-cli prompt "地球と月の距離はどのくらいですか？"

# ストリーミングで応答を取得
llm-cli prompt "音楽を発見したロボットの短編小説を教えてください。" --stream
```

## 設定

`llm-cli` のすべての設定は、`~/.config/llm-cli/config.json` にある単一の設定ファイルで管理されます。`llm-cli profile edit` でこのファイルを直接編集することもできますが、`profile` サブコマンド群を使用することが推奨されます。

### プロバイダー別のセットアップ

#### 1. Ollama

Ollamaをデフォルトのアドレス（`http://localhost:11434`）で実行している場合、`llm-cli` は追加設定なしで動作します。`default` プロファイルがこの設定に最適化されています。

Ollamaで取得した特定のモデルを使用するには：
```bash
# defaultプロファイルに切り替え（まだの場合）
llm-cli profile use default

# 使用したいモデルを設定
llm-cli profile set model "llama3"
```

#### 2. LM Studio (およびその他のOpenAI互換API)

LM Studioを使用するには、まずローカルサーバーを起動する必要があります。

1.  **サーバーの起動**: LM Studioで、「Local Server」タブ（`<->` アイコン）に移動します。
2.  **モデルのロード**: モデルを選択してロードし、準備が完了するのを待ちます。
3.  **サーバーの開始**: 「Start Server」ボタンをクリックします。上部に表示されるサーバーURL（例: `http://localhost:1234/v1`）を控えておきます。

次に、`llm-cli` がこのサーバーを使用するように設定します。

```bash
# LM Studio用に新しいプロファイルを追加
llm-cli profile add lmstudio

# プロバイダーを "openai" に設定
llm-cli profile set provider openai

# エンドポイントをLM StudioのURLに設定
llm-cli profile set endpoint "http://localhost:1234/v1"

# ローカルサーバーではモデル名は任意の場合が多いですが、設定は必須です。
# 通常はLM Studioのモデル識別子を使用できます。
llm-cli profile set model "gemma-2-9b-it"

# (任意) OpenAI互換APIが認証を必要とする場合、APIキーを設定
# llm-cli profile set api_key "YOUR_API_KEY"

# 新しく作成したプロファイルに切り替え
llm-cli profile use lmstudio
```

これで、LM Studioのモデルにプロンプトを送信できます。

#### 3. Amazon Bedrock

Amazon Bedrockを利用するには、有効なAWS認証情報とリージョンの指定が必要です。

**認証情報の優先順位:**
1.  `llm-cli` プロファイルに直接設定された認証情報（`aws_access_key_id`, `aws_secret_access_key`）。
2.  標準のAWS SDK認証情報チェーン（環境変数、共有認証情報ファイル、IAMロールなど）。

**設定手順:**

```bash
# Bedrock用に新しいプロファイルを追加
llm-cli profile add bedrock-nova

# プロバイダーを "bedrock" に設定
llm-cli profile set provider bedrock

# 使用したいモデルのモデルIDを設定
llm-cli profile set model "amazon.nova-lite-v1:0"

# モデルを呼び出すAWSリージョンを設定
llm-cli profile set aws_region "us-east-1"

# (任意) 必要に応じて認証情報を直接設定
# llm-cli profile set aws_access_key_id "YOUR_KEY_ID"
# llm-cli profile set aws_secret_access_key "YOUR_SECRET_KEY"

# Bedrockプロファイルに切り替え
llm-cli profile use bedrock-nova
```

**必要なIAMポリシー:**
AWS IDには、Bedrockモデルを呼び出す権限が必要です。

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "bedrock:InvokeModel",
                "bedrock:InvokeModelWithResponseStream"
            ],
            "Resource": "arn:aws:bedrock:us-east-1::foundation-model/amazon.nova-lite-v1:0"
        }
    ]
}
```
*注意: ベストプラクティスとして、`Resource` は必要な特定のモデルに限定することを強く推奨します。*

#### 4. Google Cloud Vertex AI

Google Cloud Vertex AIを利用するには、GCPプロジェクトの設定と認証情報の準備が必要です。

**事前準備:**
1.  Vertex AI を使用する Google Cloud Platform (GCP) プロジェクトが作成済みであること。
2.  対象のGCPプロジェクトで **Vertex AI API** が有効になっていること。
3.  サービスアカウントキーを作成し、**JSON** 形式でダウンロードします。このキーファイルは安全な場所に保管してください。
    *   サービスアカウントには **「Vertex AI ユーザー」** ロールを付与します。

**設定手順:**

```bash
# Vertex AI用に新しいプロファイルを追加（一括設定）
llm-cli profile add my-vertex-ai \
  --provider vertexai \
  --model gemini-1.5-pro-001 \
  --project-id "your-gcp-project-id" \
  --location "us-central1" \
  --credentials-file "~/path/to/your/service-account-key.json"

# 新しく作成したプロファイルに切り替え
llm-cli profile use my-vertex-ai
```

**注意:** `credentials-file` には、サービスアカウントキーのJSONファイルへのパスを `~` を含む形式または絶対パスで指定してください。

**必要なIAMロール:**
サービスアカウントには、Vertex AIモデルを呼び出す権限が必要です。
*   `Vertex AI ユーザー` ロール

**システムプロンプトの扱い:**
Vertex AIのGenAI SDKでは、システムプロンプトに直接対応する機能がありません。そのため、`llm-cli` では、チャットの最初のメッセージとしてシステムプロンプトの内容を送信し、その後にユーザープロンプトの内容を送信することで、擬似的にシステムプロンプトに対応しています。

## コマンドリファレンス


### `llm-cli prompt`

現在アクティブなLLMにプロンプトを送信します。

| フラグ                 | 短縮形 | 説明                                                 |
| -------------------- | ------ | ---------------------------------------------------- |
| `--user-prompt`      | `-p`   | モデルに送信するメインのプロンプトテキスト。         |
| `--user-prompt-file` | `-f`   | ユーザープロンプトを含むファイルへのパス。`-`で標準入力。 |
| `--system-prompt`    | `-P`   | モデルへのオプションのシステムレベルの指示。         |
| `--system-prompt-file`| `-F`   | システムプロンプトを含むファイルへのパス。           |
| `--stream`           |        | 応答をリアルタイムストリームとして表示するかどうか。 |

*プロンプト用フラグが指定されない場合、最初の位置引数がプロンプトとして使用されます。それも無い場合は、標準入力から読み込まれます。*

### `llm-cli profile`

設定プロファイルを管理します。

| サブコマンド | 説明                                                               |
| ---------- | ------------------------------------------------------------------ |
| `list`     | 利用可能な全プロファイルとアクティブなプロファイルを表示します。     |
| `use`      | アクティブなプロファイルを切り替えます。`llm-cli profile use <profile-name>` |
| `add`      | 新しいプロファイルを作成します。パラメータを指定しない場合、デフォルトプロファイルの設定をコピーします。`llm-cli profile add <new-name> [--provider <provider>] [--model <model>] [...]` |
| `set`      | 現在のプロファイルのキーを変更します。`llm-cli profile set <key> <value>` |
| `remove`   | プロファイルを削除します。`llm-cli profile remove <profile-name>`     |
| `edit`     | `config.json` ファイルをデフォルトのテキストエディタで開きます。     |

## コントリビューションと開発

新しい機能の追加やバグ修正などのコントリビューションを歓迎します。

**macOS開発者への注意:** macOSでビルドする場合、`make build` は自動的にユニバーサルバイナリ（`amd64` と `arm64` の両方のアーキテクチャをサポート）を生成し、幅広い互換性を確保します。

コントリビューションに興味がある方は、[コントリビューションガイド](./CONTRIBUTING.md)をご覧ください。

## 謝辞

このプロジェクトは、GoogleのAIアシスタント「Gemini」をコーディングパートナーとして開発されました。

## ライセンス

このプロジェクトはMITライセンスです。詳細は`LICENSE`ファイルをご覧ください。

## セキュリティ

セキュリティ脆弱性の報告方法については、[セキュリティポリシー](./SECURITY.md)をご覧ください。
