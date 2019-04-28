package job

// Job struct
type Job struct {
	ID       int64  `json:"id"`
	JobsCla  string `json:"jobs_classification"`
	Jid      int64  `json:"-"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Slogan   string `json:"slogan,omitempty"`
	CateID   int    `json:"-"`
	AddrID   int    `json:"-"`
	Duty     string `json:"duty"`
	Demand   string `json:"demand"`
}

// Category struct
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type int8   `json:"type"`
}
