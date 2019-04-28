package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/dm2/model"
	"go-common/library/log"
)

func (s *Service) trackSubject(c context.Context, m *model.BinlogMsg) (err error) {
	nw := &model.Subject{}
	if err = json.Unmarshal(m.New, &nw); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", m.New, err)
		return
	}
	switch m.Action {
	case "insert":
		if err = s.dao.AddSubjectCache(c, nw); err != nil {
			log.Error("s.dao.AddSubjectCache(%v) error(%v)", nw, err)
			return
		}
	case "delete":
		if err = s.dao.DelSubjectCache(c, nw.Type, nw.Oid); err != nil {
			log.Error("s.dao.DelSubjectCahce(%v) error(%v)", nw, err)
			return
		}
	case "update":
		old := model.Subject{}
		if err = json.Unmarshal(m.Old, &old); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", m.Old, err)
			return
		}
		if err = s.dao.AddSubjectCache(c, nw); err != nil { // 全量缓存subject
			log.Error("s.dao.AddSubjectCache(%v) error(%v)", nw, err)
			return
		}
		if nw.Childpool != old.Childpool || nw.Maxlimit != old.Maxlimit || nw.State != old.State {
			// 立刻刷新全段弹幕缓存
			flush := &model.Flush{Oid: nw.Oid, Type: nw.Type, Force: true}
			s.flushDmCache(c, flush)
			// 立刻刷新分段弹幕缓存
			s.flushXMLSegCache(c, nw)
		}
	}
	return
}

func (s *Service) trackIndex(c context.Context, m *model.BinlogMsg) (err error) {
	if m.Action != "update" {
		return
	}
	dm := &model.DM{}
	old := &model.DM{}
	if err = json.Unmarshal(m.New, &dm); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", m.New, err)
		return
	}
	if err = json.Unmarshal(m.Old, &old); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", m.Old, err)
		return
	}
	s.asyncAddRecent(c, dm) // 更新up主最新1000条弹幕
	s.asyncAddFlushDM(c, &model.Flush{
		Type:  dm.Type,
		Oid:   dm.Oid,
		Force: true,
	}) // 刷新全段弹幕
	sub, err := s.subject(c, dm.Type, dm.Oid)
	if err != nil {
		return
	}
	p, err := s.pageinfo(c, sub.Pid, dm)
	if err != nil {
		return
	}
	if dm.NeedUpdateSpecial(old) {
		if err = s.specialLocationUpdate(c, dm.Type, dm.Oid); err != nil {
			return
		}
	}
	s.dao.DelIdxContentCaches(c, dm.Type, dm.Oid, dm.ID) // 删除content cache
	s.asyncAddFlushDMSeg(c, &model.FlushDMSeg{
		Type:  dm.Type,
		Oid:   dm.Oid,
		Force: true,
		Page:  p,
	})
	return
}

func (s *Service) trackVideoup(c context.Context, aid int64) (err error) {
	var (
		retry  = 5
		tp     = model.SubTypeVideo
		videos []*model.Video
	)
	for i := 0; i < retry; i++ {
		if videos, err = s.dao.Videos(c, aid); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Error("track video failed,aid(%d),error(%v)", aid, err)
		return
	}
	for _, v := range videos {
		for i := 0; i < retry; i++ {
			if err = s.syncVideo(c, tp, v); err == nil {
				break
			}
			time.Sleep(time.Second)
		}
	}
	return
}

func (s *Service) syncVideo(c context.Context, tp int32, v *model.Video) (err error) {
	log.Info("sync video:%+v", v)
	sub, err := s.dao.Subject(c, tp, v.Cid)
	if err != nil {
		return
	}
	if sub == nil {
		if v.XCodeState >= model.VideoXcodeHDFinish {
			// 生成弹幕蒙版
			var attr int32
			for _, mid := range s.maskMid {
				if mid == v.Mid {
					if err = s.dao.GenerateMask(c, v.Cid, mid, model.MaskPlatAll, model.MaskPriorityHgih, v.Aid, 0, 0); err != nil {
						break
					}
					attr = attr | (model.AttrYes << model.AttrSubMaskOpen)
					break
				}
			}
			if _, err = s.dao.AddSubject(c, tp, v.Cid, v.Aid, v.Mid, s.maxlimit(v.Duration), attr); err != nil {
				return
			}
		}

	} else {
		if sub.Mid != v.Mid {
			if _, err = s.dao.UpdateSubMid(c, tp, v.Cid, v.Mid); err != nil {
				return
			}
			if err = s.updateSubtilte(c, tp, v); err != nil {
				log.Error("updateSubtilte(params:%+v),error(%v)", v, err)
				return
			}
		}
	}
	return
}

func (s *Service) updateSubtilte(c context.Context, tp int32, v *model.Video) (err error) {
	var (
		subtitles []*model.Subtitle
		subtitle  *model.Subtitle
	)
	if subtitles, err = s.dao.GetSubtitles(c, tp, v.Cid); err != nil {
		log.Error("updateSubtilte(params:%+v),error(%v)", v, err)
		return
	}
	for _, subtitle = range subtitles {
		subtitle.UpMid = v.Mid
		if err = s.dao.UpdateSubtitle(c, subtitle); err != nil {
			log.Error("updateSubtilte(params:%+v),error(%v)", v, err)
			return
		}
		s.dao.DelSubtitleCache(c, v.Cid, subtitle.ID)
		if subtitle.Status == model.SubtitleStatusDraft || subtitle.Status == model.SubtitleStatusToAudit {
			s.dao.DelSubtitleDraftCache(c, v.Cid, tp, subtitle.Mid, subtitle.Lan)
		}
	}
	s.dao.DelVideoSubtitleCache(c, v.Cid, tp)
	return
}

func (s *Service) maxlimit(duration int64) (limit int64) {
	switch {
	case duration == 0:
		limit = 1500
	case duration > 3600:
		limit = 8000
	case duration > 2400:
		limit = 6000
	case duration > 900:
		limit = 3000
	case duration > 600:
		limit = 1500
	case duration > 150:
		limit = 1000
	case duration > 60:
		limit = 500
	case duration > 30:
		limit = 300
	case duration <= 30:
		limit = 100
	}
	return
}
