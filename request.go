package jsonrpc

type Request struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id,omitempty"`
}

func NewRequest(method string, params interface{}, id int) Request {
	request := Request {
		"2.0",
		method,
		params,
		id,
	}

	return request
}