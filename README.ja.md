# llm-cli

`llm-cli` は、ローカル（Ollama, LM Studio）またはリモートのLLM（将来的にはOpenAIなど）と、コマンドラインから直接対話するためのCLIツールです。

## 特徴

*   **マルチプロバイダー対応**: Ollama, LM Studio (OpenAI互換API) に対応。
*   **プロファイル管理**: 複数のLLM設定（エンドポイント、モデルなど）をプロファイルとして保存し、簡単に切り替え可能。
*   **柔軟な入力**: コマンドライン引数、ファイル、標準入力（パイプ）からプロンプトを渡せます。
*   **ストリーミング表示**: LLMからの応答をリアルタイムで表示します。
*   **Goによるシングルバイナリ**: 設定ファイル以外は単一の実行ファイルで動作し、簡単に配布できます。

## 使い方

### プロンプトの送信 (必須)

```bash
# シンプルなプロンプト
llm-cli ask --prompt "日本の首都はどこですか？" # --prompt または --prompt-file が必須です

# システムプロンプト付き
llm-cli ask --prompt "自己紹介して" --system-prompt "あなたは猫です。語尾にニャンを付けて話してください。"

# ストリーミング表示
llm-cli ask --prompt "1から100まで数えてください" --stream

# ファイルからプロンプトを読み込む (または標準入力からパイプ)
llm-cli ask --prompt-file ./my_prompt.txt 

# パイプで渡す
echo "この文章を要約して" | llm-cli ask
```

### プロファイルの管理

```bash
# プロファイルの一覧表示
llm-cli profile list

# 新しいプロファイルの追加 (defaultプロファイルをコピーして作成)
llm-cli profile add my-new-profile

# 使用するプロファイルの切り替え
llm-cli profile use my-new-profile

# 現在のプロファイルの設定を変更
llm-cli profile set model "new-model-name"
llm-cli profile set endpoint "http://my-endpoint/v1"

# プロファイルの削除
llm-cli profile remove my-new-profile

# 設定ファイルを直接編集
llm-cli profile edit
```

### Amazon Bedrock の設定

Amazon Bedrock を利用するには、AWSの認証情報とリージョン設定が必要です。
認証情報は、プロファイルに直接設定するか、AWS SDKのデフォルトの認証情報プロバイダチェーン（環境変数、IAMロールなど）を利用できます。

**Bedrockプロファイルの例:**

```bash
# 新しいBedrockプロファイルを追加
llm-cli profile add bedrock-claude

# プロバイダーをbedrockに設定 (NovaモデルはMessages APIを使用します)
llm-cli profile set provider bedrock

# モデルIDを設定 (例: Amazon Nova Lite v1)
llm-cli profile set model amazon.nova-lite-v1:0

# AWSリージョンを設定 (例: ap-northeast-1)
llm-cli profile set aws_region ap-northeast-1

# アクセスキーIDとシークレットアクセスキーを直接設定する場合 (非推奨: 環境変数やIAMロールを推奨)
llm-cli profile set aws_access_key_id YOUR_AWS_ACCESS_KEY_ID
llm-cli profile set aws_secret_access_key YOUR_AWS_SECRET_ACCESS_KEY

# 設定後、このプロファイルに切り替える
llm-cli profile use bedrock-claude # または llm-cli profile use bedrock
```

**認証情報の優先順位:**

1.  `llm-cli` プロファイルに直接設定された `aws_access_key_id` と `aws_secret_access_key`
2.  AWS SDKのデフォルトの認証情報プロバイダチェーン（環境変数 `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`、IAMロールなど）

#### 必要なIAMポリシー

Amazon Bedrockのモデルを呼び出すには、AWSの認証情報に適切なIAMポリシーが付与されている必要があります。最小限必要なアクションは `bedrock:InvokeModel` および `bedrock:InvokeModelWithResponseStream` です。

**最小限のIAMポリシーの例:**

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
            "Resource": "arn:aws:bedrock:ap-northeast-1::foundation-model/amazon.nova*"
        }
    ]
}
```

**注意**: `<your-aws-region>` と `<your-model-id>` は、実際に使用するリージョンとモデルIDに置き換えてください。セキュリティのベストプラクティスとして、`Resource` は可能な限り具体的なモデルに限定することを強く推奨します。複数のモデルを使用する場合は、`"Resource": "arn:aws:bedrock:<your-aws-region>::/foundation-model/*"` のようにワイルドカードを使用することもできますが、その場合はアクセス権が広がることに注意してください。

## 設定

設定は `~/.config/llm-cli/config.json` に保存されます。`profile` コマンド群で管理できますが、`profile edit` で直接編集することも可能です。

**セキュリティに関する注意**: APIキーなどの機密情報は、設定ファイルに平文で保存されます。このファイルへのアクセスは、ご自身の責任で管理してください。

## 謝辞

このプロジェクトは、GoogleのAIアシスタント「Gemini」をコーディングパートナーとして開発されました。

## ライセンス

このプロジェクトはMITライセンスです。詳細は`LICENSE`ファイルをご覧ください。