package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/coupon/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// AddAllowanceBatchInfo add allowance batch info.
func (s *Service) AddAllowanceBatchInfo(c context.Context, b *model.CouponBatchInfo) (token string, err error) {
	if b.ExpireDay < 0 && (b.StartTime <= 0 && b.ExpireTime <= 0) {
		err = ecode.CouponExpireErr
		return
	}
	if b.ExpireDay < 0 && b.StartTime >= b.ExpireTime {
		err = ecode.CouPonBatchTimeErr
		return
	}
	if b.Amount >= b.FullAmount {
		err = ecode.CouponAmountErr
		return
	}
	b.BatchToken = s.token()
	b.Ctime = xtime.Time(time.Now().Unix())
	if _, err = s.dao.AddAllowanceBatchInfo(c, b); err != nil {
		err = errors.WithStack(err)
	}
	token = b.BatchToken
	return
}

//UpdateAllowanceBatchInfo update allowance batch info.
func (s *Service) UpdateAllowanceBatchInfo(c context.Context, b *model.CouponBatchInfo) (err error) {
	if _, err = s.dao.UpdateAllowanceBatchInfo(c, b); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//UpdateCodeBatchInfo update code batch info.
func (s *Service) UpdateCodeBatchInfo(c context.Context, b *model.CouponBatchInfo) (err error) {
	var data *model.CouponBatchInfo
	if data, err = s.BatchInfoByID(c, b.ID); err != nil {
		return
	}
	if data == nil || batchState(data) != model.CodeBatchUsable {
		return ecode.RequestErr
	}
	if _, err = s.dao.UpdateCodeBatchInfo(c, b); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//UpdateBatchStatus update batch status.
func (s *Service) UpdateBatchStatus(c context.Context, status int8, operator string, id int64) (err error) {
	var data *model.CouponBatchInfo
	if data, err = s.BatchInfoByID(c, id); err != nil {
		return
	}
	if data == nil ||
		(status == model.CodeBatchBlock && batchState(data) != model.CodeBatchUsable) ||
		(status == model.CodeBatchUsable && batchState(data) != model.CodeBatchBlock) {
		return ecode.RequestErr
	}
	if _, err = s.dao.UpdateBatchStatus(c, status, operator, id); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//BatchInfo allowance batch info.
func (s *Service) BatchInfo(c context.Context, token string) (res *model.CouponBatchInfo, err error) {
	if res, err = s.dao.BatchInfo(c, token); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//BatchInfoByID allowance batch info.
func (s *Service) BatchInfoByID(c context.Context, id int64) (*model.CouponBatchInfo, error) {
	return s.dao.BatchInfoByID(c, id)
}

// AllowanceSalary allowance salary.
func (s *Service) AllowanceSalary(c context.Context, f multipart.File, h *multipart.FileHeader, mids []int64, token, msgType string) (count int, err error) {
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
	if bi.ExpireDay < 0 && bi.ExpireTime < time.Now().Unix() {
		err = ecode.CouponBatchExpireTimeErr
		return
	}
	if bi.MaxCount != _notLimitSalary && len(mids) > int(bi.MaxCount-bi.CurrentCount) {
		err = ecode.CouponBatchSalaryLimitErr
		return
	}
	s.RunSalaryCoupon(c, mids, token, bi.AppID, model.CouponAllowance, model.AdminSalaryOrigin, msgType)
	count = len(mids)
	return
}

// ReadCsv read csv file
func (s *Service) ReadCsv(f multipart.File, h *multipart.FileHeader) (mids []int64, err error) {
	var (
		mid     int64
		records [][]string
	)
	mids = []int64{}
	defer f.Close()
	if h != nil && !strings.HasSuffix(h.Filename, ".csv") {
		err = ecode.CouponUpdateFileErr
		return
	}
	r := csv.NewReader(f)
	records, err = r.ReadAll()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range records {
		if len(v) <= 0 {
			continue
		}
		if mid, err = strconv.ParseInt(v[0], 10, 64); err != nil {
			err = errors.WithStack(err)
			break
		}
		mids = append(mids, mid)
	}
	return
}

// //AllowancePage allowance page.
// func (s *Service) AllowancePage(c context.Context, arg *model.ArgAllowanceSearch) (res *model.PageCouponInfo, err error) {
// 	var page *model.SearchData
// 	res = &model.PageCouponInfo{}
// 	if page, err = s.dao.AllowancePage(c, arg); err != nil {
// 		err = errors.WithStack(err)
// 		return
// 	}
// 	if page != nil && page.Data != nil && page.Data.Page != nil {
// 		res.Count = page.Data.Page.Total
// 		res.Item = page.Data.Result
// 	}
// 	return
// }

//AllowanceList allowance list.
func (s *Service) AllowanceList(c context.Context, arg *model.ArgAllowanceSearch) (res []*model.CouponAllowanceInfo, err error) {
	if res, err = s.dao.AllowanceList(c, arg); err != nil {
		err = errors.WithStack(err)
	}
	for _, v := range res {
		if v.State == model.NotUsed && v.ExpireTime < time.Now().Unix() {
			v.State = model.Expire
		}
	}
	return
}

//UpdateAllowanceState update allowance state.
func (s *Service) UpdateAllowanceState(c context.Context, mid int64, state int8, token string) (err error) {
	var (
		cp         *model.CouponAllowanceInfo
		changeType int8
	)
	if cp, err = s.dao.AllowanceByToken(c, mid, token); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cp == nil {
		err = ecode.CouPonTokenNotFoundErr
		return
	}
	if state == model.Block {
		changeType = model.AllowanceBlock
	} else {
		changeType = model.AllowanceUnBlock
	}
	if err = s.UpdateAllowanceCoupon(c, mid, state, token, cp.Ver, changeType, cp); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UpdateAllowanceCoupon update coupon info.
func (s *Service) UpdateAllowanceCoupon(c context.Context, mid int64, state int8, token string, ver int64, changeType int8, cp *model.CouponAllowanceInfo) (err error) {
	var (
		tx  *sql.Tx
		aff int64
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("%+v", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback %+v", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit %+v", err)
		}
	}()
	if aff, err = s.dao.UpdateAllowanceStatus(c, tx, state, mid, token, ver); err != nil {
		err = errors.WithStack(err)
	}
	if aff != 1 {
		err = fmt.Errorf("coupon update faild")
		return
	}
	l := &model.CouponAllowanceChangeLog{
		CouponToken: cp.CouponToken,
		Mid:         cp.Mid,
		State:       state,
		Ctime:       xtime.Time(time.Now().Unix()),
		OrderNO:     cp.OrderNO,
		ChangeType:  changeType,
	}
	if aff, err = s.dao.InsertCouponAllowanceHistory(c, tx, l); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aff != 1 {
		err = fmt.Errorf("add change log faild")
		return
	}
	s.dao.DelCouponAllowancesKey(c, cp.Mid, model.NotUsed)
	s.dao.DelCouponAllowancesKey(c, cp.Mid, model.InUse)
	return
}
