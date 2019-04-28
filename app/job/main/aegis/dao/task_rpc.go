package dao

import (
	"context"

	relmod "go-common/app/service/main/relation/model"
	uprpc "go-common/app/service/main/up/api/v1"
	"go-common/library/log"

	terrors "github.com/pkg/errors"
)

// FansCount 粉丝数
func (d *Dao) FansCount(c context.Context, mid int64) (fans int64, err error) {
	if d.c.Debug {
		return 10086, nil
	}
	arg := &relmod.ArgMid{Mid: mid}
	stat, err := d.relRPC.Stat(c, arg)
	if err != nil || stat == nil {
		log.Error("FansCount error(%v)", terrors.WithStack(err))
		return
	}
	fans = stat.Follower
	return
}

// UpSpecial 分组信息
func (d *Dao) UpSpecial(c context.Context, mid int64) (groupids []int64, err error) {
	if d.c.Debug {
		return
	}
	req := &uprpc.UpSpecialReq{Mid: mid}
	var reply *uprpc.UpSpecialReply
	if reply, err = d.upRPC.UpSpecial(c, req); err != nil || reply == nil {
		log.Error("UpSpecial(%d) error(%+v)", mid, err)
		return
	}
	if reply.UpSpecial != nil {
		groupids = reply.UpSpecial.GroupIDs
	}
	return
}
