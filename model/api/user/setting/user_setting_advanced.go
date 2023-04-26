package setting

type Advanced struct {
	SystemMessage string `json:"systemMessage"`
	Temperature   string `json:"temperature"`
	TopP          string `json:"top_p"`
}
