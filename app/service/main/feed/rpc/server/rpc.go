package server

import (
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/feed/conf"
	feedmdl "go-common/app/service/main/feed/model"
	"go-common/app/service/main/feed/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC struct info
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

// AddArc add archive when archive passed. purge cache.
func (r *RPC) AddArc(c context.Context, arg *feedmdl.ArgArc, res *struct{}) (err error) {
	err = r.s.AddArc(c, arg.Mid, arg.Aid, arg.PubDate, arg.RealIP)
	return
}

// DelArc delete archive when archive not passed. purge cache.
func (r *RPC) DelArc(c context.Context, arg *feedmdl.ArgAidMid, res *struct{}) (err error) {
	err = r.s.DelArc(c, arg.Mid, arg.Aid, arg.RealIP)
	return
}

// PurgeFeedCache purge cache when attention/unattention upper
func (r *RPC) PurgeFeedCache(c context.Context, arg *feedmdl.ArgMid, res *struct{}) (err error) {
	err = r.s.PurgeFeedCache(c, arg.Mid, arg.RealIP)
	return
}

// AppFeed receive ArgMid contains mid and real ip, then return app feed.
func (r *RPC) AppFeed(c context.Context, arg *feedmdl.ArgFeed, res *[]*feedmdl.Feed) (err error) {
	*res, err = r.s.Feed(c, true, arg.Mid, arg.Pn, arg.Ps, arg.RealIP)
	return
}

// WebFeed receive ArgMid contains mid and real ip, then return app feed.
func (r *RPC) WebFeed(c context.Context, arg *feedmdl.ArgFeed, res *[]*feedmdl.Feed) (err error) {
	*res, err = r.s.Feed(c, false, arg.Mid, arg.Pn, arg.Ps, arg.RealIP)
	return
}

// Fold receive ArgFold contains mid, then return upper's fold archives.
func (r *RPC) Fold(c context.Context, arg *feedmdl.ArgFold, res *[]*feedmdl.Feed) (err error) {
	*res, err = r.s.Fold(c, arg.Mid, arg.Aid, arg.RealIP)
	return
}

// AppUnreadCount receive ArgUnreadCount contains mid, then return unread count.
func (r *RPC) AppUnreadCount(c context.Context, arg *feedmdl.ArgUnreadCount, res *int) (err error) {
	*res, err = r.s.UnreadCount(c, true, arg.WithoutBangumi, arg.Mid, arg.RealIP)
	return
}

// WebUnreadCount receive ArgUnreadCount contains mid, then return unread count.
func (r *RPC) WebUnreadCount(c context.Context, arg *feedmdl.ArgMid, res *int) (err error) {
	withoutBangumi := false
	*res, err = r.s.UnreadCount(c, false, withoutBangumi, arg.Mid, arg.RealIP)
	return
}

// ArchiveFeed receive ArgFeed contains mid and real ip
func (r *RPC) ArchiveFeed(c context.Context, arg *feedmdl.ArgFeed, res *[]*feedmdl.Feed) (err error) {
	*res, err = r.s.ArchiveFeed(c, arg.Mid, arg.Pn, arg.Ps, arg.RealIP)
	return
}

// BangumiFeed receive ArgFeed contains mid and real ip
func (r *RPC) BangumiFeed(c context.Context, arg *feedmdl.ArgFeed, res *[]*feedmdl.Feed) (err error) {
	*res, err = r.s.BangumiFeed(c, arg.Mid, arg.Pn, arg.Ps, arg.RealIP)
	return
}

// ChangeArcUpper refresh feed cache when change archive's author
func (r *RPC) ChangeArcUpper(c context.Context, arg *feedmdl.ArgChangeUpper, res *struct{}) (err error) {
	err = r.s.ChangeAuthor(c, arg.Aid, arg.OldMid, arg.NewMid, arg.RealIP)
	return
}

// ArticleFeed receive ArgFeed contains mid and real ip
func (r *RPC) ArticleFeed(c context.Context, arg *feedmdl.ArgFeed, res *[]*artmdl.Meta) (err error) {
	*res, err = r.s.ArticleFeed(c, arg.Mid, arg.Pn, arg.Ps, arg.RealIP)
	return
}

// ArticleUnreadCount receive ArgUnreadCount contains mid, then return unread count.
func (r *RPC) ArticleUnreadCount(c context.Context, arg *feedmdl.ArgMid, res *int) (err error) {
	*res, err = r.s.ArticleUnreadCount(c, arg.Mid, arg.RealIP)
	return
}
