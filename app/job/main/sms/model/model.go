package model

import (
	"container/ring"
	"context"
	"encoding/json"
	"sync"

	smsmdl "go-common/app/service/main/sms/model"
)

const (
	// SmsPrefix .
	SmsPrefix = "【哔哩哔哩】"
	// SmsSuffix .
	SmsSuffix = " 回TD退订"
	// SmsSuffixChuangLan .
	SmsSuffixChuangLan = " 退订回T"
)

// Provider service provider
type Provider interface {
	// SendSms send sms
	GetPid() int32
	// SendSms send sms
	SendSms(context.Context, *smsmdl.ModelSend) (string, error)
	// SendActSms send act sms
	SendActSms(context.Context, *smsmdl.ModelSend) (string, error)
	// SendBatchActSms send batch act sms
	SendBatchActSms(context.Context, *smsmdl.ModelSend) (string, error)
	// SendInternationalSms send international sms
	SendInternationalSms(context.Context, *smsmdl.ModelSend) (string, error)
}

// Message .
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// UserMobile .
type UserMobile struct {
	CountryCode string `json:"code"`
	Mobile      string `json:"tel"`
}

// ConcurrentRing thread-safe ring
type ConcurrentRing struct {
	*ring.Ring
	sync.Mutex
}

// NewConcurrentRing .
func NewConcurrentRing(length int) *ConcurrentRing {
	return &ConcurrentRing{Ring: ring.New(length)}
}
