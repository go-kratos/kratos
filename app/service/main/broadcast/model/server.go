package model

// ServerInfo server info.
type ServerInfo struct {
	Region      string   `json:"region"`
	Server      string   `json:"server"`
	IPCount     int32    `json:"ip_count"`
	ConnCount   int32    `json:"conn_count"`
	RoomIPCount int32    `json:"room_ips"`
	Weight      int32    `json:"weight"`
	Updated     int64    `json:"updated"`
	IPAddrs     []string `json:"ip_addrs"`
	IPAddrsV6   []string `json:"ip_addrs_v6"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Overseas    bool     `json:"overseas"`
}
