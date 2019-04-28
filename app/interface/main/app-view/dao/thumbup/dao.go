package thumbup

import (
	"context"
	"fmt"

	"go-common/app/interface/main/app-view/conf"
	thumbup "go-common/app/service/main/thumbup/api"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is tag dao
type Dao struct {
	thumbupGRPC thumbup.ThumbupClient
}

// New initial tag dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{}
	var err error
	d.thumbupGRPC, err = thumbup.NewClient(c.ThumbupClient)
	if err != nil {
		panic(fmt.Sprintf("thumbup NewClient error(%v)", err))
	}
	return
}

// Like is like view.
func (d *Dao) Like(c context.Context, mid, upMid int64, business string, messageID int64, typ int8) (res *thumbup.LikeReply, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &thumbup.LikeReq{Mid: mid, UpMid: upMid, Business: business, MessageID: messageID, Action: thumbup.Action(typ), IP: ip}
	if res, err = d.thumbupGRPC.Like(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// HasLike user has like
func (d *Dao) HasLike(c context.Context, mid int64, business string, messageIDs []int64) (res *thumbup.HasLikeReply, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &thumbup.HasLikeReq{Mid: mid, Business: business, MessageIds: messageIDs, IP: ip}
	if res, err = d.thumbupGRPC.HasLike(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
