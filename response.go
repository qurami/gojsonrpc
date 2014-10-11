package jsonrpc

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Jsonrpc string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   ResponseError `json:"error,omitempty"`
	Id      int           `json:"id,omitempty"`
}

func NewResponse() Response {
	return Response{}
}

func (this *Response) hasError() bool {
	if this.Error.Code == 0 && this.Error.Message == "" {
		return false
	}

	return true
}