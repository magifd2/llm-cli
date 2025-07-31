# llm-cli

`llm-cli` は、ローカル（Ollama, LM Studio）またはリモートのLLM（将来的にはOpenAIなど）と、コマンドラインから直接対話するためのCLIツールです。

## 特徴

*   **マルチプロバイダー対応**: Ollama, LM Studio (OpenAI互換API) に対応。
*   **プロファイル管理**: 複数のLLM設定（エンドポイント、モデルなど）をプロファイルとして保存し、簡単に切り替え可能。
*   **柔軟な入力**: コマンドライン引数、ファイル、標準入力（パイプ）からプロンプトを渡せます。
*   **ストリーミング表示**: LLMからの応答をリアルタイムで表示します。
*   **Goによるシングルバイナリ**: 設定ファイル以外は単一の実行ファイルで動作し、簡単に配布できます。

## 使い方

### プロンプトの送信

```bash
# シンプルなプロンプト
llm-cli ask --prompt "日本の首都はどこですか？"

# システムプロンプト付き
llm-cli ask --prompt "自己紹介して" --system-prompt "あなたは猫です。語尾にニャンを付けて話してください。"

# ストリーミング表示
llm-cli ask --prompt "1から100まで数えてください" --stream

# ファイルからプロンプトを読み込む
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

## 設定

設定は `~/.config/llm-cli/config.json` に保存されます。`profile` コマンド群で管理できますが、`profile edit` で直接編集することも可能です。

**セキュリティに関する注意**: APIキーなどの機密情報は、設定ファイルに平文で保存されます。このファイルへのアクセスは、ご自身の責任で管理してください。

## 謝辞

このプロジェクトは、GoogleのAIアシスタント「Gemini」をコーディングパートナーとして開発されました。

## ライセンス

このプロジェクトはMITライセンスです。詳細は`LICENSE`ファイルをご覧ください。
