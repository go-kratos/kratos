package dao

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/playlist/model"
	pjmdl "go-common/app/job/main/playlist/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
)

var _defaultAdd = int64(1)

// PubView adds a view count.
func (d *Dao) PubView(c context.Context, pid, aid, view int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	view += _defaultAdd
	msg := &pjmdl.StatM{
		Type:      model.PlDBusType,
		ID:        pid,
		Aid:       aid,
		IP:        ip,
		Count:     &view,
		Timestamp: xtime.Time(time.Now().Unix()),
	}
	if err = d.viewDbus.Send(c, strconv.FormatInt(pid, 10), msg); err != nil {
		PromError("databus:发送浏览量", "d.viewDbus.Send(%+v) error(%v)", msg, err)
		return
	}
	log.Info("s.PubView( pid: %d, aid: %d, ip: %s, view: %d)", msg.ID, msg.Aid, msg.IP, *msg.Count)
	return
}

// PubShare adds a share count.
func (d *Dao) PubShare(c context.Context, pid, aid, share int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	share += _defaultAdd
	msg := &pjmdl.StatM{
		Type:      model.PlDBusType,
		ID:        pid,
		Aid:       aid,
		IP:        ip,
		Count:     &share,
		Timestamp: xtime.Time(time.Now().Unix()),
	}
	if err = d.shareDbus.Send(c, strconv.FormatInt(pid, 10), msg); err != nil {
		PromError("databus:发送分享数", "d.shareDbus.Send(%+v) error(%v)", msg, err)
		return
	}
	log.Info("s.PubShare( pid: %d, aid: %d, ip: %s, share: %d)", msg.ID, msg.Aid, msg.IP, *msg.Count)
	return
}
