package gojsonrpc

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client executes JSON RPC calls to remote servers.
type Client struct {
	serverURL  string
	httpClient *http.Client
}

// NewClient returns a newly istantiated Client pointing to the given url.
func NewClient(url string) *Client {
	httpClient := &http.Client{
		Timeout: time.Duration(time.Duration(defaultTimeout) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyFromEnvironment,
		},
	}

	return &Client{
		serverURL:  url,
		httpClient: httpClient,
	}
}

// SetTimeout sets the client timeout to the given value.
func (c *Client) SetTimeout(timeout int) {
	c.httpClient.Timeout = time.Duration(timeout) * time.Second
}

// SetHTTPProxyURL tells the client to use the given proxyURL as proxy address.
func (c *Client) SetHTTPProxyURL(proxyURL *url.URL) {
	c.httpClient.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
}

// RunOptions represents options that can be used to configure a Run
// operation.
type RunOptions struct {
	AdditionalHeaders map[string]string
}

// Run executes the given method having the given params setting the response
// value in the given result interface.
func (c *Client) Run(method string, params interface{}, result interface{}, opts ...RunOptions) error {
	request := NewRequest(method, params, RandInt(10000000, 99999999))

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}

	httpResponse, err := c.sendJSONRequest(jsonRequest, opts...)
	if err != nil {
		return err
	}

	jsonRPCResponse := NewResponse()
	jsonRPCResponse.Result = result

	err = json.Unmarshal(httpResponse, &jsonRPCResponse)
	if err != nil {
		return err
	}

	if jsonRPCResponse.hasError() {
		return errors.New(jsonRPCResponse.Error.Message)
	}

	return nil
}

// NotifyOptions represents options that can be used to configure a Notify
// operation.
type NotifyOptions = RunOptions

// Notify executes the given method with the given parameters.
// Doesn't expect any result.
func (c *Client) Notify(method string, params interface{}, opts ...NotifyOptions) error {
	request := NewRequest(method, params, 0)

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = c.sendJSONRequest(jsonRequest, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) sendJSONRequest(jsonRequest []byte, opts ...RunOptions) ([]byte, error) {
	httpRequest, err := http.NewRequest("POST", c.serverURL, strings.NewReader(string(jsonRequest)))
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Content-Length", "")
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Connection", "close")

	// Apply additional headers
	for _, o := range opts {
		for key, value := range o.AdditionalHeaders {
			httpRequest.Header.Set(key, value)
		}
	}

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()

	httpResponseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return httpResponseBody, err
	}

	if httpResponse.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("received HTTP status %d with body %s", httpResponse.StatusCode, httpResponseBody)
	}

	return httpResponseBody, nil
}
