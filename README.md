# go-togetherai

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-togetherai.svg)](https://pkg.go.dev/github.com/psyb0t/go-togetherai)
[![CI](https://github.com/psyb0t/go-togetherai/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-togetherai/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-togetherai/badges/version.svg)](https://github.com/psyb0t/go-togetherai/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-togetherai/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-togetherai/badges/importers.svg)](https://github.com/psyb0t/go-togetherai/blob/badges/importers.md)

Unofficial Go client for [Together AI](https://docs.together.ai), generated from the OpenAPI spec Together AI publishes.

| Package | API | Generator |
|---|---|---|
| `together` | REST, `https://api.together.ai/v1` | [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) |

`client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-togetherai
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. This lists the available models:

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	together "github.com/psyb0t/go-togetherai"
)

func main() {
	auth := together.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+os.Getenv("TOGETHER_API_KEY"))

		return nil
	})

	client, err := together.NewClientWithResponses("https://api.together.ai/v1", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.ModelsWithResponse(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	for _, model := range *res.JSON200 {
		fmt.Println(model.Id)
	}
}
```

`example_test.go` runs the same code against a local test server, so the example keeps compiling as the spec changes.

The full API, with the spec's descriptions as doc comments, is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-togetherai) and in `go doc github.com/psyb0t/go-togetherai` from inside your module.

## Auth

Together AI takes a bearer API key, created in the API keys section of your account settings. The client sends no credentials by itself; add the `Authorization` header with a `RequestEditorFn` as above, either for every request through `WithRequestEditorFn` or for one call by passing it as the last argument of that method.

## Agent integrations

The [skill](.agents/skills/go-togetherai) teaches a coding agent how to use go-togetherai when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-togetherai@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-togetherai@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-togetherai
```

## Regenerating

Together AI's spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `openapi.yaml` from `https://docs.together.ai/openapi.yaml`.
- `go tool patchspec` (`internal/cmd/patchspec`) writes `api/openapi.patched.yaml` with two fixes oapi-codegen needs. It drops discriminators that have no `mapping`, since their oneOf members are inline and oapi-codegen can't map them. It also drops the inline enum from array items inside a oneOf or anyOf, which oapi-codegen otherwise names the same as the array type. Those unions still decode through their `As*` and `From*` methods.
- oapi-codegen generates `client.gen.go` from the patched spec with `config.yaml`.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Together AI's published API specification, which is MIT licensed. The specification itself is not redistributed here.
