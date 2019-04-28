package service

import (
	"context"

	"go-common/app/service/main/passport/model"
)

var (
	_emptyFaceRecords = make([]*model.FaceApply, 0)
)

// FaceApplies get face applies range from and to.
func (s *Service) FaceApplies(c context.Context, mid int64, from, to int64, status, operator string) (res []*model.FaceApply, err error) {
	if from > to {
		return _emptyFaceRecords, nil
	}
	if res, err = s.d.FaceApplies(c, mid, from, to, status, operator); err != nil {
		return
	}
	if len(res) == 0 {
		res = _emptyFaceRecords
	}
	return
}
