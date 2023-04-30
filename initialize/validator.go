package initialize

import (
	"chatgpt-web-go/global"
	"github.com/go-playground/validator/v10"
)

func init() {
	global.Validate = validator.New()
}
