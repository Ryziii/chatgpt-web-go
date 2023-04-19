package response

var (
	OK       = NewResponse(200, "OK")
	Fail     = NewResponse(500, "Fail")
	NotExist = NewResponse(10001, "不存在")
	NotAuth  = NewResponse(10002, "未授权")
)
