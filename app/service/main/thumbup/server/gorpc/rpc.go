package server

import (
	"go-common/app/service/main/thumbup/conf"
	"go-common/app/service/main/thumbup/model"
	"go-common/app/service/main/thumbup/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC .
type RPC struct {
	s *service.Service
}

// New creates rpc server.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
	svr = rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Ping checks connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// Auth check connection success.
func (r *RPC) Auth(c context.Context, arg *rpc.Auth, res *struct{}) (err error) {
	return
}

// Like add like/dislike
func (r *RPC) Like(c context.Context, arg *model.ArgLike, res *struct{}) (err error) {
	_, err = r.s.Like(c, arg.Business, arg.Mid, arg.OriginID, arg.MessageID, arg.Type, arg.UpMid)
	return
}

// LikeWithStats add like/dislike and return the stats info
func (r *RPC) LikeWithStats(ctx context.Context, arg *model.ArgLike, res *model.Stats) (err error) {
	*res, err = r.s.Like(ctx, arg.Business, arg.Mid, arg.OriginID, arg.MessageID, arg.Type, arg.UpMid)
	return
}

// Stats return stats message
func (r *RPC) Stats(c context.Context, arg *model.ArgStats, res *map[int64]*model.Stats) (err error) {
	*res, err = r.s.Stats(c, arg.Business, arg.OriginID, arg.MessageIDs)
	return
}

// UserLikes user likes list
func (r *RPC) UserLikes(c context.Context, arg *model.ArgUserLikes, res *[]*model.ItemLikeRecord) (err error) {
	*res, err = r.s.UserLikes(c, arg.Business, arg.Mid, arg.Pn, arg.Ps)
	return
}

// UserDislikes user dislikes list
func (r *RPC) UserDislikes(c context.Context, arg *model.ArgUserLikes, res *[]*model.ItemLikeRecord) (err error) {
	*res, err = r.s.UserDislikes(c, arg.Business, arg.Mid, arg.Pn, arg.Ps)
	return
}

// ItemLikes item likes list
func (r *RPC) ItemLikes(c context.Context, arg *model.ArgItemLikes, res *[]*model.UserLikeRecord) (err error) {
	*res, err = r.s.ItemLikes(c, arg.Business, arg.OriginID, arg.MessageID, arg.Pn, arg.Ps, arg.Mid)
	return
}

// ItemDislikes item dislikes list
func (r *RPC) ItemDislikes(c context.Context, arg *model.ArgItemLikes, res *[]*model.UserLikeRecord) (err error) {
	*res, err = r.s.ItemDislikes(c, arg.Business, arg.OriginID, arg.MessageID, arg.Pn, arg.Ps, arg.Mid)
	return
}

// HasLike query user has liked something
func (r *RPC) HasLike(c context.Context, arg *model.ArgHasLike, res *map[int64]int8) (err error) {
	*res, _, err = r.s.HasLike(c, arg.Business, arg.Mid, arg.MessageIDs)
	return
}

// StatsWithLike return stats and like state
func (r *RPC) StatsWithLike(c context.Context, arg *model.ArgStatsWithLike, res *map[int64]*model.StatsWithLike) (err error) {
	*res, err = r.s.StatsWithLike(c, arg.Business, arg.Mid, arg.OriginID, arg.MessageIDs)
	return
}

// UserTotalLike user item list with total count
func (r *RPC) UserTotalLike(c context.Context, arg *model.ArgUserLikes, res **model.UserTotalLike) (err error) {
	*res, err = r.s.UserTotalLike(c, arg.Business, arg.Mid, arg.Pn, arg.Ps)
	return
}

// UpdateCount update count
func (r *RPC) UpdateCount(c context.Context, arg *model.ArgUpdateCount, res *struct{}) (err error) {
	err = r.s.UpdateCount(c, arg.Business, arg.OriginID, arg.MessageID, arg.LikeChange, arg.DislikeChange, arg.RealIP, arg.Operator)
	return
}

// RawStats get stat changes
func (r *RPC) RawStats(c context.Context, arg *model.ArgRawStats, res *model.RawStats) (err error) {
	*res, err = r.s.RawStats(c, arg.Business, arg.OriginID, arg.MessageID)
	return
}
