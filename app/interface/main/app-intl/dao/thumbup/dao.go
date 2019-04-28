package thumbup

import (
	"context"

	"go-common/app/interface/main/app-intl/conf"
	thumbup "go-common/app/service/main/thumbup/model"
	thumbuprpc "go-common/app/service/main/thumbup/rpc/client"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is tag dao
type Dao struct {
	thumbupRPC *thumbuprpc.Service
}

// New initial tag dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		thumbupRPC: thumbuprpc.New(c.ThumbupRPC),
	}
	return
}

// Like is like view.
func (d *Dao) Like(c context.Context, mid, upMid int64, business string, messageID int64, typ int8) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &thumbup.ArgLike{Mid: mid, UpMid: upMid, Business: business, MessageID: messageID, Type: typ, RealIP: ip}
	return d.thumbupRPC.Like(c, arg)
}

// LikeWithStat is like with stat.
func (d *Dao) LikeWithStat(c context.Context, mid, upMid int64, business string, messageID int64, typ int8) (stat *thumbup.Stats, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &thumbup.ArgLike{Mid: mid, UpMid: upMid, Business: business, MessageID: messageID, Type: typ, RealIP: ip}
	return d.thumbupRPC.LikeWithStats(c, arg)
}

// HasLike user has like
func (d *Dao) HasLike(c context.Context, mid int64, business string, messageIDs []int64) (res map[int64]int8, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &thumbup.ArgHasLike{Mid: mid, MessageIDs: messageIDs, Business: business, RealIP: ip}
	if res, err = d.thumbupRPC.HasLike(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// Stat is
func (d *Dao) Stat(c context.Context, mid int64, business string, messageIDs []int64) (res map[int64]*thumbup.Stats, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &thumbup.ArgStats{Business: business, MessageIDs: messageIDs, RealIP: ip}
	if res, err = d.thumbupRPC.Stats(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
