package goaviatrix

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	adminEmail = `{
		"return": true,
		"results": "admin email address has been successfully added"
	  }`
)

func TestSetAdminEmail(t *testing.T) {
	tf := "loginRespSuccess.json"
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		params := r.URL.Query()
		if r.Form.Get("action") == "login" {
			w.Write([]byte(fixture(tf)))
		}
		if len(params) != 0 {
			assert.Equal(t, "57e098ed708a8", params["CID"][0])
			assert.Equal(t, "add_admin_email_addr", params["action"][0])
			w.Write([]byte(adminEmail))
		}
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	client, err := NewClient("testuser", "testing123!", "localhost", SetHTTPClient(httpClient), BaseURL(server.URL+"/v1/api"))
	if err != nil {
		fmt.Println("unable to create client")
	}
	assert.Nil(t, err)
	assert.Equal(t, "57e098ed708a8", client.CID)

	client.SetAdminEmail("dan.huenink@test.com")

}
