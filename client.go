package gojsonrpc

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client executes JSON RPC calls to remote servers.
type Client struct {
	URL     string
	Timeout int
}

// NewClient returns a newly istantiated Client pointing to the given url.
func NewClient(url string) *Client {
	client := new(Client)

	client.URL = url
	client.Timeout = defaultTimeout

	return client
}

func (c *Client) sendJSONRequest(jsonRequest []byte) ([]byte, error) {
	var jsonResponse []byte

	httpClient := &http.Client{
		Timeout: time.Duration(time.Duration(c.Timeout) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	httpRequest, err := http.NewRequest("POST", c.URL, strings.NewReader(string(jsonRequest)))
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Content-Length", "")
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Connection", "close")

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return jsonResponse, err
	}

	defer httpResponse.Body.Close()

	jsonResponse, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return jsonResponse, err
	}

	return jsonResponse, nil
}

// SetTimeout sets the client timeout to the given value.
func (c *Client) SetTimeout(timeout int) {
	c.Timeout = timeout
}

// Run executes the given method having the given params setting the response
// value in the given result interface.
func (c *Client) Run(method string, params interface{}, result interface{}) error {
	request := NewRequest(method, params, RandInt(10000000, 99999999))

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}

	jsonResponse, err := c.sendJSONRequest(jsonRequest)
	if err != nil {
		return err
	}

	response := NewResponse()
	response.Result = result

	err = json.Unmarshal(jsonResponse, &response)
	if err != nil {
		return err
	}

	if response.hasError() {
		return errors.New(response.Error.Message)
	}

	return nil
}

// Notify executes the given method with the given parameters.
// Doesn't expect any result.
func (c *Client) Notify(method string, params interface{}) error {
	request := NewRequest(method, params, 0)

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = c.sendJSONRequest(jsonRequest)
	if err != nil {
		return err
	}

	return nil
}
