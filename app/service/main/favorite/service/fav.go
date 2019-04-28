package service

import (
	"context"
	"sort"
	"sync"
	"time"

	"go-common/app/service/main/favorite/conf"
	"go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

var (
	_emptyFavRelations = []*model.Favorite{}
	_emptyUsers        = []*model.User{}
	_emptyFolders      = []*model.Folder{}
)

// UserFolders return folders by mid.
func (s *Service) UserFolders(c context.Context, typ int8, mid, vmid, oid int64, otype int8) (res []*model.Folder, err error) {
	res = make([]*model.Folder, 0)
	isSelf := true
	if vmid > 0 && vmid != mid {
		isSelf = false
		mid = vmid
	}
	foldersMap, err := s.userFolders(c, typ, mid)
	if err != nil {
		return
	}
	if len(foldersMap) == 0 {
		if foldersMap, err = s.favDao.UserFolders(c, typ, mid); err != nil {
			log.Error("s.favDao.UserFolders(%d,%d) error(%v)", typ, mid, err)
			return
		}
		// if user has no folder and self == true then create default folder
		var folder *model.Folder
		if len(foldersMap) == 0 && isSelf {
			if folder, err = s.initDefaultFolder(c, typ, mid); err != nil {
				return
			}
			if folder == nil {
				res = _emptyFolders
			} else {
				res = []*model.Folder{folder}
			}
			s.cache.Do(c, func(c context.Context) {
				if e := s.favDao.AddFidsRedis(c, typ, mid, folder); e != nil {
					log.Error("s.favDao.AddFidsRedis(%d,%d,%v) error(%v)", typ, mid, folder, e)
				}
				if e := s.favDao.SetFoldersMc(c, folder); e != nil {
					log.Error("favDao.SetFoldersMc(%v) error(%v)", folder, e)
				}
			})
			return
		}
	}
	// if this resource has been faved
	var faveds map[int64]struct{}
	if oid > 0 {
		isFaved, err1 := s.IsFavored(c, otype, mid, oid)
		if err1 != nil {
			log.Error("s.IsFavored(%d,%d,%d) error(%v)", otype, mid, oid, err)
			return
		}
		if isFaved {
			var fids []int64
			if fids, err = s.favedFids(c, otype, mid, oid); err != nil {
				log.Error("s.favedFids(%d,%d,%d) error(%v)", otype, mid, oid, err)
				return
			}
			faveds = make(map[int64]struct{}, len(fids))
			for _, fid := range fids {
				faveds[fid] = struct{}{}
			}
		}
	}
	for _, f := range foldersMap {
		if !isSelf && !f.IsPublic() {
			continue
		}
	}
	for _, f := range foldersMap {
		if !isSelf && !f.IsPublic() {
			continue
		}
		if _, ok := faveds[f.ID]; ok {
			f.Favored = 1
		}
	}

	// folder sort
	fst, err := s.folderSort(c, typ, mid)
	if err != nil {
		log.Error("s.folderSort(%d,%d) error(%v)", typ, mid, err)
		return
	}
	res, update := fst.SortFolders(foldersMap, isSelf)
	if update {
		s.SetFolderSort(c, typ, mid, fst.Sort)
	}
	return
}

// SetFolderSort set folder sort.
func (s *Service) SetFolderSort(c context.Context, typ int8, mid int64, fids []int64) (err error) {
	var now = time.Now().Unix()
	fst := &model.FolderSort{
		Type:  typ,
		Mid:   mid,
		Sort:  fids,
		CTime: xtime.Time(now),
		MTime: xtime.Time(now),
	}
	count, err := s.favDao.FolderCnt(c, typ, mid)
	if err != nil {
		log.Error("s.favDao.FolderCnt(%d,%d) error(%v)", typ, mid, err)
		return
	}
	if len(fids) != count {
		if err = s.favDao.DelFidsRedis(c, typ, mid); err != nil {
			log.Error("s.favDao.DelFavsCache(%d) error(%v)", mid, err)
		}
		log.Warn("sort not equal user Folders mid:%d sort:%d count:%d", mid, len(fids), count)
		err = ecode.FavFolderSortErr
		return
	}
	folder, err := s.folder(c, typ, mid, fids[0])
	if err != nil {
		log.Error("s.folder(%d, %d) error(%v)", mid, fids[0], err)
		return
	}
	if folder == nil {
		err = ecode.FavFolderNotExist
		return
	}
	if !folder.IsDefault() {
		err = ecode.FavFolderSortErr
		return
	}
	if _, err = s.favDao.UpFolderSort(c, fst); err != nil {
		log.Error("s.favDao.UpFolderSort(%v) error(%v)", fst, err)
		return
	}
	s.cache.Do(c, func(c context.Context) {
		if err := s.favDao.SetFolderSortMc(c, fst); err != nil {
			log.Error("s.favDao.SetFolderSortMc(%v) error(%v)", fst, err)
		}
	})
	return
}

func (s *Service) folderStats(c context.Context, fvmids []*model.ArgFVmid) (map[int64]*model.Folder, error) {
	fsMap, missFvmids, err := s.favDao.FolderStatsMc(c, fvmids)
	if err != nil {
		log.Error("s.favDao.FolderStatsMc(%v) error(%v)", fvmids, err)
		return nil, err
	}
	if len(missFvmids) > 0 {
		stats, err := s.favDao.FolderStats(c, missFvmids)
		if err != nil {
			return nil, err
		}
		for key, value := range stats {
			fsMap[key] = value
		}
		s.cache.Do(c, func(c context.Context) {
			if err := s.favDao.SetFolderStatsMc(c, stats); err != nil {
				log.Error("favDao.SetFoldersMc(%v) error(%v)", stats, err)
			}
		})
	}
	return fsMap, nil
}

// Folders return folders by mid.
func (s *Service) Folders(c context.Context, typ int8, mid int64, fvmids []*model.ArgFVmid) (res []*model.Folder, err error) {
	res = make([]*model.Folder, 0)
	foldersMap, err := s.folders(c, typ, mid, fvmids)
	if err != nil {
		return
	}
	fsMap, err := s.folderStats(c, fvmids)
	if err != nil {
		log.Error("s.favDao.FoldersMc(%d,%v) error(%v)", typ, fvmids, err)
		return
	}

	for _, f := range foldersMap {
		fsMap := fsMap[f.MediaID()]
		if fsMap != nil {
			f.PlayCount = fsMap.PlayCount
			f.ShareCount = fsMap.ShareCount
			f.FavedCount = fsMap.FavedCount
		}
		f.Cover = model.CompleteURL(f.Cover)
		res = append(res, f)
	}
	return
}

func (s *Service) recentOids(c context.Context, typ int8, mid int64, fids []int64) (rctFidsMap map[int64][]int64, err error) {
	rctFidsMap, missFids, err := s.favDao.RecentOidsCache(c, typ, mid, fids)
	if err != nil {
		log.Error("s.favDao.RecentOidsCache(%d,%d,%v) error(%v)", typ, mid, fids, err)
		return nil, nil
	}
	if len(missFids) == 0 {
		return
	}
	g := new(errgroup.Group)
	if len(missFids) > 8 {
		g.GOMAXPROCS(8)
	}
	mux := new(sync.Mutex)
	for _, mf := range missFids {
		fid := mf
		g.Go(func() error {
			oids, err1 := s.favDao.RecentOids(c, mid, fid, typ)
			mux.Lock()
			if err1 != nil {
				rctFidsMap[fid] = make([]int64, 0)
			} else {
				rctFidsMap[fid] = oids
			}
			mux.Unlock()
			return nil
		})
	}
	err = g.Wait()
	return
}

// folderSort return folder sort.
func (s *Service) folderSort(c context.Context, typ int8, mid int64) (fst *model.FolderSort, err error) {
	if fst, err = s.favDao.FolderSortMc(c, typ, mid); err != nil {
		log.Error("s.favDao.FolderSortMc(%d,%d) error(%v)", typ, mid, err)
		err = nil
	}
	if fst != nil {
		return
	}
	if fst, err = s.favDao.FolderSort(c, typ, mid); err != nil {
		log.Error("s.favDao.FolderSort(%d,%d) error(%v)", typ, mid, err)
		return
	}
	if fst == nil {
		fst = &model.FolderSort{Type: typ, Mid: mid}
	}
	s.cache.Do(c, func(c context.Context) {
		if err1 := s.favDao.SetFolderSortMc(c, fst); err1 != nil {
			log.Error("s.favDao.SetFolderSortMc(%v) error(%v)", fst, err1)
		}
	})
	return
}

// folders return folders
func (s *Service) folders(c context.Context, typ int8, mid int64, fvmids []*model.ArgFVmid) (fsMap map[string]*model.Folder, err error) {
	if len(fvmids) == 0 {
		return
	}
	fsMap, missFvmids, err := s.favDao.FoldersMc(c, fvmids)
	if err != nil {
		log.Error("s.favDao.FoldersMc(%d,%v) error(%v)", typ, fvmids, err)
		return
	}
	var dbFoldersMap map[string]*model.Folder
	if len(missFvmids) != 0 {
		if dbFoldersMap, err = s.favDao.Folders(c, missFvmids); err != nil {
			log.Error("s.favDao.Folders(%d,%d,%v) error(%v)", typ, mid, fvmids, err)
			return
		}
		for k, v := range dbFoldersMap {
			fsMap[k] = v
		}
		g := new(errgroup.Group)
		if len(missFvmids) > 8 {
			g.GOMAXPROCS(8)
		}
		for _, folder := range dbFoldersMap {
			if folder != nil {
				f := folder
				f.RecentOids = []int64{}
				f.RecentRes = []*model.Resource{}
				g.Go(func() error {
					recents, err1 := s.favDao.RecentRes(c, f.Mid, f.ID)
					if err1 != nil {
						return err1
					}
					if len(recents) > 0 {
						f.RecentRes = recents
						for _, res := range recents {
							if res.Typ == int32(typ) {
								f.RecentOids = append(f.RecentOids, res.Oid)
							}
						}
					}
					return nil
				})
			}
		}
		err = g.Wait()
		if err != nil {
			return
		}
		var folders []*model.Folder
		for _, v := range dbFoldersMap {
			if v != nil {
				folders = append(folders, v)
			}
		}
		s.cache.Do(c, func(c context.Context) {
			if err := s.favDao.SetFoldersMc(c, folders...); err != nil {
				log.Error("favDao.SetFoldersMc(%v) error(%v)", folders, err)
			}
		})
	}
	return
}

// userfolders return user's folders
func (s *Service) userFolders(c context.Context, typ int8, mid int64) (folders map[int64]*model.Folder, err error) {
	var ok bool
	if ok, err = s.favDao.ExpireFolder(c, typ, mid); err != nil {
		log.Error("s.favDao.ExpireFolder(%d,%d) error(%v)", typ, mid, err)
		return
	}
	folders = make(map[int64]*model.Folder, 2)
	if ok {
		var fids []int64
		if fids, err = s.favDao.FidsRedis(c, typ, mid); err != nil {
			return
		}
		fvmids := make([]*model.ArgFVmid, 0, len(fids))
		for _, fid := range fids {
			fvmids = append(fvmids, &model.ArgFVmid{Fid: fid, Vmid: mid})
		}
		var folderMap map[string]*model.Folder
		if folderMap, err = s.folders(c, typ, mid, fvmids); err != nil {
			return
		}
		for _, f := range folderMap {
			folders[f.ID] = f
		}
		return
	}
	if folders, err = s.favDao.UserFolders(c, typ, mid); err != nil {
		log.Error("s.favDao.UserFolders(%d,%d) error(%v)", typ, mid, err)
		return
	}
	g := new(errgroup.Group)
	if len(folders) > 8 {
		g.GOMAXPROCS(8)
	}
	for _, folder := range folders {
		f := folder
		if f != nil {
			f.RecentOids = []int64{}
			f.RecentRes = []*model.Resource{}
			g.Go(func() error {
				recents, err1 := s.favDao.RecentRes(c, f.Mid, f.ID)
				if err1 != nil {
					return err1
				}
				if len(recents) > 0 {
					f.RecentRes = recents
					for _, res := range recents {
						if res.Typ == int32(typ) {
							f.RecentOids = append(f.RecentOids, res.Oid)
						}
					}
				}
				return nil
			})
		}
	}
	err = g.Wait()
	if err != nil {
		return
	}
	var fs []*model.Folder
	for _, v := range folders {
		f := *v
		fs = append(fs, &f)
	}
	s.cache.Do(c, func(c context.Context) {
		if err := s.favDao.AddFidsRedis(c, typ, mid, fs...); err != nil {
			log.Error("s.favDao.AddFidsRedis(%d,%d,%v) error(%v)", typ, mid, fs, err)
		}
		if err := s.favDao.SetFoldersMc(c, fs...); err != nil {
			log.Error("favDao.SetFoldersMc(%v) error(%v)", fs, err)
		}
	})
	return
}

// AddFolder add a folder.
func (s *Service) AddFolder(c context.Context, typ int8, mid int64, name, description, cover string, public int32, cookie, accessKey string) (fid int64, err error) {
	var hitSensitive bool
	if err = model.CheckType(typ); err != nil {
		log.Error("model.CheckArg(%d) error(%v)", typ, err)
		err = ecode.RequestErr
		return
	}
	if err = s.checkUser(c, mid); err != nil {
		log.Error("s.checkUser(mid:%d) error(%v)", mid, err)
		return
	}
	if err = s.filter(c, name); err != nil {
		if ecode.FavHitSensitive.Equal(err) {
			err = nil
			hitSensitive = true
		} else {
			log.Error("s.filter(name:%s) error(%v)", name, err)
			return
		}
	}
	var count int
	if count, err = s.favDao.FolderCnt(c, typ, mid); err != nil {
		log.Error("s.favDao.FolderCnt(%d,$d) error(%v)", typ, mid, err)
		return

	}
	// this means there is no favorite for user
	if count == 0 && typ != model.TypePlayVideo {
		s.initDefaultFolder(c, typ, mid)
	}
	maxFolders := conf.Conf.Fav.MaxFolders
	maxNameLen := conf.Conf.Fav.MaxNameLen
	maxDescLen := conf.Conf.Fav.MaxDescLen
	if typ == model.TypePlayVideo {
		maxFolders = conf.Conf.Platform.MaxFolders
		maxNameLen = conf.Conf.Platform.MaxNameLen
		maxDescLen = conf.Conf.Platform.MaxDescLen
	}
	if count >= maxFolders {
		log.Warn("The number of folders can not be more than %d", maxFolders)
		err = ecode.FavMaxFolderCount
		return
	}
	if len([]rune(name)) > maxNameLen {
		log.Warn("arg name(%s) it's length more than %d", name, maxNameLen)
		err = ecode.FavNameTooLong
		return
	}
	if len([]rune(description)) > maxDescLen {
		log.Warn("arg description(%s) it's length more than %d", name, maxDescLen)
		err = ecode.FavDescTooLang
		return
	}
	if public > 1 {
		public = 0
	}
	public = model.AttrNormalPublic | public
	if typ == model.TypeVideo {
		public = public | model.AttrBitNeedAudit
		if hitSensitive {
			//public = public | model.AttrBitHitSensitive
		}
	}
	now := time.Now()
	if typ != model.TypeVideo {
		cover = model.CleanURL(cover)
	}
	f := &model.Folder{
		Type:        typ,
		Mid:         mid,
		Name:        name,
		Cover:       cover,
		Description: description,
		Count:       0,
		Attr:        public,
		State:       model.StateNormal,
		CTime:       xtime.Time(now.Unix()),
		MTime:       xtime.Time(now.Unix()),
	}
	if name != "" {
		f.AttrSet(1, model.AttrBitName)
	}
	if description != "" {
		f.AttrSet(1, model.AttrBitDesc)
	}
	if cover != "" {
		f.AttrSet(1, model.AttrBitCover)
	}
	if typ != model.TypeVideo {
		var folder *model.Folder
		if folder, err = s.favDao.FolderByName(c, typ, mid, name); err != nil {
			log.Error("s.favDao.FolderByName(%d,%d,%s) error(%v)", typ, mid, name, err)
			return
		}
		if folder != nil {
			log.Warn("folder name(%s) exist", name)
			err = ecode.FavFolderExist
			return
		}
	}
	if fid, err = s.favDao.AddFolder(c, f); err != nil {
		log.Error("s.favDao.AddFolder(%v) error(%v)", f, err)
		return
	}
	if fid == 0 {
		log.Warn("mysql.Insert(%v) error(%v)", f, err)
		err = ecode.FavFolderExist
		return
	}
	f.ID = fid
	s.cache.Do(c, func(c context.Context) {
		if ok, err := s.favDao.ExpireFolder(c, f.Type, f.Mid); err == nil && ok {
			if err := s.favDao.AddFidsRedis(c, f.Type, f.Mid, f); err != nil {
				log.Error("favDao.AddFidsRedis(%d,%d,%v) error(%v)", f.Type, f.Mid, f, err)
			}
		}
		if err := s.favDao.SetFoldersMc(c, f); err != nil {
			log.Error("favDao.SetFoldersMc(%v) error(%v)", f, err)
		}
	})
	return
}

// SortFavs .
func (s *Service) SortFavs(c context.Context, typ int8, fid, mid int64, sorts []model.SortFav) (err error) {
	var folder *model.Folder
	if folder, err = s.folder(c, typ, mid, fid); err != nil {
		return
	}
	if folder.Count > 1000 {
		return ecode.FavMaxVideoCount
	}
	s.favDao.PubSortFavs(c, typ, mid, fid, sorts)
	return
}

// UpdateFolder update folder info.
func (s *Service) UpdateFolder(c context.Context, typ int8, fid, mid int64, name, description, cover string, private int32, attr *int32, state *int8) (err error) {
	var audit int32 // 0:默认已过审
	var folder *model.Folder
	if attr == nil {
		if err = s.checkUser(c, mid); err != nil {
			return
		}
		if folder, err = s.folder(c, typ, mid, fid); err != nil {
			return
		}
		if err = s.filter(c, name); err != nil {
			if ecode.FavHitSensitive.Equal(err) {
				err = nil
				//folder.AttrSet(1, model.AttrBitSensitive)
			} else {
				return
			}
		}
	} else {
		if folder, err = s.folderAll(c, typ, mid, fid); err != nil {
			return
		}
	}
	if state != nil && folder.IsDefault() && *state == 1 {
		err = ecode.FavCanNotDelDefault
		return
	}
	fid = folder.ID
	if attr == nil && (folder.Cover != cover || folder.Description != description || folder.Name != name) {
		//待审核
		audit = 1
		if folder.Name != name && name != "" {
			folder.AttrSet(1, model.AttrBitName)
		}
		if folder.Description != description && description != "" {
			folder.AttrSet(1, model.AttrBitDesc)
		}
		if folder.Cover != cover && cover != "" {
			folder.AttrSet(1, model.AttrBitCover)
		}
	}
	maxNameLen := conf.Conf.Fav.MaxNameLen
	maxDescLen := conf.Conf.Fav.MaxDescLen
	if typ == model.TypePlayVideo {
		maxNameLen = conf.Conf.Platform.MaxNameLen
		maxDescLen = conf.Conf.Platform.MaxDescLen
	}
	if len([]rune(name)) > maxNameLen {
		log.Warn("arg name(%s) it's length more than %d", name, maxNameLen)
		err = ecode.FavNameTooLong
		return
	}
	if len([]rune(description)) > maxDescLen {
		log.Warn("arg description(%s) it's length more than %d", name, maxDescLen)
		err = ecode.FavDescTooLang
		return
	}
	if attr == nil {
		if private > 1 {
			private = 0
		}
		folder.AttrSet(private, model.AttrBitPublic)
		folder.AttrSet(audit, model.AttrBitAudit)
	} else {
		folder.Attr = *attr
	}
	if state != nil {
		folder.State = *state
	}
	if typ != model.TypeVideo {
		cover = model.CleanURL(cover)
	}
	f := &model.Folder{
		ID:          fid,
		Type:        typ,
		Mid:         mid,
		Name:        name,
		Cover:       cover,
		Description: description,
		Attr:        folder.Attr,
		MTime:       xtime.Time(time.Now().Unix()),
		State:       folder.State,
	}
	var rows int64
	if rows, err = s.favDao.UpdateFolder(c, f); err != nil {
		log.Error("s.favDao.UpdateFolder(%v) error(%v)", f, err)
		return
	}
	if rows == 0 {
		log.Warn("mysql.Insert(%v) error(%v)", f, err)
		err = ecode.FavFolderExist
		return
	}
	folder.Name = f.Name
	folder.Description = f.Description
	folder.Cover = f.Cover
	folder.Attr = f.Attr
	folder.MTime = f.MTime
	if state != nil {
		// delete
		if *state == 1 {
			s.favDao.PubDelFolder(c, typ, mid, fid, folder.Attr, int64(folder.MTime))
			s.cache.Do(c, func(c context.Context) {
				if ok, err := s.favDao.ExpireFolder(c, folder.Type, folder.Mid); err == nil && ok {
					if err := s.favDao.RemFidsRedis(c, folder.Type, folder.Mid, folder); err != nil {
						log.Error("s.favDao.RemFidsRedis(%d,%d,%v) error(%v)", folder.Type, folder.Mid, folder, err)
					}
				}
				if err := s.favDao.DelFolderMc(c, folder.Type, folder.Mid, folder.ID); err != nil {
					log.Error("s.favDao.DelFolderMc(%d,%d,%d) error(%v)", folder.Type, folder.Mid, folder.ID, err)
				}
			})
		} else {
			s.cache.Do(c, func(c context.Context) {
				if ok, err := s.favDao.ExpireFolder(c, folder.Type, folder.Mid); err == nil && ok {
					if err := s.favDao.AddFidsRedis(c, folder.Type, folder.Mid, folder); err != nil {
						log.Error("favDao.AddFidsRedis(%d,%d,%v) error(%v)", folder.Type, folder.Mid, folder, err)
					}
				}
				if err := s.favDao.SetFoldersMc(c, folder); err != nil {
					log.Error("favDao.SetFoldersMc(%v) error(%v)", folder, err)
				}
			})
		}
	} else {
		s.cache.Do(c, func(c context.Context) {
			if err := s.favDao.SetFoldersMc(c, folder); err != nil {
				log.Error("s.favDao.SetFoldersMc(%v) error(%v)", folder, err)
			}
		})
	}
	return
}

// folder return a folder by fid.
func (s *Service) folderAll(c context.Context, typ int8, mid, fid int64) (folder *model.Folder, err error) {
	if typ == model.TypePlayVideo && fid == 0 {
		err = ecode.FavFolderNotExist
		return
	}
	// if fid query didn't set then return default favorite
	if fid == 0 {
		return s.defaultFolder(c, typ, mid)
	}
	if folder, err = s.favDao.FolderMc(c, typ, mid, fid); err != nil {
		log.Error("s.favDao.FolderMc(%d,%d,%d) error(%v)", typ, mid, fid, err)
		err = nil
	}
	if folder == nil {
		if folder, err = s.favDao.Folder(c, typ, mid, fid); err != nil {
			log.Error("favDao.Folder(%d,%d) error(%v) or folder is nil", mid, fid, err)
			return
		}
		if folder != nil {
			var recent []*model.Resource
			if recent, err = s.favDao.RecentRes(c, mid, fid); err != nil {
				log.Error(" s.favDao.RecentRes(%d,%d) error(%v) or folder is nil", mid, fid, err)
				return
			}
			folder.RecentOids = []int64{}
			folder.RecentRes = []*model.Resource{}
			if len(recent) > 0 {
				folder.RecentRes = recent
				for _, res := range recent {
					if res.Typ == int32(typ) {
						folder.RecentOids = append(folder.RecentOids, res.Oid)
					}
				}
			}
			s.cache.Do(c, func(c context.Context) {
				if e := s.favDao.SetFoldersMc(c, folder); e != nil {
					log.Error("favDao.SetFoldersMc(%v) error(%v)", folder, e)
				}
			})
		}
	}
	if folder == nil {
		err = ecode.FavFolderNotExist
	}
	return
}

// folder return a folder by fid.
func (s *Service) folder(c context.Context, typ int8, mid, fid int64) (folder *model.Folder, err error) {
	if typ == model.TypePlayVideo && fid == 0 {
		err = ecode.FavFolderNotExist
		return
	}
	// if fid query didn't set then return default favorite
	if fid == 0 {
		return s.defaultFolder(c, typ, mid)
	}
	if folder, err = s.favDao.FolderMc(c, typ, mid, fid); err != nil {
		log.Error("s.favDao.FolderMc(%d,%d,%d) error(%v)", typ, mid, fid, err)
		err = nil
	}
	if folder == nil {
		if folder, err = s.favDao.Folder(c, typ, mid, fid); err != nil {
			log.Error("favDao.Folder(%d,%d) error(%v) or folder is nil", mid, fid, err)
			return
		}
		if folder != nil {
			var recent []*model.Resource
			if recent, err = s.favDao.RecentRes(c, mid, fid); err != nil {
				log.Error(" s.favDao.RecentRes(%d,%d) error(%v) or folder is nil", mid, fid, err)
				return
			}
			folder.RecentOids = []int64{}
			folder.RecentRes = []*model.Resource{}
			if len(recent) > 0 {
				folder.RecentRes = recent
				for _, res := range recent {
					if res.Typ == int32(typ) {
						folder.RecentOids = append(folder.RecentOids, res.Oid)
					}
				}
			}
			s.cache.Do(c, func(c context.Context) {
				if e := s.favDao.SetFoldersMc(c, folder); e != nil {
					log.Error("favDao.SetFoldersMc(%v) error(%v)", folder, e)
				}
			})
		}
	}
	if folder == nil || folder.State != model.StateNormal {
		err = ecode.FavFolderNotExist
	}
	return
}

// UserFolder return a folder.
func (s *Service) UserFolder(c context.Context, typ int8, mid, vmid, fid int64) (*model.Folder, error) {
	f, err := s.Folder(c, typ, mid, vmid, fid)
	if err != nil {
		return nil, err
	}
	if (vmid > 0 && mid != vmid) && f.AttrVal(model.AttrBitPublic) != model.AttrIsPublic {
		return nil, ecode.FavFolderNoPublic
	}
	return f, err
}

// Folder return a folder by fid, include deleted folder.
func (s *Service) Folder(c context.Context, typ int8, mid, uid, fid int64) (folder *model.Folder, err error) {
	if uid > 0 && mid != uid {
		mid = uid
	}
	if folder, err = s.folder(c, typ, mid, fid); err != nil {
		return
	}
	folder.Cover = model.CompleteURL(folder.Cover)
	return
}

// defaultFolder return default folder of user.
func (s *Service) defaultFolder(c context.Context, typ int8, mid int64) (folder *model.Folder, err error) {
	fsMap, err := s.userFolders(c, typ, mid)
	if err != nil {
		return
	}
	for _, folder = range fsMap {
		if folder.IsDefault() {
			return
		}
	}
	// means that the fid is not existed and need init a default folder
	return s.initDefaultFolder(c, typ, mid)

}

// initDefaultFolder init user's default folder.
func (s *Service) initDefaultFolder(c context.Context, typ int8, mid int64) (folder *model.Folder, err error) {
	if typ == model.TypePlayVideo {
		return nil, nil
	}
	if err = model.CheckType(typ); err != nil {
		log.Error("model.CheckArg(%d) error(%v)", typ, err)
		err = ecode.RequestErr
		return
	}
	now := time.Now()
	folder = &model.Folder{
		Type:  typ,
		Mid:   mid,
		Name:  model.InitFolderName,
		Attr:  model.AttrDefaultPublic,
		State: model.StateNormal,
		CTime: xtime.Time(now.Unix()),
		MTime: xtime.Time(now.Unix()),
	}
	if folder.ID, err = s.favDao.AddFolder(c, folder); err != nil {
		return
	}
	if folder.ID == 0 {
		log.Warn("initDefaultFolder failed mid:%d", mid)
		err = ecode.FavFolderNotExist
		return
	}
	// init a default folder need cache
	s.cache.Do(c, func(c context.Context) {
		if ok, err := s.favDao.ExpireFolder(c, folder.Type, folder.Mid); err == nil && ok {
			if err := s.favDao.AddFidsRedis(c, typ, mid, folder); err != nil {
				log.Error("s.favDao.AddFidsRedis(%d,%d,%v) error(%v)", typ, mid, folder, err)
			}
		}
	})
	return
}

// FavoritesAll return favorieds info by fid.
func (s *Service) FavoritesAll(c context.Context, tp int8, mid, uid, fid int64, tid, tv, pn, ps int, keyword, order string) (res *model.Favorites, err error) {
	if keyword != "" || tid != 0 || tv != 0 || (order != "" && order != model.SortMtime) {
		return s.Relations(c, tp, mid, uid, fid, tid, tv, pn, ps, keyword, order)
	}
	var (
		ok     bool
		folder *model.Folder
		favs   []*model.Favorite
		start  = (pn - 1) * ps
		end    = start + ps - 1
		oldMid = mid
	)
	res = &model.Favorites{}
	res.Page.Num = pn
	res.Page.Size = ps
	if uid > 0 {
		mid = uid
	}
	if folder, err = s.folder(c, tp, mid, fid); err != nil {
		log.Error("s.folder(%d,%d,%d) error(%v)", tp, mid, fid, err)
		return
	}
	if !folder.Access(oldMid) {
		err = ecode.AccessDenied
		return
	}
	fid = folder.ID // NOTE if init a default folder fid=0
	if ok, err = s.favDao.ExpireAllRelations(c, mid, fid); err != nil {
		log.Error("s.favDao.ExpireRelations(%d,%d) error(%v)", fid, err)
		return
	}
	if ok {
		if res.Page.Count, err = s.favDao.CntAllRelationsCache(c, mid, fid); err != nil {
			log.Error("s.favDao.CntRelationsCache(%d,%d) error(%v)", mid, fid, err)
			return
		}
		if favs, err = s.favDao.FolderAllRelationsCache(c, tp, mid, fid, start, end); err != nil {
			log.Error("s.favDao.FolderRelationsCache(%d,%d,%d,%d,%d) error(%v)", tp, mid, fid, start, end, err)
			return
		}
		if res.Page.Count != folder.Count {
			log.Warn("del_dirt_cache %d,%d,%d,%d", mid, fid, res.Page.Count, folder.Count)
			s.favDao.DelAllRelationsCache(c, mid, fid)
		}
	} else {
		if res.Page.Count, err = s.favDao.CntAllRelations(c, mid, fid); err != nil {
			log.Error("s.favDao.CntRelations(%d,%d,%d) error(%v)", fid, mid, err)
			return
		}
		if res.Page.Count > 0 {
			if favs, err = s.favDao.FolderAllRelations(c, mid, fid, start, ps); err != nil {
				log.Error("s.favDao.FolderRelations(%d,%d,%d,%d,%d) error(%v)", tp, mid, fid, start, end, err)
				return
			}
		}
		s.favDao.PubInitAllFolderRelations(c, tp, mid, fid)
	}
	if res.Page.Count == 0 || len(favs) == 0 {
		res.List = _emptyFavRelations
	} else {
		res.List = favs
	}
	return
}

// Favorites return favorieds info by fid.
func (s *Service) Favorites(c context.Context, tp int8, mid, uid, fid int64, tid, tv, pn, ps int, keyword, order string) (res *model.Favorites, err error) {
	if keyword != "" || tid != 0 || tv != 0 || (order != "" && order != model.SortMtime) {
		return s.Relations(c, tp, mid, uid, fid, tid, tv, pn, ps, keyword, order)
	}
	var (
		ok     bool
		folder *model.Folder
		favs   []*model.Favorite
		start  = (pn - 1) * ps
		end    = start + ps - 1
		oldMid = mid
	)
	res = &model.Favorites{}
	res.Page.Num = pn
	res.Page.Size = ps
	if uid > 0 {
		mid = uid
	}
	if folder, err = s.folder(c, tp, mid, fid); err != nil {
		log.Error("s.folder(%d,%d,%d) error(%v)", tp, mid, fid, err)
		return
	}
	if !folder.Access(oldMid) {
		err = ecode.AccessDenied
		return
	}
	fid = folder.ID // NOTE if init a default folder fid=0
	if ok, err = s.favDao.ExpireRelations(c, mid, fid); err != nil {
		log.Error("s.favDao.ExpireRelations(%d,%d) error(%v)", fid, err)
		return
	}
	if ok {
		if res.Page.Count, err = s.favDao.CntRelationsCache(c, mid, fid); err != nil {
			log.Error("s.favDao.CntRelationsCache(%d,%d) error(%v)", mid, fid, err)
			return
		}
		if favs, err = s.favDao.FolderRelationsCache(c, tp, mid, fid, start, end); err != nil {
			log.Error("s.favDao.FolderRelationsCache(%d,%d,%d,%d,%d) error(%v)", tp, mid, fid, start, end, err)
			return
		}
	} else {
		if res.Page.Count, err = s.favDao.CntRelations(c, mid, fid, tp); err != nil {
			log.Error("s.favDao.CntRelations(%d,%d,%d) error(%v)", fid, mid, err)
			return
		}
		if res.Page.Count > 0 {
			if favs, err = s.favDao.FolderRelations(c, tp, mid, fid, start, ps); err != nil {
				log.Error("s.favDao.FolderRelations(%d,%d,%d,%d,%d) error(%v)", tp, mid, fid, start, end, err)
				return
			}
		}
		s.favDao.PubInitFolderRelations(c, tp, mid, fid)
	}
	if res.Page.Count == 0 || len(favs) == 0 {
		res.List = _emptyFavRelations
	} else {
		res.List = favs
	}
	return
}

// BatchFavs return 1000 oids by mid.
func (s *Service) BatchFavs(c context.Context, typ int8, mid int64, limit int) (oids []int64, err error) {
	oids = make([]int64, 0, limit)
	un, err := s.favDao.FavedBit(c, typ, mid)
	// un mean user do not had any fav
	if err != nil || un {
		return
	}
	if oids, err = s.favDao.BatchOidsMc(c, typ, mid); err != nil {
		log.Error("s.favDao.BatchOidsMc(%d,%d) error(%v)", typ, mid, err)
		return
	}
	if len(oids) > 0 {
		return
	}
	if oids, err = s.favDao.BatchOidsRedis(c, typ, mid, limit); err != nil {
		log.Error("s.favDao.BatchOidsRedis(%d,%d,%d) error(%v)", typ, mid, limit, err)
		return
	}
	if len(oids) > 0 {
		return
	}
	if oids, err = s.favDao.BatchOids(c, typ, mid, limit); err != nil {
		log.Error("s.favDao.BatchOids(%d,%d,%d) error(%v)", typ, mid, limit, err)
		return
	}
	if len(oids) == 0 {
		if err1 := s.favDao.SetUnFavedBit(c, typ, mid); err1 != nil {
			log.Error("s.favDao.SetUnFavedBit(%d,%d,%v) error(%v)", typ, mid, err1)
			return
		}
	}
	s.cache.Do(c, func(c context.Context) {
		if err1 := s.favDao.SetBatchOidsMc(c, typ, mid, oids); err1 != nil {
			log.Error("s.favDao.SetBatchOidsMc(%d,%d,%v) error(%v)", typ, mid, oids, err1)
		}
	})
	return
}

// AddFav add a favorite into folder.
func (s *Service) AddFav(c context.Context, tp int8, mid, fid, oid int64, otype int8) (err error) {
	var (
		faved  bool
		folder *model.Folder
	)
	if err = model.CheckArg(tp, oid); err != nil {
		log.Error("model.CheckArg(%d,%d) error(%v)", tp, oid, err)
		err = ecode.RequestErr
		return
	}
	if otype == 0 {
		otype = tp
	}
	if folder, err = s.folder(c, tp, mid, fid); err != nil {
		return
	}
	if folder.IsLimited(1, s.conf.Fav.DefaultFolderLimit, s.conf.Fav.NormalFolderLimit) {
		err = ecode.FavResourceOverflow
		return
	}
	fid = folder.ID // NOTE if init a default folder fid=0
	if faved, err = s.isFavedByFid(c, otype, mid, oid, fid); err != nil {
		log.Error("s.isFavedByFid(%d,%d,%d,%d) error(%v)", otype, mid, oid, fid, err)
		return
	}
	if faved {
		err = ecode.FavResourceExist
		return
	}
	now := time.Now()
	s.favDao.PubAddFav(c, tp, mid, fid, oid, folder.Attr, now.Unix(), otype)
	return
}

// DelFav delete a favorite.
func (s *Service) DelFav(c context.Context, tp int8, mid, fid, oid int64, otype int8) (err error) {
	var (
		faved  bool
		folder *model.Folder
	)
	if err = model.CheckArg(tp, oid); err != nil {
		log.Error("model.CheckArg(%d,%d) error(%v)", tp, oid, err)
		err = ecode.RequestErr
		return
	}
	if otype == 0 {
		otype = tp
	}
	if folder, err = s.folder(c, tp, mid, fid); err != nil {
		return
	}
	fid = folder.ID
	if faved, err = s.isFavedByFid(c, otype, mid, oid, fid); err != nil {
		log.Error("s.isFavedByFid(%d,%d,%d,%d) error(%v)", tp, mid, oid, fid, err)
		return
	}
	if !faved {
		err = ecode.FavResourceAlreadyDel
		return
	}
	var (
		rf    *model.Favorite
		ftime = time.Now().Unix()
	)
	if rf, err = s.favDao.Relation(c, otype, mid, fid, oid); err != nil {
		return
	}
	if rf == nil {
		log.Warn("delFav(%d,%d,%d,%d) not found", tp, mid, fid, oid)
		err = ecode.FavResourceAlreadyDel
		return
	}
	s.favDao.PubDelFav(c, tp, mid, fid, oid, folder.Attr, ftime, otype)
	return
}

// MultiAddFavs add multi favorite.
func (s *Service) MultiAddFavs(c context.Context, typ int8, mid, fid int64, oids []int64) (err error) {
	if len(oids) == 0 {
		log.Error("MultiAddFavs error(oids len equal zero)", typ, err)
		err = ecode.RequestErr
		return
	}
	if len(oids) == 1 {
		return s.AddFav(c, typ, mid, fid, oids[0], typ)
	}
	var (
		folder *model.Folder
	)
	if err = model.CheckType(typ); err != nil {
		log.Error("model.CheckType(%d) error(%v)", typ, err)
		err = ecode.RequestErr
		return
	}
	if folder, err = s.folder(c, typ, mid, fid); err != nil {
		return
	}
	if folder.IsLimited(len(oids), s.conf.Fav.DefaultFolderLimit, s.conf.Fav.NormalFolderLimit) {
		err = ecode.FavResourceOverflow
		return
	}
	fid = folder.ID
	rows, err := s.favDao.MultiAddRelations(c, typ, mid, fid, oids)
	if err != nil {
		log.Error("s.favDao.MultiAddRelations(%d,%d,%d,%v) error(%v)", typ, mid, fid, oids, err)
		return
	}
	if rows < 1 {
		log.Warn("add relations type(%d) mid(%d) fid(%d) oids(%v) have no del", typ, mid, fid, oids)
		err = ecode.FavResourceExist
		return
	}
	s.cache.Do(c, func(c context.Context) {
		if err := s.favDao.DelRelationFidsMc(c, typ, mid, oids...); err != nil {
			log.Error("s.favDao.DelRelationFidsMc(%d,%d,%v) error(%v)", typ, mid, oids, err)
		}
		if err := s.favDao.DelRecentOidsMc(c, typ, mid); err != nil {
			log.Error("s.favDao.DelRecentOidsMc(%d,%d) error(%v)", typ, mid, err)
		}
		if err := s.favDao.DelRecentResMc(c, model.TypeVideo, mid); err != nil {
			log.Error("s.favDao.DelRecentResMc(%d,%d) error(%v)", typ, mid, err)
		}
		if err := s.favDao.DelRelationOidsCache(c, typ, mid); err != nil {
			log.Error("s.favDao.DelRelationOidsCache(%d,%d) error(%v)", typ, mid, err)
		}
		if err := s.favDao.DelRelationsCache(c, mid, fid); err != nil {
			log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, fid, err)
		}
		if err := s.favDao.DelAllRelationsCache(c, mid, fid); err != nil {
			log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, fid, err)
		}
	})
	s.favDao.PubMultiAddFavs(c, typ, mid, fid, rows, folder.Attr, oids, time.Now().Unix())
	return
}

// MultiDelFavs delete multi favorite.
func (s *Service) MultiDelFavs(c context.Context, typ int8, mid, fid int64, oids []int64) (err error) {
	if len(oids) == 1 {
		return s.DelFav(c, typ, mid, fid, oids[0], typ)
	}
	var f *model.Folder
	if f, err = s.folder(c, typ, mid, fid); err != nil {
		return
	}
	fid = f.ID
	rows, err := s.favDao.MultiDelRelations(c, typ, mid, fid, oids)
	if err != nil {
		log.Error("s.favDao.MultiDelRelations(%d,%d,%d,%v) error(%v)", typ, mid, fid, oids, err)
		return
	}
	if rows < 1 {
		log.Warn("del relations type(%d) mid(%d) fid(%d) oids(%v) have no del", typ, mid, fid, oids)
		err = ecode.FavResourceAlreadyDel
		return
	}
	s.cache.Do(c, func(c context.Context) {
		if err := s.favDao.DelRelationFidsMc(c, typ, mid, oids...); err != nil {
			log.Error("s.favDao.DelRelationFidsMc(%d,%d,%v) error(%v)", typ, mid, oids, err)
		}
		if err := s.favDao.DelRecentOidsMc(c, typ, mid); err != nil {
			log.Error("s.favDao.DelRecentOidsMc(%d,%d) error(%v)", typ, mid, err)
		}
		if err := s.favDao.DelRecentResMc(c, model.TypeVideo, mid); err != nil {
			log.Error("s.favDao.DelRecentResMc(%d,%d) error(%v)", typ, mid, err)
		}
		if err := s.favDao.DelRelationOidsCache(c, typ, mid); err != nil {
			log.Error("s.favDao.DelRelationOidsCache(%d,%d) error(%v)", typ, mid, err)
		}
		if err := s.favDao.DelRelationsCache(c, mid, fid); err != nil {
			log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, fid, err)
		}
		if err := s.favDao.DelAllRelationsCache(c, mid, fid); err != nil {
			log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, fid, err)
		}
	})
	s.favDao.PubMultiDelFavs(c, typ, mid, fid, rows, f.Attr, oids, time.Now().Unix())
	return
}

// IsFavored check if oid faved by user
func (s *Service) IsFavored(c context.Context, typ int8, mid, oid int64) (faved bool, err error) {
	un, err := s.favDao.FavedBit(c, typ, mid)
	// un mean user do not had any fav
	if err != nil || un {
		return
	}
	ok, err := s.favDao.ExpireRelationOids(c, typ, mid)
	if err != nil {
		log.Error("s.favDao.ExpireRelationOids(%d,%d) error(%v)", typ, mid, err)
		return
	}
	if ok {
		if faved, err = s.favDao.IsFavedCache(c, typ, mid, oid); err != nil {
			log.Error("s.favDao.IsFavedCache(%d,%d,%d) error(%v)", typ, mid, oid, err)
			return
		}
		return
	}
	fids, err := s.favDao.RelationFidsByOid(c, typ, mid, oid)
	if err != nil {
		log.Error("s.favDao.RelationFidsByOid(%d,%d,%d) error(%v)", typ, mid, oid, err)
		return
	}
	if len(fids) != 0 {
		faved = true
	}
	s.favDao.PubInitRelationFids(c, typ, mid)
	return
}

// IsFavoreds check if oids faved by user
func (s *Service) IsFavoreds(c context.Context, typ int8, mid int64, oids []int64) (favedMap map[int64]bool, err error) {
	favedMap = make(map[int64]bool, len(oids))
	un, err := s.favDao.FavedBit(c, typ, mid)
	// un mean user do not had any fav
	if err != nil || un {
		return
	}
	ok, err := s.favDao.ExpireRelationOids(c, typ, mid)
	if err != nil {
		log.Error("s.favDao.ExpireRelationOids(%d,%d) error(%v)", typ, mid, err)
		return
	}
	if ok {
		if favedMap, err = s.favDao.IsFavedsCache(c, typ, mid, oids); err != nil {
			log.Error("s.favDao.IsFavedsCache(%d,%d,%v) error(%v)", typ, mid, oids, err)
			return
		}
		return
	}
	fidsMap, err := s.favDao.RelationFidsByOids(c, typ, mid, oids)
	if err != nil {
		log.Error("s.favDao.RelationFidsByOids(%d,%d,%v) error(%v)", typ, mid, oids, err)
		return
	}
	for oid, fids := range fidsMap {
		if len(fids) != 0 {
			favedMap[oid] = true
		}
	}
	s.favDao.PubInitRelationFids(c, typ, mid)
	return
}

// IsFavedByFid check if oid faved by the fid
func (s *Service) IsFavedByFid(c context.Context, tp int8, mid, oid, fid int64) (faved bool, err error) {
	un, err := s.favDao.FavedBit(c, tp, mid)
	// un mean user do not had any fav
	if err != nil || un {
		return
	}
	fids, err := s.favedFids(c, tp, mid, oid)
	if err != nil {
		log.Error("s.favedFids(%d,%d,%d) error(%v)", tp, mid, oid, err)
		return
	}
	if len(fids) == 0 {
		return
	}
	for _, id := range fids {
		if id == fid {
			faved = true
		}
	}
	return
}

// isFavedByFid
func (s *Service) isFavedByFid(c context.Context, tp int8, mid, oid, fid int64) (faved bool, err error) {
	fids, err := s.favedFids(c, tp, mid, oid)
	if err != nil {
		log.Error("s.favedFids(%d,%d,%d) error(%v)", tp, mid, oid, err)
		return
	}
	if len(fids) == 0 {
		return false, nil
	}
	for _, ofid := range fids {
		if ofid == fid {
			return true, nil
		}
	}
	return
}

// favedFids determine oid whether or not collected by user.
func (s *Service) favedFids(c context.Context, typ int8, mid, oid int64) (fids []int64, err error) {
	if fids, err = s.favDao.RelaitonFidsMc(c, typ, mid, oid); err != nil {
		log.Error("s.favDao.RelaitonFidsMc(%d,%d,%d) error(%v)", typ, mid, oid, err)
	}
	if len(fids) == 0 {
		if fids, err = s.favDao.RelationFidsByOid(c, typ, mid, oid); err != nil {
			log.Error("favDao.RelationFidsByOid(%d,%d,%d) error(%v)", typ, mid, oid, err)
			return
		}
		if len(fids) != 0 {
			s.cache.Do(c, func(c context.Context) {
				if err1 := s.favDao.SetRelaitonFidsMc(c, typ, mid, oid, fids); err1 != nil {
					log.Error("favDao.SetRelaitonFidsMc(%d,%d,%d,%v) error(%v)", typ, mid, oid, fids, err1)
				}
			})
		}
	}
	return
}

// DelFolder delete Folder and push databus msg to del favs in folder.
func (s *Service) DelFolder(c context.Context, typ int8, mid, fid int64) (err error) {
	var f *model.Folder
	if f, err = s.folder(c, typ, mid, fid); err != nil {
		return
	}
	fid = f.ID
	if f.IsDefault() {
		err = ecode.FavCanNotDelDefault
		return
	}
	// delete folder from database
	rows, err := s.favDao.DelFolder(c, typ, mid, fid)
	if err != nil {
		log.Error("s.favDao.DelFolder(%d,%d,%d) error(%v)", typ, mid, fid, err)
		return
	}
	if rows < 1 {
		log.Warn("Del folder mid(%d) fid(%d) have delete failed", mid, fid)
		err = ecode.FavFloderAlreadyDel
		return
	}
	// delete resource async by databus
	now := time.Now().Unix()
	s.favDao.PubDelFolder(c, typ, mid, fid, f.Attr, now)
	f.MTime = xtime.Time(now)
	s.cache.Do(c, func(c context.Context) {
		if ok, err := s.favDao.ExpireFolder(c, f.Type, f.Mid); err == nil && ok {
			if err := s.favDao.RemFidsRedis(c, f.Type, f.Mid, f); err != nil {
				log.Error("s.favDao.RemFidsRedis(%d,%d,%v) error(%v)", f.Type, f.Mid, f, err)
			}
		}
		if err := s.favDao.DelFolderMc(c, f.Type, f.Mid, f.ID); err != nil {
			log.Error("s.favDao.DelFolderMc(%d,%d,%d) error(%v)", f.Type, f.Mid, f.ID, err)
		}
	})
	return
}

// CntUserFolders count user's folders.
func (s *Service) CntUserFolders(c context.Context, typ int8, mid, vimd int64) (count int, err error) {
	if vimd > 0 && mid != vimd {
		mid = vimd
	}
	fsMap, err := s.userFolders(c, typ, mid)
	if err != nil {
		return
	}
	if len(fsMap) == 0 {
		return
	}
	for _, folder := range fsMap {
		if folder.State == model.StateNormal {
			if (vimd > 0 && mid != vimd) && folder.AttrVal(model.AttrBitPublic) != model.AttrIsPublic {
				continue
			}
			count++
		}
	}
	return
}

// UpFolderName update user's folder name.
func (s *Service) UpFolderName(c context.Context, typ int8, mid, fid int64, name, cookie, accessKey string) (err error) {
	if err = s.checkRealname(c, mid); err != nil {
		return
	}
	if err = s.filter(c, name); err != nil {
		return
	}
	f, err := s.folder(c, typ, mid, fid)
	if err != nil {
		return
	}

	fid = f.ID
	rows, err := s.favDao.UpFolderName(c, typ, mid, fid, name)
	if err != nil {
		log.Error("s.favDao.UpFolderName(%d, %d,%d, %s) error(%v)", typ, mid, fid, name, err)
		return
	}
	if rows < 1 {
		err = nil
		return
	}
	f.AttrSet(1, model.AttrBitAudit)
	f.AttrSet(1, model.AttrBitName)
	if _, err = s.favDao.UpFolderAttr(c, typ, mid, fid, f.Attr); err != nil {
		log.Error("s.favDao.UpFolderAttr(%d, %d,%d, %d) error(%v)", typ, mid, fid, f.Attr, err)
		return
	}
	f.Name = name
	s.cache.Do(c, func(c context.Context) {
		if err := s.favDao.SetFoldersMc(c, f); err != nil {
			log.Error("s.favDao.SetFoldersMc(%d,%d,%v) error(%v)", f.ID, f.Mid, err)
		}
	})
	return
}

// UpFolderAttr update folder's attr.
func (s *Service) UpFolderAttr(c context.Context, typ int8, mid, fid int64, public int32) (err error) {
	f, err := s.folder(c, typ, mid, fid)
	if err != nil {
		return
	}
	if f == nil {
		err = ecode.FavFolderNotExist
		return
	}
	if public > 1 {
		public = 0
	}
	if public == f.AttrVal(model.AttrBitPublic) {
		return
	}
	f.AttrSet(public, model.AttrBitPublic)
	fid = f.ID
	if _, err = s.favDao.UpFolderAttr(c, typ, mid, fid, f.Attr); err != nil {
		log.Error("s.favDao.UpFolderAttr(%d,%d,%d,%d) error(%v)", typ, mid, fid, f.Attr, err)
		return
	}
	if f.IsDefault() && f.Name == "__default__" {
		f.Name = model.DefaultFolderName
	}
	s.cache.Do(c, func(c context.Context) {
		if err := s.favDao.SetFoldersMc(c, f); err != nil {
			log.Error("s.favDao.SetFoldersMc(%d,%d,%d,%d) error(%v)", f.Type, f.Mid, f.ID, f.Attr, err)
		}
	})
	return
}

// MoveFavs move resources from old folder to new folder
func (s *Service) MoveFavs(c context.Context, typ int8, mid, oldfid, newfid int64, oids []int64) (err error) {
	var (
		nf, of        *model.Folder
		rows, rowsDel int64
	)
	if of, err = s.folder(c, typ, mid, oldfid); err != nil {
		return
	}
	oldfid = of.ID
	if nf, err = s.folder(c, typ, mid, newfid); err != nil {
		return
	}
	newfid = nf.ID
	if newfid == oldfid {
		// means one fid = 0 the other is default favorite
		return
	}
	if nf.IsLimited(len(oids), s.conf.Fav.DefaultFolderLimit, s.conf.Fav.NormalFolderLimit) {
		err = ecode.FavResourceOverflow
		return
	}
	tx, err := s.favDao.BeginTran(c)
	if err != nil {
		log.Error("s.favDao.BeginTran() error(%v)", err)
		return
	}
	if rows, err = s.favDao.TxCopyRelations(tx, typ, mid, mid, oldfid, newfid, oids); err != nil {
		log.Error("s.favDao.MoveRelations(%d,%d,%d,%d,%v) error(%v)", typ, mid, oldfid, newfid, oids, err)
		tx.Rollback()
		return
	}
	if rowsDel, err = s.favDao.TxMultiDelRelations(tx, typ, mid, oldfid, oids); err != nil {
		log.Error("s.favDao.TxMultiDelRelations(%d,%d,%d,%v) error(%v)", typ, mid, oldfid, oids, err)
		tx.Rollback()
		return
	}
	if int(rows) < len(oids) || int(rowsDel) < len(oids) {
		log.Warn("oids not equal rows or rowsDel oids:%d,rows:%d,mid:%d error(%v)", len(oids), rows, rowsDel, mid, err)
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	if rowsDel > 0 {
		s.cache.Do(c, func(c context.Context) {
			if err := s.favDao.DelRelationFidsMc(c, typ, mid, oids...); err != nil {
				log.Error("s.favDao.DelRelationFidsMc(%d,%d,%v) error(%v)", typ, mid, oids, err)
			}
			if err := s.favDao.DelRecentOidsMc(c, typ, mid); err != nil {
				log.Error("s.favDao.DelRecentOidsMc(%d,%d) error(%v)", typ, mid, err)
			}
			if err := s.favDao.DelRecentResMc(c, model.TypeVideo, mid); err != nil {
				log.Error("s.favDao.DelRecentResMc(%d,%d) error(%v)", typ, mid, err)
			}
			if err := s.favDao.DelRelationOidsCache(c, typ, mid); err != nil {
				log.Error("s.favDao.DelRelationOidsCache(%d,%d) error(%v)", typ, mid, err)
			}
			if err := s.favDao.DelRelationsCache(c, mid, oldfid); err != nil {
				log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, oldfid, err)
			}
			if err := s.favDao.DelRelationsCache(c, mid, newfid); err != nil {
				log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, newfid, err)
			}
			if err := s.favDao.DelAllRelationsCache(c, mid, oldfid); err != nil {
				log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, oldfid, err)
			}
			if err := s.favDao.DelAllRelationsCache(c, mid, newfid); err != nil {
				log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, newfid, err)
			}
		})
	}
	s.favDao.PubMoveFavs(c, typ, mid, oldfid, newfid, rows, oids, time.Now().Unix())
	return
}

// CopyFavs copy resources from old folder to new folder
func (s *Service) CopyFavs(c context.Context, typ int8, oldmid, mid, oldfid, newfid int64, oids []int64) (err error) {
	var nf, of *model.Folder
	if oldmid == 0 {
		oldmid = mid
	}
	if of, err = s.folder(c, typ, oldmid, oldfid); err != nil {
		return
	}
	oldfid = of.ID
	if nf, err = s.folder(c, typ, mid, newfid); err != nil {
		return
	}
	newfid = nf.ID
	if newfid == oldfid {
		// means one fid = 0 the other is default favorite
		return
	}
	if nf.IsLimited(len(oids), s.conf.Fav.DefaultFolderLimit, s.conf.Fav.NormalFolderLimit) {
		err = ecode.FavResourceOverflow
		return
	}
	rows, err := s.favDao.CopyRelations(c, typ, oldmid, mid, oldfid, newfid, oids)
	if err != nil {
		log.Error("s.favDao.CopyRelations(%d,%d,%d,%d,%v) error(%v)", typ, oldmid, mid, oldfid, newfid, oids, err)
		return
	}
	if int(rows) < len(oids) {
		log.Warn("oids len not equal rows(oids:%d,rows:%d,mid:%d), error(%v)", len(oids), int(rows), mid, err)
	}
	if rows > 0 {
		s.cache.Do(c, func(c context.Context) {
			if err := s.favDao.DelRelationFidsMc(c, typ, mid, oids...); err != nil {
				log.Error("s.favDao.DelRelationFidsMc(%d,%d,%v) error(%v)", typ, mid, oids, err)
			}
			if err := s.favDao.DelRecentOidsMc(c, typ, mid); err != nil {
				log.Error("s.favDao.DelRecentOidsMc(%d,%d) error(%v)", typ, mid, err)
			}
			if err := s.favDao.DelRecentResMc(c, model.TypeVideo, mid); err != nil {
				log.Error("s.favDao.DelRecentResMc(%d,%d) error(%v)", typ, mid, err)
			}
			if err := s.favDao.DelRelationOidsCache(c, typ, mid); err != nil {
				log.Error("s.favDao.DelRelationOidsCache(%d,%d) error(%v)", typ, mid, err)
			}
			if err := s.favDao.DelRelationsCache(c, mid, newfid); err != nil {
				log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, newfid, err)
			}
			if err := s.favDao.DelAllRelationsCache(c, mid, newfid); err != nil {
				log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, newfid, err)
			}
		})
	}
	s.favDao.PubCopyFavs(c, typ, mid, oldfid, newfid, rows, oids, time.Now().Unix())
	return
}

// InDefaultFolder detemine oid whether or not in default folder.
func (s *Service) InDefaultFolder(c context.Context, typ int8, mid, oid int64) (isIn bool, err error) {
	un, err := s.favDao.FavedBit(c, typ, mid)
	// un mean user do not had any fav
	if err != nil || un {
		return
	}
	ok, err := s.favDao.ExpireRelationOids(c, typ, mid)
	if err != nil {
		log.Error("s.favDao.ExpireRelationOids(%d,%d) error(%v)", typ, mid, err)
		return
	}
	if ok {
		var faved bool
		if faved, err = s.favDao.IsFavedCache(c, typ, mid, oid); err != nil {
			log.Error("s.favDao.IsFavedCache(%d,%d,%d) error(%v)", typ, mid, oid, err)
			return
		}
		if !faved {
			return
		}
	}
	fids, err := s.favedFids(c, typ, mid, oid)
	if err != nil {
		return
	}
	fvmids := make([]*model.ArgFVmid, 0, len(fids))
	for _, fid := range fids {
		fvmids = append(fvmids, &model.ArgFVmid{Fid: fid, Vmid: mid})
	}
	var folderMap map[string]*model.Folder
	if folderMap, err = s.folders(c, typ, mid, fvmids); err != nil {
		return
	}
	for _, folder := range folderMap {
		if folder.IsDefault() {
			isIn = true
			return
		}
	}
	return
}

// UserList return fav users info by oid.
func (s *Service) UserList(c context.Context, typ int8, oid int64, pn, ps int) (res *model.UserList, err error) {
	var (
		users []*model.User
		start = (pn - 1) * ps
		end   = start + ps - 1
	)
	res = new(model.UserList)
	res.Page.Num = pn
	res.Page.Size = ps
	if res.Page.Total, err = s.favDao.CntUsers(c, typ, oid); err != nil {
		log.Error("s.favDao.CntUsers(%d,%d) error(%v)", typ, oid, err)
		return
	}
	if res.Page.Total > 0 {
		if users, err = s.favDao.Users(c, typ, oid, start, ps); err != nil {
			log.Error("s.favDao.FolderRelations(%d,%d,%d,%d) error(%v)", typ, oid, start, end, err)
			return
		}
	}
	if res.Page.Total == 0 || len(users) == 0 {
		res.List = _emptyUsers
	} else {
		res.List = users
	}
	return
}

// OidCount return resource's oid fav count.
func (s *Service) OidCount(c context.Context, typ int8, oid int64) (count int64, err error) {
	if count, err = s.favDao.OidCountMc(c, typ, oid); err != nil {
		log.Error("s.favDao.OidCountMc(%d,%d) error(%v)", typ, oid, err)
	}
	if count == 0 {
		if count, err = s.favDao.OidCount(c, typ, oid); err != nil {
			log.Error("s.favDao.OidCount(%d,%d) error(%v)", typ, oid, err)
			return
		}
		if count == 0 {
			return
		}
		s.cache.Do(c, func(c context.Context) {
			if err1 := s.favDao.SetOidCountMc(c, typ, oid, count); err1 != nil {
				log.Error("s.favDao.SetOidCountMc(%d,%d,%d) error(%v)", typ, oid, count, err1)
			}
		})
	}
	return
}

// OidsCount return resources's  fav count.
func (s *Service) OidsCount(c context.Context, typ int8, oids []int64) (counts map[int64]int64, err error) {
	counts, misOids, err := s.favDao.OidsCountMc(c, typ, oids)
	if err != nil {
		log.Error("s.favDao.OidsCountMc(%d,%v) error(%v)", typ, oids, err)
	}
	if len(counts) == 0 {
		counts = make(map[int64]int64)
	}
	if len(misOids) == 0 {
		return
	}
	dbCnts, err := s.favDao.OidsCount(c, typ, misOids)
	if err != nil {
		log.Error("s.favDao.OidsCount(%d,%v) error(%v)", typ, misOids, err)
		return
	}
	misCnts := make(map[int64]int64, len(misOids))
	for _, oid := range misOids {
		if cnt, ok := dbCnts[oid]; ok {
			counts[oid] = cnt
		} else {
			counts[oid] = 0
		}
		misCnts[oid] = counts[oid]
	}
	s.cache.Do(c, func(c context.Context) {
		if err1 := s.favDao.SetOidsCountMc(c, typ, misCnts); err1 != nil {
			log.Error("s.favDao.SetOidCountMc(%d,%v) error(%v)", typ, misCnts, err1)
		}
	})
	return
}

// PARRelationsCache return the folder all relations from redis.
func (s *Service) parRelations(c context.Context, typ int8, mid, fid int64, count int, unExpired bool) (res []*model.Favorite, err error) {
	if !unExpired {
		s.favDao.PubInitFolderRelations(c, typ, mid, fid)
		if count > s.conf.Fav.MaxParallelSize { // 减小db压力
			log.Warn("user's favs too large,count(%d)", count)
			err = ecode.FavRetryLater
			return
		}
	}
	size := s.conf.Fav.MaxDataSize
	if count <= size {
		if unExpired {
			if res, err = s.favDao.FolderRelationsCache(c, typ, mid, fid, 0, count); err != nil {
				log.Error("s.favDao.FolderRelationsCache(%d,%d,%d,%d,%d) error(%v)", typ, mid, fid, 0, count, err)
				return
			}
		} else {
			if res, err = s.favDao.FolderRelations(c, typ, mid, fid, 0, count); err != nil {
				log.Error("s.favDao.FolderRelations(%d,%d,%d,%d,%d) error(%v)", typ, mid, fid, 0, count, err)
				return
			}
		}
		return
	}
	g := new(errgroup.Group)
	mux := new(sync.Mutex)
	for i := 0; i < count; i += size {
		n := i
		g.Go(func() (err error) {
			var favs []*model.Favorite
			if unExpired {
				if favs, err = s.favDao.FolderRelationsCache(c, typ, mid, fid, n, n+size-1); err != nil {
					log.Error("s.favDao.FolderRelationsCache(%d,%d,%d,%d,%d) error(%v)", typ, mid, fid, n, n+size-1, err)
					return err
				}
			} else {
				if favs, err = s.favDao.FolderRelations(c, typ, mid, fid, n, size); err != nil {
					log.Error("s.favDao.FolderRelations(%d,%d,%d,%d,%d) error(%v)", typ, mid, fid, n, size, err)
					return err
				}
			}
			mux.Lock()
			res = append(res, favs...)
			mux.Unlock()
			return nil
		})
	}
	g.Wait()
	return
}

// PARRelationsCache return the folder all relations from redis.
func (s *Service) parAllRelations(c context.Context, typ int8, mid, fid int64, count int, unExpired bool) (res []*model.Favorite, err error) {
	size := s.conf.Fav.MaxDataSize
	if count <= size {
		if unExpired {
			res, err = s.favDao.FolderAllRelationsCache(c, typ, mid, fid, 0, count)
			if err != nil {
				log.Error("s.favDao.FolderRelationsCache(%d,%d,%d,%d,%d) error(%v)", typ, mid, fid, 0, count, err)
				return
			}
		} else {
			res, err = s.favDao.FolderAllRelations(c, mid, fid, 0, count)
			if err != nil {
				log.Error("s.favDao.FolderRelations(%d,%d,%d,%d,%d) error(%v)", typ, mid, fid, 0, count, err)
				return
			}
			s.favDao.PubInitAllFolderRelations(c, typ, mid, fid)
		}
		return
	}
	g := new(errgroup.Group)
	mux := new(sync.Mutex)
	for i := 0; i < count; i += size {
		n := i
		g.Go(func() (err error) {
			var favs []*model.Favorite
			if unExpired {
				favs, err = s.favDao.FolderAllRelationsCache(c, typ, mid, fid, n, n+size-1)
				if err != nil {
					log.Error("s.favDao.FolderRelationsCache(%d,%d,%d,%d,%d) error(%v)", typ, mid, fid, n, n+size-1, err)
					return err
				}
			} else {
				favs, err = s.favDao.FolderAllRelations(c, mid, fid, n, size)
				if err != nil {
					log.Error("s.favDao.FolderRelations(%d,%d,%d,%d,%d) error(%v)", typ, mid, fid, n, size, err)
					return err
				}
			}
			mux.Lock()
			res = append(res, favs...)
			mux.Unlock()
			return nil
		})
	}
	g.Wait()
	return
}

// only for medialist
func (s *Service) RecentResources(c context.Context, typ int8, mid int64, size int) (recents []*model.Resource, err error) {
	if typ != model.TypeVideo {
		log.Error("RecentResources only support meidalist")
		return nil, ecode.RequestErr
	}
	un1, err := s.favDao.FavedBit(c, model.TypeVideo, mid)
	// un mean user do not had any fav
	if err != nil {
		return
	}
	un2, err := s.favDao.FavedBit(c, model.TypeMusicNew, mid)
	// un mean user do not had any fav
	if err != nil {
		return
	}
	if un1 && un2 {
		return
	}
	if size > conf.Conf.Fav.MaxPagesize || size <= 0 {
		size = conf.Conf.Fav.MaxRecentSize
	}
	if recents, err = s.favDao.UserRecentResourcesMc(c, typ, mid); err != nil {
		log.Error("s.favDao.UserRecentResourcesMc(%d,%d) error(%v)", typ, mid, err)
		return
	}
	if len(recents) > 0 {
		if len(recents) > size {
			return recents[:size], nil
		}
		return
	}
	folders, err := s.userFolders(c, typ, mid)
	if err != nil {
		log.Error("s.userFolders(%d,%d) error(%v)", typ, mid, err)
		return
	}
	if len(folders) == 0 {
		return
	}
	var favs []*model.Favorite
	if len(folders) == 1 {
		for _, f := range folders {
			unExpired, err1 := s.favDao.ExpireAllRelations(c, mid, f.ID)
			if err1 != nil {
				log.Error("s.favDao.ExpireRelations(%d,%d) error(%v)", mid, f.ID, err1)
				return nil, err1
			}
			if favs, err = s.parAllRelations(c, typ, mid, f.ID, size, unExpired); err != nil {
				return nil, err
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
		if okMap, err = s.favDao.MultiExpireAllRelations(c, mid, fids); err != nil {
			log.Error("s.favDao.MultiExpireAllRelations(%v,%d) error(%v)", fids, err)
			return
		}
		g := new(errgroup.Group)
		mux := new(sync.Mutex)
		for _, f := range folders {
			unExpired := okMap[f.ID]
			fid := f.ID
			g.Go(func() error {
				gfavs, err1 := s.parAllRelations(c, typ, mid, fid, size, unExpired)
				if err1 != nil {
					return err1
				}
				mux.Lock()
				favs = append(favs, gfavs...)
				mux.Unlock()
				return nil
			})
		}
		if err = g.Wait(); err != nil {
			log.Error("g.Wait() error(%v)", err)
		}
	}
	if len(favs) == 0 {
		return
	}
	sort.Slice(favs, func(i, j int) bool {
		return favs[i].MTime > favs[j].MTime
	})
	favMap := make(map[int64]*model.Favorite)
	var recentsRes []*model.Resource
	for _, fav := range favs {
		if _, ok := favMap[fav.ResourceID()]; ok {
			continue
		}
		recentsRes = append(recentsRes, &model.Resource{Oid: fav.Oid, Typ: int32(fav.Type)})
		favMap[fav.ResourceID()] = fav
	}
	if len(recentsRes) > size {
		recents = recentsRes[:size]
	} else {
		recents = recentsRes
	}
	s.cache.Do(c, func(c context.Context) {
		if err1 := s.favDao.SetUserRecentResourcesMc(c, typ, mid, recentsRes); err1 != nil {
			log.Error("s.favDao.SetUserRecentOidsMc(%d,%d,%v) error(%v)", typ, mid, recentsRes, err1)
		}
	})
	return
}

// RecentFavs return 6 oids by mid.
func (s *Service) RecentFavs(c context.Context, typ int8, mid int64, size int) (oids []int64, err error) {
	oids = make([]int64, 0)
	un, err := s.favDao.FavedBit(c, typ, mid)
	// un mean user do not had any fav
	if err != nil || un {
		return
	}
	if size > conf.Conf.Fav.MaxPagesize || size <= 0 {
		size = conf.Conf.Fav.MaxRecentSize
	}
	if oids, err = s.favDao.UserRecentOidsMc(c, typ, mid); err != nil {
		log.Error("s.favDao.UserRecentOidsMc(%d,%d) error(%v)", typ, mid, err)
		return
	}
	if len(oids) > 0 {
		if len(oids) > size {
			return oids[:size], nil
		}
		return
	}
	folders, err := s.userFolders(c, typ, mid)
	if err != nil {
		log.Error("s.userFolders(%d,%d) error(%v)", typ, mid, err)
		return
	}
	if len(folders) == 0 {
		return
	}
	var favs []*model.Favorite
	if len(folders) == 1 {
		for _, f := range folders {
			unExpired, err1 := s.favDao.ExpireRelations(c, mid, f.ID)
			if err1 != nil {
				log.Error("s.favDao.ExpireRelations(%d,%d) error(%v)", mid, f.ID, err1)
				return nil, err1
			}
			if favs, err = s.parRelations(c, typ, mid, f.ID, size, unExpired); err != nil {
				return nil, err
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
		for _, f := range folders {
			unExpired := okMap[f.ID]
			fid := f.ID
			g.Go(func() error {
				gfavs, err1 := s.parRelations(c, typ, mid, fid, size, unExpired)
				if err1 != nil {
					return err1
				}
				mux.Lock()
				favs = append(favs, gfavs...)
				mux.Unlock()
				return nil
			})
		}
		if err = g.Wait(); err != nil {
			log.Error("g.Wait() error(%v)", err)
		}
	}
	if len(favs) == 0 {
		return
	}
	sort.Slice(favs, func(i, j int) bool {
		return favs[i].MTime > favs[j].MTime
	})
	favMap := make(map[int64]*model.Favorite)
	var t []int64
	for _, fav := range favs {
		if _, ok := favMap[fav.Oid]; ok {
			continue
		}
		t = append(t, fav.Oid)
		favMap[fav.Oid] = fav
	}
	if len(t) > size {
		oids = t[:size]
	} else {
		oids = t
	}
	s.cache.Do(c, func(c context.Context) {
		if err1 := s.favDao.SetUserRecentOidsMc(c, typ, mid, oids); err1 != nil {
			log.Error("s.favDao.SetUserRecentOidsMc(%d,%d,%v) error(%v)", typ, mid, oids, err1)
		}
	})
	return
}

// CleanState return this folder clean state.
func (s *Service) CleanState(c context.Context, typ int8, mid, fid int64) (cleanState int, err error) {
	cleanedTime, err := s.favDao.IsCleaned(c, typ, mid, fid)
	if err != nil {
		log.Error("s.favDao.IsCleaned(%d) error(%v)", mid, err)
		return
	}
	if cleanedTime < 0 {
		cleanState = model.StateCleaning
	} else if time.Now().Unix()-cleanedTime < s.cleanCDTime {
		cleanState = model.StateCleanCD
	} else {
		cleanState = model.StateAllowToClean
	}
	return
}

// CleanInvalidArcs clean invalid archives.
func (s *Service) CleanInvalidArcs(c context.Context, typ int8, mid, fid int64) (err error) {
	cleanedTime, err := s.favDao.IsCleaned(c, typ, mid, fid)
	if err != nil {
		log.Error("s.favDao.IsCleaned(%d) error(%v)", mid, err)
		return
	}
	if cleanedTime < 0 {
		log.Warn("FavCleaneInProgress cleanedTime:%d", cleanedTime)
		err = ecode.FavCleaneInProgress
		return
	}
	if time.Now().Unix()-cleanedTime < s.cleanCDTime {
		err = ecode.FavCleanedLocked
		return
	}
	if err = s.favDao.SetCleanedCache(c, typ, mid, fid, -1, s.cleanCDTime); err != nil {
		log.Error("s.favDao.SetCleanedCache(%d,%d,%d,%d) error(%v)", typ, mid, fid, s.cleanCDTime, err)
		return
	}
	s.favDao.PubClean(c, typ, mid, fid, time.Now().Unix())
	return
}
