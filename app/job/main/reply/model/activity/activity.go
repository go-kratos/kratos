package activity

// Activity activity.
type Activity struct {
	Aid   int64  `json:"id"`
	Title string `json:"name"`
	Link  string `json:"pc_url"`
	Cover string `json:"pc_cover"`
}

// Subject subject.
type Subject struct {
	Title string `json:"name"`
	Link  string `json:"act_url"`
}
