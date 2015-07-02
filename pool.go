package gojsonrpc

import (
	"errors"
	"strings"
)

type JSONRPCPool struct {
	clients chan (*Client)
}

func NewJSONRPCPool(url string, n int) *JSONRPCPool {
	pool := new(JSONRPCPool)
	pool.clients = make(chan (*Client), n)

	for i := 0; i < n; i++ {
		newClient := NewClient(url)
		pool.clients <- newClient
	}

	return pool
}

func (this *JSONRPCPool) getClient() *Client {
	return <-this.clients
}

func (this *JSONRPCPool) releaseClient(c *Client) {
	this.clients <- c
}

func (this *JSONRPCPool) Do(command, methodName string, params interface{}, result interface{}) error {
	c := this.getClient()
	defer this.releaseClient(c)

	switch strings.ToLower(command) {
	case "run":
		err := c.Run(methodName, params, result)
		if err != nil {
			return err
		}
		break
	case "notify":
		c.Notify(methodName, params)
		break
	default:
		return errors.New("Invalid JSONRPC command")
	}

	return nil
}
