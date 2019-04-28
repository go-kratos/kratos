package v1

import (
	"context"
	"strconv"

	"google.golang.org/grpc/status"

	v1pb "go-common/app/admin/live/live-admin/api/http/v1"
	"go-common/app/admin/live/live-admin/conf"
	client "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/library/ecode"
)

// GaeaService struct
type GaeaService struct {
	conf   *conf.Config
	client *client.Client
	// optionally add other properties here, such as dao
	// dao *dao.Dao
}

//NewGaeaService init
func NewGaeaService(c *conf.Config) (s *GaeaService) {
	s = &GaeaService{
		conf: c,
	}
	s.client, _ = client.NewClient(c.ResourceClient)
	return s
}

// GetConfigByKeyword implementation
// 获取team下某个keyword的配置 `internal:"true"`
func (s *GaeaService) GetConfigByKeyword(ctx context.Context, req *v1pb.GetConfigReq) (resp *v1pb.GetConfigResp, err error) {
	resp = &v1pb.GetConfigResp{}
	if "" == req.GetKeyword() || 0 == req.GetTeam() {
		err = ecode.Error(1, "参数错误")
		return
	}
	ret, err := s.client.GetConfigByKeyword(ctx, &client.GetConfigReq{
		Team:    req.GetTeam(),
		Keyword: req.GetKeyword(),
	})
	if err != nil {
		return
	}
	resp.Team = ret.Team
	resp.Keyword = ret.Keyword
	resp.Name = ret.Name
	resp.Value = ret.Value
	resp.Ctime = ret.Ctime
	resp.Mtime = ret.Mtime
	resp.Status = ret.Status
	resp.Id = ret.Id
	return
}

// SetConfigByKeyword implementation
// 设置team下某个keyword配置 `internal:"true"`
func (s *GaeaService) SetConfigByKeyword(ctx context.Context, req *v1pb.SetConfigReq) (resp *v1pb.SetConfigResp, err error) {
	resp = &v1pb.SetConfigResp{}
	if "" == req.GetKeyword() || len(req.GetKeyword()) > 16 {
		err = ecode.Error(1, "参数错误")
		return
	}
	ret, err := s.client.SetConfigByKeyword(ctx, &client.SetConfigReq{
		Team:    req.GetTeam(),
		Keyword: req.GetKeyword(),
		Value:   req.GetValue(),
		Name:    req.GetName(),
		Id:      req.GetId(),
		Status:  req.GetStatus(),
	})
	if err != nil {
		return
	}
	resp.Id = ret.Id
	return
}

// GetConfigsByParams implementation
// 管理后台根据条件获取配置 `internal:"true"`
func (s *GaeaService) GetConfigsByParams(ctx context.Context, req *v1pb.ParamsConfigReq) (resp *v1pb.ParamsConfigResp, err error) {
	resp = &v1pb.ParamsConfigResp{}
	clientResp, err := s.client.GetConfigsByParams(ctx, &client.ParamsConfigReq{
		Team:     req.GetTeam(),
		Keyword:  req.GetKeyword(),
		Name:     req.GetName(),
		Status:   req.GetStatus(),
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
		Id:       req.GetId(),
	})
	resp.TotalNum = clientResp.TotalNum
	resp.List = []*v1pb.List{}
	for _, v := range clientResp.List {
		detail := &v1pb.List{
			Id:      v.Id,
			Team:    v.Team,
			Keyword: v.Keyword,
			Name:    v.Name,
			Value:   v.Value,
			Ctime:   v.Ctime,
			Mtime:   v.Mtime,
			Status:  v.Status,
		}
		resp.List = append(resp.List, detail)
	}

	if err != nil {
		return
	}
	return
}

// FormatErr format error msg
func (s *GaeaService) FormatErr(statusCode *status.Status) (code int32, msg string) {
	gCode := statusCode.Code()
	code = 1
	if gCode == 2 {
		code, _ := strconv.Atoi(statusCode.Message())

		switch code {
		case 1:
			msg = "必要参数不正确"
		case 11:
			msg = "索引名称在分组内冲突"
		case -500:
			msg = "内部错误"
		default:
			msg = "内部错误"
		}
	} else {
		msg = "内部错误"
	}
	return
}

// GetConfigsByTeam implementation
// 获取单个team的全部配置 `internal:"true"`
func (s *GaeaService) GetConfigsByTeam(ctx context.Context, req *v1pb.TeamConfigReq) (resp *v1pb.TeamConfigResp, err error) {
	resp = &v1pb.TeamConfigResp{}
	return
}

// GetConfigsByKeyword implementation
// 通过keyword获取配置 `internal:"true"`
func (s *GaeaService) GetConfigsByKeyword(ctx context.Context, req *v1pb.GetConfigsReq) (resp *v1pb.GetConfigsResp, err error) {
	resp = &v1pb.GetConfigsResp{}
	return
}

// GetConfigsByTeams implementation
// 获取多个team下的全部配置 `internal:"true"`
func (s *GaeaService) GetConfigsByTeams(ctx context.Context, req *v1pb.TeamsConfigReq) (resp *v1pb.TeamsConfigResp, err error) {
	resp = &v1pb.TeamsConfigResp{}
	return
}
