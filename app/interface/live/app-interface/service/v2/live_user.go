package v2

import (
	"context"
	"fmt"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
	liveUserDao "go-common/app/interface/live/app-interface/dao/live_user"
	liveUserV1 "go-common/app/service/live/live_user/api/liverpc/v1"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_mmtagType = 12
	_mseaType  = 14
)

// GetIndexV2TagList ...
// 获取APP首页 - 我的个人标签列表
func (s *IndexService) GetIndexV2TagList(ctx context.Context, req *liveUserV1.UserSettingGetTagReq) (ret []*v2pb.MMyTag, err error) {
	ret = []*v2pb.MMyTag{}
	ExtraInfo := &v2pb.MyTagExtra{
		Offline: []*v2pb.OfflineTag{},
	}
	List := make([]*v2pb.MyTagItem, 0)
	res := &v2pb.MMyTag{}
	d := &liveUserDao.Dao{}
	moduleList := s.GetAllModuleInfoMapFromCache(ctx)
	module, ok := moduleList[_mmtagType]
	if false == ok || 0 == len(module) {
		log.Error("[GetUserTagList]AllMinfoMap is empty error:%+v, %+v", err, moduleList)
		return
	}
	for _, v := range module {
		if v.Type == _mmtagType {
			var data *liveUserV1.UserSettingGetTagResp_Data
			data, err = d.GetUserTagList(ctx, req)
			if err != nil {
				log.Error("[GetUserTagList]live_user.v1.getTag rpc error:%+v, %+v", err, data)
				return
			}
			ExtraInfo.IsGray = data.IsGray
			for _, offlineInfo := range data.Offline {
				ExtraInfo.Offline = append(ExtraInfo.Offline, &v2pb.OfflineTag{Id: int64(offlineInfo.Id), AreaV2Name: offlineInfo.Name})
			}

			for _, tagInfo := range data.Tags {
				link := fmt.Sprintf("http://live.bilibili.com/app/area?parent_area_id=%d&parent_area_name=%s&area_id=%d&area_name=%s", tagInfo.ParentId, tagInfo.ParentName, tagInfo.Id, tagInfo.Name)
				List = append(List, &v2pb.MyTagItem{AreaV2Id: int64(tagInfo.Id), AreaV2Name: tagInfo.Name, AreaV2ParentId: int64(tagInfo.ParentId), AreaV2ParentName: tagInfo.ParentName, Link: link, Pic: tagInfo.Pic, IsAdvice: int64(tagInfo.IsAdvice)})
				if len(List) >= 4 {
					break
				}
			}
			List = append(List, &v2pb.MyTagItem{AreaV2Id: 0, AreaV2Name: "全部标签", AreaV2ParentId: 0, AreaV2ParentName: "", Pic: "http://i0.hdslb.com/bfs/vc/ff03528785fc8c91491d79e440398484811d6d87.png", Link: "http://live.bilibili.com/app/mytag/", IsAdvice: 1})
			res.ExtraInfo = ExtraInfo
			res.List = List
			res.ModuleInfo = v
		}
		break
	}

	ret = append(ret, res)
	return
}

// GetIndexV2SeaPatrol ...
// 获取APP首页 - 我的大航海提示信息
func (s *IndexService) GetIndexV2SeaPatrol(ctx context.Context, req *liveUserV1.NoteGetReq) (ret []*v2pb.MSeaPatrol, err error) {
	ret = []*v2pb.MSeaPatrol{}
	_, isUIDSet := metadata.Value(ctx, metadata.Mid).(int64)
	if !isUIDSet {
		return
	}
	ExtraInfo := &v2pb.PicItem{}
	res := &v2pb.MSeaPatrol{}
	d := &liveUserDao.Dao{}
	moduleList := s.GetAllModuleInfoMapFromCache(ctx)
	module, ok := moduleList[_mseaType]
	if false == ok || 0 == len(module) {
		log.Error("[GetIndexV2SeaPatrol]AllMinfoMap is empty error:%+v, %+v", err, moduleList)
		return
	}
	for _, v := range module {
		if v.Type == _mseaType {
			var data *liveUserV1.NoteGetResp_Data
			data, err = d.GetDaHangHai(ctx, req)
			if err != nil {
				log.Error("[GetIndexV2SeaPatrol]live_user.v1.getNoteSea rpc error:%+v, %+v", err, data)
				return
			}
			if data.Title == "" || data.Link == "" { // 返回信息为空
				return
			}
			ExtraInfo.Title = data.Title
			ExtraInfo.Pic = data.Logo
			ExtraInfo.Link = data.Link
			ExtraInfo.Content = data.Content
			ExtraInfo.Id = 0

			res.ExtraInfo = ExtraInfo
			res.ModuleInfo = v
		}
		break
	}

	ret = append(ret, res)

	return
}
