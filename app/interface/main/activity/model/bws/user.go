package bws

// User .
type User struct {
	User         *UserInfo                     `json:"user"`
	Achievements []*UserAchieveDetail          `json:"achievements"`
	Items        map[string][]*UserPointDetail `json:"items"`
}

// UserInfo .
type UserInfo struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
	Key  string `json:"key"`
	Face string `json:"face"`
	Hp   int64  `json:"hp"`
}

// LotteryUser .
type LotteryUser struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
	Face string `json:"face"`
}

// AdminInfo .
type AdminInfo struct {
	IsAdmin bool        `json:"is_admin"`
	Point   interface{} `json:"point"`
}
