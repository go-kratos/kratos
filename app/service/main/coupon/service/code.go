package service

import (
	"context"
	"time"

	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// CaptchaToken get captcha token.
func (s *Service) CaptchaToken(c context.Context, ip string) (res *model.Token, err error) {
	return s.dao.CaptchaToken(c, s.c.Property.CaptchaBID, ip)
}

// UseCouponCode use coupon code.
func (s *Service) UseCouponCode(c context.Context, a *model.ArgUseCouponCode) (res *model.UseCouponCodeResp, err error) {
	var (
		code        *model.CouponCode
		bi          *model.CouponBatchInfo
		count       int64
		couponToken string
	)
	if err = s.dao.CaptchaVerify(c, a.Verify, a.Token, a.IP); err != nil {
		err = ecode.CouponCodeVerifyFaildErr
		return
	}
	if code, err = s.dao.CouponCode(c, a.Code); err != nil {
		return
	}
	if code == nil {
		err = ecode.CouponCodeNotFoundErr
		return
	}
	if code.State == model.CodeStateUsed {
		err = ecode.CouponCodeUsedErr
		return
	}
	if code.State == model.CodeStateBlock {
		err = ecode.CouponCodeBlockErr
		return
	}
	if bi = s.allBranchInfo[code.BatchToken]; bi == nil {
		err = ecode.CouPonBatchNotExistErr
		return
	}
	if bi.State == model.BatchStateBlock {
		err = ecode.CouponCodeBlockErr
		return
	}
	if bi.ExpireDay == -1 && bi.ExpireTime < time.Now().Unix() {
		err = ecode.CouponBatchExpireTimeErr
		return
	}
	if bi.LimitCount != -1 {
		if count, err = s.dao.CountCodeByMid(c, a.Mid, bi.BatchToken); err != nil {
			return
		}
		if bi.LimitCount <= count {
			err = ecode.CouponCodeLimitByMidErr
			return
		}
	}
	if bi.MaxCount != -1 && bi.MaxCount <= bi.CurrentCount {
		err = ecode.CouponCodeMaxLimitByMidErr
		return
	}
	if couponToken, err = s.updateCode(c, code, bi, a.Mid); err != nil {
		return
	}
	res = &model.UseCouponCodeResp{
		CouponToken:          couponToken,
		CouponAmount:         bi.Amount,
		FullAmount:           bi.FullAmount,
		ProductLimitMonth:    int32(bi.ProdLimMonth),
		ProductLimitRenewal:  int32(bi.ProdLimRenewal),
		PlatfromLimitExplain: model.PlatfromLimitExplain(bi.PlatformLimit, s.c.Platform),
	}
	return
}

// UpdateCode update code.
func (s *Service) updateCode(c context.Context, code *model.CouponCode, bi *model.CouponBatchInfo, mid int64) (couponToken string, err error) {
	var (
		tx  *sql.Tx
		aff int64
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = errors.Wrapf(err, "s.dao.BeginTran(%d)", mid)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				err = errors.Wrapf(err, "tx.Rollback(%d)", mid)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			err = errors.Wrapf(err, "tx.Commit(%d)", mid)
		}
		s.dao.DelCouponAllowancesKey(c, mid, model.NotUsed)
	}()
	if err = s.UpdateBranch(c, tx, 1, bi); err != nil {
		return
	}
	info := new(model.CouponAllowanceInfo)
	info.BatchToken = bi.BatchToken
	info.Mid = mid
	info.CouponToken = s.token()
	info.State = model.NotUsed
	info.FullAmount = bi.FullAmount
	info.AppID = bi.AppID
	info.Amount = bi.Amount
	info.Origin = model.AllowanceCodeOpen
	if bi.ExpireDay >= 0 {
		now, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
		info.StartTime = time.Now().Unix()
		info.ExpireTime = now.AddDate(0, 0, int(bi.ExpireDay+1)).Add(-1 * time.Second).Unix()
	} else {
		info.StartTime = bi.StartTime
		info.ExpireTime = bi.ExpireTime
	}
	if err = s.dao.TxAddAllowanceCoupon(tx, info); err != nil {
		return
	}
	if aff, err = s.dao.TxUpdateCodeState(tx, &model.CouponCode{
		State:       model.CodeStateUsed,
		Mid:         mid,
		CouponToken: info.CouponToken,
		Code:        code.Code,
		Ver:         code.Ver,
	}); err != nil {
		return
	}
	if aff != 1 {
		return "", ecode.CouponCodeCanNotUseErr
	}
	couponToken = info.CouponToken
	return
}
