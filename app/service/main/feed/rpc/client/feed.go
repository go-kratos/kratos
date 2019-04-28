package client

import (
	"context"

	artmdl "go-common/app/interface/openplatform/article/model"
	feedmdl "go-common/app/service/main/feed/model"
	"go-common/library/net/rpc"
)

const (
	_appFeed            = "RPC.AppFeed"
	_webFeed            = "RPC.WebFeed"
	_archiveFeed        = "RPC.ArchiveFeed"
	_bangumiFeed        = "RPC.BangumiFeed"
	_addArc             = "RPC.AddArc"
	_delArc             = "RPC.DelArc"
	_purgeFeedCache     = "RPC.PurgeFeedCache"
	_fold               = "RPC.Fold"
	_appUnreadCount     = "RPC.AppUnreadCount"
	_webUnreadCount     = "RPC.WebUnreadCount"
	_changeArcUpper     = "RPC.ChangeArcUpper"
	_articleFeed        = "RPC.ArticleFeed"
	_articleUnreadCount = "RPC.ArticleUnreadCount"
)

const (
	_appid = "community.service.feed"
)

var (
	_noArg = &struct{}{}
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

// AppFeed receive ArgMid contains mid and real ip, then init app feed.
func (s *Service) AppFeed(c context.Context, arg *feedmdl.ArgFeed) (res []*feedmdl.Feed, err error) {
	err = s.client.Call(c, _appFeed, arg, &res)
	return
}

// WebFeed receive ArgMid contains mid and real ip, then init web feed without fold.
func (s *Service) WebFeed(c context.Context, arg *feedmdl.ArgFeed) (res []*feedmdl.Feed, err error) {
	err = s.client.Call(c, _webFeed, arg, &res)
	return
}

// ArchiveFeed receive ArgMid contains mid and real ip
func (s *Service) ArchiveFeed(c context.Context, arg *feedmdl.ArgFeed) (res []*feedmdl.Feed, err error) {
	err = s.client.Call(c, _archiveFeed, arg, &res)
	return
}

// BangumiFeed receive ArgMid contains mid and real ip
func (s *Service) BangumiFeed(c context.Context, arg *feedmdl.ArgFeed) (res []*feedmdl.Feed, err error) {
	err = s.client.Call(c, _bangumiFeed, arg, &res)
	return
}

// ArticleFeed receive ArgMid and return article feed.
func (s *Service) ArticleFeed(c context.Context, arg *feedmdl.ArgFeed) (res []*artmdl.Meta, err error) {
	err = s.client.Call(c, _articleFeed, arg, &res)
	return
}

// ArticleUnreadCount return unread count of article feed.
func (s *Service) ArticleUnreadCount(c context.Context, arg *feedmdl.ArgMid) (res int, err error) {
	err = s.client.Call(c, _articleUnreadCount, arg, &res)
	return
}

// AddArc add archive when archive passed. purge cache.
func (s *Service) AddArc(c context.Context, arg *feedmdl.ArgArc) (err error) {
	err = s.client.Call(c, _addArc, arg, &struct{}{})
	return
}

// DelArc delete archive when archive not passed. purge cache.
func (s *Service) DelArc(c context.Context, arg *feedmdl.ArgAidMid) (err error) {
	err = s.client.Call(c, _delArc, arg, &struct{}{})
	return
}

// PurgeFeedCache purge cache when attention/unattention upper
func (s *Service) PurgeFeedCache(c context.Context, arg *feedmdl.ArgMid) (err error) {
	err = s.client.Call(c, _purgeFeedCache, arg, &struct{}{})
	return
}

// Fold receive ArgFold contains mid, then return upper's fold archives.
func (s *Service) Fold(c context.Context, arg *feedmdl.ArgFold) (res []*feedmdl.Feed, err error) {
	err = s.client.Call(c, _fold, arg, &res)
	return
}

// AppUnreadCount receive ArgUnreadCount contains mid, and withoutBangumi then return unread count.
func (s *Service) AppUnreadCount(c context.Context, arg *feedmdl.ArgUnreadCount) (res int, err error) {
	err = s.client.Call(c, _appUnreadCount, arg, &res)
	return
}

// WebUnreadCount receive ArgUnreadCount contains mid, then return unread count.
func (s *Service) WebUnreadCount(c context.Context, arg *feedmdl.ArgMid) (res int, err error) {
	err = s.client.Call(c, _webUnreadCount, arg, &res)
	return
}

// ChangeArcUpper refresh feed cache when change archive's author
func (s *Service) ChangeArcUpper(c context.Context, arg *feedmdl.ArgChangeUpper) (err error) {
	err = s.client.Call(c, _changeArcUpper, arg, &struct{}{})
	return
}
