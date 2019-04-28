package service

import (
	"context"
	"time"

	"github.com/json-iterator/go"

	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/conf"
)

// AppSetting .
func (s *Service) AppSetting(c context.Context, arg *v1.AppSettingRequest) (resp *v1.AppSettingResponse, err error) {
	plat := 1
	if arg.Base.Client == "ios" {
		plat = 2
	}
	appVerison, err := s.dao.FetchNewAppVersion(c, plat, arg.VersionCode)
	newVersion := uint8(0)
	if err == nil && appVerison.ID > 0 {
		newVersion = uint8(1)
	}

	// TODO：暂时先这样，快速上线，以后改(也许。。)
	ver := 1
	if arg.Base.Client == "ios" && arg.VersionCode > 100000 {
		ver = 2
	} else if arg.Base.Client == "android" && arg.VersionCode > 101000 {
		ver = 2
	}

	appResource, err := s.dao.FetchAppResource(c, plat, ver)
	current := time.Now().Unix()

	pubSetting := make(map[string]interface{})
	b, _ := jsoniter.Marshal(conf.App)
	jsoniter.Unmarshal(b, &pubSetting)
	for _, v := range appResource {
		if v.StartTime.Time().Unix() < current && current < v.EndTime.Time().Unix() {
			pubSetting["dynamic_effect"] = v.ID
		}
	}
	resp = &v1.AppSettingResponse{
		Public: pubSetting,
		Update: &v1.AppUpdate{
			NewVersion: newVersion,
			Info:       appVerison,
		},
		Resources: appResource,
	}
	return
}

// AppPackage .
func (s *Service) AppPackage(c context.Context, lastest int) (resp []*v1.AppPackage, err error) {
	resp = make([]*v1.AppPackage, 0)

	if lastest > 0 {
		result, e := s.dao.FetchNewAppVersion(c, 1, 0)
		resp = append(resp, &v1.AppPackage{
			ID:          int64(result.ID),
			Platform:    uint8(result.Platform),
			VersionName: result.Name,
			VersionCode: uint32(result.Code),
			Title:       result.Title,
			Content:     result.Content,
			Download:    result.Download,
			MD5:         result.MD5,
			Size:        int32(result.Size),
			Force:       uint8(result.Force),
			Status:      uint8(result.Status),
		})
		err = e
	} else {
		resp, err = s.dao.FetchAppPackage(c)
	}

	return
}
