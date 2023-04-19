package response

import "encoding/json"

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}

func (r *Response) WithMessage(msg string) *Response {
	r.Code = 200
	r.Status = "Success"
	r.Message = msg
	return r
}

func (r *Response) WithData(data interface{}) *Response {
	r.Code = 200
	r.Status = "Success"
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

func (r *Response) ToString() string {
	raw, _ := json.Marshal(r)
	return string(raw)
}
