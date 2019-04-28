package service

import (
	"context"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// ExchangeHandler exchange handler
type ExchangeHandler struct {
	service *WalletService
}

// NeedCheckUid need check uid
func (handler *ExchangeHandler) NeedCheckUid() bool {
	return true
}

// NeedTransactionMutex need transaction mutex
func (handler *ExchangeHandler) NeedTransactionMutex() bool {
	return true
}

// SetWalletService set wallet service
func (handler *ExchangeHandler) SetWalletService(ws *WalletService) {
	handler.service = ws
}

// BizExecute old TODO deprecated
func (handler *ExchangeHandler) BizExecute(ws *WalletService, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	return
}

func checkByExchangeParam(param *model.ExchangeParam) (r bool) {
	if !model.IsValidPlatform(param.Platform) {
		return
	}

	for {
		if param.Uid <= 0 || param.ExtendTid == "" || param.TransactionId == "" {
			break
		}

		if param.SrcCoinNum <= 0 || !model.IsValidCoinType(param.SrcCoinType) {
			break
		}
		if param.DestCoinNum <= 0 || !model.IsValidCoinType(param.DestCoinType) {
			break
		}
		r = true
		break
	}
	return r
}

// Exchange exchange
func (handler *ExchangeHandler) Exchange(c context.Context, param *model.ExchangeParam) (resp *model.MelonseedResp, err error) {
	if !checkByExchangeParam(param) {
		err = ecode.RequestErr
		return
	}

	srcSysCoinType := model.GetSysCoinType(param.SrcCoinType, param.Platform)
	srcSysCoinTypeNo := model.GetCoinTypeNumber(srcSysCoinType)
	srcIsLocalCoin := model.IsLocalCoin(srcSysCoinTypeNo)

	destSysCoinType := model.GetSysCoinType(param.DestCoinType, param.Platform)
	destSysCoinTypeNo := model.GetCoinTypeNumber(destSysCoinType)
	destIsLocalCoin := model.IsLocalCoin(destSysCoinTypeNo)

	// 不允许兑换相同货币
	if srcSysCoinTypeNo == destSysCoinTypeNo {
		err = ecode.RequestErr
		return
	}

	// 锁住tid
	err = handler.service.lockTransactionId(param.TransactionId)
	if err != nil {
		return
	}

	log.Info("##TX## Exchange stage:begin| time:%d|tid:%s|uid:%d|platform:%s|etid:%s|srcType:%s|srcNum:%d|destType:%s|destNum:%d",
		param.Timestamp, param.TransactionId, param.Uid, param.Platform, param.ExtendTid, srcSysCoinType, param.SrcCoinNum, destSysCoinType, param.DestCoinNum)

	srcCoinStream := model.CoinStreamRecord{}
	model.InjectFieldToCoinStream(&srcCoinStream, param)
	srcCoinStream.DeltaCoinNum = -param.SrcCoinNum
	srcCoinStream.CoinType = srcSysCoinTypeNo
	srcCoinStream.OpType = int32(model.EXCHANGETYPE)
	// 初始状态
	srcCoinStream.OrgCoinNum = -1
	srcCoinStream.OpResult = model.STREAM_OP_RESULT_SUB_FAILED

	destCoinStream := model.CoinStreamRecord{}
	model.InjectFieldToCoinStream(&destCoinStream, param)
	destCoinStream.DeltaCoinNum = param.DestCoinNum
	destCoinStream.CoinType = destSysCoinTypeNo
	destCoinStream.OpType = int32(model.EXCHANGETYPE)
	// 初始状态
	destCoinStream.OrgCoinNum = -1
	destCoinStream.OpResult = model.STREAM_OP_RESULT_ADD_FAILED

	exchangeRecord := model.NewExchangeSteam(param.Uid, param.TransactionId, srcSysCoinTypeNo, int32(param.SrcCoinNum), destSysCoinTypeNo, int32(param.DestCoinNum), param.Timestamp, 0)

	userLock := handler.service.getUserLock()
	// 锁用户
	lockErr := userLock.lock(param.Uid)
	if lockErr != nil {
		model.SetReasonByLockErr(lockErr, &srcCoinStream)
		err = lockErr
		handler.service.s.dao.NewCoinStreamRecord(c, &srcCoinStream)
		return

	}
	defer userLock.release()

	var realHandler RealExchangeHandler

	if srcIsLocalCoin && destIsLocalCoin {
		realHandler = getLocalPayHandler(handler.service)
	} else {
		realHandler = getOutExchangeHandler(handler.service)
	}

	wallet, err := realHandler.exchange(c, param.Uid, param.Platform, srcSysCoinTypeNo, param.SrcCoinNum, destSysCoinTypeNo, param.DestCoinNum, &srcCoinStream, &destCoinStream, exchangeRecord)
	if err == nil {
		handler.service.s.dao.DelWalletCache(c, param.Uid)
		resp = model.GetMelonByDetailWithSnapShot(wallet, param.Platform)
		handler.service.s.pubWalletChangeWithDetailSnapShot(handler.service.c, param.Uid, "exchange", param.SrcCoinNum, param.SrcCoinType, param.Platform, param.DestCoinType, param.DestCoinNum, wallet)
	}
	return

}

func getExchangeHandler(ws *WalletService) *ExchangeHandler {
	handler := ExchangeHandler{}
	handler.SetWalletService(ws)
	return &handler
}

// Exchange exchange
func (s *Service) Exchange(c context.Context, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {

	platform, _ := params[0].(string)
	arg, _ := params[1].(*model.ExchangeForm)

	param := buildExchangeParam(platform, arg, uid, basicParam)

	ws := new(WalletService)
	ws.c = c
	ws.s = s
	handler := getExchangeHandler(ws)
	handler.SetWalletService(ws)
	return handler.Exchange(c, param)
}

// RealExchangeHandler real exchange handler
type RealExchangeHandler interface {
	exchange(c context.Context, uid int64, platform string, srcSysCoinTypeNo int32, srcCoinNum int64,
		destSysCoinTypeNo int32, destCoinNum int64, srcStream *model.CoinStreamRecord, destStream *model.CoinStreamRecord, exchangeStream *model.CoinExchangeRecord) (*model.DetailWithSnapShot, error)
}
