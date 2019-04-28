package service

import (
	"context"
	"time"

	"go-common/app/admin/main/space/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

// Notice get notice data.
func (s *Service) Notice(c context.Context, mid int64) (data *model.Notice, err error) {
	data = &model.Notice{Mid: mid}
	if err = s.dao.DB.Table(data.TableName()).Where("mid=?", mid).First(&data).Error; err != nil {
		log.Error("Notice (mid:%d) error (%v)", mid, err)
		if err == ecode.NothingFound {
			err = nil
		}
	}
	return
}

// NoticeUp notice clear and forbid.
func (s *Service) NoticeUp(c context.Context, arg *model.NoticeUpArg) (err error) {
	var action string
	notice := &model.Notice{Mid: arg.Mid}
	if err = s.dao.DB.Table(notice.TableName()).Where("mid=?", arg.Mid).First(&notice).Error; err != nil {
		log.Error("NoticeForbid error (mid:%d) (%v)", arg.Mid, err)
		if err != ecode.NothingFound {
			return
		}
	}
	up := make(map[string]interface{})
	switch arg.Type {
	case model.NoticeTypeClear:
		up["notice"] = ""
		action = model.NoticeClear
	case model.NoticeTypeClearAndForbid:
		up["notice"] = ""
		up["is_forbid"] = model.NoticeForbid
		action = model.NoticeClearAndForbid
	case model.NoticeTypeUnForbid:
		up["is_forbid"] = model.NoticeNoForbid
		action = model.NoticeUnForbid
	}
	if err != ecode.NothingFound {
		if err = s.dao.DB.Table(notice.TableName()).Where("id=?", notice.ID).Update(up).Error; err != nil {
			log.Error("NoticeForbid (mid:%d) update error (%v)", arg.Mid, err)
			return
		}
	} else {
		create := &model.Notice{Mid: arg.Mid}
		if arg.Type == model.NoticeTypeClearAndForbid {
			create.IsForbid = model.NoticeForbid
		}
		if err = s.dao.DB.Table(notice.TableName()).Create(create).Error; err != nil {
			log.Error("NoticeForbid (mid:%d) insert error (%v)", arg.Mid, err)
			return
		}
	}
	if err = report.Manager(&report.ManagerInfo{
		Uname:    arg.Uname,
		UID:      arg.UID,
		Business: model.NoticeLogID,
		Type:     0,
		Oid:      arg.Mid,
		Action:   action,
		Ctime:    time.Now(),
		Content: map[string]interface{}{
			"old": notice,
		},
	}); err != nil {
		return
	}
	return
}
