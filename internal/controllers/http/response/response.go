package response

import "backend/config"

type Response struct {
	Code        config.ServiceCode     `json:"code"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Description string                 `json:"description,omitempty"`
}

// NewErrResponse generates the response with the given service code and error string
func NewErr(code config.ServiceCode, err string) Response {
	return Response{
		Code:        code,
		Description: err,
	}
}

// NewResponse generates the response structure with the given service code
// and initializes the data map.
func New(code config.ServiceCode) *Response {
	return &Response{
		Code: code,
		Data: make(map[string]interface{}),
	}
}

// AddKey adds the key to data map of the response with the given value
func (r *Response) AddKey(key string, value interface{}) *Response {
	r.Data[key] = value
	return r
}

// SetDescription sets the value of description field in the response
func (r *Response) SetDescription(value string) *Response {
	r.Description = value
	return r
}
