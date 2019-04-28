package sports

import "encoding/json"

// QqRes qq result
type QqRes struct {
	Ret     int             `json:"ret,omitempty"`
	IDlist  json.RawMessage `json:"idlist,omitempty"`
	ID      string          `json:"id"`
	Title   string          `json:"title"`
	URL     string          `json:"url"`
	Kurl    string          `json:"kurl"`
	Atype   interface{}     `json:"atype"`
	Img     json.RawMessage `json:"img"`
	Cid     string          `json:"cid"`
	Src     string          `json:"src"`
	Time    int64           `json:"time"`
	Pubtime string          `json:"pubtime"`
	Surl    string          `json:"surl"`
	Content json.RawMessage `json:"content"`
	Voteid  interface{}     `json:"voteid"`
}
