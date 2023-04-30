package global

import (
	"chatgpt-web-go/initialize/config"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Gdb       *gorm.DB
	Gredis    *redis.Client
	Cfg       *config.Config
	Gzap      *zap.Logger
	GPTClient *openai.Client
	Validate  *validator.Validate
)
