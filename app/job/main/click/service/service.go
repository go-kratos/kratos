package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-common/app/job/main/click/conf"
	"go-common/app/job/main/click/dao"
	"go-common/app/job/main/click/model"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/queue/databus"
)

const (
	_unLock = 0
	_locked = 1
)

// Service struct
type Service struct {
	c  *conf.Config
	db *dao.Dao
	// archive
	reportMergeSub   *databus.Databus
	statViewPub      *databus.Databus
	chanWg           sync.WaitGroup
	redis            *redis.Pool
	cliChan          []chan *model.ClickMsg
	closed           bool
	maxAID           int64
	gotMaxAIDTime    int64
	lockedMap        []int64
	currentLockedIdx int64
	// aid%50[aid[plat[cnt]]]
	aidMap []map[int64]*model.ClickInfo
	// send databus chan
	busChan     chan *model.StatMsg
	bigDataChan chan *model.BigDataMsg
	// forbid cache
	forbids    map[int64]map[int8]*model.Forbid
	forbidMids map[int64]struct{}
	// epid to aid map
	eTam            map[int64]int64
	etamMutex       sync.RWMutex
	infoc2          *infoc.Infoc
	arcRPC          *arcrpc.Service2
	arcDurWithMutex struct {
		Durations map[int64]*model.ArcDuration
		Mutex     sync.RWMutex
	}
	allowPlat     map[int8]struct{}
	bnjListAidMap map[int64]struct{}
}

// New is archive service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		arcRPC:         arcrpc.New2(c.ArchiveRPC),
		redis:          redis.NewPool(c.Redis),
		db:             dao.New(c),
		busChan:        make(chan *model.StatMsg, 10240),
		bigDataChan:    make(chan *model.BigDataMsg, 10240),
		reportMergeSub: databus.New(c.ReportMergeDatabus),
		statViewPub:    databus.New(c.StatViewPub),
		infoc2:         infoc.New(c.Infoc2),
		allowPlat:      make(map[int8]struct{}),
	}
	s.allowPlat[model.PlatForWeb] = struct{}{}
	s.allowPlat[model.PlatForH5] = struct{}{}
	s.allowPlat[model.PlatForOuter] = struct{}{}
	s.allowPlat[model.PlatForIos] = struct{}{}
	s.allowPlat[model.PlatForAndroid] = struct{}{}
	s.allowPlat[model.PlatForAndroidTV] = struct{}{}
	s.allowPlat[model.PlatForAutoPlayIOS] = struct{}{}
	s.allowPlat[model.PlafForAutoPlayInlineIOS] = struct{}{}
	s.allowPlat[model.PlatForAutoPlayAndroid] = struct{}{}
	s.allowPlat[model.PlatForAutoPlayInlineAndroid] = struct{}{}
	s.arcDurWithMutex.Durations = make(map[int64]*model.ArcDuration)
	s.loadConf()
	go s.confproc()
	go s.releaseAIDMap()
	for i := int64(0); i < s.c.ChanNum; i++ {
		s.aidMap = append(s.aidMap, make(map[int64]*model.ClickInfo, 300000))
		s.cliChan = append(s.cliChan, make(chan *model.ClickMsg, 256))
		s.lockedMap = append(s.lockedMap, _unLock)
	}
	for i := int64(0); i < s.c.ChanNum; i++ {
		s.chanWg.Add(1)
		go s.cliChanProc(i)
	}
	for i := 0; i < 10; i++ {
		s.chanWg.Add(1)
		go s.sendStat()
	}
	s.chanWg.Add(1)
	go s.sendBigDataMsg()
	s.chanWg.Add(1)
	go s.reportMergeSubConsumer()
	return s
}

func (s *Service) reportMergeSubConsumer() {
	defer s.chanWg.Done()
	msgs := s.reportMergeSub.Messages()
	for {
		msg, ok := <-msgs
		if !ok || s.closed {
			log.Info("s.reportMergeSub is closed")
			return
		}
		msg.Commit()
		var (
			sbs [][]byte
			err error
		)
		if err = json.Unmarshal(msg.Value, &sbs); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", msg.Value, err)
			continue
		}
		for _, bs := range sbs {
			var (
				click *model.ClickMsg
				allow bool
				now   = time.Now().Unix()
			)
			log.Info("split merged message(%s)", strings.Replace(string(bs), "\001", "|", -1))
			if click, err = s.checkMsgIllegal(bs); err != nil {
				log.Error("s.checkMsgIllegal(%s) error(%v)", strings.Replace(string(bs), "\001", "|", -1), err)
				continue
			}
			if s.maxAID > 0 && now-s.gotMaxAIDTime < 120 {
				allow = s.maxAID+300 > click.AID
			}
			if !allow {
				log.Error("maxAid(%d) currentAid(%d) not allow!!!!", s.maxAID, click.AID)
				continue
			}
			log.Info("merge consumer(%d) append to chan", click.AID)
			s.cliChan[click.AID%s.c.ChanNum] <- click
		}
	}
}

func (s *Service) loadConf() {
	var (
		forbids     map[int64]map[int8]*model.Forbid
		bnjListAids = make(map[int64]struct{})
		forbidMids  map[int64]struct{}
		etam        map[int64]int64
		maxAID      int64
		err         error
	)
	for _, aid := range s.c.BnjListAids {
		bnjListAids[aid] = struct{}{}
	}
	s.bnjListAidMap = bnjListAids
	if forbidMids, err = s.db.ForbidMids(context.Background()); err == nil {
		s.forbidMids = forbidMids
		log.Info("forbid mids(%d)", len(forbidMids))
	}
	if forbids, err = s.db.Forbids(context.TODO()); err == nil {
		s.forbids = forbids
		log.Info("forbid av(%d)", len(forbids))
	}
	if maxAID, err = s.db.MaxAID(context.TODO()); err == nil {
		s.maxAID = maxAID
		s.gotMaxAIDTime = time.Now().Unix()
	}
	if etam, err = s.db.LoadAllBangumi(context.TODO()); err == nil {
		s.etamMutex.Lock()
		s.eTam = etam
		s.etamMutex.Unlock()
	}
}

func (s *Service) releaseAIDMap() {
	for {
		time.Sleep(5 * time.Minute)
		now := time.Now()
		if (now.Hour() > 1 && now.Hour() < 6) || (now.Hour() == 6 && now.Minute() < 30) { // 2:00 to 6:30
			if s.currentLockedIdx < int64(len(s.aidMap)) {
				atomic.StoreInt64(&s.lockedMap[s.currentLockedIdx], _locked)
			}
			s.currentLockedIdx++
			continue
		}
		s.currentLockedIdx = 0
	}
}

func (s *Service) confproc() {
	for {
		time.Sleep(1 * time.Minute)
		s.loadConf()
	}
}

func (s *Service) sendBigDataMsg() {
	defer s.chanWg.Done()
	for {
		var (
			msg   *model.BigDataMsg
			msgBs []byte
			ok    bool
			err   error
			infos []interface{}
		)
		if msg, ok = <-s.bigDataChan; !ok {
			break
		}
		infos = append(infos, strconv.FormatInt(int64(msg.Tp), 10))
		for _, v := range strings.Split(msg.Info, "\001") {
			infos = append(infos, v)
		}
		log.Info("truly used %+v", infos)
		if err = s.infoc2.Info(infos...); err != nil {
			log.Error("s.infoc2.Info(%s) error(%v)", msgBs, err)
			continue
		}
	}
}

func (s *Service) sendStat() {
	defer s.chanWg.Done()
	for {
		var (
			msg *model.StatMsg
			ok  bool
			c   = context.TODO()
			err error
			key string
		)
		if msg, ok = <-s.busChan; !ok {
			break
		}
		key = strconv.FormatInt(msg.AID, 10)
		vmsg := &model.StatViewMsg{Type: "archive", ID: msg.AID, Count: msg.Click, Ts: time.Now().Unix()}
		if err = s.statViewPub.Send(c, key, vmsg); err != nil {
			log.Error("s.statViewPub.Send(%d, %+v) error(%v)", msg.AID, vmsg, err)
		}
	}
}

func (s *Service) cliChanProc(i int64) {
	defer s.chanWg.Done()
	var (
		cli     *model.ClickMsg
		cliChan = s.cliChan[i]
		ok      bool
	)
	for {
		if cli, ok = <-cliChan; !ok {
			s.countClick(context.TODO(), nil, i)
			return
		}
		var (
			rtype int8
			err   error
			c     = context.TODO()
		)
		if rtype, err = s.isAllow(c, cli); err != nil {
			log.Error("cliChanProc Err %v", err)
		}
		select {
		case s.bigDataChan <- &model.BigDataMsg{Info: string(cli.KafkaBs), Tp: rtype}:
		default:
			log.Error("s.bigDataChan is full")
		}
		if rtype == model.LogTypeForTurly {
			s.countClick(context.TODO(), cli, i)
		}
	}
}

func (s *Service) checkMsgIllegal(msg []byte) (click *model.ClickMsg, err error) {
	var (
		aid        int64
		clickMsg   []string
		plat       int64
		did        string
		buvid      string
		mid        int64
		lv         int64
		ctime      int64
		stime      int64
		epid       int64
		ip         string
		seasonType int
		userAgent  string
	)
	clickMsg = strings.Split(string(msg), "\001")
	if len(clickMsg) < 10 {
		err = errors.New("click msg error")
		return
	}
	if aid, err = strconv.ParseInt(clickMsg[1], 10, 64); err != nil {
		err = fmt.Errorf("aid(%s) error", clickMsg[1])
		return
	}
	if aid <= 0 {
		err = fmt.Errorf("wocao aid(%s) error", clickMsg[1])
		return
	}
	if plat, err = strconv.ParseInt(clickMsg[0], 10, 64); err != nil {
		err = fmt.Errorf("plat(%s) error", clickMsg[0])
		return
	}

	if _, ok := s.allowPlat[int8(plat)]; !ok {
		err = fmt.Errorf("plat(%d) is illegal", plat)
		return
	}
	userAgent = clickMsg[10]
	did = clickMsg[8]
	if did == "" {
		err = fmt.Errorf("bvID(%s) is illegal", clickMsg[8])
		return
	}
	buvid = clickMsg[11]
	if clickMsg[4] != "" && clickMsg[4] != "0" {
		if mid, err = strconv.ParseInt(clickMsg[4], 10, 64); err != nil {
			err = fmt.Errorf("mid(%s) is illegal", clickMsg[4])
			return
		}
	}
	if clickMsg[5] != "" {
		if lv, err = strconv.ParseInt(clickMsg[5], 10, 64); err != nil {
			err = fmt.Errorf("lv(%s) is illegal", clickMsg[5])
			return
		}
	}
	if ctime, err = strconv.ParseInt(clickMsg[6], 10, 64); err != nil {
		err = fmt.Errorf("ctime(%s) is illegal", clickMsg[6])
		return
	}
	if stime, err = strconv.ParseInt(clickMsg[7], 10, 64); err != nil {
		err = fmt.Errorf("stime(%s) is illegal", clickMsg[7])
		return
	}
	if ip = clickMsg[9]; ip == "" {
		err = errors.New("ip is illegal")
		return
	}
	if clickMsg[17] != "" {
		if epid, err = strconv.ParseInt(clickMsg[17], 10, 64); err != nil {
			err = fmt.Errorf("epid(%s) is illegal", clickMsg[17])
			return
		}
		if clickMsg[15] != "null" {
			if seasonType, err = strconv.Atoi(clickMsg[15]); err != nil {
				err = fmt.Errorf("seasonType(%s) is illegal", clickMsg[15])
				return
			}
		}
	}
	if strings.Contains(userAgent, "(auto_play)") ||
		strings.Contains(userAgent, "(inline_play_heartbeat)") ||
		strings.Contains(userAgent, "(inline_play_to_view)") || // to remove auto_play & inline_play_heartbeat
		strings.Contains(userAgent, "(played_time_enough)") {
		if did, err = s.getRealDid(context.TODO(), buvid, aid); err != nil || did == "" {
			err = fmt.Errorf("bvid(%s) dont have did", buvid)
			return
		}
		did = buvid
		msg = []byte(strings.Replace(string(msg), buvid, did, 1))
	}
	click = &model.ClickMsg{
		Plat:       int8(plat),
		AID:        aid,
		MID:        mid,
		Lv:         int8(lv),
		CTime:      ctime,
		STime:      stime,
		Did:        did,
		Buvid:      buvid,
		IP:         ip,
		KafkaBs:    msg,
		EpID:       epid,
		SeasonType: seasonType,
		UserAgent:  userAgent,
	}
	return
}

// ArcDuration return archive duration, manager local cache
func (s *Service) ArcDuration(c context.Context, aid int64) (duration int64) {
	var (
		ok     bool
		arcDur *model.ArcDuration
		now    = time.Now().Unix()
		err    error
	)
	// duration default
	duration = s.c.CacheConf.PGCReplayTime
	s.arcDurWithMutex.Mutex.RLock()
	arcDur, ok = s.arcDurWithMutex.Durations[aid]
	s.arcDurWithMutex.Mutex.RUnlock()
	if ok && now-arcDur.GotTime > s.c.CacheConf.ArcUpCacheTime {
		duration = arcDur.Duration
		return
	}
	var arc *api.Arc
	if arc, err = s.arcRPC.Archive3(c, &arcmdl.ArgAid2{Aid: aid}); err != nil {
		// just log
		log.Error("s.arcRPC.Archive3(%d) error(%v)", aid, err)
	} else {
		duration = arc.Duration
	}
	s.arcDurWithMutex.Mutex.Lock()
	s.arcDurWithMutex.Durations[aid] = &model.ArcDuration{Duration: duration, GotTime: now}
	s.arcDurWithMutex.Mutex.Unlock()
	return
}

// Close kafaka consumer close.
func (s *Service) Close() (err error) {
	s.closed = true
	time.Sleep(time.Second)
	for i := 0; i < len(s.cliChan); i++ {
		close(s.cliChan[i])
	}
	close(s.bigDataChan)
	s.chanWg.Wait()
	s.statViewPub.Close()
	return
}
