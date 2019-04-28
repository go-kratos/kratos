package service

import (
	"context"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

type Handler interface {
	NeedTransactionMutex() bool
	NeedCheckUid() bool
	BizExecute(ws *WalletService, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error)
}

type WalletService struct {
	handler   Handler
	s         *Service
	c         context.Context
	uid       int64
	locked    bool
	lockValue string
}

func (ws *WalletService) SetServiceHandler(handler Handler) {
	ws.handler = handler
}

func (ws *WalletService) lockTransactionId(transactionId string) (err error) {
	if transactionId == "" {
		err = ecode.RequestErr
		return
	}
	err = ws.s.dao.LockTransactionId(ws.c, transactionId)
	if err != nil {
		if ws.s.dao.IsLockFailedError(err) {
			err = ecode.TargetBlocked
		} else {
			err = ecode.ServerErr
		}
		log.Error("wallet lock tid failed : %s, err:%s", transactionId, err.Error())
		return
	}
	return
}

func (ws *WalletService) Execute(basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	// 基本检测
	if ws.handler.NeedCheckUid() && uid <= 0 {
		err = ecode.RequestErr
		return
	}
	needTransactionMutex := ws.handler.NeedTransactionMutex()
	if needTransactionMutex {
		if basicParam.TransactionId == "" {
			err = ecode.RequestErr
			return
		}
		err = ws.s.dao.LockTransactionId(ws.c, basicParam.TransactionId)
		if err != nil {
			if ws.s.dao.IsLockFailedError(err) {
				err = ecode.TargetBlocked
			} else {
				err = ecode.ServerErr
			}
			log.Error("wallet lock tid failed : %s, err:%s", basicParam.TransactionId, err.Error())
			return
		}
	}

	ws.uid = uid

	v, err = ws.handler.BizExecute(ws, basicParam, uid, params...)
	ws.unLockUser()
	return
}
func (ws *WalletService) unLockUser() {
	if ws.locked {
		ws.s.dao.UnLockUser(ws.c, ws.uid, ws.lockValue)
	}
	ws.locked = false
	ws.uid = 0
}

func (ws *WalletService) lockUser() error {
	err, gotLocked, lockValue := ws.s.dao.LockUser(ws.c, ws.uid)
	if gotLocked {
		ws.lockValue = lockValue
		ws.locked = true
	} else if err == nil {
		err = ecode.TargetBlocked
	} else {
		err = ecode.ServerErr
	}

	if err != nil {
		log.Error("wallet user lock failed : %d", ws.uid)
	}

	return err
}

func (ws *WalletService) lockSpecificUser(uid int64) error {
	ws.uid = uid
	return ws.lockUser()
}

func (ws *WalletService) getUserLock() UserLock {
	//r = &RedisUserLock{ws: ws}
	r := &NopUserLock{}
	return r
}

func buildRechargeOrPayParam(platform string, arg *model.RechargeOrPayForm, uid int64, basicParam *model.BasicParam) *model.RechargeOrPayParam {
	rechargeParam := model.RechargeOrPayParam{
		Uid: uid, CoinType: arg.CoinType, CoinNum: arg.CoinNum, ExtendTid: arg.ExtendTid, TransactionId: arg.TransactionId,
		Timestamp: arg.Timestamp, BizCode: basicParam.BizCode, Area: basicParam.Area, Source: basicParam.Source, MetaData: basicParam.MetaData,
		BizSource: basicParam.BizSource, Reason: basicParam.Reason, Version: basicParam.Version, Platform: platform,
	}
	return &rechargeParam
}

func buildExchangeParam(platform string, arg *model.ExchangeForm, uid int64, basicParam *model.BasicParam) *model.ExchangeParam {
	param := model.ExchangeParam{
		Uid: uid, SrcCoinType: arg.SrcCoinType, SrcCoinNum: arg.SrcCoinNum, DestCoinType: arg.DestCoinType, DestCoinNum: arg.DestCoinNum, ExtendTid: arg.ExtendTid, TransactionId: arg.TransactionId,
		Timestamp: arg.Timestamp, BizCode: basicParam.BizCode, Area: basicParam.Area, Source: basicParam.Source, MetaData: basicParam.MetaData,
		BizSource: basicParam.BizSource, Reason: basicParam.Reason, Version: basicParam.Version, Platform: platform,
	}
	return &param
}
