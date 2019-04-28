package service

import (
	"context"

	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/log"
)

// AddFavorite add article favorite .
func (s *Service) AddFavorite(c context.Context, mid, aid, fid int64, ip string) (err error) {
	if err = s.checkArticle(c, aid); err != nil {
		return
	}
	arg := &favmdl.ArgAdd{Type: favmdl.Article, Mid: mid, Oid: aid, Fid: fid}
	if err = s.favRPC.Add(c, arg); err != nil {
		dao.PromError("rpc:添加收藏")
		log.Error("s.favRPC.Add(%+v) error(%+v)", arg, err)
	}
	return
}

// DelFavorite del article favorite .
func (s *Service) DelFavorite(c context.Context, mid, aid, fid int64, ip string) (err error) {
	arg := &favmdl.ArgDel{Type: favmdl.Article, Mid: mid, Oid: aid, Fid: fid}
	if err = s.favRPC.Del(c, arg); err != nil {
		dao.PromError("rpc:删除收藏")
		log.Error("s.favRPC.Del(%+v) error(%+v)", arg, err)
	}
	return
}

// Favorites article favorites.
func (s *Service) Favorites(c context.Context, mid, fid int64, pn, ps int, ip string) (favs *favmdl.Favorites, err error) {
	arg := &favmdl.ArgFavs{Type: favmdl.Article, Mid: mid, Fid: fid, Pn: pn, Ps: ps}
	if favs, err = s.favRPC.Favorites(c, arg); err != nil {
		dao.PromError("rpc:获取收藏列表")
		log.Error("s.favRPC.Favorites(%+v) error(%+v)", arg, err)
	}
	return
}

// IsFav return user is fav article
func (s *Service) IsFav(c context.Context, mid, aid int64) (res bool, err error) {
	arg := &favmdl.ArgIsFav{Type: favmdl.Article, Mid: mid, Oid: aid}
	if res, err = s.favRPC.IsFav(c, arg); err != nil {
		dao.PromError("rpc:是否已收藏")
		log.Error("s.favRPC.IsFav(%+v) error(%+v)", arg, err)
	}
	return
}

// Favs gets user favorite article list.
func (s *Service) Favs(c context.Context, mid, fid int64, pn, ps int, ip string) (favs []*artmdl.Favorite, page *artmdl.Page, err error) {
	var (
		a    *artmdl.Meta
		ok   bool
		fs   *favmdl.Favorites
		aids []int64
		as   = make(map[int64]*artmdl.Meta)
		ts   = make(map[int64]int64)
	)
	favs = make([]*artmdl.Favorite, 0)
	if fs, err = s.Favorites(c, mid, fid, pn, ps, ip); err != nil {
		return
	}
	page = &artmdl.Page{
		Pn:    fs.Page.Num,
		Ps:    fs.Page.Size,
		Total: fs.Page.Count,
	}
	if len(fs.List) == 0 {
		return
	}
	aids = make([]int64, 0, len(fs.List))
	for _, v := range fs.List {
		aids = append(aids, v.Oid)
		ts[v.Oid] = v.MTime.Time().Unix()
	}
	if as, err = s.ArticleMetas(c, aids); err != nil {
		return
	}
	for _, aid := range aids {
		var (
			valid = true
			meta  *artmdl.Meta
		)
		if a, ok = as[aid]; !ok {
			meta = &artmdl.Meta{ID: aid}
			valid = false
		} else {
			meta = a
		}
		favs = append(favs, &artmdl.Favorite{
			Meta:         meta,
			FavoriteTime: ts[aid],
			Valid:        valid,
		})
	}
	return
}

// ValidFavs get valid favorites
func (s *Service) ValidFavs(c context.Context, mid, fid int64, pn, ps int, ip string) (res []*artmdl.Favorite, page *artmdl.Page, err error) {
	defer func() {
		if res == nil {
			res = make([]*artmdl.Favorite, 0)
		}
	}()
	var favs []*artmdl.Favorite
	if favs, page, err = s.Favs(c, mid, fid, pn, ps, ip); err != nil {
		return
	}
	for _, fav := range favs {
		if fav.Valid {
			res = append(res, fav)
		}
	}
	if page.Total <= ps {
		page.Total -= (len(favs) - len(res))
	}
	return
}
