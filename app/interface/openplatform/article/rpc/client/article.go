package client

import (
	"context"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/net/rpc"
)

const (
	_addArticle              = "RPC.AddArticle"
	_updateArticle           = "RPC.UpdateArticle"
	_delArticle              = "RPC.DelArticle"
	_creationUpperArticles   = "RPC.CreationUpperArticles"
	_creationArticle         = "RPC.CreationArticle"
	_categories              = "RPC.Categories"
	_categoriesMap           = "RPC.CategoriesMap"
	_setStat                 = "RPC.SetStat"
	_addArticleCache         = "RPC.AddArticleCache"
	_updateArticleCache      = "RPC.UpdateArticleCache"
	_delArticleCache         = "RPC.DelArticleCache"
	_upsArtMetas             = "RPC.UpsArtMetas"
	_artMetas                = "RPC.ArticleMetas"
	_updateRecommends        = "RPC.UpdateRecommends"
	_recommends              = "RPC.Recommends"
	_creationWithdrawArticle = "RPC.CreationWithdrawArticle"
	_upArtMetas              = "RPC.UpArtMetas"
	_addArtDraft             = "RPC.AddArtDraft"
	_delArtDraft             = "RPC.DelArtDraft"
	_artDraft                = "RPC.ArtDraft"
	_upperDrafts             = "RPC.UpperDrafts"
	_articleRemainCount      = "RPC.ArticleRemainCount"
	_delRecommendArtCache    = "RPC.DelRecommendArtCache"
	_favorites               = "RPC.Favorites"
	_updateAuthorCache       = "RPC.UpdateAuthorCache"
	_updateSortCache         = "RPC.UpdateSortCache"
	_isAuthor                = "RPC.IsAuthor"
	_newArticleCount         = "RPC.NewArticleCount"
	_hadLikesByMid           = "RPC.HadLikesByMid"
	_upMoreArts              = "RPC.UpMoreArts"
	_creationUpStat          = "RPC.CreationUpStat"
	_creationUpThirtyDayStat = "RPC.CreationUpThirtyDayStat"
	_upLists                 = "RPC.UpLists"
	_rebuildAllRC            = "RPC.RebuildAllListReadCount"
	_updateHotspots          = "RPC.UpdateHotspots"
)

const (
	_appid = "article.service"
)

var (
	_noArg   = &struct{}{}
	_noReply = &struct{}{}
)

// Service struct info.
type Service struct {
	client *rpc.Client2
}

//go:generate mockgen -source article.go  -destination mock.go -package client

// ArticleRPC article rpc.
type ArticleRPC interface {
	AddArticle(c context.Context, arg *model.ArgArticle) (id int64, err error)
	AddArticleCache(c context.Context, arg *model.ArgAid) (err error)
	UpdateArticleCache(c context.Context, arg *model.ArgAidCid) (err error)
	DelArticleCache(c context.Context, arg *model.ArgAidMid) (err error)
	UpdateArticle(c context.Context, arg *model.ArgArticle) (err error)
	CreationWithdrawArticle(c context.Context, arg *model.ArgAidMid) (err error)
	DelArticle(c context.Context, arg *model.ArgAidMid) (err error)
	CreationArticle(c context.Context, arg *model.ArgAidMid) (res *model.Article, err error)
	CreationUpperArticles(c context.Context, arg *model.ArgCreationArts) (res *model.CreationArts, err error)
	Categories(c context.Context, arg *model.ArgIP) (res *model.Categories, err error)
	CategoriesMap(c context.Context, arg *model.ArgIP) (res map[int64]*model.Category, err error)
	SetStat(c context.Context, arg *model.ArgStats) (err error)
	UpsArtMetas(c context.Context, arg *model.ArgUpsArts) (res map[int64][]*model.Meta, err error)
	ArticleMetas(c context.Context, arg *model.ArgAids) (res map[int64]*model.Meta, err error)
	UpdateRecommends(c context.Context) (err error)
	Recommends(c context.Context, arg *model.ArgRecommends) (res []*model.Meta, err error)
	UpArtMetas(c context.Context, arg *model.ArgUpArts) (res *model.UpArtMetas, err error)
	AddArtDraft(c context.Context, arg *model.ArgArticle) (id int64, err error)
	UpdateArtDraft(c context.Context, arg *model.ArgAidMid) (err error)
	DelArtDraft(c context.Context, arg *model.ArgAidMid) (err error)
	ArtDraft(c context.Context, arg *model.ArgAidMid) (res *model.Draft, err error)
	UpperDrafts(c context.Context, arg *model.ArgUpDraft) (res *model.Drafts, err error)
	ArticleRemainCount(c context.Context, arg *model.ArgMid) (res int, err error)
	DelRecommendArtCache(c context.Context, arg *model.ArgAidCid) (err error)
	Favorites(c context.Context, arg *model.ArgFav) (res []*model.Favorite, err error)
	UpdateAuthorCache(c context.Context, arg *model.ArgAuthor) (err error)
	UpdateSortCache(c context.Context, arg *model.ArgSort) (err error)
	IsAuthor(c context.Context, arg *model.ArgMid) (res bool, err error)
	NewArticleCount(c context.Context, arg *model.ArgNewArt) (res int64, err error)
	HadLikesByMid(c context.Context, arg *model.ArgMidAids) (res map[int64]int, err error)
	UpMoreArts(c context.Context, arg *model.ArgAid) (res []*model.Meta, err error)
	CreationUpStat(c context.Context, arg *model.ArgMid) (res model.UpStat, err error)
	CreationUpThirtyDayStat(c context.Context, arg *model.ArgMid) (res []*model.ThirtyDayArticle, err error)
	UpLists(c context.Context, arg *model.ArgMid) (res model.UpLists, err error)
}

// New new service instance and return.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// AddArticle adds article when article passed. purge cache.
func (s *Service) AddArticle(c context.Context, arg *model.ArgArticle) (id int64, err error) {
	err = s.client.Call(c, _addArticle, arg, &id)
	return
}

// AddArticleCache adds article when article passed. purge cache.
func (s *Service) AddArticleCache(c context.Context, arg *model.ArgAid) (err error) {
	err = s.client.Call(c, _addArticleCache, arg, _noReply)
	return
}

// UpdateArticleCache adds article when article passed. purge cache.
func (s *Service) UpdateArticleCache(c context.Context, arg *model.ArgAidCid) (err error) {
	err = s.client.Call(c, _updateArticleCache, arg, _noReply)
	return
}

// DelArticleCache adds article when article passed. purge cache.
func (s *Service) DelArticleCache(c context.Context, arg *model.ArgAidMid) (err error) {
	err = s.client.Call(c, _delArticleCache, arg, _noReply)
	return
}

// UpdateArticle updates article when article passed. purge cache.
func (s *Service) UpdateArticle(c context.Context, arg *model.ArgArticle) (err error) {
	err = s.client.Call(c, _updateArticle, arg, _noReply)
	return
}

// CreationWithdrawArticle author withdraw article.
func (s *Service) CreationWithdrawArticle(c context.Context, arg *model.ArgAidMid) (err error) {
	err = s.client.Call(c, _creationWithdrawArticle, arg, _noReply)
	return
}

// DelArticle drops article when article not passed. purge cache.
func (s *Service) DelArticle(c context.Context, arg *model.ArgAidMid) (err error) {
	err = s.client.Call(c, _delArticle, arg, _noReply)
	return
}

// CreationArticle gets article's meta.
func (s *Service) CreationArticle(c context.Context, arg *model.ArgAidMid) (res *model.Article, err error) {
	err = s.client.Call(c, _creationArticle, arg, &res)
	return
}

// CreationUpperArticles gets article's meta.
func (s *Service) CreationUpperArticles(c context.Context, arg *model.ArgCreationArts) (res *model.CreationArts, err error) {
	err = s.client.Call(c, _creationUpperArticles, arg, &res)
	return
}

// Categories list categories of article
func (s *Service) Categories(c context.Context, arg *model.ArgIP) (res *model.Categories, err error) {
	err = s.client.Call(c, _categories, arg, &res)
	return
}

// CategoriesMap list categories of article map
func (s *Service) CategoriesMap(c context.Context, arg *model.ArgIP) (res map[int64]*model.Category, err error) {
	err = s.client.Call(c, _categoriesMap, arg, &res)
	return
}

// SetStat updates article's stat cache.
func (s *Service) SetStat(c context.Context, arg *model.ArgStats) (err error) {
	err = s.client.Call(c, _setStat, arg, _noReply)
	return
}

// UpsArtMetas list passed article meta of ups
func (s *Service) UpsArtMetas(c context.Context, arg *model.ArgUpsArts) (res map[int64][]*model.Meta, err error) {
	err = s.client.Call(c, _upsArtMetas, arg, &res)
	return
}

// ArticleMetas get article metas by aids
func (s *Service) ArticleMetas(c context.Context, arg *model.ArgAids) (res map[int64]*model.Meta, err error) {
	err = s.client.Call(c, _artMetas, arg, &res)
	return
}

// UpdateRecommends updates recommended articles.
func (s *Service) UpdateRecommends(c context.Context) (err error) {
	err = s.client.Call(c, _updateRecommends, _noArg, _noReply)
	return
}

// Recommends list recommend articles
func (s *Service) Recommends(c context.Context, arg *model.ArgRecommends) (res []*model.Meta, err error) {
	err = s.client.Call(c, _recommends, arg, &res)
	return
}

// UpArtMetas list up's article list
func (s *Service) UpArtMetas(c context.Context, arg *model.ArgUpArts) (res *model.UpArtMetas, err error) {
	err = s.client.Call(c, _upArtMetas, arg, &res)
	return
}

// AddArtDraft add article draft.
func (s *Service) AddArtDraft(c context.Context, arg *model.ArgArticle) (id int64, err error) {
	err = s.client.Call(c, _addArtDraft, arg, &id)
	return
}

// DelArtDraft deletes draft.
func (s *Service) DelArtDraft(c context.Context, arg *model.ArgAidMid) (err error) {
	err = s.client.Call(c, _delArtDraft, arg, _noReply)
	return
}

// ArtDraft get article draft by id
func (s *Service) ArtDraft(c context.Context, arg *model.ArgAidMid) (res *model.Draft, err error) {
	err = s.client.Call(c, _artDraft, arg, &res)
	return
}

// UpperDrafts get article drafts by mid
func (s *Service) UpperDrafts(c context.Context, arg *model.ArgUpDraft) (res *model.Drafts, err error) {
	err = s.client.Call(c, _upperDrafts, arg, &res)
	return
}

// ArticleRemainCount returns the number that user could be use to posting new articles.
func (s *Service) ArticleRemainCount(c context.Context, arg *model.ArgMid) (res int, err error) {
	err = s.client.Call(c, _articleRemainCount, arg, &res)
	return
}

// DelRecommendArtCache del recommend article cache
func (s *Service) DelRecommendArtCache(c context.Context, arg *model.ArgAidCid) (err error) {
	err = s.client.Boardcast(c, _delRecommendArtCache, arg, _noReply)
	return
}

// Favorites list user's favorite articles
func (s *Service) Favorites(c context.Context, arg *model.ArgFav) (res []*model.Favorite, err error) {
	err = s.client.Call(c, _favorites, arg, &res)
	return
}

// UpdateAuthorCache update author cache
func (s *Service) UpdateAuthorCache(c context.Context, arg *model.ArgAuthor) (err error) {
	err = s.client.Call(c, _updateAuthorCache, arg, _noReply)
	return
}

// UpdateSortCache update sort cache
func (s *Service) UpdateSortCache(c context.Context, arg *model.ArgSort) (err error) {
	err = s.client.Call(c, _updateSortCache, arg, _noReply)
	return
}

// IsAuthor checks that whether user has permission to write model.
func (s *Service) IsAuthor(c context.Context, arg *model.ArgMid) (res bool, err error) {
	err = s.client.Call(c, _isAuthor, arg, &res)
	return
}

// NewArticleCount get new article count since given pubtime
func (s *Service) NewArticleCount(c context.Context, arg *model.ArgNewArt) (res int64, err error) {
	err = s.client.Call(c, _newArticleCount, arg, &res)
	return
}

// HadLikesByMid check user if has liked articles
func (s *Service) HadLikesByMid(c context.Context, arg *model.ArgMidAids) (res map[int64]int, err error) {
	err = s.client.Call(c, _hadLikesByMid, arg, &res)
	return
}

// UpMoreArts get upper more arts
func (s *Service) UpMoreArts(c context.Context, arg *model.ArgAid) (res []*model.Meta, err error) {
	err = s.client.Call(c, _upMoreArts, arg, &res)
	return
}

// CreationUpStat creation up stat
func (s *Service) CreationUpStat(c context.Context, arg *model.ArgMid) (res model.UpStat, err error) {
	err = s.client.Call(c, _creationUpStat, arg, &res)
	return
}

// CreationUpThirtyDayStat creation up thirty day stat
func (s *Service) CreationUpThirtyDayStat(c context.Context, arg *model.ArgMid) (res []*model.ThirtyDayArticle, err error) {
	err = s.client.Call(c, _creationUpThirtyDayStat, arg, &res)
	return
}

// UpLists get upper article lists
func (s *Service) UpLists(c context.Context, arg *model.ArgMid) (res model.UpLists, err error) {
	err = s.client.Call(c, _upLists, arg, &res)
	return
}

// RebuildAllListReadCount rebuild all list read count
func (s *Service) RebuildAllListReadCount(c context.Context) (err error) {
	err = s.client.Call(c, _rebuildAllRC, _noArg, _noReply)
	return
}

// UpdateHotspots update hotspots
func (s *Service) UpdateHotspots(c context.Context, arg *model.ArgForce) (err error) {
	err = s.client.Call(c, _updateHotspots, arg, _noReply)
	return
}
