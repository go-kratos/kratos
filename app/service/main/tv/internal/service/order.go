package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/service/main/account/api"
	"go-common/app/service/main/tv/internal/model"
	"go-common/app/service/main/tv/internal/service/validator"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// makeOrderNo get order id
func (s *Service) makeOrderNo() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%05d", s.r.Int63n(99999)))
	b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
	b.WriteString(time.Now().Format("060102150405"))
	return b.String()
}

func (s *Service) createOrder(c context.Context, mid int64, token string, platform int8, paymentType string, clientIp string) (pi *model.PayInfo, err error) {
	var (
		payParam     *model.PayParam
		panel        *model.PanelPriceConfig
		mvipPanels   []*model.PanelPriceConfig
		mvip         *model.MainVip
		payOrder     *model.PayOrder
		payingOrders []*model.PayOrder
		account      *api.Info
		ui           *model.UserInfo
		tx           *xsql.Tx
	)
	if payParam, err = s.dao.CachePayParamByToken(c, token); payParam == nil || err != nil {
		if err == nil {
			err = ecode.TVIPTokenErr
		}
		log.Error("s.dao.CachePayParamByToken(%s) err(%v)", token, err)
		return
	}
	if panel, err = s.PanelPriceConfigByPid(c, payParam.Pid); panel == nil || err != nil {
		if err == nil {
			err = ecode.TVIPPanelNotFound
		}
		log.Error("s.PanelPriceConfigByPid(%d) err(%v)", payParam.Pid, err)
		return
	}
	if mid == -1 {
		mid = payParam.Mid
	}
	key := fmt.Sprintf("LOCK:ORDER:%d", mid)
	val := strconv.Itoa(time.Now().Nanosecond())
	if err = s.dao.Lock(c, key, val); err != nil {
		log.Error("s.dao.Lock(%s, %s) err(%+v)", key, err)
		return
	}
	defer func() {
		if err := s.dao.Unlock(c, key, val); err != nil {
			log.Error("s.dao.Unlock(%s, %s) err(%+v)", key, val, err)
		}
	}()
	if mvip, err = s.dao.MainVip(c, mid); err != nil {
		log.Error("s.dao.MainVip(%d) err(%v)", mid, err)
		return
	}
	if payOrder, err = s.dao.PayOrderByOrderNo(c, payParam.OrderNo); err != nil {
		log.Error("s.dao.PayOrderByOrderNo(%s) err(%v)", payParam.OrderNo, err)
		return
	}
	if ui, err = s.dao.RawUserInfoByMid(c, mid); err != nil {
		log.Error("s.dao.RawUserInfoByMid(%d) err(%+v)", mid, err)
		return
	}
	from := xtime.Time(time.Now().Add(-time.Duration(s.c.PAY.OrderRateFromDuration)).Unix())
	to := xtime.Time(time.Now().Unix())
	if mvipPanels, err = s.PanelPriceConfigsBySuitType(c, model.SuitTypeMvip); err != nil {
		log.Error("s.PanelPriceConfigsBySuitType(%d) err(%+v)", model.SuitTypeMvip, err)
		return
	}
	if payingOrders, _, err = s.dao.PayOrdersByMidAndStatusAndCtime(c, mid, model.PayOrderStatusPaying, from, to, 1, 100); err != nil {
		log.Error("s.dao.PayOrdersByMidAndStatusAndCtime(%d, %d, %v, %v) err(%+v)", mid, model.PayOrderStatusPaying, from, to, err)
		return
	}
	log.Info("payingOrders mid: %d from: %+v to: %+v orders: %+v", mid, from, to, payingOrders)
	cv := &validator.CreateOrderValidator{
		PayOrder:          payOrder,
		PaymentType:       paymentType,
		PayParam:          payParam,
		MVip:              mvip,
		Panel:             panel,
		Ui:                ui,
		PayingOrders:      payingOrders,
		MVipPanels:        mvipPanels,
		MVipRateMaxNumber: s.c.PAY.OrderRateMaxNumber,
	}
	if err = cv.Validate(); err != nil {
		log.Error("cv.Validate err(%v)", err)
		return
	}
	if payParam.OrderNo == "" {
		payParam.OrderNo = s.makeOrderNo()
	}
	payOrder = &model.PayOrder{
		Platform:    platform,
		Mid:         mid,
		Status:      1,
		PaymentType: paymentType,
		Ver:         1,
		Token:       token,
	}
	payOrder.CopyFromPayParam(payParam)
	payOrder.CopyFromPanel(panel)
	ystCreateOrderReq := &model.YstCreateOrderReq{GUID: payParam.Guid, ClientIp: clientIp}
	ystCreateOrderReq.CopyFromPayOrder(payOrder)
	if account, err = s.dao.AccountInfo(c, mid); err != nil {
		return
	}
	ystCreateOrderReq.CopyFromAccount(account)
	ystOrder, err := s.dao.CreateYstOrder(c, ystCreateOrderReq)
	if err != nil {
		log.Error("s.dao.CreateYstOrder(%+v) err(%+v)", ystCreateOrderReq, err)
		return
	}
	payParam.Status = model.PayOrderStatusPaying
	payParam.OrderNo = payOrder.OrderNo
	payOrder.ThirdTradeNo = ystOrder.TraceNo
	tx, err = s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran() err(%+v)", err)
		return
	}
	defer func() {
		s.dao.EndTran(tx, err)
	}()
	if _, err = s.dao.TxInsertPayOrder(c, tx, payOrder); err != nil {
		log.Error("s.dao.TxInsertPayOrder(%+v) err(%+v)", payOrder, err)
		return
	}
	s.flushPayParamAsync(c, token, payParam)
	pi = &model.PayInfo{CodeUrl: ystOrder.CodeUrl}
	pi.CopyFromPayOrder(payOrder)
	return pi, nil
}

func (s *Service) CreateOrder(c context.Context, token string, platform int8, paymentType string, clientIp string) (pi *model.PayInfo, err error) {
	return s.createOrder(c, -1, token, platform, paymentType, clientIp)
}

func (s *Service) CreateGuestOrder(c context.Context, mid int64, token string, platform int8, paymentType string, clientIp string) (pi *model.PayInfo, err error) {
	return s.createOrder(c, mid, token, platform, paymentType, clientIp)
}

func (s *Service) MakeUpOrderStatus() error {
	var (
		err error
		res []*model.PayOrder
		std time.Duration
		etd time.Duration
	)
	c := context.TODO()
	now := time.Now()
	if std, err = time.ParseDuration(s.c.Ticker.UnpaidDurationEtime); err != nil {
		return err
	}
	stime := now.Add(-std)
	if etd, err = time.ParseDuration(s.c.Ticker.UnpaidDurationStime); err != nil {
		return err
	}
	etime := now.Add(-etd)
	if res, err = s.dao.UnpaidNotCallbackOrder(c, xtime.Time(stime.Unix()), xtime.Time(etime.Unix()), 1, 500); err != nil {
		log.Error("s.dao.UnpaidNoCallbackOrder err(%v)", err)
		return err
	}

	for _, payOrder := range res {
		var (
			ystOrder *model.YstOrderStateReply
		)
		ystOrderReq := &model.YstOrderStateReq{
			SeqNo:   payOrder.OrderNo,
			TraceNo: payOrder.ThirdTradeNo,
		}
		if ystOrder, err = s.dao.YstOrderState(c, ystOrderReq); err != nil {
			log.Error("s.dao.YstOrderState(%+v) err(%+v)", ystOrderReq, err)
			continue
		}
		if ystOrder.PayStatus == model.YstPayStatusPaied {
			s.mission(func() {
				if err = s.paySuccess(c, ystOrder.PayStatus, payOrder); err != nil {
					log.Error("s.MakeUpPaySuccessStatus(%s, %+v) err(%+v)", ystOrder.PayStatus, payOrder, err)
				}
			})
		} else if ystOrder.PayStatus == model.YstPayStatusPending {
			s.mission(func() {
				if err = s.payFail(c, ystOrder.PayStatus, payOrder); err != nil {
					log.Error("s.MakeUpPayFailStatus(%s, %+v) err(%+v)", ystOrder.PayStatus, payOrder, err)
				}
			})

		}
	}
	return err
}
