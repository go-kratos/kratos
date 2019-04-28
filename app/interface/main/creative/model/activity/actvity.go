package activity

const (
	//CancelState 取消活动
	CancelState = -1
	//JoinState 参加活动
	JoinState = 0
)

// Activity for activiy list.
type Activity struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Tags     string `json:"tags"`
	Hot      int8   `json:"hot"`
	ActURL   string `json:"act_url"`
	Protocol string `json:"protocol"`
	Type     int    `json:"type"`
	New      int8   `json:"new"`
	Comment  string `json:"comment"`
	STime    string `json:"stime"`
}

// Like for Like
type Like struct {
	Count int `json:"count"`
}

// Subject for Subject
type Subject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Protocol str
type Protocol struct {
	ID       string `json:"id"`
	Protocol string `json:"protocol"`
	Tags     string `json:"tags"`
	Types    string `json:"types"`
}

// ActWithTP str
type ActWithTP struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Tags     string `json:"tags"`
	Types    string `json:"types"`
	Hot      int8   `json:"hot"`
	ActURL   string `json:"act_url"`
	Protocol string `json:"protocol"`
	Type     int    `json:"type"`
}
