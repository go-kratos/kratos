package relation

import (
	"context"

	"go-common/app/interface/main/app-interface/conf"
	relation "go-common/app/service/main/relation/model"
	relationrpc "go-common/app/service/main/relation/rpc/client"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

type Dao struct {
	relationRPC *relationrpc.Service
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		relationRPC: relationrpc.New(c.RelationRPC),
	}
	return
}

// Stat get mid relation stat
func (d *Dao) Stat(c context.Context, mid int64) (stat *relation.Stat, err error) {
	stat, err = d.relationRPC.Stat(c, &relation.ArgMid{Mid: mid})
	if err != nil {
		err = errors.Wrapf(err, "%v", mid)
	}
	return
}

func (d *Dao) FollowersUnread(c context.Context, vmid int64) (res bool, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &relation.ArgMid{Mid: vmid, RealIP: ip}
	if res, err = d.relationRPC.FollowersUnread(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

func (d *Dao) Followings(c context.Context, vmid int64) (res []*relation.Following, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &relation.ArgMid{Mid: vmid, RealIP: ip}
	if res, err = d.relationRPC.Followings(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

func (d *Dao) Relations(c context.Context, mid int64, fids []int64) (res map[int64]*relation.Following, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &relation.ArgRelations{Mid: mid, Fids: fids, RealIP: ip}
	if res, err = d.relationRPC.Relations(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

func (d *Dao) Tag(c context.Context, mid, tid int64) (res []int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &relation.ArgTagId{Mid: mid, TagId: tid, RealIP: ip}
	if res, err = d.relationRPC.Tag(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

func (d *Dao) FollowersUnreadCount(c context.Context, mid int64) (res int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &relation.ArgMid{Mid: mid, RealIP: ip}
	if res, err = d.relationRPC.FollowersUnreadCount(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
