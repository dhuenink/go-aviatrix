package goaviatrix

import (
	"encoding/json"
	"fmt"
)

// AdminEmailRequest contains the data to set the admin Email
// see also:
//   SetAdminEmail()
type AdminEmailRequest struct {
	APIRequest
	Email string `form:"admin_email" url:"admin_email"`
}

// LoginProcRequest contains the data to get the Admin Email
// see Also:
//   GetAdminEmail()
type LoginProcRequest struct {
	Action   string `form:"action" url:"action"`
	Username string `form:"username" url:"username"`
	Password string `form:"password" url:"password"`
}

// AdminEmailResponse contains the results of the SetAdminEmail Function
// see also:
//   SetAdminEmail()
type AdminEmailResponse struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

// LoginProcResponse contains the results of the GetAdminEmail Function
// see Also:
//   GetAdminEmail()
type LoginProcResponse struct {
	AdminEmail   string `json:"admin_email"`
	InitialSetup bool   `json:"initial_setup"`
}

// SetAdminEmail sets the admin eamil address in the controller and checks the response for errors
// Arguments:
//   adminEmail string "test@test.com"
// Returns:
//   error if any
func (c *Client) SetAdminEmail(adminEmail string) error {
	debug("[TRACE] Setting admin email to '%s'", adminEmail)
	admin := new(AdminEmailRequest)
	admin.Email = adminEmail
	admin.Action = "add_admin_email_addr"
	admin.CID = c.CID
	_, _, err := c.Do("GET", admin)
	if err != nil {
		return err
	}

	return nil
}

// GetAdminEmail logins into the controller using the sepcified username and password and returns the
// admin email address that is currently set
// Arguments:
//   username string
//   password string
// Returns:
//   string containing the admin email address
//   error if any
func (c *Client) GetAdminEmail(username string, password string) (string, error) {
	debug("[TRACE] Getting admin email")
	path := fmt.Sprintf("https://%s/v1/backend1", c.ControllerIP)
	admin := new(LoginProcRequest)
	admin.Action = "login_proc"
	admin.Username = username
	admin.Password = password
	resp, err := c.Post(path, admin)
	if err != nil {
		return "", err
	}
	var data LoginProcResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.AdminEmail, nil
}
