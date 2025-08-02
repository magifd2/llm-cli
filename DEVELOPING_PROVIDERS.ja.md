# 開発者ガイド: 新規プロバイダーの追加

このガイドでは、`llm-cli` に新しいLLMプロバイダーのサポートを追加する方法を説明します。

## `Provider` インターフェース

プロバイダーシステムの中心となるのが、`internal/llm/provider.go` に定義されている `Provider` インターフェースです。新しいプロバイダーは、このインターフェースを実装する必要があります。

```go
package llm

import (
	"context"
)

// Provider は、さまざまなLLMと対話するためのインターフェースを定義します。
type Provider interface {
	// Chat は、標準的な（ストリーミングではない）リクエストをLLMに送信します。
	Chat(systemPrompt, userPrompt string) (string, error)

	// ChatStream は、ストリーミングリクエストをLLMに送信します。
	ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error
}
```

### メソッドの詳細

#### `Chat(systemPrompt, userPrompt string) (string, error)`

*   このメソッドは、単純なリクエスト・レスポンスのサイクルを処理します。
*   `systemPrompt`（提供されている場合）と `userPrompt` をLLMのAPIに送信する必要があります。
*   完全なレスポンスを受信するまで処理をブロックしなければなりません。
*   完全な応答テキストを `string` として返す必要があります。
*   エラー（ネットワーク、APIエラーなど）が発生した場合は、`error` を返す必要があります。

#### `ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error`

*   このメソッドは、リアルタイムのストリーミング応答を処理します。
*   プロンプトをLLMのストリーミングAPIエンドポイントに送信します。
*   応答のチャンク（トークン）を受信するたびに、それらを `string` として `responseChan` に送信する必要があります。
*   **重要な規約**: `ChatStream` の実装は、`responseChan` を**決して閉じてはなりません（closeしない）**。チャネルのライフサイクルは、呼び出し元である `cmd/prompt.go` によって管理されます。あなたの実装は、単にデータをチャネルに送信するだけです。
*   いずれかの時点（ストリームの前または最中）でエラーが発生した場合、関数は処理を停止し、`error` を返す必要があります。
*   ユーザーからのキャンセルリクエスト（例: Ctrl+C）を処理するために、`context.Context` を尊重する必要があります。

---

## ステップ・バイ・ステップ実装ガイド

新しいプロバイダーを作成し、統合する方法は以下の通りです。

### ステップ1: プロバイダーファイルの作成

`internal/llm/` ディレクトリに新しいファイルを作成します。例: `internal/llm/my_provider.go`

### ステップ2: インターフェースの実装

新しいファイルで、プロバイダー用の構造体を定義し、2つの必須メソッドを実装します。以下のテンプレートを開始点として使用できます。

```go
package llm

import (
	"context"
	"fmt"

	appconfig "github.com/magifd2/llm-cli/internal/config"
)

// MyProvider は、私たちの新しいサービス用のProviderインターフェースを実装します。
type MyProvider struct {
	Profile appconfig.Profile
}

// Chat は、MyProviderの非ストリーミングリクエストを処理します。
func (p *MyProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	// TODO: プロバイダーのAPIを呼び出すロジックを実装します。
	// 1. プロンプトを使用してリクエストボディを構築します。
	// 2. APIエンドポイント（p.Profile.Endpoint）にHTTPリクエストを送信します。
	// 3. エラーをチェックしながらAPIレスポンスを処理します。
	// 4. レスポンスボディをパースしてメッセージの内容を抽出します。
	// 5. 内容とnilエラーを返します。

	return "", fmt.Errorf("Chat not implemented for MyProvider")
}

// ChatStream は、MyProviderのストリーミングリクエストを処理します。
func (p *MyProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	// TODO: ストリーミングのロジックを実装します。
	// 1. ストリーミング応答用のリクエストを構築します。
	// 2. HTTPリクエストを送信します。
	// 3. ストリームを開始する前にAPIエラーをチェックします。
	// 4. レスポンスボディを行ごと、またはチャンクごとに読み取ります。
	// 5. 各チャンクをパースし、テキストコンテンツを responseChan に送信します。
	// 6. キャンセルのためにコンテキストを尊重します（例: 読み取りループ内）。
	// 7. エラーが発生した場合は、すぐにそれを返します。

	return fmt.Errorf("ChatStream not implemented for MyProvider")
}

```

### ステップ3: プロバイダーの有効化

最後に、CLIが新しいプロバイダーを認識できるようにします。`cmd/prompt.go` を開き、`Run` 関数内の `switch` 文を見つけます。プロバイダー用の新しい `case` を追加してください。

```go
// cmd/prompt.go

// ...
        var provider llm.Provider
        switch activeProfile.Provider {
        case "ollama":
            provider = &llm.OllamaProvider{Profile: activeProfile}
        case "openai":
            provider = &llm.OpenAIProvider{Profile: activeProfile}
        case "bedrock":
            // ... (Bedrockのロジック)

        // 新しいプロバイダーをここに追加
        case "my_provider": // この文字列は設定ファイル内の 'provider' の値と一致する必要があります
            provider = &llm.MyProvider{Profile: activeProfile}

        default:
            fmt.Fprintf(os.Stderr, "警告: プロバイダー '%s' は認識されません...\n", activeProfile.Provider)
            provider = &llm.MockProvider{}
        }
// ...
```

これらの手順の後、ユーザーはプロファイルで `provider: my_provider` と設定することで、あなたの新しい実装を使用できるようになります。

### ステップ4: Vertex AI プロバイダーの実装

Google Cloud Vertex AI プロバイダーは、`internal/llm/vertexai.go` に実装されています。このプロバイダーは、`cloud.google.com/go/vertexai/genai` ライブラリを使用して Vertex AI API と対話します。

**認証:**
認証は、プロファイルで `credentials_file` が指定されている場合はそのサービスアカウントキーを使用し、指定されていない場合は Application Default Credentials (ADC) を使用します。`credentials_file` には `~` (チルダ) を含むパスを指定でき、これは実行時にユーザーのホームディレクトリに展開されます。

**必要なプロファイル設定:**
Vertex AI プロバイダーを使用するには、以下の設定がプロファイルに必要です。

*   `provider`: `vertexai`
*   `model`: 使用する Vertex AI モデルのID (例: `gemini-1.5-pro-001`)
*   `project_id`: GCP プロジェクトID
*   `location`: Vertex AI エンドポイントのリージョン (例: `us-central1`)
*   `credentials_file`: (任意) サービスアカウントキーのJSONファイルへのパス

**実装のポイント:**
*   `newVertexAIClient` 関数内で、`project_id` と `location` の検証が行われます。
*   `credentials_file` が指定されている場合、`option.WithCredentialsFile()` を使用して認証を行います。
*   `Chat` および `ChatStream` メソッドは、`genai.GenerativeModel` を使用してコンテンツを生成します。
*   `SystemInstruction` (単数形) フィールドを使用してシステムプロンプトを設定します。

```
