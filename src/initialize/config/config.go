package config

import (
	"time"
)

type Config struct {
	RunMode  string   `mapstructure:"run_mode"`
	Server   Server   `mapstructure:"server"`
	App      App      `mapstructure:"app"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
}
type Server struct {
	HTTPPort     int           `mapstructure:"http_port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type App struct {
	PageSize   int    `mapstructure:"page_size"`
	JwtSecret  string `mapstructure:"jwt_secret"`
	JwtExpires int    `mapstructure:"jwt_expires"`
}

type Database struct {
	Type        string `mapstructure:"type"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Host        string `mapstructure:"host"`
	Name        string `mapstructure:"name"`
	TablePrefix string `mapstructure:"table_prefix"`
}

type Redis struct {
	Host        string        `mapstructure:"host"`
	Password    string        `mapstructure:"password"`
	MaxIdle     int           `mapstructure:"max_idle"`
	MaxActive   int           `mapstructure:"max_active"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout" metadata:"idle_timeout"`
}
