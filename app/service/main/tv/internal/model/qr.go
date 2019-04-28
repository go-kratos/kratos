package model

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"

	xtime "go-common/library/time"
)

// QR represents pay qr info.
type QR struct {
	ExpireAt xtime.Time
	URL      string
	Token    string
}

// PayParam represents pay params.
type PayParam struct {
	Mid        int64
	Pid        int32
	BuyNum     int32
	Guid       string
	AppChannel string
	Status     int8
	OrderNo    string
	ExpireAt   xtime.Time
}

// MD5 calculates md5 of pay params.
func (p *PayParam) MD5() string {
	m := md5.New()
	io.WriteString(m, strconv.Itoa(int(p.Mid)))
	io.WriteString(m, strconv.Itoa(int(p.Pid)))
	io.WriteString(m, strconv.Itoa(int(p.ExpireAt)))
	return fmt.Sprintf("%x", string(m.Sum(nil)))
}

// IsExpired return true if pay param is expired.
func (p *PayParam) IsExpired() bool {
	return int64(p.ExpireAt) < time.Now().Unix()
}
