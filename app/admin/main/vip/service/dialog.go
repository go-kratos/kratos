package service

import (
	"context"
	"time"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// DialogAll .
func (s *Service) DialogAll(c context.Context, appID, platform int64, status string) (res []*model.ConfDialogList, err error) {
	var list []*model.ConfDialog
	if list, err = s.dao.DialogAll(c, appID, platform, status); err != nil {
		return
	}
	res = make([]*model.ConfDialogList, len(list))
	for i, v := range list {
		res[i] = &model.ConfDialogList{ConfDialog: v, Status: dialogStatus(v)}
	}
	return
}

func dialogStatus(v *model.ConfDialog) (status string) {
	curr := time.Now()
	switch {
	case v.Stage && v.StartTime.Time().After(curr):
		return "padding"
	case v.Stage && v.StartTime.Time().Before(curr) && (int64(v.EndTime) == 0 || v.EndTime.Time().After(curr)):
		return "active"
	case !v.Stage || (int64(v.EndTime) > 0 && v.EndTime.Time().Before(curr)):
		return "inactive"
	default:
		return ""
	}
}

// DialogByID .
func (s *Service) DialogByID(c context.Context, arg *model.ArgID) (dlg *model.ConfDialog, err error) {
	return s.dao.DialogByID(c, arg.ID)
}

// DialogSave .
func (s *Service) DialogSave(c context.Context, arg *model.ConfDialog) (eff int64, err error) {
	var db *model.ConfDialog
	if arg.ID != 0 {
		if db, err = s.dao.DialogByID(c, arg.ID); err != nil {
			return
		}
		if db != nil {
			arg.Ctime = db.Ctime
			arg.Stage = db.Stage
		}
	}
	// PRD: 如果平台、app_id完全相同时，生效时间不可冲突，否则不允许保存
	var exist []*model.ConfDialog
	if exist, err = s.dao.DialogBy(c, arg.AppID, arg.Platform, arg.ID); err != nil {
		return
	}
	for _, v := range exist {
		if v.StartTime.Time().Unix() <= arg.EndTime.Time().Unix() && v.EndTime.Time().Unix() >= arg.StartTime.Time().Unix() {
			err = ecode.VipDialogConflictErr
			return
		}
	}

	return s.dao.DialogSave(c, arg)
}

// DialogEnable .
func (s *Service) DialogEnable(c context.Context, arg *model.ConfDialog) (eff int64, err error) {
	var db *model.ConfDialog
	if arg.ID != 0 {
		if db, err = s.dao.DialogByID(c, arg.ID); err != nil {
			return
		}
	}
	if db == nil || db.StartTime.Time().After(time.Now()) {
		err = ecode.VipDialogTimeErr
		return
	}
	return s.dao.DialogEnable(c, arg)
}

// DialogDel .
func (s *Service) DialogDel(c context.Context, arg *model.ArgID, operator string) (eff int64, err error) {
	eff, err = s.dao.DialogDel(c, arg.ID)
	log.Warn("user(%s) delete dialog(%d) effect row(%d)", operator, arg.ID, eff)
	return
}
