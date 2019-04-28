package model

import "time"

// MainStream 备用流
type MainStream struct {
	ID             int64     `json:"id,omitempty"`
	RoomID         int64     `json:"room_id,omitempty"`
	StreamName     string    `json:"stream_name,omitempty"`
	Key            string    `json:"key,omitempty"`
	DefaultVendor  int64     `json:"default_vendor,omitempty"`
	OriginUpstream int64     `json:"origin_upstream,omitempty"`
	Streaming      int64     `json:"streaming,omitempty"`
	LastStreamTime time.Time `json:"last_stream_time,omitempty"`
	//第一位预留 第二位是否开启蒙版直播流 第三位wmask蒙版流开播/关播 第四位mmask蒙版流开播/关播
	Options int64 `json:"options,omitempty"`
	Status  int32 `json:"status,omitempty"`
}
