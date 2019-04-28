package model

import "time"

// SpecialGift SpecialGift
type SpecialGift struct {
	ID          int64     `json:"id"`
	UID         int64     `json:"uid"`
	GiftID      int64     `json:"gift_id"`
	GiftNum     int64     `json:"gift_num"`
	RoomID      int64     `json:"room_id"`
	CreateTime  time.Time `json:"create_time"`
	CustomField string    `json:"custom_fields"`
	Status      int       `json:"status"`
}
