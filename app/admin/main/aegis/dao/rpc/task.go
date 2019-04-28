package rpc

import (
	"context"
	"errors"

	"go-common/app/admin/main/aegis/model"
	acc "go-common/app/service/main/account/api"
	relmod "go-common/app/service/main/relation/model"
	uprpc "go-common/app/service/main/up/api/v1"
	"go-common/library/log"
	"go-common/library/net/metadata"

	terrors "github.com/pkg/errors"
)

//ERROR
var (
	ErrEmptyReply = errors.New("rpc empty reply")
)

// FansCount 粉丝数
func (d *Dao) FansCount(c context.Context, mids []int64) (fans map[int64]int64, err error) {
	arg := &relmod.ArgMids{Mids: mids}
	stats, err := d.relRPC.Stats(c, arg)
	if err != nil {
		log.Error("FansCount error(%v)", terrors.WithStack(err))
		return
	}
	fans = make(map[int64]int64)
	for mid, item := range stats {
		fans[mid] = item.Follower
	}
	log.Info("FansCount fans(%+v)", fans)
	return
}

// UserInfos 提供给资源列表批量查
func (d *Dao) UserInfos(c context.Context, mids []int64) (res map[int64]*model.UserInfo, err error) {
	arg1 := &relmod.ArgMids{Mids: mids}
	stats, err := d.relRPC.Stats(c, arg1)
	if err != nil {
		log.Error("Stats error(%v)", terrors.WithStack(err))
		return
	}
	midsReq := &acc.MidsReq{
		Mids:   mids,
		RealIp: metadata.String(c, metadata.RemoteIP),
	}
	cardsreply, err := d.AccountClient.Cards3(c, midsReq)
	if err != nil {
		log.Error("Cards3(%+v) error(%v)", mids, terrors.WithStack(err))
		return
	}
	if cardsreply == nil {
		err = ErrEmptyReply
		log.Error("Cards3(%+v) error(%v)", mids, terrors.WithStack(err))
		return
	}
	cards := cardsreply.Cards
	res = make(map[int64]*model.UserInfo)
	for _, mid := range mids {
		userinfo := &model.UserInfo{Mid: mid}
		if card, ok := cards[mid]; ok {
			userinfo.Name = card.Name
			userinfo.Official = card.Official
		}
		if stat, ok := stats[mid]; ok {
			userinfo.Follower = stat.Follower
		}
		res[mid] = userinfo
	}
	return
}

// Profile get account.
func (d *Dao) Profile(c context.Context, mid int64) (userinfo *model.UserInfo, err error) {
	if mid <= 0 {
		return
	}
	midReq := &acc.MidReq{
		Mid:    mid,
		RealIp: metadata.String(c, metadata.RemoteIP),
	}

	res, err := d.AccountClient.ProfileWithStat3(c, midReq)
	if err != nil {
		log.Error("d.acc.ProfileWithStat3() error(%v) arg(%+v)", err, midReq)
	}
	if res == nil {
		err = ErrEmptyReply
		log.Error("ProfileWithStat3(%+v) error(%v)", mid, terrors.WithStack(err))
		return
	}
	userinfo = &model.UserInfo{
		Mid:      mid,
		Follower: res.Follower,
	}
	if res.Profile != nil {
		userinfo.Official = res.Profile.Official
		userinfo.Name = res.Profile.Name
	}

	return
}

// Info3 get Name.
func (d *Dao) Info3(c context.Context, mid int64) (res *acc.Info, err error) {
	midReq := &acc.MidReq{
		Mid:    mid,
		RealIp: metadata.String(c, metadata.RemoteIP),
	}
	info, err := d.AccountClient.Info3(c, midReq)
	if err != nil {
		log.Error("query info3 failed,mid(%v), err(%v)", mid, err)
		return
	}
	if info == nil {
		err = ErrEmptyReply
		log.Error("Info3(%+v) error(%v)", mid, terrors.WithStack(err))
		return
	}
	res = info.Info
	log.Info("getUserInfo userbase (%v)", res)
	return
}

// UpSpecial 分组信息
func (d *Dao) UpSpecial(c context.Context, mid int64) (ups *uprpc.UpSpecial, err error) {
	midReq := &uprpc.UpSpecialReq{
		Mid: mid,
	}
	var reply *uprpc.UpSpecialReply
	if reply, err = d.UpClient.UpSpecial(c, midReq); err != nil {
		log.Error("UpSpecial(%d) error(%v)", mid, terrors.WithStack(err))
		return
	}
	if reply == nil {
		err = ErrEmptyReply
		log.Error("UpSpecial(%+v) error(%v)", mid, terrors.WithStack(err))
		return
	}
	ups = reply.UpSpecial
	return
}

//UpsSpecial 分组信息
func (d *Dao) UpsSpecial(c context.Context, mids []int64) (ups map[int64]*uprpc.UpSpecial, err error) {
	midReq := &uprpc.UpsSpecialReq{
		Mids: mids,
	}
	var reply *uprpc.UpsSpecialReply
	if reply, err = d.UpClient.UpsSpecial(c, midReq); err != nil {
		log.Error("UpsSpecial(%d) error(%v)", mids, terrors.WithStack(err))
		return
	}
	if reply == nil {
		err = ErrEmptyReply
		log.Error("UpsSpecial(%+v) error(%v)", mids, terrors.WithStack(err))
		return
	}
	ups = reply.UpSpecials
	return
}

//UpGroups 所有分组
func (d *Dao) UpGroups(c context.Context) (upgs map[int64]*uprpc.UpGroup, err error) {
	noReq := &uprpc.NoArgReq{}
	var reply *uprpc.UpGroupsReply
	if reply, err = d.UpClient.UpGroups(c, noReq); err != nil {
		log.Error("UpGroups error(%v)", terrors.WithStack(err))
		return
	}
	if reply == nil {
		err = ErrEmptyReply
		log.Error("UpGroups error(%v)", terrors.WithStack(err))
		return
	}
	upgs = reply.UpGroups
	return
}
