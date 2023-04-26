package gpt

type ApiTypeEnum int

const (
	ApiKey ApiTypeEnum = iota
	AccessToken
)

var apiTypeMessage = [...]string{
	"ApiKey",
	"AccessToken",
}

func (t ApiTypeEnum) String() string {
	return apiTypeMessage[t]
}
