package model

// ClientStatusResp response clientstatus request just for debug
type ClientStatusResp struct {
	QueueLen int             `json:"queue_len"`
	Clients  []*ClientStatus `json:"clients"`
}

// ClientStatus client status
type ClientStatus struct {
	Addr     string `json:"addr"`
	UpTime   int64  `json:"up_time"`
	ErrCount int64  `json:"err_count"`
	Rate     int64  `json:"rate"`
}
