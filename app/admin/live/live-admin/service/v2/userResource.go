package v2

import (
	"context"
	v2pb "go-common/app/admin/live/live-admin/api/http/v2"
	"go-common/app/admin/live/live-admin/conf"
	v2rspb "go-common/app/service/live/resource/api/grpc/v2"
	"go-common/library/log"
)

// UserResourceService struct
type UserResourceService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	v2rsCli *v2rspb.Client
}

//NewUserResourceService init
func NewUserResourceService(c *conf.Config) (s *UserResourceService) {
	s = &UserResourceService{
		conf: c,
	}

	var svc *v2rspb.Client
	var err error

	log.Info("ResourceServiceV2 Init: %+v", s.conf.ResourceClientV2)
	if svc, err = v2rspb.NewClient(s.conf.ResourceClientV2); err != nil {
		panic(err)
	}
	s.v2rsCli = svc
	return s
}

// Add implementation
// Add 添加资源接口
// `method:"POST" internal:"true" `
func (s *UserResourceService) Add(ctx context.Context, req *v2pb.UserResourceAddReq) (resp *v2pb.UserResourceAddResp, err error) {
	respRPC, err := s.v2rsCli.Add(ctx, &v2rspb.AddReq{
		ResType: req.ResType,
		Title:   req.Title,
		Url:     req.Url,
		Weight:  req.Weight,
		Creator: req.Creator,
	})
	if err == nil {
		resp = &v2pb.UserResourceAddResp{
			Id:       respRPC.Id,
			CustomId: respRPC.CustomId,
		}
	}
	return
}

// Edit implementation
// Edit 编辑现有资源
// `method:"POST" internal:"true" `
func (s *UserResourceService) Edit(ctx context.Context, req *v2pb.UserResourceEditReq) (resp *v2pb.UserResourceEditResp, err error) {
	resp = &v2pb.UserResourceEditResp{}
	_, err = s.v2rsCli.Edit(ctx, &v2rspb.EditReq{
		ResType:  req.ResType,
		Title:    req.Title,
		Url:      req.Url,
		Weight:   req.Weight,
		CustomId: req.CustomId,
	})
	return
}

// Get implementation
// Get 获取资源列表
// `method:"GET" internal:"true" `
func (s *UserResourceService) Get(ctx context.Context, req *v2pb.UserResourceListReq) (resp *v2pb.UserResourceListResp, err error) {
	respRPC, err := s.v2rsCli.List(ctx, &v2rspb.ListReq{
		ResType:  req.ResType,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err == nil {
		resp = &v2pb.UserResourceListResp{
			CurrentPage: respRPC.CurrentPage,
			TotalCount:  respRPC.TotalCount,
			List:        convertRPCListRes(respRPC.List),
		}
	}
	return
}

// SetStatus implementation
// SetStatus 更改资源状态
// `method:"POST" internal:"true" `
func (s *UserResourceService) SetStatus(ctx context.Context, req *v2pb.UserResourceSetStatusReq) (resp *v2pb.UserResourceSetStatusResp, err error) {
	resp = &v2pb.UserResourceSetStatusResp{}
	_, err = s.v2rsCli.SetStatus(ctx, &v2rspb.SetStatusReq{
		ResType:  req.ResType,
		CustomId: req.CustomId,
		Status:   req.Status,
	})
	return
}

// GetSingle implementation
// Query 请求单个资源
func (s *UserResourceService) GetSingle(ctx context.Context, req *v2pb.UserResourceGetSingleReq) (resp *v2pb.UserResourceGetSingleResp, err error) {
	respRPC, err := s.v2rsCli.Query(ctx, &v2rspb.QueryReq{
		CustomId: req.CustomId,
		ResType:  req.ResType,
	})
	if err == nil {
		resp = &v2pb.UserResourceGetSingleResp{
			Id:       respRPC.Id,
			ResType:  respRPC.ResType,
			CustomId: respRPC.CustomId,
			Title:    respRPC.Title,
			Url:      respRPC.Url,
			Weight:   respRPC.Weight,
			Creator:  respRPC.Creator,
			Status:   respRPC.Status,
			Ctime:    respRPC.Ctime,
			Mtime:    respRPC.Mtime,
		}
	}
	return
}

func convertRPCListRes(RPCList []*v2rspb.ListResp_List) (HTTPList []*v2pb.UserResourceListResp_List) {
	HTTPList = make([]*v2pb.UserResourceListResp_List, len(RPCList))
	for index, RPCListItem := range RPCList {
		HTTPList[index] = &v2pb.UserResourceListResp_List{
			Id:       RPCListItem.Id,
			ResType:  RPCListItem.ResType,
			CustomId: RPCListItem.CustomId,
			Title:    RPCListItem.Title,
			Url:      RPCListItem.Url,
			Weight:   RPCListItem.Weight,
			Creator:  RPCListItem.Creator,
			Status:   RPCListItem.Status,
			Ctime:    RPCListItem.Ctime,
			Mtime:    RPCListItem.Mtime,
		}
	}
	return
}
