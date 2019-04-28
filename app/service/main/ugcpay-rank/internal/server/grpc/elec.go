package grpc

import (
	"context"

	"go-common/app/service/main/ugcpay-rank/api/v1"
	"go-common/app/service/main/ugcpay-rank/internal/conf"
	"go-common/app/service/main/ugcpay-rank/internal/model"
	"go-common/app/service/main/ugcpay-rank/internal/service"
	"go-common/app/service/main/ugcpay-rank/internal/service/rank"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"

	"github.com/golang/protobuf/ptypes/empty"
)

// New Identify warden rpc server
func New(cfg *warden.ServerConfig, s *service.Service) *warden.Server {
	w := warden.NewServer(cfg)
	v1.RegisterUGCPayRankServer(w.Server(), &UGCRankServer{s})
	ws, err := w.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

// UGCRankServer .
type UGCRankServer struct {
	svr *service.Service
}

var _ v1.UGCPayRankServer = &UGCRankServer{}

// RankElecAllAV .
func (u *UGCRankServer) RankElecAllAV(ctx context.Context, req *v1.RankElecAVReq) (resp *v1.RankElecAVResp, err error) {
	if req.AVID <= 0 {
		return
	}
	if req.RankSize <= 0 || req.RankSize > conf.Conf.Biz.ElecAVRankSize {
		req.RankSize = conf.Conf.Biz.ElecAVRankSize
	}
	r, err := u.svr.ElecTotalRankAV(ctx, req.UPMID, req.AVID, req.RankSize)
	if err != nil {
		return
	}
	resp = &v1.RankElecAVResp{
		AV: r,
	}
	return
}

// RankElecMonthAV .
func (u *UGCRankServer) RankElecMonthAV(ctx context.Context, req *v1.RankElecAVReq) (resp *v1.RankElecAVResp, err error) {
	if req.AVID <= 0 {
		return
	}
	if req.RankSize <= 0 || req.RankSize > conf.Conf.Biz.ElecAVRankSize {
		req.RankSize = conf.Conf.Biz.ElecAVRankSize
	}
	r, err := u.svr.ElecMonthRankAV(ctx, req.UPMID, req.AVID, req.RankSize)
	if err != nil {
		return
	}
	resp = &v1.RankElecAVResp{
		AV: r,
	}
	return
}

// RankElecMonthUP .
func (u *UGCRankServer) RankElecMonthUP(ctx context.Context, req *v1.RankElecUPReq) (resp *v1.RankElecUPResp, err error) {
	if req.UPMID <= 0 {
		return
	}
	if req.RankSize <= 0 || req.RankSize > conf.Conf.Biz.ElecUPRankSize {
		req.RankSize = conf.Conf.Biz.ElecUPRankSize
	}
	r, err := u.svr.ElecMonthRankUP(ctx, req.UPMID, req.RankSize)
	if err != nil {
		return
	}
	resp = &v1.RankElecUPResp{
		UP: r,
	}
	return
}

// RankElecMonth .
func (u *UGCRankServer) RankElecMonth(ctx context.Context, req *v1.RankElecMonthReq) (resp *v1.RankElecMonthResp, err error) {
	var (
		up *model.RankElecUPProto
		av *model.RankElecAVProto
	)
	if req.RankSize <= 0 || req.RankSize > conf.Conf.Biz.ElecAVRankSize {
		req.RankSize = conf.Conf.Biz.ElecAVRankSize
	}
	if req.UPMID > 0 {
		if up, err = u.svr.ElecMonthRankUP(ctx, req.UPMID, req.RankSize); err != nil {
			return
		}
	}
	if req.AVID > 0 {
		if av, err = u.svr.ElecMonthRankAV(ctx, req.UPMID, req.AVID, req.RankSize); err != nil {
			return
		}
	}
	resp = &v1.RankElecMonthResp{
		UP: up,
		AV: av,
	}
	return
}

// RankElecUpdateOrder .
func (u *UGCRankServer) RankElecUpdateOrder(ctx context.Context, req *v1.RankElecUpdateOrderReq) (reply *empty.Empty, err error) {
	reply = &empty.Empty{}
	var (
		prs                   = make([]rank.PrepRank, 0)
		elecMonthlyPrepUPRank = rank.NewElecPrepUPRank(req.UPMID, conf.Conf.Biz.ElecUPRankSize, req.Ver, u.svr.Dao)
		// elecTotalPrepUPRank   = rank.NewElecPrepUPRank(req.UPMID, conf.Conf.Biz.ElecUPRankSize, 0, u.svr.Dao) 因为敖厂长，所以不能及时更新up总榜，否则很容易db超时
	)
	prs = append(prs, elecMonthlyPrepUPRank)
	if req.AVID != 0 {
		elecMonthlyPrepAVRank := rank.NewElecPrepAVRank(req.AVID, conf.Conf.Biz.ElecAVRankSize, req.Ver, u.svr.Dao)
		elecTotalPrepAVRank := rank.NewElecPrepAVRank(req.AVID, conf.Conf.Biz.ElecAVRankSize, 0, u.svr.Dao)
		prs = append(prs, elecMonthlyPrepAVRank, elecTotalPrepAVRank)
	}
	// 更新缓存Prep_Rank
	u.svr.Asyncer.Do(ctx, func(ctx context.Context) {
		for _, pr := range prs {
			if theErr := u.svr.UpdateElecPrepRankFromOrder(ctx, pr, req.PayMID, req.Fee); theErr != nil {
				log.Error("update prep_rank failed, pr: %s, payMID: %d, fee: %d, err: %+v", pr, req.PayMID, req.Fee, theErr)
				continue
			}
		}
	})
	return
}

// RankElecUpdateMessage .
func (u *UGCRankServer) RankElecUpdateMessage(ctx context.Context, req *v1.RankElecUpdateMessageReq) (reply *empty.Empty, err error) {
	reply = &empty.Empty{}
	// 更新缓存Prep_Rank
	var (
		prs = make([]rank.PrepRank, 0)
	)
	if req.AVID != 0 {
		prs = append(prs, rank.NewElecPrepAVRank(req.AVID, conf.Conf.Biz.ElecAVRankSize, req.Ver, u.svr.Dao))
		prs = append(prs, rank.NewElecPrepAVRank(req.AVID, conf.Conf.Biz.ElecAVRankSize, 0, u.svr.Dao))
	}
	if req.UPMID != 0 {
		prs = append(prs, rank.NewElecPrepUPRank(req.UPMID, conf.Conf.Biz.ElecUPRankSize, req.Ver, u.svr.Dao))
	}
	u.svr.Asyncer.Do(ctx, func(ctx context.Context) {
		for _, pr := range prs {
			if theErr := u.svr.UpdateElecPrepRankFromMessage(ctx, pr, req.PayMID, req.Message, req.Hidden); theErr != nil {
				log.Error("update prep_rank failed, pr: %s, payMID: %d, message: %s, hidden: %t, err: %+v", pr, req.PayMID, req.Message, req.Hidden, theErr)
				continue
			}
		}
	})
	return
}
