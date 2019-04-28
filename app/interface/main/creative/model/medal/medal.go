package medal

// Status for medal status
type Status struct {
	Enable  bool  `json:"enable"`
	MedalID int64 `json:"medal_id"`
}

// Medal for fan medal
type Medal struct {
	UID          string `json:"uid"`
	MedalName    string `json:"medal_name"`
	LiveStatus   string `json:"live_status"`
	MasterStatus string `json:"master_status"`
	TimeToChange int64  `json:"time_able_change"`
	RenameStatus int8   `json:"rename_status"`
	Status       string `json:"status"`
	Reason       string `json:"reason"`
	Elec         int64  `json:"charge_num"`
	Coin         int64  `json:"coin_num"`
}

// RecentFans for recent list
type RecentFans struct {
	FansID      int    `json:"fans_id"`
	FansName    string `json:"fans_name"`
	HeadURL     string `json:"head_url"`
	CTime       string `json:"ctime"`
	ReceiveTime string `json:"receive_time"`
}

// FansRank for fans rank
type FansRank struct {
	UID        int64  `json:"uid"`
	Rank       int64  `json:"rank"`
	Score      int64  `json:"score"`
	Level      int64  `json:"level"`
	Uname      string `json:"uname"`
	MedalName  string `json:"medal_name"`
	Special    string `json:"special"`
	MedalColor int64  `json:"medal_color"`
	Face       string `json:"face"`
}
