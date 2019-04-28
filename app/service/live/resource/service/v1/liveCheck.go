package v1

import (
	"context"
	"encoding/json"

	v1pb "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/app/service/live/resource/conf"
	"go-common/app/service/live/resource/dao"
	"go-common/library/ecode"
	"go-common/library/log"
)

// LiveCheckService struct
type LiveCheckService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewLiveCheckService init
func NewLiveCheckService(c *conf.Config) (s *LiveCheckService) {
	s = &LiveCheckService{
		conf: c,
		dao:  dao.New(c),
	}

	return s
}

// LiveCheck implementation
// 客户端获取能否直播接口
func (s *LiveCheckService) LiveCheck(ctx context.Context, req *v1pb.LiveCheckReq) (resp *v1pb.LiveCheckResp, err error) {
	resp = &v1pb.LiveCheckResp{}
	log.Info("[LiveCheck] platform is %s system is %s mobile is %s", req.Platform, req.System, req.Mobile)
	if req.Platform == "" || req.Mobile == "" || req.System == "" {
		err = ecode.ResourceParamErr
		return
	}
	resp.IsLive = s.dao.GetLiveCheck(ctx, req.Platform, req.System, req.Mobile)
	return
}

// GetLiveCheckList implementation
// 后台查询所有配置设备黑名单
func (s *LiveCheckService) GetLiveCheckList(ctx context.Context, req *v1pb.GetLiveCheckListReq) (resp *v1pb.GetLiveCheckListResp, err error) {
	resp = &v1pb.GetLiveCheckListResp{}
	value, err := s.dao.GetLiveCheckList(ctx)
	if err != nil {
		err = ecode.GetConfAdminErr
		return
	}
	if value == "" {
		err = ecode.GetConfAdminErr
		return
	}
	err = json.Unmarshal([]byte(value), resp)
	if err != nil {
		log.Error("[LiveCheck] admin get live check conf error by wrong format")
		err = ecode.GetConfAdminErr
		return
	}
	return
}

// AddLiveCheck implementation
// 后台添加能否直播设备黑名单
func (s *LiveCheckService) AddLiveCheck(ctx context.Context, req *v1pb.AddLiveCheckReq) (resp *v1pb.AddLiveCheckResp, err error) {
	resp = &v1pb.AddLiveCheckResp{}
	value := req.LiveCheck
	if value == "" {
		err = ecode.ResourceParamErr
		return
	}
	list := &v1pb.GetLiveCheckListResp{}
	err = json.Unmarshal([]byte(value), list)
	if err != nil {
		err = ecode.ResourceParamErr
		return
	}
	jsonStr, _ := json.Marshal(list)
	err = s.dao.SetLiveCheck(ctx, string(jsonStr))
	if err != nil {
		err = ecode.SetConfAdminErr
		return
	}
	return
}
