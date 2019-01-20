package goaviatrix

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccountUser(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("action") == "login" {
			w.Write([]byte(fixture("loginRespSuccess.json")))
		}
		if r.Form.Get("action") == "add_account_user" {
			assert.Equal(t, "57e098ed708a8", r.Form.Get("CID"))
			assert.Equal(t, "api_user1", r.Form.Get("username"))
			assert.Equal(t, "devtest", r.Form.Get("account_name"))
			assert.Equal(t, "api_user1@test.com", r.Form.Get("email"))
			assert.Equal(t, "test123!", r.Form.Get("password"))
			w.Write([]byte(fixture("createAccountUser.json")))
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

func TestUpdateAccountUser(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("action") == "login" {
			w.Write([]byte(fixture("loginRespSuccess.json")))
		}
		if r.Form.Get("action") == "edit_account_user" {
			assert.Equal(t, "57e098ed708a8", r.Form.Get("CID"))
			assert.Equal(t, "api_user1", r.Form.Get("username"))
			assert.Equal(t, "devtest", r.Form.Get("account_name"))
			assert.Equal(t, "account_name", r.Form.Get("what"))
			w.Write([]byte(fixture("updateAccountUser.json")))
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

	user := AccountUserEdit{
		UserName:    "api_user1",
		What:        "account_name",
		AccountName: "devtest",
	}
	if err := client.UpdateAccountUserObject(&user); err != nil {
		fmt.Println("unable to update user account")
	}
	assert.Nil(t, err)

}

func TestDeleteAccountUser(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		params := r.URL.Query()
		if r.Form.Get("action") == "login" {
			w.Write([]byte(fixture("loginRespSuccess.json")))
		}
		if len(params) != 0 {
			assert.Equal(t, "57e098ed708a8", params["CID"][0])
			assert.Equal(t, "delete_account_user", params["action"][0])
			w.Write([]byte(fixture("deleteAccountUser.json")))
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
		UserName: "api_user2",
	}
	if err := client.DeleteAccountUser(&user); err != nil {
		fmt.Printf("error: %v", err)
	}
	assert.Nil(t, err)

}

func TestListAccountUser(t *testing.T) {
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

	users, err := client.ListAccountUsers()
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	assert.Equal(t, "user1", (*users)[0].UserName)
	assert.Equal(t, "acct1", (*users)[0].AccountName)
	assert.Equal(t, "user1@example.com", (*users)[0].Email)
	assert.Equal(t, "user2", (*users)[1].UserName)
	assert.Equal(t, "acct1", (*users)[1].AccountName)
	assert.Equal(t, "user2@example.com", (*users)[1].Email)
}
