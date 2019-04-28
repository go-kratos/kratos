package block

// AccountNotify .
type AccountNotify struct {
	UID    int64  `json:"mid"`
	Type   string `json:"type"`
	Action string `json:"action"`
}
