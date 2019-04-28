package model

//AddFreeGift AddFreeGift
type AddFreeGift struct {
	UID      int64  `json:"uid"`
	GiftID   int64  `json:"gift_id"`
	GiftNum  int64  `json:"gift_num"`
	ExpireAt int64  `json:"expire_at"`
	Source   string `json:"source"`
	MsgID    string `json:"msg_id"`
}
