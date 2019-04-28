package service

import (
	"context"
	"encoding/json"

	arcModel "go-common/app/job/main/archive/model/archive"
	"go-common/app/job/main/tag/model"
	"go-common/library/log"
)

// 稿件更新
func (s *Service) upTagArcCache(c context.Context, newArc, oldArc *model.Archive) (err error) {
	var tids []int64
	if tids, err = s.dao.Resources(c, oldArc.Aid, 3); err != nil {
		log.Info("s.dao.Resourcess(%d) error(%v)", oldArc.Aid, err)
		return
	}
	if len(tids) == 0 {
		log.Error("s.dao.Resources(%d) tids is null", oldArc.Aid)
		return
	}
	if err = s.dao.RemTidArcCache(c, oldArc.Aid, tids...); err != nil {
		log.Info("s.dao.RemTidArcCache(%v,%v) error(%v)", oldArc, tids, err)
	}
	if err = s.dao.RemoveRegionNewArcCache(c, oldArc, tids...); err != nil {
		log.Error("d.RemoveRegionNewArcCache(%v,%v) err(%v)", oldArc, tids, err)
	}
	if prid, ok := s.rpMap[int64(oldArc.TypeID)]; ok {
		if err = s.dao.RemTagPridArcCache(c, oldArc.Aid, prid, tids); err != nil {
			log.Error("d.RemTagPridArcCache(%v,%d,%v) err(%v)", oldArc, prid, tids, err)
		}
	}
	err = s.addNew(c, tids, newArc)
	return
}

// 插入更新
func (s *Service) insertArcCache(c context.Context, newArc *model.Archive) (err error) {
	var tids []int64
	if tids, err = s.dao.Resources(c, newArc.Aid, 3); err != nil {
		log.Info("s.dao.Resourcess(%d) error(%v)", newArc.Aid, err)
		return
	}
	if len(tids) == 0 {
		log.Error("s.dao.Resources(%d) tids is null", newArc.Aid)
		return
	}
	err = s.addNew(c, tids, newArc)
	return
}

func (s *Service) addNew(c context.Context, tids []int64, newArc *model.Archive) (err error) {
	if newArc.State >= arcModel.StateOpen || newArc.State == arcModel.StateForbidFixed {
		// tag的最新视频
		if err = s.dao.AddTagNewArcCache(c, newArc, tids...); err != nil {
			log.Error("s.dao.AddTagNewArcCache(%v,%v) error(%v)", newArc, tids, err)
		}
		// 二级分区热门tag的最新视频
		if err = s.dao.AddRegionTagNewArcCache(c, newArc, tids...); err != nil {
			log.Error("s.dao.AddRegionTagNewArcCache(%v,%v) error(%v)", newArc, tids, err)
		}
		// 一级分区下tag的最新视频
		if prid, ok := s.rpMap[int64(newArc.TypeID)]; ok {
			if err = s.dao.AddTagPridArcCache(c, newArc, prid, tids); err != nil {
				log.Error("s.dao.AddTagPridArcCache(%v,%d,%v) error(%v)", newArc, prid, tids, err)
			}
		}
		log.Info("aid(%d) success", newArc.Aid)
	}
	return
}

func (s *Service) parseAllMsg(msg *model.Message) (new, old *model.Archive, err error) {
	new = &model.Archive{}
	old = &model.Archive{}
	if err = json.Unmarshal(msg.New, new); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", msg.New, err)
		return
	}
	if err = json.Unmarshal(msg.Old, old); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", msg.New, err)
	}
	return
}

func (s *Service) parseNewMsg(msg *model.Message) (new *model.Archive, err error) {
	new = &model.Archive{}
	if err = json.Unmarshal(msg.New, new); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", msg.New, err)
	}
	return
}

// resTagBind .
func (s *Service) resTagBind(c context.Context, rt *model.ResTag) (err error) {
	_, err = s.dao.InsertResTags(c, rt)
	s.cacheCh.Save(func() {
		s.dao.AddTagResCache(context.Background(), rt)
		s.dao.DelResOidCache(context.Background(), rt.Oid, rt.Type)
		s.dao.DelTagResourceCache(context.Background(), rt.Oid, rt.Type)
	})
	return
}

// resTagUntied .
func (s *Service) resTagDelete(c context.Context, rt *model.ResTag) (err error) {
	_, err = s.dao.UpdateResTags(c, rt)
	s.cacheCh.Save(func() {
		s.dao.RemoveTagResCache(context.Background(), rt)
		s.dao.DelResOidCache(context.Background(), rt.Oid, rt.Type)
		s.dao.DelTagResourceCache(context.Background(), rt.Oid, rt.Type)
	})
	return
}
