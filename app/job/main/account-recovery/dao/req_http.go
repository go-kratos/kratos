package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/job/main/account-recovery/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// CompareInfo compare info
func (d *Dao) CompareInfo(c context.Context, rid int64) (err error) {
	params := url.Values{}
	params.Set("rid", strconv.FormatInt(rid, 10))
	res := new(model.CommonResq)
	if err = d.httpClient.Post(c, d.c.AccRecover.CompareURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("CompareInfo HTTP request err %v", err)
		return
	}
	if res.Code != 0 {
		log.Error("CompareInfo server err_code %d", res.Code)
		err = ecode.Int(int(res.Code))
		return
	}
	return
}

// SendMail send mail
func (d *Dao) SendMail(c context.Context, rid, status int64) (err error) {
	params := url.Values{}
	params.Set("rid", strconv.FormatInt(rid, 10))
	params.Set("status", strconv.FormatInt(status, 10))
	res := new(model.CommonResq)
	if err = d.httpClient.Post(c, d.c.AccRecover.SendMailURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("SendMail HTTP request err %v", err)
		return
	}
	if res.Code != 0 {
		log.Error("SendMail server err_code %d", res.Code)
		err = ecode.Int(int(res.Code))
		return
	}
	return
}
