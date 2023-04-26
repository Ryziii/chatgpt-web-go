package api

type Session struct {
	Auth  bool   `json:"auth"`
	Model string `json:"model"`
}
