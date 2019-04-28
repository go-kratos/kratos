package model

type RuleProtocol struct {
	Cond      string          `json:"cond"`
	Key       string          `json:"key"`
	ConfType  string          `json:"type"`
	Max       int64           `json:"max"`
	Min       int64           `json:"min"`
	TopV      int64           `json:"top_v"`
	StringV   string          `json:"string_v"`
	Condition []*RuleProtocol `json:"conditions"`
}

type RecPoolConf struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	ConfType   int     `json:"type"`
	Rules      string  `json:"rules"`
	Priority   int64   `json:"priority"`
	Percent    int64 `json:"percent"`
	ModuleType int64   `json:"module_type"`
	Position   int64   `json:"position"`
}

type RecWhiteList struct {
	RoomId int `json:"room_id"`
}

type RoomData = struct {
	RoomId          int    `json:"room_id"`
	Title           string `json:"title"`
	PopularityCount int    `json:"popularity_count"`
	Keyframe        string `json:"Keyframe"`
	Cover           string `json:"cover"`
	ParentAreaId    int    `json:"parent_area_id"`
	ParenAreaName   string `json:"parent_area_name"`
	AreaId          int    `json:"area_id"`
	AreaName        string `json:"area_name"`
}
