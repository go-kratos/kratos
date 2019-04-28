package bplus

import (
	"context"
	"strconv"

	"go-common/app/interface/main/app-interface/model/space"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// NotifyContribute .
func (d *Dao) NotifyContribute(c context.Context, vmid int64, attrs *space.Attrs, ctime xtime.Time) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	value := struct {
		Vmid  int64        `json:"vmid"`
		Attrs *space.Attrs `json:"attrs"`
		CTime xtime.Time   `json:"ctime"`
		IP    string       `json:"ip"`
	}{vmid, attrs, ctime, ip}
	if err = d.pub.Send(c, strconv.FormatInt(vmid, 10), value); err != nil {
		err = errors.Wrapf(err, "%v", value)
	}
	return
}
