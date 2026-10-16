// Package together is a Go client for the Together AI REST API, generated
// with oapi-codegen from the OpenAPI document Together AI publishes.
package together

//go:generate bash scripts/fetch-spec.sh
//go:generate go tool patchspec
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.patched.yaml
