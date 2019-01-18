package goaviatrix

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	loginFailed = `{
		"return": false,
		"reason": "User name\/password does not match"
	  
	  }`
)

var server *httptest.Server

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	server = httptest.NewTLSServer(handler)

	c := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return c, server.Close
}

func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestNewClientSuccess(t *testing.T) {
	tf := "loginRespSuccess.json"
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		r.ParseForm()
		assert.Equal(t, "login", r.Form.Get("action"))
		assert.Equal(t, "testing123!", r.Form.Get("password"))
		assert.Equal(t, "testuser", r.Form.Get("username"))
		//assert.Equal(t, "secret", r.Header.Get("Secret"))
		w.Write([]byte(fixture(tf)))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	client, err := NewClient("testuser", "testing123!", "localhost", SetHTTPClient(httpClient), BaseURL(server.URL+"/v1/api"))
	if err != nil {
		fmt.Println("unable to create client")
	}
	assert.Nil(t, err)
	assert.Equal(t, "57e098ed708a8", client.CID)

}

func TestNewClientBadUserOrPassword(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		r.ParseForm()
		assert.Equal(t, "login", r.Form.Get("action"))
		assert.Equal(t, "testing123!", r.Form.Get("password"))
		assert.Equal(t, "testuser", r.Form.Get("username"))
		w.Write([]byte(loginFailed))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	_, err := NewClient("testuser", "testing123!", "127.0.0.1", SetHTTPClient(httpClient), BaseURL(server.URL+"/v1/api"))
	assert.NotNil(t, err)
}

func TestNewClientNoControllerIP(t *testing.T) {
	_, err := NewClient("testuser", "testing123", "")
	assert.NotNil(t, err)
}
func TestNewClientInvalidControllerIP(t *testing.T) {
	_, err := NewClient("testuser", "testing123", "10.256.0.1")
	assert.NotNil(t, err)
}

func TestNewClientInvalidControllerHostName(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping InvalidControllerIP in short mode")
	}
	_, err := NewClient("testuser", "testing123", "xserfsd")
	assert.NotNil(t, err)
}
