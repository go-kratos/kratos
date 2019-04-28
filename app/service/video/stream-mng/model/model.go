package model

import "go-common/library/time"

// StreamBase
type StreamBase struct {
	StreamName      string  `json:"stream_name,omitempty"`
	DefaultUpStream int64   `json:"default_upstream,omitempty"`
	Origin          int64   `json:"origin,omitempty"`
	Forward         []int64 `json:"forward,omitempty"`
	Type            int     `json:"type,omitempty"`
	Key             string  `json:"-"`
	Options         int64   `json:"options,omitempty"`
	Wmask           bool    `json:"wmask,omitempty"`
	Mmask           bool    `json:"mmask,omitempty"`
}

// StreamFullInfo，
type StreamFullInfo struct {
	RoomID     int64         `json:"room_id,omitempty"`
	Hot        int64         `json:"hot"`
	StreamName string        `json:"stream_name,omitempty"`
	Origin     int64         `json:"origin,omitempty"`
	Forward    []int64       `json:"forward,omitempty"`
	List       []*StreamBase `json:"list,omitempty"`
}

// StreamChangeLog 修改cdnlog
type StreamChangeLog struct {
	RoomID      int64     `json:"room_id,omitempty"`
	FromOrigin  int64     `json:"from_origin,omitempty"`
	ToOrigin    int64     `json:"to_origin,omitempty"`
	Source      string    `json:"source,omitempty"`
	OperateName string    `json:"operate_name,omitempty"`
	Reason      string    `json:"reason,omitempty"`
	CTime       time.Time `json:"ctime,omitempty"`
}

// StreamStatus 流状态
type StreamStatus struct {
	RoomID          int64  `json:"room_id,omitempty"`
	StreamName      string `json:"stream_name,omitempty"`
	DefaultUpStream int64  `json:"default_upstream,omitempty"`
	DefaultChange   bool   `json:"default_change,omitempty"`
	Origin          int64  `json:"origin,omitempty"`
	OriginChange    bool   `json:"origin_change,omitempty"`
	Forward         int64  `json:"forward,omitempty"`
	ForwardChange   bool   `json:"forward_change,omitempty"`
	Key             string `json:"key,omitempty"`
	Add             bool   `json:"add,omitempty"`
	Options         int64  `json:"options,omitempty"`
	OptionsChange   bool   `json:"options_change,omitempty"`
}
