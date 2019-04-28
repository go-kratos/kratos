package dao

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/net/http/blademaster/middleware/permit"
)

const (
	_sessionlen   = 32
	_sessionLife  = 2592000
	_dsbCaller    = "manager-go"
	_dsbVerifyURL = "http://dashboard-mng.bilibili.co/api/session/verify"
)

// VerifyDsb .
func (d *Dao) VerifyDsb(ctx context.Context, sid string) (res string, err error) {
	params := url.Values{}
	params.Set("session_id", sid)
	params.Set("encrypt", "md5")
	params.Set("caller", _dsbCaller)
	var dsbRes struct {
		Code     int    `json:"code"`
		UserName string `json:"username"`
	}
	if err = d.dsbClient.Get(ctx, _dsbVerifyURL, "", params, &dsbRes); err != nil {
		return
	}
	if ecode.Int(dsbRes.Code) != ecode.OK {
		err = ecode.Int(dsbRes.Code)
		return
	}
	res = dsbRes.UserName
	return
}

// NewSession .
func (d *Dao) NewSession(ctx context.Context) (res *permit.Session) {
	b := make([]byte, _sessionlen)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return
	}
	res = &permit.Session{
		Sid:    hex.EncodeToString(b),
		Values: make(map[string]interface{}),
	}
	return
}
