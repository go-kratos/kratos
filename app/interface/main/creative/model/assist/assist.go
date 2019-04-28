package assist

import "go-common/library/time"

var (
	// ActEnum action enum
	ActEnum = map[int8]map[int8]string{
		1: {
			1: "删除评论",
			2: "隐藏评论",
		},
		2: {
			1: "删除弹幕",
			2: "屏蔽弹幕",
			3: "保护弹幕",
			4: "拉黑用户",
			5: "移动弹幕到字幕池",
			6: "忽略字幕池的弹幕",
			7: "取消拉黑用户",
		},
	}
)

// Assist is Assists model.
type Assist struct {
	AssistMid    int64                 `json:"assist_mid"`
	Banned       int8                  `json:"banned"`
	AssistAvatar string                `json:"assist_avatar"`
	AssistName   string                `json:"assist_name"`
	Rights       *Rights               `json:"rights"`
	CTime        time.Time             `json:"ctime"`
	MTime        time.Time             `json:"mtime"`
	Total        map[int8]map[int8]int `json:"total"`
}

// AssistLog is single record for assist done
type AssistLog struct {
	ID           int64     `json:"id"`
	Mid          int64     `json:"mid"`
	AssistMid    int64     `json:"assist_mid"`
	AssistAvatar string    `json:"assist_avatar"`
	AssistName   string    `json:"assist_name"`
	Type         int8      `json:"type"`
	Action       int8      `json:"action"`
	SubjectID    int64     `json:"subject_id"`
	ObjectID     string    `json:"object_id"`
	Detail       string    `json:"detail"`
	State        int8      `json:"state"`
	CTime        time.Time `json:"ctime"`
}

// LiveAssist is single record for assist done
type LiveAssist struct {
	AssistMid int64     `json:"uid"`
	RoomID    int64     `json:"roomid"`
	CTime     time.Time `json:"-"`
	Datetime  string    `json:"ctime"`
}

// Rights main and live status
type Rights struct {
	Main int8 `json:"main"`
	Live int8 `json:"live"`
}
