package rpc

import (
	"go-common/app/service/main/relation/conf"
	"go-common/app/service/main/relation/model"
	"go-common/app/service/main/relation/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC rpc
type RPC struct {
	s *service.Service
}

// New new rpc server.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
	svr = rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Ping check connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// Relation relation
func (r *RPC) Relation(c context.Context, a *model.ArgRelation, res *model.Following) (err error) {
	var f *model.Following
	if f, err = r.s.Relation(c, a.Mid, a.Fid); err == nil && f != nil {
		*res = *f
	}
	return
}

// Relations relations
func (r *RPC) Relations(c context.Context, a *model.ArgRelations, res *map[int64]*model.Following) (err error) {
	*res, err = r.s.Relations(c, a.Mid, a.Fids)
	return
}

// Stat stat
func (r *RPC) Stat(c context.Context, a *model.ArgMid, res *model.Stat) (err error) {
	var st *model.Stat
	if st, err = r.s.Stat(c, a.Mid); err == nil && st != nil {
		*res = *st
	}
	return
}

// Stats stats
func (r *RPC) Stats(c context.Context, a *model.ArgMids, res *map[int64]*model.Stat) (err error) {
	*res, err = r.s.Stats(c, a.Mids)
	return
}

// Attentions attentions
func (r *RPC) Attentions(c context.Context, a *model.ArgMid, res *[]*model.Following) (err error) {
	*res, err = r.s.Attentions(c, a.Mid)
	return
}

// Followings followings
func (r *RPC) Followings(c context.Context, a *model.ArgMid, res *[]*model.Following) (err error) {
	*res, err = r.s.Followings(c, a.Mid)
	return
}

// AddFollowing add following
func (r *RPC) AddFollowing(c context.Context, a *model.ArgFollowing, res *struct{}) (err error) {
	err = r.s.AddFollowing(c, a.Mid, a.Fid, a.Source, a.Infoc)
	return
}

// DelFollowing del following
func (r *RPC) DelFollowing(c context.Context, a *model.ArgFollowing, res *struct{}) (err error) {
	err = r.s.DelFollowing(c, a.Mid, a.Fid, a.Source, a.Infoc)
	return
}

// Whispers whispers
func (r *RPC) Whispers(c context.Context, a *model.ArgMid, res *[]*model.Following) (err error) {
	*res, err = r.s.Whispers(c, a.Mid)
	return
}

// AddWhisper add whisper
func (r *RPC) AddWhisper(c context.Context, a *model.ArgFollowing, res *struct{}) (err error) {
	err = r.s.AddWhisper(c, a.Mid, a.Fid, a.Source, a.Infoc)
	return
}

// DelWhisper del whisper
func (r *RPC) DelWhisper(c context.Context, a *model.ArgFollowing, res *struct{}) (err error) {
	err = r.s.DelWhisper(c, a.Mid, a.Fid, a.Source, a.Infoc)
	return
}

// Blacks blacks
func (r *RPC) Blacks(c context.Context, a *model.ArgMid, res *[]*model.Following) (err error) {
	*res, err = r.s.Blacks(c, a.Mid)
	return
}

// AddBlack add black
func (r *RPC) AddBlack(c context.Context, a *model.ArgFollowing, res *struct{}) (err error) {
	err = r.s.AddBlack(c, a.Mid, a.Fid, a.Source, a.Infoc)
	return
}

// DelBlack del black
func (r *RPC) DelBlack(c context.Context, a *model.ArgFollowing, res *struct{}) (err error) {
	err = r.s.DelBlack(c, a.Mid, a.Fid, a.Source, a.Infoc)
	return
}

// Followers followers
func (r *RPC) Followers(c context.Context, a *model.ArgMid, res *[]*model.Following) (err error) {
	*res, err = r.s.Followers(c, a.Mid)
	return
}

// DelFollower del Follower
func (r *RPC) DelFollower(c context.Context, a *model.ArgFollowing, res *struct{}) (err error) {
	err = r.s.DelFollowing(c, a.Fid, a.Mid, a.Source, a.Infoc)
	return
}

// Tag tag
func (r *RPC) Tag(c context.Context, a *model.ArgTagId, res *[]int64) (err error) {
	*res, err = r.s.Tag(c, a.Mid, a.TagId, a.RealIP)
	return
}

// Tags tags
func (r *RPC) Tags(c context.Context, a *model.ArgMid, res *[]*model.TagCount) (err error) {
	*res, err = r.s.Tags(c, a.Mid, a.RealIP)
	return
}

// UserTag user tag
func (r *RPC) UserTag(c context.Context, a *model.ArgRelation, res *map[int64]string) (err error) {
	*res, err = r.s.UserTag(c, a.Mid, a.Fid, a.RealIP)
	return
}

// CreateTag create tag
func (r *RPC) CreateTag(c context.Context, a *model.ArgTag, res *int64) (err error) {
	*res, err = r.s.CreateTag(c, a.Mid, a.Tag, a.RealIP)
	return
}

// UpdateTag update tag
func (r *RPC) UpdateTag(c context.Context, a *model.ArgTagUpdate, res *struct{}) (err error) {
	err = r.s.UpdateTag(c, a.Mid, a.TagId, a.New, a.RealIP)
	return
}

// DelTag del tag
func (r *RPC) DelTag(c context.Context, a *model.ArgTagDel, res *struct{}) (err error) {
	err = r.s.DelTag(c, a.Mid, a.TagId, a.RealIP)
	return
}

// TagsAddUsers tags add users
func (r *RPC) TagsAddUsers(c context.Context, a *model.ArgTagsMoveUsers, res *struct{}) (err error) {
	err = r.s.TagsAddUsers(c, a.Mid, a.AfterTagIds, a.Fids, a.RealIP)
	return
}

// TagsCopyUsers tags copy users
func (r *RPC) TagsCopyUsers(c context.Context, a *model.ArgTagsMoveUsers, res *struct{}) (err error) {
	err = r.s.TagsMoveUsers(c, a.Mid, a.BeforeID, a.AfterTagIds, a.Fids, a.RealIP)
	return
}

// TagsMoveUsers tags move users
func (r *RPC) TagsMoveUsers(c context.Context, a *model.ArgTagsMoveUsers, res *struct{}) (err error) {
	err = r.s.TagsMoveUsers(c, a.Mid, a.BeforeID, a.AfterTagIds, a.Fids, a.RealIP)
	return
}

// Prompt rpc prompt.
func (r *RPC) Prompt(c context.Context, m *model.ArgPrompt, res *bool) (err error) {
	*res, err = r.s.Prompt(c, m)
	return
}

// ClosePrompt close prompt.
func (r *RPC) ClosePrompt(c context.Context, m *model.ArgPrompt, res *struct{}) (err error) {
	return r.s.ClosePrompt(c, m)
}

// AddSpecial add user to special.
func (r *RPC) AddSpecial(c context.Context, m *model.ArgFollowing, res *struct{}) (err error) {
	return r.s.AddSpecial(c, m.Mid, m.Fid)
}

// DelSpecial del user from sepcial.
func (r *RPC) DelSpecial(c context.Context, m *model.ArgFollowing, res *struct{}) (err error) {
	return r.s.DelSpecial(c, m.Mid, m.Fid)
}

// Special get user specail list.
func (r *RPC) Special(c context.Context, m *model.ArgMid, res *[]int64) (err error) {
	*res, err = r.s.Special(c, m.Mid)
	return
}

// FollowersUnread is
func (r *RPC) FollowersUnread(c context.Context, arg *model.ArgMid, res *bool) (err error) {
	*res, err = r.s.Unread(c, arg.Mid)
	return
}

// FollowersUnreadCount is
func (r *RPC) FollowersUnreadCount(c context.Context, arg *model.ArgMid, res *int64) (err error) {
	*res, err = r.s.UnreadCount(c, arg.Mid)
	return
}

// AchieveGet is
func (r *RPC) AchieveGet(c context.Context, arg *model.ArgAchieveGet, res *model.AchieveGetReply) error {
	reply, err := r.s.AchieveGet(c, arg)
	if err != nil {
		return err
	}
	*res = *reply
	return nil
}

// Achieve is
func (r *RPC) Achieve(c context.Context, arg *model.ArgAchieve, res *model.Achieve) error {
	reply, err := r.s.Achieve(c, arg)
	if err != nil {
		return err
	}
	*res = *reply
	return nil
}

// ResetFollowersUnread is
func (r *RPC) ResetFollowersUnread(c context.Context, arg *model.ArgMid, res *struct{}) (err error) {
	err = r.s.ResetUnread(c, arg.Mid)
	return
}

// ResetFollowersUnreadCount is
func (r *RPC) ResetFollowersUnreadCount(c context.Context, arg *model.ArgMid, res *struct{}) (err error) {
	err = r.s.ResetUnreadCount(c, arg.Mid)
	return
}

// DisableFollowerNotify set followerNotify as disabled.
func (r *RPC) DisableFollowerNotify(c context.Context, arg *model.ArgMid, res *struct{}) (err error) {
	err = r.s.DisableFollowerNotify(c, arg)
	return
}

// EnableFollowerNotify set followerNotify as enabled.
func (r *RPC) EnableFollowerNotify(c context.Context, arg *model.ArgMid, res *struct{}) (err error) {
	err = r.s.EnableFollowerNotify(c, arg)
	return
}

// FollowerNotifySetting get member follower notify setting
func (r *RPC) FollowerNotifySetting(c context.Context, arg *model.ArgMid, res *model.FollowerNotifySetting) (err error) {
	rely, err := r.s.FollowerNotifySetting(c, arg)
	if err != nil {
		return
	}
	*res = *rely
	return
}

// SameFollowings is
func (r *RPC) SameFollowings(c context.Context, arg *model.ArgSameFollowing, res *[]*model.Following) error {
	reply, err := r.s.SameFollowings(c, arg)
	if err != nil {
		return err
	}
	*res = reply
	return nil
}
