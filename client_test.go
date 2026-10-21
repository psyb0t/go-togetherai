package together

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAPIKey = "test-key"

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testAPIKey)

			return nil
		},
	)

	client, err := NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func TestModels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		params    *SupportedModelsServiceListSupportedModelsParams
		wantQuery string
	}{
		{name: "no params", params: nil, wantQuery: ""},
		{
			name: "search",
			params: &SupportedModelsServiceListSupportedModelsParams{
				Search: new("llama"),
			},
			wantQuery: "search=llama",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := newTestClient(t,
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, "/v2/supported-models", r.URL.Path)
					assert.Equal(t, tt.wantQuery, r.URL.RawQuery)
					assert.Equal(t, "Bearer "+testAPIKey,
						r.Header.Get("Authorization"))

					w.Header().Set("Content-Type", "application/json")
					_, _ = io.WriteString(w,
						`{"object":"list",`+
							`"data":[{"id":"org/model-a"}]}`)
				})

			res, err := client.
				SupportedModelsServiceListSupportedModelsWithResponse(
					context.Background(), tt.params)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode())
			require.NotNil(t, res.JSON200)
			require.Len(t, res.JSON200.Data, 1)
			assert.Equal(t, "org/model-a", res.JSON200.Data[0].Id)
		})
	}
}

func TestChatCompletions(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)
		assert.Equal(t, "Bearer "+testAPIKey, r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "org/model-a", body["model"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"cmpl-1","object":"chat.completion",`+
			`"created":1,"model":"org/model-a","prompt":[],"choices":[`+
			`{"index":0,"message":{"role":"assistant","content":"hi"}}]}`)
	})

	var req ChatCompletionRequest
	require.NoError(t, json.Unmarshal([]byte(`{"model":"org/model-a",`+
		`"messages":[{"role":"user","content":"hello"}]}`), &req))

	res, err := client.ChatCompletionsWithResponse(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	assert.Equal(t, "cmpl-1", res.JSON200.Id)
	require.Len(t, res.JSON200.Choices, 1)
	require.NotNil(t, res.JSON200.Choices[0].Message)
	assert.Equal(t, "hi", res.JSON200.Choices[0].Message.Content)
}

func TestChatCompletionsError(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":{"message":"bad model",`+
			`"type":"invalid_request_error"}}`)
	})

	res, err := client.ChatCompletionsWithResponse(
		context.Background(), ChatCompletionRequest{Model: "nope"})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode())
	require.NotNil(t, res.JSON400)
	assert.Equal(t, "bad model", res.JSON400.Error.Message)
}
