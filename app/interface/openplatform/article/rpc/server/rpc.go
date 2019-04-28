package server

import (
	"go-common/app/interface/openplatform/article/model"
	"go-common/app/interface/openplatform/article/service"
	"go-common/library/ecode"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC .
type RPC struct {
	s *service.Service
}

// New creates rpc server.
func New(s *service.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
	svr = rpc.NewServer(nil)
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

// AddArticle adds article when article passed. purge cache.
func (r *RPC) AddArticle(c context.Context, arg *model.ArgArticle, res *int64) (err error) {
	var (
		aid int64
		a   = model.TransformArticle(arg)
	)
	if aid, err = r.s.AddArticle(c, a, arg.ActivityID, arg.ListID, arg.RealIP); err != nil {
		return
	}
	*res = aid
	return
}

// UpdateArticle updates article when article passed. purge cache.
func (r *RPC) UpdateArticle(c context.Context, arg *model.ArgArticle, res *struct{}) (err error) {
	var a = model.TransformArticle(arg)
	if err = r.s.UpdateArticle(c, a, arg.ActivityID, arg.ListID, arg.RealIP); err != nil {
		return
	}
	return
}

// CreationArticle gets article info.
func (r *RPC) CreationArticle(c context.Context, arg *model.ArgAidMid, res *model.Article) (err error) {
	var rr *model.Article
	if rr, err = r.s.CreationArticle(c, arg.Aid, arg.Mid); err == nil {
		*res = *rr
	}
	return
}

// DelArticle drops article when article not passed. purge cache.
func (r *RPC) DelArticle(c context.Context, arg *model.ArgAidMid, res *struct{}) error {
	return r.s.DelArticle(c, arg.Aid, arg.Mid)
}

// AddArticleCache adds article cache.
func (r *RPC) AddArticleCache(c context.Context, arg *model.ArgAid, res *struct{}) (err error) {
	err = r.s.AddArticleCache(c, arg.Aid)
	return
}

// UpdateArticleCache adds article cache.
func (r *RPC) UpdateArticleCache(c context.Context, arg *model.ArgAidCid, res *struct{}) (err error) {
	err = r.s.UpdateArticleCache(c, arg.Aid, arg.Cid)
	return
}

// DelArticleCache adds article cache.
func (r *RPC) DelArticleCache(c context.Context, arg *model.ArgAidMid, res *struct{}) (err error) {
	err = r.s.DelArticleCache(c, arg.Mid, arg.Aid)
	return
}

// CreationWithdrawArticle author withdraw model.
func (r *RPC) CreationWithdrawArticle(c context.Context, arg *model.ArgAidMid, res *struct{}) (err error) {
	err = r.s.CreationWithdrawArticle(c, arg.Mid, arg.Aid)
	return
}

// Categories list article categories
func (r *RPC) Categories(c context.Context, arg *model.ArgIP, res *model.Categories) (err error) {
	*res, err = r.s.ListCategories(c, arg.RealIP)
	return
}

// CategoriesMap list article categories map
func (r *RPC) CategoriesMap(c context.Context, arg *model.ArgIP, res *map[int64]*model.Category) (err error) {
	*res, err = r.s.ListCategoriesMap(c, arg.RealIP)
	return
}

// CreationUpperArticles gets upper's article list for creation center.
func (r *RPC) CreationUpperArticles(c context.Context, arg *model.ArgCreationArts, res *model.CreationArts) (err error) {
	var rr *model.CreationArts
	if rr, err = r.s.CreationUpperArticlesMeta(c, arg.Mid, arg.Group, arg.Category, arg.Sort, arg.Pn, arg.Ps, arg.RealIP); err == nil {
		*res = *rr
	}
	return
}

// SetStat set all stat cache(redis)
func (r *RPC) SetStat(c context.Context, arg *model.ArgStats, res *struct{}) (err error) {
	err = r.s.SetStat(c, arg.Aid, arg.Stats)
	return
}

// UpsArtMetas list passed article meta of ups
func (r *RPC) UpsArtMetas(c context.Context, arg *model.ArgUpsArts, res *map[int64][]*model.Meta) (err error) {
	if arg.Pn < 1 {
		arg.Pn = 1
	}
	var (
		start = (arg.Pn - 1) * arg.Ps
		end   = start + arg.Ps - 1
	)
	*res, err = r.s.UpsArticleMetas(c, arg.Mids, start, end)
	return
}

// UpArtMetas list up article metas
func (r *RPC) UpArtMetas(c context.Context, arg *model.ArgUpArts, res *model.UpArtMetas) (err error) {
	if arg.Pn < 1 {
		arg.Pn = 1
	}
	if arg.Ps <= 0 {
		err = ecode.RequestErr
		return
	}
	var rr *model.UpArtMetas
	if rr, err = r.s.UpArticleMetas(c, arg.Mid, arg.Pn, arg.Ps, arg.Sort); err == nil {
		*res = *rr
	}
	return
}

// ArticleMetas list article metas
func (r *RPC) ArticleMetas(c context.Context, arg *model.ArgAids, res *map[int64]*model.Meta) (err error) {
	*res, err = r.s.ArticleMetas(c, arg.Aids)
	return
}

// UpdateRecommends refresh recommend data.
func (r *RPC) UpdateRecommends(c context.Context, arg *model.ArgIP, res *struct{}) (err error) {
	err = r.s.UpdateRecommends(c)
	return
}

// Recommends list recommend articles
func (r *RPC) Recommends(c context.Context, arg *model.ArgRecommends, res *[]*model.Meta) (err error) {
	rs, err := r.s.Recommends(c, arg.Cid, arg.Pn, arg.Ps, arg.Aids, arg.Sort)
	if err != nil {
		return
	}
	for _, r := range rs {
		*res = append(*res, &r.Meta)
	}
	return
}

// AddArtDraft adds or updates draft.
func (r *RPC) AddArtDraft(c context.Context, arg *model.ArgArticle, res *int64) (err error) {
	d := model.TransformDraft(arg)
	*res, err = r.s.AddArtDraft(c, d)
	return
}

// DelArtDraft .
func (r *RPC) DelArtDraft(c context.Context, arg *model.ArgAidMid, res *struct{}) (err error) {
	err = r.s.DelArtDraft(c, arg.Aid, arg.Mid)
	return
}

// ArtDraft get article draft
func (r *RPC) ArtDraft(c context.Context, arg *model.ArgAidMid, res *model.Draft) (err error) {
	var v *model.Draft
	if v, err = r.s.ArtDraft(c, arg.Aid, arg.Mid); err == nil && v != nil {
		*res = *v
	}
	return
}

// UpperDrafts get article drafts by mid
func (r *RPC) UpperDrafts(c context.Context, arg *model.ArgUpDraft, res *model.Drafts) (err error) {
	var v *model.Drafts
	if v, err = r.s.UpperDrafts(c, arg.Mid, arg.Pn, arg.Ps); err == nil && v != nil {
		*res = *v
	}
	return
}

// ArticleRemainCount returns the number that user could be use to posting new articles.
func (r *RPC) ArticleRemainCount(c context.Context, arg *model.ArgMid, res *int) (err error) {
	*res, err = r.s.ArticleRemainCount(c, arg.Mid)
	return
}

// DelRecommendArtCache del recommend article cache
func (r *RPC) DelRecommendArtCache(c context.Context, arg *model.ArgAidCid, res *struct{}) (err error) {
	err = r.s.DelRecommendArtCache(c, arg.Aid, arg.Cid)
	return
}

// Favorites list user's favorite articles
func (r *RPC) Favorites(c context.Context, arg *model.ArgFav, res *[]*model.Favorite) (err error) {
	*res, _, err = r.s.ValidFavs(c, arg.Mid, 0, arg.Pn, arg.Ps, arg.RealIP)
	return
}

// UpdateAuthorCache update author cache
func (r *RPC) UpdateAuthorCache(c context.Context, arg *model.ArgAuthor, res *struct{}) (err error) {
	err = r.s.UpdateAuthorCache(c, arg.Mid)
	return
}

// UpdateSortCache update sort cache
func (r *RPC) UpdateSortCache(c context.Context, arg *model.ArgSort, res *struct{}) (err error) {
	err = r.s.UpdateSortCache(c, arg.Aid, arg.Changed, arg.RealIP)
	return
}

// IsAuthor checks that whether user has permission to write model.
func (r *RPC) IsAuthor(c context.Context, arg *model.ArgMid, res *bool) (err error) {
	*res, _, err = r.s.IsAuthor(c, arg.Mid)
	return
}

// NewArticleCount get new article count since given pubtime
func (r *RPC) NewArticleCount(c context.Context, arg *model.ArgNewArt, res *int64) (err error) {
	*res, err = r.s.NewArticleCount(c, arg.PubTime)
	return
}

// HadLikesByMid check user if has liked articles
func (r *RPC) HadLikesByMid(c context.Context, arg *model.ArgMidAids, res *map[int64]int8) (err error) {
	*res, err = r.s.HadLikesByMid(c, arg.Mid, arg.Aids)
	return
}

// UpMoreArts get upper more arts
func (r *RPC) UpMoreArts(c context.Context, arg *model.ArgAid, res *[]*model.Meta) (err error) {
	*res, err = r.s.MoreArts(c, arg.Aid)
	return
}

// CreationUpStat creation up stat
func (r *RPC) CreationUpStat(c context.Context, arg *model.ArgMid, res *model.UpStat) (err error) {
	*res, err = r.s.UpStat(c, arg.Mid)
	return
}

// CreationUpThirtyDayStat creation up thirty day stat
func (r *RPC) CreationUpThirtyDayStat(c context.Context, arg *model.ArgMid, res *[]*model.ThirtyDayArticle) (err error) {
	*res, err = r.s.UpThirtyDayStat(c, arg.Mid)
	return
}

// UpLists get upper article lists
func (r *RPC) UpLists(c context.Context, arg *model.ArgMid, res *model.UpLists) (err error) {
	*res, err = r.s.UpLists(c, arg.Mid, model.ListSortPtime)
	return
}

// RebuildAllListReadCount rebuild all list read count
func (r *RPC) RebuildAllListReadCount(c context.Context, arg *struct{}, res *struct{}) (err error) {
	r.s.RebuildAllListReadCount(c)
	return
}

// UpdateHotspots update hotspots
func (r *RPC) UpdateHotspots(c context.Context, arg *model.ArgForce, res *struct{}) (err error) {
	err = r.s.UpdateHotspots(arg.Force)
	return
}
