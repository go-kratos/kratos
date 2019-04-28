package service

import (
	"context"
	"html/template"
	"math"

	musicmdl "go-common/app/interface/main/favorite/model"
	"go-common/app/service/main/archive/api"
	pb "go-common/app/service/main/favorite/api"
	"go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var (
	_emptyArchives     = []*model.FavArchive{}
	_emptyVideoFolders = []*model.VideoFolder{}
)

// newCovers get three cover of each fid
func (s *Service) newCovers(c context.Context, mid int64, recents map[int64][]*model.Resource) (fcvs map[int64][]*model.Cover, err error) {
	var (
		fids, misFids []int64
	)
	for fid := range recents {
		fids = append(fids, fid)
	}
	if fcvs, misFids, err = s.videoDao.NewCoversCache(c, mid, fids); err != nil {
		log.Error("s.videoDao.NewCoversCache(fids %v) err(%v)", fids, err)
		return
	}
	// get miss cover from db
	if len(misFids) > 0 {
		var (
			allAids []int64
			allMids []int64
			as      map[int64]*api.Arc
			ms      map[int64]*musicmdl.Music
		)
		for _, fid := range misFids {
			if resources, ok := recents[fid]; ok {
				for _, res := range resources {
					if int8(res.Typ) == model.TypeVideo {
						allAids = append(allAids, res.Oid)
					} else {
						allMids = append(allMids, res.Oid)
					}
				}
			}
		}
		if len(allAids) > 0 {
			if as, err = s.arcsRPC(c, allAids); err != nil {
				log.Error("s.arcsRPC(aids %v),err(%v)", allAids, err)
				return
			}
		}
		if len(allMids) > 0 {
			if ms, err = s.musicDao.MusicMap(c, allMids); err != nil {
				log.Error("s.musicMap(allMids %v),err(%v)", allMids, err)
				return
			}
		}
		// set miss fid's cover
		for _, misFid := range misFids {
			fid := misFid
			cvs := make([]*model.Cover, 0)
			for _, res := range recents[fid] {
				cv := &model.Cover{}
				cv.Aid = res.Oid
				cv.Type = res.Typ
				if int8(res.Typ) == model.TypeVideo {
					if arc, ok := as[res.Oid]; ok {
						if !arc.IsNormal() {
							continue
						}
						cv.Pic = arc.Pic
						cvs = append(cvs, cv)
					}
				} else {
					if music, ok := ms[res.Oid]; ok {
						cv.Pic = music.Cover
						cvs = append(cvs, cv)
					}
				}
			}
			if err := s.videoDao.SetNewCoverCache(c, mid, fid, cvs); err != nil {
				log.Error("s.videoDao.SetNewCoverCache(%d,%d,%v) error(%v)", mid, fid, cvs, err)
			}
			fcvs[fid] = cvs
		}
	}
	return
}

func (s *Service) arcsRPC(c context.Context, aids []int64) (arcsMap map[int64]*api.Arc, err error) {
	var (
		batch = s.conf.Fav.MaxPagesize
		arcs  map[int64]*api.Arc
	)
	arcsMap = make(map[int64]*api.Arc, len(aids))
	for len(aids) > 0 {
		if len(aids) < batch {
			batch = len(aids)
		}
		arcs, err = s.ArcsRPC(c, aids[:batch])
		if err != nil {
			log.Error("s.ArcsRPC(%v) error(%v)", aids[:batch], err)
			return
		}
		for k, v := range arcs {
			arcsMap[k] = v
		}
		aids = aids[batch:]
	}
	return
}

func (s *Service) normalArcs(c context.Context, aids []int64) (as []*api.Arc, err error) {
	arcs, err := s.ArcsRPC(c, aids)
	if err != nil {
		log.Error("s.ArcsRPC(%v) error(%v)", aids, err)
		return
	}
	as = make([]*api.Arc, 0, len(aids))
	for _, v := range aids {
		if a, ok := arcs[v]; ok {
			as = append(as, a)
		}
	}
	return
}

// FavVideo get fav videos from search or db
func (s *Service) FavVideo(c context.Context, mid, vmid, uid, fid int64, keyword, order string, tid, pn, ps int) (sv *model.SearchArchive, err error) {
	if order != "click" && order != "pubdate" {
		order = model.SortMtime
	}
	if order == "pubdate" {
		order = model.SortPubtime
	}
	sv = new(model.SearchArchive)
	var favs *model.Favorites
	if favs, err = s.FavoritesRPC(c, model.TypeVideo, mid, vmid, fid, tid, keyword, order, pn, ps); err != nil || favs == nil {
		log.Error("s.FavoritesRPC(%d,%d,%d,%d,%d,%d) error(%v)", model.TypeVideo, mid, vmid, fid, pn, ps, err)
		return
	}
	if len(favs.List) == 0 {
		sv.Result = nil
		sv.Archives = _emptyArchives
		return
	}
	if err = s.newFillArchives(c, favs.List, sv); err != nil {
		log.Error("s.newFillArchives error(%v)", err)
	}
	sv.Mid = uid
	sv.Fid = fid
	sv.Tid = tid
	sv.Order = order
	sv.Keyword = keyword
	sv.PageCount = int(math.Ceil(float64(favs.Page.Count) / float64(s.conf.Fav.MaxPagesize)))
	sv.PageSize = s.conf.Fav.MaxPagesize
	sv.Total = favs.Page.Count
	sv.NumPages = 0
	sv.NumResults = 0
	return
}

// TidList get video folder type names from search.
func (s *Service) TidList(c context.Context, mid, vmid, uid, fid int64) (res []*model.Partition, err error) {
	tidCounts, err := s.TlistsRPC(c, model.TypeVideo, mid, vmid, fid)
	if err != nil {
		log.Error("s.TlistsRPC(%d,%d,%d,%d) error(%v)", model.TypeVideo, mid, vmid, fid, err)
		return
	}
	types, err := s.TypesRPC(c)
	if err != nil {
		log.Error("s.TypesRPC() error(%v)", err)
		return
	}
	for _, t := range tidCounts {
		if t.Tid == 0 {
			continue
		}
		if v, ok := types[int16(t.Tid)]; ok {
			t.Name = v.Name
		}
		res = append(res, t)
	}
	return
}

func (s *Service) archive(c context.Context, aid int64) (arc *api.Arc, err error) {
	arc, err = s.ArcRPC(c, aid)
	if err != nil {
		log.Error("s.ArcRPC(%d), error(%v)", aid, err)
	}
	if !arc.IsNormal() {
		err = ecode.ArchiveNotExist
	}
	return
}

// FavFolders get mid user's favorites.
func (s *Service) FavFolders(c context.Context, mid, vmid, uid, aid int64, isSelf bool, mediaList bool, fromWeb bool) (res []*model.VideoFolder, err error) {
	var fs []*model.Folder
	typ := model.TypeVideo
	ip := metadata.String(c, metadata.RemoteIP)
	if fs, err = s.AllFoldersRPC(c, typ, mid, vmid, aid, ip); err != nil {
		log.Error("s.AllFoldersRPC(%d,%d,%d,%d,%s) error(%v)", typ, mid, vmid, ip, err)
		return
	}
	if len(fs) == 0 {
		res = _emptyVideoFolders
		return
	}
	faids := make(map[int64][]*model.Resource, len(fs))
	for _, f := range fs {
		faids[f.ID] = f.RecentRes
	}
	var covers map[int64][]*model.Cover
	if covers, err = s.newCovers(c, uid, faids); err != nil {
		log.Error("s.newCovers(%d,%v) error(%v)", uid, faids, err)
		err = nil
	}
	//兼容老的缓存，后期下掉
	for _, cover := range covers {
		for _, co := range cover {
			if co.Type == 0 {
				co.Type = int32(model.TypeVideo)
			}
		}
	}
	for _, f := range fs {
		maxCount := model.DefaultFolderLimit
		if !f.IsDefault() {
			maxCount = model.NormalFolderLimit
		}
		name := f.Name
		if fromWeb { // web端html转义
			name = template.HTMLEscapeString(name)
		}
		if mediaList {
			if f.IsDefault() && name == "默认收藏夹" {
				name = "默认播单"
			}
		}
		cover := []*model.Cover{}
		if mediaList {
			if f.Cover != "" {
				cover = []*model.Cover{{
					Type: 0,
					Pic:  f.Cover,
				}}
			} else {
				cover = covers[f.ID]
			}
		} else {
			for _, co := range covers[f.ID] {
				if int8(co.Type) == model.TypeVideo {
					cover = append(cover, co)
				}
			}
		}
		vf := &model.VideoFolder{
			MediaId:  f.ID*100 + f.Mid%100,
			Fid:      f.ID,
			Mid:      f.Mid,
			Name:     name,
			MaxCount: maxCount,
			CurCount: f.Count,
			Favoured: f.Favored,
			State:    int8(f.Attr & 3),
			CTime:    f.CTime,
			MTime:    f.MTime,
			Cover:    cover,
		}
		res = append(res, vf)
	}

	return
}

// AddFavFolder add a new favorite folder
func (s *Service) AddFavFolder(c context.Context, mid int64, name, cookie, accessKey string, state int32) (fid int64, err error) {
	var reply *pb.AddFolderReply
	reply, err = s.favClient.AddFolder(c, &pb.AddFolderReq{
		Typ:       int32(model.TypeVideo),
		Mid:       mid,
		Name:      name,
		Cookie:    cookie,
		AccessKey: accessKey,
		Public:    state,
	})
	if err != nil {
		return
	}
	fid = reply.Fid
	return
}

// UpFavName update favorite name.
func (s *Service) UpFavName(c context.Context, mid, fid int64, name, cookie, accessKey string) (err error) {
	_, err = s.favClient.UpFolderName(c, &pb.UpFolderNameReq{
		Typ:       int32(model.TypeVideo),
		Fid:       fid,
		Mid:       mid,
		Name:      name,
		Cookie:    cookie,
		AccessKey: accessKey,
	})
	return
}

// SetVideoFolderSort set folder sort.
func (s *Service) SetVideoFolderSort(c context.Context, mid int64, fids []int64) (err error) {
	_, err = s.favClient.SetFolderSort(c, &pb.SetFolderSortReq{
		Typ:  int32(model.TypeVideo),
		Mid:  mid,
		Fids: fids,
	})
	return
}

// UpFavState update folder state.
func (s *Service) UpFavState(c context.Context, mid, fid int64, public int32, cookie string, accessKey string) (err error) {
	_, err = s.favClient.UpFolderAttr(c, &pb.UpFolderAttrReq{
		Typ:    int32(model.TypeVideo),
		Fid:    fid,
		Mid:    mid,
		Public: public,
	})
	return
}

// DelVideoFolder delete favFolder and push databus msg to del videos in folder.
func (s *Service) DelVideoFolder(c context.Context, mid, fid int64) (err error) {
	_, err = s.favClient.DelFolder(c, &pb.DelFolderReq{
		Typ: int32(model.TypeVideo),
		Mid: mid,
		Fid: fid,
	})
	return
}

// RecentArcs return the newest archives in all folder.
func (s *Service) RecentArcs(c context.Context, mid int64, pageNum, pageSize int) (sv *model.SearchArchive, err error) {
	aids, err := s.RecentsRPC(c, model.TypeVideo, mid, pageSize)
	if err != nil {
		return
	}
	sv = new(model.SearchArchive)
	if len(aids) == 0 {
		sv.Result = nil
		sv.Archives = _emptyArchives
		return
	}
	archives, err := s.normalArcs(c, aids)
	if err != nil {
		log.Error("s.NormalArchives(%v) error(%v)", sv, err)
		return
	}
	if err = s.fillArchives(c, sv); err != nil {
		log.Error("s.RecentArcs err(%v)", err)
	}
	farchives := make([]*model.FavArchive, 0, len(archives))
	for _, arc := range archives {
		farchive := new(model.FavArchive)
		farchive.Arc = arc
		farchives = append(farchives, farchive)
	}
	sv.Result = nil
	sv.Archives = farchives
	sv.Mid = mid
	sv.PageCount = sv.NumPages
	sv.Total = sv.NumResults
	sv.NumPages = 0
	sv.NumResults = 0
	return
}

func (s *Service) fillArchives(c context.Context, sv *model.SearchArchive) (err error) {
	aids := make([]int64, 0, len(sv.Result))
	searchArcs := make(map[int64]*model.SearchArchiveResult, len(sv.Result))
	for _, v := range sv.Result {
		aids = append(aids, v.ID)
		searchArcs[v.ID] = &model.SearchArchiveResult{
			FavTime: v.FavTime,
			Title:   v.Title,
			Play:    v.Play,
		}
	}
	archives, err := s.normalArcs(c, aids)
	if err != nil {
		log.Error("s.NormalArchives(%v) error(%v)", sv, err)
		return
	}
	farchives := make([]*model.FavArchive, 0, len(archives))
	for _, arc := range archives {
		var farchive = &model.FavArchive{}
		farchive.Arc = arc
		farchive.FavAt = searchArcs[arc.Aid].FavTime
		farchive.HighlightTitle = searchArcs[arc.Aid].Title
		farchive.PlayNum = searchArcs[arc.Aid].Play
		farchives = append(farchives, farchive)
	}
	sv.Result = nil
	sv.Archives = farchives
	return
}

func (s *Service) newFillArchives(c context.Context, favorites []*model.Favorite, sv *model.SearchArchive) (err error) {
	aids := make([]int64, 0, len(sv.Result))
	favoriteArcs := make(map[int64]*model.Favorite, len(sv.Result))
	for _, v := range favorites {
		aids = append(aids, v.Oid)
		favoriteArcs[v.Oid] = v
	}
	archives, err := s.normalArcs(c, aids)
	if err != nil {
		log.Error("s.normalArcs(%v) error(%v)", aids, err)
		return
	}
	farchives := make([]*model.FavArchive, 0, len(archives))
	for _, arc := range archives {
		var farchive = &model.FavArchive{}
		farchive.Arc = arc
		farchive.FavAt = int64(favoriteArcs[arc.Aid].MTime)
		farchive.HighlightTitle = arc.Title
		farchives = append(farchives, farchive)
	}
	sv.Result = nil
	sv.Archives = farchives
	return
}

// AddArc add a archive into folder.
func (s *Service) AddArc(c context.Context, mid, fid, aid int64, ck, ak string) (err error) {
	if _, err = s.archive(c, aid); err != nil {
		return
	}
	_, err = s.favClient.AddFav(c, &pb.AddFavReq{
		Tp:  int32(model.TypeVideo),
		Mid: mid,
		Fid: fid,
		Oid: aid,
	})
	if err != nil {
		return
	}
	s.cache.Do(c, func(c context.Context) {
		if err := s.videoDao.DelCoverCache(c, mid, fid); err != nil {
			log.Error("s.videoDao.DelCoverCache(%d,%d) error(%v)", mid, fid, err)
		}
	})
	return
}

// AddArcToFolders add a archive into multi folders.
func (s *Service) AddArcToFolders(c context.Context, mid, aid int64, fids []int64, ck, ak string) (err error) {
	if len(fids) == 0 {
		err = s.AddArc(c, mid, 0, aid, ck, ak)
	}
	for _, fid := range fids {
		err = s.AddArc(c, mid, fid, aid, ck, ak)
	}
	return
}

// DelArc delete a archive from favorite.
func (s *Service) DelArc(c context.Context, mid, fid, aid int64) (err error) {
	_, err = s.favClient.DelFav(c, &pb.DelFavReq{
		Tp:  int32(model.TypeVideo),
		Mid: mid,
		Fid: fid,
		Oid: aid,
	})
	if err != nil {
		return
	}
	s.cache.Do(c, func(c context.Context) {
		if err := s.videoDao.DelCoverCache(c, mid, fid); err != nil {
			log.Error("s.videoDao.DelCoverCache(%d,%d) error(%v)", mid, fid, err)
		}
	})
	return
}

// DelArcs delete some archives from favorite.
func (s *Service) DelArcs(c context.Context, mid, fid int64, aids []int64) (err error) {
	_, err = s.favClient.MultiDel(c, &pb.MultiDelReq{
		Typ:  int32(model.TypeVideo),
		Mid:  mid,
		Fid:  fid,
		Oids: aids,
	})
	if err != nil {
		return
	}
	s.cache.Do(c, func(c context.Context) {
		if err := s.videoDao.DelCoverCache(c, mid, fid); err != nil {
			log.Error("s.videoDao.DelCoverCache(%d,%d) error(%v)", mid, fid, err)
		}
	})
	return
}

// MoveArcs move archives from old favorite to new favorite
func (s *Service) MoveArcs(c context.Context, mid, oldfid, newfid int64, aids []int64) (err error) {
	if len(aids) == 0 || oldfid == newfid {
		return
	}
	_, err = s.favClient.MoveFavs(c, &pb.MoveFavsReq{
		Typ:    int32(model.TypeVideo),
		Mid:    mid,
		OldFid: oldfid,
		NewFid: newfid,
		Oids:   aids,
	})
	if err != nil {
		return
	}
	s.cache.Do(c, func(c context.Context) {
		if err := s.videoDao.DelCoverCache(c, mid, newfid); err != nil {
			log.Error("s.videoDao.DelCoverCache(%d,%d) error(%v)", mid, newfid, err)
		}
		if err := s.videoDao.DelCoverCache(c, mid, oldfid); err != nil {
			log.Error("s.videoDao.DelCoverCache(%d,%d) error(%v)", mid, oldfid, err)
		}
	})
	return
}

// CopyArcs copy archives to other favorite.
func (s *Service) CopyArcs(c context.Context, mid, oldmid, oldfid, newfid int64, aids []int64) (err error) {
	if len(aids) == 0 || oldfid == newfid {
		return
	}
	_, err = s.favClient.CopyFavs(c, &pb.CopyFavsReq{
		Typ:    int32(model.TypeVideo),
		OldMid: oldmid,
		Mid:    mid,
		OldFid: oldfid,
		NewFid: newfid,
		Oids:   aids,
	})
	if err != nil {
		return
	}
	s.cache.Do(c, func(c context.Context) {
		if err := s.videoDao.DelCoverCache(c, mid, newfid); err != nil {
			log.Error("s.videoDao.DelCoverCache(%d,%d) error(%v)", mid, newfid, err)
		}
	})
	return
}

// IsFaveds check if aids faved by user
func (s *Service) IsFaveds(c context.Context, mid int64, aids []int64) (m map[int64]bool, err error) {
	if m, err = s.IsFavsRPC(c, model.TypeVideo, mid, aids); err != nil {
		log.Error("s.IsFavsRPC(%d,%d,%v) error(%v)", model.TypeVideo, mid, aids, err)
	}
	return
}

// IsFaved check if aid faved by user
func (s *Service) IsFaved(c context.Context, mid, aid int64) (faved bool, count int, err error) {
	if faved, err = s.IsFavRPC(c, model.TypeVideo, mid, aid); err != nil {
		log.Error("s.IsFavsRPC(%d,%d,%d) error(%v)", model.TypeVideo, mid, aid, err)
	}
	count = 1
	return
}

// InDef detemine aid whether or not archive in default folder.
func (s *Service) InDef(c context.Context, mid, aid int64) (isin bool, err error) {
	if isin, err = s.InDefaultRPC(c, model.TypeVideo, mid, aid); err != nil {
		log.Error("s.InDefaultRPC(%d,%d,%d) error(%v)", model.TypeVideo, mid, aid, err)
	}
	return
}

// CleanState return this folder clean state.
func (s *Service) CleanState(c context.Context, mid, fid int64) (cleanState int, err error) {
	reply, err := s.favClient.CleanState(c, &pb.CleanStateReq{
		Typ: int32(model.TypeVideo),
		Mid: mid,
		Fid: fid,
	})
	if err != nil {
		return
	}
	cleanState = int(reply.CleanState)
	return
}

// CleanInvalidArcs clean invalid archives.
func (s *Service) CleanInvalidArcs(c context.Context, mid, fid int64) (err error) {
	_, err = s.favClient.CleanInvalidFavs(c, &pb.CleanInvalidFavsReq{
		Typ: int32(model.TypeVideo),
		Mid: mid,
		Fid: fid,
	})
	return
}
