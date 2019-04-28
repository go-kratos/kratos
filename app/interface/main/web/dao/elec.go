package dao

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// ElecShow elec show.
func (d *Dao) ElecShow(c context.Context, mid, aid, loginID int64) (rs *model.ElecShow, err error) {
	var (
		params   = url.Values{}
		remoteIP = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	if loginID > 0 {
		params.Set("login_mid", strconv.FormatInt(loginID, 10))
	}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("act", "appkey")
	var res struct {
		Code int             `json:"code"`
		Data *model.ElecShow `json:"data"`
	}
	if err = d.httpR.Get(c, d.elecShowURL, remoteIP, params, &res); err != nil {
		log.Error("ElecShow url(%s) error(%v)", d.elecShowURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		return
	}
	rs = res.Data
	return
}
