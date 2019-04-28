package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"
)

func (s *Service) autoRenewPay(c context.Context, n *model.VipUserInfoMsg, o *model.VipUserInfoMsg) (err error) {
	// wechat autorenew 扣款重放条件：
	// 1.自动续费用户
	// 2.非IAP支付方式
	// 3.用户状态发生变化
	// 4.旧状态为未过期
	// 5.新状态为过期
	// 6.新旧类型都是不是NotVip type
	if n.IsAutoRenew == model.AutoRenew &&
		o.IsAutoRenew == model.AutoRenew &&
		n.PayChannelID != model.IAPChannelID &&
		o.PayChannelID != model.IAPChannelID &&
		n.Status != o.Status &&
		o.Status == model.VipStatusNotOverTime &&
		n.Status == model.VipStatusOverTime &&
		n.Type != model.NotVip && o.Type != model.NotVip {
		_, err = s.dao.AutoRenewPay(c, n.Mid)
	}
	return
}

func (s *Service) retryautorenewpayproc() {
	defer s.waiter.Done()
	msgs := s.autoRenewdDatabus.Messages()
	var err error
	for {
		msg, ok := <-msgs
		if !ok {
			log.Warn("[service.retryautorenewpayproc|vip] dataConsumer has been closed.")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("retryautorenewpayproc msg.Commit err(%+v)", err)
		}
		log.Info("cur consumer retryautorenewpayproc(%v)", string(msg.Value))
		v := &model.Message{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("retryautorenewpayproc json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		if v.Table != _tableUserInfo || v.Action != _updateAction {
			continue
		}
		n := new(model.VipUserInfoMsg)
		if err = json.Unmarshal(v.New, n); err != nil {
			log.Error("retryautorenewpayproc json.Unmarshal val(%v) error(%v)", string(v.New), err)
			continue
		}
		o := new(model.VipUserInfoMsg)
		if err = json.Unmarshal(v.Old, o); err != nil {
			log.Error("retryautorenewpayproc json.Unmarshal val(%v) error(%v)", string(v.Old), err)
			continue
		}
		s.autoRenewPay(context.Background(), n, o)
	}
}
