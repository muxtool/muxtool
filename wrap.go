package muxtool

import (
	"context"
	"encoding/json"
	"net/http"
)

// Wrap API handler and returns standard http.HandlerFunc function
func Wrap[RQ any, RS any](handler func(ctx context.Context, request *RQ) (RS, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(RQ)
		richifyRequest(req, r)
		switch r.Method {
		case http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodPut:
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}
		resp, err := handler(r.Context(), req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		statusCode := http.StatusOK
		contentType := "application/json"
		var body []byte

		if v, ok := (any)(resp).(WithContentType); ok {
			contentType = v.ContentType()
		}
		if v, ok := (any)(resp).(WithHTTPStatus); ok {
			statusCode = v.Status()
		}
		if v, ok := (any)(resp).(Renderer); ok {
			body, err = v.Render()
		} else {
			body, err = json.Marshal(resp)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", contentType)
		w.Write(body)
	}
}

func richifyRequest[RQ any](req *RQ, baseRequest *http.Request) {
	if v, ok := (any)(req).(WithHeader); ok {
		v.WithHeader(baseRequest.Header)
	}
	if v, ok := (any)(req).(WithMethod); ok {
		v.WithMethod(baseRequest.Method)
	}
}

type NilRequest struct{}

// Optional interfaces for request type

// WithHeader sets headers to request
type WithHeader interface {
	WithHeader(header http.Header)
}

// WithMethod sets method to request
type WithMethod interface {
	WithMethod(method string)
}

// Optional interfaces for response type

// Renderer renders response to byte slice
type Renderer interface {
	Render() ([]byte, error)
}

// WithContentType returns custom content type for response
type WithContentType interface {
	ContentType() string
}

// WithHTTPStatus returns custom status code
type WithHTTPStatus interface {
	Status() int
}
