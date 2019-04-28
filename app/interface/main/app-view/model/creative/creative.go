package creative

// FollowSwitch get auto follow switch from creative.
type FollowSwitch struct {
	State int8 `json:"state"`
}

// PlayerFollow for player auto follow.
type PlayerFollow struct {
	Show bool `json:"show"`
}

// Points is
type Points struct {
	Type    int    `json:"type"`
	From    int64  `json:"from"`
	To      int64  `json:"to"`
	Content string `json:"content"`
}

// Points is
type Bgm struct {
	Sid     int64  `json:"sid"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	JumpURL string `json:"jump_url"`
}
