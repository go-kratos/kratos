package service

import (
	"context"
	"fmt"
	"time"

	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const ROOM_LEN = 300

//实时数据处理逻辑,生成list
func (s *Service) minuteDataToCacheList() {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()
	ctx = GetTraceLogCtx(ctx, "minuteDataToCacheList")
	log.Infow(ctx, "log", "minuteDataToCacheList_start")
	//获取需要生成的数据的content列表
	contentMap, err := s.dao.GetContentMap(ctx)
	if err != nil {
		log.Errorw(ctx, "log", fmt.Sprintf("minuteDataToCacheList_GetContentMap_error:reply=%v;err=%v", contentMap, err))
		return
	}
	//获取全量开播房间
	allLiveingRoomIds, err := s.dao.GetAllLiveRoomIds(ctx)
	if allLiveingRoomIds == nil || err != nil {
		log.Errorw(ctx, "log", fmt.Sprintf("minuteDataToCacheList_allLiveingRoomIds_error:reply=%v;err=%v", allLiveingRoomIds, err))
		return
	}
	eg := errgroup.Group{}
	for content := range contentMap {
		log.Infow(ctx, fmt.Sprintf("minuteDataToCacheList_start:%s", content))
		eg.Go(func(contentParam string) func() error {
			slice := make([]int64, 0)
			for i := 0; i < len(allLiveingRoomIds); {
				end := ROOM_LEN + i
				if ROOM_LEN+i >= len(allLiveingRoomIds) {
					end = len(allLiveingRoomIds)
				}
				slice = allLiveingRoomIds[i:end]
				if len(slice) <= 0 {
					break
				} else {
					s.dao.CreateCacheList(ctx, slice, contentParam)

				}
				i = end
			}
			log.Infow(ctx, "log", fmt.Sprintf("minuteDataToCacheList_end_content=%s;err=%v", contentParam, err))
			return nil
		}(content))
	}
	eg.Wait()
	log.Infow(ctx, "log", "minuteDataToCacheList_end")
	return
}

func (s *Service) minuteDataToDB() {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	defer cancel()
	ctx = GetTraceLogCtx(ctx, "minuteDataToDB")
	log.Infow(ctx, "log", "minuteDataToDB_start")
	//获取需要生成的数据的content列表
	contentMap, err := s.dao.GetContentMap(ctx)
	if err != nil {
		log.Errorw(ctx, "log", fmt.Sprintf("data_allLiveingRoomIds_error:reply=%v;err=%v", contentMap, err))
		return
	}
	//获取全量开播房间
	allLiveingRoomIds, err := s.dao.GetAllLiveRoomIds(ctx)
	if allLiveingRoomIds == nil || err != nil {
		log.Errorw(ctx, "log", fmt.Sprintf("data_allLiveingRoomIds_error:reply=%v;err=%v", allLiveingRoomIds, err))
		return
	}
	eg := errgroup.Group{}
	for content := range contentMap {
		log.Info("minuteDataToCacheList_start:" + content)
		eg.Go(func(contentParam string) func() error {
			return func() (err error) {
				for _, roomId := range allLiveingRoomIds {
					s.dao.CreateDBData(ctx, []int64{int64(roomId)}, contentParam)
				}
				log.Infow(ctx, "log", fmt.Sprintf("minuteDataToCacheList_end_content=%s;err=%v", contentParam, err))
				return
			}
		}(content))
	}
	eg.Wait()
	log.Infow(ctx, "log", fmt.Sprintf("minuteDataToDB_end"))
	return
}
