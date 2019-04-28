package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/service/main/spy/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_auditInfo = "/api/member/audit"
)

// AuditInfo get user audit info
func (d *Dao) AuditInfo(c context.Context, mid int64, remoteIP string) (rs *model.AuditInfo, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int             `json:"code"`
		Data model.AuditInfo `json:"data"`
	}
	if err = d.httpClient.Get(c, d.auditInfoURI, remoteIP, params, &res); err != nil {
		log.Error("account get audit info url(%s) error(%v)", d.auditInfoURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("account get audit info url(%s) error(%v)", d.auditInfoURI+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	rs = &res.Data
	log.Info("GET AuditInfoURL suc url(%s) resp(%v)", d.auditInfoURI+"?"+params.Encode(), res)
	return
}
