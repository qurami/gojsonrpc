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

type Client struct {
	Url     string
	Timeout int
}

func NewClient(url string) *Client {
	client := new(Client)

	client.Url = url
	client.Timeout = DEFAULT_TIMEOUT

	return client
}

func (this *Client) sendJsonRequest(jsonRequest []byte) ([]byte, error) {
	var jsonResponse []byte

	httpClient := &http.Client{
		Timeout: time.Duration(time.Duration(this.Timeout) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	httpRequest, err := http.NewRequest("POST", this.Url, strings.NewReader(string(jsonRequest)))
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

func (this *Client) SetTimeout(timeout int) {
	this.Timeout = timeout
}

func (this *Client) Run(method string, params interface{}, result interface{}) error {
	request := NewRequest(method, params, RandInt(10000000, 99999999))

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}

	jsonResponse, err := this.sendJsonRequest(jsonRequest)
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

func (this *Client) Notify(method string, params interface{}) error {
	request := NewRequest(method, params, 0)

	requestJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = this.sendJsonRequest(requestJson)
	if err != nil {
		return err
	}

	return nil
}
