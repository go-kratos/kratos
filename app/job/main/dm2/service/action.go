package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/dm2/model"
	"go-common/library/log"
)

func (s *Service) actionAct(c context.Context, act *model.Action) (err error) {
	switch act.Action {
	case model.ActFlushDM:
		fc := new(model.Flush)
		if err = json.Unmarshal(act.Data, &fc); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", act.Data, err)
			return
		}
		s.asyncAddFlushDM(c, fc)
	case model.ActFlushDMSeg:
		fc := new(model.FlushDMSeg)
		if err = json.Unmarshal(act.Data, &fc); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", act.Data, err)
			return
		}
		if fc.Page == nil {
			log.Error("s.ActFlushDMSeg(+%v) error page nil", fc)
			return
		}
		// async flush cache
		s.asyncAddFlushDMSeg(c, fc)
	case model.ActAddDM:
		var (
			dm  = &model.DM{}
			sub *model.Subject
		)
		if err = json.Unmarshal(act.Data, &dm); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", act.Data, err)
			return
		}
		if sub, err = s.subject(c, dm.Type, dm.Oid); err != nil {
			return
		}
		if err = s.actionAddDM(c, sub, dm); err != nil {
			log.Error("s.actionAddDM(+%v) error(%v)", dm, err)
			return
		}
		if dm.State == model.StateNormal || dm.State == model.StateMonitorAfter {
			// 1. 创作中心最新1000条弹幕
			s.asyncAddRecent(c, dm)
			// 2. 刷新全段弹幕,NOTE 忽略redis缓存报错
			if ok, _ := s.dao.ExpireDMCache(c, dm.Type, dm.Oid); ok {
				s.dao.AddDMCache(c, dm)
			}
			s.asyncAddFlushDM(c, &model.Flush{
				Type:  dm.Type,
				Oid:   dm.Oid,
				Force: false,
			})
			// 3. 刷新分段弹幕缓存,NOTE 忽略redis缓存报错
			var p *model.Page
			if p, err = s.pageinfo(c, sub.Pid, dm); err != nil {
				return
			}
			switch dm.Pool {
			case model.PoolNormal:
				if ok, _ := s.dao.ExpireDMID(c, dm.Type, dm.Oid, p.Total, p.Num); ok {
					s.dao.AddDMIDCache(c, dm.Type, dm.Oid, p.Total, p.Num, dm.ID)
				}
			case model.PoolSubtitle:
				if ok, _ := s.dao.ExpireDMIDSubtitle(c, dm.Type, dm.Oid); ok {
					s.dao.AddDMIDSubtitleCache(c, dm.Type, dm.Oid, dm)
				}
			case model.PoolSpecial:
				if err = s.specialLocationUpdate(c, dm.Type, dm.Oid); err != nil {
					return
				}
				// TODO add cache
			default:
				return
			}
			s.dao.AddIdxContentCaches(c, dm.Type, dm.Oid, dm)
			s.asyncAddFlushDMSeg(c, &model.FlushDMSeg{
				Type:  dm.Type,
				Oid:   dm.Oid,
				Force: false,
				Page:  p,
			})
		}
		s.bnjDmCount(c, sub, dm)
	}
	return
}

func (s *Service) actionFlushDM(c context.Context, tp int32, oid int64, force bool) (err error) {
	sub, err := s.subject(c, tp, oid)
	if err != nil {
		return
	}
	if force {
		s.dao.DelDMCache(c, tp, oid) // delete redis cache,ignore error
	}
	xml, err := s.genXML(c, sub) // generate xml from redis or database
	if err != nil {
		log.Error("s.genXML(%d) error(%v)", oid, err)
		return
	}
	data, err := s.gzflate(xml, 4)
	if err != nil {
		log.Error("s.gzflate(type:%d,oid:%d) error(%v)", tp, oid, err)
		return
	}
	if err = s.dao.AddXMLCache(c, sub.Oid, data); err != nil {
		return
	}
	log.Info("actionFlushDM type:%d,oid:%d fore:%v", tp, oid, force)
	return
}

// actionAddDM add dm index and content to db by transaction.
func (s *Service) actionAddDM(c context.Context, sub *model.Subject, dm *model.DM) (err error) {
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		return
	}
	// special dm
	if dm.Pool == model.PoolSpecial && dm.ContentSpe != nil {
		if _, err = s.dao.TxAddContentSpecial(tx, dm.ContentSpe); err != nil {
			return tx.Rollback()
		}
	}
	if _, err = s.dao.TxAddContent(tx, dm.Oid, dm.Content); err != nil {
		return tx.Rollback()
	}
	if _, err = s.dao.TxAddIndex(tx, dm); err != nil {
		return tx.Rollback()
	}
	if dm.State == model.StateMonitorBefore || dm.State == model.StateMonitorAfter {
		if _, err = s.dao.TxIncrSubMCount(tx, dm.Type, dm.Oid); err != nil {
			return tx.Rollback()
		}
	}
	var count int64
	if dm.State == model.StateNormal || dm.State == model.StateMonitorAfter || dm.State == model.StateHide {
		count = 1
		if sub.Childpool == model.PoolNormal && dm.Pool != model.PoolNormal {
			sub.Childpool = 1
		}
	}
	if _, err = s.dao.TxIncrSubjectCount(tx, sub.Type, sub.Oid, 1, count, sub.Childpool); err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}

// actionFlushXMLDmSeg flush xml dm seg
func (s *Service) actionFlushXMLDmSeg(c context.Context, tp int32, oid int64, p *model.Page, force bool) (err error) {
	var (
		sub      *model.Subject
		duration int64
		seg      *model.Segment
	)
	if sub, err = s.subject(c, tp, oid); err != nil {
		return
	}
	if force {
		if err = s.dao.DelDMIDCache(c, tp, oid, p.Total, p.Num); err != nil {
			return
		}
		if sub.Childpool > 0 {
			s.dao.DelDMIDSubtitleCache(c, tp, oid)
		}
	}
	if duration, err = s.videoDuration(c, sub.Pid, sub.Oid); err != nil {
		return
	}
	ps, _ := model.SegmentPoint(p.Num, duration)
	if seg, err = s.segmentInfo(c, sub.Pid, sub.Oid, ps, duration); err != nil {
		return
	}
	res, err := s.dmSegXML(c, sub, seg)
	if err != nil {
		return
	}
	if err = s.dao.SetXMLSegCache(c, tp, oid, seg.Cnt, seg.Num, res); err != nil {
		return
	}
	log.Info("actionFlushXMLDmSeg type:%d,oid:%d,seg:%+v", tp, oid, seg)
	return
}

func (s *Service) flushDmSegCache(c context.Context, fc *model.FlushDMSeg) (err error) {
	if fc.Page == nil {
		return
	}
	if err = s.actionFlushXMLDmSeg(c, fc.Type, fc.Oid, fc.Page, fc.Force); err != nil {
		return
	}
	return
}

func (s *Service) flushDmCache(c context.Context, fc *model.Flush) (err error) {
	if err = s.actionFlushDM(c, fc.Type, fc.Oid, fc.Force); err != nil {
		return
	}
	if err = s.dao.DelAjaxDMCache(c, fc.Oid); err != nil {
		return
	}
	return
}
