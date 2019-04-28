package offer

const (
	ActionActive = "active"
)

type Retry struct {
	Action string `json:"action,omitempty"`
	Data   *Data  `json:"data,omitempty"`
}

type Data struct {
	OS        string `json:"os,omitempty"`
	IMEI      string `json:"imei,omitempty"`
	Androidid string `json:"androidid,omitempty"`
	Mac       string `json:"mac,omitempty"`
}
