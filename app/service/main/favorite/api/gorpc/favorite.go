package favorite

import (
	"context"

	"go-common/app/service/main/favorite/model"
	"go-common/library/net/rpc"
)

const (
	_allFolders     = "RPC.AllFolders"
	_folder         = "RPC.Folder"
	_folders        = "RPC.Folders"
	_addFolder      = "RPC.AddFolder"
	_delFolder      = "RPC.DelFolder"
	_updateFolder   = "RPC.UpdateFolder"
	_favorites      = "RPC.Favorites"
	_add            = "RPC.Add"
	_del            = "RPC.Del"
	_multiAdd       = "RPC.MultiAdd"
	_multiDel       = "RPC.MultiDel"
	_isFav          = "RPC.IsFav"
	_isFavs         = "RPC.IsFavs"
	_isInDefault    = "RPC.InDefault"
	_isFavedByFid   = "RPC.IsFavedByFid"
	_cntUserFolders = "RPC.CntUserFolders"
	_users          = "RPC.Users"
	_tlists         = "RPC.Tlists"
	_recents        = "RPC.Recents"
	// fav v2
	_addVideo = "RPC.AddVideo"
)

const (
	_appid = "community.service.favorite"
)

var (
	_noRes = &struct{}{}
)

type Service struct {
	client *rpc.Client2
}

func New2(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

func (s *Service) AllFolders(c context.Context, arg *model.ArgAllFolders) (res []*model.Folder, err error) {
	err = s.client.Call(c, _allFolders, arg, &res)
	return
}

func (s *Service) Folder(c context.Context, arg *model.ArgFolder) (res *model.Folder, err error) {
	err = s.client.Call(c, _folder, arg, &res)
	return
}

func (s *Service) Folders(c context.Context, arg *model.ArgFolders) (res []*model.Folder, err error) {
	err = s.client.Call(c, _folders, arg, &res)
	return
}

func (s *Service) AddFolder(c context.Context, arg *model.ArgAddFolder) (fid int64, err error) {
	err = s.client.Call(c, _addFolder, arg, &fid)
	return
}

func (s *Service) UpdateFolder(c context.Context, arg *model.ArgUpdateFolder) (err error) {
	err = s.client.Call(c, _updateFolder, arg, _noRes)
	return
}

func (s *Service) DelFolder(c context.Context, arg *model.ArgDelFolder) (err error) {
	err = s.client.Call(c, _delFolder, arg, _noRes)
	return
}

func (s *Service) Favorites(c context.Context, arg *model.ArgFavs) (res *model.Favorites, err error) {
	err = s.client.Call(c, _favorites, arg, &res)
	return
}

func (s *Service) Add(c context.Context, arg *model.ArgAdd) (err error) {
	err = s.client.Call(c, _add, arg, _noRes)
	return
}

func (s *Service) Del(c context.Context, arg *model.ArgDel) (err error) {
	err = s.client.Call(c, _del, arg, _noRes)
	return
}

func (s *Service) MultiAdd(c context.Context, arg *model.ArgMultiAdd) (err error) {
	err = s.client.Call(c, _multiAdd, arg, _noRes)
	return
}

func (s *Service) MultiDel(c context.Context, arg *model.ArgMultiDel) (err error) {
	err = s.client.Call(c, _multiDel, arg, _noRes)
	return
}

func (s *Service) InDefault(c context.Context, arg *model.ArgInDefaultFolder) (faved bool, err error) {
	err = s.client.Call(c, _isInDefault, arg, &faved)
	return
}

func (s *Service) IsFav(c context.Context, arg *model.ArgIsFav) (faved bool, err error) {
	err = s.client.Call(c, _isFav, arg, &faved)
	return
}

func (s *Service) IsFavedByFid(c context.Context, arg *model.ArgIsFavedByFid) (faved bool, err error) {
	err = s.client.Call(c, _isFavedByFid, arg, &faved)
	return
}

func (s *Service) CntUserFolders(c context.Context, arg *model.ArgCntUserFolders) (count int, err error) {
	err = s.client.Call(c, _cntUserFolders, arg, &count)
	return
}

func (s *Service) Users(c context.Context, arg *model.ArgUsers) (res *model.UserList, err error) {
	err = s.client.Call(c, _users, arg, &res)
	return
}

// AddVideo add video fav.
func (s *Service) AddVideo(c context.Context, arg *model.ArgAddVideo) (err error) {
	err = s.client.Call(c, _addVideo, arg, _noRes)
	return
}

// IsFavs .
func (s *Service) IsFavs(c context.Context, arg *model.ArgIsFavs) (res map[int64]bool, err error) {
	err = s.client.Call(c, _isFavs, arg, &res)
	return
}

func (s *Service) Tlists(c context.Context, arg *model.ArgTlists) (res []*model.Partition, err error) {
	err = s.client.Call(c, _tlists, arg, &res)
	return
}

func (s *Service) Recents(c context.Context, arg *model.ArgRecents) (res []int64, err error) {
	err = s.client.Call(c, _recents, arg, &res)
	return
}
