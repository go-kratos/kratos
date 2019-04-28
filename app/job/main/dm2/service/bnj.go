package service

import (
	"context"
	"encoding/json"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"go-common/app/job/main/dm2/model"
	"go-common/app/service/main/archive/api"
	arcMdl "go-common/app/service/main/archive/model/archive"
	filterMdl "go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

var (
	msgRegex = regexp.MustCompile(`^(\s|\xE3\x80\x80)*$`) // 全文仅空格

	_bnjDmMsgLen = 100

	_dateFormat = "2006-01-02 15:04:05"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func (s *Service) initBnj() {
	var err error
	if s.conf.BNJ.Aid <= 0 {
		return
	}
	s.bnjAid = s.conf.BNJ.Aid
	//bnj count
	if s.conf.BNJ.BnjCounter != nil {
		bnjSubAids := make(map[int64]struct{})
		for _, aid := range s.conf.BNJ.BnjCounter.SubAids {
			bnjSubAids[aid] = struct{}{}
		}
		s.bnjSubAids = bnjSubAids
	}
	// bnj danmu
	s.bnjVideos(context.TODO())
	s.bnjLiveConfig(context.TODO())
	go func() {
		ticker := time.NewTicker(time.Second * 30)
		for range ticker.C {
			s.bnjVideos(context.TODO())
			s.bnjLiveConfig(context.TODO())
		}
	}()
	s.bnjIgnoreRate = s.conf.BNJ.BnjLiveDanmu.IgnoreRate
	s.bnjIgnoreBeginTime = time.Duration(s.conf.BNJ.BnjLiveDanmu.IgnoreBegin)
	s.bnjIgnoreEndTime = time.Duration(s.conf.BNJ.BnjLiveDanmu.IgnoreEnd)
	s.bnjliveRoomID = s.conf.BNJ.BnjLiveDanmu.RoomID
	s.bnjUserLevel = s.conf.BNJ.BnjLiveDanmu.Level
	if s.bnjStart, err = time.ParseInLocation(_dateFormat, s.conf.BNJ.BnjLiveDanmu.Start, time.Now().Location()); err != nil {
		panic(err)
	}
	s.bnjCsmr = databus.New(s.conf.Databus.BnjCsmr)
	log.Info("bnj init start:%v room_id:%v", s.bnjStart.String(), s.conf.BNJ.BnjLiveDanmu.RoomID)
	go s.bnjProc()
}

func (s *Service) bnjProc() {
	var (
		err error
		c   = context.Background()
	)
	for {
		msg, ok := <-s.bnjCsmr.Messages()
		if !ok {
			log.Error("bnj bnjProc consumer exit")
			return
		}
		log.Info("bnj partition:%d,offset:%d,key:%s,value:%s", msg.Partition, msg.Offset, msg.Key, msg.Value)
		m := &model.LiveDanmu{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		if err = s.bnjLiveDanmu(c, m); err != nil {
			log.Error("bnj bnjLiveDanmu(msg:%+v),error(%v)", m, err)
			continue
		}
		if err = msg.Commit(); err != nil {
			log.Error("commit offset(%v) error(%v)", msg, err)
		}
	}
}

func (s *Service) bnjVideos(c context.Context) (err error) {
	var (
		videos []*model.Video
	)
	if videos, err = s.dao.Videos(c, s.bnjAid); err != nil {
		log.Error("bnj bnjVideos(aid:%v) error(%v)", s.bnjAid, err)
		return
	}
	if len(videos) >= 4 {
		videos = videos[:4]
	}
	for _, video := range videos {
		if err = s.syncBnjVideo(c, model.SubTypeVideo, video); err != nil {
			log.Error("bnj syncBnjVideo(video:%+v) error(%v)", video, err)
			return
		}
	}
	s.bnjArcVideos = videos
	return
}

func (s *Service) syncBnjVideo(c context.Context, tp int32, v *model.Video) (err error) {
	sub, err := s.dao.Subject(c, tp, v.Cid)
	if err != nil {
		return
	}
	if sub == nil {
		if v.XCodeState >= model.VideoXcodeHDFinish {
			if _, err = s.dao.AddSubject(c, tp, v.Cid, v.Aid, v.Mid, s.maxlimit(v.Duration), 0); err != nil {
				return
			}
		}
	} else {
		if sub.Mid != v.Mid {
			if _, err = s.dao.UpdateSubMid(c, tp, v.Cid, v.Mid); err != nil {
				return
			}
		}
	}
	return
}

// bnjDmCount laji bnj count
func (s *Service) bnjDmCount(c context.Context, sub *model.Subject, dm *model.DM) (err error) {
	var (
		dmid     int64
		pages    []*api.Page
		chosen   *api.Page
		choseSub *model.Subject
	)
	if _, ok := s.bnjSubAids[sub.Pid]; !ok {
		return
	}
	if pages, err = s.arcRPC.Page3(c, &arcMdl.ArgAid2{
		Aid:    s.bnjAid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}); err != nil {
		log.Error("bnjDmCount Page3(aid:%v) error(%v)", sub.Pid, err)
		return
	}
	if len(pages) <= 0 {
		return
	}
	idx := time.Now().Unix() % int64(len(pages))
	if chosen = pages[idx]; chosen == nil {
		return
	}
	if choseSub, err = s.subject(c, model.SubTypeVideo, chosen.Cid); err != nil {
		return
	}
	if dmid, err = s.genDMID(c); err != nil {
		log.Error("bnjDmCount genDMID() error(%v)", err)
		return
	}
	forkDM := &model.DM{
		ID:       dmid,
		Type:     model.SubTypeVideo,
		Oid:      chosen.Cid,
		Mid:      dm.Mid,
		Progress: int32((chosen.Duration + 1) * 1000),
		Pool:     dm.Pool,
		State:    model.StateAdminDelete,
		Ctime:    dm.Ctime,
		Mtime:    dm.Mtime,
		Content: &model.Content{
			ID:       dmid,
			FontSize: dm.Content.FontSize,
			Color:    dm.Content.Color,
			Mode:     dm.Content.Mode,
			IP:       dm.Content.IP,
			Plat:     dm.Content.Plat,
			Msg:      dm.Content.Msg,
			Ctime:    dm.Content.Ctime,
			Mtime:    dm.Content.Mtime,
		},
	}
	if dm.Pool == model.PoolSpecial {
		forkDM.ContentSpe = &model.ContentSpecial{
			ID:    dmid,
			Msg:   dm.ContentSpe.Msg,
			Ctime: dm.ContentSpe.Ctime,
			Mtime: dm.ContentSpe.Mtime,
		}
	}
	if err = s.bnjAddDM(c, choseSub, forkDM); err != nil {
		return
	}
	return
}

// bnjAddDM add dm index and content to db by transaction.
func (s *Service) bnjAddDM(c context.Context, sub *model.Subject, dm *model.DM) (err error) {
	if dm.State != model.StateAdminDelete {
		return
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		return
	}
	// special dm
	if dm.Pool == model.PoolSpecial && dm.ContentSpe != nil {
		if _, err = s.dao.TxAddContentSpecial(tx, dm.ContentSpe); err != nil {
			return tx.Rollback()
		}
	}
	if _, err = s.dao.TxAddContent(tx, dm.Oid, dm.Content); err != nil {
		return tx.Rollback()
	}
	if _, err = s.dao.TxAddIndex(tx, dm); err != nil {
		return tx.Rollback()
	}
	if _, err = s.dao.TxIncrSubjectCount(tx, sub.Type, sub.Oid, 1, 0, sub.Childpool); err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}

func (s *Service) genDMID(c context.Context) (dmid int64, err error) {
	if dmid, err = s.seqRPC.ID(c, s.seqArg); err != nil {
		log.Error("seqRPC.ID() error(%v)", err)
		return
	}
	return
}

// bnjLiveDanmu laji live to video
// TODO stime
func (s *Service) bnjLiveDanmu(c context.Context, liveDanmu *model.LiveDanmu) (err error) {
	var (
		cid, dmid int64
		progress  float64
	)
	// ignore time before
	if time.Since(s.bnjStart) < 0 {
		return
	}
	// limit
	if liveDanmu == nil || s.bnjliveRoomID <= 0 || s.bnjliveRoomID != liveDanmu.RoomID || liveDanmu.MsgType != model.LiveDanmuMsgTypeNormal {
		return
	}
	if liveDanmu.UserLevel < s.bnjUserLevel {
		return
	}
	if s.bnjIgnoreRate <= 0 || rand.Int63n(s.bnjIgnoreRate) != 0 {
		return
	}
	if cid, progress, err = s.pickBnjVideo(c, liveDanmu.Time); err != nil {
		return
	}
	// ignore illegal progress
	if progress <= 0 {
		return
	}
	if err = s.checkBnjDmMsg(c, liveDanmu.Content); err != nil {
		log.Error("bnj bnjLiveDanmu checkBnjDmMsg(liveDanmu:%+v) error(%v)", liveDanmu, err)
		return
	}
	if dmid, err = s.genDMID(c); err != nil {
		log.Error("bnj bnjLiveDanmu genDMID() error(%v)", err)
		return
	}
	now := time.Now().Unix()
	forkDM := &model.DM{
		ID:       dmid,
		Type:     model.SubTypeVideo,
		Oid:      cid,
		Mid:      liveDanmu.UID,
		Progress: int32(progress * 1000),
		Pool:     model.PoolNormal,
		State:    model.StateMonitorAfter,
		Ctime:    model.ConvertStime(time.Now()),
		Mtime:    model.ConvertStime(time.Now()),
		Content: &model.Content{
			ID:       dmid,
			FontSize: 25,
			Color:    16777215,
			Mode:     model.ModeRolling,
			Plat:     0,
			Msg:      liveDanmu.Content,
			Ctime:    xtime.Time(now),
			Mtime:    xtime.Time(now),
		},
	}
	if err = s.bnjCheckFilterService(c, forkDM); err != nil {
		log.Error("s.bnjCheckFilterService(%+v) error(%v)", forkDM, err)
		return
	}
	var (
		bs []byte
	)
	if bs, err = json.Marshal(forkDM); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", forkDM, err)
		return
	}
	act := &model.Action{
		Action: model.ActAddDM,
		Data:   bs,
	}
	if err = s.actionAct(c, act); err != nil {
		log.Error("s.actionAddDM(%+v) error(%v)", liveDanmu, err)
		return
	}
	return
}

func (s *Service) pickBnjVideo(c context.Context, timestamp int64) (cid int64, progress float64, err error) {
	var (
		idx   int
		video *model.Video
	)
	progress = float64(timestamp - s.bnjStart.Unix())
	for idx, video = range s.bnjArcVideos {
		if progress > float64(video.Duration) {
			progress = progress - float64(video.Duration)
			continue
		}
		// ignore p1 start
		if idx != 0 && progress < s.bnjIgnoreBeginTime.Seconds() {
			err = ecode.DMProgressTooBig
			return
		}
		if float64(video.Duration)-progress < s.bnjIgnoreEndTime.Seconds() {
			err = ecode.DMProgressTooBig
			return
		}
		if progress >= 0 {
			progress = progress + float64(rand.Int31n(1000)/1000)
		}
		cid = video.Cid
		return
	}
	err = ecode.DMProgressTooBig
	return
}

func (s *Service) bnjCheckFilterService(c context.Context, dm *model.DM) (err error) {
	var (
		filterReply *filterMdl.FilterReply
	)
	if filterReply, err = s.filterRPC.Filter(c, &filterMdl.FilterReq{
		Area:    "danmu",
		Message: dm.Content.Msg,
		Id:      dm.ID,
		Oid:     dm.Oid,
		Mid:     dm.Mid,
	}); err != nil {
		log.Error("checkFilterService(dm:%+v),err(%v)", dm, err)
		return
	}
	if filterReply.Level > 0 || filterReply.Limit == model.SpamBlack || filterReply.Limit == model.SpamOverflow {
		dm.State = model.StateFilter
		log.Info("bnj filter service delete(dmid:%d,data:+%v)", dm.ID, filterReply)
	}
	return
}

func (s *Service) checkBnjDmMsg(c context.Context, msg string) (err error) {
	var (
		msgLen = len([]rune(msg))
	)
	if msgRegex.MatchString(msg) { // 空白弹幕
		err = ecode.DMMsgIlleagel
		return
	}
	if msgLen > _bnjDmMsgLen {
		err = ecode.DMMsgTooLong
		return
	}
	if strings.Contains(msg, `\n`) || strings.Contains(msg, `/n`) {
		err = ecode.DMMsgIlleagel
		return
	}
	return
}

func (s *Service) bnjLiveConfig(c context.Context) (err error) {
	var (
		bnjConfig *model.BnjLiveConfig
		start     time.Time
	)
	if bnjConfig, err = s.dao.BnjConfig(c); err != nil {
		log.Error("bnjLiveConfig error current:%v err:%+v", time.Now().String(), err)
		return
	}
	if bnjConfig == nil {
		log.Error("bnjLiveConfig error current:%v bnjConfig nil", time.Now().String())
		return
	}
	if start, err = time.ParseInLocation(_dateFormat, bnjConfig.DanmuDtarTime, time.Now().Location()); err != nil {
		log.Error("bnjLiveConfig start time error current:%v config:%+v", time.Now().String(), bnjConfig)
		return
	}
	if bnjConfig.CommentID <= 0 || bnjConfig.RoomID <= 0 {
		log.Info("bnjLiveConfig illegal current:%v config:%+v", time.Now().String(), bnjConfig)
		return
	}
	s.bnjAid = bnjConfig.CommentID
	s.bnjliveRoomID = bnjConfig.RoomID
	s.bnjStart = start
	log.Info("bnjLiveConfig ok current:%v config:%+v", time.Now().String(), bnjConfig)
	return
}
