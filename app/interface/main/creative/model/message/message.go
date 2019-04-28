package message

// Message str
type Message struct {
	ID        int64  `json:"id"`
	TimeStamp int64  `json:"timestamp"`
	TimeAt    string `json:"time_at"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}
