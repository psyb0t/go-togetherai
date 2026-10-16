# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-16

- Added `together`, a REST client generated with oapi-codegen from Together AI's OpenAPI document (spec version 2.0.0, OpenAPI 3.1.0).
- `make generate` downloads the current spec from docs.together.ai, patches the two shapes oapi-codegen can't read (discriminators without a mapping, inline enums on array items inside a union) and regenerates the client. The spec is not committed.
