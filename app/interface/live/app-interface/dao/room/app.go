package room

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"go-common/app/interface/live/app-interface/conf"
	cDao "go-common/app/interface/live/app-interface/dao"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
)

// 根据moduleId查common房间列表
func (d *Dao) GetListByIds(ctx context.Context, moduleIds []int64) (respModuleList map[int64]*roomV2.AppIndexGetRoomListByIdsResp_RoomList, err error) {
	getModuleListByIdsTimeOut := time.Duration(conf.GetTimeout("getModuleListByIds", 100)) * time.Millisecond
	moduleList, err := cDao.RoomApi.V2AppIndex.GetRoomListByIds(rpcCtx.WithTimeout(ctx, getModuleListByIdsTimeOut), &roomV2.AppIndexGetRoomListByIdsReq{
		Ids: moduleIds,
	})
	if err != nil {
		log.Error("[GetListByIds]get all module ids rpc error, room.v1.AppIndex.GetModuleIds, error:%+v", err)
		err = errors.New("GET Module List Rpc error")
		return
	}

	if moduleList.Code != 0 || moduleList.Data == nil {
		log.Error("[GetListByIds]get all module ids return data error, code, %d, msg: %s", moduleList.Code, moduleList.Msg)
		err = errors.New("GET Module List return error")
		return
	}

	if len(moduleList.Data) <= 0 {
		log.Error("[GetListByIds]get all module ids empty error")
		err = errors.New("GET Module List empty error")
		return
	}

	respModuleList = moduleList.Data
	return
}

// 获取所有模块基础信息
func (d *Dao) GetAllModuleInfo(ctx context.Context, moduleId int64) (moduleInfoList []*roomV2.AppIndexGetBaseMInfoListResp_ModuleInfo, err error) {
	moduleInfoList = make([]*roomV2.AppIndexGetBaseMInfoListResp_ModuleInfo, 0)
	param := &roomV2.AppIndexGetBaseMInfoListReq{}
	if moduleId != 0 {
		param.ModuleId = moduleId
	}
	moduleInfoListOut, err := cDao.RoomApi.V2AppIndex.GetBaseMInfoList(ctx, param)

	if err != nil {
		log.Error("[GetAllModuleInfo]RoomApi.V2AppIndex.GetBaseMInfoList rpc error, RoomApi.V2AppIndex.GetBaseMInfoList, error:%+v", err)
		err = errors.New("RoomApi.V2AppIndex.GetBaseMInfoList rpc error")
		return
	}

	if moduleInfoListOut.Code != 0 || moduleInfoListOut.Data == nil || len(moduleInfoListOut.Data) <= 0 {
		log.Error("[GetAllModuleInfo]RoomApi.V2AppIndex.GetBaseMInfoList return data error, code, %d, msg: %s, error:%+v", moduleInfoListOut.Code, moduleInfoListOut.Msg, err)
		err = errors.New("RoomApi.V2AppIndex.GetBaseMInfoList return data error")
		return
	}

	moduleInfoList = moduleInfoListOut.Data

	return
}

// 根据模块id获取分区入口信息
func (d *Dao) GetAreaEntrance(ctx context.Context, ids []int64) (result map[int64]*roomV2.AppIndexGetPicListByIdsResp_ItemList, err error) {
	result = make(map[int64]*roomV2.AppIndexGetPicListByIdsResp_ItemList, 0)
	areaEntranceOut, err := cDao.RoomApi.V2AppIndex.GetPicListByIds(ctx, &roomV2.AppIndexGetPicListByIdsReq{Ids: ids})

	if err != nil {
		log.Error("[GetAreaEntrance]RoomApi.V2AppIndex.GetPicListByIds rpc error:%+v", err)
		err = errors.New("RoomApi.V2AppIndex.GetPicListByIds rpc error")
		return
	}

	if areaEntranceOut.Code != 0 || areaEntranceOut.Data == nil {
		log.Error("[GetAreaEntrance]RoomApi.V2AppIndex.GetPicListByIds return data error, code, %d, msg: %s, error:%+v", areaEntranceOut.Code, areaEntranceOut.Msg, err)
		err = errors.New("RoomApi.V2AppIndex.GetPicListByIds return data error")
		return
	}

	result = areaEntranceOut.Data

	return
}

// 根据分区ids获取房间列表
func (d *Dao) GetMultiRoomList(ctx context.Context, areaIds string, platform string) (result map[int64][]*roomV2.AppIndexGetMultiRoomListResp_RoomList, err error) {
	result = make(map[int64][]*roomV2.AppIndexGetMultiRoomListResp_RoomList)
	multiRoomListOut, err := cDao.RoomApi.V2AppIndex.GetMultiRoomList(ctx, &roomV2.AppIndexGetMultiRoomListReq{
		AreaIds:  areaIds,
		Platform: platform,
	})
	if err != nil {
		log.Error("[GetMultiRoomList]RoomApi.V2AppIndex.GetMultiRoomList rpc error:%+v", err)
		err = errors.New("RoomApi.V2AppIndex.GetMultiRoomList rpc error")
		return
	}

	if multiRoomListOut.Code != 0 || multiRoomListOut.Data == nil {
		log.Error("[GetMultiRoomList]RoomApi.V2AppIndex.GetMultiRoomList return data error, code, %d, msg: %s", multiRoomListOut.Code, multiRoomListOut.Msg)
		err = errors.New("RoomApi.V2AppIndex.GetMultiRoomList return data error")
		return
	}

	if len(multiRoomListOut.Data) <= 0 {
		log.Error("[GetMultiRoomList]RoomApi.V2AppIndex.GetMultiRoomList empty error")
		err = errors.New("RoomApi.V2AppIndex.GetMultiRoomList empty error")
		return
	}

	for _, item := range multiRoomListOut.Data {
		if item != nil {
			result[item.Id] = item.List
		}
	}
	return
}
