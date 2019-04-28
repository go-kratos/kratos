package dao

import (
	"context"
	"strconv"
	"time"

	"go-common/app/service/main/share/model"
)

// PubShare .
func (d *Dao) PubShare(c context.Context, p *model.ShareParams) (err error) {
	msg := &model.MIDShare{
		OID:  p.OID,
		MID:  p.MID,
		TP:   p.TP,
		Time: time.Now().Unix(),
	}
	return d.databus.Send(c, strconv.FormatInt(p.MID, 10), &msg)
}

// PubStatShare .
func (d *Dao) PubStatShare(c context.Context, typ string, oid, count int64) (err error) {
	msg := &model.ArchiveShare{
		Type:  typ,
		ID:    oid,
		Count: int(count),
		Ts:    time.Now().Unix(),
	}
	return d.archiveDatabus.Send(c, strconv.FormatInt(oid, 10), &msg)
}
