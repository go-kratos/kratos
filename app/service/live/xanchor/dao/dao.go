package dao

import (
	"context"

	v1pb "go-common/app/service/live/xanchor/api/grpc/v1"
	"go-common/app/service/live/xanchor/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
)

type refValue struct {
	field string
	v     interface{}
}

// Dao dao
type Dao struct {
	c         *conf.Config
	redis     *redis.Pool
	db        *xsql.DB
	dbLiveApp *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:         c,
		redis:     redis.NewPool(c.Redis),
		db:        xsql.NewMySQL(c.MySQL),
		dbLiveApp: xsql.NewMySQL(c.LiveAppMySQL),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}

// FetchRoomByIDs implementation
// FetchRoomByIDs 查询房间信息
func (d *Dao) FetchRoomByIDs(ctx context.Context, req *v1pb.RoomByIDsReq) (resp *v1pb.RoomByIDsResp, err error) {
	return d.fetchRoomByIDs(ctx, req)
}

// RoomOnlineList implementation
// RoomOnlineList 在线房间列表
func (d *Dao) RoomOnlineList(ctx context.Context, req *v1pb.RoomOnlineListReq) (resp *v1pb.RoomOnlineListResp, err error) {
	return d.roomOnlineList(ctx, req)
}

// RoomCreate implementation
// RoomCreate 房间创建
func (d *Dao) RoomCreate(ctx context.Context, req *v1pb.RoomCreateReq) (resp *v1pb.RoomCreateResp, err error) {
	return d.roomCreate(ctx, req)
}

// RoomUpdate implementation
// RoomUpdate 房间更新
func (d *Dao) RoomUpdate(ctx context.Context, req *v1pb.RoomUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return d.roomUpdate(ctx, req)
}

// RoomBatchUpdate implementation
// RoomBatchUpdate 房间更新
func (d *Dao) RoomBatchUpdate(ctx context.Context, req *v1pb.RoomBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return d.roomBatchUpdate(ctx, req)
}

// RoomExtendUpdate implementation
// RoomExtendUpdate 房间更新
func (d *Dao) RoomExtendUpdate(ctx context.Context, req *v1pb.RoomExtendUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return d.roomExtendUpdate(ctx, req)
}

// RoomExtendBatchUpdate implementation
// RoomExtendBatchUpdate 房间更新
func (d *Dao) RoomExtendBatchUpdate(ctx context.Context, req *v1pb.RoomExtendBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return d.roomExtendBatchUpdate(ctx, req)
}

// RoomExtendIncre implementation
// RoomExtendIncre 房间增量更新
func (d *Dao) RoomExtendIncre(ctx context.Context, req *v1pb.RoomExtendIncreReq) (resp *v1pb.UpdateResp, err error) {
	return d.roomExtendIncre(ctx, req)
}

// RoomExtendBatchIncre implementation
// RoomExtendBatchIncre 房间增量更新
func (d *Dao) RoomExtendBatchIncre(ctx context.Context, req *v1pb.RoomExtendBatchIncreReq) (resp *v1pb.UpdateResp, err error) {
	return d.roomExtendBatchIncre(ctx, req)
}

// RoomTagSet implementation
// RoomTagSet 房间Tag更新
func (d *Dao) RoomTagSet(ctx context.Context, req *v1pb.RoomTagSetReq) (resp *v1pb.UpdateResp, err error) {
	return d.roomTagSet(ctx, req)
}

// AnchorUpdate implementation
// AnchorUpdate 主播更新
func (d *Dao) AnchorUpdate(ctx context.Context, req *v1pb.AnchorUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return d.anchorUpdate(ctx, req)
}

// AnchorBatchUpdate implementation
// AnchorBatchUpdate 主播更新
func (d *Dao) AnchorBatchUpdate(ctx context.Context, req *v1pb.AnchorBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	return d.anchorBatchUpdate(ctx, req)
}

// AnchorIncre implementation
// AnchorIncre 主播增量更新
func (d *Dao) AnchorIncre(ctx context.Context, req *v1pb.AnchorIncreReq) (resp *v1pb.UpdateResp, err error) {
	return d.anchorIncre(ctx, req)
}

// AnchorBatchIncre implementation
// AnchorBatchIncre 主播增量更新
func (d *Dao) AnchorBatchIncre(ctx context.Context, req *v1pb.AnchorBatchIncreReq) (resp *v1pb.UpdateResp, err error) {
	return d.anchorBatchIncre(ctx, req)
}

// AnchorTagSet implementation
// AnchorTagSet 主播Tag更新
func (d *Dao) AnchorTagSet(ctx context.Context, req *v1pb.AnchorTagSetReq) (resp *v1pb.UpdateResp, err error) {
	return d.anchorTagSet(ctx, req)
}
