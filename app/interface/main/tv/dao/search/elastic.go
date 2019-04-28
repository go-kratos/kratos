package search

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	seaMdl "go-common/app/interface/main/tv/model/search"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
)

func yearTrans(year string) (stime, etime string, err error) {
	yearExp := regexp.MustCompile(`^([\d]{4})-([\d]{4})$`)
	params := yearExp.FindStringSubmatch(year)
	if len(params) < 3 {
		err = ecode.RequestErr
		return
	}
	return params[1] + "-01-01 00:00:00", params[2] + "-12-31 23:59:59", nil
}

// PgcIdx treats the pgc index request and call the ES to get the result
func (d *Dao) PgcIdx(c context.Context, req *seaMdl.ReqPgcIdx) (data *seaMdl.EsPgcResult, err error) {
	var (
		syear, eyear string
		cfg          = d.conf.Cfg.EsIdx.PgcIdx
		r            = d.esClient.NewRequest(cfg.Business).Index(cfg.Index).
				Fields("season_id").WhereEq("is_deleted", 0).WhereEq("status", 0).
				WhereRange("season_id", 0, nil, elastic.RangeScopeLoRc)
	)
	if req.SeasonType > 0 {
		r.WhereEq("season_type", req.SeasonType)
	}
	if req.ProducerID > 0 {
		r.WhereEq("producer_id", req.ProducerID)
	}
	if !req.IsAllStr(req.Year) {
		if syear, eyear, err = yearTrans(req.Year); err != nil {
			log.Warn("PgcIdx Request Year %s, Err %v", req.Year, err)
			return
		}
		r.WhereRange("release_date", syear, eyear, elastic.RangeScopeLcRc)
	}
	if !req.IsAllStr(req.PubDate) {
		if syear, eyear, err = yearTrans(req.PubDate); err != nil {
			log.Warn("PgcIdx Request PubDate %s, Err %v", req.PubDate, err)
			return
		}
		r.WhereRange("pub_time", syear, eyear, elastic.RangeScopeLcRc)
	}
	if req.StyleID > 0 {
		r.WhereEq("style_id", req.StyleID)
	}
	if req.SeasonMonth > 0 {
		r.WhereEq("season_month", req.SeasonMonth)
	}
	if !req.IsAll(req.SeasonStatus) {
		r.WhereIn("pay_status", req.SeasonStatus)
	}
	if !req.IsAll(req.Copyright) {
		r.WhereIn("copyright_info", req.Copyright)
	}
	if !req.IsAllStr(req.IsFinish) {
		isFin, _ := strconv.Atoi(req.IsFinish)
		r.WhereEq("is_finish", isFin)
	}
	if !req.IsAll(req.Area) {
		r.WhereIn("area_id", req.Area)
	}
	if req.SeasonVersion > 0 {
		r.WhereEq("season_version", req.SeasonVersion)
	}
	r.Ps(req.Ps).Pn(int(req.Pn)).Order(req.PgcOrder(), seaMdl.IdxSort(req.Sort))
	if err = r.Scan(c, &data); err != nil {
		log.Error("PgcIdx:Scan params(%s) error(%v)", r.Params(), err)
		return
	}
	if data == nil || data.Page == nil {
		err = fmt.Errorf("data or data.Page nil")
		log.Error("PgcIdx params(%s) error(%v)", r.Params(), err)
		return
	}
	data.Page.GetPageNb() // calculate page number
	return
}

// UgcIdx treats the ugc index request and call the ES to get the result
func (d *Dao) UgcIdx(c context.Context, req *seaMdl.SrvUgcIdx) (data *seaMdl.EsUgcResult, err error) {
	if len(req.TIDs) == 0 {
		err = ecode.RequestErr
		return
	}
	var (
		cfg = d.conf.Cfg.EsIdx.UgcIdx
		r   = d.esClient.NewRequest(cfg.Business).Index(cfg.Index).WhereEq("deleted", 0).
			WhereEq("valid", 1).WhereEq("result", 1).WhereIn("typeid", req.TIDs)
	)
	if pub := req.PubTime; pub != nil {
		r.WhereRange("pubtime", pub.STime, pub.ETime, elastic.RangeScopeLcRc)
	}
	r.Ps(req.Ps).Pn(int(req.Pn)).Order(req.UgcOrder(), seaMdl.IdxSort(req.Sort))
	if err = r.Scan(c, &data); err != nil {
		log.Error("PgcIdx:Scan params(%s) error(%v)", r.Params(), err)
		return
	}
	if data == nil || data.Page == nil {
		err = fmt.Errorf("data or data.Page nil")
		log.Error("PgcIdx params(%s) error(%v)", r.Params(), err)
		return
	}
	data.Page.GetPageNb() // calculate page number
	return
}
