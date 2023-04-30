package auth

type Logout struct {
	Token string `json:"token" required:"true"`
}
