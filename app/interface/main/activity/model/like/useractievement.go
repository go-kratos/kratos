package like

// MissionAward .
type MissionAward struct {
	ID    int64  `json:"id"`
	Award int64  `json:"award"`
	Image string `json:"image"`
	Name  string `json:"name"`
}

// MissionFriends .
type MissionFriends struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
	Face string `json:"face"`
}

// MissionRank .
type MissionRank struct {
	Lid   int64 `json:"lid"`
	Score int64 `json:"score"`
	Rank  int64 `json:"rank"`
}

// MissionLikeAct .
type MissionLikeAct struct {
	Mlid    int64    `json:"mlid"`
	Lottery *Lottery `json:"lottery"`
}

// Lottery .
type Lottery struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Name           string `json:"name"`
		Sponsors       string `json:"sponsors"`
		SponsorsLogo   string `json:"sponsors_logo"`
		GiftID         int64  `json:"gift_id"`
		MessageTitle   string `json:"message_title"`
		MessageContent string `json:"message_content"`
	} `json:"data"`
}

// MissionInfo .
type MissionInfo struct {
	HasHelp int64 `json:"has_help"`
	HasBuff int64 `json:"has_buff"`
}
