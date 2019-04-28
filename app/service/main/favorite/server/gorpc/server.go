package rpc

import (
	"go-common/app/service/main/favorite/conf"
	"go-common/app/service/main/favorite/model"
	"go-common/app/service/main/favorite/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC favorite rpc.
type RPC struct {
	c *conf.Config
	s *service.Service
}

// New init rpc.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{
		c: c,
		s: s,
	}
	svr = rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Ping check connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// Folder return folder by mid.
func (r *RPC) Folder(c context.Context, a *model.ArgFolder, res *model.Folder) (err error) {
	var fl *model.Folder
	if fl, err = r.s.Folder(c, a.Type, a.Mid, a.Vmid, a.Fid); err == nil {
		*res = *fl
	}
	return
}

// Folders return folder by mid.
func (r *RPC) Folders(c context.Context, a *model.ArgFolders, res *[]*model.Folder) (err error) {
	*res, err = r.s.Folders(c, a.Type, a.Mid, a.FVmids)
	return
}

// AllFolders return users folders.
func (r *RPC) AllFolders(c context.Context, a *model.ArgAllFolders, res *[]*model.Folder) (err error) {
	*res, err = r.s.UserFolders(c, a.Type, a.Mid, a.Vmid, a.Oid, a.Type)
	return
}

// AddFolder add a Folder.
func (r *RPC) AddFolder(c context.Context, a *model.ArgAddFolder, res *int64) (err error) {
	var fid int64
	if fid, err = r.s.AddFolder(c, a.Type, a.Mid, a.Name, a.Description, a.Cover, int32(a.Public), a.Cookie, a.AccessKey); err == nil {
		*res = fid
	}
	return
}

// UpdateFolder update a Folder.
func (r *RPC) UpdateFolder(c context.Context, a *model.ArgUpdateFolder, res *struct{}) (err error) {
	err = r.s.UpdateFolder(c, a.Type, a.Fid, a.Mid, a.Name, a.Description, a.Cover, int32(a.Public), nil, nil)
	return
}

// DelFolder del a folder.
func (r *RPC) DelFolder(c context.Context, a *model.ArgDelFolder, res *struct{}) (err error) {
	err = r.s.DelFolder(c, a.Type, a.Mid, a.Fid)
	return
}

// Favorites return favorites by mid.
func (r *RPC) Favorites(c context.Context, a *model.ArgFavs, res *model.Favorites) (err error) {
	var fs *model.Favorites
	if fs, err = r.s.Favorites(c, a.Type, a.Mid, a.Vmid, a.Fid, a.Tid, a.Tv, a.Pn, a.Ps, a.Keyword, a.Order); err == nil {
		*res = *fs
	}
	return
}

// Add add a favorite relation.
func (r *RPC) Add(c context.Context, a *model.ArgAdd, res *struct{}) (err error) {
	err = r.s.AddFav(c, a.Type, a.Mid, a.Fid, a.Oid, a.Type)
	return
}

// Del del a favorite relation.
func (r *RPC) Del(c context.Context, a *model.ArgDel, res *struct{}) (err error) {
	err = r.s.DelFav(c, a.Type, a.Mid, a.Fid, a.Oid, a.Type)
	return
}

// Adds add a resource to folders.
func (r *RPC) Adds(c context.Context, a *model.ArgAdds, res *struct{}) (err error) {
	for _, fid := range a.Fids {
		err = r.s.AddFav(c, a.Type, a.Mid, fid, a.Oid, a.Type)
	}
	return
}

// Dels del a resource in fodlers.
func (r *RPC) Dels(c context.Context, a *model.ArgDels, res *struct{}) (err error) {
	for _, fid := range a.Fids {
		err = r.s.DelFav(c, a.Type, a.Mid, fid, a.Oid, a.Type)
	}
	return
}

// MultiAdd multi add favorite relations.
func (r *RPC) MultiAdd(c context.Context, a *model.ArgMultiAdd, res *struct{}) (err error) {
	err = r.s.MultiAddFavs(c, a.Type, a.Mid, a.Fid, a.Oids)
	return
}

// MultiDel multi del favorite relations.
func (r *RPC) MultiDel(c context.Context, a *model.ArgMultiDel, res *struct{}) (err error) {
	err = r.s.MultiDelFavs(c, a.Type, a.Mid, a.Fid, a.Oids)
	return
}

// IsFav check favorited relation.
func (r *RPC) IsFav(c context.Context, a *model.ArgIsFav, faved *bool) (err error) {
	*faved, err = r.s.IsFavored(c, a.Type, a.Mid, a.Oid)
	return
}

// IsFavs return favored relation map.
func (r *RPC) IsFavs(c context.Context, a *model.ArgIsFavs, res *map[int64]bool) (err error) {
	*res, err = r.s.IsFavoreds(c, a.Type, a.Mid, a.Oids)
	return
}

// InDefault return favored in default folder.
func (r *RPC) InDefault(c context.Context, a *model.ArgInDefaultFolder, in *bool) (err error) {
	*in, err = r.s.InDefaultFolder(c, a.Type, a.Mid, a.Oid)
	return
}

// IsFavedByFid check the oid and fid relation.
func (r *RPC) IsFavedByFid(c context.Context, a *model.ArgIsFavedByFid, faved *bool) (err error) {
	*faved, err = r.s.IsFavedByFid(c, a.Type, a.Mid, a.Oid, a.Fid)
	return
}

// CntUserFolders count user's folders.
func (r *RPC) CntUserFolders(c context.Context, a *model.ArgCntUserFolders, count *int) (err error) {
	*count, err = r.s.CntUserFolders(c, a.Type, a.Mid, a.Vmid)
	return
}

// Users return favored users by mid.
func (r *RPC) Users(c context.Context, a *model.ArgUsers, res *model.UserList) (err error) {
	var us *model.UserList
	if us, err = r.s.UserList(c, a.Type, a.Oid, a.Pn, a.Ps); err == nil {
		*res = *us
	}
	return
}

// under v2 ===

// AddVideo add a favorite video.
func (r *RPC) AddVideo(c context.Context, a *model.ArgAddVideo, res *struct{}) (err error) {
	for _, fid := range a.Fids {
		err = r.s.AddFav(c, model.TypeVideo, a.Mid, fid, a.Aid, model.TypeVideo)
	}
	return
}
