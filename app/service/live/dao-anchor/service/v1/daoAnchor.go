package v1

import (
	"context"

	v1pb "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/app/service/live/dao-anchor/conf"
	"go-common/app/service/live/dao-anchor/dao"
)

// DaoAnchorService struct
type DaoAnchorService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewDaoAnchorService init
func NewDaoAnchorService(c *conf.Config) (s *DaoAnchorService) {
	s = &DaoAnchorService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

// FetchRoomByIDs implementation
// FetchRoomByIDs 查询房间信息
func (s *DaoAnchorService) FetchRoomByIDs(ctx context.Context, req *v1pb.RoomByIDsReq) (resp *v1pb.RoomByIDsResp, err error) {
	return s.dao.FetchRoomByIDs(ctx, req)
}

// RoomOnlineList implementation
// RoomOnlineList 在线房间列表
func (s *DaoAnchorService) RoomOnlineList(ctx context.Context, req *v1pb.RoomOnlineListReq) (resp *v1pb.RoomOnlineListResp, err error) {
	return s.dao.RoomOnlineList(ctx, req)
}

// RoomOnlineListByArea implementation
// RoomOnlineListByArea 分区在线房间列表
func (s *DaoAnchorService) RoomOnlineListByArea(ctx context.Context, req *v1pb.RoomOnlineListByAreaReq) (resp *v1pb.RoomOnlineListByAreaResp, err error) {
	return s.dao.RoomOnlineListByArea(ctx, req)
}

// RoomOnlineListByAttrs implementation
// RoomOnlineListByAttrs 在线房间维度信息(不传attrs，不查询attr)
func (s *DaoAnchorService) RoomOnlineListByAttrs(ctx context.Context, req *v1pb.RoomOnlineListByAttrsReq) (resp *v1pb.RoomOnlineListByAttrsResp, err error) {
	return s.dao.RoomOnlineListByAttrs(ctx, req)
}

// RoomCreate implementation
// RoomCreate 房间创建
func (s *DaoAnchorService) RoomCreate(ctx context.Context, req *v1pb.RoomCreateReq) (resp *v1pb.RoomCreateResp, err error) {
	return s.dao.RoomCreate(ctx, req)
}

// RoomUpdate implementation
// RoomUpdate 房间信息更新
func (s *DaoAnchorService) RoomUpdate(ctx context.Context, req *v1pb.RoomUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomUpdate(ctx, req)
}

// RoomBatchUpdate implementation
// RoomBatchUpdate 房间信息批量更新
func (s *DaoAnchorService) RoomBatchUpdate(ctx context.Context, req *v1pb.RoomBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomBatchUpdate(ctx, req)
}

// RoomExtendUpdate implementation
// RoomExtendUpdate 房间扩展信息更新
func (s *DaoAnchorService) RoomExtendUpdate(ctx context.Context, req *v1pb.RoomExtendUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomExtendUpdate(ctx, req)
}

// RoomExtendBatchUpdate implementation
// RoomExtendBatchUpdate 房间扩展信息批量更新
func (s *DaoAnchorService) RoomExtendBatchUpdate(ctx context.Context, req *v1pb.RoomExtendBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomExtendBatchUpdate(ctx, req)
}

// RoomExtendIncre implementation
// RoomExtendIncre 房间扩展信息增量更新
func (s *DaoAnchorService) RoomExtendIncre(ctx context.Context, req *v1pb.RoomExtendIncreReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomExtendIncre(ctx, req)
}

// RoomExtendBatchIncre implementation
// RoomExtendBatchIncre 房间扩展信息批量增量更新
func (s *DaoAnchorService) RoomExtendBatchIncre(ctx context.Context, req *v1pb.RoomExtendBatchIncreReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomExtendBatchIncre(ctx, req)
}

// RoomTagCreate implementation
// RoomTagCreate 房间Tag创建
func (s *DaoAnchorService) RoomTagCreate(ctx context.Context, req *v1pb.RoomTagCreateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomTagCreate(ctx, req)
}

// RoomAttrCreate implementation
// RoomAttrCreate 房间Attr创建
func (s *DaoAnchorService) RoomAttrCreate(ctx context.Context, req *v1pb.RoomAttrCreateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomAttrCreate(ctx, req)
}

// RoomAttrSetEx implementation
// RoomAttrSetEx 房间Attr更新
func (s *DaoAnchorService) RoomAttrSetEx(ctx context.Context, req *v1pb.RoomAttrSetExReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomAttrSetEx(ctx, req)
}

// AnchorUpdate implementation
// AnchorUpdate 主播信息更新
func (s *DaoAnchorService) AnchorUpdate(ctx context.Context, req *v1pb.AnchorUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.AnchorUpdate(ctx, req)
}

// AnchorBatchUpdate implementation
// AnchorBatchUpdate 主播信息批量更新
func (s *DaoAnchorService) AnchorBatchUpdate(ctx context.Context, req *v1pb.AnchorBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.AnchorBatchUpdate(ctx, req)
}

// AnchorIncre implementation
// AnchorIncre 主播信息增量更新
func (s *DaoAnchorService) AnchorIncre(ctx context.Context, req *v1pb.AnchorIncreReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.AnchorIncre(ctx, req)
}

// AnchorBatchIncre implementation
// AnchorBatchIncre 主播信息批量增量更新
func (s *DaoAnchorService) AnchorBatchIncre(ctx context.Context, req *v1pb.AnchorBatchIncreReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.AnchorBatchIncre(ctx, req)
}

// FetchAreas implementation
// FetchAreas 根据父分区号查询子分区
func (s *DaoAnchorService) FetchAreas(ctx context.Context, req *v1pb.FetchAreasReq) (resp *v1pb.FetchAreasResp, err error) {
	return s.dao.FetchAreas(ctx, req)
}

// FetchAttrByIDs implementation
// FetchAttrByIDs 批量根据房间号查询指标
func (s *DaoAnchorService) FetchAttrByIDs(ctx context.Context, req *v1pb.FetchAttrByIDsReq) (resp *v1pb.FetchAttrByIDsResp, err error) {
	return s.dao.FetchAttrByIDs(ctx, req)
}

// DeleteAttr implementation
// DeleteAttr 删除一个指标
func (s *DaoAnchorService) DeleteAttr(ctx context.Context, req *v1pb.DeleteAttrReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.DeleteAttr(ctx, req)
}
