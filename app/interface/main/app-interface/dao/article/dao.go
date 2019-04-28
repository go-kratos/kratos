package article

import (
	"context"

	"go-common/app/interface/main/app-interface/conf"
	article "go-common/app/interface/openplatform/article/model"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is atticle dao
type Dao struct {
	artRPC *artrpc.Service
}

// New initial tag dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		artRPC: artrpc.New(c.ArticleRPC),
	}
	return
}

// UpArticles get article data from api.
func (d *Dao) UpArticles(c context.Context, mid int64, pn, ps int) (ams []*article.Meta, count int, err error) {
	var (
		res *article.UpArtMetas
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	arg := &article.ArgUpArts{Mid: mid, Pn: pn, Ps: ps, RealIP: ip}
	if res, err = d.artRPC.UpArtMetas(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	if res != nil {
		ams = res.Articles
		count = res.Count
	}
	return
}

// Favorites get article data from api.
func (d *Dao) Favorites(c context.Context, mid int64, pn, ps int) (res []*article.Favorite, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &article.ArgFav{Mid: mid, Pn: pn, Ps: ps, RealIP: ip}
	if res, err = d.artRPC.Favorites(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

func (d *Dao) Articles(c context.Context, aids []int64) (res map[int64]*article.Meta, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &article.ArgAids{Aids: aids, RealIP: ip}
	if res, err = d.artRPC.ArticleMetas(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

func (d *Dao) UpLists(c context.Context, mid int64) (lists []*article.List, count int, err error) {
	var (
		res article.UpLists
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	arg := &article.ArgMid{Mid: mid, RealIP: ip}
	if res, err = d.artRPC.UpLists(c, arg); err != nil {
		err = errors.Wrapf(err, "%+v", arg)
		return
	}
	lists = res.Lists
	count = res.Total
	return
}
