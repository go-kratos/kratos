package model

//Msg for databus consume.
type Msg struct {
	MID         int64 `json:"mid"`
	From        int   `json:"from"`
	IsAuthor    int   `json:"is_author"`
	TimeStamp   int64 `json:"timestamp"`
	ConsumeTime int64 `json:"consume_time,omitempty"`
}
