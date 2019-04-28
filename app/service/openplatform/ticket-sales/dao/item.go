package dao

import (
	"context"
	itemv1 "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
)

//ItemBillInfo 获取商品信息
func (d *Dao) ItemBillInfo(ctx context.Context, itemIDs []int64, scIDs []int64, tkIDs []int64) (*itemv1.BillReply, error) {
	req := &itemv1.BillRequest{
		IDs:   itemIDs,
		ScIDs: scIDs,
		TkIDs: tkIDs,
	}
	return d.itemClient.BillInfo(ctx, req)
}
