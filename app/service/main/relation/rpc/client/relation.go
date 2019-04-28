package relation

import (
	"context"

	"go-common/app/service/main/relation/model"
	"go-common/library/net/rpc"
)

const (
	// relation
	_relation       = "RPC.Relation"
	_relations      = "RPC.Relations"
	_stat           = "RPC.Stat"
	_stats          = "RPC.Stats"
	_followings     = "RPC.Followings"
	_sameFollowings = "RPC.SameFollowings"
	_whispers       = "RPC.Whispers"
	_blacks         = "RPC.Blacks"
	_followers      = "RPC.Followers"
	_addFollowing   = "RPC.AddFollowing"
	_delFollowing   = "RPC.DelFollowing"
	_addWhisper     = "RPC.AddWhisper"
	_delWhisper     = "RPC.DelWhisper"
	_addBlack       = "RPC.AddBlack"
	_delBlack       = "RPC.DelBlack"
	_delFollower    = "RPC.DelFollower"

	// relation tag
	_tag           = "RPC.Tag"
	_tags          = "RPC.Tags"
	_userTag       = "RPC.UserTag"
	_createTag     = "RPC.CreateTag"
	_updateTag     = "RPC.UpdateTag"
	_delTag        = "RPC.DelTag"
	_tagsAddUsers  = "RPC.TagsAddUsers"
	_tagsCopyUsers = "RPC.TagsCopyUsers"
	_tagsMoveUsers = "RPC.TagsMoveUsers"
	_AddSpecial    = "RPC.AddSpecial"
	_DelSpecial    = "RPC.DelSpecial"
	_Special       = "RPC.Special"

	// prompt
	_prompt      = "RPC.Prompt"
	_closePrompt = "RPC.ClosePrompt"

	// followers incr notify
	_FollowersUnread           = "RPC.FollowersUnread"
	_FollowersUnreadCount      = "RPC.FollowersUnreadCount"
	_ResetFollowersUnread      = "RPC.ResetFollowersUnread"
	_ResetFollowersUnreadCount = "RPC.ResetFollowersUnreadCount"
	_DisableFollowerNotify     = "RPC.DisableFollowerNotify"
	_EnableFollowerNotify      = "RPC.EnableFollowerNotify"
	_FollowerNotifySetting     = "RPC.FollowerNotifySetting"

	// achieve
	_AchieveGet = "RPC.AchieveGet"
	_Achieve    = "RPC.Achieve"
)

var (
	_noRes = &struct{}{}
)

const (
	_appid = "account.service.relation"
)

// Service struct info.
type Service struct {
	client *rpc.Client2
}

// New new service instance and return.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Relation get relation info.
func (s *Service) Relation(c context.Context, arg *model.ArgRelation) (res *model.Following, err error) {
	res = new(model.Following)
	err = s.client.Call(c, _relation, arg, res)
	return
}

// Relations get relation infos.
func (s *Service) Relations(c context.Context, arg *model.ArgRelations) (res map[int64]*model.Following, err error) {
	err = s.client.Call(c, _relations, arg, &res)
	return
}

// Followings get followings infos.
func (s *Service) Followings(c context.Context, arg *model.ArgMid) (res []*model.Following, err error) {
	err = s.client.Call(c, _followings, arg, &res)
	return
}

// Whispers get whispers infos.
func (s *Service) Whispers(c context.Context, arg *model.ArgMid) (res []*model.Following, err error) {
	err = s.client.Call(c, _whispers, arg, &res)
	return
}

// Blacks get black list.
func (s *Service) Blacks(c context.Context, arg *model.ArgMid) (res []*model.Following, err error) {
	err = s.client.Call(c, _blacks, arg, &res)
	return
}

// Followers get followers list.
func (s *Service) Followers(c context.Context, arg *model.ArgMid) (res []*model.Following, err error) {
	err = s.client.Call(c, _followers, arg, &res)
	return
}

// Stat get user relation stat.
func (s *Service) Stat(c context.Context, arg *model.ArgMid) (res *model.Stat, err error) {
	res = new(model.Stat)
	err = s.client.Call(c, _stat, arg, &res)
	return
}

// Stats get users relation stat.
func (s *Service) Stats(c context.Context, arg *model.ArgMids) (res map[int64]*model.Stat, err error) {
	err = s.client.Call(c, _stats, arg, &res)
	return
}

// ModifyRelation modify user relation.
func (s *Service) ModifyRelation(c context.Context, arg *model.ArgFollowing) (err error) {
	switch arg.Action {
	case model.ActAddBlack:
		err = s.client.Call(c, _addBlack, arg, _noRes)
	case model.ActAddFollowing:
		err = s.client.Call(c, _addFollowing, arg, _noRes)
	case model.ActAddWhisper:
		err = s.client.Call(c, _addWhisper, arg, _noRes)
	case model.ActDelBalck:
		err = s.client.Call(c, _delBlack, arg, _noRes)
	case model.ActDelFollower:
		err = s.client.Call(c, _delFollower, arg, _noRes)
	case model.ActDelFollowing:
		err = s.client.Call(c, _delFollowing, arg, _noRes)
	case model.ActDelWhisper:
		err = s.client.Call(c, _delWhisper, arg, _noRes)
	default:

	}
	return
}

// Tag tag
func (s *Service) Tag(c context.Context, arg *model.ArgTagId) (res []int64, err error) {
	err = s.client.Call(c, _tag, arg, &res)
	return
}

// Tags tags
func (s *Service) Tags(c context.Context, arg *model.ArgMid) (res []*model.TagCount, err error) {
	err = s.client.Call(c, _tags, arg, &res)
	return
}

// UserTag user tag
func (s *Service) UserTag(c context.Context, arg *model.ArgRelation) (res map[int64]string, err error) {
	err = s.client.Call(c, _userTag, arg, &res)
	return
}

// CreateTag create tag
func (s *Service) CreateTag(c context.Context, arg *model.ArgTag) (res int64, err error) {
	err = s.client.Call(c, _createTag, arg, &res)
	return
}

// UpdateTag update tag
func (s *Service) UpdateTag(c context.Context, arg *model.ArgTagUpdate) (err error) {
	err = s.client.Call(c, _updateTag, arg, _noRes)
	return
}

// DelTag del tag
func (s *Service) DelTag(c context.Context, arg *model.ArgTagDel) (err error) {
	err = s.client.Call(c, _delTag, arg, _noRes)
	return
}

// TagsAddUsers tags add users
func (s *Service) TagsAddUsers(c context.Context, arg *model.ArgTagsMoveUsers) (err error) {
	err = s.client.Call(c, _tagsAddUsers, arg, _noRes)
	return
}

// TagsCopyUsers tags copy users
func (s *Service) TagsCopyUsers(c context.Context, arg *model.ArgTagsMoveUsers) (err error) {
	err = s.client.Call(c, _tagsCopyUsers, arg, _noRes)
	return
}

// TagsMoveUsers tags move users
func (s *Service) TagsMoveUsers(c context.Context, arg *model.ArgTagsMoveUsers) (err error) {
	err = s.client.Call(c, _tagsMoveUsers, arg, _noRes)
	return
}

// Prompt rpc rpompt client
func (s *Service) Prompt(c context.Context, arg *model.ArgPrompt) (b bool, err error) {
	err = s.client.Call(c, _prompt, arg, &b)
	return
}

// ClosePrompt close prompt client.
func (s *Service) ClosePrompt(c context.Context, arg *model.ArgPrompt) (err error) {
	err = s.client.Call(c, _closePrompt, arg, _noRes)
	return
}

// AddSpecial add specail.
func (s *Service) AddSpecial(c context.Context, arg *model.ArgFollowing) (err error) {
	err = s.client.Call(c, _AddSpecial, arg, &_noRes)
	return
}

// DelSpecial del special.
func (s *Service) DelSpecial(c context.Context, arg *model.ArgFollowing) (err error) {
	err = s.client.Call(c, _DelSpecial, arg, &_noRes)
	return
}

// Special get special.
func (s *Service) Special(c context.Context, arg *model.ArgMid) (res []int64, err error) {
	err = s.client.Call(c, _Special, arg, &res)
	return
}

// FollowersUnread check unread status, for the 'show red point' function.
func (s *Service) FollowersUnread(c context.Context, arg *model.ArgMid) (show bool, err error) {
	err = s.client.Call(c, _FollowersUnread, arg, &show)
	return
}

// FollowersUnreadCount unread count.
func (s *Service) FollowersUnreadCount(c context.Context, arg *model.ArgMid) (count int64, err error) {
	err = s.client.Call(c, _FollowersUnreadCount, arg, &count)
	return
}

// AchieveGet is
func (s *Service) AchieveGet(c context.Context, arg *model.ArgAchieveGet) (*model.AchieveGetReply, error) {
	reply := &model.AchieveGetReply{}
	err := s.client.Call(c, _AchieveGet, arg, &reply)
	return reply, err
}

// Achieve is
func (s *Service) Achieve(c context.Context, arg *model.ArgAchieve) (*model.Achieve, error) {
	reply := &model.Achieve{}
	err := s.client.Call(c, _Achieve, arg, &reply)
	return reply, err
}

// ResetFollowersUnread is
func (s *Service) ResetFollowersUnread(c context.Context, arg *model.ArgMid) (err error) {
	err = s.client.Call(c, _ResetFollowersUnread, arg, &_noRes)
	return
}

// ResetFollowersUnreadCount is
func (s *Service) ResetFollowersUnreadCount(c context.Context, arg *model.ArgMid) (err error) {
	err = s.client.Call(c, _ResetFollowersUnreadCount, arg, &_noRes)
	return
}

// DisableFollowerNotify set followerNotify as disabled.
func (s *Service) DisableFollowerNotify(c context.Context, arg *model.ArgMid) (err error) {
	err = s.client.Call(c, _DisableFollowerNotify, arg, &_noRes)
	return
}

// EnableFollowerNotify set followerNotify as disabled.
func (s *Service) EnableFollowerNotify(c context.Context, arg *model.ArgMid) (err error) {
	err = s.client.Call(c, _EnableFollowerNotify, arg, &_noRes)
	return
}

// FollowerNotifySetting get followerNotify setting.
func (s *Service) FollowerNotifySetting(c context.Context, arg *model.ArgMid) (followerNotify *model.FollowerNotifySetting, err error) {
	followerNotify = &model.FollowerNotifySetting{}
	err = s.client.Call(c, _FollowerNotifySetting, arg, followerNotify)
	return
}

// SameFollowings is
func (s *Service) SameFollowings(c context.Context, arg *model.ArgSameFollowing) (res []*model.Following, err error) {
	err = s.client.Call(c, _sameFollowings, arg, &res)
	return
}
