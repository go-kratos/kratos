package v1

import (
	"context"
	v1pb "go-common/app/interface/live/web-room/api/http/v1"
	"go-common/app/interface/live/web-room/conf"
	"go-common/app/interface/live/web-room/dao"
	"go-common/library/log"
)

// RoomAdminService struct
type RoomAdminService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	dao *dao.Dao
}

//NewRoomAdminService init
func NewRoomAdminService(c *conf.Config) (s *RoomAdminService) {
	s = &RoomAdminService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

// History 相关服务

// GetByRoom implementation
// 获取主播拥有的的所有房管, 无需登录态
// `method:"GET"
func (s *RoomAdminService) GetByRoom(ctx context.Context, req *v1pb.RoomAdminGetByRoomReq) (resp *v1pb.RoomAdminGetByRoomResp, err error) {
	// 默认值
	resp = &v1pb.RoomAdminGetByRoomResp{
		Page: &v1pb.RoomAdminGetByRoomResp_Page{
			Page:       1,
			PageSize:   1,
			TotalPage:  1,
			TotalCount: 0,
		},
		Data: []*v1pb.RoomAdminGetByRoomResp_Data{},
	}

	roomID := req.GetRoomid()
	page := req.GetPage()
	pageSize := req.GetPageSize()

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = int64(10)
	}

	ret, err := s.dao.GetByRoomIDPage(ctx, roomID, page, pageSize)
	if ret == nil {
		log.Info("call GetByAnchor nil mid(%v) err (%v)", roomID, err)
		return
	}

	if err != nil {
		return
	}

	if nil != ret.Page {
		resp.Page = ret.Page
	}
	if nil != ret.Data {
		resp.Data = ret.Data
	}

	return
}
