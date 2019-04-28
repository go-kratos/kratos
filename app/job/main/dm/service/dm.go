package service

import (
	"context"

	"go-common/app/job/main/dm/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// flushTrimQueue 将数据库数据填充到redis中
func (s *Service) flushTrimQueue(c context.Context, tp int32, oid int64) (err error) {
	var (
		dms   []*model.DM
		trims []*model.Trim
	)
	if dms, err = s.dao.DMInfos(c, tp, oid); err != nil {
		return
	}
	for _, dm := range dms {
		// NOTE 只有普通弹幕会被顶掉
		if dm.Pool == model.PoolNormal && (dm.State == model.StateNormal || dm.State == model.StateMonitorAfter) {
			trim := &model.Trim{ID: dm.ID, Attr: dm.AttrVal(model.AttrProtect)}
			trims = append(trims, trim)
		}
	}
	return s.dao.FlushTrimCache(c, tp, oid, trims)
}

// addTrimQueue add dm index redis trim queue and return segment need flush.
func (s *Service) addTrimQueue(c context.Context, tp int32, oid, maxlimit int64, dms ...*model.DM) (err error) {
	var (
		ok             bool
		trimCnt, count int64
		trims          []*model.Trim
		dmids          []int64
	)
	for _, dm := range dms {
		// NOTE 只有普通弹幕并且弹幕状态处于正常或者先发后审状态的弹幕会被放入顶队列
		if dm.Pool == model.PoolNormal && dm.NeedDisplay() {
			trim := &model.Trim{ID: dm.ID, Attr: dm.AttrVal(model.AttrProtect)}
			trims = append(trims, trim)
		}
	}
	if len(trims) == 0 {
		return
	}
	if ok, err = s.dao.ExpireTrimQueue(c, tp, oid); err != nil {
		return
	}
	if !ok {
		if err = s.flushTrimQueue(c, tp, oid); err != nil {
			return
		}
	}
	if count, err = s.dao.AddTrimQueueCache(c, tp, oid, trims); err != nil {
		return
	}
	// NOTE 对于满弹幕的视频，始终保持两倍的候选弹幕集
	if trimCnt = count - 2*maxlimit; trimCnt > 0 {
		if dmids, err = s.dao.TrimCache(c, tp, oid, trimCnt); err != nil || len(trims) == 0 {
			return
		}
		if len(dmids) == 0 {
			return
		}
		if _, err = s.dao.UpdateDMStates(c, oid, dmids, model.StateHide); err != nil {
			return
		}
		if err = s.dao.DelIdxContentCaches(c, tp, oid, dmids...); err != nil {
			return
		}
		log.Info("oid:%d,trimCnt:%d,trims:%v", oid, len(dmids), dmids)
	}
	return
}

// recoverDM delete a dm and recover a hide state dm from db.
func (s *Service) recoverDM(c context.Context, typ int32, oid, rcvCnt int64) (dms []*model.DM, err error) {
	if dms, err = s.dao.DMHides(c, typ, oid, rcvCnt); err != nil {
		return
	}
	if len(dms) > 0 {
		var dmids []int64
		for _, dm := range dms {
			dmids = append(dmids, dm.ID)
			dm.State = model.StateNormal
		}
		if _, err = s.dao.UpdateDMStates(c, oid, dmids, model.StateNormal); err != nil {
			return
		}
		log.Info("recoverDM oid:%d dmids:%v", oid, dmids)
	}
	return
}

func (s *Service) subject(c context.Context, tp int32, oid int64) (sub *model.Subject, err error) {
	var cache = true
	if sub, err = s.dao.SubjectCache(c, tp, oid); err != nil {
		err = nil
		cache = false
	}
	if sub == nil {
		if sub, err = s.dao.Subject(c, tp, oid); err != nil {
			return
		}
		if sub == nil {
			sub = &model.Subject{
				Type: tp,
				Oid:  oid,
			}
		}
		if cache {
			s.cache.Do(c, func(ctx context.Context) {
				s.dao.SetSubjectCache(ctx, sub)
			})
		}
	}
	if sub.ID == 0 {
		err = ecode.NothingFound
		return
	}
	return
}
