package reply

// Member ReplyMember
type Member struct {
	*Info
	FansDetail *FansDetail `json:"fans_detail"`
	Following  int16       `json:"following"` //是否关注
}

// FansDetail FansDetail
type FansDetail struct {
	UID       int64  `json:"uid"`
	MedalID   int32  `json:"medal_id"`      //勋章id
	MedalName string `json:"medal_name"`    //勋章名称
	Score     int32  `json:"score"`         //当前总经验值
	Level     int8   `json:"level"`         //level等级
	Intimacy  int32  `json:"intimacy"`      //当前亲密度
	Status    int8   `json:"master_status"` //佩戴状态1:佩戴中0:未佩戴
	Received  int8   `json:"is_receive"`    //是否领取0:未领取1:已领取
}
