package response

import "encoding/json"

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}

func (r *Response) WithMessage(msg string) *Response {
	r.Message = msg
	return r
}

func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

func NewResponse(code int, msg string) *Response {
	return &Response{
		Code:    code,
		Data:    nil,
		Message: msg,
		Status:  "Success",
	}
}
func FailResponse(code int, msg string) *Response {
	return &Response{
		Code:    code,
		Data:    nil,
		Message: msg,
		Status:  "Fail",
	}
}

func (r *Response) ToString() string {
	raw, _ := json.Marshal(r)
	return string(raw)
}
