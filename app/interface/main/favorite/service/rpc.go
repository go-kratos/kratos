package service

import (
	"context"

	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	favpb "go-common/app/service/main/favorite/api"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/time"
)

// ArcRPC find archive by rpc
func (s *Service) ArcRPC(c context.Context, aid int64) (a *api.Arc, err error) {
	argAid := &arcmdl.ArgAid2{
		Aid: aid,
	}
	if a, err = s.arcRPC.Archive3(c, argAid); err != nil {
		log.Error("s.arcRPC.Archive3(%v), error(%v)", argAid, err)
	}
	if !a.IsNormal() {
		err = ecode.ArchiveNotExist
	}
	return
}

// ArcsRPC find archives by rpc
func (s *Service) ArcsRPC(c context.Context, aids []int64) (as map[int64]*api.Arc, err error) {
	if len(aids) == 0 {
		return
	}
	argAids := &arcmdl.ArgAids2{
		Aids: aids,
	}
	if as, err = s.arcRPC.Archives3(c, argAids); err != nil {
		log.Error("s.arcRPC.Archives3(%v, archives), err(%v)", argAids, err)
	}
	return
}

// AddFavRPC add favorite .
func (s *Service) AddFavRPC(c context.Context, typ int8, mid, aid, fid int64) (err error) {
	arg := &favpb.AddFavReq{Tp: int32(typ), Mid: mid, Oid: aid, Fid: fid}
	_, err = s.favClient.AddFav(c, arg)
	if err != nil {
		log.Error("s.favClient.AddFav(%+v) error(%v)", arg, err)
	}
	return
}

// DelFavRPC del favorite .
func (s *Service) DelFavRPC(c context.Context, typ int8, mid, aid, fid int64) (err error) {
	arg := &favpb.DelFavReq{Tp: int32(typ), Mid: mid, Oid: aid, Fid: fid}
	if _, err = s.favClient.DelFav(c, arg); err != nil {
		log.Error("s.favClient.DelFavRPC(%+v) error(%v)", arg, err)
	}
	return
}

// FavoritesRPC favorites list.
func (s *Service) FavoritesRPC(c context.Context, typ int8, mid, vmid, fid int64, tid int, keyword, order string, pn, ps int) (favs *favmdl.Favorites, err error) {
	arg := &favpb.FavoritesReq{Tp: int32(typ), Mid: mid, Uid: vmid, Fid: fid, Tid: int32(tid), Keyword: keyword, Order: order, Pn: int32(pn), Ps: int32(ps)}
	var reply *favpb.FavoritesReply
	if reply, err = s.favClient.Favorites(c, arg); err != nil {
		log.Error("s.favClient.Favorites(%+v) error(%v)", arg, err)
		return
	}
	favs = &favmdl.Favorites{}
	favs.Page.Count = int(reply.Res.Page.Count)
	favs.Page.Num = int(reply.Res.Page.Num)
	favs.Page.Size = int(reply.Res.Page.Size_)
	for _, data := range reply.Res.List {
		favs.List = append(favs.List, &favmdl.Favorite{
			ID:    data.Id,
			Oid:   data.Oid,
			Mid:   data.Mid,
			Fid:   data.Fid,
			Type:  int8(data.Type),
			State: int8(data.State),
			CTime: time.Time(data.Ctime),
			MTime: time.Time(data.Mtime),
		})
	}
	return
}

// IsFavByFidRPC return user whether favored
func (s *Service) IsFavByFidRPC(c context.Context, typ int8, mid, aid int64, fid int64) (res bool, err error) {
	arg := &favpb.IsFavoredByFidReq{Type: int32(typ), Mid: mid, Oid: aid, Fid: fid}
	var reply *favpb.IsFavoredReply
	if reply, err = s.favClient.IsFavoredByFid(c, arg); err != nil {
		log.Error("s.favClient.IsFavoredByFid(%+v) error(%v)", arg, err)
		return
	}
	res = reply.Faved
	return
}

// IsFavRPC return user whether favored
func (s *Service) IsFavRPC(c context.Context, typ int8, mid, aid int64) (res bool, err error) {
	arg := &favpb.IsFavoredReq{Typ: int32(typ), Mid: mid, Oid: aid}
	var reply *favpb.IsFavoredReply
	if reply, err = s.favClient.IsFavored(c, arg); err != nil {
		log.Error("s.favClient.IsFavored(%+v) error(%v)", arg, err)
		return
	}
	res = reply.Faved
	return
}

// InDefaultRPC return aid whether in default fodler.
func (s *Service) InDefaultRPC(c context.Context, typ int8, mid, aid int64) (res bool, err error) {
	arg := &favpb.InDefaultFolderReq{Typ: int32(typ), Mid: mid, Oid: aid}
	var reply *favpb.InDefaultFolderReply
	if reply, err = s.favClient.InDefault(c, arg); err != nil {
		log.Error("s.favClient.IsFavored(%+v) error(%v)", arg, err)
		return
	}
	res = reply.IsIn
	return
}

// IsFavsRPC return user's oids whether favored
func (s *Service) IsFavsRPC(c context.Context, typ int8, mid int64, aids []int64) (res map[int64]bool, err error) {
	arg := &favpb.IsFavoredsReq{Typ: int32(typ), Mid: mid, Oids: aids}
	var reply *favpb.IsFavoredsReply
	if reply, err = s.favClient.IsFavoreds(c, arg); err != nil {
		log.Error("s.favClient.IsFavoreds(%+v) error(%v)", arg, err)
		return
	}
	res = reply.Faveds
	return
}

// AllFoldersRPC user's folders list.
func (s *Service) AllFoldersRPC(c context.Context, typ int8, mid, vmid, oid int64, ip string) (fs []*favmdl.Folder, err error) {
	arg := &favpb.UserFoldersReq{Typ: int32(typ), Mid: mid, Vmid: vmid, Oid: oid}
	var reply *favpb.UserFoldersReply
	if reply, err = s.favClient.UserFolders(c, arg); err != nil {
		log.Error("s.favClient.UserFolders(%+v) error(%v)", arg, err)
		return
	}
	fs = reply.GetRes()
	return
}

// FolderRPC user's folder .
func (s *Service) FolderRPC(c context.Context, typ int8, fid, mid, vmid int64) (f *favmdl.Folder, err error) {
	arg := &favpb.UserFolderReq{Typ: int32(typ), Mid: mid, Vmid: vmid, Fid: fid}
	var reply *favpb.UserFolderReply
	if reply, err = s.favClient.UserFolder(c, arg); err != nil {
		log.Error("s.favClient.UserFolder(%+v) error(%v)", arg, err)
		return
	}
	f = reply.GetRes()
	return
}

// TlistsRPC archive type list.
func (s *Service) TlistsRPC(c context.Context, typ int8, mid, vmid, fid int64) (ps []*favmdl.Partition, err error) {
	arg := &favpb.TlistsReq{Tp: int32(typ), Mid: mid, Uid: vmid, Fid: fid}
	var reply *favpb.TlistsReply
	if reply, err = s.favClient.Tlists(c, arg); err != nil {
		log.Error("s.favClient.Tlists(%+v) error(%v)", arg, err)
		return
	}
	ps = make([]*favmdl.Partition, 0)
	for _, v := range reply.Res {
		ps = append(ps, &favmdl.Partition{
			Tid:   int(v.Tid),
			Name:  v.Name,
			Count: int(v.Count),
		})
	}
	return
}

// RecentsRPC recent favs .
func (s *Service) RecentsRPC(c context.Context, typ int8, mid int64, size int) (aids []int64, err error) {
	arg := &favpb.RecentFavsReq{Tp: int32(typ), Mid: mid, Size_: int32(size)}
	var reply *favpb.RecentFavsReply
	if reply, err = s.favClient.RecentFavs(c, arg); err != nil {
		log.Error("s.favClient.RecentFavs(%+v) error(%v)", arg, err)
		return
	}
	aids = reply.GetRes()
	return
}

// TypesRPC find all archives's type by rpc
func (s *Service) TypesRPC(c context.Context) (ats map[int16]*arcmdl.ArcType, err error) {
	if ats, err = s.arcRPC.Types2(c); err != nil {
		log.Error("s.arcRPC.Types2(), error(%v)", err)
	}
	return
}
