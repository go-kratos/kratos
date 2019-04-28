package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/job/bbq/video/model"
	videov1 "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/log"
)

//videoRepositorySub video_repository subscription .
func (s *Service) videoRepositoryBinlogSub() {
	msgs := s.videoRep.Messages()
	for {
		var err error
		msg, ok := <-msgs
		if !ok {
			log.Info("video_repository databus Consumer exit")
			return
		}
		res := &model.DatabusRes{}
		log.Infov(context.Background(), log.KV("log", fmt.Sprintf("canal message %s", string(msg.Value))))
		if err = json.Unmarshal(msg.Value, &res); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			msg.Commit()
			continue
		}

		if res.Table != "video_repository" || (res.Action != "update" && res.Action != "insert") {
			msg.Commit()
			continue
		}

		//unserialize databus struct
		var vNew, vOld *model.VideoRepRaw
		if res.Action == "insert" || res.Action == "update" {
			if err = json.Unmarshal(res.New, &vNew); err != nil {
				log.Error("video unmarshal err(%v) data[%s]", err, string(res.New))
				msg.Commit()
				continue
			}
		}
		if res.Action == "update" {
			if err = json.Unmarshal(res.Old, &vOld); err != nil {
				log.Error("video unmarshal err(%v) data[%s]", err, string(res.Old))
				msg.Commit()
				continue
			}
		}

		if res.Action == "insert" {
			for i := 0; i < _retryTimes; i++ {
				if err = s.PepareResource(vNew); err == nil {
					break
				}
			}
		}
		msg.Commit()
	}
}

//PepareResource ...
func (s *Service) PepareResource(vNew *model.VideoRepRaw) (err error) {
	var (
		ctx  = context.Background()
		SVID int64
		row  *model.VideoRepRaw
	)
	//bbq/cms video not trans to bvc
	if vNew.From == model.VideoFromBBQ || vNew.From == model.VideoFromCMS {
		return
	}
	if row, err = s.dao.RawVideoByID(ctx, vNew.ID); err != nil {
		return
	}
	if row.SVID > 0 {
		SVID = row.SVID
	} else {
		req := &videov1.CreateIDRequest{
			Mid: vNew.MID,
		}
		var rep *videov1.CreateIDResponse
		if rep, err = s.dao.VideoClient.CreateID(ctx, req); err != nil {
			log.Error("Numbering device return err:%v", err)
			return
		}
		if err = s.dao.UpdateSvid(context.Background(), vNew.ID, rep.NewId); err != nil {
			return
		}
		SVID = rep.NewId
	}
	reqBvc := &videov1.BVideoTransRequset{
		SVID: SVID,
		CID:  vNew.CID,
	}
	log.Info("bvc trans commit req:%v", reqBvc)
	if _, err = s.dao.VideoClient.BVCTransCommit(ctx, reqBvc); err != nil {
		log.Error("BVCTransCommit err :%v,req:%v", err, reqBvc)
	}
	s.dao.UpdateSyncStatus(ctx, SVID, model.SourceRequest)
	return
}
