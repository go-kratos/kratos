package service

import (
	"context"
	"mime/multipart"
	"time"

	"go-common/app/admin/main/coupon/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

//CouponViewBatchAdd view batch add.
func (s *Service) CouponViewBatchAdd(c context.Context, arg *model.ArgCouponViewBatch) (err error) {
	if arg.StartTime > arg.ExpireTime {
		err = ecode.CouPonBatchTimeErr
		return
	}
	arg.BatchToken = s.token()
	arg.CouponType = model.CouponVideo
	if err = s.dao.AddViewBatch(c, arg); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//CouponViewbatchSave view batch save.
func (s *Service) CouponViewbatchSave(c context.Context, arg *model.ArgCouponViewBatch) (err error) {
	if err = s.dao.UpdateViewBatch(c, arg); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//CouponViewBlock view block
func (s *Service) CouponViewBlock(c context.Context, mid int64, couponToken string) (err error) {
	var (
		r  *model.CouponInfo
		tx *sql.Tx
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
		s.dao.DelCouponTypeCache(context.Background(), mid, model.CouponVideo)
	}()
	if r, err = s.dao.CouponViewInfo(c, couponToken, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if r.State != model.NotUsed {
		err = ecode.CouponInfoStateBlockErr
		return
	}
	if err = s.dao.TxUpdateViewInfo(tx, model.Block, couponToken, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	olog := new(model.CouponChangeLog)
	olog.Mid = mid
	olog.CouponToken = couponToken
	olog.State = model.Block
	if err = s.dao.TxCouponViewLog(tx, olog); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//CouponViewUnblock view unblock
func (s *Service) CouponViewUnblock(c context.Context, mid int64, couponToken string) (err error) {
	var (
		r  *model.CouponInfo
		tx *sql.Tx
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
		s.dao.DelCouponTypeCache(context.Background(), mid, model.CouponVideo)
	}()
	if r, err = s.dao.CouponViewInfo(c, couponToken, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if r.State != model.Block {
		err = ecode.CouponInfoStateUnblockErr
		return
	}
	if err = s.dao.TxUpdateViewInfo(tx, model.NotUsed, couponToken, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	olog := new(model.CouponChangeLog)
	olog.Mid = mid
	olog.CouponToken = couponToken
	olog.State = model.NotUsed
	if err = s.dao.TxCouponViewLog(tx, olog); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//CouponViewList view list.
func (s *Service) CouponViewList(c context.Context, arg *model.ArgSearchCouponView) (res []*model.CouponInfo, count int64, err error) {
	var (
		bls []*model.CouponBatchInfo
	)
	if arg.AppID > 0 || len(arg.BatchToken) > 0 {
		if bls, err = s.dao.BatchViewList(c, arg.AppID, arg.BatchToken, model.CouponVideo); err != nil {
			err = errors.WithStack(err)
			return
		}
		if len(bls) == 0 {
			return
		}
		for _, v := range bls {
			arg.BatchTokens = append(arg.BatchTokens, v.BatchToken)
		}
	}
	if count, err = s.dao.SearchViewCouponCount(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	if count <= 0 {
		return
	}
	if res, err = s.dao.SearchViewCouponInfo(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	now := time.Now()
	for _, v := range res {
		if v.OID == 0 {
			v.Mtime = 0
		} else {
			var pgc *model.PGCInfoResq
			if pgc, err = s.dao.GetPGCInfo(c, v.OID); err != nil {
				log.Error("get pgc info coupon%+v  error:%+v", v, err)
				continue
			}
			if pgc != nil {
				v.Title = pgc.Title
			}
		}
		if len(v.BatchToken) > 0 {
			var cbi *model.CouponBatchInfo
			if cbi, err = s.dao.BatchInfo(c, v.BatchToken); err != nil {
				err = errors.WithStack(err)
				return
			}
			v.BatchName = cbi.Name
		}
		if v.ExpireTime < now.Unix() && v.State == model.NotUsed {
			v.State = model.Expire
		}
	}
	return
}

//CouponViewSalary view salary.
func (s *Service) CouponViewSalary(c context.Context, f multipart.File, h *multipart.FileHeader, mids []int64, token string) (count int, err error) {
	var bi *model.CouponBatchInfo
	if len(mids) == 0 {
		if mids, err = s.ReadCsv(f, h); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if len(mids) == 0 {
		err = ecode.CouponBatchSalaryCountZeroErr
		return
	}
	if len(mids) > _maxSalaryCount {
		err = ecode.CouponBatchSalaryLimitErr
		return
	}
	if bi, err = s.dao.BatchInfo(c, token); err != nil {
		err = errors.WithStack(err)
		return
	}
	if bi == nil {
		err = ecode.CouPonTokenNotFoundErr
		return
	}
	if bi.State != model.BatchStateNormal {
		err = ecode.CouPonHadBlockErr
		return
	}
	if bi.MaxCount != _notLimitSalary && len(mids) > int(bi.MaxCount-bi.CurrentCount) {
		err = ecode.CouponBatchSalaryLimitErr
		return
	}
	s.RunSalaryCoupon(c, mids, token, bi.AppID, model.CouponVideo, model.SystemAdminSalary, "")
	count = len(mids)
	return
}
