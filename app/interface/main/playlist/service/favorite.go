package service

import (
	"context"
	"sort"
	"strconv"
	"time"

	"go-common/app/interface/main/playlist/dao"
	"go-common/app/interface/main/playlist/model"
	accwarden "go-common/app/service/main/account/api"
	arcmdl "go-common/app/service/main/archive/api"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
)

const (
	_first       = 1
	_sortDefault = 0
	_sortByMTime = 1
	_sortByView  = 2
)

var _empPlaylists = make([]*model.Playlist, 0)

// White  playlist white list.
func (s *Service) White(c context.Context, mid int64) (res map[string]bool, err error) {
	_, power := s.allowMids[mid]
	res = make(map[string]bool, 1)
	res["power"] = power
	return
}

// Add add playlist.
func (s *Service) Add(c context.Context, mid int64, public int8, name, description, cover, cookie, accessKey string) (pid int64, err error) {
	var (
		fid int64
		ts  = time.Now()
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	if _, ok := s.allowMids[mid]; !ok {
		err = ecode.PlDenied
		return
	}
	arg := &favmdl.ArgAddFolder{Type: favmdl.TypePlayVideo, Mid: mid, Name: name, Description: description, Cover: cover, Public: public, Cookie: cookie, AccessKey: accessKey, RealIP: ip}
	if fid, err = s.fav.AddFolder(c, arg); err != nil {
		dao.PromError("添加播单rpc错误", "s.fav.AddFolder(%v) error(%v)", arg, err)
		return
	}
	if pid, err = s.dao.Add(c, mid, fid); err != nil {
		log.Error("s.dao.Add(%d,%d) error(%v)", mid, fid, err)
	} else if pid > 0 {
		s.cache.Save(func() {
			stat := &model.PlStat{ID: pid, Mid: mid, Fid: fid, MTime: xtime.Time(ts.Unix())}
			s.dao.SetPlStatCache(context.Background(), mid, pid, stat)
		})
		if err = s.dao.RegReply(c, pid, mid); err != nil {
			err = nil
		}
	}
	return
}

// Del delete playlist.
func (s *Service) Del(c context.Context, mid, pid int64) (err error) {
	var (
		affected int64
		stat     *model.PlStat
		ip       = metadata.String(c, metadata.RemoteIP)
	)
	if stat, err = s.plByPid(c, pid); err != nil {
		log.Error("s.plByPid(%d,%d) error(%v)", mid, pid, err)
		return
	}
	arg := &favmdl.ArgDelFolder{Type: favmdl.TypePlayVideo, Mid: mid, Fid: stat.Fid, RealIP: ip}
	if err = s.fav.DelFolder(c, arg); err != nil {
		dao.PromError("删除播单rpc错误", "s.fav.DelFolder(%+v) error(%v)", arg, err)
		return
	}
	if affected, err = s.dao.Del(c, pid); err != nil {
		log.Error("s.dao.Del(%d) error(%v)", pid, err)
		return
	} else if affected > 0 {
		s.dao.DelPlCache(c, mid, pid)
	}
	return
}

// Update update playlist.
func (s *Service) Update(c context.Context, mid, pid int64, public int8, name, description, cover, cookie, accessKey string) (err error) {
	var (
		stat *model.PlStat
		ip   = metadata.String(c, metadata.RemoteIP)
	)
	if stat, err = s.plByPid(c, pid); err != nil {
		log.Error("s.plByPid(%d) error(%v)", pid, err)
		return
	}
	arg := &favmdl.ArgUpdateFolder{Type: favmdl.TypePlayVideo, Fid: stat.Fid, Mid: mid, Name: name, Description: description, Cover: cover, Public: public, Cookie: cookie, AccessKey: accessKey, RealIP: ip}
	if err = s.fav.UpdateFolder(c, arg); err != nil {
		dao.PromError("更新播单rpc错误", "s.fav.UpdateFolder(%+v) error(%v)", arg, err)
		return
	}
	s.updatePlTime(c, mid, pid)
	return
}

func (s *Service) updatePlTime(c context.Context, mid, pid int64) (err error) {
	var (
		affected int64
		ts       = time.Now()
		stat     *model.PlStat
	)
	if affected, err = s.dao.Update(c, pid); err != nil {
		err = nil
		log.Error("s.dao.Update(%d) error(%v)", pid, err)
		return
	} else if affected > 0 {
		s.cache.Save(func() {
			if stat, err = s.plByPid(context.Background(), pid); err != nil {
				err = nil
			} else {
				stat.MTime = xtime.Time(ts.Unix())
				s.dao.SetPlStatCache(context.Background(), mid, pid, stat)
			}
		})
	}
	return
}

// Info playlist stat info.
func (s *Service) Info(c context.Context, mid, pid int64) (res *model.Playlist, err error) {
	var (
		fav       *favmdl.Folder
		stat      *model.PlStat
		infoReply *accwarden.InfoReply
		isFav     bool
		ip        = metadata.String(c, metadata.RemoteIP)
	)
	if stat, err = s.plByPid(c, pid); err != nil {
		return
	}
	if stat == nil || stat.ID == 0 {
		err = ecode.PlNotExist
		dao.PromError("Info:播单不存在", "s.plByPid(%d) error(%v)", pid, stat)
		return
	}
	arg := &favmdl.ArgFolder{Type: favmdl.TypePlayVideo, Fid: stat.Fid, Mid: stat.Mid, RealIP: ip}
	if fav, err = s.fav.Folder(c, arg); err != nil || fav == nil {
		dao.PromError("Info Forder:rpc错误", "s.fav.Folder(%+v) error(%v)", arg, err)
		return
	}
	if fav.State == favmdl.StateIsDel {
		err = ecode.PlNotExist
		dao.PromError("InfoFav:播单不存在", "s.fav.Folder(%d) error(%v)", pid, err)
		return
	}
	// author
	if infoReply, err = s.accClient.Info3(c, &accwarden.MidReq{Mid: fav.Mid, RealIp: ip}); err != nil {
		dao.PromError("账号Info:grpc错误", "s.accClient.Info3 error(%v)", err)
		return
	}
	if mid > 0 {
		if isFav, err = s.fav.IsFav(c, &favmdl.ArgIsFav{Type: favmdl.TypePlayList, Mid: mid, Oid: pid, RealIP: ip}); err != nil {
			log.Error("s.fav.IsFav(%d,%d) error(%d)", mid, pid, err)
			err = nil
		}
	}
	owner := &arcmdl.Author{Mid: fav.Mid, Name: infoReply.Info.Name, Face: infoReply.Info.Face}
	fav.MTime = stat.MTime
	res = &model.Playlist{Pid: pid, Folder: fav, Stat: &model.Stat{Pid: stat.ID, View: stat.View, Reply: stat.Reply, Fav: stat.Fav, Share: stat.Share}, Author: owner, IsFavorite: isFav}
	return
}

func (s *Service) plInfo(c context.Context, mid, pid int64, ip string) (res *model.PlStat, err error) {
	var fav *favmdl.Folder
	if res, err = s.plByPid(c, pid); err != nil {
		return
	}
	if res == nil || res.ID == 0 {
		err = ecode.PlNotExist
		log.Error("s.plByPid(%d) res(%v)", pid, res)
		return
	}
	arg := &favmdl.ArgFolder{Type: favmdl.TypePlayVideo, Fid: res.Fid, Mid: res.Mid, RealIP: ip}
	if fav, err = s.fav.Folder(c, arg); err != nil || fav == nil {
		log.Error("s.fav.Folder(%+v) error(%v)", arg, err)
		return
	}
	if fav.State == favmdl.StateIsDel {
		err = ecode.PlNotExist
		log.Error("s.fav.Folder(%d) state(%d)", pid, fav.State)
		return
	}
	if mid > 0 && mid != res.Mid {
		err = ecode.PlNotUser
	}
	return
}

// List playlist.
func (s *Service) List(c context.Context, mid int64, pn, ps, sortType int) (res []*model.Playlist, count int, err error) {
	var (
		start   = (pn - 1) * ps
		end     = start + ps - 1
		plStats []*model.PlStat
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	if plStats, err = s.plsByMid(c, mid); err != nil {
		return
	}
	count = len(plStats)
	if count == 0 || count < start {
		res = _empPlaylists
		return
	}
	switch sortType {
	case _sortDefault, _sortByMTime:
		sort.Slice(plStats, func(i, j int) bool { return plStats[i].MTime > plStats[j].MTime })
	case _sortByView:
		sort.Slice(plStats, func(i, j int) bool { return plStats[i].View > plStats[j].View })
	}
	if count > end {
		plStats = plStats[start : end+1]
	} else {
		plStats = plStats[start:]
	}
	res, err = s.batchFav(c, mid, plStats, ip)
	return
}

//AddFavorite add playlist to favorite.
func (s *Service) AddFavorite(c context.Context, mid, pid int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if _, err = s.Info(c, 0, pid); err != nil {
		return
	}
	arg := &favmdl.ArgAdd{Type: favmdl.TypePlayList, Mid: mid, Oid: pid, Fid: 0, RealIP: ip}
	if err = s.fav.Add(c, arg); err != nil {
		dao.PromError("rpc:添加播单收藏", "s.fav.Add(%+v) error(%v)", arg, err)
	}
	return
}

// DelFavorite del playlist from favorite.
func (s *Service) DelFavorite(c context.Context, mid, pid int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &favmdl.ArgDel{Type: favmdl.TypePlayList, Mid: mid, Oid: pid, Fid: 0, RealIP: ip}
	if err = s.fav.Del(c, arg); err != nil {
		dao.PromError("rpc:删除播单收藏", "s.fav.Del(%+v) error(%v)", arg, err)
	}
	return
}

//ListFavorite playlist list.
func (s *Service) ListFavorite(c context.Context, mid, vmid int64, pn, ps, sortType int) (res []*model.Playlist, count int, err error) {
	var (
		plStats []*model.PlStat
		favRes  *favmdl.Favorites
		pids    []int64
		tmpFavs map[int64]*favmdl.Favorite
		tmpRs   []*model.Playlist
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	arg := &favmdl.ArgFavs{Type: favmdl.TypePlayList, Mid: mid, Vmid: vmid, Fid: 0, Pn: pn, Ps: ps, RealIP: ip}
	if favRes, err = s.fav.Favorites(c, arg); err != nil {
		dao.PromError("rpc:播单收藏列表", "s.fav.Favorites(%+v) error(%v)", arg, err)
		return
	}
	if favRes == nil || len(favRes.List) == 0 {
		res = _empPlaylists
		return
	}
	tmpFavs = make(map[int64]*favmdl.Favorite)
	for _, fav := range favRes.List {
		pids = append(pids, fav.Oid)
		tmpFavs[fav.Oid] = fav
	}
	if plStats, err = s.plsByPid(c, pids); err != nil {
		return
	}
	count = favRes.Page.Count
	tmpRs, err = s.batchFav(c, mid, plStats, ip)
	for _, v := range tmpRs {
		v.FavoriteTime = tmpFavs[v.Pid].MTime
		res = append(res, v)
	}
	return
}

func (s *Service) batchFav(c context.Context, uid int64, plStats []*model.PlStat, ip string) (res []*model.Playlist, err error) {
	var (
		fVMids   []*favmdl.ArgFVmid
		tmpStats map[string]*model.PlStat
		favRes   []*favmdl.Folder
		stat     *model.Stat
	)
	tmpStats = make(map[string]*model.PlStat)
	for _, v := range plStats {
		statKey := strconv.FormatInt(v.Mid, 10) + "_" + strconv.FormatInt(v.Fid, 10)
		tmpStats[statKey] = &model.PlStat{ID: v.ID, Mid: v.Mid, Fid: v.Fid, View: v.View, Reply: v.Reply, Fav: v.Fav, Share: v.Share, MTime: v.MTime}
		fVMids = append(fVMids, &favmdl.ArgFVmid{Fid: v.Fid, Vmid: v.Mid})
	}
	arg := &favmdl.ArgFolders{Type: favmdl.TypePlayVideo, Mid: uid, FVmids: fVMids, RealIP: ip}
	if favRes, err = s.fav.Folders(c, arg); err != nil {
		dao.PromError("rpc:批量获取播单列表", "s.fav.Folders(%+v) error(%v)", arg, err)
		return
	}
	for _, fav := range favRes {
		statKey := strconv.FormatInt(fav.Mid, 10) + "_" + strconv.FormatInt(fav.ID, 10)
		plStat := tmpStats[statKey]
		stat = &model.Stat{Pid: plStat.ID, View: plStat.View, Fav: plStat.Fav, Reply: plStat.Reply, Share: plStat.Share}
		fav.MTime = plStat.MTime
		res = append(res, &model.Playlist{Pid: plStat.ID, Folder: fav, Stat: stat})
	}
	return
}
