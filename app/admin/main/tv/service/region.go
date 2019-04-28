package service

import (
	"context"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"
)

// RegList region list .
func (s *Service) RegList(ctx context.Context, arg *model.Param) (res []*model.RegList, err error) {
	var reg []*model.RegDB
	res = make([]*model.RegList, 0)
	if reg, err = s.dao.RegList(ctx, arg); err != nil {
		log.Error("s.dao.RegList error(%v)", err)
		return
	}
	for _, v := range reg {
		res = append(res, v.ToList())
	}
	return
}

// AddReg add region .
func (s *Service) AddReg(ctx context.Context, title, itype, itid, rank string) (err error) {
	if err = s.dao.AddReg(ctx, title, itype, itid, rank); err != nil {
		log.Error("s.dao.AddReg error(%v)", err)
	}
	return
}

// EditReg edit region .
func (s *Service) EditReg(ctx context.Context, pid, title, itype, itid string) (err error) {
	if err = s.dao.EditReg(ctx, pid, title, itype, itid); err != nil {
		log.Error("s.dao.EditReg error(%v)", err)
	}
	return
}

// UpState publish or not .
func (s *Service) UpState(ctx context.Context, pids []int, state string) (err error) {
	if err = s.dao.UpState(ctx, pids, state); err != nil {
		log.Error("s.dao.UpState error(%v)", err)
	}
	return
}

// RegSort .
func (s *Service) RegSort(ctx context.Context, ids []int) (err error) {
	order := 0
	for _, v := range ids {
		if !s.isExist(v) {
			log.Error("id is not exit! id(%d) error(%v)", v, err)
			return
		}
	}
	tx := s.DB.Begin()
	for _, v := range ids {
		order += 1
		if err = tx.Model(&model.RegDB{}).Where("id=?", v).Update(map[string]int{"rank": order}).Error; err != nil {
			log.Error("RegSort Update error(%v)", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

func (s *Service) isExist(id int) (f bool) {
	var (
		err error
		a   = &model.RegDB{}
	)
	if err = s.DB.Where("id=?", id).Where("deleted=0").Find(a).Error; err != nil {
		log.Error("isExist s.DB.Where error(%s)")
		return
	}
	if a.ID != 0 {
		return true
	}
	return false
}
