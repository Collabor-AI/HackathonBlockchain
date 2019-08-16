package transports


import (
	"context"
	"encoding/json"
	"fmt"
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
	m.Handle("/printBlockchain", httptransport.NewServer(
		endpoints.PrintBlockchainEndpoint,
		decodeHTTPPrintBlockchainRequest,
		encodeHTTPPrintBlockchainResponse,
	))
	m.Handle("/AddBlock", httptransport.NewServer(
		endpoints.AddBlockEndpoint,
		decodeHTTPAddBlockRequest,
		encodeHTTPAddBlockResponse,
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

	// respBytes, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}



func decodeHTTPPrintBlockchainRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.PrintBlockchainRequest
	return req, nil
}

func encodeHTTPPrintBlockchainResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		return nil
	}
	fmt.Printf("Print Blockchain: %+v", response)
	respBytes, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(respBytes)
}

func decodeHTTPAddBlockRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AddBlockRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}


func encodeHTTPAddBlockResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		return nil
	}

	// respBytes, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
