package transport

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeJSONRequest[T any](_ context.Context, r *http.Request) (interface{}, error) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func DecodeDefaultRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}
