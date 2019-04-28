package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/job/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AsoSt aso user status.
type AsoSt struct {
	Code int64
	Data struct {
		Email    string `json:"email"`
		Telphone string `json:"telphone"`
		SafeQs   int8   `json:"safe_question"`
		Spacesta int8   `json:"spacesta"`
	} `json:"data"`
}

const (
	_asoURI        = "http://passport.bilibili.co/intranet/acc/queryByMid"
	_identityURL   = "http://account.bilibili.co/api/identify/info"
	_updateFaceURL = "http://account.bilibili.co/api/internal/member/updateFace"
)

// AsoStatus get aso bind status.
func (d *Dao) AsoStatus(c context.Context, mid int64) (aso *model.MemberAso, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	res := new(struct {
		Code int64           `json:"code"`
		Data model.MemberAso `json:"data"`
	})
	if err = d.client.Get(c, _asoURI, "", params, res); err != nil {
		log.Error("get aso err %v uri %s ", err, _asoURI+params.Encode())
		return
	}
	if res.Code != 0 {
		err = ecode.RequestErr
		log.Error("asotatus err  res %v", res)
		return
	}
	aso = &res.Data
	return
}

// IdentifyStaus identify status.
func (d *Dao) IdentifyStaus(c context.Context, mid int64) (b bool, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int64 `json:"code"`
		Data struct {
			Identify int8 `json:"identify"` // 0
		} `json:"data"`
	}
	if err = d.client.Get(c, _identityURL, "", params, &res); err != nil {
		log.Error("get aso err %v uri %s ", err, _asoURI+params.Encode())
		return
	}
	if res.Code == 0 && res.Data.Identify == 0 {
		b = true
	}
	return
}

// UpdateAccFace update acc face.
func (d *Dao) UpdateAccFace(c context.Context, mid int64, face string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("face", face)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, _updateFaceURL, "", params, &res); err != nil {
		log.Error("update account face err %v uri %s ", err, _updateFaceURL+params.Encode())
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("UpdateAccFace mid(%d) error(%v)", mid, err)
		return
	}
	return
}
