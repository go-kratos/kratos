package model

const (

	//State0 视频未审核
	State0 = 0

	//State1 视频安全审核通过
	State1 = 1

	//State2 待冷启动回查
	State2 = 2

	//State3 回查可放出
	State3 = 3

	//State4 视频优质
	State4 = 4

	//State5 视频精选
	State5 = 5
)

//Record4Dup ...
type Record4Dup struct {
	SVID int64  `json:"svid"`
	MID  string `json:"mid"`
	Tag  string `json:"tag"`
}

//Tag ...
type Tag struct {
	TagName string
	TagType int64
	TagID   int64
}
