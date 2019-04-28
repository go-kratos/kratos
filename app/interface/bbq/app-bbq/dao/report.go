package dao

import (
	"context"
	"errors"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/log"
	xhttp "net/http"
	"net/url"

	"strconv"
)

//ReportUser rType=0 face;rType=1 name;
func (d *Dao) ReportUser(c context.Context, rType int, mid int64, rmid int64, reason string) (err error) {
	var (
		params url.Values
		req    *xhttp.Request
		res    model.HTTPRpcRes
	)
	params = url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("rmid", strconv.FormatInt(rmid, 10))
	params.Set("type", strconv.Itoa(rType))
	params.Set("reason", reason)
	if req, err = d.httpClient.NewRequest("POST", d.c.URLs["cms_report"], "", params); err != nil {
		log.Errorv(c, log.KV("event", "ReportUser d.httpClient.NewRequest failed"), log.KV("err", err))
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Errorv(c, log.KV("event", "cms ReportUser http req failed"), log.KV("err", err), log.KV("req", req))
		return
	}
	if res.Code != 0 {
		log.Errorv(c, log.KV("event", "cms ReportUser res.code err"), log.KV("err", err))
		err = errors.New("cms ReportUser return err")
		return
	}
	return
}

//ReportDanmu ..
func (d *Dao) ReportDanmu(c context.Context, danmu int64, mid int64, rmid int64, reason string) (err error) {
	var (
		params url.Values
		req    *xhttp.Request
		res    model.HTTPRpcRes
	)
	params = url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("rmid", strconv.FormatInt(rmid, 10))
	params.Set("bullet_id", strconv.FormatInt(danmu, 10))
	params.Set("reason", reason)
	if req, err = d.httpClient.NewRequest("POST", d.c.URLs["cms_report_bullet"], "", params); err != nil {
		log.Errorv(c, log.KV("event", "ReportDanmu d.httpClient.NewRequest failed"), log.KV("err", err))
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Errorv(c, log.KV("event", "cms ReportDanmu http req failed"), log.KV("err", err), log.KV("req", req))
		return
	}
	if res.Code != 0 {
		log.Errorv(c, log.KV("event", "cms ReportDanmu res.code err"), log.KV("err", err))
		err = errors.New("cms ReportDanmu return err")
		return
	}
	return
}

//ReportVideo ..
func (d *Dao) ReportVideo(c context.Context, svid int64, rmid int64, reason string) (err error) {
	var (
		params url.Values
		req    *xhttp.Request
		res    model.HTTPRpcRes
	)
	params = url.Values{}
	params.Set("rmid", strconv.FormatInt(rmid, 10))
	params.Set("svid", strconv.FormatInt(svid, 10))
	params.Set("reason", reason)
	if req, err = d.httpClient.NewRequest("POST", d.c.URLs["cms_report_video"], "", params); err != nil {
		log.Errorv(c, log.KV("event", "ReportVideo d.httpClient.NewRequest failed"), log.KV("err", err))
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Errorv(c, log.KV("event", "cms ReportVideo http req failed"), log.KV("err", err), log.KV("req", req))
		return
	}
	if res.Code != 0 {
		log.Errorv(c, log.KV("event", "cms ReportVideo res.code err"), log.KV("err", err))
		err = errors.New("cms ReportVideo return err")
		return
	}
	return
}
