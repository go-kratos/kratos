package service

import (
	"context"
	"encoding/json"
	"runtime"
	"strconv"
	"sync"
	"time"

	actrpc "go-common/app/interface/main/activity/rpc/client"
	artmdl "go-common/app/interface/openplatform/article/model"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	"go-common/app/job/main/activity/conf"
	"go-common/app/job/main/activity/dao/bnj"
	"go-common/app/job/main/activity/dao/dm"
	"go-common/app/job/main/activity/dao/kfc"
	"go-common/app/job/main/activity/dao/like"
	kfcmdl "go-common/app/job/main/activity/model/kfc"
	l "go-common/app/job/main/activity/model/like"
	"go-common/app/job/main/activity/model/match"
	"go-common/app/service/main/account/api"
	arcapi "go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	comarcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/coin/api/gorpc"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/robfig/cron"
)

const (
	//_startVoteP       = 2
	//_startVote        = 3
	//_endingVote       = 4
	//_endVote          = 5
	//_goOn             = 6
	//_next             = 7
	_matchObjTable    = "act_matchs_object"
	_subjectTable     = "act_subject"
	_likesTable       = "likes"
	_likeContentTable = "like_content"
	_likeActionTable  = "like_action"
	//_vipActOrderTable = "vip_order_activity_record"
	_objectPieceSize = 100
	_retryTimes      = 3
	_typeArc         = "archive"
	_typeArt         = "article"
	_sharding        = 10
)

// Service service
type Service struct {
	c      *conf.Config
	dao    *like.Dao
	bnj    *bnj.Dao
	dm     *dm.Dao
	kfcDao *kfc.Dao
	// waiter
	waiter sync.WaitGroup
	closed bool
	// cache: type, upper
	// arc rpc
	arcRPC     *arcrpc.Service2
	coinRPC    *coin.Service
	actRPC     *actrpc.Service
	articleRPC *artrpc.Service
	//grpc
	accClient api.AccountClient
	// databus
	actSub *databus.Databus
	bnjSub *databus.Databus
	// vip binlog databus
	//vipSub      *databus.Databus
	kfcSub      *databus.Databus
	kfcActionCh []chan *kfcmdl.CouponMsg
	kfcShare    int
	subActionCh []chan *l.Action
	actionSM    []map[int64]*l.LastTmStat
	// bnj
	bnjMaxSecond    int64
	bnjLessSecond   int64
	bnjTimeFinish   int64
	bnjMsgFlagMap   map[int]int64
	bnjMsgFlagMu    sync.Mutex
	bnjWxMsgFlagMap map[int]int64
	bnjWxMsgFlagMu  sync.Mutex
	// cron
	cron *cron.Cron
}

// New is archive service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		dao:        like.New(c),
		dm:         dm.New(c),
		bnj:        bnj.New(c),
		kfcDao:     kfc.New(c),
		arcRPC:     arcrpc.New2(c.ArchiveRPC),
		articleRPC: artrpc.New(c.ArticleRPC),
		coinRPC:    coin.New(c.CoinRPC),
		actRPC:     actrpc.New(c.ActRPC),
		actSub:     databus.New(c.ActSub),
		bnjSub:     databus.New(c.BnjSub),
		//vipSub:     databus.New(c.VipSub),
		kfcSub: databus.New(c.KfcSub),
		cron:   cron.New(),
	}
	var err error
	if s.accClient, err = api.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	if s.c.Bnj2019.MsgSpec != "" {
		if err = s.cron.AddFunc(s.c.Bnj2019.MsgSpec, s.cronInformationMessage); err != nil {
			panic(err)
		}
		log.Info("cronInformationMessage init")
		s.cron.Start()
	}
	//time.Sleep(2 * time.Second)
	//subject, err := s.sub(context.Background(), c.Rule.BroadcastSid)
	//if err != nil {
	//	log.Error("error(%v)", err)
	//	return
	//}
	//log.Info("start-subject")
	//log.Info("subject(%v)", subject)
	//log.Info("end-subject")
	//if subject != nil {
	//	go s.genesis(subject)
	//}
	s.bnjMsgFlagMap = make(map[int]int64, len(bnjSteps))
	s.bnjWxMsgFlagMap = make(map[int]int64, len(bnjSteps))
	for _, step := range bnjSteps {
		s.bnjMsgFlagMap[step] = 0
		s.bnjWxMsgFlagMap[step] = 0
	}
	for i := 0; i < _sharding; i++ {
		s.subActionCh = append(s.subActionCh, make(chan *l.Action, 10240))
		s.actionSM = append(s.actionSM, map[int64]*l.LastTmStat{})
		s.waiter.Add(1)
		go s.actionDealProc(i)
	}
	//s.waiter.Add(1)
	//go s.vipCanal()
	s.waiter.Add(1)
	go s.consumeCanal()
	go s.subjectStat(s.c.Rule.ArcObjStatSid, _typeArc)
	go s.subjectStat(s.c.Rule.ArtObjStatSid, _typeArt)
	go s.kingStoryTotalStat(s.c.Rule.KingStorySid)
	go s.subsRankproc()
	go s.initBnjSecond()
	if runtime.NumCPU() <= 4 {
		s.kfcShare = 4
	} else if runtime.NumCPU() > 32 {
		s.kfcShare = 32
	} else {
		s.kfcShare = runtime.NumCPU()
	}
	for j := 0; j < s.kfcShare; j++ {
		s.kfcActionCh = append(s.kfcActionCh, make(chan *kfcmdl.CouponMsg, 10240))
		s.waiter.Add(1)
		go s.kfcActionDeal(j)
	}
	s.waiter.Add(1)
	go s.kfcCanal()
	return s
}

func (s *Service) likeArc(c context.Context, sub *l.Subject) (res *l.Subject, err error) {
	if sub != nil {
		if sub.ID == 0 {
			res = nil
		} else {
			res = sub
			var (
				ok   bool
				arcs map[int64]*arcapi.Arc
				aids []int64
			)
			for _, l := range res.List {
				aids = append(aids, l.Wid)
			}
			argAids := &comarcmdl.ArgAids2{
				Aids: aids,
			}
			if arcs, err = s.arcRPC.Archives3(c, argAids); err != nil {
				log.Error("s.arcRPC.Archives(arcAids:(%v), arcs), err(%v)", aids, err)
				return
			}
			for _, l := range res.List {
				if l.Archive, ok = arcs[l.Wid]; !ok {
					log.Info("s.arcs.wid:(%d), (%v)", l.Wid, ok)
					continue
				}
			}
		}
	}
	return
}

//func (s *Service) genesis(l *l.Subject) {
//	var (
//		index                         int
//		nowTime, stage, yes, no, next int64
//		err                           error
//		c                             = context.Background()
//	)
//	log.Info("st")
//	for {
//		if time.Now().Unix() >= l.Stime.Time().Unix() {
//			break
//		}
//		time.Sleep(time.Second)
//	}
//	lstime := map[string]interface{}{
//		"aid":   0,
//		"time":  0,
//		"index": 0,
//		"stage": 0,
//	}
//	go s.inLtime(lstime, l.ID)
//	for i, a := range l.List {
//		arg1 := &comarcmdl.ArgAid2{Aid: a.Archive.Aid}
//		arc, errRPC := s.arcRPC.Archive3(c, arg1)
//		if errRPC != nil {
//			log.Error("act-job s.arcRPC.Archive3(%v) error(%v)", arg1, errRPC)
//			errRPC = nil
//			continue
//		}
//		if arc.State < 0 && arc.State != -6 {
//			log.Error("act-job s.arcRPC.Archive3(%v) stat err", arg1)
//		}
//		index = i
//		nowTime = 0
//		stage = 0
//		next = 0
//		time.Sleep(time.Duration(l.Ltime) * time.Second)
//		log.Info("aid sleep")
//		aidTime := time.Now().Unix()
//		ltime := map[string]interface{}{
//			"aid":    a.Archive.Aid,
//			"time":   aidTime,
//			"index":  i,
//			"stage":  stage,
//			"title":  arc.Title,
//			"author": arc.Author.Name,
//			"tname":  arc.TypeName,
//		}
//		go s.inLtime(ltime, l.ID)
//		for {
//			fmt.Println(stage)
//			log.Info("s.stage{%d}", stage)
//			//演出开始
//			ltime := map[string]interface{}{
//				"aid":    a.Archive.Aid,
//				"time":   aidTime,
//				"index":  i,
//				"stage":  stage,
//				"title":  arc.Title,
//				"author": arc.Author.Name,
//				"tname":  arc.TypeName,
//			}
//			go s.inLtime(ltime, l.ID)
//			tp := l.Interval - l.Ltime
//			nowTime += tp
//			time.Sleep(time.Duration(tp) * time.Second)
//			//预备投票
//			sdm := &dmm.ActDM{Act: _startVoteP, Aid: a.Archive.Aid, Next: next, Yes: 0, No: 0, Stage: stage, Title: arc.Title, Author: arc.Author.Name, Tname: arc.TypeName}
//			go s.brodcast(sdm)
//			log.Info("act:1")
//			go s.dao.CreateSelection(c, a.Archive.Aid, stage)
//			nowTime += l.Ltime
//			time.Sleep(time.Duration(l.Ltime) * time.Second)
//			//投票开始
//			interval := &dmm.ActDM{Act: _startVote, Aid: a.Archive.Aid, Next: next, Yes: 0, No: 0, Stage: stage, Title: arc.Title, Author: arc.Author.Name, Tname: arc.TypeName}
//			go s.brodcast(interval)
//			log.Info("act:2")
//			tl := l.Tlimit - l.Ltime
//			nowTime += tl
//			time.Sleep(time.Duration(tl) * time.Second)
//			//投票预结束
//			intervalP := &dmm.ActDM{Act: _endingVote, Aid: a.Archive.Aid, Next: next, Yes: 0, No: 0, Stage: stage, Title: arc.Title, Author: arc.Author.Name, Tname: arc.TypeName}
//			log.Info("act:3")
//			go s.brodcast(intervalP)
//			nowTime += l.Ltime
//			time.Sleep(time.Duration(l.Ltime) * time.Second)
//			//投票结果
//			if yes, no, err = s.dao.Selection(c, a.Archive.Aid, stage); err != nil {
//				log.Error("s.dao.Selection() error(%v)", err)
//				return
//			}
//			goNext := true
//			if yes != 0 || no != 0 {
//				goNext = (float64(no)/float64(yes+no)*100 > 40)
//			}
//			if goNext {
//				next = 1
//			} else {
//				next = 0
//			}
//			intervalEnd := &dmm.ActDM{Act: _endVote, Aid: a.Archive.Aid, Next: next, Yes: yes, No: no, Stage: stage, Title: arc.Title, Author: arc.Author.Name, Tname: arc.TypeName}
//			go s.brodcast(intervalEnd)
//			log.Info("act:4")
//			nowTime += l.Ltime
//			time.Sleep(time.Duration(l.Ltime) * time.Second)
//			go s.inOnlinelog(c, l.ID, a.Archive.Aid, stage, yes, no)
//			if goNext {
//				tlimit := &dmm.ActDM{Act: _next, Aid: a.Archive.Aid, Next: next, Yes: yes, No: no, Stage: stage, Title: arc.Title, Author: arc.Author.Name, Tname: arc.TypeName}
//				go s.brodcast(tlimit)
//				log.Info("act:6")
//				nowTime += l.Ltime
//				ldtime := map[string]interface{}{
//					"aid":    0,
//					"time":   time.Now().Unix() + l.Ltime,
//					"index":  i + 1,
//					"stage":  0,
//					"title":  arc.Title,
//					"author": arc.Author.Name,
//					"tname":  arc.TypeName,
//				}
//				go s.inLtime(ldtime, l.ID)
//				time.Sleep(time.Duration(l.Ltime) * time.Second)
//				break
//			}
//			tlimit := &dmm.ActDM{Act: _goOn, Aid: a.Archive.Aid, Next: next, Yes: yes, No: no, Stage: stage, Title: arc.Title, Author: arc.Author.Name, Tname: arc.TypeName}
//			log.Info("bro:%v", tlimit)
//			go s.brodcast(tlimit)
//			log.Info("act:5")
//			//投票结果判断
//			stage++
//			if a.Archive.Duration-nowTime < 60 {
//				go s.inOnlinelog(c, l.ID, a.Archive.Aid, 100, 0, 0)
//				time.Sleep(time.Duration(a.Archive.Duration-nowTime) * time.Second)
//				break
//			}
//		}
//	}
//	ltime := map[string]interface{}{
//		"aid":    0,
//		"time":   time.Now().Unix(),
//		"index":  index + 1,
//		"stage":  0,
//		"title":  "",
//		"author": "",
//		"tname":  "",
//	}
//	go s.inLtime(ltime, l.ID)
//	log.Info("end")
//}
//
//func (s *Service) inOnlinelog(c context.Context, sid, aid, stage, yes, no int64) {
//	if row, err := s.dao.InOnlinelog(c, sid, aid, stage, yes, no); err != nil {
//		log.Error("s.dao.inOnlinelog, err(%v) row(%v)", err, row)
//	}
//}
//
//func (s *Service) inLtime(lt map[string]interface{}, sid int64) {
//	var v, err = json.Marshal(lt)
//	if err != nil {
//		log.Error("s.genesis.inLtime.json.Marshal(dm:(%v)), err(%v)", v, err)
//		return
//	}
//	s.dao.RbSet(context.Background(), "ltime:"+strconv.FormatInt(sid, 10), v)
//}
//
//func (s *Service) brodcast(d *dmm.ActDM) {
//	var ds, err = json.Marshal(d)
//	if err != nil {
//		log.Error("s.genesis.json.Marshal(dm:(%v)), err(%v)", d, err)
//		return
//	}
//	var m = &dmm.Broadcast{
//		RoomID: s.c.Rule.BroadcastCid,
//		CMD:    dmm.BroadcastCMDACT,
//		Info:   ds,
//	}
//	s.dm.Broadcast(context.Background(), m)
//}
//
//func (s *Service) sub(c context.Context, sid int64) (res *l.Subject, err error) {
//	var (
//		eg errgroup.Group
//		ls []*l.Like
//	)
//	eg.Go(func() (err error) {
//		res, err = s.dao.Subject(c, sid)
//		return
//	})
//	eg.Go(func() (err error) {
//		ls, err = s.dao.Like(c, sid)
//		return
//	})
//	if err = eg.Wait(); err != nil {
//		log.Error("eg.Wait error(%v)", err)
//		return
//	}
//	if res != nil {
//		res.List = ls
//	}
//	if res, err = s.likeArc(c, res); err != nil {
//		return
//	}
//	return
//}

//func (s *Service) vipCanal() {
//	defer s.waiter.Done()
//	var c = context.Background()
//	for {
//		msg, ok := <-s.vipSub.Messages()
//		if !ok {
//			log.Info("databus:activity-job vip binlog consumer exit!")
//			return
//		}
//		msg.Commit()
//		m := &match.Message{}
//		if err := json.Unmarshal(msg.Value, m); err != nil {
//			log.Error("json.Unmarshal(%s) error(%+v)", msg.Value, err)
//			continue
//		}
//		switch m.Table {
//		case _vipActOrderTable:
//			if m.Action == match.ActInsert {
//				s.addElemeLottery(c, m.New)
//			}
//		}
//		log.Info("vipCanal key:%s partition:%d offset:%d table:%s", msg.Key, msg.Partition, msg.Offset, m.Table)
//	}
//}

func (s *Service) consumeCanal() {
	defer s.waiter.Done()
	var c = context.Background()
	for {
		msg, ok := <-s.actSub.Messages()
		if !ok {
			log.Info("databus: activity-job binlog consumer exit!")
			return
		}
		msg.Commit()
		m := &match.Message{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%+v)", msg.Value, err)
			continue
		}
		switch m.Table {
		case _matchObjTable:
			if m.Action == match.ActUpdate {
				s.upMatchUser(c, m.New, m.Old)
			}
		case _subjectTable:
			if m.Action == match.ActInsert || m.Action == match.ActUpdate {
				s.upSubject(c, m.New)
			} else if m.Action == match.ActDelete {
				s.upSubject(c, m.Old)
			}
		case _likesTable:
			if m.Action == match.ActInsert {
				s.AddLike(c, m.New)
			} else if m.Action == match.ActUpdate {
				s.UpLike(c, m.New, m.Old)
			} else if m.Action == match.ActDelete {
				s.DelLike(c, m.Old)
			}
		case _likeContentTable:
			if m.Action == match.ActDelete {
				s.upLikeContent(c, m.Old)
			} else {
				s.upLikeContent(c, m.New)
			}
		case _likeActionTable:
			if m.Action == match.ActInsert {
				s.actionProc(c, m.New)
			}
		}
		log.Info("consumeCanal key:%s partition:%d offset:%d table:%s", msg.Key, msg.Partition, msg.Offset, m.Table)
	}
}

func (s *Service) kfcCanal() {
	defer s.waiter.Done()
	for {
		if s.closed {
			return
		}
		msg, ok := <-s.kfcSub.Messages()
		if !ok {
			log.Info("databus: activity-job binlog consumer exit!")
			return
		}
		msg.Commit()
		m := &kfcmdl.CouponMsg{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("kfcCanal:json.Unmarshal(%s) error(%+v)", msg.Value, err)
			continue
		}
		j := m.CouponID % int64(s.kfcShare)
		select {
		case s.kfcActionCh[j] <- m:
		default:
			log.Info("kfcCanal cache full (%d)", j)
		}
		log.Info("kfcCanal key:%s partition:%d offset:%d value %s goroutine(%d) all(%d)", msg.Key, msg.Partition, msg.Offset, msg.Value, j, s.kfcShare)
	}
}

func (s *Service) subjectStat(sid int64, typ string) {
	var err error
	for {
		if s.closed {
			return
		}
		var (
			statLike int64
			likeCnt  int
			likes    []*l.Like
		)
		if sid <= 0 {
			log.Warn("conf sid == 0 typ(%s)", typ)
			time.Sleep(time.Duration(s.c.Interval.ObjStatInterval))
			continue
		}
		if likeCnt, err = s.dao.LikeCnt(context.Background(), sid); err != nil {
			log.Error("s.dao.LikeCnt(sid:%d) error(%v)", sid, err)
			time.Sleep(time.Duration(s.c.Interval.ObjStatInterval))
			continue
		}
		if likeCnt == 0 {
			log.Warn("s.dao.LikeCnt(sid:%d) likeCnt == 0", sid)
			time.Sleep(time.Duration(s.c.Interval.ObjStatInterval))
			continue
		}
		for i := 0; i < likeCnt; i += _objectPieceSize {
			if likes, err = s.likeList(context.Background(), sid, i, _objectPieceSize, _retryTimes); err != nil {
				log.Error("objectStatproc s.likeList(%d,%d,%d) error(%+v)", sid, i, _objectPieceSize, err)
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				var aids []int64
				for _, v := range likes {
					if v.Wid > 0 {
						aids = append(aids, v.Wid)
					}
				}
				switch typ {
				case _typeArc:
					var arcs map[int64]*arcapi.Arc
					if arcs, err = s.arcs(context.Background(), aids, _retryTimes); err != nil {
						log.Error("objectStatproc s.arcs(%v) error(%v)", aids, err)
						time.Sleep(100 * time.Millisecond)
						continue
					} else {
						for _, aid := range aids {
							if arc, ok := arcs[aid]; ok && arc.IsNormal() {
								statLike += int64(arc.Stat.Like)
							} else {
								likeCnt--
							}
						}
					}
				case _typeArt:
					var arts map[int64]*artmdl.Meta
					if arts, err = s.arts(context.Background(), aids, _retryTimes); err != nil {
						log.Error("objectStatproc s.arcs(%v) error(%v)", aids, err)
						time.Sleep(100 * time.Microsecond)
						continue
					} else {
						for _, aid := range aids {
							if art, ok := arts[aid]; ok && art.IsNormal() {
								statLike += art.Stats.Like
							} else {
								likeCnt--
							}
						}
					}
				}
			}
		}
		if err = s.setSubjectStat(context.Background(), sid, &l.SubjectTotalStat{SumLike: statLike}, likeCnt, _retryTimes); err != nil {
			log.Error("objectStatproc s.setObjectStat(%d,%d) error(%+v)", sid, statLike, err)
			time.Sleep(time.Duration(s.c.Interval.ObjStatInterval))
			continue
		}
		time.Sleep(time.Duration(s.c.Interval.ObjStatInterval))
	}
}

func (s *Service) likeList(c context.Context, sid int64, offset, limit, retryCnt int) (list []*l.Like, err error) {
	for i := 0; i < retryCnt; i++ {
		if list, err = s.dao.LikeList(c, sid, offset, limit); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func (s *Service) webDataList(c context.Context, vid int64, offset, limit, retryCnt int) (list []*l.WebData, err error) {
	for i := 0; i < retryCnt; i++ {
		if list, err = s.dao.WebDataList(c, vid, offset, limit); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func (s *Service) arts(c context.Context, aids []int64, retryCnt int) (arcs map[int64]*artmdl.Meta, err error) {
	for i := 0; i < retryCnt; i++ {
		if arcs, err = s.articleRPC.ArticleMetas(c, &artmdl.ArgAids{Aids: aids}); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func (s *Service) kingStoryTotalStat(vid int64) {
	var err error
	for {
		if s.closed {
			return
		}
		var (
			statView, statLike, statFav, StatCoin int64
			likeCnt                               int
			likes                                 []*l.WebData
		)
		if vid <= 0 {
			log.Warn("conf vid == 0")
			time.Sleep(time.Duration(s.c.Interval.KingStoryInterval))
			continue
		}
		if likeCnt, err = s.dao.WebDataCnt(context.Background(), vid); err != nil {
			log.Error("kingStoryTotalStat s.dao.WebDataCnt(sid:%d) error(%v)", vid, err)
			time.Sleep(time.Duration(s.c.Interval.KingStoryInterval))
			continue
		}
		if likeCnt == 0 {
			log.Warn("kingStoryTotalStat s.dao.LikeCnt(sid:%d) likeCnt == 0", vid)
			time.Sleep(time.Duration(s.c.Interval.KingStoryInterval))
			continue
		}
		for i := 0; i < likeCnt; i += _objectPieceSize {
			if likes, err = s.webDataList(context.Background(), vid, i, _objectPieceSize, _retryTimes); err != nil {
				log.Error("kingStoryTotalStat s.webDataList(%d,%d,%d) error(%+v)", vid, i, _objectPieceSize, err)
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				var (
					aids      []int64
					aidStruct = new(struct {
						Aid string `json:"AID"`
					})
				)
				for _, v := range likes {
					if v.Data != "" {
						if e := json.Unmarshal([]byte(v.Data), &aidStruct); e != nil {
							log.Warn("kingStoryTotalStat json.Unmarshal(%s) error(%v)", v.Data, e)
							continue
						}
						if aid, e := strconv.ParseInt(aidStruct.Aid, 10, 64); e != nil {
							log.Warn("kingStoryTotalStat strconv.ParseInt(%s) error(%v)", aidStruct, e)
							continue
						} else {
							aids = append(aids, aid)
						}
					}
				}
				var arcs map[int64]*arcapi.Arc
				if arcs, err = s.arcs(context.Background(), aids, _retryTimes); err != nil {
					log.Error("kingStoryTotalStat s.arcs(%v) error(%v)", aids, err)
					time.Sleep(100 * time.Millisecond)
					continue
				} else {
					for _, aid := range aids {
						if arc, ok := arcs[aid]; ok && arc.IsNormal() {
							statView += int64(arc.Stat.View)
							statLike += int64(arc.Stat.Like)
							statFav += int64(arc.Stat.Fav)
							StatCoin += int64(arc.Stat.Coin)
						}
					}
				}
			}
		}
		if err = s.setSubjectStat(context.Background(), vid, &l.SubjectTotalStat{SumView: statView, SumLike: statLike, SumFav: statFav, SumCoin: StatCoin}, likeCnt, _retryTimes); err != nil {
			log.Error("kingStoryTotalStat s.setObjectStat(%d,%d) error(%+v)", vid, statLike, err)
			time.Sleep(time.Duration(s.c.Interval.KingStoryInterval))
			continue
		}
		time.Sleep(time.Duration(s.c.Interval.KingStoryInterval))
	}
}

func (s *Service) setSubjectStat(c context.Context, sid int64, stat *l.SubjectTotalStat, count, retryCnt int) (err error) {
	for i := 0; i < retryCnt; i++ {
		if err = s.dao.SetObjectStat(c, sid, stat, count); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

// Ping reports the heath of services.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close kafaka consumer close.
func (s *Service) Close() (err error) {
	defer s.waiter.Wait()
	s.closed = true
	if s.bnjTimeFinish == 0 {
		s.bnj.AddCacheLessTime(context.Background(), s.bnjLessSecond)
	}
	s.cron.Stop()
	s.dao.Close()
	s.actSub.Close()
	s.bnjSub.Close()
	//s.vipSub.Close()
	s.kfcSub.Close()
	s.closed = true
	return
}
