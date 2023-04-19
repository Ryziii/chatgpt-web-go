package setting

type Usage struct {
	ApiModel     string `json:"apiModel"`
	Usage        string `json:"usage"`
	ReverseProxy string `json:"reverseProxy"`
	TimeoutMs    int    `json:"timeoutMs"`
	SocksProxy   string `json:"socksProxy"`
	HttpsProxy   string `json:"httpsProxy"`
}
