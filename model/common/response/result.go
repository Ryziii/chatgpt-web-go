package response

var (
	OK       = NewResponse(200, "OK")
	Fail     = FailResponse(500, "Fail")
	NotExist = FailResponse(10001, "不存在")
	NotAuth  = FailResponse(10002, "未授权")
)
