package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/dm/model"
	"go-common/library/log"
)

// trackDMMeta 顶弹幕逻辑  保持pool0的弹幕池只有maxlimt*2的数量
func (s *Service) trackDMMeta(c context.Context, m *model.BinlogMsg) (err error) {
	var (
		sub *model.Subject
		nw  = &model.DM{}
	)
	if err = json.Unmarshal(m.New, &nw); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", m.New, err)
		return
	}
	if sub, err = s.subject(c, model.SubTypeVideo, nw.Oid); err != nil {
		log.Error("s.subject(%d) error(%v)", nw.Oid, err)
		return
	}
	if sub == nil {
		err = errSubNotExist
		return
	}
	switch m.Action {
	case model.SyncInsert:
		if sub.Count >= sub.Maxlimit {
			if err = s.addTrimQueue(c, nw.Type, nw.Oid, sub.Maxlimit, nw); err != nil {
				log.Error("s.addTrimQueue(%v) error(%v)", nw, err)
				return err
			}
		}
	case model.SyncUpdate:
		old := &model.DM{}
		if err = json.Unmarshal(m.Old, &old); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", m.Old, err)
			return
		}
		if nw.NeedStateNormal(old) {
			nw.State = model.StateNormal
			if _, err = s.dao.UpdateDM(c, nw); err != nil {
				log.Error("dao.UpdateDM(%v) error(%v)", nw, err)
				return err
			}
		}
		if sub.Count >= sub.Maxlimit {
			dms := make([]*model.DM, 0)
			if isDelOperation(nw, old) {
				if err = s.dao.ZRemTrimCache(c, nw.Type, nw.Oid, nw.ID); err != nil {
					return
				}
				if dms, err = s.recoverDM(c, nw.Type, nw.Oid, 1); err != nil {
					log.Error("s.recoverIdx(%d) error(%v)", nw.Oid, err)
					return
				}
			}
			dms = append(dms, nw)
			if err = s.addTrimQueue(c, nw.Type, nw.Oid, sub.Maxlimit, dms...); err != nil {
				log.Error("s.addTrimQueue(%v) error(%v)", dms, err)
				return
			}
		}
	case model.SyncDelete:
	}
	return
}

func isDelOperation(nw, old *model.DM) bool {
	if nw.State != model.StateHide && old.NeedDisplay() && !nw.NeedDisplay() { // 弹幕从展示变为非展示状态
		return true
	}
	if nw.Pool != old.Pool && (nw.Pool == model.PoolSpecial || nw.Pool == model.PoolSubtitle) {
		return true
	}
	return false
}
