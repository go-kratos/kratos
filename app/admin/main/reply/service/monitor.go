package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/reply/model"
	accmdl "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

// MonitorStats return monitor stats.
func (s *Service) MonitorStats(c context.Context, mode, page, pageSize int64, adminIDs, sort, order, startTime, endTime string) (res *model.StatsMonitorResult, err error) {
	return s.dao.MonitorStats(c, mode, page, pageSize, adminIDs, sort, order, startTime, endTime)
}

// MonitorSearch return monitor result from search.
func (s *Service) MonitorSearch(c context.Context, sp *model.SearchMonitorParams, page, pageSize int64) (res *model.SearchMonitorResult, err error) {
	if res, err = s.dao.SearchMonitor(c, sp, page, pageSize); err != nil {
		log.Error("s.dao.SearchMonitor(%v,%d,%d) error(%v)", sp, page, pageSize, err)
	}
	return
}

// UpMonitorState set monitor state into subject attr.
func (s *Service) UpMonitorState(c context.Context, adminID int64, adName string, oid int64, typ, state int32, remark string) (err error) {
	sub, err := s.dao.Subject(c, oid, typ)
	if err != nil {
		return
	}
	if sub == nil {
		err = ecode.NothingFound
		return
	}
	var logState int32
	switch state {
	case model.MonitorClose:
		if sub.AttrVal(model.SubAttrMonitor) == model.AttrYes {
			logState = model.AdminOperSubMonitorClose
		} else if sub.AttrVal(model.SubAttrAudit) == model.AttrYes {
			logState = model.AdminOperSubAuditClose
		} else {
			err = ecode.ReplyIllegalSubState
			return
		}
		sub.AttrSet(model.AttrNo, model.SubAttrMonitor)
		sub.AttrSet(model.AttrNo, model.SubAttrAudit)
	case model.MonitorOpen:
		sub.AttrSet(model.AttrYes, model.SubAttrMonitor)
		sub.AttrSet(model.AttrNo, model.SubAttrAudit)
		logState = model.AdminOperSubMonitorOpen
	case model.MonitorAudit:
		sub.AttrSet(model.AttrNo, model.SubAttrMonitor)
		sub.AttrSet(model.AttrYes, model.SubAttrAudit)
		logState = model.AdminOperSubAuditOpen
	default:
		err = ecode.RequestErr
		return
	}
	// update attr
	now := time.Now()
	if _, err = s.dao.UpSubjectAttr(c, oid, typ, sub.Attr, now); err != nil {
		log.Error("s.dao.UpSubjectAttr(%d,%d,%d,%d) state:%d error(%v)", typ, oid, sub.Attr, state, err)
		return
	}
	if err = s.dao.DelSubjectCache(c, oid, typ); err != nil {
		log.Error("MonitorState del subject cache error(%v)", err)
	}
	// update search index
	if err = s.dao.UpSearchMonitor(c, sub, remark); err != nil {
		log.Error("s.dao.UpdateMonitor(%v) error(%v)", sub, err)
		return
	}
	s.dao.AddAdminLog(c, []int64{oid}, []int64{0}, adminID, typ, model.AdminIsNew, model.AdminIsNotReport, logState, fmt.Sprintf("修改监控状态为: %d", state), remark, now)
	report.Manager(&report.ManagerInfo{
		UID:      adminID,
		Uname:    adName,
		Business: 41,
		Type:     int(typ),
		Oid:      oid,
		Ctime:    now,
		Action:   model.ReportActionReplyMonitor,
		Content: map[string]interface{}{
			"remark": remark,
		},
		Index: []interface{}{sub.Mid, logState, state},
	})
	return
}

// MointorLog MointorLog
func (s *Service) MointorLog(c context.Context, sp model.LogSearchParam) (result *model.MonitorLogResult, err error) {
	var (
		mids     []int64
		userInfo map[int64]*accmdl.Info
	)
	adNames := map[int64]string{}

	result = &model.MonitorLogResult{
		Logs: []*model.MonitorLog{},
	}

	sp.Action = "monitor"
	reportData, err := s.dao.ReportLog(c, sp)
	if err != nil {
		return
	}
	result.Page = reportData.Page
	result.Sort = reportData.Sort
	result.Order = reportData.Order

	for i, data := range reportData.Result {
		mid := data.Index0
		reportData.Result[i].OidStr = strconv.FormatInt(reportData.Result[i].Oid, 10)
		logState := data.Index1
		state := data.Index2
		title, link, _ := s.TitleLink(c, data.Oid, data.Type)
		var extra map[string]string
		if data.Content != "" {
			err = json.Unmarshal([]byte(data.Content), &extra)
			if err != nil {
				log.Error("MointorLog unmarshal failed!err:=%v", err)
				return
			}
		}
		if data.AdminName == "" {
			adNames[data.AdminID] = ""
		}
		result.Logs = append(result.Logs, &model.MonitorLog{
			Mid:         mid,
			AdminID:     data.AdminID,
			AdminName:   data.AdminName,
			Oid:         data.Oid,
			OidStr:      data.OidStr,
			Type:        data.Type,
			Remark:      extra["remark"],
			CTime:       data.Ctime,
			LogState:    logState,
			State:       state,
			Title:       title,
			RedirectURL: link,
		})
		mids = append(mids, mid)
	}
	if len(adNames) > 0 {
		s.dao.AdminName(c, adNames)
	}
	if len(mids) > 0 {
		var infosReply *accmdl.InfosReply
		infosReply, err = s.accSrv.Infos3(c, &accmdl.MidsReq{Mids: mids})
		if err != nil {
			log.Error(" s.accSrv.Infos3 (%v) error(%v)", mids, err)
			err = nil
			return
		}
		userInfo = infosReply.Infos
	}

	for _, log := range result.Logs {
		if user, ok := userInfo[log.Mid]; ok {
			log.UserName = user.GetName()
		}
		if log.AdminName == "" {
			log.AdminName = adNames[log.AdminID]
		}
	}
	return
}
