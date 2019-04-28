package offer

type ActiveMsg struct {
	OS        string `json:"os,omitempty"`
	IMEI      string `json:"imei,omitempty"`
	Androidid string `json:"androidid,omitempty"`
	Mac       string `json:"mac,omitempty"`
}
