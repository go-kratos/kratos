package service

import (
	"context"
	"go-common/app/service/live/wallet/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

type OutExchangeHandler struct {
	service *WalletService
}

func (handler *OutExchangeHandler) SetWalletService(ws *WalletService) {
	handler.service = ws
}

func (handler *OutExchangeHandler) exchange(c context.Context, uid int64, platform string, srcSysCoinTypeNo int32, srcCoinNum int64,
	destSysCoinTypeNo int32, destCoinNum int64, srcStream *model.CoinStreamRecord, destStream *model.CoinStreamRecord, exchangeStream *model.CoinExchangeRecord) (resp *model.DetailWithSnapShot, err error) {

	payResp, err := handler.pay(c, uid, platform, srcSysCoinTypeNo, srcCoinNum, destCoinNum, srcStream)
	if err != nil {
		resp = payResp
		return
	}

	rechargeResp, err := handler.recharge(c, uid, platform, destSysCoinTypeNo, destCoinNum, destStream, exchangeStream)
	if err != nil {
		resp = rechargeResp
		// 走到这里说明pay成功了，但是recharge失败，且pay的是直播内部货币(非硬币) recharge的是硬币 那么直接认为成功
		// 原因： 硬币没有提供幂等性接口，无法知道知道最后是否成功，如果认为失败，可能会影响业务 比如用户会多次去扣银瓜子加硬币
		// 如果发生这样的事情 ， 则直播认为成功，用户最后发现硬币没有加上，通过反馈的方式人工接入
		// 如果这样的事情发生太多 则需要推动硬币提供幂等接口
		if destSysCoinTypeNo == model.SysCoinTypeMetal {
			resp = payResp
			err = nil
		}
		return
	}

	if payResp == nil {
		resp = rechargeResp
	} else {
		resp = payResp
	}
	return
}
func (handler *OutExchangeHandler) pay(c context.Context, uid int64, platform string, sysCoinTypeNo int32, coinNum int64, destCoinNum int64, stream *model.CoinStreamRecord) (resp *model.DetailWithSnapShot, err error) {

	if sysCoinTypeNo == model.SysCoinTypeMetal {
		stream.OrgCoinNum = -1
		_, err = handler.service.s.dao.ConsumeCoin(handler.service.c, int(coinNum), uid, sysCoinTypeNo, destCoinNum, true, nil)
		err = handleMetalResp(err, stream, uid, sysCoinTypeNo, coinNum, handler, model.STREAM_OP_RESULT_SUB_SUCC)
		return
	}

	// 本地货币
	payHandler := getPayHandler(handler.service)
	resp, err = payHandler.localPay(uid, platform, sysCoinTypeNo, coinNum, stream)
	return
}

func (handler *OutExchangeHandler) recharge(c context.Context, uid int64, platform string, sysCoinTypeNo int32, coinNum int64, stream *model.CoinStreamRecord, exchangeStream *model.CoinExchangeRecord) (resp *model.DetailWithSnapShot, err error) {

	if sysCoinTypeNo == model.SysCoinTypeMetal {
		stream.OrgCoinNum = -1
		var success bool
		success, _, err = handler.service.s.dao.ModifyMetal(handler.service.c, uid, coinNum, 0, nil)
		err = handleMetalResp(err, stream, uid, sysCoinTypeNo, coinNum, handler, model.STREAM_OP_RESULT_ADD_SUCC)
		if success {
			_, insertErr := handler.service.s.dao.NewCoinExchangeRecord(c, exchangeStream)
			if insertErr != nil {
				log.Error("tx#exchange#metal handle success but insert exchange stream :%s", insertErr.Error())
			}
		}
		return
	}

	// 本地货币
	rechargeHandler := getRechargeHandler(handler.service)

	dao := handler.service.s.dao
	v, err := dao.DoTx(handler.service.c, func(conn *sql.Tx) (v interface{}, err error) {
		v, err = rechargeHandler.rechargeInTx(conn, uid, platform, sysCoinTypeNo, coinNum, stream)
		if err == nil {
			_, err = dao.NewCoinExchangeRecordInTx(conn, exchangeStream)
		}
		return
	})
	if err == nil {
		resp = v.(*model.DetailWithSnapShot)
	}

	return
}

func handleMetalResp(err error, stream *model.CoinStreamRecord, uid int64, sysCoinTypeNo int32, coinNum int64, handler *OutExchangeHandler, successResult int32) error {
	if err != nil {
		// 虽然上面检查了余额并加锁，但是如硬币这样的主站货币live这边依然无法保证能够锁住依然可能会有余额不够的情况
		if err == ecode.CoinNotEnough {
			stream.SetOpReason(model.STREAM_OP_REASON_NOT_ENOUGH_COIN)
		} else { // 更新失败
			stream.SetOpReason(model.STREAM_OP_REASON_EXECUTE_FAILED)
			log.Error("tx#oper exchange sub uid :%d,sysCoinTypeNo:%d,coinNum:%d err:%s", uid, sysCoinTypeNo, coinNum, err.Error())
			err = ecode.ServerErr
		}
	} else { // 成功
		stream.OpResult = successResult
	}
	_, insertErr := handler.service.s.dao.NewCoinStreamRecord(handler.service.c, stream)
	if insertErr != nil {
		log.Error("tx#exchange#metal handle success but insert coin stream :%s", insertErr.Error())
	}
	return err
}

func getOutExchangeHandler(ws *WalletService) *OutExchangeHandler {
	handler := OutExchangeHandler{service: ws}
	return &handler
}
