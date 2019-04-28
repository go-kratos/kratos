package jpush

// Payload .
type Payload struct {
	CID          string      `json:"cid,omitempty"`
	Platform     interface{} `json:"platform"`
	Audience     interface{} `json:"audience"`
	Notification interface{} `json:"notification,omitempty"`
	Message      interface{} `json:"message,omitempty"`
	Options      *Option     `json:"options,omitempty"`
}
