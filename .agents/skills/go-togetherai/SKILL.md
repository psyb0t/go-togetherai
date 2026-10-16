---
name: go-togetherai
description: Unofficial Go client for the Together AI REST API, generated with oapi-codegen. `together.NewClientWithResponses(baseURL, opts...)` from `github.com/psyb0t/go-togetherai` returns a `*together.ClientWithResponses` with one typed `<Operation>WithResponse` method per endpoint (models, chat completions, completions, embeddings, images, fine-tuning, endpoints and more). Auth is a bearer API key set through `together.WithRequestEditorFn`; responses decode into `JSON200` and per-status `ErrorData` fields. Use when Go code needs to call Together AI or already imports go-togetherai.
homepage: https://github.com/psyb0t/go-togetherai
user-invocable: true
metadata:
  openclaw:
    emoji: "🤝"
    requires:
      bins:
        - go
      env:
        - TOGETHER_API_KEY
permissions:
  filesystem: "Guides edits to the Go source of the current project. The library reads and writes no files."
  shell: "The Go toolchain (`go get`, `go build`, `go test`) to add and compile the dependency."
  network: "HTTPS to the Together AI API at api.together.ai, and `go get` to your Go module proxy."
---

# go-togetherai

Typed Go client for the Together AI API. Every endpoint in Together AI's OpenAPI document is a method on `ClientWithResponses`.

For install, client options, more operations and response handling, see [references/setup.md](references/setup.md).

## Security & safety

- Read the API key from an environment variable such as `TOGETHER_API_KEY`. Never hardcode it, commit it or log it, and never print request headers.
- Requests go only to the base URL you pass, `https://api.together.ai/v1`. Prompts and files you send are processed by Together AI.
- Fine-tuning jobs, dedicated endpoints and clusters cost money, and delete endpoints remove data. Call them only when the user asked for that action.

## When to use

- Go code needs to list models, run chat completions, completions or embeddings, or manage files, fine-tunes or endpoints on Together AI.
- The project already imports `github.com/psyb0t/go-togetherai`.

## When NOT to use

- You want streaming token by token. The generated methods read the whole response body; a streamed reply needs your own SSE reader on top of the raw `Client`.
- The code is not Go.

## Usage

```go
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	together "github.com/psyb0t/go-togetherai"
)

auth := together.WithRequestEditorFn(
	func(_ context.Context, req *http.Request) error {
		req.Header.Set(
			"Authorization",
			"Bearer "+os.Getenv("TOGETHER_API_KEY"),
		)

		return nil
	},
)

client, err := together.NewClientWithResponses(
	"https://api.together.ai/v1",
	auth,
)
if err != nil {
	return err
}

models, err := client.ModelsWithResponse(ctx, nil)
if err != nil {
	return err // transport failure
}

if models.JSON200 != nil {
	for _, model := range *models.JSON200 {
		fmt.Println(model.Id)
	}
}
```

A chat completion. `Messages` is a union type, so the shortest way to build the request is to unmarshal JSON into it:

```go
var req together.ChatCompletionRequest

err = json.Unmarshal([]byte(`{
	"model": "<model-id>",
	"messages": [{"role": "user", "content": "hello"}]
}`), &req)
if err != nil {
	return err
}

res, err := client.ChatCompletionsWithResponse(ctx, req)
if err != nil {
	return err
}

if res.JSON200 == nil {
	if res.JSON400 != nil {
		return fmt.Errorf("together: %s", res.JSON400.Error.Message)
	}

	return fmt.Errorf("together: status %d", res.StatusCode())
}

for _, choice := range res.JSON200.Choices {
	if choice.Message != nil {
		fmt.Println(choice.Message.Content)
	}
}
```

A method returns an error only when the request could not be sent or the body could not be decoded. An API error is a normal response: check `StatusCode()` and read the matching `JSON4xx` or `JSON5xx` field, an `*ErrorData`.
