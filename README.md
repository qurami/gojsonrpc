# gojsonrpc

This package provides a client to JSON RPC services.

## Usage example

```go
package main

import (
	"log"
	"github.com/qurami/gojsonrpc"
)

func main() {
	client := gojsonrpc.NewClient("http://mock.rpcservice.url")
	
	// you can optionally set a HTTP proxy for the connection
	proxyURL, _ := "http://proxy.url:3128"
	client.SetHTTPProxy(proxyURL)
	
	// you can also optionally set the connection timeout
	client.SetTimeout(120)

	args := map[string]interface{}{}
	names := make([]string, 0)

	err := client.Run("GetNames", args, &names)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(names)
}
```
