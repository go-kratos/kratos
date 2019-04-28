// Package server generate by warden_gen
package server

import (
	"context"
	"go-common/library/ecode"

	pb "go-common/app/service/main/favorite/api"
	"go-common/app/service/main/favorite/model"
	service "go-common/app/service/main/favorite/service"
	"go-common/library/net/rpc/warden"

	empty "github.com/golang/protobuf/ptypes/empty"
)

// New Favorite warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterFavoriteServer(ws.Server(), &server{svr})
	return ws
}

type server struct {
	svr *service.Service
}

var _ pb.FavoriteServer = &server{}

func (s *server) SortFavs(ctx context.Context, req *pb.SortFavsReq) (*empty.Empty, error) {
	var sorts []model.SortFav
	for _, data := range req.Sorts {
		sorts = append(sorts, model.SortFav{
			Pre:    data.Pre,
			Insert: data.Insert,
		})
	}
	if len(sorts) == 0 {
		return nil, ecode.RequestErr
	}
	s.svr.SortFavs(ctx, int8(req.Typ), req.Fid, req.Mid, sorts)
	return &empty.Empty{}, nil
}

func (s *server) AdminUpdateFolder(ctx context.Context, req *pb.AdminUpdateFolderReq) (*empty.Empty, error) {
	attr := int32(req.Attr)
	state := int8(req.State)
	err := s.svr.UpdateFolder(ctx, int8(req.Typ), req.Fid, req.Mid, req.Name, req.Description, req.Cover, 0, &attr, &state)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) UpFolderName(ctx context.Context, req *pb.UpFolderNameReq) (*empty.Empty, error) {
	err := s.svr.UpFolderName(ctx, int8(req.Typ), req.Mid, req.Fid, req.Name, req.Cookie, req.AccessKey)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) FavDelFolders(ctx context.Context, req *pb.FavDelFoldersReq) (*empty.Empty, error) {
	for _, fid := range req.Fids {
		err := s.svr.DelFav(ctx, int8(req.Typ), req.Mid, fid, req.Oid, int8(req.Otype))
		if err != nil {
			return nil, err
		}
	}
	return &empty.Empty{}, nil
}

func (s *server) FavAddFolders(ctx context.Context, req *pb.FavAddFoldersReq) (*empty.Empty, error) {
	for _, fid := range req.Fids {
		err := s.svr.AddFav(ctx, int8(req.Typ), req.Mid, fid, req.Oid, int8(req.Otype))
		if err != nil {
			return nil, err
		}
	}
	return &empty.Empty{}, nil
}

func (s *server) UpFolderAttr(ctx context.Context, req *pb.UpFolderAttrReq) (*empty.Empty, error) {
	err := s.svr.UpFolderAttr(ctx, int8(req.Typ), req.Mid, req.Fid, req.Public)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) CleanState(ctx context.Context, req *pb.CleanStateReq) (*pb.CleanStateReply, error) {
	state, err := s.svr.CleanState(ctx, int8(req.Typ), req.Mid, req.Fid)
	if err != nil {
		return nil, err
	}
	return &pb.CleanStateReply{CleanState: int32(state)}, nil
}

func (s *server) CleanInvalidFavs(ctx context.Context, req *pb.CleanInvalidFavsReq) (*empty.Empty, error) {
	err := s.svr.CleanInvalidArcs(ctx, int8(req.Typ), req.Mid, req.Fid)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) CopyFavs(ctx context.Context, req *pb.CopyFavsReq) (*empty.Empty, error) {
	err := s.svr.CopyFavs(ctx, int8(req.Typ), req.OldMid, req.Mid, req.OldFid, req.NewFid, req.Oids)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) MoveFavs(ctx context.Context, req *pb.MoveFavsReq) (*empty.Empty, error) {
	err := s.svr.MoveFavs(ctx, int8(req.Typ), req.Mid, req.OldFid, req.NewFid, req.Oids)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) SetFolderSort(ctx context.Context, req *pb.SetFolderSortReq) (*empty.Empty, error) {
	err := s.svr.SetFolderSort(ctx, int8(req.Typ), req.Mid, req.Fids)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) FavedUsers(ctx context.Context, req *pb.FavedUsersReq) (*pb.FavedUsersReply, error) {
	list, err := s.svr.UserList(ctx, int8(req.Type), req.Oid, int(req.Pn), int(req.Ps))
	if err != nil {
		return nil, err
	}
	reply := &pb.FavedUsersReply{
		Page: &pb.ModelPage{
			Num:   int32(list.Page.Num),
			Size_: int32(list.Page.Size),
			Count: int32(list.Page.Total),
		},
	}
	for _, data := range list.List {
		reply.User = append(reply.User, &pb.User{
			Id:    data.ID,
			Oid:   data.Oid,
			Mid:   data.Mid,
			Typ:   int32(data.Type),
			State: int32(data.State),
			Ctime: int64(data.CTime),
			Mtime: int64(data.MTime),
		})
	}
	return reply, nil
}

func (s *server) CntUserFolders(ctx context.Context, req *pb.CntUserFoldersReq) (*pb.CntUserFoldersReply, error) {
	count, err := s.svr.CntUserFolders(ctx, int8(req.Typ), req.Mid, req.Vmid)
	if err != nil {
		return nil, err
	}
	return &pb.CntUserFoldersReply{Count: int32(count)}, nil
}

func (s *server) InDefault(ctx context.Context, req *pb.InDefaultFolderReq) (*pb.InDefaultFolderReply, error) {
	isIn, err := s.svr.InDefaultFolder(ctx, int8(req.Typ), req.Mid, req.Oid)
	if err != nil {
		return nil, err
	}
	return &pb.InDefaultFolderReply{IsIn: isIn}, nil
}

func (s *server) MultiDel(ctx context.Context, req *pb.MultiDelReq) (*empty.Empty, error) {
	err := s.svr.MultiDelFavs(ctx, int8(req.Typ), req.Mid, req.Fid, req.Oids)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) MultiAdd(ctx context.Context, req *pb.MultiAddReq) (*empty.Empty, error) {
	err := s.svr.MultiAddFavs(ctx, int8(req.Typ), req.Mid, req.Fid, req.Oids)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) DelFolder(ctx context.Context, req *pb.DelFolderReq) (*empty.Empty, error) {
	err := s.svr.DelFolder(ctx, int8(req.Typ), req.Mid, req.Fid)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) UpdateFolder(ctx context.Context, req *pb.UpdateFolderReq) (*empty.Empty, error) {
	err := s.svr.UpdateFolder(ctx, int8(req.Typ), req.Fid, req.Mid, req.Name, req.Description, req.Cover, req.Public, nil, nil)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *server) AddFolder(ctx context.Context, req *pb.AddFolderReq) (*pb.AddFolderReply, error) {
	fid, err := s.svr.AddFolder(ctx, int8(req.Typ), req.Mid, req.Name, req.Description, req.Cover, req.Public, req.Cookie, req.AccessKey)
	if err != nil {
		return nil, err
	}
	return &pb.AddFolderReply{Fid: fid}, nil
}

// IsFavoredByFid return folders by mid.
func (s *server) IsFavoredByFid(ctx context.Context, req *pb.IsFavoredByFidReq) (*pb.IsFavoredReply, error) {
	ok, err := s.svr.IsFavedByFid(ctx, int8(req.Type), req.Mid, req.Oid, req.Fid)
	return &pb.IsFavoredReply{Faved: ok}, err
}

// UserFolders return folders by mid.
func (s *server) UserFolders(ctx context.Context, req *pb.UserFoldersReq) (*pb.UserFoldersReply, error) {
	if req.Otype <= 0 {
		req.Otype = req.Typ
	}
	uf, err := s.svr.UserFolders(ctx, int8(req.Typ), req.Mid, req.Vmid, req.Oid, int8(req.Otype))
	return &pb.UserFoldersReply{Res: uf}, err
}

// UserFolder return one folder by mid.
func (s *server) UserFolder(ctx context.Context, req *pb.UserFolderReq) (*pb.UserFolderReply, error) {
	f, err := s.svr.UserFolder(ctx, int8(req.Typ), req.Mid, req.Vmid, req.Fid)
	return &pb.UserFolderReply{Res: f}, err
}

// Tlists return partitions .
func (s *server) Tlists(ctx context.Context, req *pb.TlistsReq) (*pb.TlistsReply, error) {
	ts, err := s.svr.Tlists(ctx, int8(req.Tp), req.Mid, req.Uid, req.Fid)
	if err != nil {
		return nil, err
	}
	reply := &pb.TlistsReply{}
	for _, v := range ts {
		reply.Res = append(reply.Res, &pb.ModelPartition{
			Tid:   int32(v.Tid),
			Name:  v.Name,
			Count: int32(v.Count),
		})
	}
	return reply, err
}

// RecentFavs return favs by mid.
func (s *server) RecentFavs(ctx context.Context, req *pb.RecentFavsReq) (*pb.RecentFavsReply, error) {
	ids, err := s.svr.RecentFavs(ctx, int8(req.Tp), req.Mid, int(req.Size_))
	return &pb.RecentFavsReply{Res: ids}, err
}

// RecentFavs return favs by mid.
func (s *server) RecentResources(ctx context.Context, req *pb.RecentResourcesReq) (*pb.RecentResourcesReply, error) {
	recents, err := s.svr.RecentResources(ctx, int8(req.Tp), req.Mid, int(req.Size_))
	return &pb.RecentResourcesReply{Res: recents}, err
}

// Favorites return favorieds info by fid.
func (s *server) Favorites(ctx context.Context, req *pb.FavoritesReq) (*pb.FavoritesReply, error) {
	f, err := s.svr.Favorites(ctx, int8(req.Tp), req.Mid, req.Uid, req.Fid, int(req.Tid), int(req.Tv), int(req.Pn), int(req.Ps), req.Keyword, req.Order)
	if err != nil {
		return nil, err
	}
	reply := &pb.FavoritesReply{Res: &pb.ModelFavorites{Page: &pb.ModelPage{
		Num:   int32(f.Page.Num),
		Size_: int32(f.Page.Size),
		Count: int32(f.Page.Count),
	}}}
	for _, data := range f.List {
		reply.Res.List = append(reply.Res.List, &pb.ModelFavorite{
			Id:    data.ID,
			Oid:   data.Oid,
			Mid:   data.Mid,
			Fid:   data.Fid,
			Type:  int32(data.Type),
			State: int32(data.State),
			Ctime: int64(data.CTime),
			Mtime: int64(data.MTime),
		})
	}
	return reply, err
}

// AddFav add a favorite into folder.
func (s *server) AddFav(ctx context.Context, req *pb.AddFavReq) (*pb.AddFavReply, error) {
	err := s.svr.AddFav(ctx, int8(req.Tp), req.Mid, req.Fid, req.Oid, int8(req.Otype))
	return &pb.AddFavReply{}, err
}

// DelFav delete a favorite.
func (s *server) DelFav(ctx context.Context, req *pb.DelFavReq) (*pb.DelFavReply, error) {
	err := s.svr.DelFav(ctx, int8(req.Tp), req.Mid, req.Fid, req.Oid, int8(req.Otype))
	return &pb.DelFavReply{}, err
}

// IsFavored check if oid faved by user
func (s *server) IsFavored(ctx context.Context, req *pb.IsFavoredReq) (*pb.IsFavoredReply, error) {
	is, err := s.svr.IsFavored(ctx, int8(req.Typ), req.Mid, req.Oid)
	return &pb.IsFavoredReply{Faved: is}, err
}

// IsFavoreds check if oids faved by user
func (s *server) IsFavoreds(ctx context.Context, req *pb.IsFavoredsReq) (*pb.IsFavoredsReply, error) {
	res, err := s.svr.IsFavoreds(ctx, int8(req.Typ), req.Mid, req.Oids)
	return &pb.IsFavoredsReply{Faveds: res}, err
}

// FavoritesAll return favorieds info by fid.
func (s *server) FavoritesAll(ctx context.Context, req *pb.FavoritesReq) (*pb.FavoritesReply, error) {
	if req.Pn <= 0 {
		req.Pn = 1
	}
	if req.Ps <= 0 {
		req.Ps = 20
	}
	var f *model.Favorites
	var err error
	if int8(req.Tp) != model.TypeVideo {
		f, err = s.svr.Favorites(ctx, int8(req.Tp), req.Mid, req.Uid, req.Fid, int(req.Tid), int(req.Tv), int(req.Pn), int(req.Ps), req.Keyword, req.Order)
	} else {
		f, err = s.svr.FavoritesAll(ctx, int8(req.Tp), req.Mid, req.Uid, req.Fid, int(req.Tid), int(req.Tv), int(req.Pn), int(req.Ps), req.Keyword, req.Order)
	}
	if err != nil {
		return nil, err
	}
	reply := &pb.FavoritesReply{Res: &pb.ModelFavorites{Page: &pb.ModelPage{
		Num:   int32(f.Page.Num),
		Size_: int32(f.Page.Size),
		Count: int32(f.Page.Count),
	}}}
	for _, data := range f.List {
		reply.Res.List = append(reply.Res.List, &pb.ModelFavorite{
			Id:    data.ID,
			Oid:   data.Oid,
			Mid:   data.Mid,
			Fid:   data.Fid,
			Type:  int32(data.Type),
			State: int32(data.State),
			Ctime: int64(data.CTime),
			Mtime: int64(data.MTime),
		})
	}
	return reply, err
}

// Folders return folders by ids.
func (s *server) Folders(ctx context.Context, req *pb.FoldersReq) (*pb.FoldersReply, error) {
	var afids []*model.ArgFVmid
	for _, id := range req.Ids {
		afids = append(afids, &model.ArgFVmid{
			Fid:  id.Fid,
			Vmid: id.Mid % 100,
		})
	}
	res, err := s.svr.Folders(ctx, int8(req.Typ), req.Mid, afids)
	if err != nil {
		return nil, err
	}
	return &pb.FoldersReply{Res: res}, nil
}
