package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	codesLimit = 100
)

//BatchInfo get resource Batch Info by id.
func (s *Service) BatchInfo(c context.Context, id int64, appkey string) (r *model.VipResourceBatch, bis *model.VipBusinessInfo, err error) {
	if r, err = s.dao.SelVipResourceBatch(c, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if r == nil {
		err = ecode.VipBatchIDErr
		return
	}
	t := time.Now().Unix()
	if t < r.StartTime.Time().Unix() || t > r.EndTime.Time().Unix() {
		err = ecode.VipBatchTTLErr
		return
	}
	if bis, err = s.BusinessByPool(c, r.PoolID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if bis == nil {
		err = ecode.VipBusinessErr
		return
	}
	if bis.Status == model.VipBusinessStatusClose {
		err = ecode.VipBusinessStatusErr
		return
	}
	if bis.AppKey != appkey {
		err = ecode.VipBusinessErr
	}
	return
}

//BusinessByPool get pool info by id.
func (s *Service) BusinessByPool(c context.Context, poolID int64) (r *model.VipBusinessInfo, err error) {
	var (
		pool *model.VipResourcePool
	)
	if pool, err = s.dao.SelResourcePool(c, poolID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if pool == nil {
		err = ecode.VipPoolIDErr
		return
	}
	t := time.Now().Unix()
	if pool.StartTime.Time().Unix() > t || t > pool.EndTime.Time().Unix() {
		err = ecode.VipPoolTTLErr
		return
	}
	r, err = s.dao.SelBusiness(c, pool.BusinessID)
	return
}

//ResourceBatchOpenVip .
func (s *Service) ResourceBatchOpenVip(c context.Context, arg *model.ArgUseBatch) (err error) {
	var (
		vch *model.VipChangeHistory
		bi  *model.VipResourceBatch
		bis *model.VipBusinessInfo
	)
	if vch, err = s.dao.OldVipchangeHistory(c, arg.OrderNo, arg.BatchID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if vch != nil {
		log.Info("order already handler success......")
		return
	}
	key := fmt.Sprintf("lock:o:%s,i:%d", arg.OrderNo, arg.BatchID)
	if succeed := s.dao.AddTransferLock(c, key); !succeed {
		log.Info("this order handlering,please wait arg:%+v", arg)
		return
	}
	defer func() {
		s.dao.DelCache(context.Background(), key)
	}()
	if bi, bis, err = s.BatchInfo(c, arg.BatchID, arg.Appkey); err != nil {
		err = errors.WithStack(err)
		return
	}
	if bi.SurplusCount <= 0 {
		err = ecode.VipBatchNotEnoughErr
		return
	}
	if bis.BusinessType == model.BizTypeOut {
		if err = s.checkSign(arg.ToMap(), bis.Secret); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if err = s.resourceBatchOpenVip(c, arg, bi); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (s *Service) resourceBatchOpenVip(c context.Context, arg *model.ArgUseBatch, bi *model.VipResourceBatch) (err error) {
	var (
		tx  *sql.Tx
		hv  *model.OldHandlerVip
		ip  = metadata.String(c, metadata.RemoteIP)
		eff int64
	)
	if tx, err = s.dao.OldStartTx(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
		s.dao.DelVipInfoCache(context.Background(), arg.Mid)
		s.cache(func() {
			s.oldCleanCacheAndNotify(context.Background(), hv, ip)
		})
		s.asyncBcoin(func() {
			s.OldProcesserHandler(context.Background(), hv, ip)
		})
	}()
	ver := bi.Ver
	bi.Ver = bi.Ver + 1
	bi.DirectUseCount++
	bi.SurplusCount--
	if eff, err = s.dao.UpdateBatchCount(c, tx, bi, ver); err != nil {
		return
	}
	if eff <= 0 {
		err = ecode.VipBatchNotEnoughErr
		return
	}
	r := new(model.VipChangeBo)
	r.BatchID = arg.BatchID
	r.ChangeType = model.ChangeTypeSystem
	r.Remark = arg.Remark
	r.RelationID = arg.OrderNo
	r.Days = bi.Unit
	r.Mid = arg.Mid
	if hv, err = s.OldUpdateVipWithHistory(c, tx, r); err != nil {
		return
	}
	return
}

// OpenCode open code.
func (s *Service) OpenCode(c context.Context, codeStr string, mid int64) (code *model.VipResourceCode, err error) {
	var (
		count     int
		batchCode *model.VipResourceBatchCode
		data      *model.CommonResq
		tx        *sql.Tx
		eff       int64
	)
	if count, err = s.dao.GetOpenCodeCount(c, mid); err != nil {
		log.Error("mc error(%v)", err)
		err = nil
	}
	if count > 20 {
		err = ecode.VipOpenCodeCountErr
		return
	}
	defer func() {
		count++
		if err1 := s.dao.SetOpenCode(c, mid, count); err1 != nil {
			log.Error("set open code error(%+v)", err1)
		}
	}()
	if code, err = s.dao.SelCode(c, codeStr); err != nil {
		err = errors.WithStack(err)
		return
	}
	if code == nil {
		err = ecode.VipCodeNotExitErr
		return
	}

	if code.Status == model.CodeUse {
		err = ecode.VipCodeUsedErr
		return
	}

	if code.Status == model.CodeFrozen {
		err = ecode.VipCodeFrozenErr
		return
	}
	if batchCode, err = s.checkCodeInfo(c, codeStr, code, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if tx, err = s.dao.StartTx(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err == nil || err == ecode.VipCodeUsedErr {
			err1 := err
			if err = tx.Commit(); err != nil {
				log.Error("commit error(%v)", err)
				return
			}
			if err1 == ecode.VipCodeUsedErr {
				err = ecode.VipCodeUsedErr
			}
		} else {
			tx.Rollback()
		}
	}()
	if eff, err = s.dao.TxUpdateCodeStatus(tx, code.ID, model.CodeUse); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff == 0 {
		err = ecode.VipOpenCodeErr
		return
	}
	if eff, err = s.dao.TxUpdateCode(tx, code.ID, mid, xtime.Time(time.Now().Unix())); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff == 0 {
		err = ecode.VipOpenCodeErr
		return
	}
	batchCode.SurplusCount--
	if eff, err = s.dao.TxUpdateBatchCode(tx, batchCode.ID, batchCode.SurplusCount); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff == 0 {
		err = ecode.VipOpenCodeErr
		return
	}
	code.Days = batchCode.Unit
	if data, err = s.dao.OpenCode(c, mid, code.BatchCodeID, batchCode.Unit, "激活码开通", code.Code); err != nil {
		log.Error("open code error(%+v)", err)
		var data *model.OpenCodeResp
		time.Sleep(time.Millisecond * 200)
		if data, err = s.dao.GetOpenInfo(context.TODO(), code.Code); err != nil {
			log.Error("get open info (%+v)", err)
			err = nil
			return
		}
		if data.Code == 0 && data.Data == int64(model.CodeUnUse) {
			err = ecode.VipOpenCodeErr
			return
		}
		err = nil
		return
	}
	if int(data.Code) == ecode.VipCodeUsedErr.Code() {
		err = ecode.VipCodeUsedErr
		return
	} else if data.Code != 0 {
		err = ecode.VipOpenCodeErr
		return
	}
	//if hv, err = s.openCodeVip(c, tx, mid, code.BatchCodeID, batchCode.Unit, "激活码开通", code.Code); err != nil {
	//	err = errors.WithStack(err)
	//	return
	//}
	s.cache(func() {
		s.dao.DelVipInfoCache(context.TODO(), mid)
	})
	//s.asyncBcoin(func() {
	//	s.ProcesserHandler(context.TODO(), hv, ip)
	//})

	return
}

func (s *Service) checkCodeInfo(c context.Context, codeStr string, code *model.VipResourceCode, mid int64) (batchCode *model.VipResourceBatchCode, err error) {
	var (
		pool  *model.VipResourcePool
		bis   *model.VipBusinessInfo
		vip   *model.VipInfoResp
		res   *model.PassportDetail
		count int64
	)
	if batchCode, err = s.dao.SelBatchCode(c, code.BatchCodeID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if batchCode == nil {
		err = ecode.VipCodeNotExitErr
		return
	}
	if batchCode.Status == model.BatchFrozen {
		err = ecode.VipCodeFrozenErr
		return
	}
	if batchCode.Type == model.OnlyNotVip {
		if vip, err = s.ByMid(c, mid); err != nil {
			err = errors.WithStack(err)
			return
		}
		if vip.VipStatus != model.VipStatusOverTime {
			err = ecode.VipOpenOnlyNotVipErr
			return
		}
	}
	now := time.Now()
	if batchCode.StartTime.Time().After(now) {
		err = ecode.VipCodeNotStartErr
		return
	}
	if batchCode.EndTime.Time().Before(now) {
		err = ecode.VipCodeTTLErr
		return
	}
	if pool, err = s.dao.SelNewResourcePool(c, batchCode.PoolID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if pool == nil {
		err = ecode.VipCodeNotExitErr
		return
	}
	if pool.StartTime.Time().After(now) {
		err = ecode.VipCodeNotStartErr
		return
	}
	if pool.EndTime.Time().Before(now) {
		err = ecode.VipCodeTTLErr
		return
	}
	if bis, err = s.dao.SelNewBusiness(c, batchCode.BusinessID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if bis == nil {
		err = ecode.VipCodeNotExitErr
		return
	}
	if bis.Status == model.StatusClose {
		err = ecode.VipCodeFrozenErr
		return
	}
	if batchCode.MaxCount > 0 {
		if count, err = s.dao.SelBatchCount(c, batchCode.ID, mid); err != nil {
			err = errors.WithStack(err)
			return
		}
		if count >= batchCode.MaxCount {
			err = ecode.VipBatchMaxCountErr
			return
		}
	}
	if batchCode.LimitDay >= 0 {
		if res, err = s.dao.GetPassportDetail(c, mid); err != nil {
			err = errors.WithStack(err)
			return
		}
		location, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
		jointime, _ := time.ParseInLocation("2006-01-02", time.Unix(int64(res.JoinTime), 0).Format("2006-01-02"), time.Local)
		day := location.Sub(jointime).Hours() / 24
		if int64(day) > batchCode.LimitDay {
			err = ecode.VipBatchLimitDayErr
			return
		}
	}
	return
}

//CodeInfo get code info.
func (s *Service) CodeInfo(c context.Context, codeStr string) (code *model.VipResourceCode, err error) {
	var (
		batchCode *model.VipResourceBatchCode
	)
	if code, err = s.dao.SelCode(c, codeStr); err != nil {
		err = errors.WithStack(err)
		return
	}
	if code == nil {
		return
	}
	if batchCode, err = s.dao.SelBatchCode(c, code.BatchCodeID); err != nil {
		err = errors.WithStack(err)
		return
	}
	code.Days = batchCode.Unit
	return
}

//CodeInfos get code infos.
func (s *Service) CodeInfos(c context.Context, codes []string) (cs []*model.VipResourceCode, err error) {
	var (
		batchCodeIDs []int64
		batchCodes   []*model.VipResourceBatchCode
		batchCodeMap = make(map[int64]*model.VipResourceBatchCode)
	)
	if len(codes) > codesLimit {
		err = ecode.VipCodeLimitErr
		return
	}
	if cs, err = s.dao.SelCodes(c, codes); err != nil {
		err = errors.WithStack(err)
		return
	}

	for _, v := range cs {
		batchCodeIDs = append(batchCodeIDs, v.BatchCodeID)
	}
	if batchCodes, err = s.dao.SelBatchCodes(c, batchCodeIDs); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range batchCodes {
		batchCodeMap[v.ID] = v
	}
	for _, v := range cs {
		if batchCodeMap[v.BatchCodeID] != nil {
			batchCode := batchCodeMap[v.BatchCodeID]
			v.Days = batchCode.Unit
		}
	}
	return
}

//WebToken get web token.
func (s *Service) WebToken(c context.Context) (token *model.Token, err error) {
	var (
		tokenReq *model.TokenResq
	)
	if tokenReq, err = s.dao.GetToken(c, s.c.Property.TokenBID, metadata.String(c, metadata.RemoteIP)); err != nil {
		err = errors.WithStack(err)
		return
	}
	token = tokenReq.Data
	return
}

//Belong get belong info.
func (s *Service) Belong(c context.Context, mid int64) (cs []*model.VipResourceCode, err error) {
	var (
		codes []string
	)
	if codes, err = s.dao.SelCodesByBMid(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cs, err = s.CodeInfos(c, codes); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//Verify verify code.
func (s *Service) Verify(c context.Context, token, code string) (t *model.TokenResq, err error) {
	if t, err = s.dao.Verify(c, code, token, metadata.String(c, metadata.RemoteIP)); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//Actives get active datas.
func (s *Service) Actives(c context.Context, relations []string) (res []*model.VipActiveShow, err error) {
	if res, err = s.dao.SelActives(c, relations); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//CodeOpened .
func (s *Service) CodeOpened(c context.Context, arg *model.ArgCodeOpened) (res []*model.CodeInfoResp, err error) {
	var (
		biz    *model.VipBusinessInfo
		rbc    []*model.VipResourceBatchCode
		rbcIDs []int64
	)
	if biz, err = s.dao.SelNewBusinessByAppkey(c, arg.BisAppkey); err != nil {
		err = errors.WithStack(err)
		return
	}
	if biz == nil {
		err = ecode.VipBusinessNotExitErr
		return
	}
	if biz.Status == model.StatusClose {
		err = ecode.VipBusinessStatusErr
		return
	}
	if biz.BusinessType == model.BizTypeIn {
		err = ecode.VipBisTypeErr
		return
	}
	if err = s.checkSign(arg.ToMap(), biz.Secret); err != nil {
		err = errors.WithStack(err)
		return
	}
	if rbc, err = s.dao.SelBatchCodesByBisID(c, biz.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(rbc) <= 0 {
		return
	}
	for _, v := range rbc {
		rbcIDs = append(rbcIDs, v.ID)
	}
	if res, err = s.dao.SelCodeOpened(c, rbcIDs, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (s *Service) checkSign(arg map[string]interface{}, secret string) (err error) {
	var signval string
	val := url.Values{}
	for k, v := range arg {
		if k == "sign" {
			signval = fmt.Sprintf("%v", v)
			continue
		}
		val.Add(k, fmt.Sprintf("%v", v))
	}
	if err = s.sign(val, signval, secret); err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (s *Service) sign(val url.Values, signVal, secret string) (err error) {
	encode := val.Encode()
	if strings.IndexByte(encode, '+') > -1 {
		encode = strings.Replace(encode, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(encode)
	b.WriteString(secret)
	hash := md5.New()
	hash.Write(b.Bytes())
	sum := hash.Sum(nil)
	sign := hex.EncodeToString(sum)
	if sign != signVal {
		err = ecode.SignCheckErr
	}
	return
}
