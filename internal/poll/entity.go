package poll

type Option struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

type Poll struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Options  []Option `json:"options"`
}
