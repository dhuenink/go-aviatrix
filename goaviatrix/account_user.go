package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AccountUser contains the elements necessary for creating, getting and deleting user.
// accounts on the controller
// See Also:
//   CreateAccountUser()
//   GetAccountUser()
//   DeleteAccountUser()
type AccountUser struct {
	CID         string `form:"CID,omitempty"`
	Action      string `form:"action,omitempty"`
	UserName    string `form:"username,omitempty" json:"user_name,omitempty"`
	AccountName string `form:"account_name,omitempty" json:"acct_names,omitempty"`
	Email       string `form:"email,omitempty" json:"user_email,omitempty"`
	Password    string `form:"password,omitempty" json:"password,omitempty"`
}

// AccountUserEdit contains the elements necessary for updating or editing as existing account.
// See Also:
//   UpdateAccountUserObject()
type AccountUserEdit struct {
	CID         string `form:"CID,omitempty"`
	Action      string `form:"action,omitempty"`
	UserName    string `form:"username,omitempty" json:"user_name,omitempty"`
	AccountName string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Email       string `form:"email,omitempty" json:"email,omitempty"`
	What        string `form:"what,omitempty" json:"what,omitempty"`
	OldPassword string `form:"old_password,omitempty" json:"old_password,omitempty"`
	NewPassword string `form:"new_password,omitempty" json:"new_password,omitempty"`
}

// AccountUserListResp contains the http response object returned from listing the user accounts.
// See Also:
// ListAccountUsers()
type AccountUserListResp struct {
	Return          bool          `json:"return"`
	AccountUserList []AccountUser `json:"results"`
	Reason          string        `json:"reason"`
}

// CreateAccountUser does an http POST request to add a new user account to the controller.
//
// Required Arguments:
//   user *AccountUser
// Returns:
//   error if any
func (c *Client) CreateAccountUser(user *AccountUser) error {
	user.CID = c.CID
	user.Action = "add_account_user"
	resp, err := c.Post(c.baseURL, user)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

// GetAccountUser does an http GET request to retrieve a specific user account from the controller and
// decodes the json response into the AccountUserListResponse struct then it loops through the data and
// returns the specific user requested or an error if the user was not found.
//
// Required Arguments:
//   user *AccountUser
// Returns:
//   *AccountUser
//   error if any
func (c *Client) GetAccountUser(user *AccountUser) (*AccountUser, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_account_users", c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data AccountUserListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if !data.Return {
		return nil, errors.New(data.Reason)
	}
	users := data.AccountUserList
	for i := range users {
		if users[i].UserName == user.UserName && users[i].AccountName == user.AccountName {
			debug("[INFO] Found Aviatrix user account %s", user.UserName)
			return &users[i], nil
		}
	}
	debug("Couldn't find Aviatrix user account %s", user.UserName)
	return nil, ErrNotFound

}

// UpdateAccountUserObject does an http POST request with the user data in the AccountUserEdit struct and updates
// the user account in the controller, it then checks the response and will return an error if the update was not
// successful.
//
// Required Arguments:
//   user *AccountUserEdit struct
// Returns:
//   error if any
func (c *Client) UpdateAccountUserObject(user *AccountUserEdit) error {
	user.CID = c.CID
	user.Action = "edit_account_user"
	resp, err := c.Post(c.baseURL, user)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

// DeleteAccountUser does an http GET request given a user account and deletes it from the conroller, then checks
// the response and will return an error if it was unsuccessful.
//
// Required Arguments:
//    user *AccountUser struct
// Returns:
//    error if any
func (c *Client) DeleteAccountUser(user *AccountUser) error {
	path := c.baseURL + fmt.Sprintf("?action=delete_account_user&CID=%s&username=%s", c.CID, user.UserName)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

// ListAccountUsers does an http GET request to retrieve all the user accounts from the controller and
// decodes the json response into the AccountUserListResponse struct and returns the list of users or
// an error if the request was not successful.
//
// Required Arguments:
//   None
// Returns:
//   *[]AccountUser
//   error if any
func (c *Client) ListAccountUsers() (*[]AccountUser, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_account_users", c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data AccountUserListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if !data.Return {
		return nil, errors.New(data.Reason)
	}

	users := data.AccountUserList
	return &users, nil

}
