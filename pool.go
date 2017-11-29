package gojsonrpc

import (
	"errors"
	"strings"
)

// Pool executes JSON RPC calls to remote servers using a pool of Clients.
type Pool struct {
	clients chan (*Client)
}

// NewPool returns a newly initialized Pool of clients pointing to the given url.
func NewPool(url string, n int, additionalParameters ...interface{}) *Pool {
	pool := new(Pool)
	pool.clients = make(chan (*Client), n)
	timeout := defaultTimeout

	if len(additionalParameters) == 1 {
		timeoutIntVal, ok := additionalParameters[0].(int)
		if ok {
			timeout = timeoutIntVal
		}
	}

	for i := 0; i < n; i++ {
		newClient := NewClient(url)
		newClient.SetTimeout(timeout)
		pool.clients <- newClient
	}

	return pool
}

func (p *Pool) getClient() *Client {
	return <-p.clients
}

func (p *Pool) releaseClient(c *Client) {
	p.clients <- c
}

// Do executes the given command (Run or Notify) using the given methodName and params
// and building the response in the given result interface.
func (p *Pool) Do(command, methodName string, params interface{}, result interface{}) error {
	c := p.getClient()
	defer p.releaseClient(c)

	switch strings.ToLower(command) {
	case "run":
		err := c.Run(methodName, params, result)
		if err != nil {
			return err
		}
		break
	case "notify":
		err := c.Notify(methodName, params)
		if err != nil {
			return err
		}
		break
	default:
		return errors.New("Invalid JSONRPC command")
	}

	return nil
}
