package v1

import (
	"context"
	"encoding/json"
	"regexp"
	"time"

	v1pb "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/app/service/live/resource/conf"
	"go-common/app/service/live/resource/dao"
	"go-common/app/service/live/resource/model"
	"go-common/library/ecode"
)

// ResourceService struct
type ResourceService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewResourceService init
func NewResourceService(c *conf.Config) (s *ResourceService) {
	s = &ResourceService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

func checkURL(s string) (r bool) {
	reg, _ := regexp.Compile("^(https|http)://")
	r = reg.MatchString(s)
	return
}

// Add implementation
// Add 添加资源接口
func (s *ResourceService) Add(ctx context.Context, req *v1pb.AddReq) (resp *v1pb.AddResp, err error) {
	type device struct {
		Platform string `json:"platform"`
		Build    int64  `json:"build"`
		Limit    int64  `json:"limit"`
	}
	deviceArrs := make([]device, 0)
	e := json.Unmarshal([]byte(req.Device), &deviceArrs)
	if e != nil {
		err = ecode.DeviceError
		return
	}

	imageArr := make(map[string]interface{})
	if checkURL(req.ImageUrl) == false {
		err = ecode.CheckURLErr
		return
	}
	imageArr["imageUrl"] = req.ImageUrl
	if req.JumpPath != "" {
		if checkURL(req.JumpPath) == false {
			err = ecode.CheckURLErr
			return
		}
		imageArr["jumpPath"] = req.JumpPath
	}
	if req.JumpPathType != 0 {
		imageArr["jumpPathType"] = req.JumpPathType
	}
	if req.JumpTime != 0 {
		imageArr["jumpTime"] = req.JumpTime
	}

	b, err := json.Marshal(imageArr)
	if err != nil {
		err = ecode.ResourceParamErr
		return
	}

	resp = &v1pb.AddResp{}
	loc, _ := time.LoadLocation("Local")
	startTime, stErr := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, loc)
	if stErr != nil {
		err = ecode.TimeForErr
		return
	}
	endTime, etErr := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, loc)
	if etErr != nil {
		err = ecode.TimeForErr
		return
	}

	for _, da := range deviceArrs {
		existInfo, errSe := s.dao.SelectByTypeAndPlatform(ctx, req.Type, da.Platform)
		if errSe != nil {
			err = ecode.AddResourceErr
			return
		}
		if existInfo != nil {
			err = ecode.RepdAddErr
			return
		}
	}

	for _, deviceArr := range deviceArrs {
		insert := &model.Resource{
			Platform:  deviceArr.Platform,
			Build:     deviceArr.Build,
			LimitType: deviceArr.Limit,
			StartTime: startTime,
			EndTime:   endTime,
			Type:      req.Type,
			Title:     req.Title,
			ImageInfo: string(b),
		}
		reply, _ := s.dao.AddResource(ctx, insert)
		if reply <= 0 {
			err = ecode.AddResourceErr
			return
		}
		resp.Id = append(resp.Id, reply)
	}
	return
}

// AddEx implementation
// AddEx 添加资源接口
func (s *ResourceService) AddEx(ctx context.Context, req *v1pb.AddReq) (resp *v1pb.AddResp, err error) {
	type device struct {
		Platform string `json:"platform"`
		Build    int64  `json:"build"`
		Limit    int64  `json:"limit"`
	}
	deviceArrs := make([]device, 0)
	e := json.Unmarshal([]byte(req.Device), &deviceArrs)
	if e != nil {
		err = ecode.DeviceError
		return
	}

	imageArr := make(map[string]interface{})
	if checkURL(req.ImageUrl) == false {
		err = ecode.CheckURLErr
		return
	}
	imageArr["imageUrl"] = req.ImageUrl
	if req.JumpPath != "" {
		if checkURL(req.JumpPath) == false {
			err = ecode.CheckURLErr
			return
		}
		imageArr["jumpPath"] = req.JumpPath
	}
	if req.JumpPathType != 0 {
		imageArr["jumpPathType"] = req.JumpPathType
	}
	if req.JumpTime != 0 {
		imageArr["jumpTime"] = req.JumpTime
	}

	b, err := json.Marshal(imageArr)
	if err != nil {
		err = ecode.ResourceParamErr
		return
	}

	resp = &v1pb.AddResp{}
	loc, _ := time.LoadLocation("Local")
	startTime, stErr := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, loc)
	if stErr != nil {
		err = ecode.TimeForErr
		return
	}
	endTime, etErr := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, loc)
	if etErr != nil {
		err = ecode.TimeForErr
		return
	}

	for _, deviceArr := range deviceArrs {
		insert := &model.Resource{
			Platform:  deviceArr.Platform,
			Build:     deviceArr.Build,
			LimitType: deviceArr.Limit,
			StartTime: startTime,
			EndTime:   endTime,
			Type:      req.Type,
			Title:     req.Title,
			ImageInfo: string(b),
		}
		reply, _ := s.dao.AddResource(ctx, insert)
		if reply <= 0 {
			err = ecode.AddResourceErr
			return
		}
		resp.Id = append(resp.Id, reply)
	}
	return
}

//Edit implementation
//Edit 编辑资源接口
func (s *ResourceService) Edit(ctx context.Context, req *v1pb.EditReq) (resp *v1pb.EditResp, err error) {
	resp = &v1pb.EditResp{}
	update := make(map[string]interface{})
	resourceInfo, err := s.dao.SelectById(ctx, req.Id)
	if err != nil || resourceInfo == nil {
		err = ecode.SeltResErr
		return
	}
	if resourceInfo.ID == 0 {
		err = ecode.SeltResErr
		return
	}

	imageInfo := resourceInfo.ImageInfo
	imageInfoArr := make(map[string]interface{})
	e := json.Unmarshal([]byte(imageInfo), &imageInfoArr)
	if e != nil {
		err = ecode.ResourceParamErr
		return
	}
	if req.Title != "" {
		update["title"] = req.Title
	}

	loc, _ := time.LoadLocation("Local")
	if req.StartTime != "" {
		update["start_time"], _ = time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, loc)
	}
	if req.EndTime != "" {
		update["end_time"], _ = time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, loc)
	}

	if req.JumpPath != "" {
		if checkURL(req.JumpPath) == false {
			err = ecode.CheckURLErr
			return
		}
		imageInfoArr["jumpPath"] = req.JumpPath
	}
	if req.JumpPathType != 0 {
		imageInfoArr["jumpPathType"] = req.JumpPathType
	}
	if req.JumpTime != 0 {
		imageInfoArr["jumpTime"] = req.JumpTime
	}
	if req.ImageUrl != "" {
		if checkURL(req.ImageUrl) == false {
			err = ecode.CheckURLErr
			return
		}
		imageInfoArr["imageUrl"] = req.ImageUrl
	}
	if imageInfoArr != nil {
		b, e := json.Marshal(imageInfoArr)
		if e != nil || b == nil {
			err = ecode.ResourceParamErr
			return
		}
		update["image_info"] = string(b)
	}
	reply, err := s.dao.EditResource(ctx, req.Id, update)
	if err != nil || reply <= 0 {
		err = ecode.EditResErr
		return
	}
	return
}

// Offline implementation
// Offline 下线资源接口
func (s *ResourceService) Offline(ctx context.Context, req *v1pb.OfflineReq) (resp *v1pb.OfflineResp, err error) {
	resp = &v1pb.OfflineResp{}
	_, err = s.dao.OfflineResource(ctx, req.Id)
	if err != nil {
		err = ecode.OfflineResErr
		return
	}
	return
}

// GetList implementation
// GetList 获取资源列表
func (s *ResourceService) GetList(ctx context.Context, req *v1pb.GetListReq) (resp *v1pb.GetListResp, err error) {
	resp = &v1pb.GetListResp{}
	reply, err := s.dao.GetResourceList(ctx, req.Type, req.Page, req.PageSize)
	count, err := s.dao.GetDBCount(ctx, req.Type)
	if err != nil {
		err = ecode.GetListResErr
		return
	}
	for _, v := range reply {
		type updateImage struct {
			JumpPath     string `json:"JumpPath"`
			JumpPathType int64  `json:"JumpPathType"`
			JumpTime     int64  `json:"JumpTime"`
			ImageUrl     string `json:"ImageUrl"`
		}
		ImageArr := &updateImage{}
		e := json.Unmarshal([]byte(v.ImageInfo), ImageArr)
		if e != nil {
			continue
		}
		var status int64
		if v.EndTime.Unix() <= time.Now().Unix() {
			status = -1
		}
		if (v.StartTime.Unix() <= time.Now().Unix()) && v.EndTime.Unix() > time.Now().Unix() {
			status = 1
		}
		list := &v1pb.GetListResp_List{}
		list.Id = v.ID
		list.ImageUrl = ImageArr.ImageUrl
		list.JumpPath = ImageArr.JumpPath
		list.StartTime = v.StartTime.Format("2006-01-02 15:04:05")
		list.EndTime = v.EndTime.Format("2006-01-02 15:04:05")
		list.DevicePlatform = v.Platform
		list.DeviceBuild = v.Build
		list.DeviceLimit = v.LimitType
		list.Title = v.Title
		list.Status = status
		list.JumpPathType = ImageArr.JumpPathType
		list.JumpTime = ImageArr.JumpTime
		resp.List = append(resp.List, list)
	}
	resp.CurrentPage = req.Page
	resp.TotalCount = count
	return
}

// GetPlatformList implementation
// 获取平台列表
func (s *ResourceService) GetPlatformList(ctx context.Context, req *v1pb.GetPlatformListReq) (resp *v1pb.GetPlatformListResp, err error) {
	resp = &v1pb.GetPlatformListResp{}
	pList := []string{"ios", "ios_link", "android", "android_link", "ipad", "pc_link"}
	resp = &v1pb.GetPlatformListResp{
		Platform: pList,
	}
	return
}

// GetListEx implementation
// GetListEx 获取资源列表
func (s *ResourceService) GetListEx(ctx context.Context, req *v1pb.GetListExReq) (resp *v1pb.GetListExResp, err error) {
	resp = &v1pb.GetListExResp{}
	reply, count, err := s.dao.GetResourceListEx(ctx, req.Type, req.Page, req.PageSize, req.DevicePlatform, req.Status, req.StartTime, req.EndTime)
	if err != nil {
		err = ecode.GetListResErr
		return
	}
	for _, v := range reply {
		type updateImage struct {
			JumpPath     string `json:"JumpPath"`
			JumpPathType int64  `json:"JumpPathType"`
			JumpTime     int64  `json:"JumpTime"`
			ImageUrl     string `json:"ImageUrl"`
		}
		ImageArr := &updateImage{}
		e := json.Unmarshal([]byte(v.ImageInfo), ImageArr)
		if e != nil {
			continue
		}
		var status int64
		if v.EndTime.Unix() <= time.Now().Unix() {
			status = -1
		}
		if (v.StartTime.Unix() <= time.Now().Unix()) && v.EndTime.Unix() > time.Now().Unix() {
			status = 1
		}
		list := &v1pb.GetListExResp_List{}
		list.Id = v.ID
		list.ImageUrl = ImageArr.ImageUrl
		list.JumpPath = ImageArr.JumpPath
		list.StartTime = v.StartTime.Format("2006-01-02 15:04:05")
		list.EndTime = v.EndTime.Format("2006-01-02 15:04:05")
		list.DevicePlatform = v.Platform
		list.DeviceBuild = v.Build
		list.DeviceLimit = v.LimitType
		list.Title = v.Title
		list.Status = status
		list.JumpPathType = ImageArr.JumpPathType
		list.JumpTime = ImageArr.JumpTime
		list.Type = v.Type
		resp.List = append(resp.List, list)
	}
	resp.CurrentPage = req.Page
	resp.TotalCount = count
	return
}
