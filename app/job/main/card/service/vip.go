package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/card/model"
	cardapi "go-common/app/service/main/card/api/grpc/v1"
	cardmol "go-common/app/service/main/card/model"
	vipmol "go-common/app/service/main/vip/model"
	"go-common/library/log"
)

// ChangeEquipTime change vip equip time.
func (s *Service) ChangeEquipTime(c context.Context, v *model.VipReq) (err error) {
	var res *cardapi.UserCardReply
	if res, err = s.cardRPC.UserCard(c, &cardapi.UserCardReq{Mid: v.Mid}); err != nil {
		return
	}
	if res.Res == nil ||
		res.Res.Id == 0 ||
		res.Res.CardType != cardmol.CardTypeVip {
		return
	}
	var expire int64
	switch {
	case v.VipType == vipmol.NotVip || v.VipStatus == vipmol.Expire:
		expire = time.Now().Unix()
	case v.VipOverdueTime != res.Res.ExpireTime:
		expire = v.VipOverdueTime
	default:
	}
	if expire == 0 {
		return
	}
	if err = s.dao.UpdateExpireTime(c, expire, v.Mid); err != nil {
		return
	}
	err = s.dao.DelCacheEquip(c, v.Mid)
	return
}

func (s *Service) vipchangeproc() {
	defer s.waiter.Done()
	msgs := s.vipConsumer.Messages()
	var err error
	for {
		msg, ok := <-msgs
		if !ok {
			log.Warn("[service.dataConsume|vip] dataConsumer has been closed.")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%+v)", err)
		}
		log.Info("cur consumer vipchangeproc(%v)", string(msg.Value))
		v := &model.MsgCanal{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		if v.Table != _tableUserInfo || v.Action != _updateAction {
			continue
		}
		n := new(model.VipUserInfoMsg)
		if err = json.Unmarshal(v.New, n); err != nil {
			log.Error("vipchangeproc json.Unmarshal val(%v) error(%v)", string(v.New), err)
			continue
		}
		o := new(model.VipUserInfoMsg)
		if err = json.Unmarshal(v.Old, o); err != nil {
			log.Error("vipchangeproc json.Unmarshal val(%v) error(%v)", string(v.Old), err)
			continue
		}
		if n.VipStatus == o.VipStatus &&
			n.VipType == o.VipType &&
			n.VipOverdueTime == o.VipOverdueTime {
			continue
		}
		var duetime time.Time
		if duetime, err = time.ParseInLocation("2006-01-02 15:04:05", n.VipOverdueTime, time.Local); err != nil {
			log.Error("vipchangeproc ParseInLocation val(%s) error(%v)", n.VipOverdueTime, err)
			continue
		}
		if err = s.ChangeEquipTime(context.Background(), &model.VipReq{
			Mid:            n.Mid,
			VipType:        n.VipType,
			VipStatus:      n.VipStatus,
			VipOverdueTime: duetime.Unix(),
		}); err != nil {
			log.Error("ChangeEquipTime val(%+v) error(%v)", n, err)
			continue
		}
	}
}
