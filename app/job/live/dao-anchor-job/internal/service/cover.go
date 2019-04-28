package service

import (
	"context"
	"strings"
	"time"

	"go-common/app/job/live/dao-anchor-job/internal/dao"

	"go-common/library/sync/errgroup"

	"go-common/library/log"
)

//封面图/关键帧相关脚本

const ROOM_LEN_KEY_FRAME = 500

//updateKeyFrame  更新关键帧
func (s *Service) updateKeyFrame() {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()
	log.Info("updateKeyFrame_start")
	//获取全量开播房间
	allLiveingRoom, err := s.dao.GetAllLiveRoomIds(ctx)
	if allLiveingRoom == nil || err != nil {
		log.Error("updateKeyFrame_allLiveingRoom_error:reply=%v;err=%v", allLiveingRoom, err)
		return
	}
	slice := make([]int64, 0)
	eg := errgroup.Group{}
	for i := 0; i < len(allLiveingRoom); {
		end := ROOM_LEN_KEY_FRAME + i
		if (ROOM_LEN_KEY_FRAME + i) >= len(allLiveingRoom) {
			end = len(allLiveingRoom)
		}
		slice = allLiveingRoom[i:end]
		if len(slice) <= 0 {
			break
		} else {
			eg.Go(func(sliceParam []int64) func() error {
				return func() (err error) {
					for _, roomId := range sliceParam {
						coverUrl, err := s.dealKeyFrame(ctx, roomId)
						if err != nil {
							time.Sleep(time.Second)
							log.Error("updateKeyFrame_deal_error:roomId=%d;ketFrame=%s", roomId, coverUrl)
							continue
						}
						if coverUrl == "" {
							continue
						}
						//更新关键帧
						coverUrlArr := strings.Split(coverUrl, "?")
						coverUrl = coverUrlArr[0] + "?" + time.Now().Format("01021504")
						s.dao.UpdateRoomEx(ctx, roomId, []string{"keyframe"}, coverUrl)
						time.Sleep(time.Millisecond * 10)
					}
					return
				}
			}(slice))
		}
		i = end
	}
	eg.Wait()
	log.Info("updateKeyFrame_end")
	return
}

func (s *Service) dealKeyFrame(ctx context.Context, roomId int64) (coverUrl string, err error) {
	//二次确认是否关播，关播不再做
	roomInfos, err := s.dao.GetInfosByRoomIds(ctx, []int64{roomId}, []string{"live_status"})
	if err != nil {
		log.Error("updateKeyFrame_GetInfosByRoomIds_error:room_id=%d;err=%v", roomId, err)
		return
	}
	roomInfo := roomInfos[roomId]
	//未开播，不更新关键帧
	if roomInfo == nil || roomInfo.LiveStatus != dao.LIVE_OPEN {
		return
	}
	//判断是否为pk房间,pk房间不更新关键帧
	pkReply, err := s.dao.GetPkStatus(ctx, roomId)
	if err != nil {
		log.Error("updateKeyFrame_GetPkStatus_error:room_id=%d", roomId)
		return
	}
	if pkReply.PkStatus > 0 {
		return
	}
	//获取关键帧
	startTime := time.Now().Add(-time.Minute)
	endTime := time.Now()
	pics, err := s.dao.GetPicsByRoomId(ctx, roomId, startTime, endTime)
	if err != nil || pics == nil || pics[0] == "" {
		log.Warn("updateKeyFrame_GetPicsByRoomId_error:room_id=%d;pics=%v;err=%v", roomId, pics, err)
		return
	}
	//上传至bfs
	reply, err := s.dao.ImgDownload(ctx, pics[0])
	if err != nil || reply == nil {
		log.Warn("updateKeyFrame_ImgDownload_error:room_id=%d;pic=%s;err=%v;reply=%v", roomId, pics[0], err, reply)
		return
	}
	coverUrl, err = s.dao.ImgUpload(ctx, roomId, pics[0], reply)
	if err != nil || coverUrl == "" {
		log.Error("updateKeyFrame_ImgUploadBfs_error:room_id=%d;pic=%s", roomId, pics[0])
		return
	}
	return
}
