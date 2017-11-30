package gojsonrpc

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client executes JSON RPC calls to remote servers.
type Client struct {
	URL string

	timeout    int
	proxyURL   string
	httpClient *http.Client
}

// NewClient returns a newly istantiated Client pointing to the given url.
func NewClient(url string) *Client {
	client := &Client{
		URL:      url,
		timeout:  defaultTimeout,
		proxyURL: "",
	}
	client.setHTTPClient()

	return client
}

// SetTimeout sets the client timeout to the given value.
func (c *Client) SetTimeout(timeout int) {
	c.timeout = timeout
	c.setHTTPClient()
}

// SetHTTPProxy tells the client to use the given httpProxyURL as proxy address.
func (c *Client) SetHTTPProxy(httpProxyURL string) {
	c.proxyURL = httpProxyURL
	c.setHTTPClient()
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

func (c *Client) setHTTPClient() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy: http.ProxyFromEnvironment,
	}

	if parsedProxyURL, err := url.Parse(c.proxyURL); c.proxyURL != "" && err == nil {
		transport.Proxy = http.ProxyURL(parsedProxyURL)
	}

	newHTTPClient := &http.Client{
		Timeout:   time.Duration(time.Duration(c.timeout) * time.Second),
		Transport: transport,
	}

	c.httpClient = newHTTPClient
}

func (c *Client) sendJSONRequest(jsonRequest []byte) ([]byte, error) {
	var jsonResponse []byte

	httpRequest, err := http.NewRequest("POST", c.URL, strings.NewReader(string(jsonRequest)))
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Content-Length", "")
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Connection", "close")

	httpResponse, err := c.httpClient.Do(httpRequest)
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
