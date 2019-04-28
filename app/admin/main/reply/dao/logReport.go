package dao

import (
	"context"
	"strconv"

	"go-common/app/admin/main/reply/model"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
	"strings"
	"time"
)

const (
	// log api
	_logURL = "/x/admin/search/log"
)

// ReportLogData report log data
type ReportLogData struct {
	Sort   string             `json:"sort"`
	Order  string             `json:"order"`
	Page   model.Page         `json:"page"`
	Result []*ReportLogResult `json:"result"`
}

// ReportLogResult ReportLogResult
type ReportLogResult struct {
	AdminName string `json:"uname"`
	AdminID   int64  `json:"uid"`
	Business  int64  `json:"business"`
	Type      int32  `json:"type"`
	Oid       int64  `json:"oid"`
	OidStr    string `json:"oid_str"`
	Action    string `json:"action"`
	Ctime     string `json:"ctime"`
	Index0    int64  `json:"int_0"`
	Index1    int64  `json:"int_1"`
	Index2    int64  `json:"int_2"`
	Content   string `json:"extra_data"`
}

// ReportLog get notice info.
func (d *Dao) ReportLog(c context.Context, sp model.LogSearchParam) (data *ReportLogData, err error) {
	var (
		// log_audit评论的索引是按时间分区的，最早的数据是2018年的，所以这里最早时间是写死的2018年
		stime = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
		etime = time.Now()
	)
	if sp.CtimeFrom != "" {
		if stime, err = time.Parse(model.DateFormat, sp.CtimeFrom); err != nil {
			log.Error("time.Parse(%v) error", sp.CtimeFrom)
			stime = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
		}
	}
	if sp.CtimeTo != "" {
		if etime, err = time.Parse(model.DateFormat, sp.CtimeTo); err != nil {
			log.Error("time.Parse(%v) error", sp.CtimeTo)
			etime = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
		}
	}
	action := strings.Split(sp.Action, ",")
	r := d.es.NewRequest("log_audit").IndexByTime("log_audit_41", elastic.IndexTypeYear, stime, etime).WhereIn("action", action).
		WhereRange("ctime", stime.Format(model.DateFormat), etime.Format(model.DateFormat), elastic.RangeScopeLcRc)
	if sp.Pn <= 0 {
		sp.Pn = 1
	}
	if sp.Ps <= 0 {
		sp.Ps = 20
	}
	if sp.Order == "" {
		sp.Order = "ctime"
	}
	if sp.Sort == "" {
		sp.Sort = "desc"
	}
	r = r.Order(sp.Order, sp.Sort).Pn(int(sp.Pn)).Ps(int(sp.Ps))
	if sp.Oid > 0 {
		r = r.WhereEq("oid", strconv.FormatInt(sp.Oid, 10))
	}
	if sp.Mid > 0 {
		r = r.WhereEq("int_0", sp.Mid)
	}
	if sp.Type > 0 {
		r = r.WhereEq("type", sp.Type)
	}
	if sp.Other > 0 {
		r = r.WhereEq("int_1", sp.Other)
	}
	log.Warn(r.Params())
	err = r.Scan(c, &data)
	if err != nil {
		log.Error("r.Scan(%v) error(%v)", c, err)
		return
	}
	if data == nil {
		err = ecode.ServerErr
		log.Error("log_audit error")
		return
	}
	return
}
