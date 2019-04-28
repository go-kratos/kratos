package model

// EventMessage event from bigdata.
type EventMessage struct {
	Time      string      `json:"time"`
	IP        string      `json:"ip"`
	Service   string      `json:"service"`
	Event     string      `json:"event"`
	ActiveMid int64       `json:"active_mid"`
	TargetMid int64       `json:"target_mid"`
	TargetID  int64       `json:"target_id"`
	Args      interface{} `json:"args"`
	Result    string      `json:"result"`
	Effect    string      `json:"effect"`
	RiskLevel int8        `json:"risk_level"`
}

// SpyStatMessage sin stat message.
type SpyStatMessage struct {
	TargetMid int64  `json:"target_mid"`
	TargetID  int64  `json:"target_id"`
	Type      int8   `json:"type"`
	Quantity  int64  `json:"quantity"`
	Time      int64  `json:"time"`
	EventName string `json:"event_name"`
	UUID      string `json:"uuid"`
}
