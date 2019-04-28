package search

import (
	"context"
	"time"

	"fmt"
	model "go-common/app/interface/main/reply/model/reply"
	es "go-common/library/database/elastic"
	"go-common/library/log"
)

var (
	recordStates = []int64{0, 1, 2, 5, 6, 7, 9, 11}
)

type recordResult struct {
	Page struct {
		Num   int32 `json:"num"`
		Size  int32 `json:"size"`
		Total int32 `json:"total"`
	} `json:"page"`
	Result []*model.Record `json:"result"`
}

// RecordPaginate return a page of records from es.
func (d *Dao) RecordPaginate(c context.Context, types []int64, mid, stime, etime int64, order, sort string, pn, ps int32) (records []*model.Record, total int32, err error) {
	r := d.es.NewRequest("reply_record").Index(fmt.Sprintf("%s_%d", "replyrecord", mid%100)).
		WhereRange("ctime", time.Unix(stime, 0).Format(model.RecordTimeLayout), time.Unix(etime, 0).Format(model.RecordTimeLayout), es.RangeScopeLcRc).
		WhereIn("state", recordStates).
		WhereEq("mid", mid).
		Order(order, sort).Pn(int(pn)).Ps(int(ps))
	if len(types) > 0 {
		r = r.WhereIn("type", types)
	}
	var res recordResult
	err = r.Scan(c, &res)
	if err != nil {
		log.Error("r.Scan(%v) error(%v)", c, err)
		return
	}
	records = res.Result
	total = int32(res.Page.Total)
	return
}
