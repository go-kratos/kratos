package service

import (
	"context"
	"math"

	daoAnchorV1 "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_video = 33
	_music = 34
)

var _fields = []string{
	"room_id",
	"uid",
	"title",
	"popularity_count",
	"keyframe",
	"cover",
	"parent_area_id",
	"parent_area_name",
	"area_id",
	"area_name",
}

func (s *Service) setRecInfoCache(ctx context.Context, roomIds []int64) (filterIds []int64) {
	filterIds = make([]int64, 0)
	chunkSize := 100
	// 批次
	chunkNum := int(math.Ceil(float64(len(roomIds)) / float64(chunkSize)))
	chunkIds := make([][]int64, chunkNum)
	wg := errgroup.Group{}
	for i := 1; i <= chunkNum; i++ {
		x := i
		wg.Go(func() error {
			chunkRoomIds := make([]int64, 10)
			if x == chunkNum {
				chunkRoomIds = roomIds[(x-1)*chunkSize:]
			} else {
				chunkRoomIds = roomIds[(x-1)*chunkSize : x*chunkSize]
			}
			resp, err := s.daoAnchor.FetchRoomByIDs(ctx, &daoAnchorV1.RoomByIDsReq{
				RoomIds: chunkRoomIds,
				Fields:  _fields,
			})
			if err != nil {
				log.Error("[setRecInfoCache]FetchRoomByIDs_error:%+v", err)
				return nil
			}

			if resp == nil || len(resp.RoomDataSet) == 0 {
				log.Info("[setRecInfoCache]FetchRoomByIDs_empty")
				return nil
			}
			filterRoomData := make(map[int64]*daoAnchorV1.RoomData)
			for roomId, data := range resp.RoomDataSet {
				if data == nil {
					continue
				}
				if data.Cover == "" && data.Keyframe == "" {
					log.Info("[setRecInfoCache]emptyCoverOrKeyFrame, roomId:%d ,cover:%s, keyframe:%s", data.RoomId, data.Cover, data.Keyframe)
					continue
				}
				if data.AreaName == "" || data.ParentAreaName == "" {
					log.Info("[setRecInfoCache]emptyAreaName, roomId:%d ,areaName:%s, parentAreaName:%s", data.RoomId, data.AreaName, data.ParentAreaName)
					continue
				}
				if data.AreaId == _video || data.AreaId == _music {
					log.Info("[setRecInfoCache]musicOrVideoArea, roomId:%d ,area:%d", data.RoomId, data.AreaId)
					continue
				}
				if s.isBlackRoomID(data.RoomId) {
					log.Info("[setRecInfoCache]is IndexBlackRoomID, roomId:%d ", data.RoomId)
					continue
				}
				chunkIds[x-1] = append(chunkIds[x-1], roomId)
				filterRoomData[roomId] = data
			}
			s.dao.SetRecInfoCache(ctx, filterRoomData)
			return nil
		})
	}

	err := wg.Wait()
	if err != nil {
		log.Error("[setRecInfoCache]waitError:%+v", err)
	}
	for _, ids := range chunkIds {
		for _, id := range ids {
			filterIds = append(filterIds, id)
		}
	}
	return
}
