package dao

import (
	"context"

	ugcpay_rank "go-common/app/service/main/ugcpay-rank/api/v1"
)

// RankElecUpdateOrder .
func (d *Dao) RankElecUpdateOrder(ctx context.Context, avID, upMID, payMID, ver, fee int64) (err error) {
	var (
		req = &ugcpay_rank.RankElecUpdateOrderReq{
			AVID:   avID,
			UPMID:  upMID,
			PayMID: payMID,
			Ver:    ver,
			Fee:    fee,
		}
	)
	_, err = d.ugcPayRankAPI.RankElecUpdateOrder(ctx, req)
	return
}

// RankElecUpdateMessage .
func (d *Dao) RankElecUpdateMessage(ctx context.Context, avID, upMID, payMID, ver int64, message string, hidden bool) (err error) {
	var (
		req = &ugcpay_rank.RankElecUpdateMessageReq{
			AVID:    avID,
			UPMID:   upMID,
			PayMID:  payMID,
			Ver:     ver,
			Message: message,
			Hidden:  hidden,
		}
	)
	_, err = d.ugcPayRankAPI.RankElecUpdateMessage(ctx, req)
	return
}
