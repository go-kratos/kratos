package service

import (
	"context"
	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"time"
)

// dealAddPay fn
func (s *Service) dealAddPay(c context.Context, pay *archive.Pay, aid, mid int64, ip string) (err error) {
	// step 1: call API
	if err = s.pay.AssReg(c, mid, aid, pay.Price, ip); err != nil {
		log.Error("s.pay.AssReg mid(%d)|aid(%d)|pay(%+v) error(%v)", mid, aid, pay, err)
		return
	}
	// step 2: add protocol user log
	var index = []interface{}{mid, pay.ProtocolID}
	uInfo := &report.UserInfo{
		Mid:      mid,
		Business: archive.UgcpayAddarcProtocol,
		Type:     1,
		Oid:      aid,
		Action:   "add",
		Ctime:    time.Now(),
		IP:       ip,
		Index:    index,
	}
	uInfo.Content = map[string]interface{}{
		"content":     pay,
		"protocol_id": pay.ProtocolID,
		"mid":         mid,
	}
	report.User(uInfo)
	log.Warn("sendLog dealAddPay protocol info (%+v)", uInfo)
	return
}

// dealAdjustPay fn
func (s *Service) dealAdjustPay(c context.Context, pay *archive.Pay, aid, mid int64, ip string) (err error) {
	// step 1: call API
	if err = s.pay.AssReg(c, mid, aid, pay.Price, ip); err != nil {
		log.Error("s.pay.AssReg mid(%d)|aid(%d)|pay(%+v) error(%v)", mid, aid, pay, err)
		return
	}
	// step 2: add protocol user log
	var index = []interface{}{mid, pay.ProtocolID}
	uInfo := &report.UserInfo{
		Mid:      mid,
		Business: archive.UgcpayAddarcProtocol,
		Type:     1,
		Oid:      aid,
		Action:   "edit",
		Ctime:    time.Now(),
		IP:       ip,
		Index:    index,
	}
	uInfo.Content = map[string]interface{}{
		"content":     pay,
		"protocol_id": pay.ProtocolID,
		"mid":         mid,
	}
	report.User(uInfo)
	log.Warn("sendLog dealAdjustPay protocol info (%+v)", uInfo)
	return
}

// 参与UGC付费之前必须接受当前最新的投稿协议
func (s *Service) checkPayProtocol(c context.Context, pay *archive.Pay, mid int64) (err error) {
	if pay != nil {
		accept, _ := s.pay.UserAcceptProtocol(c, pay.ProtocolID, mid)
		if accept {
			return
		}
		if !accept && pay.ProtocolAccept == 0 {
			log.Error("s.rejectUgcProtocolBefore (%+v),(%+v)", pay, err)
			err = ecode.VideoupPayProtocolLimit
			return
		}
	}
	return
}

// 付费提交校验
func (s *Service) checkPayLimit(c context.Context, ap *archive.ArcParam) (err error) {
	pay := ap.Pay
	if pay != nil {
		if _, ok := s.exemptUgcPayUps[ap.Mid]; !ok { // 用户灰度
			log.Error("s.checkAddPayLimit VideoupPayUserNotAllow (%d),(%+v)", ap.Mid, err)
			err = ecode.VideoupPayUserNotAllow
			return
		}
		if pay.Open == 1 {
			if ap.Copyright != archive.CopyrightOriginal { // 创作类型
				log.Error("s.checkAddPayLimit VideoupPayCopyrightErr (%+v),(%+v)", ap, err)
				err = ecode.VideoupPayCopyrightErr
				return
			}
			if pay.Price > 1000 || pay.Price < 1 { //开启之后的定价必须合理有效
				log.Error("s.checkAddPayLimit VideoupPayPriceErr (%+v),(%+v)", ap, err)
				err = ecode.VideoupPayPriceErr
				return
			}
			ap.UgcPay = 1
		} else {
			ap.UgcPay = 0
		}
	}
	return
}

// 一起检测是否和商单以及私单冲突，只能三选一
func (s *Service) checkPayWithOrder(c context.Context, porder *archive.Porder, pay *archive.Pay, orderID, mid int64) (err error) {
	joinPorder := porder != nil && porder.FlowID > 0
	joinPay := pay != nil && pay.Open == 1
	if joinPay && (joinPorder || orderID > 0) {
		log.Error("s.checkAddPayWithOrder VideoupPayCommericalLimit (%d)|(%+v)|(%+v)|(%+v),(%+v)", mid, porder, pay, orderID, err)
		err = ecode.VideoupPayCommericalLimit
		return
	}
	return
}
