package service

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"go-common/app/service/main/favorite/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// Relations return relations info by fid.
func (s *Service) Relations(c context.Context, typ int8, mid, uid, fid int64, tid, tv, pn, ps int, keyword, order string) (res *model.Favorites, err error) {
	if uid > 0 {
		mid = uid
	}
	if order == "" {
		order = model.SortMtime
	}
	res = new(model.Favorites)
	res.Page.Num = pn
	res.Page.Size = ps
	un, err := s.favDao.FavedBit(c, typ, mid)
	if err != nil || un {
		return
	}
	folders := make(map[int64]*model.Folder)
	if fid == int64(model.AllFidFlag) {
		if folders, err = s.userFolders(c, typ, mid); err != nil {
			log.Error("s.userFolders(%d,%d) error(%v)", typ, mid, err)
			return nil, err
		}
	} else {
		fd, err1 := s.folder(c, typ, mid, fid)
		if err1 != nil {
			log.Error("s.folder(%d,%d,%d) error(%v)", typ, mid, fid, err1)
			return nil, err1
		}
		folders[fd.ID] = fd
	}
	var oids []int64
	favs := make([]*model.Favorite, 0)
	favMap := make(map[int64]*model.Favorite)
	if len(folders) == 1 {
		for _, f := range folders {
			unExpired, err1 := s.favDao.ExpireRelations(c, mid, f.ID)
			if err1 != nil {
				log.Error("s.favDao.ExpireRelations(%d,%d) error(%v)", mid, f.ID, err1)
				return nil, err1
			}
			if favs, err = s.parRelations(c, typ, mid, f.ID, f.Count, unExpired); err != nil {
				return nil, err
			}
			for _, f := range favs {
				oids = append(oids, f.Oid)
				favMap[f.Oid] = f
			}
		}
	} else {
		var (
			fids  []int64
			okMap map[int64]bool
		)
		for fid := range folders {
			fids = append(fids, fid)
		}
		if okMap, err = s.favDao.MultiExpireRelations(c, mid, fids); err != nil {
			log.Error("s.favDao.MultiExpireRelations(%v,%d) error(%v)", fids, err)
			return
		}
		g := new(errgroup.Group)
		mux := new(sync.Mutex)
		for _, v := range folders {
			f := v
			unExpired := okMap[f.ID]
			g.Go(func() error {
				gfavs, err1 := s.parRelations(c, typ, mid, f.ID, f.Count, unExpired)
				if err1 != nil {
					log.Error("s.parRelations(%d,%d,%d,%d,%t) error(%v)", typ, mid, f.ID, f.Count, unExpired, err1)
				}
				mux.Lock()
				favs = append(favs, gfavs...)
				mux.Unlock()
				return nil
			})
		}
		if err = g.Wait(); err != nil {
			return
		}
		sort.Slice(favs, func(i, j int) bool {
			return favs[i].MTime < favs[j].MTime
		})
		ts := make([]*model.Favorite, 0)
		for _, f := range favs {
			if _, ok := favMap[f.Oid]; ok {
				continue
			}
			ts = append(ts, f)
			oids = append(oids, f.Oid)
			favMap[f.Oid] = f
		}
		favs = ts
	}
	favs = make([]*model.Favorite, 0)
	res.List = _emptyFavRelations
	if keyword != "" { // riot-search
		resp, err1 := s.favDao.SearchArcs(c, keyword, oids, pn, ps)
		if err1 != nil || resp == nil {
			return res, err1
		}
		res.Page.Count = resp.Page.Total
		for _, oid := range resp.IDs {
			if v, ok := favMap[int64(oid)]; ok {
				favs = append(favs, v)
			}
		}
	} else { // rank-service
		rankPn := pn
		rankPs := ps
		if order == model.SortMtime {
			rankPn = 1
			rankPs = len(oids)
		}
		resp, err1 := s.SortArcsRPC(c, tid, tv, rankPn, rankPs, order, oids)
		if err1 != nil || resp == nil {
			return res, err1
		}
		res.Page.Count = resp.Page.Total
		for _, oid := range resp.Result {
			if v, ok := favMap[oid]; ok {
				favs = append(favs, v)
			}
		}
		if order == model.SortMtime {
			sort.Slice(favs, func(i, j int) bool {
				return favs[i].MTime > favs[j].MTime
			})
			start := (pn - 1) * ps
			end := pn * ps
			if start > len(favs) {
				res.List = make([]*model.Favorite, 0)
				return
			}
			if end > len(favs) {
				end = len(favs)
			}
			res.List = favs[start:end]
			return
		}
	}
	res.List = favs
	return
}

// Tlists return lists info by fid.
func (s *Service) Tlists(c context.Context, typ int8, mid, uid, fid int64) (res []*model.Partition, err error) {
	if uid > 0 {
		mid = uid
	}
	res = make([]*model.Partition, 0)
	un, err := s.favDao.FavedBit(c, typ, mid)
	if err != nil || un {
		return
	}
	f, err := s.folder(c, typ, mid, fid)
	if err != nil {
		log.Error("s.folder(%d,%d,%d) error(%v)", typ, mid, fid, err)
		return nil, err
	}
	var favs []*model.Favorite
	unExpired, err := s.favDao.ExpireRelations(c, mid, f.ID)
	if err != nil {
		log.Error("s.favDao.ExpireRelations(%d,%d) error(%v)", mid, f.ID, err)
		return
	}
	if favs, err = s.parRelations(c, typ, mid, f.ID, f.Count, unExpired); err != nil {
		return
	}
	var oids []int64
	for _, fav := range favs {
		oids = append(oids, fav.Oid)
	}
	typeids, err := s.TypeidsRPC(c, oids)
	if err != nil || typeids == nil {
		return
	}
	for _, v := range typeids.List {
		p := new(model.Partition)
		p.Count = v.Count
		v, _ := strconv.Atoi(v.Key)
		p.Tid = v
		res = append(res, p)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Tid < res[j].Tid
	})
	return
}
