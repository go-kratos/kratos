package pgc

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	appDao "go-common/app/job/main/tv/dao/app"
	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/ecode"
	"go-common/library/log"
	timex "go-common/library/time"
)

type cntFunc func(ctx context.Context) (count int, err error)
type refreshFunc func(ctx context.Context, LastID int, nbData int) (myLast int, err error)
type reqCachePro struct {
	cnt     cntFunc
	proName string
	refresh refreshFunc
	ps      int
}

func (s *Service) cacheProducer(ctx context.Context, req *reqCachePro) (err error) {
	var (
		count    int
		pagesize = req.ps
		maxID    = 0 // the max ID of the latest piece
		begin    = time.Now()
	)
	if count, err = req.cnt(ctx); err != nil {
		log.Error("[%s] CountEP error [%v]", req.proName, err)
		return
	}
	nbPiece := appDao.NumPce(count, pagesize)
	log.Info("[%s] NumPiece %d, Pagesize %d", req.proName, nbPiece, pagesize)
	for i := 0; i < nbPiece; i++ {
		newMaxID, errR := req.refresh(ctx, maxID, pagesize)
		if errR != nil {
			log.Error("[%s] Pick Piece %d Error, Ignore it", req.proName, i)
			continue
		}
		if newMaxID > maxID {
			maxID = newMaxID
		} else { // fatal error
			log.Error("[%s] MaxID is not increasing! [%d,%d]", req.proName, newMaxID, maxID)
			return
		}
		time.Sleep(time.Duration(s.c.UgcSync.Frequency.ProducerFre)) // pause after each piece produced
		log.Info("[%s] Pagesize %d, Num of piece %d, Time Already %v", req.proName, pagesize, i, time.Since(begin))
	}
	log.Info("[%s] Finish! Pagesize %d, Num of piece %d, Time %v", req.proName, pagesize, nbPiece, time.Since(begin))
	return
}

// refreshCache refreshes the cache of ugc and pgc
func (s *Service) refreshCache() {
	var (
		ctx   = context.Background()
		begin = time.Now()
		pgcPS = s.c.PlayControl.PieceSize
		reqEp = &reqCachePro{
			cnt:     s.dao.CountEP,
			proName: "epProducer",
			refresh: s.dao.RefreshEPMC,
			ps:      pgcPS,
		}
		reqSn = &reqCachePro{
			cnt:     s.dao.CountSeason,
			proName: "snProducer",
			refresh: s.dao.RefreshSnMC,
			ps:      pgcPS,
		}
	)
	if err := s.cacheProducer(ctx, reqEp); err != nil {
		log.Error("reqEp Err %v", err)
		return
	}
	if err := s.cacheProducer(ctx, reqSn); err != nil {
		log.Error("reqSn Err %v", err)
	}
	log.Info("refreshCache Finish, Time %v", time.Since(begin))
}

// stock EP&Season auth info and intervention info in MC
func (s *Service) stockContent(jsonstr json.RawMessage, tableName string) (err error) {
	// season stock in MC
	if tableName == "tv_ep_season" {
		sn := &model.DatabusSeason{}
		if err = json.Unmarshal(jsonstr, sn); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", jsonstr, err)
			return
		}
		if reflect.DeepEqual(sn.Old, sn.New) { // if media fields not modified, no need to update
			log.Info("SeasonID %d No need to update", sn.New.ID)
			return
		}
		return s.stockSeason(sn)
		// ep stock in MC
	} else if tableName == "tv_content" {
		ep := &model.DatabusEP{}
		if err = json.Unmarshal(jsonstr, ep); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", jsonstr, err)
			return
		}
		if reflect.DeepEqual(ep.Old, ep.New) { // if media fields not modified, no need to update
			log.Info("Epid %d No need to update", ep.New.EPID)
			return
		}
		return s.stockEP(ep)
	} else {
		return fmt.Errorf("Databus Msg (%s) - Incorrect Table (%s) ", jsonstr, tableName)
	}
}

func (s *Service) composeSnCMS(sn *model.MediaSn) *model.SeasonCMS {
	var (
		epid, order int
		err         error
		playtime    int64
	)
	if epid, order, err = s.dao.NewestOrder(ctx, sn.ID); err != nil {
		log.Warn("stockSeason NewestOrder Sid: %d, Err %v", sn.ID, err)
	}
	if playtime, err = appDao.TimeTrans(sn.Playtime); err != nil {
		log.Warn("stockSeason Playtime Sid: %d, Err %v", sn.ID, err)
	}
	return &model.SeasonCMS{
		SeasonID:    int(sn.ID),
		Cover:       sn.Cover,
		Desc:        sn.Desc,
		Title:       sn.Title,
		UpInfo:      sn.UpInfo,
		Category:    sn.Category,
		Area:        sn.Area,
		Playtime:    timex.Time(playtime),
		Role:        sn.Role,
		Staff:       sn.Staff,
		TotalNum:    sn.TotalNum,
		Style:       sn.Style,
		NewestOrder: order,
		NewestEPID:  epid,
		PayStatus:   sn.Status, // databus sn logic
	}
}

// treat the databus season msg, stock the auth & media info in MC
func (s *Service) stockSeason(sn *model.DatabusSeason) (err error) {
	var (
		snSub   *model.TVEpSeason
		snAuth  = sn.New.ToSimple()      // auth info in MC
		snMedia = s.composeSnCMS(sn.New) // media info in MC
	)
	s.batchFilter(ctx, []*model.SeasonCMS{snMedia})                     // treat the newest NB logic
	if sn.New.Check == _seasonPassed && sn.Old.Check == _seasonPassed { // keep already passed logic
		if snSub, err = s.dao.Season(ctx, int(sn.New.ID)); err != nil {
			return
		}
		s.addRetrySn(snSub)
	}
	if err = s.dao.SetSeason(ctx, snAuth); err != nil { // auth
		log.Error("SetSeason error(%v)", snAuth, err)
		return
	}
	if err = s.dao.SetSnCMSCache(ctx, snMedia); err != nil { // media
		log.Error("SetSnCMSCache error(%v)", snMedia, err)
		return
	}
	if err = s.listMtn(sn.Old, sn.New); err != nil { // maintenance of the zone list in Redis
		log.Error("stockContent listMtn error(%v)", sn.New, err)
	}
	return
}

// treat the databus ep msg, stock the auth & media info in MC
func (s *Service) stockEP(ep *model.DatabusEP) (err error) {
	var (
		epAuth  = ep.New.ToSimple()
		epMedia = ep.New.ToCMS()
		epSub   *model.Content
	)
	if ep.New.State == _epPassed && ep.Old.State == _epPassed { // keep already passed logic
		if epSub, err = s.dao.Cont(ctx, ep.New.EPID); err != nil {
			return
		}
		s.addRetryEp(epSub)
	}
	if err = s.dao.SetEP(ctx, epAuth); err != nil { // set ep auth MC
		return
	}
	if err = s.dao.SetEpCMSCache(ctx, epMedia); err != nil { // set ep media MC
		return
	}
	err = s.updateSnCMS(epAuth.SeasonID)
	return
}

// updateSnCMS picks the season info from DB and update the CMS cache
func (s *Service) updateSnCMS(sid int) (err error) {
	var snMedia *model.SeasonCMS
	if snMedia, err = s.dao.PickSeason(ctx, sid); err != nil { // pick season cms info
		log.Error("stockEP PickSeason Sid: %d, Err: %v", sid, err)
		return
	}
	if snMedia == nil { // season info not found
		err = ecode.NothingFound
		log.Error("stockEP PickSeason Sid: %d, Err: %v", sid, err)
		return
	}
	s.batchFilter(ctx, []*model.SeasonCMS{snMedia})
	if err = s.dao.SetSnCMSCache(ctx, snMedia); err != nil { // ep update, we also consider to update its season info for the "latest" info
		log.Error("SetSnCMSCache error(%v)", snMedia, err)
	}
	return
}

// consume Databus message; because daily modification is not many, so use simple loop
func (s *Service) consumeContent() {
	defer s.waiterConsumer.Done()
	for {
		msg, ok := <-s.contentSub.Messages()
		if !ok {
			log.Info("databus: tv-job ep/season consumer exit!")
			return
		}
		msg.Commit()
		s.treatMsg(msg.Value)
		time.Sleep(1 * time.Millisecond)
	}
}

func (s *Service) treatMsg(msg json.RawMessage) {
	m := &model.DatabusRes{}
	log.Info("[ConsumeContent] New Message: %s", msg)
	if err := json.Unmarshal(msg, m); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", msg, err)
		return
	}
	if m.Action == "delete" {
		log.Info("[ConsumeContent] Content Deletion, We ignore:<%v>,<%v>", m, msg)
		return
	}
	if err := s.stockContent(msg, m.Table); err != nil {
		log.Error("stockContent.(%s,%s), error(%v)", msg, m.Table, err)
		return
	}
}
