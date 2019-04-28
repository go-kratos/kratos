package v1

import (
	"context"

	v1pb "go-common/app/service/live/xanchor/api/grpc/v1"
	"go-common/app/service/live/xanchor/conf"
	"go-common/app/service/live/xanchor/dao"
)

// XAnchorService struct
type XAnchorService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewXAnchorService init
func NewXAnchorService(c *conf.Config) (s *XAnchorService) {
	s = &XAnchorService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

// FetchRoomByIDs implementation
// FetchRoomByIDs 查询房间信息
func (s *XAnchorService) FetchRoomByIDs(ctx context.Context, req *v1pb.RoomByIDsReq) (resp *v1pb.RoomByIDsResp, err error) {
	return s.dao.FetchRoomByIDs(ctx, req)
}

// RoomOnlineList implementation
// RoomOnlineList 在线房间列表
func (s *XAnchorService) RoomOnlineList(ctx context.Context, req *v1pb.RoomOnlineListReq) (resp *v1pb.RoomOnlineListResp, err error) {
	return s.dao.RoomOnlineList(ctx, req)
}

// RoomCreate implementation
// RoomCreate 房间创建
func (s *XAnchorService) RoomCreate(ctx context.Context, req *v1pb.RoomCreateReq) (resp *v1pb.RoomCreateResp, err error) {
	return s.dao.RoomCreate(ctx, req)
}

// RoomUpdate implementation
// RoomUpdate 房间信息更新
func (s *XAnchorService) RoomUpdate(ctx context.Context, req *v1pb.RoomUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomUpdate(ctx, req)
}

// RoomBatchUpdate implementation
// RoomBatchUpdate 房间信息批量更新
func (s *XAnchorService) RoomBatchUpdate(ctx context.Context, req *v1pb.RoomBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomBatchUpdate(ctx, req)
}

// RoomExtendUpdate implementation
// RoomExtendUpdate 房间扩展信息更新
func (s *XAnchorService) RoomExtendUpdate(ctx context.Context, req *v1pb.RoomExtendUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomExtendUpdate(ctx, req)
}

// RoomExtendBatchUpdate implementation
// RoomExtendBatchUpdate 房间扩展信息批量更新
func (s *XAnchorService) RoomExtendBatchUpdate(ctx context.Context, req *v1pb.RoomExtendBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomExtendBatchUpdate(ctx, req)
}

// RoomExtendIncre implementation
// RoomExtendIncre 房间扩展信息增量更新
func (s *XAnchorService) RoomExtendIncre(ctx context.Context, req *v1pb.RoomExtendIncreReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomExtendIncre(ctx, req)
}

// RoomExtendBatchIncre implementation
// RoomExtendBatchIncre 房间扩展信息批量增量更新
func (s *XAnchorService) RoomExtendBatchIncre(ctx context.Context, req *v1pb.RoomExtendBatchIncreReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomExtendBatchIncre(ctx, req)
}

// RoomTagSet implementation
// RoomTagSet 房间Tag更新
func (s *XAnchorService) RoomTagSet(ctx context.Context, req *v1pb.RoomTagSetReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.RoomTagSet(ctx, req)
}

// AnchorUpdate implementation
// AnchorUpdate 主播信息更新
func (s *XAnchorService) AnchorUpdate(ctx context.Context, req *v1pb.AnchorUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.AnchorUpdate(ctx, req)
}

// AnchorBatchUpdate implementation
// AnchorBatchUpdate 主播信息批量更新
func (s *XAnchorService) AnchorBatchUpdate(ctx context.Context, req *v1pb.AnchorBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.AnchorBatchUpdate(ctx, req)
}

// AnchorIncre implementation
// AnchorIncre 主播信息增量更新
func (s *XAnchorService) AnchorIncre(ctx context.Context, req *v1pb.AnchorIncreReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.AnchorIncre(ctx, req)
}

// AnchorBatchIncre implementation
// AnchorBatchIncre 主播信息批量增量更新
func (s *XAnchorService) AnchorBatchIncre(ctx context.Context, req *v1pb.AnchorBatchIncreReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.AnchorBatchIncre(ctx, req)
}

// AnchorTagSet implementation
// AnchorTagSet 主播Tag更新
func (s *XAnchorService) AnchorTagSet(ctx context.Context, req *v1pb.AnchorTagSetReq) (resp *v1pb.UpdateResp, err error) {
	return s.dao.AnchorTagSet(ctx, req)
}
