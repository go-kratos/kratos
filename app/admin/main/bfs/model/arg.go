package model

// ArgCluster .
type ArgCluster struct {
	Cluster string `form:"cluster" validate:"required"`
}

// ArgAddVolume add volume
type ArgAddVolume struct {
	Group string `form:"group" validate:"required"`
	Num   int64  `form:"num" validate:"required"`
}

// ArgAddFreeVolume add free volume
type ArgAddFreeVolume struct {
	Group string `form:"group" validate:"required"`
	Dir   string `form:"dir" validate:"required"`
	Num   int64  `form:"num" validate:"required"`
}

// ArgCompact group compact
type ArgCompact struct {
	Group string `form:"group" validate:"required"`
	Vid   int64  `form:"vid"`
}

// ArgGroupStatus group status
type ArgGroupStatus struct {
	Group  string `form:"group" validate:"required"`
	Status string `form:"status" validate:"required"`
}

// RespRack .
type RespRack struct {
	Racks map[string]*Rack `json:"racks"`
}

// RespGroup .
type RespGroup struct {
	Groups map[string]*Group `json:"groups"`
}

// RespVolume .
type RespVolume struct {
	Volumes map[string]*VolumeState `json:"volumes"`
}

// RespTotal .
type RespTotal struct {
	Space     int64 `json:"space"`
	FreeSpace int64 `json:"free_space"`
	Groups    int64 `json:"groups"`
	Stores    int64 `json:"stores"`
	Volumes   int64 `json:"volumes"`
}
