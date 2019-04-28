package service

import (
	"go-common/app/admin/main/tv/model"
	bm "go-common/library/net/http/blademaster"
)

const (
	_isDeleted = 1
)

//SetSearInterRank set search intervene rank
func (s *Service) SetSearInterRank(c *bm.Context, rank []*model.OutSearchInter) (err error) {
	err = s.dao.SetSearchInterv(c, rank)
	return
}

//GetSearInterRank get search intervene rank
func (s *Service) GetSearInterRank(c *bm.Context) (rank []*model.OutSearchInter, err error) {
	rank, err = s.dao.GetSearchInterv(c)
	return
}

//GetSearInterList get search intervene list
func (s *Service) GetSearInterList(c *bm.Context, pn, ps int) (items []*model.SearInter, total int, err error) {
	//rank, err = s.dao.GetSearchInterv(c)
	//return
	start := (pn - 1) * ps
	db := s.DB.Where("deleted!=?", _isDeleted).Order("rank ASC")
	if err = db.Model(&model.SearInter{}).Offset(start).Limit(ps).Find(&items).Error; err != nil {
		return
	}
	s.DB.Model(&model.SearInter{}).Where("deleted!=?", _isDeleted).Count(&total)
	return
}

//GetSearInterCount get search intervene count
func (s *Service) GetSearInterCount(c *bm.Context) (total int, err error) {
	if err = s.DB.Model(&model.SearInter{}).Where("deleted!=?", _isDeleted).Count(&total).Error; err != nil {
		return
	}
	return
}

//AddSearInter add search intervene
func (s *Service) AddSearInter(c *bm.Context, si *model.SearInter) (err error) {
	if err = s.DB.Create(si).Error; err != nil {
		return
	}
	return
}

//UpdateSearInter update search intervene
func (s *Service) UpdateSearInter(c *bm.Context, id int64, searchword string) (err error) {
	if err = s.DB.Model(&model.SearInter{}).Where("id=?", id).Update("searchword", searchword).Error; err != nil {
		return
	}
	return
}

//DelSearInter delete search intervene
func (s *Service) DelSearInter(c *bm.Context, id int64) (err error) {
	if err = s.DB.Model(&model.SearInter{}).Where("id=?", id).Update("deleted", 1).Error; err != nil {
		return
	}
	return
}

//RankSearInter set search intervene new rank
func (s *Service) RankSearInter(c *bm.Context, idsArr []string) (err error) {
	tx := s.DB.Begin()
	for k, v := range idsArr {
		newRank := k + 1
		id := v
		if errDB := s.DB.Model(&model.SearInter{}).Where("id=?", id).Update("rank", newRank).Error; errDB != nil {
			err = errDB
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

//GetSearInterPublish get search intervene publish status
func (s *Service) GetSearInterPublish(c *bm.Context) (items []*model.SearInter, err error) {
	limit := s.c.Cfg.SearInterMax
	db := s.DB.Where("deleted!=?", _isDeleted).Order("rank ASC")
	if err = db.Model(&model.SearInter{}).Limit(limit).Find(&items).Error; err != nil {
		return
	}
	return
}

//GetMaxRank get search intervene max rank
func (s *Service) GetMaxRank(c *bm.Context) (items model.SearInter, err error) {
	db := s.DB.Where("deleted!=?", _isDeleted).Order("rank DESC")
	if err = db.Model(&model.SearInter{}).Limit(1).Find(&items).Error; err != nil {
		return
	}
	return
}

//SetPublishState set publish status
func (s *Service) SetPublishState(c *bm.Context, state *model.PublishStatus) (err error) {
	err = s.dao.SetPublishCache(c, state)
	return
}

//GetPublishState get publish status
func (s *Service) GetPublishState(c *bm.Context) (state *model.PublishStatus, err error) {
	state, err = s.dao.GetPublishCache(c)
	return
}
