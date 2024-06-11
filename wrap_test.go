package muxtool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleWrap() {
	rr := httptest.NewRecorder()

	// Sample request
	req := reqHello{
		Name: "NeonXP",
	}
	b, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPost, "/hello", bytes.NewReader(b))

	// Handler
	mux := http.NewServeMux()
	// Handle wrapped `handleHello(context.Context, *reqHello) (*respHello, error)`
	mux.Handle("POST /hello", Wrap(handleHello))

	mux.ServeHTTP(rr, request)

	fmt.Println(rr.Body.String())
	// Output: {"message":"Hello, NeonXP!"}
}

type reqHello struct {
	Name string `json:"name"`
}

type respHello struct {
	Message string `json:"message"`
}

func handleHello(ctx context.Context, req *reqHello) (*respHello, error) {
	return &respHello{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}
