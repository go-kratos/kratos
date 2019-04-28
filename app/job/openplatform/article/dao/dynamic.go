package dao

import (
	"context"
	"strconv"

	"go-common/app/job/openplatform/article/model"
	"go-common/library/log"
)

const _dynamicArt = 64

// PubDynamic pub dynamic
func (d *Dao) PubDynamic(c context.Context, mid int64, aid int64, show bool, comment string, ts int64, dynamicIntro string) (err error) {
	msg := &model.DynamicMsg{}
	msg.Card.Type = _dynamicArt
	msg.Card.Rid = aid
	msg.Card.OwnerID = mid
	if show {
		msg.Card.Show = 1
	}
	msg.Card.Comment = comment
	msg.Card.Ts = ts
	msg.Card.Dynamic = dynamicIntro
	if err = d.dynamicDbus.Send(c, strconv.FormatInt(aid, 10), msg); err != nil {
		PromError("dynamic:发送动态消息")
		log.Error("dynamic: d.SendPubDynamic(%+v) error(%+v)", msg, err)
		return
	}
	PromInfo("databus:发送动态消息")
	log.Info("dynamic: dao.PubDynamic(%+v)", msg)
	return
}
