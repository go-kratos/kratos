package dao

import (
	"context"

	api "go-common/app/service/main/ugcpay-rank/api/v1"
)

// RankElecMonth .
func (d *Dao) RankElecMonth(ctx context.Context, avID, upMID int64, rankSize int) (resp *api.RankElecMonthResp, err error) {
	var (
		req = &api.RankElecMonthReq{
			AVID:     avID,
			UPMID:    upMID,
			RankSize: rankSize,
		}
	)
	return d.ugcPayRankAPI.RankElecMonth(ctx, req)
}

// RankElecMonthUP .
func (d *Dao) RankElecMonthUP(ctx context.Context, upMID int64, rankSize int) (rank *api.RankElecUPResp, err error) {
	var (
		req = &api.RankElecUPReq{
			UPMID:    upMID,
			RankSize: rankSize,
		}
	)
	return d.ugcPayRankAPI.RankElecMonthUP(ctx, req)
}

// RankElecAllAV .
func (d *Dao) RankElecAllAV(ctx context.Context, upMID int64, avID int64, rankSize int) (rank *api.RankElecAVResp, err error) {
	var (
		req = &api.RankElecAVReq{
			UPMID:    upMID,
			AVID:     avID,
			RankSize: rankSize,
		}
	)
	return d.ugcPayRankAPI.RankElecAllAV(ctx, req)
}
