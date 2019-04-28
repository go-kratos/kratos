package service

import (
	"context"

	apihttp "go-common/app/interface/main/ugcpay-rank/api/http"
	apirpc "go-common/app/service/main/ugcpay-rank/api/v1"
)

// RankElecMonthUP .
func (s *Service) RankElecMonthUP(ctx context.Context, upMID int64, rankSize int) (rank *apihttp.RespRankElecMonthUP, err error) {
	var (
		resp *apirpc.RankElecUPResp
	)
	if resp, err = s.dao.RankElecMonthUP(ctx, upMID, rankSize); err != nil {
		return
	}
	rank = &apihttp.RespRankElecMonthUP{}
	rank.Parse(resp.UP)
	return
}

// RankElecMonth .
func (s *Service) RankElecMonth(ctx context.Context, upMID, avID int64, rankSize int) (rank *apihttp.RespRankElecMonth, err error) {
	var (
		resp *apirpc.RankElecMonthResp
	)
	if resp, err = s.dao.RankElecMonth(ctx, avID, upMID, rankSize); err != nil {
		return
	}
	rank = &apihttp.RespRankElecMonth{}
	rank.Parse(resp.AV, resp.UP)
	return
}

// RankElecAllAV .
func (s *Service) RankElecAllAV(ctx context.Context, upMID int64, avID int64, rankSize int) (rank *apihttp.RespRankElecAllAV, err error) {
	var (
		resp *apirpc.RankElecAVResp
	)
	if resp, err = s.dao.RankElecAllAV(ctx, upMID, avID, rankSize); err != nil {
		return
	}
	rank = &apihttp.RespRankElecAllAV{}
	rank.Parse(resp.AV)
	return
}
