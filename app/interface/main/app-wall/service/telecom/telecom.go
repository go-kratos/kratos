package telecom

import (
	"bytes"
	"context"
	"crypto/des"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"go-common/app/interface/main/app-wall/model/telecom"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_sendSMSCaptcha = `{"captcha":"%v"}`
	_sendSMSFlow    = `{"flow":"%v"}`
	_flowPackageID  = 279
)

// InOrdersSync insert OrdersSync
func (s *Service) InOrdersSync(c context.Context, ip string, u *telecom.TelecomOrderJson) (err error) {
	if !s.iplimit(_telecomKey, ip) {
		err = ecode.AccessDenied
		return
	}
	if u == nil || u.Detail == nil {
		err = ecode.NothingFound
		return
	}
	var (
		result    int64
		requestNo int
		phoneStr  string
		detail    = u.Detail
	)
	detail.TelecomJSONChange()
	requestNo, _ = strconv.Atoi(u.RequestNo)
	if detail.PhoneID == "" {
		if phoneStr, err = s.dao.PayPhone(c, int64(requestNo)); err != nil {
			log.Error("telecom_s.dao.PayPhone error(%v)", err)
			return
		}
	} else {
		phoneStr = detail.PhoneID
	}
	if result, err = s.dao.InOrderSync(c, requestNo, u.ResultType, phoneStr, detail); err != nil || result == 0 {
		log.Error("telecom_s.dao.OrdersSync (%v) error(%v) or result==0", u, err)
		return
	}
	if detail.OrderStatus == 3 {
		phoneInt, _ := strconv.Atoi(phoneStr)
		if err = s.dao.SendTelecomSMS(c, phoneInt, s.smsOrderTemplateOK); err != nil {
			log.Error("telecom_s.dao.SendTelecomSMS error(%v)", err)
			return
		}
	}
	return
}

// InRechargeSync insert RechargeSync
func (s *Service) InRechargeSync(c context.Context, ip string, u *telecom.RechargeJSON) (err error) {
	if !s.iplimit(_telecomKey, ip) {
		err = ecode.AccessDenied
		return
	}
	var result int64
	if result, err = s.dao.InRechargeSync(c, u); err != nil || result == 0 {
		log.Error("telecom_s.dao.InRechargeSync (%v) error(%v) or result==0", u, err)
		return
	}
	return
}

// telecomMessageOrder
func (s *Service) TelecomMessageSync(c context.Context, ip string, u *telecom.TelecomMessageJSON) (err error) {
	if !s.iplimit(_telecomKey, ip) {
		err = ecode.AccessDenied
		return
	}
	if u == nil {
		err = ecode.ServerErr
		return
	}
	phoneInt, _ := strconv.Atoi(u.PhoneID)
	if err = s.dao.SendTelecomSMS(c, phoneInt, s.smsMsgTemplate); err != nil {
		log.Error("telecom_s.dao.SendTelecomSMS error(%v)", err)
		return
	}
	return
}

// telecomInfo
func (s *Service) telecomInfo(c context.Context, phone int) (res map[int]*telecom.OrderInfo, err error) {
	var (
		row *telecom.OrderInfo
	)
	res = map[int]*telecom.OrderInfo{}
	if row, err = s.dao.TelecomCache(c, phone); err == nil && row != nil {
		res[phone] = row
		s.pHit.Incr("telecom_cache")
		return
	}
	if res, err = s.dao.OrdersUserFlow(c, phone); err != nil {
		log.Error("telecom_s.dao.OrdersUserFlow phone (%v) error(%v)", phone, err)
		return
	}
	s.pMiss.Incr("telecom_cache")
	if user, ok := res[phone]; ok {
		if err = s.dao.AddTelecomCache(c, phone, user); err != nil {
			log.Error("telecom_s.dao.AddTelecomCache error(%v)", err)
			return
		}
	}
	return
}

// telecomInfoByOrderID
func (s *Service) telecomInfoByOrderID(c context.Context, orderID int64) (res map[int64]*telecom.OrderInfo, err error) {
	var (
		row *telecom.OrderInfo
	)
	res = map[int64]*telecom.OrderInfo{}
	if row, err = s.dao.TelecomOrderIDCache(c, orderID); err == nil && row != nil {
		res[orderID] = row
		s.pHit.Incr("telecom_orderid_cache")
		return
	}
	if res, err = s.dao.OrdersUserByOrderID(c, orderID); err != nil {
		log.Error("telecom_s.dao.OrdersUserByOrderID phone (%v) error(%v)", orderID, err)
		return
	}
	s.pMiss.Incr("telecom_orderid_cache")
	if user, ok := res[orderID]; ok {
		if err = s.dao.AddTelecomOrderIDCache(c, orderID, user); err != nil {
			log.Error("telecom_s.dao.AddTelecomOrderIDCache error(%v)", err)
			return
		}
	}
	return
}

// TelecomPay
func (s *Service) TelecomPay(c context.Context, phone, isRepeatOrder, payChannel, payAction int, orderID int64, ipStr string) (res *telecom.Pay, msg string, err error) {
	var (
		requestNo         int64
		t                 map[int]*telecom.OrderInfo
		rcode             int
		beginTime         time.Time
		firstOrderEndtime time.Time
	)
	if requestNo, err = s.seqdao.SeqID(c); err != nil {
		log.Error("telecom_s.seqdao.SeqID error (%v)", err)
		return
	}
	if t, err = s.telecomInfo(c, phone); err != nil {
		log.Error("telecom_s.telecomInfo phone(%v) error (%v)", phone, err)
		return
	}
	if user, ok := t[phone]; ok {
		beginTime = user.Begintime.Time()
		firstOrderEndtime = user.Endtime.Time()
	}
	if res, err, msg = s.dao.PayInfo(c, requestNo, phone, isRepeatOrder, payChannel, payAction, orderID, ipStr, beginTime, firstOrderEndtime); err != nil || res == nil {
		log.Error("telecom_s.dao.PayInfo requestNo (%v) phone (%v) isRepeatOrder (%v) payChannel (%v) payAction (%v) t.OrderID (%v) ipStr (%v) error (%v)",
			requestNo, phone, isRepeatOrder, payChannel, payAction, orderID, ipStr, err)
		return
	}
	phoneStr := strconv.Itoa(phone)
	if rcode, err = s.dao.AddPayPhone(c, requestNo, phoneStr); err != nil || rcode != 1 {
		log.Error("telecom_s.dao.AddPayPhone error (%v)", err)
		return
	}
	return
}

// CancelRepeatOrder
func (s *Service) CancelRepeatOrder(c context.Context, phone int) (msg string, err error) {
	var (
		res map[int]*telecom.OrderInfo
	)
	res, err = s.telecomInfo(c, phone)
	if err != nil {
		log.Error("telecom_s.telecomInfo phone(%v) error (%v)", phone, err)
		return
	}
	user, ok := res[phone]
	if !ok {
		err = ecode.NothingFound
		msg = "订单不存在"
		return
	}
	if msg, err = s.dao.CancelRepeatOrder(c, phone, user.SignNo); err != nil {
		log.Error("telecom_s.dao.CancelRepeatOrder phone(%v) signNo(%v) error (%v)", phone, user.SignNo, err)
		return
	}
	return
}

// OrderList user order list
func (s *Service) OrderList(c context.Context, orderID int64, phone int) (res *telecom.SucOrder, msg string, err error) {
	if res, err, msg = s.dao.SucOrderList(c, phone); err != nil || res == nil {
		log.Error("telecom_s.dao.SucOrderList orderID (%v) phone (%v) error (%v)", orderID, phone, err)
		return
	}
	return
}

// PhoneFlow user flow
func (s *Service) PhoneFlow(c context.Context, orderID int64, phone int) (res *telecom.OrderFlow, msg string, err error) {
	var t *telecom.SucOrder
	if t, err, msg = s.dao.SucOrderList(c, phone); err != nil || t == nil {
		log.Error("telecom_s.dao.SucOrderList orderID (%v) phone (%v) error (%v)", orderID, phone, err)
		return
	}
	res = &telecom.OrderFlow{
		FlowBalance: t.FlowBalance,
	}
	if (t.FlowBalance/t.FlowPackageSize)*100 < s.flowPercentage {
		flow := strconv.Itoa(t.FlowBalance)
		dataJSON := fmt.Sprintf(_sendSMSFlow, flow)
		if err = s.dao.SendSMS(c, phone, s.smsFlowTemplate, dataJSON); err != nil {
			log.Error("telecom_s.dao.SendSMS phone(%v) error (%v)", phone, err)
			return
		}
		msg = "免流量量剩余不不⾜足10%，超出部分会按正常流量量资费计费"
		return
	}
	return
}

// OrderConsent user orderConsent
func (s *Service) OrderConsent(c context.Context, phone int, orderID int64, captcha string) (res *telecom.PhoneConsent, msg string, err error) {
	var (
		area       string
		t          *telecom.SucOrder
		captchaStr string
		order      *telecom.OrderPhoneState
	)
	captchaStr, err = s.dao.PhoneCode(c, phone)
	if err != nil {
		log.Error("telecom_s.dao.PhoneCode error (%v)", err)
		msg = "验证码已过期"
		return
	}
	if captchaStr == "" || captchaStr != captcha {
		err = ecode.NotModified
		msg = "验证码错误"
		return
	}
	res = &telecom.PhoneConsent{
		Consent: 0,
	}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() error {
		if area, err, msg = s.dao.PhoneArea(ctx, phone); err != nil {
			log.Error("telecom_s.dao.PhoneArea phone(%v) error (%v)", phone, err)
			return err
		}
		return nil
	})
	g.Go(func() error {
		if t, err, _ = s.dao.SucOrderList(ctx, phone); err != nil {
			if err == ecode.NothingFound {
				err = nil
			} else {
				log.Error("telecom_s.dao.SucOrderList sphone (%v) error (%v)", phone, err)
			}
			return err
		}
		return nil
	})
	if orderID > 0 {
		g.Go(func() error {
			if order, err = s.dao.OrderState(ctx, orderID); err != nil {
				log.Error("telecom_s.dao.OrderState phone(%v) orderID(%v) error(%v)", phone, orderID, err)
				return err
			}
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("telecom_errgroup.WithContext error(%v)", err)
		return
	}
	if len(area) >= 4 {
		a := area[:4]
		if _, ok := s.telecomArea[a]; !ok {
			return
		}
	}
	if order != nil && order.OrderState == 2 && order.FlowPackageID == _flowPackageID {
		res.Consent = 3
		return
	}
	if t != nil && t.OrderID > 0 {
		res.Consent = 2
		return
	}
	res.Consent = 1
	return
}

// PhoneCode
func (s *Service) PhoneCode(c context.Context, phone int, captcha string, now time.Time) (res *telecom.Pay, err error, msg string) {
	var (
		captchaStr string
		order      map[int]*telecom.OrderInfo
	)
	captchaStr, err = s.dao.PhoneCode(c, phone)
	if err != nil {
		log.Error("telecom_s.dao.PhoneCode error (%v)", err)
		msg = "验证码已过期"
		return
	}
	if captchaStr == "" || captchaStr != captcha {
		err = ecode.NotModified
		msg = "验证码错误"
		return
	}
	if order, err = s.telecomInfo(c, phone); err != nil {
		log.Error("telecom_s.telecomInfo phone(%v) error (%v)", phone, err)
		return
	}
	user, ok := order[phone]
	if !ok {
		err = ecode.NothingFound
		msg = "订单不存在"
		return
	}
	switch user.OrderState {
	case 2:
		err = ecode.NotModified
		msg = "订购中"
	case 3:
		if now.Unix() > int64(user.Endtime) {
			err = ecode.NotModified
			msg = "订单已过期"
		} else {
			res = &telecom.Pay{
				OrderID: user.OrderID,
			}
			msg = "激活成功，麻麻再也不不⽤用担⼼心我的流量量了了(/≧▽≦)/"
		}
	case 4:
		err = ecode.NotModified
		msg = "订单失败"
	case 5:
		err = ecode.NotModified
		msg = "订单已过期"
	}
	return
}

// PhoneSendSMS
func (s *Service) PhoneSendSMS(c context.Context, phone int) (err error) {
	var (
		captcha string
		rcode   int
	)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		captcha = captcha + strconv.Itoa(r.Intn(10))
	}
	dataJSON := fmt.Sprintf(_sendSMSCaptcha, captcha)
	if err = s.dao.SendSMS(c, phone, s.smsTemplate, dataJSON); err != nil {
		log.Error("telecom_s.dao.SendSMS error (%v)", err)
		return
	}
	if rcode, err = s.dao.AddPhoneCode(c, phone, captcha); err != nil || rcode != 1 {
		log.Error("telecom_s.dao.AddPhoneCode error (%v)", err)
		return
	}
	return
}

// OrderState
func (s *Service) OrderState(c context.Context, orderid int64) (t *telecom.OrderState, msg string, err error) {
	var (
		orderState *telecom.OrderPhoneState
		torder     map[int64]*telecom.OrderInfo
		userOrder  *telecom.OrderInfo
		tSucOrder  *telecom.SucOrder
		ok         bool
	)
	g, ctx := errgroup.WithContext(c)
	log.Error("userOrder.PhoneID_test")
	g.Go(func() error {
		if orderState, err = s.dao.OrderState(ctx, orderid); err != nil {
			log.Error("telecom_s.dao.OrderState orderID error(%v)", orderid, err)
			return err
		}
		return nil
	})
	g.Go(func() error {
		if torder, err = s.telecomInfoByOrderID(c, orderid); err != nil {
			log.Error("telecom_s.telecomInfoByOrderID error(%v)", orderid, err)
			return err
		}
		if userOrder, ok = torder[orderid]; ok {
			if tSucOrder, err, _ = s.dao.SucOrderList(ctx, userOrder.PhoneID); err != nil {
				if err == ecode.NothingFound {
					err = nil
				} else {
					log.Error("telecom_s.dao.SucOrderList sphone (%v) error (%v)", userOrder.PhoneID, err)
				}
				return err
			}
		}
		return nil
	})
	if err = g.Wait(); err != nil {
		log.Error("telecom_errgroup.WithContext error(%v)", err)
		return
	}
	if orderState == nil {
		t = &telecom.OrderState{
			OrderState: 1,
		}
		return
	}
	t = &telecom.OrderState{
		OrderState: orderState.OrderState,
		FlowSize:   orderState.FlowSize,
	}
	if tSucOrder != nil && tSucOrder.OrderID > 0 {
		t.FlowBalance = tSucOrder.FlowBalance
		t.FlowSize = tSucOrder.FlowPackageSize
	}
	if ok {
		t.IsRepeatorder = userOrder.IsRepeatorder
		if userOrder.IsRepeatorder == 0 {
			t.Endtime = userOrder.Endtime
		}
	}
	switch t.OrderState {
	case 6:
		if !ok {
			err = ecode.NothingFound
			msg = "订单不存在"
			return
		}
		t.Endtime = userOrder.Endtime
	}
	return
}

// telecomIp ip limit
func (s *Service) iplimit(k, ip string) bool {
	key := fmt.Sprintf(_initIPlimitKey, k, ip)
	if _, ok := s.operationIPlimit[key]; ok {
		return true
	}
	return false
}

// DesDecrypt
func (s *Service) DesDecrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = s.zeroUnPadding(out)
	return out, nil
}

// zeroUnPadding
func (s *Service) zeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
