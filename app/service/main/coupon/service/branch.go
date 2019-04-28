package service

import (
	"context"

	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// UpdateBranch update branch info.
func (s *Service) UpdateBranch(c context.Context, tx *sql.Tx, count int, bi *model.CouponBatchInfo) (err error) {
	var aff int64
	if bi.MaxCount == _maxCount {
		if aff, err = s.dao.UpdateBatchInfo(c, tx, bi.BatchToken, count); err != nil {
			err = errors.WithStack(err)
			return
		}
	} else {
		if aff, err = s.dao.UpdateBatchLimitInfo(c, tx, bi.BatchToken, count); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if aff != 1 {
		err = ecode.CouPonGrantMaxCountErr
		return
	}
	if bi.MaxCount != _maxCount {
		s.dao.IncreaseBranchCurrentCountCache(c, bi.BatchToken, uint64(count))
	}
	return
}

// CurrentCount branch current count.
func (s *Service) CurrentCount(c context.Context, token string) (count int, err error) {
	var (
		mc = true
		b  *model.CouponBatchInfo
	)
	if count, err = s.dao.BranchCurrentCountCache(c, token); err != nil {
		log.Error("token(%s) err(%+v)", token, err)
		err = nil
		mc = false
	}
	if count >= 0 {
		return
	}
	if b, err = s.dao.BatchInfo(c, token); err != nil {
		err = errors.WithStack(err)
		return
	}
	if b == nil {
		err = ecode.CouPonBatchNotExistErr
		return
	}
	count = int(b.CurrentCount)
	if mc {
		s.dao.SetBranchCurrentCountCache(c, token, count)
	}
	return
}
