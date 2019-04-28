package service

import (
	"context"

	"go-common/app/service/main/usersuit/model"
)

// PointFlag .
func (s *Service) PointFlag(c context.Context, arg *model.ArgMID) (pf *model.PointFlag, err error) {
	var pid, nid int64
	if pid, err = s.pendantDao.RedPointCache(c, arg.MID); err != nil {
		return
	}
	if nid, err = s.medalDao.RedPointCache(c, arg.MID); err != nil {
		return
	}
	pf = &model.PointFlag{
		Pendant: pid > 0,
		Medal:   nid > 0,
	}
	return
}
