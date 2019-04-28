package v1

import (
	"context"

	v1pb "go-common/app/interface/live/web-ucenter/api/http/v1"
	"go-common/app/interface/live/web-ucenter/conf"
	"go-common/app/service/live/xrewardcenter/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// AnchorTaskService struct
type AnchorTaskService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	conn v1.AnchorRewardClient
}

//NewAnchorTaskService init
func NewAnchorTaskService(c *conf.Config) (s *AnchorTaskService) {
	s = &AnchorTaskService{
		conf: c,
	}
	conn, err := v1.NewClient(conf.Conf.Warden)
	if err != nil {
		panic(err)
	}
	s.conn = conn

	return s
}

// MyReward implementation
// * (主播侧)-我的主播奖励(登录态)
func (s *AnchorTaskService) MyReward(ctx context.Context, req *v1pb.AnchorTaskMyRewardReq) (resp *v1pb.AnchorTaskMyRewardResp, err error) {
	resp = &v1pb.AnchorTaskMyRewardResp{}

	mid, _ := metadata.Value(ctx, "mid").(int64)

	if mid <= 0 {
		err = ecode.NoLogin
		return
	}

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	ret, err := s.conn.MyReward(ctx, &v1.AnchorTaskMyRewardReq{
		Page: page,
		Uid:  mid,
	})

	if err != nil {
		return
	}

	resp.Page = &v1pb.AnchorTaskMyRewardResp_Page{
		Page:       ret.GetPage().GetPage(),
		PageSize:   ret.GetPage().GetPageSize(),
		TotalPage:  ret.GetPage().GetTotalPage(),
		TotalCount: ret.GetPage().GetTotalCount(),
	}
	resp.ExpireCount = ret.GetExpireCount()

	for _, v := range ret.Data {
		resp.Data = append(resp.Data, &v1pb.AnchorTaskMyRewardResp_RewardObj{
			Id:          v.GetId(),
			RewardType:  v.GetRewardType(),
			Status:      v.GetStatus(),
			RewardId:    v.GetRewardId(),
			Name:        v.GetName(),
			Icon:        v.GetIcon(),
			AchieveTime: v.GetAchieveTime(),
			ExpireTime:  v.GetExpireTime(),
			Source:      v.GetSource(),
			RewardIntro: v.GetRewardIntro(),
		})
	}

	return
}

// UseRecord implementation
// (主播侧)-奖励使用记录(登录态)
// `midware:"auth"`
func (s *AnchorTaskService) UseRecord(ctx context.Context, req *v1pb.AnchorTaskUseRecordReq) (resp *v1pb.AnchorTaskUseRecordResp, err error) {
	resp = &v1pb.AnchorTaskUseRecordResp{}

	mid, _ := metadata.Value(ctx, "mid").(int64)

	if mid <= 0 {
		err = ecode.NoLogin
		return
	}

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	ret, err := s.conn.UseRecord(ctx, &v1.AnchorTaskUseRecordReq{
		Page: page,
		Uid:  mid,
	})

	if err != nil {
		return
	}

	resp.Page = &v1pb.AnchorTaskUseRecordResp_Page{
		Page:       ret.GetPage().GetPage(),
		PageSize:   ret.GetPage().GetPageSize(),
		TotalPage:  ret.GetPage().GetTotalPage(),
		TotalCount: ret.GetPage().GetTotalCount(),
	}
	for _, v := range ret.Data {
		resp.Data = append(resp.Data, &v1pb.AnchorTaskUseRecordResp_RewardObj{
			Id:          v.GetId(),
			RewardId:    v.GetRewardId(),
			Status:      v.GetStatus(),
			Name:        v.GetName(),
			Icon:        v.GetIcon(),
			AchieveTime: v.GetAchieveTime(),
			ExpireTime:  v.GetExpireTime(),
			Source:      v.GetSource(),
			RewardIntro: v.GetRewardIntro(),
			UseTime:     v.GetUseTime(),
		})
	}
	return
}

// UseReward implementation
// (主播侧)-使用奖励(登录态)
// `method:"POST" midware:"auth"`
func (s *AnchorTaskService) UseReward(ctx context.Context, req *v1pb.AnchorTaskUseRewardReq) (resp *v1pb.AnchorTaskUseRewardResp, err error) {

	resp = &v1pb.AnchorTaskUseRewardResp{}

	mid, _ := metadata.Value(ctx, "mid").(int64)

	if mid <= 0 {
		err = ecode.NoLogin
		return
	}

	id := req.GetId()

	if id <= 0 {
		err = ecode.ParamInvalid
		return
	}

	platform := req.GetPlatform()
	if "" == platform {
		platform = "web"
	}

	request := &v1.AnchorTaskUseRewardReq{
		Id:      id,
		Uid:     mid,
		UsePlat: platform,
	}
	ret, err := s.conn.UseReward(ctx, request)
	log.Info("useReward req(%v) ret(%v), err(%v)", request, ret, err)

	if err != nil {

		statusCode := ecode.Cause(err)
		log.Info("useReward error statusCode(%v) ret(%v), err(%+v)", statusCode, statusCode.Code(), err)

		busCode := statusCode.Code()
		msg := ""

		switch busCode {
		case 1:
			msg = "参数错误"
		case 2:
			msg = "这个奖励已经过期了呢"
		case 3:
			msg = "这个奖励已经被你使用啦~"
		case 4:
			msg = "为了更好的使用体验，请在开播状态下使用【任意门】哦"
		case 5:
			msg = "奖励不存在"
		default:
			msg = "内部错误"
		}

		err = ecode.Error(ecode.Code(busCode), msg)
		return
	}

	resp.Result = ret.GetResult()

	return
}

// IsViewed implementation
// (主播侧)-奖励和任务红点(登录态)
// `midware:"auth"`
func (s *AnchorTaskService) IsViewed(ctx context.Context, req *v1pb.AnchorTaskIsViewedReq) (resp *v1pb.AnchorTaskIsViewedResp, err error) {
	resp = &v1pb.AnchorTaskIsViewedResp{}

	mid, _ := metadata.Value(ctx, "mid").(int64)

	if mid <= 0 {
		err = ecode.NoLogin
		return
	}

	ret, err := s.conn.IsViewed(ctx, &v1.AnchorTaskIsViewedReq{
		Uid: mid,
	})
	log.Info("IsViewed req(%v) ret(%v), err(%v)", mid, ret, err)

	if err != nil {
		return
	}

	resp = &v1pb.AnchorTaskIsViewedResp{
		TaskShouldNotice:   ret.GetTaskShouldNotice(),
		ShowRewardEntry:    ret.GetShowRewardEntry(),
		RewardShouldNotice: ret.GetRewardShouldNotice(),
		TaskStatus:         ret.GetTaskStatus(),
		IsBlacked:          ret.GetIsBlacked(),
		Url:                ret.GetUrl(),
	}

	return
}

// AddReward implementation
// (主播侧)-添加主播奖励(内部接口)
// `method:"POST" internal:"true"`
func (s *AnchorTaskService) AddReward(ctx context.Context, req *v1pb.AnchorTaskAddRewardReq) (resp *v1pb.AnchorTaskAddRewardResp, err error) {
	resp = &v1pb.AnchorTaskAddRewardResp{}

	ret, err := s.conn.AddReward(ctx, &v1.AnchorTaskAddRewardReq{
		RewardId: req.GetRewardId(),
		Roomid:   req.GetRoomid(),
		Source:   req.GetSource(),
		Uid:      req.GetUid(),
		OrderId:  req.GetOrderId(),
	})

	if err != nil {
		return
	}

	resp = &v1pb.AnchorTaskAddRewardResp{
		Result: ret.GetResult(),
	}
	resp.Result = ret.GetResult()
	return
}
