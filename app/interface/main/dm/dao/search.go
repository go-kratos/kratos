package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/dm/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

// return recent two years report search index.
func (d *Dao) rptSearchIndex() string {
	year := time.Now().Year()
	return fmt.Sprintf("dmreport_%d,dmreport_%d", year-1, year)
}

// SearchReportAid 根据mid获取用户的所有有举报弹幕的稿件
func (d *Dao) SearchReportAid(c context.Context, mid int64, upOp int8, states []int8, pn, ps int64) (aids []int64, err error) {
	res := new(model.SearchReportAidResult)
	order := []map[string]string{{"arc_aid": "desc"}}
	req := d.es.NewRequest("dmreport").Fields("arc_aid").Index(d.rptSearchIndex())
	req.WhereEq("arc_mid", mid).WhereEq("up_op", upOp).WhereIn("state", states).GroupBy(elastic.EnhancedModeGroupBy, "arc_aid", order)
	req.Pn(int(pn)).Ps(int(ps))
	if err = req.Scan(c, res); err != nil {
		log.Error("search params(%s), err(%v)", req.Params(), err)
		return
	}
	if values, ok := res.Result["group_by_arc_aid"]; ok {
		for _, v := range values {
			var aid int64
			if aid, err = strconv.ParseInt(v.Key, 10, 64); err != nil {
				log.Error("strconv.ParseInt(%s) error(%v)", v.Key, err)
				return
			}
			aids = append(aids, aid)
		}
	}
	return
}

// SearchReport 根据up主id,稿件id获取举报弹幕列表
func (d *Dao) SearchReport(c context.Context, mid, aid, pn, ps int64, upOp int8, states []int64) (res *model.SearchReportResult, err error) {
	req := d.es.NewRequest("dmreport")
	req.Fields("id", "dmid", "cid", "arc_aid", "arc_typeid", "dm_owner_uid", "dm_msg", "count", "content", "up_op", "state",
		"uid", "rp_time", "reason", "dm_deleted", "arc_mid", "pool_id", "model", "score", "dm_ctime", "ctime", "mtime")
	req.Index(d.rptSearchIndex())
	if aid > 0 {
		req.WhereEq("arc_aid", aid)
	}
	if mid > 0 {
		req.WhereEq("arc_mid", mid)
	}
	if len(states) > 0 {
		req.WhereIn("state", states)
	}
	req.WhereNot(elastic.NotTypeEq, "dm_owner_uid").WhereEq("dm_owner_uid", 0)
	req.WhereEq("up_op", upOp)
	req.Order("rp_time", "desc")
	req.Pn(int(pn)).Ps(int(ps))
	res = &model.SearchReportResult{}
	if err = req.Scan(c, res); err != nil {
		log.Error("req.Scan() search params(%s), err(%v)", req.Params(), err)
	}
	return
}

// UpdateSearchReport update report search index.
func (d *Dao) UpdateSearchReport(c context.Context, rpts []*model.UptSearchReport) (err error) {
	up := d.es.NewUpdate("dmreport").Insert() // if data not exist insert,else update
	for _, rpt := range rpts {
		t, err1 := time.ParseInLocation("2006-01-02 15:04:05", rpt.Ctime, time.Local)
		if err1 != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", rpt.Ctime, err1)
			return err1
		}
		up.AddData(fmt.Sprintf("dmreport_%d", t.Year()), rpt)
	}
	if err = up.Do(c); err != nil {
		log.Error("update.Do() error(%v)", err)
	}
	return
}
