package dao

import (
	"context"
	"strconv"
	"time"

	"go-common/app/job/main/thumbup/model"
	"go-common/library/log"
)

// PubStatDatabus .
func (d *Dao) PubStatDatabus(c context.Context, business string, mid int64, s *model.Stats, upMid int64) (err error) {
	msg := &model.StatMsg{
		Type:         business,
		ID:           s.ID,
		Count:        s.Likes,
		Timestamp:    time.Now().Unix(),
		OriginID:     s.OriginID,
		DislikeCount: s.Dislikes,
		Mid:          mid,
		UpMid:        upMid,
	}
	if err = d.statDbus.Send(c, strconv.FormatInt(s.ID, 10), msg); err != nil {
		log.Error("d.statDbus.Send error(%v)", err)
		return
	}
	log.Info("pub stat databus success params(%+v)", msg)
	return
}
