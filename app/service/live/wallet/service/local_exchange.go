package service

import (
	"context"
	"go-common/app/service/live/wallet/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
)

type LocalExchangeHandler struct {
	service *WalletService
}

func (handler *LocalExchangeHandler) SetWalletService(ws *WalletService) {
	handler.service = ws
}

func (handler *LocalExchangeHandler) exchange(c context.Context, uid int64, platform string, srcSysCoinTypeNo int32, srcCoinNum int64,
	destSysCoinTypeNo int32, destCoinNum int64, srcStream *model.CoinStreamRecord, destStream *model.CoinStreamRecord, exchangeStream *model.CoinExchangeRecord) (resp *model.DetailWithSnapShot, err error) {

	dao := handler.service.s.dao
	v, err := dao.DoTx(handler.service.c, func(conn *sql.Tx) (v interface{}, err error) {
		return handler.exchangeInTx(conn, uid, platform, srcSysCoinTypeNo, srcCoinNum, destSysCoinTypeNo, destCoinNum, srcStream, destStream, exchangeStream)
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
func (handler *LocalExchangeHandler) exchangeInTx(tx *sql.Tx, uid int64, platform string, srcSysCoinTypeNo int32, srcCoinNum int64,
	destSysCoinTypeNo int32, destCoinNum int64, srcStream *model.CoinStreamRecord, destSteam *model.CoinStreamRecord, exchangeStream *model.CoinExchangeRecord) (wrapper *DetailWithSnapShotWrapper, err error) {

	wrapper = new(DetailWithSnapShotWrapper)
	dao := handler.service.s.dao

	// 获取数据 for update
	wallet, err := dao.WalletForUpdate(tx, uid)
	if err != nil {
		return
	}
	srcCurCoin := model.GetCoinByDetailWithSnapShot(srcSysCoinTypeNo, wallet)
	destCurCoin := model.GetCoinByDetailWithSnapShot(destSysCoinTypeNo, wallet)
	srcStream.OrgCoinNum = srcCurCoin
	destSteam.OrgCoinNum = destCurCoin
	needWriteSecondStream := false
	defer func() {
		if err == nil {
			_, err = dao.NewCoinStreamRecordInTx(tx, srcStream)
			if err == nil && needWriteSecondStream {
				_, err = dao.NewCoinStreamRecordInTx(tx, destSteam)
			}
			if err == nil && wrapper.logicErr == nil {
				_, err = dao.NewCoinExchangeRecordInTx(tx, exchangeStream)
			}
		}
	}()
	if srcCurCoin < srcCoinNum {
		wrapper.logicErr = ecode.CoinNotEnough
		srcStream.OpReason = model.STREAM_OP_REASON_NOT_ENOUGH_COIN
		return
	}

	_, err = dao.ExchangeCoinInTx(tx, uid, srcSysCoinTypeNo, srcCoinNum, destSysCoinTypeNo, destCoinNum, wallet)
	if err != nil {
		return
	}
	needWriteSecondStream = true

	srcStream.OpResult = model.STREAM_OP_RESULT_SUB_SUCC
	destSteam.OpResult = model.STREAM_OP_RESULT_ADD_SUCC
	model.ModifyCoinInDetailWithSnapShot(wallet, srcSysCoinTypeNo, -srcCoinNum)
	model.ModifyCoinInDetailWithSnapShot(wallet, destSysCoinTypeNo, destCoinNum)
	wrapper.resp = wallet

	return
}

func getLocalPayHandler(ws *WalletService) *LocalExchangeHandler {
	handler := LocalExchangeHandler{service: ws}
	return &handler
}
