package dao

import (
	"context"
	"strconv"
	"time"

	"go-common/app/service/main/thumbup/model"
	"go-common/library/log"
)

// PubStatDatabus pub stat databus
func (d *Dao) PubStatDatabus(c context.Context, business string, mid int64, s *model.Stats, upMid int64) (err error) {
	msg := &model.StatMsg{Type: business, ID: s.ID, Count: s.Likes, Timestamp: time.Now().Unix(), OriginID: s.OriginID, DislikeCount: s.Dislikes, Mid: mid, UpMid: upMid}
	if err = d.statDbus.Send(c, strconv.FormatInt(s.ID, 10), msg); err != nil {
		log.Error("d.databus.Send error(%v)", err)
		PromError("databus:stat")
		return
	}
	log.Info("s.PubStatDatabus (%+v)", msg)
	return
}

// PubLikeDatabus .
func (d *Dao) PubLikeDatabus(c context.Context, p *model.LikeMsg) (err error) {
	if err = d.likeDbus.Send(c, strconv.FormatInt(p.Mid, 10), p); err != nil {
		log.Error("d.likeDbus.Send error(%v)", err)
		PromError("databus:like")
		return
	}
	log.Info("s.PubLikeDatabus success (%+v)", p)
	return
}

// PubItemMsg .
func (d *Dao) PubItemMsg(c context.Context, business string, originID, messageID int64, state int8) (err error) {
	msg := &model.ItemMsg{
		State:     state,
		Business:  business,
		OriginID:  originID,
		MessageID: messageID,
	}
	if err = d.itemDbus.Send(c, strconv.FormatInt(messageID, 10), msg); err != nil {
		log.Error("d.PubItemMsg.databus.Send error(%v)", err)
		PromError("databus:item")
		return
	}
	log.Info("s.PubItemMsg success (%+v)", msg)
	return
}

// PubUserMsg .
func (d *Dao) PubUserMsg(c context.Context, business string, mid int64, state int8) (err error) {
	msg := &model.UserMsg{
		Mid:      mid,
		State:    state,
		Business: business,
	}
	if err = d.userDbus.Send(c, strconv.FormatInt(mid, 10), msg); err != nil {
		log.Error("d.PubUserMsg.databus.Send error(%v)", err)
		PromError("databus:user")
		return
	}
	log.Info("s.PubUserMsg success (%+v)", msg)
	return
}

// AddCacheUserLikeList .
func (d *Dao) AddCacheUserLikeList(c context.Context, mid int64, miss []*model.ItemLikeRecord, businessID int64, state int8) (err error) {
	err = d.PubUserMsg(c, d.BusinessIDMap[businessID].Name, mid, state)
	return
}
