# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.4.0, 2026-10-28

- Breaking: the type of `GPUClusterInfo.GpuType` changed from `interface{}` to the new `GPUClusterInfoGpuType`. Code that reads or sets the field as an untyped value must use the new type.
- Added the `GPUClusterInfoGpuType` enum with the constants `GPUClusterInfoGpuTypeB200SXM`, `GPUClusterInfoGpuTypeB300SXM`, `GPUClusterInfoGpuTypeGB200SXM`, `GPUClusterInfoGpuTypeH100SXM`, `GPUClusterInfoGpuTypeH100SXMINF`, `GPUClusterInfoGpuTypeH200SXM`, `GPUClusterInfoGpuTypeL40PCIE` and `GPUClusterInfoGpuTypeRTX6000PCI`, and a `Valid` method.
- Added the `GB200_SXM` GPU type: `GPUClusterCreateRequestCommonGpuTypeGB200SXM`, `GPUClusterCreateRequestLegacyNvidiaGpuTypeGB200SXM` and `GPUClusterCreateRequestNvidiaVersionGpuTypeGB200SXM`.
- Added the `ExtendModelResourcesExpiry` operation: `Client.ExtendModelResourcesExpiry`, `Client.ExtendModelResourcesExpiryWithBody`, `ClientWithResponses.ExtendModelResourcesExpiryWithResponse` and `ClientWithResponses.ExtendModelResourcesExpiryWithBodyWithResponse`, plus `NewExtendModelResourcesExpiryRequest`, `NewExtendModelResourcesExpiryRequestWithBody`, `ParseExtendModelResourcesExpiryResponse`, `ExtendModelResourcesExpiryJSONRequestBody`, `ExtendModelResourcesExpiryResponse` and `RLExtendModelResourcesExpiryRequest`.
- Breaking: `ClientInterface` and `ClientWithResponsesInterface` gained the four new `ExtendModelResourcesExpiry` methods, so custom implementations or mocks of them must add the methods.
- Added `InactiveTimeout` to `RLCreateModelResourcesRequest` and `RLModelResources`, and `ExpiresAt` to `RLModelResources`.
- Added `TtlSeconds` to `CreateSecretRequest` and `UpdateSecretRequest`, `TtlSeconds` and `ExpiresAt` to `SecretResponseItem`, and `ValueFromSecretExpiresAt` to `EnvironmentVariable`.

## v0.3.3, 2026-10-28

- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.3.2, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.

## v0.3.1, 2026-10-23

- Added the optional `AppliedCredit` field to `BillingUsageLineItem`. It holds the credit that covers the usage and is nil for uncovered usage.
- Added the `BillingUsageLineItemAppliedCreditType` type with the constants `BillingUsageLineItemAppliedCreditTypeCREDIT`, `BillingUsageLineItemAppliedCreditTypePOSTPAID` and `BillingUsageLineItemAppliedCreditTypePREPAID`, and a `Valid` method.

## v0.3.0, 2026-10-22

- Breaking: the upstream spec dropped the optional `remediation_id` query parameter from the create remediation operation, so `RemediationServiceCreateRemediationParams` was removed.
- Breaking: the `params *RemediationServiceCreateRemediationParams` argument was removed from `Client.RemediationServiceCreateRemediation`, `Client.RemediationServiceCreateRemediationWithBody`, `ClientWithResponses.RemediationServiceCreateRemediationWithResponse` and `ClientWithResponses.RemediationServiceCreateRemediationWithBodyWithResponse`. The `ClientInterface` and `ClientWithResponsesInterface` entries changed the same way.
- Breaking: `NewRemediationServiceCreateRemediationRequest` and `NewRemediationServiceCreateRemediationRequestWithBody` no longer take the params argument.

## v0.2.0, 2026-10-21

- Breaking: the generated client now follows the new upstream spec, which puts the version in each path. Create the client with the base URL `https://api.together.ai` and drop `/v1` or `/v2`. One client now reaches both versions.
- Breaking: methods without a version in their name were removed and replaced by versioned ones. For example `Client.GetFineTunesModelsSupported`, `Client.GetFiles`, `Client.GetBatches` and `Client.GetDeployments` became `GetV1FineTunesModelsSupported`, `GetV1Files`, `GetV1Batches` and `GetV1Deployments`. The `WithResponse` variants and the `ClientInterface` and `ClientWithResponsesInterface` entries changed the same way, as did the matching params types such as `GetFineTunesModelsLimitsParams`.
- Breaking: the upstream spec dropped `GET /v1/models`, so `Client.Models`, `ClientWithResponses.ModelsWithResponse` and `ModelsParams` are gone. `SupportedModelsServiceListSupportedModelsWithResponse` calls `/v2/supported-models` and returns a different shape, with the models in `JSON200.Data`.
- Breaking: `ChatCompletionMessageParam.FromChatCompletionFunctionMessageParam` and `MergeChatCompletionFunctionMessageParam` were removed.
- Breaking: the `FunctionCall` field was removed from the delta in `ChatCompletionChunk.Choices`.
- Added many `V1` operations, for example `GetV1FilesIdContent`, `PostV1FineTunesEstimatePrice`, `PostV1ImagesGenerations` and `PostV1DeploymentsIdRollback`. See the type definitions for the full list.
- The module now requires Go 1.26.10, which fixes `GO-2026-6809` in `crypto/x509`. The development image is pinned to the same version.
- The README, the agent skill and the tests use the new base URL and the supported models call.

## v0.1.3, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.
- The setup reference says that operations under `/organizations`, `/public`, `/supported-models`, `/videos` and some `/projects` paths use the `https://api.together.ai/v2` server, which the generated client does not switch to on its own, so they need a second client with that base URL.

## v0.1.2, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.1, 2026-10-16

- Added the `go-togetherai` agent skill. It installs with `claude plugin install go-togetherai@psyb0t` in Claude Code, `codex plugin add go-togetherai@psyb0t` in Codex and `openclaw skills install @psyb0t/go-togetherai` from ClawHub.

## v0.1.0, 2026-10-16

- Added `together`, a REST client generated with oapi-codegen from Together AI's OpenAPI document (spec version 2.0.0, OpenAPI 3.1.0).
- `make generate` downloads the current spec from docs.together.ai, patches the two shapes oapi-codegen can't read (discriminators without a mapping, inline enums on array items inside a union) and regenerates the client. The spec is not committed.
