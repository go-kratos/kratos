package service

import (
	"context"
	"go-common/app/service/live/wallet/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

type ModifyHandler struct {
	service *WalletService
}

func (handler *ModifyHandler) NeedCheckUid() bool {
	return true
}
func (handler *ModifyHandler) NeedTransactionMutex() bool {
	return true
}

func (handler *ModifyHandler) SetWalletService(ws *WalletService) {
	handler.service = ws
}

// old TODO deprecated
func (handler *ModifyHandler) BizExecute(ws *WalletService, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	return
}

func (handler *ModifyHandler) Modify(c context.Context, param *model.RechargeOrPayParam) (resp *model.MelonseedResp, err error) {
	// 校验参数
	if param.Uid <= 0 ||
		param.CoinNum == 0 ||
		param.ExtendTid == "" ||
		!model.IsValidCoinType(param.CoinType) ||
		!model.IsValidPlatform(param.Platform) {
		err = ecode.RequestErr
		return
	}

	sysCoinType := model.GetSysCoinType(param.CoinType, param.Platform)
	sysCoinTypeNo := model.GetCoinTypeNumber(sysCoinType)

	if !model.IsLocalCoin(sysCoinTypeNo) {
		err = ecode.RequestErr
		return
	}

	// 锁住tid
	err = handler.service.lockTransactionId(param.TransactionId)
	if err != nil {
		return
	}
	log.Info("TAdd# opr: modify, tid:%s,etid:%s,uid,%d,platform:%s,type:%d,num:%d,time:%d",
		param.TransactionId, param.GetExtendTid(), param.GetUid(), param.GetPlatform(), sysCoinTypeNo, param.CoinNum, param.Timestamp)

	var serviceType model.ServiceType
	var failedResult int32

	if param.CoinNum < 0 {
		serviceType = model.PAYTYPE
		failedResult = model.STREAM_OP_RESULT_SUB_FAILED
	} else {
		serviceType = model.RECHARGETYPE
		failedResult = model.STREAM_OP_RESULT_ADD_FAILED
	}
	coinStream := model.CoinStreamRecord{}
	model.InjectFieldToCoinStream(&coinStream, param)
	coinStream.DeltaCoinNum = param.CoinNum
	coinStream.CoinType = sysCoinTypeNo
	coinStream.OpType = int32(serviceType)
	coinStream.OpResult = failedResult

	userLock := handler.service.getUserLock()
	// 锁用户
	lockErr := userLock.lock(param.Uid)
	if lockErr != nil {
		model.SetReasonByLockErr(lockErr, &coinStream)
		err = lockErr
		handler.service.s.dao.NewCoinStreamRecord(c, &coinStream)
		return

	}
	defer userLock.release()

	// 实际的db操作
	wallet, err := handler.modify(param.Uid, param.Platform, sysCoinTypeNo, param.CoinNum, &coinStream)

	if err == nil {
		handler.service.s.dao.DelWalletCache(c, param.Uid)
		resp = model.GetMelonByDetailWithSnapShot(wallet, param.Platform)
		handler.service.s.pubWalletChangeWithDetailSnapShot(c,
			param.Uid, "modify", param.CoinNum, param.CoinType, param.Platform, "", 0, wallet)
		log.Info("tx#oper modify success uid :%d,CoinTypeNo:%d,coinNum:%d tid:%s", param.Uid, sysCoinTypeNo, param.CoinNum, param.TransactionId)
	}
	return
}
func (handler *ModifyHandler) modify(uid int64, platform string, sysCoinTypeNo int32, coinNum int64, stream *model.CoinStreamRecord) (resp *model.DetailWithSnapShot, err error) {
	dao := handler.service.s.dao
	v, err := dao.DoTx(handler.service.c, func(conn *sql.Tx) (v interface{}, err error) {
		return handler.modifyInTx(conn, uid, platform, sysCoinTypeNo, coinNum, stream)
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
func (handler *ModifyHandler) modifyInTx(tx *sql.Tx, uid int64, platform string, sysCoinTypeNo int32, coinNum int64, stream *model.CoinStreamRecord) (wrapper *DetailWithSnapShotWrapper, err error) {
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
	if coinNum < 0 {
		if curCoin < -coinNum {
			wrapper.logicErr = ecode.CoinNotEnough
			stream.OpReason = model.STREAM_OP_REASON_NOT_ENOUGH_COIN
			return
		}
	}
	_, err = dao.ModifyCoinInTx(tx, uid, sysCoinTypeNo, coinNum, wallet)
	if err != nil {
		return
	}

	var succResult int32

	if coinNum < 0 {
		succResult = model.STREAM_OP_RESULT_SUB_SUCC
	} else {
		succResult = model.STREAM_OP_RESULT_ADD_SUCC
	}
	stream.OpResult = succResult
	model.ModifyCoinInDetailWithSnapShot(wallet, sysCoinTypeNo, coinNum)

	wrapper.resp = wallet
	return

}

func (s *Service) Modify(c context.Context, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	platform, _ := params[0].(string)
	arg, _ := params[1].(*model.RechargeOrPayForm)
	rechargeParam := buildRechargeOrPayParam(platform, arg, uid, basicParam)

	handler := ModifyHandler{}
	ws := new(WalletService)
	ws.c = c
	ws.s = s
	ws.SetServiceHandler(&handler)
	handler.SetWalletService(ws)
	return handler.Modify(c, rechargeParam)
}
