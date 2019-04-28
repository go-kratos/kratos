package ugc

import (
	"encoding/json"
	"time"

	model "go-common/app/job/main/tv/model/pgc"
	ugcmdl "go-common/app/job/main/tv/model/ugc"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func (s *Service) repCidproc() {
	defer s.waiter.Done()
	var toRepCids []int64
	for {
		cid, ok := <-s.repCidCh
		if !ok {
			log.Warn("[repCidproc] channel quit")
			return
		}
		toRepCids = append(toRepCids, cid)
		if len(toRepCids) < s.c.UgcSync.Batch.ReportCidPS { // not enough cid, stay waiting
			time.Sleep(5 * time.Second)
			continue
		}
		goCids := make([]int64, len(toRepCids))
		copy(goCids, toRepCids)
		toRepCids = []int64{}
		if err := s.reportCids(goCids); err != nil {
			log.Error("reportCids Cids %v, Err %v", goCids, err)
			continue
		}
	}
}

func (s *Service) audCidproc() {
	defer s.waiter.Done()
	var (
		toAudAids = make(map[int64]int)
		ps        = s.c.UgcSync.Batch.ReportCidPS
	)
	for {
		aids, ok := <-s.audAidCh
		if !ok {
			log.Warn("[audCidproc] channel quit")
			return
		}
		for _, aid := range aids {
			toAudAids[aid] = 1 // use map to remove duplicated aids
		}
		if len(toAudAids) < ps { // not enough cid, stay waiting
			time.Sleep(3 * time.Second)
			continue
		}
		distinctAIDs := pickKeys(toAudAids)
		toAudAids = make(map[int64]int)
		if err := s.wrapSyncLic(ctx, distinctAIDs); err != nil {
			log.Error("audCidproc Aids %v, Err %v", distinctAIDs, err)
			continue
		}
		log.Info("audCidproc Apply %d Aids: %v", len(aids), distinctAIDs)
	}
}

func (s *Service) reshelfArcproc() {
	defer s.waiter.Done()
	var (
		reshelfAids = make(map[int64]int)
		ps          = s.c.UgcSync.Batch.ReshelfPS
	)
	for {
		aid, ok := <-s.reshelfAidCh
		if !ok {
			log.Warn("[reshelfAid] channel quit")
			return
		}
		reshelfAids[aid] = 1       // use map to remove duplicated aids
		if len(reshelfAids) < ps { // not enough cid, stay waiting
			time.Sleep(3 * time.Second)
			continue
		}
		distinctAIDs := pickKeys(reshelfAids)
		reshelfAids = make(map[int64]int)
		if offAids, err := s.cmsDao.OffArcs(ctx, distinctAIDs); err != nil {
			log.Error("reshelfAid OffArcs Aids %v, Err %v", distinctAIDs, err)
			continue
		} else if len(offAids) == 0 {
			log.Warn("reshelfAid OffArcs Origin Aids %v, after filter it's empty", distinctAIDs)
			continue
		} else {
			if err = s.cmsDao.ReshelfArcs(ctx, offAids); err != nil {
				log.Error("reshelfAid OffAids %v, ReshelfArcs Err %v", offAids, err)
				continue
			}
			log.Info("reshelfAid Apply %d Aids: %v", len(offAids), offAids)
		}
	}
}

func (s *Service) reportCids(cids []int64) (err error) {
	var cidReq []*ugcmdl.CidReq
	for _, v := range cids {
		cidReq = append(cidReq, &ugcmdl.CidReq{CID: v})
	}
	for i := 0; i < _apiRetry; i++ {
		if err = s.dao.RepCidBatch(ctx, cidReq); err == nil {
			break
		}
	}
	if err != nil { // 3 times still error
		log.Error("ReportCid Cids %v Err %v", cids, err)
		return
	}
	err = s.dao.FinishReport(ctx, cids)
	log.Info("ReportCids %v, Len %d, Succ!", cids, len(cids))
	return
}

// consume Databus message; beause daily modification is not many, so use simple loop
func (s *Service) consumeVideo() {
	defer s.waiter.Done()
	for {
		msg, ok := <-s.ugcSub.Messages()
		if !ok {
			log.Info("databus: tv-job video consumer exit!")
			return
		}
		msg.Commit()
		time.Sleep(1 * time.Millisecond)
	Loop:
		for {
			select {
			case s.consumerLimit <- struct{}{}: // would block if already 2 goroutines:
				go s.UgcDbus(msg)
				break Loop
			default:
				log.Warn("consumeVideo thread Full!!!")
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// UgcDbus def.
func (s *Service) UgcDbus(msg *databus.Message) {
	m := &model.DatabusRes{}
	log.Info("[consumeVideo] New Message: %s", msg)
	if err := json.Unmarshal(msg.Value, m); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
		<-s.consumerLimit // clean the space for new consumer to begin
		return
	}
	if m.Action == "delete" {
		log.Info("[consumeVideo] Video Deletion, We ignore:<%v>,<%v>", m, msg.Value)
		<-s.consumerLimit // clean the space for new consumer to begin
		return
	}
	if m.Table == "ugc_video" {
		s.videoDatabus(msg.Value)
	} else if m.Table == "ugc_archive" {
		s.arcDatabus(msg.Value)
	} else {
		log.Error("[consumeVideo] Wrong Table Name: ", m.Table)
	}
	<-s.consumerLimit // clean the space for new consumer to begin
}
