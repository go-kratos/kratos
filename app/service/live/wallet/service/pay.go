package service

import (
	"context"
	"go-common/app/service/live/wallet/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// PayHandler pay handler
type PayHandler struct {
	service *WalletService
}

// NeedCheckUid need check uid
func (handler *PayHandler) NeedCheckUid() bool {
	return true
}

// NeedTransactionMutex need transaction mutex
func (handler *PayHandler) NeedTransactionMutex() bool {
	return true
}

// SetWalletService set wallet service
func (handler *PayHandler) SetWalletService(ws *WalletService) {
	handler.service = ws
}

// BizExecute biz execute
func (handler *PayHandler) BizExecute(ws *WalletService, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	platform, _ := params[0].(string)
	if !model.IsValidPlatform(platform) {
		err = ecode.RequestErr
		return
	}

	arg, _ := params[1].(*model.RechargeOrPayForm)

	if uid <= 0 || arg.CoinNum <= 0 || arg.ExtendTid == "" || !model.IsValidCoinType(arg.CoinType) {
		err = ecode.RequestErr
		return
	}
	var reason interface{}
	if len(params) > 2 {
		reason = params[2]
	} else {
		reason = nil
	}

	sysCoinType := model.GetSysCoinType(arg.CoinType, platform)
	sysCoinTypeNo := model.GetCoinTypeNumber(sysCoinType)

	log.Info("TAdd# opr: pay, tid:%s,etid:%s,uid,%d,platform:%s,type:%d,num:%d,time:%d",
		arg.TransactionId, arg.ExtendTid, uid, platform, sysCoinTypeNo, arg.CoinNum, arg.Timestamp)

	coinStreamRecord := model.NewCoinStream(uid, arg.TransactionId, arg.ExtendTid, sysCoinTypeNo, -1*arg.CoinNum, int32(model.PAYTYPE), arg.Timestamp,
		basicParam.BizCode, basicParam.Area, basicParam.Source, basicParam.BizSource, basicParam.MetaData)
	model.AddMoreParam2CoinStream(coinStreamRecord, basicParam, platform)

	// 初始状态
	coinStreamRecord.OrgCoinNum = -1
	coinStreamRecord.OpResult = model.STREAM_OP_RESULT_SUB_FAILED

	var (
		executeRes                 = false
		originCoin     interface{} = -1
		originUserCoin int64       // 用于写入数据库中的org_coin_num字段
	)

	// 获取交易前的用户数据用来1:校验参数　2: 用于最后返回数据
	originMelon, dbErr := ws.s.dao.Melonseed(ws.c, uid)

	for {
		// 查询出错则结束交易
		if dbErr != nil {
			err = ecode.ServerErr
			log.Error("wallet user get coin failed : uid :%d", uid)
			coinStreamRecord.OpReason = model.STREAM_OP_REASON_PRE_QUERY_FAILED
			break
		}

		if model.IsLocalCoin(sysCoinTypeNo) {
			originCoin = model.GetCoinByMelonseed(sysCoinTypeNo, originMelon)
		} else {
			// 如果硬币支付需要单独获取数据，查询失败结束交易
			originCoin, err = ws.s.dao.GetMetal(ws.c, uid)
			if err != nil {
				log.Error("wallet user get metal failed : uid :%d", uid)
				coinStreamRecord.OpReason = model.STREAM_OP_REASON_PRE_QUERY_FAILED
				originCoin = -1
				err = ecode.ServerErr
				break
			}
		}

		// 校验参数　货币不足结束交易
		if !model.CompareCoin(originCoin, arg.CoinNum) {
			coinStreamRecord.OpReason = model.STREAM_OP_REASON_NOT_ENOUGH_COIN
			err = ecode.CoinNotEnough
			break
		}

		// 对用户加锁，加锁失败结束交易
		lockErr := ws.lockUser()
		if lockErr != nil {

			if lockErr == ecode.TargetBlocked {
				coinStreamRecord.OpReason = model.STREAM_OP_REASON_LOCK_FAILED
			} else {
				coinStreamRecord.OpReason = model.STREAM_OP_REASON_LOCK_ERROR
			}
			err = lockErr
			// 跳出
			break
		}

		// 锁定成功再次获取数据
		var lockAfterCoin interface{}
		lockAfterCoin, err = ws.s.dao.GetCoin(ws.c, sysCoinTypeNo, uid)
		if err != nil {
			log.Error("wallet user get coin failed after lock : uid :%d,coinType:%d", uid, sysCoinTypeNo)
			err = ecode.ServerErr
			coinStreamRecord.OpReason = model.STREAM_OP_REASON_QUERY_FAILED
			break
		}

		if !model.CompareCoin(lockAfterCoin, arg.CoinNum) {
			coinStreamRecord.OpReason = model.STREAM_OP_REASON_NOT_ENOUGH_COIN
			err = ecode.CoinNotEnough
			break
		}
		// 赋予锁后的数据
		originCoin = lockAfterCoin

		// 更新数据
		var success bool
		success, err = ws.s.dao.ConsumeCoin(ws.c, int(arg.CoinNum), uid, sysCoinTypeNo, 0, true, reason)
		if !success {
			// 虽然上面检查了余额并加锁，但是如硬币这样的主站货币live这边依然无法保证能够锁住依然可能会有余额不够的情况
			if err == ecode.CoinNotEnough {
				coinStreamRecord.OpReason = model.STREAM_OP_REASON_NOT_ENOUGH_COIN
				err = ecode.CoinNotEnough
				break
			}
			if err != nil { // 更新失败
				log.Error("tx#oper pay update recharge uid :%d,sysCoinTypeNo:%d,coinNum:%d err:%s", uid, sysCoinTypeNo, arg.CoinNum, err.Error())
				err = ecode.PayFailed
				break
			}

			var consumeAfterCoin interface{}
			consumeAfterCoin, err = ws.s.dao.GetCoin(ws.c, sysCoinTypeNo, uid)
			if err != nil {
				err = ecode.PayFailed
				coinStreamRecord.OpReason = model.STREAM_OP_REASON_POST_QUERY_FAILED
				break
			}

			if model.SubCoin(originCoin, consumeAfterCoin) == arg.CoinNum {
				executeRes = true
				break
			}
			err = ecode.PayFailed
			break

		}

		log.Info("tx#oper pay success uid :%d,sysCoinTypeNo:%d,coinNum:%d ", uid, sysCoinTypeNo, arg.CoinNum)
		executeRes = true

		goto payEnd
	}

payEnd:
	originUserCoin = model.GetDbFitCoin(originCoin)
	coinStreamRecord.OrgCoinNum = originUserCoin

	if executeRes {
		coinStreamRecord.OpResult = model.STREAM_OP_RESULT_SUB_SUCC
		model.IncrMelonseedCoin(originMelon, 0-arg.CoinNum, sysCoinTypeNo)
		v = model.GetMelonseedResp(platform, originMelon)
		if model.IsLocalCoin(sysCoinTypeNo) {
			ws.s.pubWalletChange(ws.c, uid, "pay", arg.CoinNum, arg.CoinType, platform, "", 0)
		}
	}

	ws.s.dao.NewCoinStreamRecord(ws.c, coinStreamRecord)

	return

}

// LocalPay local pay
func (handler *PayHandler) LocalPay(c context.Context, payParam *model.RechargeOrPayParam) (resp *model.MelonseedResp, err error) {
	// 校验参数
	if payParam.Uid <= 0 ||
		payParam.CoinNum <= 0 ||
		payParam.ExtendTid == "" ||
		!model.IsValidCoinType(payParam.CoinType) ||
		!model.IsValidPlatform(payParam.Platform) {
		err = ecode.RequestErr
		return
	}

	sysCoinType := model.GetSysCoinType(payParam.CoinType, payParam.Platform)
	sysCoinTypeNo := model.GetCoinTypeNumber(sysCoinType)

	if !model.IsLocalCoin(sysCoinTypeNo) {
		err = ecode.RequestErr
		return
	}

	// 锁住tid
	err = handler.service.lockTransactionId(payParam.TransactionId)
	if err != nil {
		return
	}
	log.Info("TAdd# opr: pay, tid:%s,etid:%s,uid,%d,platform:%s,type:%d,num:%d,time:%d",
		payParam.TransactionId, payParam.GetExtendTid(), payParam.GetUid(), payParam.GetPlatform(), sysCoinTypeNo, payParam.CoinNum, payParam.Timestamp)

	coinStream := model.CoinStreamRecord{}
	model.InjectFieldToCoinStream(&coinStream, payParam)
	coinStream.DeltaCoinNum = -payParam.CoinNum
	coinStream.CoinType = sysCoinTypeNo
	coinStream.OpType = int32(model.PAYTYPE)
	// 初始状态
	coinStream.OrgCoinNum = -1
	coinStream.OpResult = model.STREAM_OP_RESULT_SUB_FAILED

	userLock := handler.service.getUserLock()
	// 锁用户
	lockErr := userLock.lock(payParam.Uid)
	if lockErr != nil {
		model.SetReasonByLockErr(lockErr, &coinStream)
		err = lockErr
		handler.service.s.dao.NewCoinStreamRecord(c, &coinStream)
		return

	}
	defer userLock.release()

	wallet, err := handler.localPay(payParam.Uid, payParam.Platform, sysCoinTypeNo, payParam.CoinNum, &coinStream)

	if err == nil {
		handler.service.s.dao.DelWalletCache(c, payParam.Uid)
		resp = model.GetMelonByDetailWithSnapShot(wallet, payParam.Platform)
		handler.service.s.pubWalletChangeWithDetailSnapShot(c,
			payParam.Uid, "pay", payParam.CoinNum, payParam.CoinType, payParam.Platform, "", 0, wallet)
	}
	return
}

// DetailWithSnapShotWrapper detail with snapshot wrapper
type DetailWithSnapShotWrapper struct {
	resp     *model.DetailWithSnapShot
	logicErr error
}

func (handler *PayHandler) localPay(uid int64, platform string, sysCoinTypeNo int32, coinNum int64, stream *model.CoinStreamRecord) (resp *model.DetailWithSnapShot, err error) {
	dao := handler.service.s.dao
	v, err := dao.DoTx(handler.service.c, func(conn *sql.Tx) (v interface{}, err error) {
		return handler.localPayInTx(conn, uid, platform, sysCoinTypeNo, coinNum, stream)
	})
	if err != nil {
		return
	}
	wrapper := v.(*DetailWithSnapShotWrapper)
	if wrapper.logicErr == nil && err == nil {
		resp = wrapper.resp
	}
	if wrapper.logicErr != nil {
		err = wrapper.logicErr
	}

	return
}
func (handler *PayHandler) localPayInTx(tx *sql.Tx, uid int64, platform string, sysCoinTypeNo int32, coinNum int64, stream *model.CoinStreamRecord) (wrapper *DetailWithSnapShotWrapper, err error) {
	wrapper = new(DetailWithSnapShotWrapper)
	dao := handler.service.s.dao

	// 获取数据 for update
	wallet, err := dao.WalletForUpdate(tx, uid)
	if err != nil {
		return
	}
	curCoin := model.GetCoinByDetailWithSnapShot(sysCoinTypeNo, wallet)
	stream.OrgCoinNum = curCoin
	defer func() {
		if err == nil {
			_, err = dao.NewCoinStreamRecordInTx(tx, stream)
		}
	}()
	if curCoin < coinNum {
		wrapper.logicErr = ecode.CoinNotEnough
		stream.OpReason = model.STREAM_OP_REASON_NOT_ENOUGH_COIN
		return
	}
	// 减款
	_, err = dao.PayCoinInTx(tx, uid, sysCoinTypeNo, coinNum, wallet)
	if err != nil {
		return
	}

	stream.OpResult = model.STREAM_OP_RESULT_SUB_SUCC
	model.ModifyCoinInDetailWithSnapShot(wallet, sysCoinTypeNo, -coinNum)

	wrapper.resp = wallet
	return

}

// Pay pay
func (s *Service) Pay(c context.Context, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	platform, _ := params[0].(string)
	if !model.IsValidPlatform(platform) {
		err = ecode.RequestErr
		return
	}
	arg, _ := params[1].(*model.RechargeOrPayForm)
	if !model.IsValidCoinType(arg.CoinType) {
		err = ecode.RequestErr
		return
	}
	sysCoinType := model.GetSysCoinType(arg.CoinType, platform)
	sysCoinTypeNo := model.GetCoinTypeNumber(sysCoinType)
	if model.IsLocalCoin(sysCoinTypeNo) {
		// 走新的
		platform, _ := params[0].(string)
		arg, _ := params[1].(*model.RechargeOrPayForm)
		payParam := buildRechargeOrPayParam(platform, arg, uid, basicParam)
		ws := new(WalletService)
		ws.c = c
		ws.s = s
		handler := getPayHandler(ws)
		return handler.LocalPay(c, payParam)

	} else {
		handler := PayHandler{}
		return s.execByHandler(&handler, c, basicParam, uid, params...)
	}
}

func getPayHandler(ws *WalletService) *PayHandler {
	handler := PayHandler{}
	handler.SetWalletService(ws)
	return &handler
}
