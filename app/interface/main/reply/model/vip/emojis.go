package vip

// Emoji Emoji
type Emoji struct {
	Pid    int64   `json:"pid"`
	Pname  string  `json:"pname"`
	Pstate int8    `json:"pstate"`
	Purl   string  `json:"purl"`
	Emojis []*Face `json:"emojis"`
}

// Face Face
type Face struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	State  int8   `json:"state"`
	Remark string `json:"remark"`
}
