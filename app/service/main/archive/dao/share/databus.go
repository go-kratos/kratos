package share

import (
	"context"
	"strconv"
	"time"

	"go-common/library/log"
)

// PubStatDatabus pub share count into databus.
func (d *Dao) PubStatDatabus(c context.Context, aid int64, share int) (err error) {
	type stat struct {
		Type  string `json:"type"`
		ID    int64  `json:"id"`
		Count int    `json:"count"`
		Ts    int64  `json:"timestamp"`
	}
	if err = d.statDbus.Send(c, strconv.FormatInt(aid, 10), &stat{Type: "archive", ID: aid, Count: share, Ts: time.Now().Unix()}); err != nil {
		log.Error("d.databus.Send error(%v)", err)
	}
	return
}

// PubShare pub first share to databus
func (d *Dao) PubShare(c context.Context, aid int64, mid int64, ip string) (err error) {
	type share struct {
		Event string `json:"event"`
		Mid   int64  `json:"mid"`
		IP    string `json:"ip"`
		Ts    int64  `json:"ts"`
	}
	if err = d.shareDbus.Send(c, strconv.FormatInt(mid, 10), &share{Event: "share", Mid: mid, IP: ip, Ts: time.Now().Unix()}); err != nil {
		log.Error("d.shareDbus.Send error(%v)", err)
	}
	return
}
