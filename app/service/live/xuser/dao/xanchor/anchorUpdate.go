package xanchor

import (
	"context"
	XanchorV1 "go-common/app/service/live/dao-anchor/api/grpc/v1"
)

// UpdateAnchorInfo 更新主播经验db
func (d *Dao) UpdateAnchorInfo(ctx context.Context, params *XanchorV1.AnchorIncreReq) (err error) {
	_, err = d.xuserGRPC.AnchorIncre(ctx, params)
	if err != nil {
		return
	}
	return
}
