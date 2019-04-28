package jpush

import (
	"encoding/json"
)

// Payload .
type Payload struct {
	Platform     interface{}  `json:"platform"`
	Audience     interface{}  `json:"audience"`
	Notification interface{}  `json:"notification,omitempty"`
	Message      interface{}  `json:"message,omitempty"`
	Options      *Option      `json:"options,omitempty"`
	Callback     *CallbackReq `json:"callback,omitempty"`
}

// NewPayload .
func NewPayload() *Payload {
	return &Payload{
		Options: &Option{},
	}
}

// SetPlatform .
func (p *Payload) SetPlatform(plat *Platform) {
	p.Platform = plat.OS
}

// SetAudience .
func (p *Payload) SetAudience(ad *Audience) {
	p.Audience = ad.Object
}

// SetOptions .
func (p *Payload) SetOptions(o *Option) {
	p.Options = o
}

// SetMessage .
func (p *Payload) SetMessage(m *Message) {
	p.Message = m
}

// SetNotice .
func (p *Payload) SetNotice(notice *Notice) {
	p.Notification = notice
}

// SetCallbackReq .
func (p *Payload) SetCallbackReq(cb *CallbackReq) {
	p.Callback = cb
}

// ToBytes .
func (p *Payload) ToBytes() ([]byte, error) {
	content, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return content, nil
}
