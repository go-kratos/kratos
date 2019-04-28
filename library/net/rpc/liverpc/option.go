package liverpc

import (
	"time"
)

// CallOption ...
type CallOption interface {
	before(*callInfo)
	after(*callInfo)
}

type callInfo struct {
	Header      *Header
	HTTP        *HTTP
	DialTimeout time.Duration
	Timeout     time.Duration
}

// TimeoutOption is timeout for a specific call
type TimeoutOption struct {
	DialTimeout time.Duration
	Timeout     time.Duration
}

func (t TimeoutOption) before(info *callInfo) {
	info.DialTimeout = t.DialTimeout
	info.Timeout = t.Timeout
}

func (t TimeoutOption) after(*callInfo) {
}

// HeaderOption contains Header for liverpc
type HeaderOption struct {
	Header *Header
}

func (h HeaderOption) before(info *callInfo) {
	info.Header = h.Header
}

func (h HeaderOption) after(*callInfo) {
}

// HTTPOption contains HTTP for liverpc
type HTTPOption struct {
	HTTP *HTTP
}

func (h HTTPOption) before(info *callInfo) {
	info.HTTP = h.HTTP
}

func (h HTTPOption) after(*callInfo) {
}
