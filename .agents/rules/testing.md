# Testing

Refines `~/.claude/rule-details/go/writing-tests.md`.

- **TST1** Tests never call Together AI. They point the client at an `httptest` server through `newTestClient` in `client_test.go`.
- **TST2** Every test calls `t.Parallel()`. Assert the request method, path, auth header and body, then the decoded `JSONnnn` field.
- **TST3** `example_test.go` holds the README usage example against a local server, with an `// Output:` block, so it compiles and runs. Keep the two in sync.
- **TST4** `make test` runs `go test -race` and must pass.
