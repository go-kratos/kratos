package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/service/main/tv/internal/model"
	"go-common/app/service/main/tv/internal/service/validator"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func (s *Service) flushPayParamAsync(c context.Context, token string, payParam *model.PayParam) {
	s.mission(func() {
		if err := s.dao.UpdateCachePayParam(context.TODO(), token, payParam); err != nil {
			log.Error("s.dao.UpdateCachePayParam(%s, %+v) err(%+v)", token, payParam, err)
		}
	})
}

func (s *Service) flushUserInfoAsync(c context.Context, mid int64) {
	s.mission(func() {
		// delete user info cache
		if err := s.dao.DelCacheUserInfoByMid(context.TODO(), mid); err != nil {
			log.Error("s.dao.DelCacheUserInfoByMid(%d) err(%+v)", mid, err)
		}
	})
}

func (s *Service) giveMVipGiftAsync(c context.Context, mid int64, pid int32, orderNo string) {
	s.mission(func() {
		if err := s.GiveMVipGift(context.TODO(), mid, pid, orderNo); err != nil {
			log.Error("s.GiveMVipGift(%d, %d, %s)", mid, pid, orderNo)
		}
	})
}

func (s *Service) initUserInfo(c context.Context, tx *xsql.Tx, mid int64) (ui *model.UserInfo, err error) {
	ui = &model.UserInfo{
		Mid:     mid,
		VipType: model.VipTypeVip,
	}
	if _, err = s.dao.TxInsertUserInfo(c, tx, ui); err != nil {
		log.Info("s.dao.TxInsertUserInfo(%+v) err(%+v)", ui, err)
		return ui, err
	}
	return ui, nil
}

func (s *Service) PayOrder(c context.Context, id int) (po *model.PayOrder, err error) {
	return s.dao.PayOrderByID(c, id)
}

func (s *Service) YstOrderState(c context.Context, seqNo string, traceNo string) (res *model.YstOrderStateReply, err error) {
	ystOrderReq := &model.YstOrderStateReq{
		SeqNo:   seqNo,
		TraceNo: traceNo,
	}
	return s.dao.YstOrderState(c, ystOrderReq)
}

func (s *Service) PayPending(c context.Context, req *model.YstPayCallbackReq, payOrder *model.PayOrder) error {
	log.Info("s.PayCallback.PayPending(%+v, %+v)", req, payOrder)
	return nil
}

func (s *Service) payFail(c context.Context, status string, payOrder *model.PayOrder) error {
	log.Info("s.PayCallback.PayFail(%s, %+v)", status, payOrder)
	var (
		payParam *model.PayParam
		//ystOrder *model.YstOrderStateReply
		tx       *xsql.Tx
		err      error
	)
	if payOrder.Status == model.PayOrderStatusFail {
		log.Info("s.PayFail(%s, %+v) msg(DuplicatedCallback)", status, payOrder)
		return nil
	}
	//if ystOrder, err = s.YstOrderState(c, payOrder.OrderNo, payOrder.ThirdTradeNo); err != nil {
	//	log.Error("s.YstOrderState(%s, %s) err(%+v)", payOrder.OrderNo, payOrder.ThirdTradeNo, err)
	//	return err
	//}
	//if ystOrder.PayStatus != model.YstPayStatusPending {
	//	return ecode.TVIPYstRequestErr
	//}
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	payOrder.Status = model.PayOrderStatusFail
	if err = s.dao.TxUpdatePayOrder(c, tx, payOrder); err != nil {
		log.Error("s.dao.TxUpdatePayOrder(%v) err(%v)", payOrder, err)
		return err
	}
	if payOrder.OrderType == model.PayOrderTypeSub {
		return nil
	}
	if payParam, err = s.dao.CachePayParamByToken(c, payOrder.Token); err != nil {
		log.Info("s.dao.CachePayParamByToken(%s) err(%+v)", payOrder.Token, err)
		return err
	}
	payParam.Status = model.PayOrderStatusFail
	payParam.OrderNo = payOrder.OrderNo
	s.flushPayParamAsync(c, payOrder.Token, payParam)
	return nil
}

func (s *Service) paySuccess(c context.Context, status string, payOrder *model.PayOrder) error {
	log.Info("s.PayCallback.paySuccess(%s, %+v)", status, payOrder)
	var (
		tx       *xsql.Tx
		ui       *model.UserInfo
		uch      *model.UserChangeHistory
		uc       *model.UserContract
		payParam *model.PayParam
		panel    *model.PanelPriceConfig
		//ystOrder *model.YstOrderStateReply
		err      error
	)
	if payOrder.Status == model.PayOrderStatusSuccess {
		log.Info("s.PaySuccess(%s, %+v) msg(DuplicatedCallback)", status, payOrder)
		return nil
	}
	//ystOrderReq := &model.YstOrderStateReq{
	//	SeqNo:   payOrder.OrderNo,
	//	TraceNo: payOrder.ThirdTradeNo,
	//}
	//if ystOrder, err = s.dao.YstOrderState(c, ystOrderReq); err != nil {
	//	log.Error("s.dao.YstOrderState(%+v) err(%+v)", ystOrderReq, err)
	//	return err
	//}
	//if ystOrder.PayStatus != model.YstPayStatusPending {
	//	return ecode.TVIPYstRequestErr
	//}
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	payOrder.Status = model.PayOrderStatusSuccess
	if err = s.dao.TxUpdatePayOrder(c, tx, payOrder); err != nil {
		log.Error("s.dao.TxUpdatePayOrder(%v) err(%v)", payOrder, err)
		return err
	}
	if panel, err = s.PanelPriceConfigByProductId(c, payOrder.ProductId); err != nil {
		log.Info("s.dao.PanelPriceConfigByProductId(%s) err(%v)", payOrder.ProductId, err)
		return err
	}
	if ui, err = s.dao.RawUserInfoByMid(c, payOrder.Mid); err != nil {
		log.Info("s.dao.RawInfoByMid(%d) err(%v)", payOrder.Mid, err)
		return err
	}
	if ui == nil {
		if ui, err = s.initUserInfo(c, tx, payOrder.Mid); err != nil {
			log.Info("s.initUserInfo(%d) err(%+v)", payOrder.Mid, err)
			return err
		}
	}
	ui.CopyFromPayOrder(payOrder)
	ui.CopyFromPanel(panel)
	buyDuration := int64(time.Hour) * int64(payOrder.BuyMonths) * 24 * 31
	if err = s.incrVipDuration(c, tx, ui, time.Duration(buyDuration)); err != nil {
		log.Error("s.incrVipDuration(%v, %v) err(%v)", ui, buyDuration, err)
		return err
	}
	uch = &model.UserChangeHistory{}
	uch.CopyFromPayOrder(payOrder)
	if _, err = s.dao.TxInsertUserChangeHistory(c, tx, uch); err != nil {
		log.Error("s.dao.TxInsertUserChangeHistory(%v) err(%v)", uch, err)
		return err
	}
	// Note: send contract request to yst when a user buys contracted package with ali pay.
	// TODO: replace contract id by contract code
	if panel.IsContracted() && payOrder.PaymentType == model.PaymentTypeWechat {
		uc = &model.UserContract{
			Mid:     payOrder.Mid,
			OrderNo: payOrder.OrderNo,
			//ContractId: req.ContractId,
		}
		if _, err = s.dao.TxInsertUserContract(c, tx, uc); err != nil {
			log.Error("s.dao.TxInsertUserContract(%v) err(%v)", uc, err)
			return err
		}
	}

	if payOrder.OrderType == model.PayOrderTypeNormal {
		if payParam, err = s.dao.CachePayParamByToken(c, payOrder.Token); err != nil {
			log.Info("s.dao.CachePayParamByToken(%s) err(%+v)", payOrder.Token, err)
			return err
		}
		payParam.Status = model.PayOrderStatusSuccess
		payParam.OrderNo = payOrder.OrderNo
		s.flushPayParamAsync(c, payOrder.Token, payParam)
	}
	s.flushUserInfoAsync(c, payOrder.Mid)
	if panel.SuitType != model.SuitTypeMvip {
		// NOTE: 存在多送和少送的可能性
		// 多送：发起赠送请求，db commit 失败
		// 少送：发起赠送请求失败，db commit 成功
		s.giveMVipGiftAsync(c, payOrder.Mid, panel.PidOrId(), payOrder.OrderNo)
	}
	return nil
}

func (s *Service) PayCallback(c context.Context, req *model.YstPayCallbackReq) (res *model.YstPayCallbackReply) {
	var (
		err error
	)
	res = &model.YstPayCallbackReply{TraceNo: req.TraceNo}
	payOrder, err := s.dao.PayOrderByOrderNo(c, req.SeqNo)
	if err != nil {
		log.Error("s.dao.PayOrderByOrderNo(%s) err(%v)", req.SeqNo, err)
		res.Result = model.YstResultSysErr
		res.Msg = err.Error()
		return
	}
	pv := &validator.PayCallbackValidator{
		Signer:         s.dao.Signer(),
		CallbackReq:    req,
		PayOrder:       payOrder,
		ExpireDuration: s.c.PAY.PayExpireDuration,
	}
	if err = pv.Validate(); err != nil {
		log.Error("payCallbackValidator.Validate() err(%v)", err)
		res.Result = model.YstResultFail
		res.Msg = err.Error()
		return
	}
	payOrder.ThirdTradeNo = req.TraceNo
	switch req.TradeState {
	case model.YstTradeStateSuccess:
		err = s.paySuccess(c, req.TradeState, payOrder)
	case model.YstTradeStateClosed, model.YstTradeStatePayFail:
		err = s.payFail(c, req.TradeState, payOrder)
	default:
		err = ecode.TVIPYstUnknownTradeState
	}
	if err != nil {
		res.Result = model.YstResultSysErr
		res.Msg = err.Error()
	} else {
		res.Result = model.YstResultSuccess
		res.Msg = "ok"
	}
	return
}

func (s *Service) GiveMVipGift(c context.Context, mid int64, pid int32, orderNo string) error {
	batchId, ok := s.c.MVIP.BatchIdsMap[strconv.Itoa(int(pid))]
	if !ok {
		log.Error("s.c.MVIP.BatchIdsMap(%d) err(UnknownBatchId)", pid)
		return ecode.TVIPBatchIdNotFound
	}
	return s.dao.GiveMVipGift(c, mid, batchId, orderNo)
}

// TODO: query user info from tv_user_contract.
func (s *Service) UserInfoByContractCode(c context.Context, contractCode string) (ui *model.UserInfo, err error) {
	return nil, nil
}

func (s *Service) SignContract(c context.Context, contractCode string, contractId string, remark string) (err error) {
	var (
		uc *model.UserContract
		ui *model.UserInfo
		tx *xsql.Tx
	)
	if ui, err = s.dao.RawUserInfoByMid(c, uc.Mid); err != nil {
		return
	}
	if ui == nil {
		return ecode.RequestErr
	}
	if uc, err = s.dao.UserContractByContractId(c, contractId); err != nil {
		return
	}
	if uc != nil {
		return nil
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		err = s.dao.EndTran(tx, err)
	}()

	ui.PayType = model.VipPayTypeSub
	if err = s.dao.TxUpdateUserPayType(c, tx, ui); err != nil {
		log.Error("s.dao.TxUpdateUserPayType(%+v) err(%v)", ui, err)
		return
	}
	uc = &model.UserContract{
		Mid:        ui.Mid,
		ContractId: contractId,
	}
	if _, err = s.dao.TxInsertUserContract(c, tx, uc); err != nil {
		log.Error("s.dao.TxInsertUserContract(%d) err(%+v)", uc.ID, err)
		return
	}
	uch := &model.UserChangeHistory{
		Mid:        ui.Mid,
		ChangeType: model.UserChangeTypeSignContract,
		ChangeTime: xtime.Time(time.Now().Unix()),
		Remark:     remark,
	}
	if _, err = s.dao.TxInsertUserChangeHistory(c, tx, uch); err != nil {
		log.Error("s.dao.TxInsertUserChangeHistory(%+v) err(%+v)", uch, err)
		return
	}
	s.flushUserInfoAsync(c, ui.Mid)
	return nil
}

func (s *Service) CancelContract(c context.Context, contractCode string, contractId string, remark string) (err error) {
	var (
		uc *model.UserContract
		ui *model.UserInfo
		tx *xsql.Tx
	)
	if ui, err = s.dao.RawUserInfoByMid(c, uc.Mid); err != nil {
		return
	}
	if ui == nil {
		return ecode.RequestErr
	}
	if uc, err = s.dao.UserContractByContractId(c, contractId); err != nil {
		return
	}
	if uc == nil {
		log.Error("s.CancelContract(%d, %s, %s) err(UserContractNotFound)", ui.Mid, contractId, remark)
		return ecode.TVIPNotContracted
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		err = s.dao.EndTran(tx, err)
	}()
	ui.PayType = model.VipPayTypeNormal
	if err = s.dao.TxUpdateUserPayType(c, tx, ui); err != nil {
		log.Error("s.dao.TxUpdateUserPayType(%+v) err(%v)", ui, err)
		return
	}
	if err = s.dao.TxDeleteUserContract(c, tx, uc.ID); err != nil {
		log.Error("s.dao.TxDeleteUserContract(%d) err(%+v)", uc.ID, err)
		return
	}
	uch := &model.UserChangeHistory{
		Mid:        ui.Mid,
		ChangeType: model.UserChangeTypeCancelContract,
		ChangeTime: xtime.Time(time.Now().Unix()),
		Remark:     remark,
	}
	if _, err = s.dao.TxInsertUserChangeHistory(c, tx, uch); err != nil {
		log.Error("s.dao.TxInsertUserChangeHistory(%+v) err(%+v)", uch, err)
		return
	}
	s.flushUserInfoAsync(c, ui.Mid)
	return nil
}

func (s *Service) WxContractCallback(c context.Context, req *model.WxContractCallbackReq) (res *model.WxContractCallbackReply) {
	var err error
	if req.ChangeType == model.YST_CONTRACT_TYPE_SIGN {
		err = s.SignContract(c, req.ContractCode, req.ContractId, req.ContractTerminationMode)
	} else {
		err = s.CancelContract(c, req.ContractCode, req.ContractId, req.ContractTerminationMode)
	}
	res = &model.WxContractCallbackReply{
		ContractId: req.ContractId,
	}
	if err != nil {
		res.Result = model.YstResultSysErr
		res.Msg = err.Error()
	} else {
		res.Result = model.YstResultSuccess
		res.Msg = "ok"
	}
	return res
}
