package service

import (
	"context"

	"go-common/app/interface/main/web/conf"
	"go-common/app/interface/main/web/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	accmdl "go-common/app/service/main/account/api"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_sortNew  = 1
	_firstPn  = 1
	_samplePn = 1
	_samplePs = 1
	_cacheCnt = 20
)

var (
	_emptyArticleList = make([]*model.Meta, 0)
	_emptyAuthorList  = make([]*model.Info, 0)
	_emptyArtMetas    = make([]*artmdl.Meta, 0)
)

// ArticleList get article list.
func (s *Service) ArticleList(c context.Context, rid, mid int64, sort, pn, ps int, aids []int64) (res []*model.Meta, err error) {
	var (
		artMetas []*artmdl.Meta
		ip       = metadata.String(c, metadata.RemoteIP)
	)
	if pn == _firstPn {
		var arts []*artmdl.Meta
		arg := &artmdl.ArgRecommends{Aids: aids, Cid: rid, Pn: _firstPn, Ps: _cacheCnt, Sort: sort, RealIP: ip}
		if arts, err = s.art.Recommends(c, arg); err != nil {
			log.Error("s.art.Recommends(%d,%d,%d,%d) error(%v)", rid, pn, ps, sort, err)
			err = nil
		} else if len(arts) > 0 {
			s.cache.Do(c, func(c context.Context) {
				s.dao.SetArticleListCache(c, rid, sort, arts)
			})
		} else {
			arts, err = s.dao.ArticleListCache(c, rid, sort)
		}
		if len(arts) > ps {
			artMetas = arts[:ps-1]
		} else {
			artMetas = arts
		}
	} else {
		arg := &artmdl.ArgRecommends{Aids: aids, Cid: rid, Pn: pn, Ps: ps, Sort: sort, RealIP: ip}
		if artMetas, err = s.art.Recommends(c, arg); err != nil {
			log.Error("s.art.Recommends(%d,%d,%d,%d) error(%v)", rid, pn, ps, sort, err)
			return
		}
	}
	if len(artMetas) == 0 {
		res = _emptyArticleList
	} else {
		var item *model.Meta
		if mid > 0 {
			var (
				likes map[int64]int
				aids  []int64
			)
			for _, art := range artMetas {
				if art != nil {
					aids = append(aids, art.ID)
				}
			}
			if likes, err = s.art.HadLikesByMid(c, &artmdl.ArgMidAids{Mid: mid, Aids: aids, RealIP: ip}); err != nil {
				log.Error("s.art.HadLikesByMid(%d,%v) error(%v)", mid, aids, err)
				err = nil
			} else {
				for _, art := range artMetas {
					if art != nil {
						if like, ok := likes[art.ID]; ok {
							item = &model.Meta{Meta: art, Like: like}
						} else {
							item = &model.Meta{Meta: art}
						}
						res = append(res, item)
					}
				}
			}
		} else {
			for _, art := range artMetas {
				if art != nil {
					res = append(res, &model.Meta{Meta: art})
				}
			}
		}
	}
	return
}

// ArticleUpList get article up list.
func (s *Service) ArticleUpList(c context.Context, mid int64) (res []*model.Info, err error) {
	if res, err = s.articleUps(c, mid); err != nil {
		err = nil
	} else if len(res) > 0 {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetArticleUpListCache(c, res)
		})
		return
	}
	res, err = s.dao.ArticleUpListCache(c)
	if len(res) == 0 {
		res = _emptyAuthorList
	}
	return
}

// Categories get article categories list
func (s *Service) Categories(c context.Context) (res *artmdl.Categories, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if res, err = s.art.Categories(c, &artmdl.ArgIP{RealIP: ip}); err != nil {
		log.Error("s.art.Categories error(%v)", err)
	}
	return
}

func (s *Service) articleUps(c context.Context, mid int64) (res []*model.Info, err error) {
	var (
		mids       []int64
		list       []*artmdl.Meta
		cardsReply *accmdl.CardsReply
		relaReply  *accmdl.RelationsReply
		ip         = metadata.String(c, metadata.RemoteIP)
	)
	res = make([]*model.Info, 0)
	arg := &artmdl.ArgRecommends{Sort: _sortNew, Pn: 1, Ps: conf.Conf.Rule.ArtUpListGetCnt, RealIP: ip}
	if list, err = s.art.Recommends(c, arg); err != nil {
		log.Error("s.art.Recommends() error(%v)", err)
		return
	}
	listMap := make(map[int64]*artmdl.Meta, conf.Conf.Rule.ArtUpListCnt)
	for _, v := range list {
		if len(listMap) == conf.Conf.Rule.ArtUpListCnt {
			break
		}
		if _, ok := listMap[v.Author.Mid]; ok {
			continue
		}
		listMap[v.Author.Mid] = v
		mids = append(mids, v.Author.Mid)
	}
	if cardsReply, err = s.accClient.Cards3(c, &accmdl.MidsReq{Mids: mids}); err != nil {
		log.Error("s.accClient.Cards3(%v) error(%v)", mids, err)
		return
	}
	if mid > 0 {
		if relaReply, err = s.accClient.Relations3(c, &accmdl.RelationsReq{Mid: mid, Owners: mids, RealIp: ip}); err != nil {
			log.Error("s.accClient.Relations3(%d,%v) error(%v)", mid, mids, err)
			err = nil
		}
	}
	for _, mid := range mids {
		if card, ok := cardsReply.Cards[mid]; ok {
			info := &model.Info{ID: listMap[mid].ID, Title: listMap[mid].Title, PublishTime: listMap[mid].PublishTime}
			info.FromCard(card)
			if relaReply != nil {
				if relation, ok := relaReply.Relations[mid]; ok {
					info.Following = relation.Following
				}
			}
			res = append(res, info)
		}
	}
	return
}

// NewCount get new publish article count
func (s *Service) NewCount(c context.Context, pubTime int64) (count int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if count, err = s.art.NewArticleCount(c, &artmdl.ArgNewArt{PubTime: pubTime, RealIP: ip}); err != nil {
		log.Error("s.art.NewArticleCount(%d) error(%v)", pubTime, err)
	}
	return
}

// UpMoreArts get up more articles
func (s *Service) UpMoreArts(c context.Context, aid int64) (res []*artmdl.Meta, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if res, err = s.art.UpMoreArts(c, &artmdl.ArgAid{Aid: aid, RealIP: ip}); err != nil {
		log.Error("s.art.UpMoreArts(%d) error(%v)", aid, err)
		return
	}
	if len(res) == 0 {
		res = _emptyArtMetas
	}
	return
}
