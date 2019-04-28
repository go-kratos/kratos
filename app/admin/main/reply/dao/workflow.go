package dao

import (
	"context"
	"go-common/library/ecode"
	"net/url"
	"strconv"

	"go-common/library/log"
)

const (
	_workflowDel = "http://api.bilibili.co/x/internal/workflow/appeal/v3/delete"
)

// FilterContent get filtered contents by ids.
func (d *Dao) DelReport(c context.Context, oid int64, rpid int64) (err error) {
	var res struct {
		Code int `json:"code"`
	}
	params := url.Values{}
	params.Set("business", "13")
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("eid", strconv.FormatInt(rpid, 10))
	if err = d.httpClient.Post(c, _workflowDel, "", params, &res); err != nil {
		log.Error("DelReport(%s,%s) error(%v)", _workflowDel, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("DelReport(%s,%s) error(%v)", _workflowDel, params.Encode(), err)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("DelReport success!(%s,%s,%d)", _workflowDel, params.Encode())
	return
}
