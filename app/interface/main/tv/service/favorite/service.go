package favorite

import (
	"context"

	"go-common/app/interface/main/tv/conf"
	"go-common/app/interface/main/tv/dao/archive"
	"go-common/app/interface/main/tv/dao/favorite"
	"go-common/app/interface/main/tv/model"
	arcwar "go-common/app/service/main/archive/api"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Service .
type Service struct {
	conf   *conf.Config
	dao    *favorite.Dao
	arcDao *archive.Dao
}

// New .
func New(c *conf.Config) *Service {
	srv := &Service{
		conf:   c,
		dao:    favorite.New(c),
		arcDao: archive.New(c),
	}
	return srv
}

const (
	_ActAdd = 1
	_ActDel = 2
)

// Favorites picks one page of the member's favorites
func (s *Service) Favorites(ctx context.Context, req *model.ReqFav) (resM *model.FavMList, err error) {
	var (
		res     *favmdl.Favorites
		arcs    map[int64]*arcwar.Arc
		aids    []int64
		pageNum int
	)
	resM = &model.FavMList{}
	resM.Page.Size = s.conf.Cfg.FavPs
	if res, err = s.dao.FavoriteV3(ctx, req.MID, req.Pn); err != nil { // pick favorite original data
		log.Error("FavoriteV3 Mid %d, Pn %d, GetFav Err %v", req.MID, req.Pn, err)
		return
	}
	if len(res.List) == 0 {
		return
	}
	resM.Page = res.Page
	// temp logic because client misuses the count as the number of pages
	if resM.Page.Count%resM.Page.Size == 0 {
		pageNum = resM.Page.Count / resM.Page.Size
	} else {
		pageNum = resM.Page.Count/resM.Page.Size + 1
	}
	resM.Page.Count = pageNum
	// temp logic
	for _, v := range res.List { // combine aids and get the archive info
		aids = append(aids, v.Oid)
	}
	if arcs, err = s.arcDao.Archives(ctx, aids); err != nil {
		log.Error("FavoriteV3 Mid %d, Pn %d, GetArc Err #%v", req.MID, req.Pn, err)
		return
	}
	for _, v := range res.List { // arrange the final result
		if arc, ok := arcs[v.Oid]; ok {
			resM.List = append(resM.List, arc)
		} else {
			log.Warn("FavoriteV3 Mid %d, Pn %d, Miss Arc Info %d", req.MID, req.Pn, v.Oid)
		}
	}
	return
}

// FavAct is favorite action, add or delete
func (s *Service) FavAct(ctx context.Context, req *model.ReqFavAct) (err error) {
	if req.Action == _ActAdd {
		return s.dao.FavAdd(ctx, req.MID, req.AID)
	} else if req.Action == _ActDel {
		return s.dao.FavDel(ctx, req.MID, req.AID)
	}
	return ecode.RequestErr
}
