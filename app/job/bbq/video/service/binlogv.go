package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/job/bbq/video/model"
	"go-common/app/service/bbq/common"
	topic "go-common/app/service/bbq/topic/api"
	"go-common/library/log"
	"strconv"
)

// videoConsumeproc 视频表消费
func (s *Service) videoBinlogSub() {
	var msgs = s.videoSub.Messages()
	for {
		var err error
		msg, ok := <-msgs
		if !ok {
			log.Info("userCanal databus Consumer exit")
			return
		}
		res := &model.DatabusRes{}
		log.Infov(context.Background(), log.KV("log", fmt.Sprintf("canal message %s", string(msg.Value))))
		if err = json.Unmarshal(msg.Value, &res); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			msg.Commit()
			continue
		}

		if res.Table != "video" || (res.Action != "update" && res.Action != "insert") {
			msg.Commit()
			continue
		}

		var vNew, vOld *model.VideoRaw
		if res.Action == "insert" || res.Action == "update" {
			if err = json.Unmarshal(res.New, &vNew); err != nil {
				log.Error("video unmarshal err(%v) data[%s]", err, string(res.New))
				continue
			}
		}
		if res.Action == "update" {
			if err = json.Unmarshal(res.Old, &vOld); err != nil {
				log.Error("video unmarshal err(%v) data[%s]", err, string(res.Old))
				continue
			}
		}

		//idempotent consume
		for i := 0; i < _retryTimes; i++ {
			//fetured video state subscription
			if err = s.VideoStateSub(vNew, vOld); err == nil {
				break
			}
		}
		//s.UpdateCms(context.Background(), vNew)
		//register comment
		for i := 0; i < _retryTimes; i++ {
			//merge related information subscription
			if err = s.CommentReg(context.Background(), vNew.SVID, model.StateActive); err == nil {
				break
			}
		}

		if res.Action == "insert" {
			for i := 0; i < _retryTimes; i++ {
				log.V(1).Infow(context.Background(), "log", "merge up info", "retry_time", i, "mid", vNew.MID, "svid", vNew.SVID)
				//merge related information subscription
				if err = s.MergeUpInfoSub(vNew); err == nil {
					break
				}
			}
		}

		//unidempotent consume
		if res.Action == "update" {
			s.UpdateStaInfoSub(vNew, vOld)
		} else if res.Action == "insert" {
			s.AddSVTotal(vNew)
		}
		msg.Commit()
	}
}

//UpdateCms ..
func (s *Service) UpdateCms(c context.Context, vNew *model.VideoRaw) (err error) {
	if err = s.dao.UpdateCms(c, vNew); err != nil {
		log.Warnw(c, "event", fmt.Sprintf("updateCms err:%v,param:%v", err, vNew))
	}
	return
}

// VideoStateSub 视频状态变更消费
func (s *Service) VideoStateSub(vNew *model.VideoRaw, vOld *model.VideoRaw) (err error) {
	log.Infow(context.Background(), "log", "one video state sub", "svid", vNew.SVID)

	s.SaveVideo2ES(strconv.Itoa(int(vNew.SVID)))

	if vOld == nil || vNew.State != vOld.State {
		var ids []int64
		ids, err = s.dao.GetRecallOpVideo(context.Background())
		if err != nil {
			log.Warnw(context.Background(), "log", "get recall op video fail")
			return
		}
		needSetRecallOpVideo := true
		if vNew.State == _selection {
			for _, id := range ids {
				if id == vNew.SVID {
					needSetRecallOpVideo = false
					break
				}
			}
			ids = append(ids, vNew.SVID)
		} else if vNew.State != _selection {
			index := -1
			for i, id := range ids {
				if id == vNew.SVID {
					index = i
					break
				}
			}
			if index != -1 {
				ids = append(ids[:index], ids[index+1:]...)
			} else {
				needSetRecallOpVideo = false
			}
		}
		if needSetRecallOpVideo {
			if err = s.dao.SetRecallOpVideo(context.Background(), ids); err != nil {
				log.Warnw(context.Background(), "log", "get recall op video fail")
				return
			}
		}
	}

	// 话题状态变更
	needUpdateTopicVideoState := false
	topicState := topic.TopicVideoStateUnAvailable
	if vOld == nil {
		needUpdateTopicVideoState = true
		// update topic video
		if common.IsTopicSvStateAvailable(int64(vNew.State)) {
			topicState = topic.TopicVideoStateAvailable
		}
	} else {
		if common.IsTopicSvStateAvailable(int64(vNew.State)) != common.IsTopicSvStateAvailable(int64(vOld.State)) {
			needUpdateTopicVideoState = true
			if common.IsTopicSvStateAvailable(int64(vNew.State)) {
				topicState = topic.TopicVideoStateAvailable
			}
		}
	}
	if needUpdateTopicVideoState {
		_, err = s.topicClient.UpdateVideoState(context.Background(), &topic.UpdateVideoStateReq{Svid: vNew.SVID, State: int32(topicState)})
		log.Infow(context.Background(), "log", "update topic video state", "svid", vNew.SVID, "new_state", vNew.State)
		if err != nil {
			log.Warnw(context.Background(), "log", "update topic video state", "svid", vNew.SVID, "new_state", vNew.State)
			return
		}
	}
	return
}

//MergeUpInfoSub ..
func (s *Service) MergeUpInfoSub(vNew *model.VideoRaw) (err error) {
	mid := vNew.MID
	if err = s.dao.MergeUpInfo(mid); err != nil {
		log.Error("MergeUpInfo failed,err:%v,mid:%d", err, mid)
	}
	return
}

//UpdateStaInfoSub ...
func (s *Service) UpdateStaInfoSub(vNew *model.VideoRaw, vOld *model.VideoRaw) {
	if vOld == nil || vNew.State == vOld.State {
		return
	}
	if vNew.State == model.VideoStInactive {
		s.dao.UpdateUVSt(vNew.MID, "unshelf_av_total")
	} else if vNew.State == model.VideoStDeleted {
		s.dao.UpdateUVStDel(vNew.MID, "av_total")
	}

	if vOld.State == model.VideoStInactive {
		s.dao.UpdateUVStDel(vNew.MID, "unshelf_av_total")
	} else if vNew.State == model.VideoStDeleted {
		s.dao.UpdateUVSt(vOld.MID, "av_total")
	}
}

//AddSVTotal ...
func (s *Service) AddSVTotal(vNew *model.VideoRaw) {
	s.dao.AddSVTotal(vNew.MID)
}
