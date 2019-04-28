package dao

import (
	"context"
	"strconv"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
)

var _defaultAdd = int64(1)

// PubView adds a view count.
func (d *Dao) PubView(c context.Context, mid int64, aid int64, ip string, cheat *artmdl.CheatInfo) (err error) {
	msg := &artmdl.StatMsg{
		Aid:       aid,
		Mid:       mid,
		IP:        ip,
		View:      &_defaultAdd,
		CheatInfo: cheat,
	}
	if err = d.statDbus.Send(c, strconv.FormatInt(aid, 10), msg); err != nil {
		PromError("databus:发送浏览")
		log.Error("d.databus.SendView(%+v) error(%+v)", msg, err)
		return
	}
	PromInfo("databus:发送浏览")
	log.Info("s.PubView(mid: %v, aid: %v, ip: %v, cheat: %+v)", msg.Mid, msg.Aid, msg.IP, cheat)
	return
}

// PubShare add share count
func (d *Dao) PubShare(c context.Context, mid int64, aid int64, ip string) (err error) {
	msg := &artmdl.StatMsg{
		Aid:   aid,
		Mid:   mid,
		IP:    ip,
		Share: &_defaultAdd,
	}
	if err = d.statDbus.Send(c, strconv.FormatInt(aid, 10), msg); err != nil {
		PromError("databus:发送分享")
		log.Error("d.databus.SendShare(%+v) error(%+v)", msg, err)
		return
	}
	PromInfo("databus:发送分享")
	log.Info("s.PubShare(mid: %v, aid: %v, ip: %v)", msg.Mid, msg.Aid, msg.IP)
	return
}
