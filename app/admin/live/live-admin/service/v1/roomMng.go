package v1

import (
	"context"
	"fmt"
	v1pb "go-common/app/admin/live/live-admin/api/http/v1"
	"go-common/app/admin/live/live-admin/conf"
	"go-common/app/admin/live/live-admin/dao"
	relationV1 "go-common/app/service/live/relation/api/liverpc/v1"
	liveRPCCli "go-common/app/service/live/room/api/liverpc"
	v1liveRPCpb "go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"os"
	"time"

	streamPb "go-common/app/service/video/stream-mng/api/v1"
	"go-common/library/ecode"
)

const timeFormat = "2006-01-02 15:04:05"

// RoomMngService struct
type RoomMngService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	rpcCli    *liveRPCCli.Client
	streamCli streamPb.StreamClient
	dao       *dao.Dao
}

//NewRoomMngService init
func NewRoomMngService(c *conf.Config) (s *RoomMngService) {
	s = &RoomMngService{
		conf: c,
		dao:  dao.New(c),
	}

	log.Info("RoomMngLiveRPCClient Init: %+v", s.conf.RoomMngClient)
	s.rpcCli = liveRPCCli.New(s.conf.RoomMngClient)

	log.Info("Stream Mng Client Init: %+v", s.conf.StreamMngClient)
	if svc, err := streamPb.NewClient(s.conf.StreamMngClient); err != nil {
		panic(err)
	} else {
		s.streamCli = svc
	}

	return s
}

// GetSecondVerifyListWithPics implementation
// 获取带有图片地址的二次审核列表
// `method:"GET" internal:"true" `
func (s *RoomMngService) GetSecondVerifyListWithPics(ctx context.Context, req *v1pb.RoomMngGetSecondVerifyListReq) (resp *v1pb.RoomMngGetSecondVerifyListResp, err error) {
	if req.Pagesize == 0 {
		req.Pagesize = 30
	}
	rpcResp, err := s.rpcCli.V1RoomMng.GetSecondVerifyList(ctx, &v1liveRPCpb.RoomMngGetSecondVerifyListReq{
		RoomId:   req.RoomId,
		Area:     req.Area,
		Page:     req.Page,
		Pagesize: req.Pagesize,
		Biz:      req.Biz,
	})

	if err != nil {
		return
	}
	if rpcResp.Code == 0 {
		result, getPicErr := s.getExtraInfo(ctx, rpcResp.Data.Result)
		if getPicErr != nil {
			err = getPicErr
		} else {
			resp = &v1pb.RoomMngGetSecondVerifyListResp{
				Count:    rpcResp.Data.Count,
				Page:     rpcResp.Data.Page,
				Pagesize: rpcResp.Data.Pagesize,
				Result:   result,
			}
		}
	} else {
		err = ecode.Error(ecode.ServerErr, fmt.Sprintf("Room v1 RoomMng LiveRPC 业务错误: { Code: %d ; Msg: %s }", rpcResp.Code, rpcResp.Msg))
	}
	return
}

func (s *RoomMngService) getExtraInfo(ctx context.Context, list []*v1liveRPCpb.RoomMngGetSecondVerifyListResp_Result) (respList []*v1pb.RoomMngGetSecondVerifyListResp_Result, err error) {
	picRes := make([][]string, len(list))
	fcRes := make([]int64, len(list))
	g, _ := errgroup.WithContext(ctx)

	for i := 0; i < len(list); i++ {
		x := i

		// 获取指定时间内的截图
		g.Go(func() error {
			picRes[x] = s.getSingleRoomPic(ctx, list[x].RoomId, list[x].BreakTime)
			return nil
		})

		// 获取粉丝计数
		g.Go(func() (feedErr error) {
			fcRes[x] = s.getUserFC(ctx, list[x].Uid)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Error("get ExtraData Error (%+v)", err)
	}

	respList = make([]*v1pb.RoomMngGetSecondVerifyListResp_Result, len(list))
	for index, RPCRespItem := range list {
		// 当截图张数大于 0 且有证据图片时，将截图第一张替换为证据图片
		if len(picRes[index]) > 0 && RPCRespItem.ProofImg != "" {
			picRes[index][0] = RPCRespItem.ProofImg
		}
		respList[index] = &v1pb.RoomMngGetSecondVerifyListResp_Result{
			Id:              RPCRespItem.Id,
			RecentCutTimes:  RPCRespItem.RecentCutTimes,
			RecentWarnTimes: RPCRespItem.RecentWarnTimes,
			Uname:           RPCRespItem.Uname,
			RoomId:          RPCRespItem.RoomId,
			Uid:             RPCRespItem.Uid,
			Title:           RPCRespItem.Title,
			AreaV2Name:      RPCRespItem.AreaV2Name,
			Fc:              fcRes[index],
			WarnReason:      RPCRespItem.WarnReason,
			BreakTime:       RPCRespItem.BreakTime,
			Pics:            picRes[index],
			WarnTimes:       RPCRespItem.WarnTimes,
		}
	}
	return
}

func (s *RoomMngService) getSingleRoomPic(ctx context.Context, roomID int64, breakTime string) (picList []string) {
	// uat 环境写死 11891462 房间号
	if os.Getenv("DEPLOY_ENV") == "uat" {
		roomID = 11891462
	}
	// 计算结束时间点，规定为 5 分钟后
	startTime, _ := time.Parse(timeFormat, breakTime)
	endTimeStr := startTime.Add(5 * 60 * 1000 * 1000 * 1000).Format(timeFormat)

	RPCResp, err := s.streamCli.GetSingleScreeShot(ctx, &streamPb.GetSingleScreeShotReq{
		RoomId:    roomID,
		StartTime: breakTime,
		EndTime:   endTimeStr,
	})
	// log.Info("res: (%+v) err:(%+v) \n", RPCResp, err)
	if err != nil {
		log.Error("Get Pic Fail: error(%v) roomId(%d) startTime(%s) endTime(%s)", err, roomID, breakTime, endTimeStr)
		picList = []string{
			"https://static.hdslb.com/error/very_sorry.png",
		}
	} else {
		picList = RPCResp.List
	}

	return
}

func (s *RoomMngService) getUserFC(ctx context.Context, UID int64) (resp int64) {
	feedResp, feedErr := s.dao.Relation.V1Feed.GetUserFc(ctx, &relationV1.FeedGetUserFcReq{Follow: UID})
	if feedErr != nil || feedResp.Data == nil {
		resp = 0
		log.Error("Get FC Error for UID(%d) with Error(%+v)", UID, feedErr)
	} else {
		resp = feedResp.Data.Fc
	}
	return
}
