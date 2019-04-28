package v1

import (
	"context"
	"strings"

	v1pb "go-common/app/admin/live/live-admin/api/http/v1"
	"go-common/app/admin/live/live-admin/conf"
	rspb "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/library/log"
)

// ResourceService struct
type ResourceService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	rsCli *rspb.Client
}

//NewResourceService init
func NewResourceService(c *conf.Config) (s *ResourceService) {
	s = &ResourceService{
		conf: c,
	}
	var svc *rspb.Client
	var err error

	log.Info("ResourceService Init: %+v", s.conf.ResourceClient)
	if svc, err = rspb.NewClient(s.conf.ResourceClient); err != nil {
		panic(err)
	}
	s.rsCli = svc
	return s
}

// Add implementation
// Add 添加资源接口
// `method:"POST" internal:"true"
func (s *ResourceService) Add(ctx context.Context, req *v1pb.AddReq) (resp *v1pb.AddResp, err error) {
	respRPC, error := s.rsCli.Add(ctx, &rspb.AddReq{
		Platform:     req.Platform,
		Title:        req.Title,
		JumpPath:     req.JumpPath,
		JumpPathType: req.JumpPathType,
		JumpTime:     req.JumpTime,
		Type:         req.Type,
		Device:       req.Device,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		ImageUrl:     req.ImageUrl,
	})
	err = error
	if error == nil {
		resp = &v1pb.AddResp{
			Id: respRPC.Id,
		}
	}
	return
}

// AddEx implementation
// AddEx 添加资源接口
// `method:"POST" internal:"true"
func (s *ResourceService) AddEx(ctx context.Context, req *v1pb.AddReq) (resp *v1pb.AddResp, err error) {
	respRPC, error := s.rsCli.AddEx(ctx, &rspb.AddReq{
		Platform:     req.Platform,
		Title:        req.Title,
		JumpPath:     req.JumpPath,
		JumpPathType: req.JumpPathType,
		JumpTime:     req.JumpTime,
		Type:         req.Type,
		Device:       req.Device,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		ImageUrl:     req.ImageUrl,
	})
	err = error
	if error == nil {
		resp = &v1pb.AddResp{
			Id: respRPC.Id,
		}
	}
	return
}

// Edit implementation
// Edit 编辑资源接口
// `method:"POST" internal:"true" `
func (s *ResourceService) Edit(ctx context.Context, req *v1pb.EditReq) (resp *v1pb.EditResp, err error) {
	_, err = s.rsCli.Edit(ctx, &rspb.EditReq{
		Platform:     req.Platform,
		Id:           req.Id,
		Title:        req.Title,
		JumpPath:     req.JumpPath,
		JumpTime:     req.JumpTime,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		ImageUrl:     req.ImageUrl,
		JumpPathType: req.JumpPathType,
	})
	return
}

// Offline implementation
// Offline 下线资源接口
// `method:"POST" internal:"true" `
func (s *ResourceService) Offline(ctx context.Context, req *v1pb.OfflineReq) (resp *v1pb.OfflineResp, err error) {
	_, err = s.rsCli.Offline(ctx, &rspb.OfflineReq{
		Platform: req.Platform,
		Id:       req.Id,
	})
	if err == nil {
		resp = &v1pb.OfflineResp{}
	}
	return
}

// GetList implementation
// GetList 获取资源列表
// `method:"GET" internal:"true" `
func (s *ResourceService) GetList(ctx context.Context, req *v1pb.GetListReq) (resp *v1pb.GetListResp, err error) {
	RPCReq := &rspb.GetListReq{
		Platform: req.Platform,
		Type:     req.Type,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	if RPCReq.Page == 0 {
		RPCReq.Page = 1
	}
	if RPCReq.PageSize == 0 {
		RPCReq.PageSize = 50
	}
	var RPCResp *rspb.GetListResp
	RPCResp, err = s.rsCli.GetList(ctx, RPCReq)
	if err == nil {
		resp = &v1pb.GetListResp{
			CurrentPage: RPCResp.CurrentPage,
			TotalCount:  RPCResp.TotalCount,
			List:        convertRPCList2HttpList(RPCResp.List),
		}
	}
	return
}

// GetPlatformList implementation
// 获取平台列表
// `method:"GET" internal:"true" `
func (s *ResourceService) GetPlatformList(ctx context.Context, req *v1pb.GetPlatformListReq) (resp *v1pb.GetPlatformListResp, err error) {
	var RPCResp *rspb.GetPlatformListResp
	RPCResp, err = s.rsCli.GetPlatformList(ctx, &rspb.GetPlatformListReq{
		Type: req.Type,
	})
	if err == nil {
		resp = &v1pb.GetPlatformListResp{
			Platform: RPCResp.Platform,
		}
	}
	return
}

// GetListEx implementation
// GetListEx 获取资源列表
// `method:"GET" internal:"true" `
func (s *ResourceService) GetListEx(ctx context.Context, req *v1pb.GetListExReq) (resp *v1pb.GetListExResp, err error) {
	var RPCResp *rspb.GetListExResp
	// 默认分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 50
	}
	if len(req.Type) > 0 {
		req.Type = strings.Split(req.Type[0], ",")
	}
	RPCResp, err = s.rsCli.GetListEx(ctx, &rspb.GetListExReq{
		Platform:       req.Platform,
		Page:           req.Page,
		PageSize:       req.PageSize,
		Type:           req.Type,
		DevicePlatform: req.DevicePlatform,
		Status:         req.Status,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
	})
	if err == nil {
		resp = &v1pb.GetListExResp{
			CurrentPage: RPCResp.CurrentPage,
			TotalCount:  RPCResp.TotalCount,
			List:        convertRPCListEx2HttpListEx(RPCResp.List),
		}
	}
	return
}

func convertRPCList2HttpList(RPCList []*rspb.GetListResp_List) (HTTPList []*v1pb.GetListResp_List) {
	HTTPList = make([]*v1pb.GetListResp_List, len(RPCList))
	for index, RPCRespItem := range RPCList {
		HTTPList[index] = &v1pb.GetListResp_List{
			Id:             RPCRespItem.Id,
			Title:          RPCRespItem.Title,
			JumpPath:       RPCRespItem.JumpPath,
			DevicePlatform: RPCRespItem.DevicePlatform,
			DeviceBuild:    RPCRespItem.DeviceBuild,
			StartTime:      RPCRespItem.StartTime,
			EndTime:        RPCRespItem.EndTime,
			Status:         RPCRespItem.Status,
			DeviceLimit:    RPCRespItem.DeviceLimit,
			ImageUrl:       RPCRespItem.ImageUrl,
			JumpPathType:   RPCRespItem.JumpPathType,
			JumpTime:       RPCRespItem.JumpTime,
		}
	}
	return
}

func convertRPCListEx2HttpListEx(RPCList []*rspb.GetListExResp_List) (HTTPList []*v1pb.GetListExResp_List) {
	HTTPList = make([]*v1pb.GetListExResp_List, len(RPCList))
	for index, RPCRespItem := range RPCList {
		HTTPList[index] = &v1pb.GetListExResp_List{
			Id:             RPCRespItem.Id,
			Title:          RPCRespItem.Title,
			JumpPath:       RPCRespItem.JumpPath,
			DevicePlatform: RPCRespItem.DevicePlatform,
			DeviceBuild:    RPCRespItem.DeviceBuild,
			StartTime:      RPCRespItem.StartTime,
			EndTime:        RPCRespItem.EndTime,
			Status:         RPCRespItem.Status,
			DeviceLimit:    RPCRespItem.DeviceLimit,
			ImageUrl:       RPCRespItem.ImageUrl,
			JumpPathType:   RPCRespItem.JumpPathType,
			JumpTime:       RPCRespItem.JumpTime,
			Type:           RPCRespItem.Type,
		}
	}
	return
}
