package model

import (
	"time"
)

// OfficialStream 正式流
type OfficialStream struct {
	ID                  int64     `orm:"pk;column(id)" json:"id,omitempty"`
	RoomID              int64     `json:"room_id,omitempty"`
	Src                 int8      `json:"src,omitempty"`
	Name                string    `json:"name,omitempty"`
	Key                 string    `json:"key,omitempty"`
	UpRank              int64     `json:"up_rank,omitempty"`
	DownRank            int64     `json:"down_rank,omitempty"`
	Status              int8      `json:"status,omitempty"`
	LastStatusUpdatedAt time.Time `json:"last_status_updated_at,omitempty"`
	CreateAt            time.Time `json:"create_at,omitempty"`
	UpdateAt            time.Time `json:"update_at,omitempty"`
}
