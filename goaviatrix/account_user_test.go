package goaviatrix

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	acctUser = `{
		"return": true,
		"results": "New user login:api_user1 has been added to account:devtest successfully - Please check email confirmation."
	  }`
)

func TestCreateAccountUser(t *testing.T) {
	tf := "loginRespSuccess.json"
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("action") == "login" {
			w.Write([]byte(fixture(tf)))
		}
		if r.Form.Get("action") == "add_account_user" {
			assert.Equal(t, "57e098ed708a8", r.Form.Get("CID"))
			assert.Equal(t, "api_user1", r.Form.Get("username"))
			assert.Equal(t, "devtest", r.Form.Get("account_name"))
			assert.Equal(t, "api_user1@test.com", r.Form.Get("email"))
			assert.Equal(t, "test123!", r.Form.Get("password"))
			w.Write([]byte(acctUser))
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

	user := AccountUser{
		UserName:    "api_user1",
		AccountName: "devtest",
		Email:       "api_user1@test.com",
		Password:    "test123!",
	}
	if err := client.CreateAccountUser(&user); err != nil {
		fmt.Println("unable to create user account")
	}
	assert.Nil(t, err)

}

func TestGetAccountUser(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		params := r.URL.Query()
		if r.Form.Get("action") == "login" {
			w.Write([]byte(fixture("loginRespSuccess.json")))
		}
		if len(params) != 0 {
			assert.Equal(t, "57e098ed708a8", params["CID"][0])
			assert.Equal(t, "list_account_users", params["action"][0])
			w.Write([]byte(fixture("listAccountUsers.json")))
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

	user := AccountUser{
		UserName:    "user1",
		AccountName: "acct1",
	}
	resp, err := client.GetAccountUser(&user)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	assert.Nil(t, err)
	assert.Equal(t, user.UserName, resp.UserName)

}
