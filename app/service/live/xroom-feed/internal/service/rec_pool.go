package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go-common/library/log"
	"go-common/library/xstr"

	pb "go-common/app/service/live/xroom-feed/api"
)

// GetList implementation
// 根据模块位置获取投放列表 position=>RoomItem
func (s *Service) GetList(ctx context.Context, req *pb.RecPoolReq) (resp *pb.RecPoolResp, err error) {
	resp = &pb.RecPoolResp{}
	//1.得到module positions_num
	module := req.ModuleType
	//6..
	positionNum := req.PositionNum
	//24..
	moduleExistRooms, err := xstr.SplitInts(req.ModuleExistRooms)
	moduleExistRoomsMap := make(map[int64]int64)
	if err != nil {
		log.Error("[recPool.GetList]moduleExistRooms SplitInts err:%+v", err)
	} else if len(moduleExistRooms) > 0 {
		for index, iRoomID := range moduleExistRooms {
			moduleExistRoomsMap[iRoomID] = int64(index) + 1
		}
	}
	otherExistRooms, err := xstr.SplitInts(req.OtherExistRooms)
	otherExistRoomsMap := make(map[int64]int64)
	if err != nil {
		log.Error("[recPool.GetList]otherExistRooms SplitInts err:%+v", err)
	} else if len(moduleExistRooms) > 0 {
		for index, iRoomID := range otherExistRooms {
			otherExistRoomsMap[iRoomID] = int64(index) + 1
		}
	}

	tidyRoomIDs := make([]int64, 0)
	tidyRoomIDsMap := make(map[int64]int64)
	tidyPoolIDsMap := make(map[int64]int64)
	existPosRoomIdMap := make(map[int64]bool)
	existCount := int64(0)
	var pos int64
	for pos = 1; pos <= positionNum; pos++ {
		modulePosition := fmt.Sprintf("%d_%d", module, pos)
		awardPoolIdsMap := make(map[int64]float64)
		//2.module_position map[string][] from 内存
		recPoolConfs := s.GetPoolConfFromMem(modulePosition)
		if len(recPoolConfs) == 0 {
			continue
		}
		//3.抽奖候选集（数据非空校验、true_percent在routine）
		for _, recPoolConf := range recPoolConfs {
			if recPoolConf == nil {
				log.Warn("[recPool.GetList]recPoolConf empty, pos:%+v", modulePosition)
				continue
			}
			if recPoolConf.TruePercent <= 0 {
				log.Warn("[recPool.GetList]recPoolConf TruePercent abnormal, recPoolConf:%+v", recPoolConf)
				continue
			}
			awardPoolIdsMap[recPoolConf.ID] = recPoolConf.TruePercent
		}
		if len(awardPoolIdsMap) <= 0 {
			continue
		}
		//4.开始抽奖
		log.Info("[recPool.GetList]awardPoolIdsMap result, module:%d; res:%+v; from:%+v", module, awardPoolIdsMap, req.From)
		resPoolId := getAwardSource(awardPoolIdsMap)
		if resPoolId <= 0 {
			continue
		}
		tidyPoolIDsMap[pos] = resPoolId

		//5.去重
		resPoolDataMap := s.dao.GetRecPoolByID(ctx, []int64{resPoolId})
		resPoolData, ok := resPoolDataMap[resPoolId]
		if !ok {
			continue
		}
		tidyPoolRoomIDs := make([]int64, 0)
		for _, iRoomID := range resPoolData {

			if ePos, ok := moduleExistRoomsMap[iRoomID]; ok && (ePos+existCount < pos) {
				//原位置已存在&&位置比投放的位置好
				continue
			}
			if _, ok := existPosRoomIdMap[iRoomID]; ok {
				//同模块前面的位置已经吐过了
				continue
			}
			if _, ok := otherExistRoomsMap[iRoomID]; ok {
				//其它模块已经吐过了
				continue
			}
			tidyPoolRoomIDs = append(tidyPoolRoomIDs, iRoomID)
		}
		//6.随机取一个
		tidyRoomIDsLen := len(tidyPoolRoomIDs)
		if tidyRoomIDsLen <= 0 {
			continue
		}
		rand.Seed(time.Now().UnixNano())
		tidyRoomID := tidyPoolRoomIDs[rand.Intn(tidyRoomIDsLen)]
		tidyRoomIDs = append(tidyRoomIDs, tidyRoomID)
		tidyRoomIDsMap[pos] = tidyRoomID
		existPosRoomIdMap[tidyRoomID] = true
		existCount++
	}
	//7.批量获取rec信息并返回
	if len(tidyRoomIDs) <= 0 {
		return
	}
	log.Info("[recPool.GetList]tidyRoomIDsMap result, module:%d; res:%+v; from:%+v", module, tidyRoomIDsMap, req.From)
	resp.List = make(map[int64]*pb.RoomItem)
	tidyRoomInfoMap := s.dao.GetRecInfoByRoomid(ctx, tidyRoomIDs)

	for pos, iRoomID := range tidyRoomIDsMap {
		rInfo, ok := tidyRoomInfoMap[iRoomID]
		if !ok {
			log.Warn("[recPool.GetList]recInfo empty:%+v from:%+v", tidyRoomInfoMap, req.From)
			continue
		}
		recType := int64(10000)
		if poolConfId, ok := tidyPoolIDsMap[pos]; ok {
			recType += poolConfId
		}
		resp.List[pos] = &pb.RoomItem{
			RoomId:          iRoomID,
			Uid:             rInfo.Uid,
			Title:           rInfo.Title,
			PopularityCount: rInfo.PopularityCount,
			Keyframe:        rInfo.KeyFrame,
			Cover:           rInfo.Cover,
			AreaId:          rInfo.AreaID,
			AreaName:        rInfo.AreaName,
			ParentAreaId:    rInfo.ParentAreaID,
			ParentAreaName:  rInfo.ParentAreaName,
			RecType:         recType,
		}
	}

	return
}

// awardMap id=>百分比/权重
func getAwardSource(awardMap map[int64]float64) (tidyID int64) {
	type awarSource struct {
		index  int64
		offset float64
		weight float64
	}

	awardSli := make([]*awarSource, 0, len(awardMap))
	var sumCount float64
	for n, c := range awardMap {
		//精度 支持小数点后两位
		c *= 100
		a := awarSource{
			index:  n,
			offset: sumCount,
			weight: c,
		}
		awardSli = append(awardSli, &a)
		sumCount += c
	}
	awardIndex := rand.Int63n(int64(sumCount))
	for _, u := range awardSli {
		if u.offset+u.weight > float64(awardIndex) {
			tidyID = u.index
			return
		}
	}
	return
}
