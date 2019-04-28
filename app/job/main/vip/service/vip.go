package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"go-common/app/job/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	iapChannelID = 100
)

func (s *Service) cleanCacheAndNotify(c context.Context, hv *model.HandlerVip) (err error) {
	s.dao.DelInfoCache(c, hv.Mid)
	if err = s.dao.SendCleanCache(c, hv); err != nil {
		return
	}
	if err = s.dao.DelVipInfoCache(c, int64(hv.Mid)); err != nil {
		log.Error("del vip info cache (mid:%v) error(%+v)", hv.Mid, err)
		return
	}
	eg, ec := errgroup.WithContext(c)
	for _, app := range s.appMap {
		ta := app
		eg.Go(func() error {
			if err = s.dao.SendAppCleanCache(ec, hv, ta); err == nil {
				log.Info("SendAppCleanCache success hv(%v) app(%v)", hv, ta)
			} else {
				ac := new(model.AppCache)
				ac.AppID = ta.ID
				ac.Mid = hv.Mid
				s.cleanAppCache <- ac
			}
			return nil
		})
	}
	if err = eg.Wait(); err != nil {
		log.Error(" eg.Wait err(%+v)", err)
	}
	err = nil
	return
}

//ScanUserInfo scan all userinfo update status
func (s *Service) ScanUserInfo(c context.Context) (err error) {
	var (
		ot        = time.Now().Format("2006-01-02 15:04:05")
		userInfos []*model.VipUserInfo
		size      = 2000
		endID     = 0
	)
	for {
		if endID, err = s.dao.SelOldUserInfoMaxID(context.TODO()); err != nil {
			time.Sleep(time.Minute * 2)
			continue
		}
		break
	}

	page := endID / size
	if endID%size != 0 {
		page++
	}
	for i := 0; i < page; {
		startID := i * size
		eID := (i + 1) * size
		if userInfos, err = s.dao.SelVipList(context.TODO(), startID, eID, ot); err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		i++
		for _, v := range userInfos {
			s.updateUserInfo(context.TODO(), v)
		}
	}
	return
}

func (s *Service) updateUserInfo(c context.Context, v *model.VipUserInfo) (err error) {
	var (
		curTime = time.Now()
		fType   = v.Type
		fStatus = v.Status
	)
	if v.AnnualVipOverdueTime.Time().Before(curTime) {
		fType = model.Vip
	}
	if v.OverdueTime.Time().Before(curTime) {
		fStatus = model.VipStatusOverTime
	}
	if fType != v.Type || fStatus != v.Status {
		v.Type = fType
		v.Status = fStatus
		if v.Status == model.VipStatusOverTime && v.PayChannelID == iapChannelID {
			v.PayType = model.Normal
		}
		if _, err = s.dao.UpdateVipUser(c, int64(v.Mid), v.Status, v.Type, v.PayType); err != nil {
			return
		}
		s.dao.DelInfoCache(c, v.Mid)
		s.dao.DelVipInfoCache(c, int64(v.Mid))
	}
	return
}

func (s *Service) handlerautorenewlogproc() {
	var (
		err error
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerautorenewlogproc panic(%v)", x)
			go s.handlerautorenewlogproc()
			log.Info("service.handlerautorenewlogproc recover")
		}
	}()
	for {
		user := <-s.handlerAutoRenewLog
		for i := 0; i <= s.c.Property.Retry; i++ {
			if err = s.handlerAutoRenewLogInfo(context.TODO(), user); err == nil {
				break
			}
			log.Error("%+v", err)
			time.Sleep(2 * time.Second)
		}

	}
}

func (s *Service) handlerAutoRenewLogInfo(c context.Context, user *model.VipUserInfo) (err error) {
	var (
		payOrder *model.VipPayOrder
		paylog   *model.VipPayOrderLog
		rlog     *model.VipPayOrderLog
	)
	if user.PayType == model.AutoRenew {
		if user.PayChannelID == iapChannelID {
			if payOrder, err = s.dao.SelPayOrderByMid(c, user.Mid, model.IAPAutoRenew, model.SUCCESS); err != nil {
				err = errors.WithStack(err)
				return
			}
			if payOrder == nil {
				err = errors.Errorf("订单号不能为空......")
				return
			}
			rlog = new(model.VipPayOrderLog)
			rlog.Mid = payOrder.Mid
			rlog.OrderNo = payOrder.OrderNo
			rlog.Status = model.SIGN
		} else {
			if payOrder, err = s.dao.SelPayOrderByMid(c, user.Mid, model.AutoRenew, model.SUCCESS); err != nil {
				err = errors.WithStack(err)
				return
			}
			if payOrder == nil {
				err = errors.Errorf("订单号不能为空......")
				return
			}
			rlog = new(model.VipPayOrderLog)
			rlog.Mid = payOrder.Mid
			rlog.OrderNo = payOrder.OrderNo
			rlog.Status = model.SIGN
		}
	} else {
		if paylog, err = s.dao.SelPayOrderLog(c, user.Mid, model.SIGN); err != nil {
			err = errors.WithStack(err)
			return
		}
		rlog = new(model.VipPayOrderLog)
		rlog.Mid = paylog.Mid
		rlog.Status = model.UNSIGN
		rlog.OrderNo = paylog.OrderNo
	}

	if rlog != nil {
		if _, err = s.dao.AddPayOrderLog(c, rlog); err != nil {
			err = errors.WithStack(err)
			return
		}
	}

	return
}

func (s *Service) handlerinsertuserinfoproc() {
	var (
		err error
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerinsertuserinfoproc panic(%v)", x)
			go s.handlerinsertuserinfoproc()
			log.Info("service.handlerinsertuserinfoproc recover")
		}
	}()
	for {
		userInfo := <-s.handlerInsertUserInfo
		for i := 0; i < s.c.Property.Retry; i++ {
			if err = s.addUserInfo(context.TODO(), userInfo); err == nil {
				s.dao.DelInfoCache(context.Background(), userInfo.Mid)
				s.dao.DelVipInfoCache(context.TODO(), userInfo.Mid)
				if s.grayScope(userInfo.Mid) {
					s.cleanCache(userInfo.Mid)
				}
				break
			}
			log.Error("add info error(%+v)", err)
		}

	}
}

func (s *Service) addUserInfo(c context.Context, ui *model.VipUserInfo) (err error) {
	var (
		tx  *sql.Tx
		udh *model.VipUserDiscountHistory
	)
	if tx, err = s.dao.StartTx(c); err != nil {
		return
	}
	defer func() {
		if err == nil {
			if err = tx.Commit(); err != nil {
				log.Error("commit(%+v)", err)
				return
			}
		} else {
			tx.Rollback()
		}
	}()
	if _, err = s.dao.AddUserInfo(tx, ui); err != nil {
		err = errors.WithStack(err)
		return
	}
	if ui.AutoRenewed == 1 {
		udh = new(model.VipUserDiscountHistory)
		udh.DiscountID = model.VipUserFirstDiscount
		udh.Status = model.DiscountUsed
		udh.Mid = ui.Mid

		if _, err = s.dao.DupUserDiscountHistory(tx, udh); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return

}

func (s *Service) updateVipUserInfo(c context.Context, ui *model.VipUserInfo) (err error) {
	var (
		tx  *sql.Tx
		udh *model.VipUserDiscountHistory
		eff int64
	)
	if tx, err = s.dao.StartTx(c); err != nil {
		return
	}
	defer func() {
		if err == nil {
			if err = tx.Commit(); err != nil {
				log.Error("commit(%+v)", err)
				return
			}
		} else {
			tx.Rollback()
		}
	}()
	if eff, err = s.dao.UpdateUserInfo(tx, ui); err != nil {
		err = errors.WithStack(err)
		return
	}
	if eff <= 0 {
		log.Warn("update vip RowsAffected 0 vip(%+v)", ui)
		return
	}
	if ui.AutoRenewed == 1 {
		udh = new(model.VipUserDiscountHistory)
		udh.DiscountID = model.VipUserFirstDiscount
		udh.Status = model.DiscountUsed
		udh.Mid = ui.Mid

		if _, err = s.dao.DupUserDiscountHistory(tx, udh); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return

}

func (s *Service) handlerfailuserinfoproc() {
	var (
		err error
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerfailuserinfoproc panic(%v)", x)
			go s.handlerfailuserinfoproc()
			log.Info("service.handlerfailuserinfoproc recover")
		}
	}()
	for {
		userInfo := <-s.handlerFailUserInfo
		_time := 0
		for {
			if err = s.updateVipUserInfo(context.TODO(), userInfo); err == nil {
				s.dao.DelInfoCache(context.Background(), userInfo.Mid)
				s.dao.DelVipInfoCache(context.TODO(), userInfo.Mid)
				if s.grayScope(userInfo.Mid) {
					s.cleanCache(userInfo.Mid)
				}
				break
			}
			log.Error("info error(%+v)", err)
			_time++
			if _time > _maxtime {
				break
			}
			time.Sleep(_sleep)
		}

	}
}

func (s *Service) handlerupdateuserinfoproc() {
	var (
		err  error
		flag bool
	)
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerupdateuserinfoproc panic(%v)", x)
			go s.handlerupdateuserinfoproc()
			log.Info("service.handlerupdateuserinfoproc recover")
		}
	}()
	for {
		userInfo := <-s.handlerUpdateUserInfo
		flag = true
		for i := 0; i < s.c.Property.Retry; i++ {
			if err = s.updateVipUserInfo(context.TODO(), userInfo); err == nil {
				s.dao.DelInfoCache(context.Background(), userInfo.Mid)
				s.dao.DelVipInfoCache(context.TODO(), userInfo.Mid)
				if s.grayScope(userInfo.Mid) {
					s.cleanCache(userInfo.Mid)
				}
				flag = false
				break
			}
			log.Error("info error(%+v)", err)
		}

		if flag {
			s.handlerFailUserInfo <- userInfo
		}

	}
}

func (s *Service) handleraddchangehistoryproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handleraddchangehistoryproc panic(%v)", x)
			go s.handleraddchangehistoryproc()
			log.Info("service.handleraddchangehistoryproc recover")
		}
	}()
	for {
		msg := <-s.handlerAddVipHistory
		history := convertMsgToHistory(msg)
		var res []*model.VipChangeHistory
		res = append(res, history)
		for i := 0; i < s.c.Property.Retry; i++ {
			if err := s.dao.AddChangeHistoryBatch(res); err == nil {
				break
			}
		}
	}
}
func convertMsgToHistory(msg *model.VipChangeHistoryMsg) (r *model.VipChangeHistory) {
	r = new(model.VipChangeHistory)
	r.Mid = msg.Mid
	r.Days = msg.Days
	r.Month = msg.Month
	r.ChangeType = msg.ChangeType
	r.OperatorID = msg.OperatorID
	r.RelationID = msg.RelationID
	r.BatchID = msg.BatchID
	r.Remark = msg.Remark
	r.ChangeTime = xtime.Time(parseTime(msg.ChangeTime).Unix())
	r.BatchCodeID = msg.BatchCodeID
	return
}

func parseTime(timeStr string) (t time.Time) {
	var err error
	if t, err = time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local); err != nil {
		t = time.Now()
	}
	return
}

func convertMsgToUserInfo(msg *model.VipUserInfoMsg) (r *model.VipUserInfo) {
	r = new(model.VipUserInfo)
	r.AnnualVipOverdueTime = xtime.Time(parseTime(msg.AnnualVipOverdueTime).Unix())
	r.Mid = msg.Mid
	r.OverdueTime = xtime.Time(parseTime(msg.OverdueTime).Unix())
	r.PayType = msg.IsAutoRenew
	r.RecentTime = xtime.Time(parseTime(msg.RecentTime).Unix())
	r.StartTime = xtime.Time(parseTime(msg.StartTime).Unix())
	r.Status = msg.Status
	r.Type = msg.Type
	r.PayChannelID = msg.PayChannelID
	r.AutoRenewed = msg.AutoRenewed
	r.IosOverdueTime = xtime.Time(parseTime(msg.IosOverdueTime).Unix())
	r.Ver = msg.Ver
	return
}

func convertUserInfoByNewMsg(msg *model.VipUserInfoNewMsg) (r *model.VipUserInfo) {
	r = new(model.VipUserInfo)
	r.AnnualVipOverdueTime = xtime.Time(parseTime(msg.AnnualVipOverdueTime).Unix())
	r.Mid = msg.Mid
	r.OverdueTime = xtime.Time(parseTime(msg.VipOverdueTime).Unix())
	r.PayType = msg.VipPayType
	r.RecentTime = xtime.Time(parseTime(msg.VipRecentTime).Unix())
	r.StartTime = xtime.Time(parseTime(msg.VipStartTime).Unix())
	r.Status = msg.VipStatus
	r.Type = msg.VipType
	r.PayChannelID = msg.PayChannelID
	r.IosOverdueTime = xtime.Time(parseTime(msg.IosOverdueTime).Unix())
	r.Ver = msg.Ver
	return
}

func convertOldToNew(old *model.VipUserInfoOld) (r *model.VipUserInfo) {
	r = new(model.VipUserInfo)
	r.AnnualVipOverdueTime = old.AnnualVipOverdueTime
	r.Mid = old.Mid
	r.OverdueTime = old.OverdueTime
	r.PayType = old.IsAutoRenew
	r.RecentTime = old.RecentTime
	r.PayChannelID = old.PayChannelID
	if old.RecentTime.Time().Unix() < 0 {
		r.RecentTime = xtime.Time(1451577600)
	}
	r.StartTime = old.StartTime
	r.Status = old.Status
	r.Type = old.Type
	r.IosOverdueTime = old.IosOverdueTime
	r.Ver = old.Ver
	return
}

//HandlerVipChangeHistory handler sync change history data
func (s *Service) HandlerVipChangeHistory() (err error) {
	var (
		newMaxID int64
		oldMaxID int64
		size     = int64(s.c.Property.BatchSize)
		startID  int64
		endID    = size
		exitMap  = make(map[string]int)
	)
	if oldMaxID, err = s.dao.SelOldChangeHistoryMaxID(context.TODO()); err != nil {
		log.Error("selOldChangeHistory error(%+v)", err)
		return
	}
	if newMaxID, err = s.dao.SelChangeHistoryMaxID(context.TODO()); err != nil {
		log.Error("selChangeHistoryMaxID error(%+v)", err)
		return
	}
	page := newMaxID / size
	if newMaxID%size != 0 {
		page++
	}
	for i := 0; i < int(page); i++ {
		startID = int64(i) * size
		endID = int64((i + 1)) * size
		if endID > newMaxID {
			endID = newMaxID
		}

		var res []*model.VipChangeHistory
		if res, err = s.dao.SelChangeHistory(context.TODO(), startID, endID); err != nil {
			log.Error("selChangeHistory(startID:%v endID:%v) error(%+v)", startID, endID, endID)
			return
		}
		for _, v := range res {
			exitMap[s.madeChangeHistoryMD5(v)] = 1
		}
	}

	page = oldMaxID / size
	if oldMaxID%size != 0 {
		page++
	}
	var batch []*model.VipChangeHistory
	for i := 0; i < int(page); i++ {
		startID = int64(i) * size
		endID = int64(i+1) * size
		if endID > oldMaxID {
			endID = oldMaxID
		}

		var res []*model.VipChangeHistory
		if res, err = s.dao.SelOldChangeHistory(context.TODO(), startID, endID); err != nil {
			log.Error("sel old change history (startID:%v endID:%v) error(%+v)", startID, endID, err)
			return
		}

		for _, v := range res {
			v.Days = s.calcDay(v)
			madeMD5 := s.madeChangeHistoryMD5(v)
			if exitMap[madeMD5] == 0 {
				batch = append(batch, v)
			}
		}
		if err = s.dao.AddChangeHistoryBatch(batch); err != nil {
			log.Error("add change history batch(%+v) error(%+v)", batch, err)
			return
		}
		batch = nil

	}
	return
}

func (s *Service) calcDay(r *model.VipChangeHistory) int32 {
	if r.Month != 0 {
		year := r.Month / 12
		month := r.Month % 12

		return int32(year)*model.VipDaysYear + int32(month)*model.VipDaysMonth
	}
	return r.Days
}

func (s *Service) madeChangeHistoryMD5(r *model.VipChangeHistory) string {
	str := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v", r.Mid, r.Remark, r.BatchID, r.RelationID, r.OperatorID, r.Days, r.ChangeTime.Time().Format("2006-01-02 15:04:05"), r.ChangeType, r.BatchCodeID)
	b := []byte(str)
	hash := md5.New()
	hash.Write(b)
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

//SyncUserInfoByMid sync user by mid.
func (s *Service) SyncUserInfoByMid(c context.Context, mid int64) (err error) {
	var (
		old  *model.VipUserInfoOld
		user *model.VipUserInfo
	)

	if old, err = s.dao.OldVipInfo(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if user, err = s.dao.SelVipUserInfo(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}

	r := convertOldToNew(old)
	r.OldVer = user.Ver

	if err = s.updateVipUserInfo(c, r); err != nil {
		err = errors.WithStack(err)
		return
	}
	// clear cache.
	s.cleanVipRetry(mid)
	return
}

// ClearUserCache clear user cache.
func (s *Service) ClearUserCache(mid int64) {
	s.cleanVipRetry(mid)
}

// ClearUserCache clear user cache.
func (s *Service) grayScope(mid int64) bool {
	return mid%10000 < s.c.Property.GrayScope
}
