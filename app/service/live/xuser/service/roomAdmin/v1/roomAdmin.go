package v1

import (
	"context"
	"github.com/pkg/errors"
	v1pb "go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/app/service/live/xuser/conf"
	"go-common/app/service/live/xuser/dao/roomAdmin"
	"go-common/library/ecode"
	"go-common/library/log"
)

// RoomAdminService struct
type RoomAdminService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *roomAdmin.Dao
}

var _pageSize = int64(20)

//NewRoomAdminService init
func NewRoomAdminService(c *conf.Config) (s *RoomAdminService) {
	s = &RoomAdminService{
		conf: c,
		dao:  roomAdmin.New(c),
	}
	return s
}

// GetByUid implementation
// 获取用户拥有的的所有房管身份
func (s *RoomAdminService) GetByUid(ctx context.Context, req *v1pb.RoomAdminGetByUidReq) (resp *v1pb.RoomAdminGetByUidResp, err error) {
	resp = &v1pb.RoomAdminGetByUidResp{}

	uid := req.GetUid()
	if uid <= 0 {
		err = ecode.NoLogin
		return
	}

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	resp, err = s.dao.GetByUidPage(ctx, uid, page, _pageSize)
	log.Info("roomadmin GetByUid uid(%v) page(%v)rst (%v) err (%v)", uid, page, resp, err)

	return
}

// Resign implementation
// 辞职房管
func (s *RoomAdminService) Resign(ctx context.Context, req *v1pb.RoomAdminResignRoomAdminReq) (resp *v1pb.RoomAdminResignRoomAdminResp, err error) {
	resp = &v1pb.RoomAdminResignRoomAdminResp{}
	uid := req.GetUid()
	if uid <= 0 {
		err = ecode.NoLogin
		return
	}
	roomId := req.GetRoomid()
	if roomId <= 0 {
		err = ecode.InvalidParam
		return
	}
	err = s.dao.Del(ctx, uid, roomId)

	if err != nil {
		err = errors.Wrap(err, "server error")
		return
	}
	log.Info("roomadmin Resign uid(%v) roomid(%v)rst (%v) err (%v)", uid, roomId, resp, err)

	return
}

// SearchForAdmin implementation
// 查询需要添加的房管
func (s *RoomAdminService) SearchForAdmin(ctx context.Context, req *v1pb.RoomAdminSearchForAdminReq) (resp *v1pb.RoomAdminSearchForAdminResp, err error) {
	resp = &v1pb.RoomAdminSearchForAdminResp{}

	keyword := req.GetKeyWord()
	uid := req.Uid

	if uid <= 0 {
		err = ecode.InvalidParam
		return
	}

	if keyword == "" {
		err = ecode.ParamInvalid
		return
	}

	resp.Data, err = s.dao.SearchForAdmin(ctx, keyword, uid)

	log.Info("roomadmin SearchForAdmin uid(%v) keyword(%v)rst (%v) err (%v)", uid, keyword, resp, err)

	return
}

// GetByAnchor implementation
// 获取主播拥有的的所有房管
func (s *RoomAdminService) GetByAnchor(ctx context.Context, req *v1pb.RoomAdminGetByAnchorReq) (resp *v1pb.RoomAdminGetByAnchorResp, err error) {
	resp = &v1pb.RoomAdminGetByAnchorResp{}
	uid := req.GetUid()
	if uid <= 0 {
		err = ecode.NoLogin
		return
	}

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	resp, _ = s.dao.GetByAnchorIdPage(ctx, uid, page, _pageSize)
	log.Info("roomadmin GetByAnchor uid(%v)page(%v) rst (%v) err (%v)", uid, page, resp, err)

	return
}

// Dismiss implementation
// 撤销房管
func (s *RoomAdminService) Dismiss(ctx context.Context, req *v1pb.RoomAdminDismissAdminReq) (resp *v1pb.RoomAdminDismissAdminResp, err error) {
	resp = &v1pb.RoomAdminDismissAdminResp{}
	uid := req.GetUid()
	if uid <= 0 {
		err = ecode.ParamInvalid
		return
	}

	anchorId := req.GetAnchorId()
	if anchorId <= 0 {
		err = ecode.NoLogin
		return
	}

	resp, err = s.dao.DismissAnchor(ctx, uid, anchorId)
	log.Info("roomadmin dismiss uid(%v)anchor(%v) rst (%v) err (%v)", uid, anchorId, resp, err)

	return
}

// Appoint implementation
// 任命房管
func (s *RoomAdminService) Appoint(ctx context.Context, req *v1pb.RoomAdminAddReq) (resp *v1pb.RoomAdminAddResp, err error) {
	resp = &v1pb.RoomAdminAddResp{}

	uid := req.GetUid()
	if uid <= 0 {
		err = ecode.ParamInvalid
		return
	}

	anchorId := req.GetAnchorId()
	if anchorId <= 0 {
		err = ecode.NoLogin
		return
	}

	resp, err = s.dao.Add(ctx, uid, anchorId)
	log.Info("roomadmin add uid(%v)anchor(%v) rst (%v) err (%v)", uid, anchorId, resp, err)

	return
}

// IsAny implementation
// 根据登录态获取功能入口是否显示, 需要登录态
func (s *RoomAdminService) IsAny(ctx context.Context, req *v1pb.RoomAdminShowEntryReq) (resp *v1pb.RoomAdminShowEntryResp, err error) {
	resp = &v1pb.RoomAdminShowEntryResp{}

	uid := req.GetUid()
	if uid <= 0 {
		err = ecode.NoLogin
		return
	}

	resp.HasAdmin, err = s.dao.HasAnyAdmin(ctx, uid)
	//spew.Dump(resp.HasAdmin, err)
	if err != nil {
		log.Error("HasAnyAdmin(%v) error(%v)", uid, err)
	}

	log.Info("roomadmin IsAny uid(%v) rst (%v) err (%v)", uid, resp, err)

	return
}

// IsAdmin implementation
// 是否房管
func (s *RoomAdminService) IsAdmin(ctx context.Context, req *v1pb.RoomAdminIsAdminReq) (resp *v1pb.RoomAdminIsAdminResp, err error) {
	//resp = &v1pb.RoomAdminIsAdminResp{}
	uid := req.GetUid()
	anchorId := req.GetAnchorId()
	roomId := req.GetRoomid()

	resp, err = s.dao.IsAdmin(ctx, uid, anchorId, roomId)
	log.Info("roomadmin IsAdmin uid(%v) anchor(%v)roomid(%v) rst (%v) err (%v)", uid, anchorId, roomId, resp, err)

	return
}

// GetByRoom implementation
// 获取主播拥有的的所有房管,房间号维度
func (s *RoomAdminService) GetByRoom(ctx context.Context, req *v1pb.RoomAdminGetByRoomReq) (resp *v1pb.RoomAdminGetByRoomResp, err error) {
	resp = &v1pb.RoomAdminGetByRoomResp{}

	roomId := req.Roomid
	if roomId <= 0 {
		err = ecode.NoLogin
		return
	}

	admins, err := s.dao.GetAllByRoomId(ctx, roomId)
	if err != nil {
		return resp, err
	}

	if admins == nil {
		return resp, nil
	}

	for _, v := range admins {
		resp.Data = append(resp.Data, &v1pb.RoomAdminGetByRoomResp_Data{
			Ctime:  v.Ctime.Time().Format("2006-01-02 15:04:05"),
			Uid:    v.Uid,
			Roomid: v.Roomid,
		})
	}
	log.Info("roomadmin GetByRoom roomid(%v) rst (%v) err (%v)", roomId, resp, err)

	return
}

// IsAdminShort implementation
// 是否房管, 不额外返回用户信息, 不判断是否主播自己
func (s *RoomAdminService) IsAdminShort(ctx context.Context, req *v1pb.RoomAdminIsAdminShortReq) (resp *v1pb.RoomAdminIsAdminShortResp, err error) {
	resp = &v1pb.RoomAdminIsAdminShortResp{}

	uid := req.GetUid()
	roomId := req.GetRoomid()

	rst, err := s.dao.IsAdminByRoomId(ctx, uid, roomId)

	log.Info("roomadmin IsAdminShort uid(%v) roomid(%v) rst (%v) err (%v)", uid, roomId, rst, err)
	resp.Result = rst
	return
}
