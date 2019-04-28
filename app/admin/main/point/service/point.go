package service

import (
	"context"

	"github.com/pkg/errors"

	"go-common/app/admin/main/point/model"
	pointmol "go-common/app/service/main/point/model"
	"go-common/library/ecode"
)

// PointConfList .
func (s *Service) PointConfList(c context.Context) (res []*model.PointConf, err error) {
	res, err = s.dao.PointConfList(c)
	if res == nil {
		return
	}
	for _, v := range res {
		v.Name = s.appMap[v.AppID]
	}
	return
}

// PointCoinInfo .
func (s *Service) PointCoinInfo(c context.Context, id int64) (res *model.PointConf, err error) {
	res, err = s.dao.PointCoinInfo(c, id)
	if res == nil {
		return
	}
	res.Name = s.appMap[res.AppID]
	return
}

// PointCoinAdd .
func (s *Service) PointCoinAdd(c context.Context, pc *model.PointConf) (id int64, err error) {
	if _, ok := s.appMap[pc.AppID]; !ok {
		err = ecode.RequestErr
		return
	}
	return s.dao.PointCoinAdd(c, pc)
}

// PointCoinEdit .
func (s *Service) PointCoinEdit(c context.Context, pc *model.PointConf) (err error) {
	var (
		info *model.PointConf
	)
	if info, err = s.PointCoinInfo(c, pc.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if info == nil {
		err = ecode.RequestErr
		return
	}
	if _, err = s.dao.PointCoinEdit(c, pc); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// PointHistory .
func (s *Service) PointHistory(c context.Context, arg *model.ArgPointHistory) (res *model.PageInfo, err error) {
	var sr *model.SearchData
	if sr, err = s.dao.PointHistory(c, arg); err != nil {
		return
	}
	res = &model.PageInfo{
		Item:        sr.Result,
		Count:       sr.Page.Total,
		CurrentPage: sr.Page.PN,
	}
	return
}

// PointAdd point add.
func (s *Service) PointAdd(c context.Context, pc *pointmol.ArgPoint) (err error) {
	var (
		status int8
	)
	if status, err = s.pointRPC.AddPoint(c, pc); err != nil {
		err = errors.Wrapf(err, "s.pointRPC.AddPoint(%v)", pc)
		return
	}
	if status != model.PointAddSuc {
		err = errors.New("point add not success")
	}
	return
}
