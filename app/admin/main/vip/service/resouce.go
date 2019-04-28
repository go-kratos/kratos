package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"time"

	"go-common/app/admin/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// QueryPool query pool
func (s *Service) QueryPool(c context.Context, r *model.ResoucePoolBo) (res []*model.VipResourcePool, count int, err error) {
	var (
		batch *model.VipResourceBatch
		PN    int
		PS    int
	)
	PN = r.PN
	PS = r.PS

	if batch, err = s.dao.SelBatchRow(c, r.BatchID); err != nil {
		return
	}
	if r.BatchID != 0 {
		r.ID = -1
	}

	if batch != nil {
		r.ID = batch.PoolID
	}

	if count, err = s.dao.SelCountPool(c, r); err != nil || count == 0 {
		return
	}
	if res, err = s.dao.SelPool(c, r, PN, PS); err != nil {
		return
	}
	return
}

// PoolInfo pool info
func (s *Service) PoolInfo(c context.Context, id int) (res *model.VipResourcePool, err error) {
	if res, err = s.dao.SelPoolRow(c, id); err != nil {
		return
	}
	return
}

// AddPool .
func (s *Service) AddPool(c context.Context, r *model.ResoucePoolBo) (err error) {
	var (
		p *model.VipResourcePool
	)
	if err = s.verificationPool(r); err != nil {
		return
	}
	if p, err = s.dao.SelPoolByName(c, r.PoolName); err != nil {
		return
	}
	if p != nil {
		err = ecode.VipPoolNameExitErr
		return
	}
	if _, err = s.dao.AddPool(c, r); err != nil {
		return
	}
	return
}

func (s *Service) verificationPool(r *model.ResoucePoolBo) (err error) {
	var (
		business *model.VipBusinessInfo
	)
	if len(r.PoolName) == 0 {
		err = ecode.VipPoolNameErr
		return
	}
	if len(r.Reason) == 0 {
		err = ecode.VipPoolReasonErr
		return
	}
	if r.StartTime <= 0 {
		err = ecode.VipPoolStartTimeErr
		return
	}
	if r.EndTime <= 0 {
		err = ecode.VipPoolEndTimeErr
		return
	}
	if r.EndTime < r.StartTime {
		err = ecode.VipPoolValidityTimeErr
		return
	}

	if business, err = s.dao.SelBusiness(context.TODO(), r.BusinessID); err != nil {
		return
	}
	if business == nil {
		err = ecode.VipBusinessNotExitErr
		return
	}
	return
}

// UpdatePool .
func (s *Service) UpdatePool(c context.Context, r *model.ResoucePoolBo) (err error) {
	var (
		pool   *model.VipResourcePool
		p      *model.VipResourcePool
		batchs []*model.VipResourceBatch
	)
	if err = s.verificationPool(r); err != nil {
		return
	}
	if pool, err = s.dao.SelPoolRow(c, r.ID); err != nil {
		return
	}
	if pool == nil {
		err = ecode.VipPoolIDErr
		return
	}
	if p, err = s.dao.SelPoolByName(c, r.PoolName); err != nil {
		return
	}
	if p != nil && pool.PoolName != p.PoolName {
		err = ecode.VipPoolNameExitErr
		return
	}
	if batchs, err = s.dao.SelBatchRows(context.TODO(), r.ID); err != nil {
		return
	}
	for _, v := range batchs {
		if !(r.StartTime <= v.StartTime && r.EndTime >= v.EndTime) {
			err = ecode.VipPoolValidityTimeErr
			return
		}
	}
	if _, err = s.dao.UpdatePool(c, r); err != nil {
		return
	}
	return
}

// BatchInfo .
func (s *Service) BatchInfo(c context.Context, id int) (res *model.VipResourceBatch, err error) {
	if res, err = s.dao.SelBatchRow(c, id); err != nil {
		return
	}
	return
}

// BatchInfoOfPool .
func (s *Service) BatchInfoOfPool(c context.Context, poolID int) (res []*model.VipResourceBatch, err error) {
	if res, err = s.dao.SelBatchRows(c, poolID); err != nil {
		return
	}
	return
}

// AddBatch .
func (s *Service) AddBatch(c context.Context, r *model.ResouceBatchBo) (err error) {
	if err = s.verificationBatch(r); err != nil {
		return
	}
	r.SurplusCount = r.Count
	if _, err = s.dao.AddBatch(c, r); err != nil {
		return
	}
	return
}

// UpdateBatch .
func (s *Service) UpdateBatch(c context.Context, id, increment int, startTime, endTime xtime.Time) (err error) {
	var (
		batch *model.VipResourceBatch
		r     = new(model.ResouceBatchBo)
	)
	r.ID = id
	r.StartTime = startTime
	r.EndTime = endTime

	if batch, err = s.dao.SelBatchRow(c, id); err != nil {
		return
	}
	if batch == nil {
		err = ecode.VipBatchIDErr
		return
	}
	r.PoolID = batch.PoolID
	if err = s.verifBatchTime(r); err != nil {
		return
	}
	if increment < 0 {
		err = ecode.VipBatchPlusResouceErr
		return
	}
	if batch.Count+increment > math.MaxInt32 || batch.Count+increment < 0 {
		err = ecode.VipBatchCountErr
		return
	}
	batch.Count += increment
	batch.SurplusCount += increment
	batch.StartTime = r.StartTime
	batch.EndTime = r.EndTime
	ver := batch.Ver
	batch.Ver++
	if _, err = s.dao.UpdateBatch(c, batch, ver); err != nil {
		return
	}
	return
}

func (s *Service) verifBatchTime(r *model.ResouceBatchBo) (err error) {
	var (
		pool *model.VipResourcePool
	)
	if pool, err = s.dao.SelPoolRow(context.TODO(), r.PoolID); err != nil {
		return
	}
	if pool == nil {
		err = ecode.VipPoolIDErr
		return
	}

	if pool.StartTime > r.StartTime || pool.EndTime < r.EndTime {
		err = ecode.VipPoolValidityTimeErr
		return
	}
	return
}
func (s *Service) verificationBatch(r *model.ResouceBatchBo) (err error) {

	if r.Unit <= 0 || r.Unit > 3660 {
		err = ecode.VipBatchUnitErr
		return
	}
	if r.Count <= 0 {
		err = ecode.VipBatchCountErr
		return
	}
	if err = s.verifBatchTime(r); err != nil {
		return
	}
	return

}

// GrandResouce grand resouce mid
func (s *Service) GrandResouce(c context.Context, remark string, batchID int64, mids []int, username string) (failMid []int, err error) {
	var (
		batch *model.VipResourceBatch
	)
	if len(remark) == 0 {
		err = ecode.VipRemarkErr
		return
	}

	if batch, err = s.dao.SelBatchRow(c, int(batchID)); err != nil {
		return
	}
	if batch == nil {
		err = ecode.VipBatchIDErr
		return
	}
	if err = s.checkBatchValid(batch); err != nil {
		return
	}
	for _, v := range mids {
		if err = s.grandMidOfResouce(c, v, int(batchID), username, remark); err != nil {
			log.Error("GrandResouce grandMidOfResouce(mid:%v,batchID:%v,username:%v,remark:%v error(%v))", v, batchID, username, remark, err)
			failMid = append(failMid, v)

		}
	}
	return
}

func (s *Service) grandMidOfResouce(c context.Context, mid, batchID int, username, remark string) (err error) {
	//var (
	//	batch *model.VipResourceBatch
	//	tx    *sql.Tx
	//	a     int64
	//	hv    *inModel.HandlerVip
	//)
	//if batch, err = s.dao.SelBatchRow(c, batchID); err != nil {
	//	return
	//}
	//if batch.SurplusCount-1 < 0 {
	//	err = ecode.VipBatchNotEnoughErr
	//	return
	//}
	//batch.DirectUseCount++
	//batch.SurplusCount--
	//ver := batch.Ver
	//batch.Ver++
	//if tx, err = s.dao.BeginTran(context.TODO()); err != nil {
	//	return
	//}
	//defer func() {
	//	if err != nil {
	//		if err = tx.Commit(); err != nil {
	//			tx.Rollback()
	//		}
	//	} else {
	//		tx.Rollback()
	//	}
	//}()
	//if a, err = s.dao.UseBatch(tx, batch, ver); err != nil {
	//	return
	//}
	//if a > 0 {
	//	if hv, err = s.exchangeVip(context.TODO(), tx, int(mid), batch.ID, batch.Unit, remark, username); err != nil {
	//		return
	//	}
	//	s.asyncBcoin(func() {
	//		s.vipRPC.BcoinProcesserHandler(context.TODO(), hv)
	//	})
	//}
	return
}

func (s *Service) checkBatchValid(batch *model.VipResourceBatch) (err error) {
	var (
		pool     *model.VipResourcePool
		business *model.VipBusinessInfo
		ct       = time.Now()
	)
	if !(batch.StartTime.Time().Unix() <= ct.Unix() && ct.Unix() <= batch.EndTime.Time().Unix()) {
		err = ecode.VipBatchTTLErr
		return
	}
	if pool, err = s.dao.SelPoolRow(context.TODO(), batch.PoolID); err != nil {
		return
	}
	if pool == nil {
		err = ecode.VipPoolIDErr
		return
	}
	if !(pool.StartTime.Time().Unix() <= ct.Unix() && ct.Unix() <= pool.EndTime.Time().Unix()) {
		err = ecode.VipPoolValidityTimeErr
		return
	}
	if business, err = s.dao.SelBusiness(context.TODO(), pool.BusinessID); err != nil {
		return
	}
	if business == nil {
		err = ecode.VipBusinessNotExitErr
		return
	}
	if business.Status == 1 {
		err = ecode.VipBusinessStatusErr
		return
	}
	return
}

// SaveBatchCode .
func (s *Service) SaveBatchCode(c context.Context, arg *model.BatchCode) (err error) {
	var batchID int64
	if arg.ID == 0 {
		if err = s.checkBatchCodeValid(c, arg); err != nil {
			return
		}
		arg.SurplusCount = arg.Count
		var tx *sql.Tx

		if tx, err = s.dao.BeginTran(c); err != nil {
			err = errors.WithStack(err)
			return
		}
		defer func() {
			if err == nil {
				if err = tx.Commit(); err != nil {
					log.Error("commimt error(%+v)", err)
				}
			} else {
				tx.Rollback()
			}
		}()

		arg.Status = 1
		if batchID, err = s.dao.TxAddBatchCode(tx, arg); err != nil {
			err = errors.WithStack(err)
			return
		}

		if err = s.createCode(tx, batchID, int(arg.Count)); err != nil {
			err = errors.WithStack(err)
			return
		}

	} else {
		var (
			bc  *model.BatchCode
			bc1 *model.BatchCode
		)
		if bc, err = s.dao.SelBatchCodeID(c, arg.ID); err != nil {
			err = errors.WithStack(err)
			return
		}
		if bc == nil {
			err = ecode.VipBatchIDErr
			return
		}
		if bc.BatchName != arg.BatchName {
			if bc1, err = s.dao.SelBatchCodeName(c, arg.BatchName); err != nil {
				err = errors.WithStack(err)
				return
			}
			if bc1 != nil {
				err = ecode.VipBatchCodeNameErr
				return
			}
		}
		bc.BatchName = arg.BatchName
		bc.Reason = arg.Reason
		bc.Price = arg.Price
		bc.Contacts = arg.Contacts
		bc.ContactsNumber = arg.ContactsNumber
		bc.Type = arg.Type
		bc.MaxCount = arg.MaxCount
		bc.LimitDay = arg.LimitDay
		if _, err = s.dao.UpdateBatchCode(c, bc); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}

func (s *Service) checkBatchCodeValid(c context.Context, arg *model.BatchCode) (err error) {
	var (
		b  *model.VipBusinessInfo
		p  *model.VipResourcePool
		bc *model.BatchCode
	)
	if b, err = s.BusinessInfo(c, int(arg.BusinessID)); err != nil {
		err = errors.WithStack(err)
		return
	}
	if b == nil {
		err = ecode.VipBusinessNotExitErr
		return
	}
	if p, err = s.PoolInfo(c, int(arg.PoolID)); err != nil {
		err = errors.WithStack(err)
		return
	}
	if p == nil {
		err = ecode.VipPoolIDErr
		return
	}

	if p.EndTime.Time().Before(arg.EndTime.Time()) || p.StartTime.Time().After(arg.StartTime.Time()) {
		err = ecode.VipPoolValidityTimeErr
		return
	}

	if bc, err = s.dao.SelBatchCodeName(c, arg.BatchName); err != nil {
		err = errors.WithStack(err)
		return
	}

	if bc != nil {
		err = ecode.VipBatchCodeNameErr
		return
	}

	if arg.Unit <= 0 || arg.Unit > 3660 {
		err = ecode.VipBatchUnitErr
		return
	}
	if arg.Count <= 0 || arg.Count > 200000 {
		err = ecode.VipBatchCodeCountErr
	}
	if arg.Price > 10000 || arg.Price < 0 {
		err = ecode.VipBatchPriceErr
		return
	}
	return
}

// FrozenCode .
func (s *Service) FrozenCode(c context.Context, codeID int64, status int8) (err error) {
	var (
		code *model.ResourceCode
	)
	if code, err = s.dao.SelCodeID(c, codeID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if code == nil {
		err = ecode.VipCodeIDErr
		return
	}
	code.Status = status
	if _, err = s.dao.UpdateCode(c, codeID, status); err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

// FrozenBatchCode .
func (s *Service) FrozenBatchCode(c context.Context, BatchCodeID int64, status int8) (err error) {
	var bc *model.BatchCode
	if bc, err = s.dao.SelBatchCodeID(c, BatchCodeID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if bc == nil {
		err = ecode.VipBatchIDErr
		return
	}
	bc.Status = status

	if _, err = s.dao.UpdateBatchCode(c, bc); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (s *Service) createCode(tx *sql.Tx, batchCodeID int64, size int) (err error) {
	var (
		hash      = md5.New()
		batchSize = 2000
		codes     = make([]*model.ResourceCode, 0)
	)

	for i := 1; i <= size; i++ {
		unix := time.Now().UnixNano()
		key := fmt.Sprintf("%v,%v,%v,%v", unix, batchCodeID, i, rand.Intn(10000000))
		hash.Write([]byte(key))
		sum := hash.Sum(nil)
		code := hex.EncodeToString(sum)
		code = code[8:24]
		r := new(model.ResourceCode)
		r.Code = code
		r.Status = model.NOTUSER
		r.BatchCodeID = batchCodeID
		codes = append(codes, r)
		if i%batchSize == 0 || i == size {
			if err = s.dao.BatchAddCode(tx, codes); err != nil {
				log.Error("batch add code %+v", err)
				return
			}
			codes = make([]*model.ResourceCode, 0)
		}
	}
	return
}

// SelBatchCodes .
func (s *Service) SelBatchCodes(c context.Context, batchIDs []int64) (res []*model.BatchCode, err error) {
	if res, err = s.dao.SelBatchCodes(c, batchIDs); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// SelCode .
func (s *Service) SelCode(c context.Context, arg *model.ArgCode, username string, curID int64, ps int) (res []*model.CodeVo, cursor int64, pre int64, err error) {
	var (
		codes      []*model.ResourceCode
		batchIDs   []int64
		batchMap   = make(map[int64]*model.BatchCode)
		batchCodes []*model.BatchCode
		linkmap    map[int64]int64
	)
	if linkmap, err = s.dao.GetSelCode(c, username); err != nil {
		err = errors.WithStack(err)
		return
	}
	fmt.Printf("cur link map(%+v) \n", linkmap)
	if len(linkmap) == 0 {
		linkmap = make(map[int64]int64)
	}
	if codes, err = s.dao.SelCode(c, arg, curID, ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(codes) > 0 {
		cursor = codes[len(codes)-1].ID
	} else {
		return
	}

	linkmap[cursor] = curID
	pre = linkmap[curID]
	if err = s.dao.SetSelCode(c, username, linkmap); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range codes {
		batchIDs = append(batchIDs, v.BatchCodeID)
	}
	if batchCodes, err = s.dao.SelBatchCodes(c, batchIDs); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range batchCodes {
		batchMap[v.ID] = v
	}
	for _, v := range codes {
		r := new(model.CodeVo)
		r.ID = v.ID
		r.BatchCodeID = v.BatchCodeID
		r.Mid = v.Mid
		r.Ctime = v.Ctime
		r.Code = v.Code
		r.Status = v.Status
		r.UseTime = v.UseTime
		batchCode := batchMap[v.BatchCodeID]
		if batchCode == nil {
			err = ecode.VipBatchIDErr
			res = nil
			return
		}
		r.Unit = batchCode.Unit
		r.BatchName = batchCode.BatchName
		r.BatchStatus = batchCode.Status
		r.StartTime = batchCode.StartTime
		r.EndTime = batchCode.EndTime
		res = append(res, r)
	}

	return
}

// ExportCode .
func (s *Service) ExportCode(c context.Context, batchID int64) (codes []string, err error) {
	var (
		rc    []*model.ResourceCode
		curID int64
		ps    = 2000
	)
	arg := new(model.ArgCode)
	arg.BatchCodeID = batchID
	arg.Status = model.NOTUSER
	for {
		if rc, err = s.dao.SelCode(c, arg, curID, ps); err != nil {
			err = errors.WithStack(err)
			return
		}
		if len(rc) == 0 {
			return
		}
		for _, v := range rc {
			codes = append(codes, v.Code)
		}
		curID = rc[len(rc)-1].ID
	}

}

// SelBatchCode .
func (s *Service) SelBatchCode(c context.Context, arg *model.ArgBatchCode, pn, ps int) (res []*model.BatchCode, total int64, err error) {
	if total, err = s.dao.SelBatchCodeCount(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res, err = s.dao.SelBatchCode(c, arg, pn, ps); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
