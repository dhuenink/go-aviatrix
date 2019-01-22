package goaviatrix

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("action") == "login" {
			w.Write([]byte(fixture("loginRespSuccess.json")))
		}
		if r.Form.Get("action") == "setup_account_profile" {
			assert.Equal(t, "57e098ed708a8", r.Form.Get("CID"))
			assert.Equal(t, "devtest", r.Form.Get("account_name"))
			assert.Equal(t, "123456789012", r.Form.Get("aws_account_number"))
			assert.Equal(t, "true", r.Form.Get("aws_iam"))
			assert.Equal(t, "arn:aws:iam::123456789012:role/aviatrix-role-app", r.Form.Get("aws_role_arn"))
			assert.Equal(t, "arn:aws:iam::123456789012:role/aviatrix-role-ec2", r.Form.Get("aws_role_ec2"))
			w.Write([]byte(fixture("createAccount.json")))
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

	acct := Account{
		AccountName:      "devtest",
		CloudType:        1,
		AwsAccountNumber: "123456789012",
		AwsIam:           "true",
		AwsRoleApp:       "arn:aws:iam::123456789012:role/aviatrix-role-app",
		AwsRoleEc2:       "arn:aws:iam::123456789012:role/aviatrix-role-ec2",
	}
	if err := client.CreateAccount(&acct); err != nil {
		fmt.Println("unable to create account")
	}
	assert.Nil(t, err)

}

func TestClient_CreateAccount(t *testing.T) {
	type fields struct {
		respFile string
		CID      string
	}
	type args struct {
		account *Account
	}

	acct := Account{
		AccountName:      "devtest",
		CloudType:        1,
		AwsAccountNumber: "123456789012",
		AwsIam:           "true",
		AwsRoleApp:       "arn:aws:iam::123456789012:role/aviatrix-role-app",
		AwsRoleEc2:       "arn:aws:iam::123456789012:role/aviatrix-role-ec2",
	}
	td := args{account: &acct}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Create Account Success",
			fields:  fields{respFile: "createAccount.json", CID: "57e098ed708a8"},
			args:    td,
			wantErr: false,
		},
		{
			name:    "Create Account Fail",
			fields:  fields{respFile: "createAccount1.json", CID: "57e098ed708a8"},
			args:    td,
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.ParseForm()
				if r.Form.Get("action") == "setup_account_profile" {
					assert.Equal(t, "57e098ed708a8", r.Form.Get("CID"))
					assert.Equal(t, "devtest", r.Form.Get("account_name"))
					assert.Equal(t, "123456789012", r.Form.Get("aws_account_number"))
					assert.Equal(t, "true", r.Form.Get("aws_iam"))
					assert.Equal(t, "arn:aws:iam::123456789012:role/aviatrix-role-app", r.Form.Get("aws_role_arn"))
					assert.Equal(t, "arn:aws:iam::123456789012:role/aviatrix-role-ec2", r.Form.Get("aws_role_ec2"))
					w.Write([]byte(fixture(tt.fields.respFile)))
				}
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()
			c := &Client{
				HTTPClient:   httpClient,
				Username:     "testuser",
				Password:     "test123!",
				CID:          tt.fields.CID,
				ControllerIP: "localhost",
				baseURL:      server.URL + "/v1/api",
			}
			if err := c.CreateAccount(tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
