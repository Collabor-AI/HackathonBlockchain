package transports


import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	httptransport "github.com/go-kit/kit/transport/http"
	"HackathonBlockchain/endpoints"
	"net/url"
)


func NewHTTPHandler(endpoints endpoints.Set) http.Handler {
	m := http.NewServeMux()
	m.Handle("/newBlockchain/", httptransport.NewServer(
		endpoints.NewBlockchainEndpoint,
		decodeHTTPNewBlockchainRequest,
		encodeHTTPNewBlockchainResponse,
	))	
	return m
}





func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}





type errorWrapper struct {
	Error string `json:"error"`
}

func decodeHTTPNewBlockchainRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.NewBlockchainRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}


func encodeHTTPNewBlockchainResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
