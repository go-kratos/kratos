package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/library/queue/databus/report"

	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) subject(c context.Context, oid int64, typ int32) (sub *model.Subject, err error) {
	if sub, err = s.dao.Subject(c, oid, typ); err != nil {
		return
	}
	if sub == nil {
		err = ecode.NothingFound
	}
	return
}

func (s *Service) subjects(c context.Context, oids []int64, typ int32) (res map[int64]*model.Subject, err error) {
	return s.dao.Subjects(c, oids, typ)
}

// Subject get subject info.
func (s *Service) Subject(c context.Context, oid int64, typ int32) (sub *model.Subject, err error) {
	if sub, err = s.dao.Subject(c, oid, typ); err != nil {
		return
	}
	// NOTE nothing found return to normal
	return
}

// ModifySubState modify subject state
func (s *Service) ModifySubState(c context.Context, adid int64, adName string, oids []int64, typ int32, state int32, remark string) (fails map[int64]string, err error) {
	var (
		action string
		sub    *model.Subject
		now    = time.Now()
	)
	switch state {
	case model.SubStateNormal:
		action = model.SujectAllow
	case model.SubStateForbid:
		action = model.SujectForbid
	default:
		err = ecode.ReplyIllegalSubState
		return
	}
	fails = make(map[int64]string)
	for _, oid := range oids {
		sub, err = s.subject(c, oid, typ)
		if err != nil {
			log.Error("rpSvr.subject(oid,%d,type,%d)error(%v)", oid, typ, err)
			fails[oid] = ecode.Cause(err).Message()
			err = nil
			continue
		}
		// subject frozen
		if sub.AttrVal(model.SubAttrFrozen) == model.AttrYes {
			fails[oid] = ecode.ReplySubjectFrozen.Message()
			continue
		}
		if _, err = s.dao.UpSubjectState(c, oid, typ, state, now); err != nil {
			log.Error("s.UpSubjectState(%d,%d) error(%v)", oid, typ, err)
			return
		}
		if err = s.dao.DelSubjectCache(c, oid, typ); err != nil {
			log.Error("s.dao.DeleteSubjectCache(%d,%d) state:%d error(%v)", oid, typ, err)
		}
		report.Manager(&report.ManagerInfo{
			UID:      adid,
			Uname:    adName,
			Business: 41,
			Type:     int(typ),
			Oid:      oid,
			Ctime:    now,
			Action:   action,
			Index:    []interface{}{sub.State, state},
			Content:  map[string]interface{}{"remark": remark},
		})
	}
	return
}

// FreezeSub freeze or unfreeze subject
func (s *Service) FreezeSub(c context.Context, adid int64, adName string, oids []int64, typ int32, freeze int32, remark string) (fails map[int64]string, err error) {
	var (
		state  int32
		action string
		sub    *model.Subject
		now    = time.Now()
	)
	fails = make(map[int64]string)
	for _, oid := range oids {
		if sub, err = s.subject(c, oid, typ); err != nil {
			log.Error("rpSvr.subject(oid,%d,type,%d)error(%v)", oid, typ, err)
			fails[oid] = ecode.Cause(err).Message()
			err = nil
			continue
		}
		// already freeze or already unfreeze
		if (sub.AttrVal(model.SubAttrFrozen) == model.AttrYes && freeze == 2) || (sub.AttrVal(model.SubAttrFrozen) == model.AttrNo && freeze != 2) {
			continue
		}
		switch freeze {
		case 0:
			// unfreeze and allow
			sub.AttrSet(model.AttrNo, model.SubAttrFrozen)
			state = model.SubStateNormal
			action = model.SujectUnfrozenAllow
		case 1:
			// unfreeze and forbid
			sub.AttrSet(model.AttrNo, model.SubAttrFrozen)
			state = model.SubStateForbid
			action = model.SujectUnfrozenForbid
		case 2:
			// freeze and forbid
			sub.AttrSet(model.AttrYes, model.SubAttrFrozen)
			state = model.SubStateForbid
			action = model.SujectFrozen
		default:
			err = ecode.ReplyIllegalSubState
			return
		}
		if _, err = s.dao.UpStateAndAttr(c, oid, typ, state, sub.Attr, now); err != nil {
			log.Error("s.UpSubjectAttr(%d,%d,%d,%d) error(%v)", oid, typ, sub.Attr, err)
			return
		}
		if err = s.dao.DelSubjectCache(c, oid, typ); err != nil {
			log.Error("s.dao.DeleteSubjectCache(%d,%d) state:%d error(%v)", oid, typ, err)
		}
		report.Manager(&report.ManagerInfo{
			UID:      adid,
			Uname:    adName,
			Business: 41,
			Type:     int(typ),
			Oid:      oid,
			Ctime:    now,
			Action:   action,
			Index:    []interface{}{sub.State, state},
			Content:  map[string]interface{}{"remark": remark},
		})
	}
	return
}

// SubjectLog returns operation logs by query parameters
func (s *Service) SubjectLog(c context.Context, sp model.LogSearchParam) (res *model.SubjectLogRes, err error) {
	res = &model.SubjectLogRes{
		Logs: []*model.SubjectLog{},
	}
	sp.Action = "subject_allow,subject_forbid,subject_frozen,subject_unfrozen_allow,subject_unfrozen_forbid"
	if sp.Pn <= 0 || sp.Ps <= 0 {
		return nil, ecode.RequestErr
	}
	reportData, err := s.dao.ReportLog(c, sp)
	if err != nil || reportData == nil {
		return
	}
	res.Page = reportData.Page
	res.Sort = reportData.Sort
	res.Order = reportData.Order
	exists := make(map[int64]bool)
	var oids []int64
	for _, data := range reportData.Result {
		if !exists[data.Oid] {
			exists[data.Oid] = true
			oids = append(oids, data.Oid)
		}
		var extra map[string]string
		err = json.Unmarshal([]byte(data.Content), &extra)
		if err != nil {
			log.Error("Subject Operation Log unmarshal failed(%v)", err)
			return
		}
		res.Logs = append(res.Logs, &model.SubjectLog{
			AdminID:   data.AdminID,
			AdminName: data.AdminName,
			Oid:       strconv.FormatInt(data.Oid, 10),
			Type:      data.Type,
			Remark:    extra["remark"],
			CTime:     data.Ctime,
			Action:    data.Action,
		})
	}
	if sp.Type != 0 && len(oids) > 0 {
		var subjects map[int64]*model.Subject
		subjects, err = s.subjects(c, oids, sp.Type)
		if err == nil {
			for _, data := range res.Logs {
				var oid int64
				oid, err = strconv.ParseInt(data.Oid, 10, 64)
				if err != nil {
					log.Error("strconv.ParseInt failed(%v)", err)
					return
				}
				if sub, ok := subjects[oid]; ok && sub != nil {
					if sub.AttrVal(model.SubAttrFrozen) == model.AttrYes {
						//frozen
						data.State = 2
					} else if sub.State == model.SubStateForbid {
						data.State = model.SubStateForbid
					} else {
						data.State = model.SubStateNormal
					}
				}
			}
		}
	}
	return
}
