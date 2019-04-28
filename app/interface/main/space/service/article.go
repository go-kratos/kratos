package service

import (
	"context"

	"go-common/app/interface/main/space/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var _emptyArticle = make([]*artmdl.Meta, 0)

// Article get articles by upMid.
func (s *Service) Article(c context.Context, mid int64, pn, ps, sort int) (res *artmdl.UpArtMetas, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if res, err = s.art.UpArtMetas(c, &artmdl.ArgUpArts{Mid: mid, Pn: pn, Ps: ps, Sort: sort, RealIP: ip}); err != nil {
		log.Error("s.art.UpArtMetas(%d,%d,%d) error(%v)", mid, pn, ps, err)
		return
	}
	if res != nil && len(res.Articles) == 0 {
		res.Articles = _emptyArticle
	}
	return
}

// UpArtStat get up all article stat.
func (s *Service) UpArtStat(c context.Context, mid int64) (data *model.UpArtStat, err error) {
	addCache := true
	if data, err = s.dao.UpArtCache(c, mid); err != nil {
		addCache = false
	} else if data != nil {
		return
	}
	if data, err = s.dao.UpArtStat(c, mid); data != nil && addCache {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetUpArtCache(c, mid, data)
		})
	}
	return
}
