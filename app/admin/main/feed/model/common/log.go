package common

const (
	//BusinessID action log business ID
	BusinessID = 204
	//LogPopularStars popular new start card log
	LogPopularStars = 0
	//LogChannelTab channel tab log
	LogChannelTab = 1
	//LogEventTopic popular event topic log
	LogEventTopic = 2
	//LogSWEBCard search web card log
	LogSWEBCard = 3
	//LogSWEB search web log
	LogSWEB = 4
)

//LogManager .
type LogManager struct {
	ID        int    `json:"id"`
	OID       int    `json:"oid"`
	Uname     string `json:"uname"`
	UID       int    `json:"uid"`
	Type      int    `json:"module"`
	ExtraData string `json:"content"`
	Action    string `json:"action"`
	CTime     string `json:"ctime"`
}

//LogSearch .
type LogSearch struct {
	ID        int    `json:"id"`
	OID       int    `json:"oid"`
	Uname     string `json:"uname"`
	UID       int    `json:"uid"`
	Type      int    `json:"type"`
	ExtraData string `json:"extra_data"`
	Action    string `json:"action"`
	CTime     string `json:"ctime"`
}

//ManagerPage .
type ManagerPage struct {
	CurrentPage int `json:"current_page"`
	TotalItems  int `json:"total_items"`
	PageSize    int `json:"page_size"`
}

//LogManagers .
type LogManagers struct {
	Item []*LogManager `json:"item"`
	Page ManagerPage   `json:"pager"`
}
