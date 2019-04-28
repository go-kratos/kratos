package v2

import (
	"context"
	"strconv"

	"go-common/app/interface/live/app-interface/dao"
	"go-common/app/interface/live/app-interface/model"
	"go-common/app/service/main/account/api"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
	"go-common/app/interface/live/app-interface/conf"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// 统一cover获取方式
func (s *IndexService) getCover(userCover string, systemCover string) (cover string) {
	if userCover != "" {
		cover = userCover
	} else {
		cover = systemCover
	}
	return
}

func (s *IndexService) ifHitSkyHorse(mid int64, device string) (hit bool) {
	if mid == 0 {
		return false
	}
	if device == "pad" {
		return false
	}
	lastMid := strconv.Itoa(int(mid % 100))
	if len(lastMid) < 2 {
		lastMid = "0" + lastMid
	}
	_, isSkyHorseGray := s.conf.SkyHorseGray[lastMid]

	return isSkyHorseGray && conf.Conf.SkyHorseStatus
}

func (s *IndexService) ifHitLiveRec(mid int64, device string) (hit bool) {
	if mid == 0 {
		return false
	}
	if device == "pad" {
		return false
	}
	lastMid := strconv.Itoa(int(mid % 100))
	if len(lastMid) < 2 {
		lastMid = "0" + lastMid
	}
	_, isLiveRec := s.conf.LiveGray[lastMid]
	return isLiveRec
}

func (s *IndexService) getRecRoomList(ctx context.Context, roomIds []int64, recPoolRoomListResp map[int64]*v2pb.CommonRoomItem, build int64, platform string, idolDuplicateMap map[int64]bool, recType int64, quality int64) (respRecRoomList []*v2pb.CommonRoomItem, err error) {
	wg, wgCtx := errgroup.WithContext(ctx)

	multiRoomListResp := make(map[int64]*roomV2.RoomGetByIdsResp_RoomInfo)
	wg.Go(func() error {
		// 天马房间基础信息,有错误cancel其他没必要执行
		fields := []string{
			"roomid",
			"title",
			"uname",
			"online",
			"cover",
			"user_cover",
			"link",
			"face",
			"area_v2_parent_id",
			"area_v2_parent_name",
			"area_v2_id",
			"area_v2_name",
			"broadcast_type",
			"uid",
		}
		multiRoomList, err := s.roomDao.GetRoomInfoByIds(wgCtx, roomIds, fields, "app-interface-skyHorseRec")
		if err != nil {
			log.Error("[getRecRoomList]getByIds error:%+v", err)
		}
		multiRoomListResp = multiRoomList
		return err
	})

	pendantRoomListResp := make(map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result)
	wg.Go(func() error {
		pendantRoomList, err := s.roomDao.GetRoomPendant(wgCtx, roomIds, "mobile_index_badge", 2)
		if err != nil {
			log.Error("[getRecRoomList]getPendantByIds error:%+v", err)
		}
		pendantRoomListResp = pendantRoomList
		return nil
	})

	err = wg.Wait()
	if err != nil {
		log.Error("[getRecRoomList]wait error(%+v)", err)
		return
	}

	respSlice := make([]*roomV2.RoomGetByIdsResp_RoomInfo, 0)
	for _, roomBaseInfo := range multiRoomListResp {
		respSlice = append(respSlice, roomBaseInfo)
	}

	respRecRoomList = make([]*v2pb.CommonRoomItem, 0)
	for i := 0; i < 6; i++ {
		if recInfo, ok := recPoolRoomListResp[int64(i+1)]; ok {
			if _, ok := idolDuplicateMap[recInfo.Roomid]; !ok {
				respRecRoomList = append(respRecRoomList, recInfo)
				continue
			}
		}

		if len(respSlice) <= 0 {
			continue
		}
		tmpItem := respSlice[0:1][0]
		respSlice = respSlice[1:]
		pendantValue, pendantBgPic, pendantBgColor := s.getPendant(tmpItem.Roomid, pendantRoomListResp)
		// 统一cover产品逻辑
		cover := s.getCover(tmpItem.UserCover, tmpItem.Cover)

		respRecRoomList = append(respRecRoomList, &v2pb.CommonRoomItem{
			Roomid:           tmpItem.Roomid,
			Title:            tmpItem.Title,
			Uname:            tmpItem.Uname,
			Online:           tmpItem.Online,
			Cover:            cover,
			Link:             "/" + strconv.Itoa(int(tmpItem.Roomid)),
			Face:             tmpItem.Face,
			AreaV2ParentId:   tmpItem.AreaV2ParentId,
			AreaV2ParentName: tmpItem.AreaV2ParentName,
			AreaV2Id:         tmpItem.AreaV2Id,
			AreaV2Name:       tmpItem.AreaV2Name,
			BroadcastType:    tmpItem.BroadcastType,
			PendentRu:        pendantValue,
			PendentRuPic:     pendantBgPic,
			PendentRuColor:   pendantBgColor,
			RecType:          recType,
		})
	}

	s.getPlayUrl(ctx, respRecRoomList, quality, build, platform)

	return
}

func (s *IndexService) getExtraDataForRoom(ctx context.Context, roomIds []int64, uids []int64, roomIdToUid map[int64]int64) (extraInfo map[int64]*model.ExtraRecInfo) {
	wg, wgCtx := errgroup.WithContext(ctx)
	extraInfo = make(map[int64]*model.ExtraRecInfo)

	userInfos := make(map[int64]*api.Info)
	wg.Go(func() error {
		userResult, err := s.accountDao.GetUserInfos(wgCtx, uids)
		if err != nil {
			log.Error("[getExtraDataForRoom]getByIds error:%+v", err)
		}
		userInfos = userResult
		return err
	})

	pendantRoomListResp := make(map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result)
	wg.Go(func() error {
		pendantRoomList, err := s.roomDao.GetRoomPendant(wgCtx, roomIds, "mobile_index_badge", 2)
		if err != nil {
			log.Error("[getExtraDataForRoom]getPendantByIds error:%+v", err)
		}
		pendantRoomListResp = pendantRoomList
		return nil
	})

	err := wg.Wait()
	if err != nil {
		log.Error("[getExtraDataForRoom]getExtraDataForRoom_waitError:%+v", err)
		return
	}

	for _, roomId := range roomIds {
		pendantValue, pendantBgPic, pendantBgColor := s.getPendant(roomId, pendantRoomListResp)
		if _, ok := extraInfo[roomId]; !ok {
			extraInfo[roomId] = &model.ExtraRecInfo{}
		}
		extraInfo[roomId].PendentRu = pendantValue
		extraInfo[roomId].PendentRuPic = pendantBgPic
		extraInfo[roomId].PendentRuColor = pendantBgColor
		if uid, ok := roomIdToUid[roomId]; ok {
			if _, ok := userInfos[uid]; ok {
				extraInfo[roomId].UName = userInfos[uid].Name
				extraInfo[roomId].Face = userInfos[uid].Face
			}
		}
	}

	return
}

func (s *IndexService) getPendant(roomId int64, pendantRoomListResp map[int64]*roomV1.RoomPendantGetPendantByIdsResp_Result) (pendantValue, pendantBgPic, pendantBgColor string) {
	if pendantRoomListResp != nil {
		if _, ok := pendantRoomListResp[roomId]; ok {
			// 移动端取value, web取name
			pendantValue = pendantRoomListResp[roomId].Value
			pendantBgPic = pendantRoomListResp[roomId].BgPic
			if pendantRoomListResp[roomId].BgColor != "" {
				pendantBgColor = pendantRoomListResp[roomId].BgColor
			} else {
				pendantBgColor = _mobileIndexBadgeColorDefault
			}
		}
	}

	return
}

func (s *IndexService) getPlayUrl(ctx context.Context, roomList []*v2pb.CommonRoomItem, quality, build int64, platform string) {
	roomIdsForPlayUrl := make([]int64, 0)
	for _, commRoomBlock := range roomList {
		roomIdsForPlayUrl = append(roomIdsForPlayUrl, commRoomBlock.Roomid)
	}
	changeRoomListPlayURLMap := dao.BvcApi.GetPlayUrlMulti(ctx, roomIdsForPlayUrl, 0, quality, build, platform)

	for _, vv := range roomList {
		if changeRoomListPlayURLMap[vv.Roomid] != nil {
			vv.AcceptQuality = changeRoomListPlayURLMap[vv.Roomid].AcceptQuality
			vv.CurrentQuality = changeRoomListPlayURLMap[vv.Roomid].CurrentQuality
			vv.PlayUrl = changeRoomListPlayURLMap[vv.Roomid].Url["h264"]
			vv.PlayUrlH265 = changeRoomListPlayURLMap[vv.Roomid].Url["h265"]
		}
	}
}
