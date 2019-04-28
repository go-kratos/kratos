package service

import (
	"context"
	"fmt"
	v1 "go-common/app/service/bbq/video/api/grpc/v1"
	httpV1 "go-common/app/service/bbq/video/api/http/v1"
	"go-common/app/service/bbq/video/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
)

// ModifyLimits .
func (s *Service) ModifyLimits(c context.Context, args *v1.ModifyLimitsRequest) (response *v1.ModifyLimitsResponse, err error) {
	response = new(v1.ModifyLimitsResponse)
	_, err = s.dao.ModifyLimits(c, args.Svid, args.LimitType, args.LimitOp)
	if err != nil {
		log.Warnw(c, "log", "modify limits fail", "args", args.String())
		return
	}
	return
}

// ListVideoInfo 视频信息列表.
func (s *Service) ListVideoInfo(ctx context.Context, v *v1.ListVideoInfoRequest) (res *v1.ListVideoInfoResponse, err error) {
	res = new(v1.ListVideoInfoResponse)
	videoBases, err := s.dao.VideoBase(ctx, v.SvIDs)
	if err != nil {
		log.Errorw(ctx, "log", "batch get video base fail")
		return
	}
	for _, videoBase := range videoBases {
		res.List = append(res.List, &v1.VideoInfo{VideoBase: videoBase})
	}
	return
}

// ImportVideo 导入视频服务.
func (s *Service) ImportVideo(ctx context.Context, v *v1.ImportVideoInfo) (res *empty.Empty, err error) {
	res = &empty.Empty{}
	err = s.dao.AddOrUpdateVideo(ctx, v)
	return
}

// SyncTag 同步标签.
func (s *Service) SyncTag(ctx context.Context, v *v1.SyncVideoTagRequest) (res *empty.Empty, err error) {
	res = &empty.Empty{}
	_, err = s.dao.AddOrUpdateTag(ctx, v.TagInfos)
	return
}

//SvStatisticsInfo ...
func (s *Service) SvStatisticsInfo(ctx context.Context, v *v1.SvStatisticsInfoReq) (res *v1.SvStatisticsInfoRes, err error) {
	res = new(v1.SvStatisticsInfoRes)
	res.SvstInfoMap = make(map[int64]*v1.SvStInfo)
	if v.SvidList == nil || len(v.SvidList) <= 0 {
		return
	}
	var data map[int64]*model.SvStInfo
	data, err = s.dao.RawVideoStatistic(ctx, v.SvidList)
	for svid, st := range data {
		res.SvstInfoMap[svid] = &v1.SvStInfo{
			Like:      st.Like,
			Play:      st.Play,
			Report:    st.Report,
			Share:     st.Share,
			Subtitles: st.Subtitles,
			Reply:     st.Reply,
			Svid:      svid,
		}
	}
	return
}

//SyncUserBase 更新userbase
func (s *Service) SyncUserBase(ctx context.Context, req *v1.SyncMidRequset) (res *v1.SyncUserBaseResponse, err error) {
	res, err = s.dao.InOrUpUserBase(ctx, req.MID)
	return
}

//SyncUserSta 更新user_statistics_hive
func (s *Service) SyncUserSta(ctx context.Context, req *v1.SyncMidRequset) (res *v1.SyncUserBaseResponse, err error) {
	res, err = s.dao.InOrUpUserSta(ctx, req.MID)
	return
}

//SyncUserBases 批量更新userbase
func (s *Service) SyncUserBases(ctx context.Context, req *v1.SyncMidsRequset) (res *v1.SyncUserBaseResponse, err error) {
	res, err = s.dao.InOrUpUserBases(ctx, req.MIDS)
	return
}

//SyncUserStas 批量更新user_statistics_hive
func (s *Service) SyncUserStas(ctx context.Context, req *v1.SyncMidsRequset) (res *v1.SyncUserBaseResponse, err error) {
	res, err = s.dao.InOrUpUserStas(ctx, req.MIDS)
	return
}

// BVCTransRes 处理BVC回调服务
func (s *Service) BVCTransRes(ctx context.Context, req *v1.BVCTransBackRequset) (err error) {
	//记录回调日志
	err = s.dao.AddOrUpdateFlowRecord(ctx, &model.BVCRecord{
		FLowID: req.FlowID,
		SVID:   req.SVID,
		Type:   req.FlowType,
	})
	if err != nil {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("AddFlowRecord err[%v]", err)))
		err = ecode.SyncBVCFail
		return
	}
	//失败回调直接返回
	if req.FlowType < 0 {
		log.Warn("BVCTrans fail [%+v]", req)
		return nil
	}
	tx, err := s.dao.BeginTran(ctx)
	if err != nil {
		return
	}
	if req.PIC.PicURL == "" || req.PIC.PicHeight == 0 || req.PIC.PicWidth == 0 {
		log.Warn("图片参数缺失 pic[%+v]", req.PIC)
		err = ecode.ReqParamErr
		return
	}
	err = s.dao.UpdateCmsSvPIC(ctx, req.SVID, req.PIC, model.SourceXcodeCover)
	if err != nil {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("UpdateSvPIC err[%v]", err)))
		return
	}
	stratAt := time.Now()
	wg := &errgroup.Group{}
	for i := range req.TransRes {
		wg.Go(func(i int) func() error {
			return func() error {
				d := req.TransRes[i]
				var bcode int64
				if _, ok := s.c.BPSCode[d.PPI]; ok {
					if codeRate, ok := s.c.BPSCode[d.PPI][d.BPS]; ok {
						bcode = int64(codeRate)
					}
				}
				if bcode == 0 {
					return ecode.UnKnownBPS
				}
				data := &model.VideoBVC{
					SVID:            req.SVID,
					Path:            d.Path,
					CodeRate:        bcode,
					ResolutionRetio: d.PPI,
					VideoCode:       d.VideoCode,
					Duration:        d.Duration,
					FileSize:        d.Filesize,
				}
				err = s.dao.TxAddOrUpdateBVCInfo(ctx, tx, data)
				if err != nil {
					log.Errorv(ctx, log.KV("log", fmt.Sprintf("TxAddOrUpdateBVCInfo err[%v]", err)))
				}
				return err
			}
		}(i))
	}
	//统一回滚
	if err = wg.Wait(); err != nil {
		tx.Rollback()
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("ErrorGroup err[%v]", err)))
		return
	}
	//统一提交事务
	if err = tx.Commit(); err != nil {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("Tx err[%v]", err)))
		err = ecode.SyncBVCFail
		return
	}
	elapsed := time.Since(stratAt)
	log.Infov(ctx,
		log.KV("log", fmt.Sprintf("BVCTransRes Sync Complete, cost[%s]", elapsed)))
	err = s.dao.CmsPub(ctx, &model.DataTopicCmsData{SVID: req.SVID})
	if err != nil {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("CmsPub Fail:[%v]", err)))
	}
	return
}

// BVCTransCommit 提交转码
func (s *Service) BVCTransCommit(ctx context.Context, req *v1.BVideoTransRequset) (*empty.Empty, error) {
	err := s.dao.CommitTrans(ctx, req)
	if err != nil {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("CommitTrans err[%v] data[%+v]", err, req)))
	}
	return &empty.Empty{}, err
}

//CheckVideoUploadSt upload video to client
func (s *Service) CheckVideoUploadSt(c context.Context, SVID int64) (err error) {
	err = s.dao.CheckSVResource(c, SVID)
	return
}

// CreateID 创建新ID
// 按十进制计算
// 63个1=9223372036854775807，共19位
// 时间戳取32个1=4294967295，共10位，且最高位4小于9
// mid%1000，共3位
// 机器标志3位
// 自增2位
// 保留1位
func (s *Service) CreateID(ctx context.Context, req *v1.CreateIDRequest) (res *v1.CreateIDResponse, err error) {
	res = new(v1.CreateIDResponse)

	mid := req.Mid
	if mid == 0 {
		err = ecode.ReqParamErr
		return
	}
	timestamp := time.Now().Unix()
	if req.Time > 0xFFFFFFFF {
		err = ecode.ReqParamErr
		return
	}
	if req.Time != 0 {
		timestamp = req.Time
	}
	// 十进制，19位的分配情况
	// 10   |1       |3   |2    |3
	// ts   |reserved|host|index|mid
	increaseID := atomic.AddInt64(&autoIncreaseID, 1)
	newIndex := increaseID % 100
	midInfo := mid % 1000

	res.NewId = (timestamp * 1000000000) + (hostHashIndex * 100000) + (newIndex * 1000) + midInfo
	log.Infov(ctx, log.KV("log",
		fmt.Sprintf("create one new id: new_id=%d, timestamp=%d, increase_id=%d, mid=%d",
			res.NewId, timestamp, increaseID, mid)))
	return
}

// VideoViewsAdd .
func (s *Service) VideoViewsAdd(c context.Context, args *httpV1.ViewsAddRequest) (response *httpV1.ViewsAddResponse, err error) {
	response = new(httpV1.ViewsAddResponse)
	affected, err := s.dao.AddVideoViews(c, args.Svid, args.Views)
	response.Affected = affected
	return
}

// VideoUnshelf .
func (s *Service) VideoUnshelf(ctx context.Context, in *v1.VideoUnshelfRequest) (*empty.Empty, error) {
	newState := model.VideoUnshelf
	if _, err := s.dao.VideoStateUpdate(ctx, in.Svid, newState); err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

// VideoDelete .
func (s *Service) VideoDelete(ctx context.Context, in *v1.VideoDeleteRequest) (*empty.Empty, error) {
	videoBase, err := s.dao.RawVideoBase(ctx, []int64{in.Svid})
	if err != nil {
		return nil, err
	}

	if v, ok := videoBase[in.Svid]; !ok || v.Mid != in.UpMid {
		return nil, ecode.VideoDelFail
	}

	newState := model.VideoDelete
	if _, err := s.dao.VideoStateUpdate(ctx, in.Svid, newState); err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}
