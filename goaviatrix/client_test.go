package goaviatrix

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
	err    error
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewTLSServer(mux)
	fmt.Println(server.URL)
	return func() {
		server.Close()
	}
}

func TestNewClient(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// ... return the JSON
	})

	client, err = NewClient("testuser", "testing123!", "127.0.0.1", &http.Client{}, BaseURL(server.URL+"/v1/api"))
	if err != nil {
		fmt.Println("unable to create client")
	}
	fmt.Println(client)
}
