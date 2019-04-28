package service

import (
	"context"
	"go-common/app/service/live/wallet/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
)

type RechargeHandler struct {
	service *WalletService
}

func (handler *RechargeHandler) NeedCheckUid() bool {
	return true
}

func (handler *RechargeHandler) NeedTransactionMutex() bool {
	return true
}

func (handler *RechargeHandler) SetWalletService(ws *WalletService) {
	handler.service = ws
}

// old TODO deprecated
func (handler *RechargeHandler) BizExecute(ws *WalletService, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	// do nothing
	return
}

func (handler *RechargeHandler) Recharge(c context.Context, rechargeParam *model.RechargeOrPayParam) (resp *model.MelonseedResp, err error) {
	// 校验参数
	if rechargeParam.Uid <= 0 ||
		rechargeParam.CoinNum <= 0 ||
		rechargeParam.ExtendTid == "" ||
		!model.IsValidCoinType(rechargeParam.CoinType) ||
		!model.IsValidPlatform(rechargeParam.Platform) {
		err = ecode.RequestErr
		return
	}

	sysCoinType := model.GetSysCoinType(rechargeParam.CoinType, rechargeParam.Platform)
	sysCoinTypeNo := model.GetCoinTypeNumber(sysCoinType)

	if !model.IsLocalCoin(sysCoinTypeNo) {
		err = ecode.RequestErr
		return
	}

	// 锁住tid
	err = handler.service.lockTransactionId(rechargeParam.TransactionId)
	if err != nil {
		return
	}

	coinStream := model.CoinStreamRecord{}
	model.InjectFieldToCoinStream(&coinStream, rechargeParam)
	coinStream.DeltaCoinNum = rechargeParam.CoinNum
	coinStream.CoinType = sysCoinTypeNo
	coinStream.OpType = int32(model.RECHARGETYPE)
	// 初始状态
	coinStream.OrgCoinNum = -1
	coinStream.OpResult = model.STREAM_OP_RESULT_ADD_FAILED

	userLock := handler.service.getUserLock()
	// 锁用户
	lockErr := userLock.lock(rechargeParam.Uid)
	if lockErr != nil {
		model.SetReasonByLockErr(lockErr, &coinStream)
		err = lockErr
		handler.service.s.dao.NewCoinStreamRecord(c, &coinStream)
		return

	}
	defer userLock.release()

	// 实际的db操作
	wallet, err := handler.recharge(rechargeParam.Uid, rechargeParam.Platform, sysCoinTypeNo, rechargeParam.CoinNum, &coinStream)

	if err == nil {
		handler.service.s.dao.DelWalletCache(c, rechargeParam.Uid)
		resp = model.GetMelonByDetailWithSnapShot(wallet, rechargeParam.Platform)
		handler.service.s.pubWalletChangeWithDetailSnapShot(c,
			rechargeParam.Uid, "recharge", rechargeParam.CoinNum, rechargeParam.CoinType, rechargeParam.Platform, "", 0, wallet)
	}
	return

}

func (handler *RechargeHandler) recharge(uid int64, platform string, sysCoinTypeNo int32, coinNum int64, coinStream *model.CoinStreamRecord) (resp *model.DetailWithSnapShot, err error) {
	dao := handler.service.s.dao
	v, err := dao.DoTx(handler.service.c, func(conn *sql.Tx) (v interface{}, err error) {
		return handler.rechargeInTx(conn, uid, platform, sysCoinTypeNo, coinNum, coinStream)
	})
	if err == nil {
		resp = v.(*model.DetailWithSnapShot)
	}

	return
}
func (handler *RechargeHandler) rechargeInTx(tx *sql.Tx, uid int64, platform string, sysCoinTypeNo int32, coinNum int64, coinStream *model.CoinStreamRecord) (resp *model.DetailWithSnapShot, err error) {
	dao := handler.service.s.dao

	// 获取数据 for update
	wallet, err := dao.WalletForUpdate(tx, uid)
	if err != nil {
		return
	}
	// 加钱并且做双余额处理
	_, err = dao.RechargeCoinInTx(tx, uid, sysCoinTypeNo, coinNum, wallet)
	if err != nil {
		return
	}

	// 写流水
	coinStream.OrgCoinNum = model.GetCoinByDetailWithSnapShot(sysCoinTypeNo, wallet)
	coinStream.OpResult = model.STREAM_OP_RESULT_ADD_SUCC
	_, err = dao.NewCoinStreamRecordInTx(tx, coinStream)
	if err != nil {
		return
	}

	model.ModifyCoinInDetailWithSnapShot(wallet, sysCoinTypeNo, coinNum)

	resp = wallet
	return
}

func (s *Service) Recharge(c context.Context, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	platform, _ := params[0].(string)
	arg, _ := params[1].(*model.RechargeOrPayForm)
	rechargeParam := buildRechargeOrPayParam(platform, arg, uid, basicParam)

	ws := new(WalletService)
	ws.c = c
	ws.s = s
	handler := getRechargeHandler(ws)
	handler.SetWalletService(ws)
	return handler.Recharge(c, rechargeParam)
}

func getRechargeHandler(ws *WalletService) *RechargeHandler {
	handler := RechargeHandler{}
	handler.SetWalletService(ws)
	return &handler
}
