package danmu

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/creative/model/danmu"
	dmMdl "go-common/app/interface/main/dm/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_dmReportUpEditURI     = "/x/internal/dm/report/up/edit"
	_dmReportUpListURI     = "/x/internal/dm/report/up/list"
	_dmReportUpArchivesURI = "/x/internal/dm/report/up/archives"
)

// ReportUpList fn
func (d *Dao) ReportUpList(c context.Context, mid, pn, ps int64, aidStr, ip string) (result []*dmMdl.RptSearch, total int64, err error) {
	result = []*dmMdl.RptSearch{}
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", aidStr)
	params.Set("page", strconv.FormatInt(pn, 10))
	params.Set("size", strconv.FormatInt(ps, 10))
	var res struct {
		Code int               `json:"code"`
		Data *dmMdl.RptSearchs `json:"data"`
	}
	if err = d.client.Get(c, d.dmReportUpListURL, ip, params, &res); err != nil {
		err = ecode.CreativeDanmuErr
		log.Error("d.ReportUpList.Get(%s,%s,%s) err(%v)", d.dmReportUpListURL, ip, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(int(res.Code))
		log.Error("d.ReportUpList.Get(%s,%s,%s) err(%v)|code(%d)", d.dmReportUpListURL, ip, params.Encode(), err, res.Code)
		return
	}
	result = res.Data.Result
	total = res.Data.Total
	return
}

// ReportUpArchives fn
func (d *Dao) ReportUpArchives(c context.Context, mid int64, ip string) (ars []*danmu.DmArc, err error) {
	ars = []*danmu.DmArc{}
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int             `json:"code"`
		Data *dmMdl.Archives `json:"data"`
	}
	if err = d.client.Get(c, d.dmReportUpArchivesURL, ip, params, &res); err != nil {
		err = ecode.CreativeDanmuErr
		log.Error("d.ReportUpArchives.Get(%s,%s,%s) err(%v)", d.dmReportUpArchivesURL, ip, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(int(res.Code))
		log.Error("d.ReportUpArchives.Get(%s,%s,%s) err(%v)|code(%d)", d.dmReportUpArchivesURL, ip, params.Encode(), err, res.Code)
		return
	}
	if res.Data != nil {
		for _, v := range res.Data.Result {
			ars = append(ars, &danmu.DmArc{
				Aid:   v.Aid,
				Title: v.Title,
			})
		}
	}
	return
}

// ReportUpEdit fn
func (d *Dao) ReportUpEdit(c context.Context, mid, dmid, cid, op int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("dmid", strconv.FormatInt(dmid, 10))
	params.Set("cid", strconv.FormatInt(cid, 10))
	params.Set("op", strconv.FormatInt(op, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.dmReportUpEditURL, ip, params, &res); err != nil {
		err = ecode.CreativeDanmuErr
		log.Error("d.dmReportUpEditURL.Post(%s,%s,%s) err(%v)", d.dmReportUpEditURL, ip, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.dmReportUpEditURL.Post(%s,%s,%s) err(%v)|code(%d)", d.dmReportUpEditURL, ip, params.Encode(), err, res.Code)
		return
	}
	return
}
