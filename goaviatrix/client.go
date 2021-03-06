package goaviatrix

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/ajg/form"
	"github.com/google/go-querystring/query"
)

// LoginResp represents the response object from the `login` action
type LoginResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
	CID     string `json:"CID"`
}

// APIResp represents the basic response from any action
type APIResp struct {
	Return bool   `json:"return"`
	Reason string `json:"reason"`
}

// APIRequest represents the basic fields for any request
type APIRequest struct {
	CID    string `form:"CID,omitempty" json:"CID" url:"CID"`
	Action string `form:"action,omitempty" json:"action" url:"action"`
}

// Client for accessing the Aviatrix Controller
type Client struct {
	HTTPClient   *http.Client
	Username     string
	Password     string
	CID          string
	ControllerIP string
	baseURL      string
}

// Option is a functional option for configuring the API client
type Option func(*Client)

// BaseURL allows overriding of API client baseURL for testing
func BaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// SetHTTPClient allows overriding of the http client for testing
func SetHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// parseOptions parses the supplied options functions and returns a configured
// *Client instance
func (c *Client) parseOptions(opts ...Option) error {
	// Range over each options function and apply it to our API type to
	// configure it. Options functions are applied in order, with any
	// conflicting options overriding earlier calls.
	for _, option := range opts {
		option(c)
	}

	return nil
}

// Login to the Aviatrix controller with the username/password provided in
// the client structure.
// Arguments:
//    None
// Returns:
//    error - if any
func (c *Client) Login() error {
	account := make(map[string]interface{})
	account["action"] = "login"
	account["username"] = c.Username
	account["password"] = c.Password

	resp, err := c.Post(c.baseURL, account)
	if err != nil {
		return err
	}
	var data LoginResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	debug("[TRACE] CID is '%s'.", data.CID)
	c.CID = data.CID
	return nil
}

// NewClient creates a Client object using the arguments provided. Logs
// in to the controller and sets up the http client.
// Required Arguments:
//   username - the controller username
//   password - the controller password
//   controllerIP - the controller IP/host
// Optional Arguments:
//   SetHTTPClient(httpClient *http.Client) - Allows passing in a custom http client
//   BaseURL(baseURL string) - Allows passing in a custom base url
// Returns:
//   Client - the newly created client
//   error - if any
// See Also:
//   init()
func NewClient(username string, password string, controllerIP string, opts ...Option) (*Client, error) {
	if ip := net.ParseIP(controllerIP); ip == nil {
		debug("[INFO] ControllerIP: %s", controllerIP)
		// assuming dns name passed; Check if host is valid
		ip, err := net.LookupIP(controllerIP)
		if err != nil || len(ip) == 0 {
			// IP addres or DNS name are invalid
			return nil, fmt.Errorf("Aviatrix: Client: Controller NotFound, Error: %v", err)
		}
		controllerIP = ip[0].String()
	}

	apiURL := "https://" + controllerIP + "/v1/api"
	client := &Client{
		Username:     username,
		Password:     password,
		ControllerIP: controllerIP,
		baseURL:      apiURL,
	}
	client.parseOptions(opts...)
	if client.HTTPClient == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client.HTTPClient = &http.Client{Transport: tr}
	}
	if err := client.Login(); err != nil {
		return nil, err
	}
	return client, nil
}

// Get issues an HTTP GET request with the given interface form-encoded.
func (c *Client) Get(path string, i interface{}) (*http.Response, error) {
	return c.Request("GET", path, i)
}

// Post issues an HTTP POST request with the given interface form-encoded.
func (c *Client) Post(path string, i interface{}) (*http.Response, error) {
	return c.Request("POST", path, i)
}

// Put issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) Put(path string, i interface{}) (*http.Response, error) {
	return c.Request("PUT", path, i)
}

// Delete issues an HTTP DELETE request.
func (c *Client) Delete(path string, i interface{}) (*http.Response, error) {
	return c.Request("GET", path, i)
}

// Do performs the HTTP request.
// Arguments:
//   verb - GET, PUT, POST, DELETE, etc
//   req  - the query string (for GET) or body for others
// Returns:
//   http.Response - the HTTP response object (body is closed)
//   []byte - the body string as a byte array
//   error - if any
func (c *Client) Do(verb string, req interface{}) (*http.Response, []byte, error) {
	var err error
	var resp *http.Response
	var url string
	var body []byte
	respdata := new(APIResp)

	// do request
	var loop int
	for {
		url = c.baseURL
		loop = loop + 1
		if verb == "GET" {
			// prepare query string
			v, _ := query.Values(req)
			url = url + "?" + v.Encode()
			resp, err = c.Request(verb, url, nil)
		} else {
			resp, err = c.Request(verb, url, req)
		}

		// check response for error
		if err != nil {
			if loop > 2 {
				return resp, nil, err
			} else {
				continue // try again
			}
		}

		debug("[TRACE] %s %s: %d", verb, url, resp.StatusCode)
		// decode the json response and look for errors to retry
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, _ = ioutil.ReadAll(resp.Body)
			if err = json.Unmarshal(body, respdata); err != nil {
				return resp, body, err
			}
			// Check if the CID has expired; if so re-login
			if respdata.Reason == "CID is invalid or expired." && loop < 2 {
				debug("[TRACE] re-login (expired CID)")
				time.Sleep(500 * time.Millisecond)
				if err = c.Login(); err != nil {
					return resp, body, err
				}
				// update the CID value in the object passed
				s := reflect.ValueOf(req).Elem()
				f := s.FieldByName("CID")
				if f.IsValid() && f.CanSet() {
					f.SetString(c.CID)
				}
				// loop around again using new CID
			} else if !respdata.Return {
				return resp, body, errors.New(respdata.Reason)
			} else {
				// Return = True; Reason is not CID expired
				return resp, body, nil
			}
		} else {
			return resp, body, errors.New("Status code")
		}
	}

	return resp, body, err
}

// Request makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) Request(verb string, path string, i interface{}) (*http.Response, error) {
	debug("[TRACE] %s %s", verb, path)
	var req *http.Request
	var err error
	if i != nil {
		buf := new(bytes.Buffer)
		if err = form.NewEncoder(buf).Encode(i); err != nil {
			return nil, err
		}
		body := buf.String()
		debug("[TRACE] %s %s Body: %s", verb, path, body)
		reader := strings.NewReader(body)
		req, err = http.NewRequest(verb, path, reader)
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	} else {
		req, err = http.NewRequest(verb, path, nil)
	}

	if err != nil {
		return nil, err
	}
	return c.HTTPClient.Do(req)
}
