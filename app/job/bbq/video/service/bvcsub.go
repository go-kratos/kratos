package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/job/bbq/video/model"
	videov1 "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/log"
	"strings"
)

//BvcTransSub ...
func (s *Service) BvcTransSub() {
	msgs := s.bvcSub.Messages()
	for {
		var (
			err error
			vr  *model.VideoRepRaw
		)
		c := context.Background()
		msg, ok := <-msgs

		//release subscription
		if s.c.SubBvcControl.Control == 2 {
			msg.Commit()
			continue
		}
		if !ok {
			log.Info("BvcTransSub databus Consumer exit")
			return
		}
		res := &model.DatabusBVCTransSub{}
		log.Infov(context.Background(), log.KV("log", fmt.Sprintf("databus message %s", string(msg.Value))))
		if err = json.Unmarshal(msg.Value, &res); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			msg.Commit()
			continue
		}
		if vr, err = s.dao.RawVideo(c, res.SVID); err != nil {
			msg.Commit()
			continue
		}
		//resource check
		if vr.SyncStatus&model.SourceXcodeCover > 0 {
			if err = s.importVideo(c, vr); err != nil {
				log.Errorw(c, "errmsg", "importVideo err", "req", vr, "err", err)
				msg.Commit()
				continue
			}
			s.syncTag(c, vr.Tag)

			s.dao.UpdateSyncStatus(c, vr.SVID, model.SourceOnshelf)
		}

		msg.Commit()
	}
}

//importVideo put video on shelf
func (s *Service) importVideo(c context.Context, vr *model.VideoRepRaw) (err error) {
	// var (
	// 	st int64
	// )
	// if vr.From == model.VideoFromBILI {
	// 	st = model.VideoStPassReview
	// } else {
	// 	st = model.VideoStPendingPassReview
	// }
	req := &videov1.ImportVideoInfo{
		AVID:        vr.AVID,
		Svid:        vr.SVID,
		MID:         vr.MID,
		CID:         vr.CID,
		SubTID:      vr.SubTID,
		TID:         vr.TID,
		Title:       vr.Title,
		Pubtime:     vr.Pubtime,
		From:        int64(vr.From),
		CoverUrl:    vr.CoverURL,
		CoverHeight: vr.CoverHeight,
		CoverWidth:  vr.CoverWidth,
		//State:         st,
		HomeImgHeight: vr.HomeImgHeight,
		HomeImgUrl:    vr.HomeImgURL,
		HomeImgWidth:  vr.HomeImgWidth,
	}
	for i := 0; i < _retryTimes; i++ {
		if _, err = s.dao.VideoClient.ImportVideo(c, req); err == nil {
			break
		}
	}
	return
}

//syncTag sync video from bilibili common tag
func (s *Service) syncTag(c context.Context, t string) (err error) {
	if t == "" {
		return
	}
	var (
		arrTag []string
		tag    []*videov1.TagInfo
	)
	arrTag = strings.Split(t, ",")
	for _, v := range arrTag {
		tmp := &videov1.TagInfo{
			TagName: v,
			TagType: 3,
		}
		tag = append(tag, tmp)
	}
	reqTag := &videov1.SyncVideoTagRequest{
		TagInfos: tag,
	}
	if _, err = s.dao.VideoClient.SyncTag(c, reqTag); err != nil {
		log.Error("sync tag err :%v,tag:%v", err, tag)
	}
	return
}
