package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/search"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

const (
	_wkfReplyLog = 11
	_wkfAuditLog = 12
)

// LastLog .
func (s *Service) LastLog(c context.Context, targets []int64, modules []int) (logs map[int64]string, err error) {
	var (
		resp   *search.AuditLogSearchCommonResult
		uids   []int64
		uNames map[int64]string
	)
	logs = make(map[int64]string)
	if len(targets) == 0 {
		return
	}

	cond := &search.AuditReportSearchCond{
		Fields:        []string{"int_1", "ctime", "str_0", "uid"},
		Business:      _wkfAuditLog,
		Type:          modules,
		Order:         "ctime",
		Sort:          "desc",
		Int1:          targets,
		Distinct:      "int_1",
		IndexTimeType: "year",
		IndexTimeFrom: time.Now().AddDate(-1, 0, 0),
		IndexTimeEnd:  time.Now(),
	}
	if resp, err = s.dao.SearchAuditReportLog(c, cond); err != nil {
		log.Error("s.dao.SearchAuditReportLog(%+v) error(%v)", cond, err)
		return
	}
	if resp == nil {
		log.Error("len resp.result == 0")
		err = ecode.Degrade
		return
	}
	// wrap uname
	for _, l := range resp.Result {
		uids = append(uids, l.UID)
	}
	if uNames, err = s.dao.BatchUNameByUID(c, uids); err != nil {
		log.Error("s.dao.SearchUNameByUid(%v) error(%v)", uids, err)
		err = nil
	}
	for _, l := range resp.Result {
		logs[l.Int1] = fmt.Sprintf("%s\n操作时间:%s\n操作人:", l.Str0, l.CTime)
		if uname, ok := uNames[l.UID]; ok {
			logs[l.Int1] = fmt.Sprint(logs[l.Int1], uname)
		} else {
			logs[l.Int1] = fmt.Sprint(logs[l.Int1], l.UID)
		}
	}
	return
}

// LastLogStat .
func (s *Service) LastLogStat(c context.Context, targets []int64, modules []int, fields []string) (logs map[int64]*search.ReportLog, err error) {
	var resp *search.AuditLogSearchCommonResult
	logs = make(map[int64]*search.ReportLog)
	if len(targets) == 0 {
		return
	}

	cond := &search.AuditReportSearchCond{
		Fields:        fields,
		Business:      _wkfAuditLog,
		Type:          modules,
		Order:         "ctime",
		Sort:          "desc",
		Int1:          targets,
		Distinct:      "int_1",
		IndexTimeType: "year",
		IndexTimeFrom: time.Now().AddDate(-1, 0, 0),
		IndexTimeEnd:  time.Now(),
	}
	if resp, err = s.dao.SearchAuditReportLog(c, cond); err != nil {
		log.Error("s.dao.SearchAuditReportLog(%+v) error(%v)", cond, err)
		return
	}
	if resp == nil {
		log.Error("len resp.result == 0")
		err = ecode.NothingFound
		return
	}

	for _, l := range resp.Result {
		logs[l.Int1] = l
	}
	return
}

// AllAuditLog search all audit log of target & modules
func (s *Service) AllAuditLog(c context.Context, target int64, modules []int) (logs []*model.WLog, err error) {
	var (
		resp   *search.AuditLogSearchCommonResult
		uids   []int64
		uNames map[int64]string
	)

	if target == 0 {
		return
	}
	cond := &search.AuditReportSearchCond{
		Fields:        []string{"int_1", "ctime", "str_0", "uid", "uname"},
		Business:      _wkfAuditLog,
		Type:          modules,
		Order:         "ctime",
		Sort:          "desc",
		Int1:          []int64{target},
		IndexTimeType: "year",
		IndexTimeFrom: time.Now().AddDate(-1, 0, 0),
		IndexTimeEnd:  time.Now(),
	}
	if resp, err = s.dao.SearchAuditReportLog(c, cond); err != nil {
		log.Error("s.dao.SearchAuditReportLog(%+v) error(%v)", cond, err)
		return
	}
	if resp == nil {
		log.Error("len resp.result == 0")
		err = ecode.Degrade
		return
	}
	// wrap uname
	for _, l := range resp.Result {
		uids = append(uids, l.UID)
	}
	if uNames, err = s.dao.BatchUNameByUID(c, uids); err != nil {
		log.Error("s.dao.SearchUNameByUid(%v) error(%v)", uids, err)
		err = nil
	}

	for _, r := range resp.Result {
		wl := &model.WLog{
			AdminID: r.UID,
			Admin:   r.UName,
			Target:  r.Int1,
			Remark:  r.Str0,
		}
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", r.CTime, time.Local)
		wl.CTime.Scan(t)
		wl.Admin = uNames[wl.AdminID]
		logs = append(logs, wl)
	}
	return
}

func (s *Service) writeAuditLog(l *model.WLog) {
	var err error
	info := &report.ManagerInfo{
		Uname:    l.Admin,
		UID:      l.AdminID,
		Business: _wkfAuditLog,
		Type:     int(l.Module),
		Oid:      l.Oid,
		Action:   "audit_log",
		Ctime:    time.Now(),
		Index:    []interface{}{l.Business, l.Target, l.TimeConsume, l.Mid, l.Remark, l.Note, l.OpType, l.PreRid},
		Content:  map[string]interface{}{"wlog": l, "param": l.Param, "mids": l.Mids},
	}
	log.Info("start report audit log target:%v oid:%v uid:%v business:%v mid:%v", l.Target, l.Oid, l.AdminID, l.Business, l.Mid)
	if err = report.Manager(info); err != nil {
		log.Error("failed to produce report.Manager(%+v), err(%v)", info, err)
	}
}

func (s *Service) writeReplyLog(l *model.WLog) {
	var err error
	info := &report.ManagerInfo{
		Uname:    l.Admin,
		UID:      l.AdminID,
		Business: _wkfReplyLog,
		Type:     int(l.Module),
		Oid:      l.Oid,
		Action:   "reply_log",
		Ctime:    time.Now(),
		Index:    []interface{}{l.Business, l.Target, l.Mid, l.Remark, l.Note},
		Content:  map[string]interface{}{"wlog": l},
	}
	log.Info("start report reply log target:%v oid:%v uid:%v business:%v mid:%v", l.Target, l.Oid, l.AdminID, l.Business, l.Mid)
	if err = report.Manager(info); err != nil {
		log.Error("failed to produce report.Manager(%+v), err(%v)", info, err)
	}
}
