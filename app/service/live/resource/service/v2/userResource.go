package v2

import (
	"context"
	v2pb "go-common/app/service/live/resource/api/grpc/v2"
	"go-common/app/service/live/resource/conf"
	"go-common/app/service/live/resource/dao"
	"go-common/app/service/live/resource/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// UserResourceService struct
type UserResourceService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewUserResourceService init
func NewUserResourceService(c *conf.Config) (s *UserResourceService) {
	s = &UserResourceService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

// Add implementation
// Add 添加资源接口
func (s *UserResourceService) Add(ctx context.Context, req *v2pb.AddReq) (resp *v2pb.AddResp, err error) {
	resp = &v2pb.AddResp{}

	// 查询新的资源，新的资源ID = 现有的最大ID + 1
	customID, err := s.dao.GetMaxCustomID(ctx, req.ResType)
	if err != nil {
		log.Error("查找最大的资源ID res_type : %d : %v", req.ResType, err)
		return
	}

	customID++

	res := &model.UserResource{
		ResType:  req.ResType,
		CustomID: customID,
		Title:    req.Title,
		URL:      req.Url,
		Weight:   req.Weight,
		Status:   0,
		Creator:  req.Creator,
	}

	// 插入新资源
	info, err := s.dao.AddUserResource(ctx, res)

	if err != nil {
		log.Error("resource.addSResource d.db.Exec err: %v", err)
	}

	resp.Id = info.ID
	resp.ResType = info.ResType
	resp.CustomId = info.CustomID
	resp.Title = info.Title
	resp.Url = info.URL
	resp.Weight = info.Weight
	resp.Creator = info.Creator
	resp.Status = info.Status
	resp.Ctime = info.Ctime.Format("2006-01-02 15:04:05")
	resp.Mtime = info.Mtime.Format("2006-01-02 15:04:05")

	return
}

// Edit 编辑现有资源
func (s *UserResourceService) Edit(ctx context.Context, req *v2pb.EditReq) (resp *v2pb.EditResp, err error) {
	resp = &v2pb.EditResp{}

	info, err := s.dao.GetUserResourceInfo(ctx, req.ResType, req.CustomId)
	if info.ID <= 0 {
		err = ecode.EditResErr
		return
	}

	update := make(map[string]interface{})

	// 名称
	if req.Title != "" {
		update["title"] = req.Title
	}

	// 链接
	if req.Url != "" {
		update["url"] = req.Url
	}

	// 权重
	if req.Weight > 0 {
		update["weight"] = req.Weight
	}

	effectRow, info, err := s.dao.EditUserResource(ctx, req.ResType, req.CustomId, update)
	if err != nil || effectRow <= 0 {
		err = ecode.EditResErr
		return
	}

	resp.Id = info.ID
	resp.ResType = info.ResType
	resp.CustomId = info.CustomID
	resp.Title = info.Title
	resp.Url = info.URL
	resp.Weight = info.Weight
	resp.Creator = info.Creator
	resp.Status = info.Status
	resp.Ctime = info.Ctime.Format("2006-01-02 15:04:05")
	resp.Mtime = info.Mtime.Format("2006-01-02 15:04:05")

	return
}

// Query 请求单个资源
func (s *UserResourceService) Query(ctx context.Context, req *v2pb.QueryReq) (resp *v2pb.QueryResp, err error) {
	resp = &v2pb.QueryResp{}

	info, err := s.dao.GetUserResourceInfo(ctx, req.ResType, req.CustomId)

	if err != nil {
		return
	}

	resp.Id = info.ID
	resp.ResType = info.ResType
	resp.CustomId = info.CustomID
	resp.Title = info.Title
	resp.Url = info.URL
	resp.Weight = info.Weight
	resp.Creator = info.Creator
	resp.Status = info.Status
	resp.Ctime = info.Ctime.Format("2006-01-02 15:04:05")
	resp.Mtime = info.Mtime.Format("2006-01-02 15:04:05")

	return
}

// List 获取资源列表
func (s *UserResourceService) List(ctx context.Context, req *v2pb.ListReq) (resp *v2pb.ListResp, err error) {
	var Page int32 = 1
	var pageSize int32 = 50

	if req.Page > 0 {
		Page = req.Page
	}

	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	resp = &v2pb.ListResp{}
	resp.CurrentPage = Page
	resp.TotalCount, _ = s.dao.GetMaxCustomID(ctx, req.ResType)

	list, err := s.dao.ListUserResourceInfo(ctx, req.ResType, Page, pageSize)
	if len(list) <= 0 {
		return
	}

	for _, v := range list {
		elem := &v2pb.ListResp_List{}
		elem.Id = v.ID
		elem.ResType = v.ResType
		elem.CustomId = v.CustomID
		elem.Title = v.Title
		elem.Url = v.URL
		elem.Weight = v.Weight
		elem.Creator = v.Creator
		elem.Status = v.Status
		elem.Ctime = v.Ctime.Format("2006-01-02 15:04:05")
		elem.Mtime = v.Mtime.Format("2006-01-02 15:04:05")

		resp.List = append(resp.List, elem)
	}

	return
}

// SetStatus implementation
// SetStatus 更改资源状态
func (s *UserResourceService) SetStatus(ctx context.Context, req *v2pb.SetStatusReq) (resp *v2pb.SetStatusResp, err error) {
	resp = &v2pb.SetStatusResp{}
	effectRow, err := s.dao.SetUserResourceStatus(ctx, req.ResType, req.CustomId, req.Status)
	if err != nil || effectRow == 0 {
		err = ecode.EditResErr
		return
	}

	return
}
