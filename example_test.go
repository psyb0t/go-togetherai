package together_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	together "github.com/psyb0t/go-togetherai"
)

// Example is the README usage example, pointed at a local test server
// instead of https://api.together.ai/v1.
func Example() {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `[{"id":"org/model-a",`+
				`"object":"model","created":1,"type":"chat"}]`)
		},
	))
	defer srv.Close()

	auth := together.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization",
				"Bearer "+os.Getenv("TOGETHER_API_KEY"))

			return nil
		},
	)

	client, err := together.NewClientWithResponses(srv.URL+"/v1", auth)
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

	// Output: org/model-a
}
