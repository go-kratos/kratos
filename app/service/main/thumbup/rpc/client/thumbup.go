package client

import (
	"context"

	"go-common/app/service/main/thumbup/model"
	"go-common/library/net/rpc"
)

const (
	_like          = "RPC.Like"
	_likeWithStats = "RPC.LikeWithStats"
	_hasLike       = "RPC.HasLike"
	_stats         = "RPC.Stats"
	_userLikes     = "RPC.UserLikes"
	_userDislikes  = "RPC.UserDislikes"
	_itemLikes     = "RPC.ItemLikes"
	_itemDislikes  = "RPC.ItemDislikes"
	_statsWithLike = "RPC.StatsWithLike"
	_userTotalLike = "RPC.UserTotalLike"
	_updateCount   = "RPC.UpdateCount"
	_rawStats      = "RPC.RawStats"
)

const (
	_appid = "community.service.thumbup"
)

var (
	// _noArg   = &struct{}{}
	_noReply = &struct{}{}
)

//go:generate mockgen -source thumbup.go  -destination mock.go -package client

// ThumbupRPC rpc interface
type ThumbupRPC interface {
	Like(c context.Context, arg *model.ArgLike) (err error)
	HasLike(c context.Context, arg *model.ArgHasLike) (res map[int64]int8, err error)
	Stats(c context.Context, arg *model.ArgStats) (res map[int64]*model.Stats, err error)
	UserLikes(c context.Context, arg *model.ArgUserLikes) (res []*model.ItemLikeRecord, err error)
	UserDislikes(c context.Context, arg *model.ArgUserLikes) (res []*model.ItemLikeRecord, err error)
	ItemLikes(c context.Context, arg *model.ArgItemLikes) (res []*model.UserLikeRecord, err error)
	ItemDislikes(c context.Context, arg *model.ArgItemLikes) (res []*model.UserLikeRecord, err error)
	StatsWithLike(c context.Context, arg *model.ArgStatsWithLike) (res map[int64]*model.StatsWithLike, err error)
	UserTotalLike(c context.Context, arg *model.ArgUserLikes) (res *model.UserTotalLike, err error)
}

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

// Like add like/dislike
func (s *Service) Like(c context.Context, arg *model.ArgLike) (err error) {
	err = s.client.Call(c, _like, arg, _noReply)
	return
}

// LikeWithStats add like/dislike and return the stats info
func (s *Service) LikeWithStats(ctx context.Context, arg *model.ArgLike) (res *model.Stats, err error) {
	err = s.client.Call(ctx, _likeWithStats, arg, &res)
	return
}

// HasLike query user has liked something
func (s *Service) HasLike(c context.Context, arg *model.ArgHasLike) (res map[int64]int8, err error) {
	err = s.client.Call(c, _hasLike, arg, &res)
	return
}

// Stats return stats message
func (s *Service) Stats(c context.Context, arg *model.ArgStats) (res map[int64]*model.Stats, err error) {
	err = s.client.Call(c, _stats, arg, &res)
	return
}

// UserLikes user likes list
func (s *Service) UserLikes(c context.Context, arg *model.ArgUserLikes) (res []*model.ItemLikeRecord, err error) {
	err = s.client.Call(c, _userLikes, arg, &res)
	return
}

// UserDislikes user dislikes list
func (s *Service) UserDislikes(c context.Context, arg *model.ArgUserLikes) (res []*model.ItemLikeRecord, err error) {
	err = s.client.Call(c, _userDislikes, arg, &res)
	return
}

// ItemLikes item likes list
func (s *Service) ItemLikes(c context.Context, arg *model.ArgItemLikes) (res []*model.UserLikeRecord, err error) {
	err = s.client.Call(c, _itemLikes, arg, &res)
	return
}

// ItemDislikes item dislikes list
func (s *Service) ItemDislikes(c context.Context, arg *model.ArgItemLikes) (res []*model.UserLikeRecord, err error) {
	err = s.client.Call(c, _itemDislikes, arg, &res)
	return
}

// StatsWithLike return stats and like state
func (s *Service) StatsWithLike(c context.Context, arg *model.ArgStatsWithLike) (res map[int64]*model.StatsWithLike, err error) {
	err = s.client.Call(c, _statsWithLike, arg, &res)
	return
}

// UserTotalLike user item list with total count
func (s *Service) UserTotalLike(c context.Context, arg *model.ArgUserLikes) (res *model.UserTotalLike, err error) {
	err = s.client.Call(c, _userTotalLike, arg, &res)
	return
}

// UpdateCount update count
func (s *Service) UpdateCount(c context.Context, arg *model.ArgUpdateCount) (err error) {
	err = s.client.Call(c, _updateCount, arg, _noReply)
	return
}

// RawStats get stat changes
func (s *Service) RawStats(c context.Context, arg *model.ArgRawStats) (res model.RawStats, err error) {
	err = s.client.Call(c, _rawStats, arg, &res)
	return
}
