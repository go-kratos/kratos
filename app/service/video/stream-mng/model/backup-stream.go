package model

import (
	"time"
)

// BackupStream 备用流
type BackupStream struct {
	ID             int64     `json:"id,omitempty"`
	RoomID         int64     `json:"room_id,omitempty"`
	StreamName     string    `json:"stream_name,omitempty"`
	Key            string    `json:"key,omitempty"`
	DefaultVendor  int64     `json:"default_vendor,omitempty"`
	OriginUpstream int64     `json:"origin_upstream,omitempty"`
	Streaming      int64     `json:"streaming,omitempty"`
	LastStreamTime time.Time `json:"last_stream_time,omitempty"`
	ExpiresAt      time.Time `json:"expires_at,omitempty"`
	Options        int64     `json:"options,omitempty"`
	Status         int32     `json:"status,omitempty"`
}
