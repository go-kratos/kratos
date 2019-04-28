package model

// TmpInfo old api will del soon
type TmpInfo struct {
	Addr     string `json:"addr"`
	ZoneID   int64  `json:"zoneId"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
}

// Info ipinfo with the smallest zone_id.
type Info struct {
	Addr        string  `json:"addr"`
	Country     string  `json:"country"`
	Province    string  `json:"province"`
	City        string  `json:"city"`
	ISP         string  `json:"isp"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ZoneID      int64   `json:"zoneId"`
	CountryCode int     `json:"country_code"`
}

// InfoComplete ipinfo with all zone_id.
type InfoComplete struct {
	Addr        string  `json:"addr"`
	Country     string  `json:"country"`
	Province    string  `json:"province"`
	City        string  `json:"city"`
	ISP         string  `json:"isp"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ZoneID      []int64 `json:"zone_id"`
	CountryCode int     `json:"country_code"`
}

// IP dont' use this, will del soon. use InfoComplete.
type IP struct {
	Addr     string  `json:"addr"`
	Country  string  `json:"country"`
	Province string  `json:"province"`
	City     string  `json:"city"`
	ISP      string  `json:"isp"`
	ZoneID   []int64 `json:"zone_id"`
}

// Version for check ip library.
type Version struct {
	UpdateTimeV4 string `json:"ipv4_flagship_ipdb_update_time"`
	NewestV4     string `json:""ipv4_flagship_ipdb_newest_url`
	StableV4     string `json:"ipv4_flagship_ipdb_stable_url"`
	UpdateTimeV6 string `json:"ipv6_flagship_ipdb_update_time"`
	NewestV6     string `json:""ipv6_flagship_ipdb_newest_url`
	StableV6     string `json:"ipv6_flagship_ipdb_stable_url"`
}
