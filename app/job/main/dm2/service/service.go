package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"go-common/app/job/main/dm2/conf"
	"go-common/app/job/main/dm2/dao"
	"go-common/app/job/main/dm2/model"
	"go-common/app/job/main/dm2/model/oplog"
	arcCli "go-common/app/service/main/archive/api/gorpc"
	filterCli "go-common/app/service/main/filter/api/grpc/v1"
	seqMdl "go-common/app/service/main/seq-server/model"
	seqCli "go-common/app/service/main/seq-server/rpc/client"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"
	"go-common/library/xstr"
)

const (
	_routineSizeDefault = 10
	_chanSize           = 10240
	_batchSize          = 1000
	_maxUpRecent        = 1000
)

// Service service struct
type Service struct {
	conf   *conf.Config
	dao    *dao.Dao
	cache  *fanout.Fanout
	arcRPC *arcCli.Service2
	// seq serer
	seqArg            *seqMdl.ArgBusiness
	seqRPC            *seqCli.Service2
	indexCsmr         *databus.Databus
	subjectCsmr       *databus.Databus
	actionCsmr        *databus.Databus
	reportCsmr        *databus.Databus
	videoupCsmr       *databus.Databus
	subtitleAuditCsmr *databus.Databus
	flushMergeChan    []chan *model.Flush
	flushSegChan      []chan *model.FlushDMSeg
	dmRecentChan      []chan *model.DM
	routineSize       int
	realname          map[int64]int64 // key：分区id，value:cid，即该分区中大于cid的视频开启实名制
	// filter service
	filterRPC         filterCli.FilterClient
	maskMid           []int64
	dmOperationLogSvc *infoc.Infoc
	opsLogCh          chan *oplog.Infoc
	// bnj
	bnjAid             int64
	bnjSubAids         map[int64]struct{}
	bnjCsmr            *databus.Databus
	bnjliveRoomID      int64
	bnjStart           time.Time
	bnjIgnoreBeginTime time.Duration
	bnjIgnoreEndTime   time.Duration
	bnjArcVideos       []*model.Video
	bnjIgnoreRate      int64
	bnjUserLevel       int32
}

// New new service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:              c,
		dao:               dao.New(c),
		cache:             fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		arcRPC:            arcCli.New2(c.ArchiveRPC),
		seqArg:            &seqMdl.ArgBusiness{BusinessID: c.Seq.BusinessID, Token: c.Seq.Token},
		seqRPC:            seqCli.New2(c.SeqRPC),
		subjectCsmr:       databus.New(c.Databus.SubjectCsmr),
		indexCsmr:         databus.New(c.Databus.IndexCsmr),
		actionCsmr:        databus.New(c.Databus.ActionCsmr),
		reportCsmr:        databus.New(c.Databus.ReportCsmr),
		videoupCsmr:       databus.New(c.Databus.VideoupCsmr),
		subtitleAuditCsmr: databus.New(c.Databus.SubtitleAuditCsmr),
		routineSize:       c.RoutineSize,
		realname:          make(map[int64]int64),
		dmOperationLogSvc: infoc.New(c.Infoc2),
		opsLogCh:          make(chan *oplog.Infoc, 1024),
	}
	if c.RoutineSize <= 0 {
		s.routineSize = _routineSizeDefault
	}
	s.flushMergeChan = make([]chan *model.Flush, s.routineSize)
	s.flushSegChan = make([]chan *model.FlushDMSeg, s.routineSize)
	s.dmRecentChan = make([]chan *model.DM, s.routineSize)
	filterRPC, err := filterCli.NewClient(c.FliterRPC)
	if err != nil {
		panic(err)
	}
	s.filterRPC = filterRPC
	for idStr, cid := range conf.Conf.Realname {
		ids, err := xstr.SplitInts(idStr)
		if err != nil {
			panic(err)
		}
		for _, id := range ids {
			if _, ok := s.realname[id]; !ok {
				s.realname[id] = cid
			}
		}
	}
	// laji bnj
	if s.conf.BNJ != nil {
		//bnj count
		s.initBnj()
	}
	//消费DMReport-T消息
	go s.reportCsmproc()
	// 消费DMAction-T消息
	go s.actionCsmproc()
	// 消费DMSubject-T消息
	go s.subjectCsmproc()
	// 消费DMMeta-T消息
	go s.indexCsmproc()
	// 消费 Videoup2Bvc消息
	go s.videoupCsmrproc()
	// 消费 字幕 提交 消息
	go s.subtitleAuditProc()
	// 刷新全段弹幕
	for i := 0; i < s.routineSize; i++ {
		flushChan := make(chan *model.Flush, _chanSize)
		s.flushMergeChan[i] = flushChan
		go s.flushmergeproc(flushChan)
	}
	// 刷新分段弹幕
	for i := 0; i < s.routineSize; i++ {
		flushSegChan := make(chan *model.FlushDMSeg, _chanSize)
		s.flushSegChan[i] = flushSegChan
		go s.flushSegproc(flushSegChan)
	}
	// 异步处理创作中心最新弹幕列表缓存
	for i := 0; i < s.routineSize; i++ {
		recentChan := make(chan *model.DM, _chanSize)
		s.dmRecentChan[i] = recentChan
		go s.dmRecentproc(recentChan)
	}
	go s.transferProc()
	// 处理热门二级分类视频的弹幕蒙版
	go s.maskProc()
	// 刷新开启蒙版mid
	s.maskMid, _ = s.dao.MaskMids(context.TODO())
	log.Info("update mask mid(%v)", s.maskMid)
	go s.maskMidProc()
	// dm task
	go s.taskResProc()
	go s.taskDelProc()
	// oplog
	go s.oplogproc()
	return
}

// Ping check thrid resource.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

func (s *Service) subjectCsmproc() {
	var (
		err          error
		c            = context.TODO()
		regexSubject = regexp.MustCompile("dm_subject_[0-9]+")
	)
	for {
		msg, ok := <-s.subjectCsmr.Messages()
		if !ok {
			log.Error("subject binlog consumer exit")
			return
		}
		m := &model.BinlogMsg{}
		if err = json.Unmarshal(msg.Value, &m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		log.Info("%s", m)
		if regexSubject.MatchString(m.Table) {
			if err = s.trackSubject(c, m); err != nil {
				log.Error("s.trackSubject(%s) error(%v)", m, err)
				continue
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("commit offset(%v) error(%v)", msg, err)
		}
	}
}

func (s *Service) indexCsmproc() {
	var (
		err        error
		c          = context.TODO()
		regexIndex = regexp.MustCompile("dm_index_[0-9]+")
	)
	for {
		msg, ok := <-s.indexCsmr.Messages()
		if !ok {
			log.Error("index binlog consumer exit")
			return
		}
		m := &model.BinlogMsg{}
		if err = json.Unmarshal(msg.Value, &m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		log.Info("%s", m)
		if regexIndex.MatchString(m.Table) {
			if err = s.trackIndex(c, m); err != nil {
				log.Error("s.traceIndex(%s) error(%v)", m, err)
				continue
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("commit offset(%v) error(%v)", msg, err)
		}
	}
}

func (s *Service) flushmergeproc(flushChan chan *model.Flush) {
	var (
		flushs = make(map[int64]*model.Flush)
		ticker = time.NewTicker(60 * time.Second)
		err    error
	)
	for {
		select {
		case flush, ok := <-flushChan:
			if !ok {
				log.Error("action channel closed")
				return
			}
			if _, ok := flushs[flush.Oid]; !ok || flush.Force { // key不存在或者需要强制刷新的
				flushs[flush.Oid] = flush
			}
			if len(flushs) < _batchSize {
				continue
			}
		case <-ticker.C:
		}
		if len(flushs) > 0 {
			for _, flush := range flushs {
				if err = s.flushDmCache(context.TODO(), flush); err != nil {
					log.Error("action:flushmergeproc,flush:%+v,error(%v)", flush, err)
				}
			}
			flushs = make(map[int64]*model.Flush)
		}
	}
}

func keySegFlush(tp int32, oid, total, num int64) string {
	return fmt.Sprintf("f_%d_%d_%d_%d", tp, oid, total, num)
}

func (s *Service) flushSegproc(ch chan *model.FlushDMSeg) {
	var (
		key    string
		merge  = make(map[string]*model.FlushDMSeg)
		ticker = time.NewTicker(60 * time.Second)
		err    error
	)
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				log.Error("action channel closed")
				return
			}
			key = keySegFlush(msg.Type, msg.Oid, msg.Page.Total, msg.Page.Num)
			if _, ok := merge[key]; !ok || msg.Force { // key不存在或者需要强制刷新的
				merge[key] = msg
			}
			if len(merge) < _batchSize {
				continue
			}
		case <-ticker.C:
		}
		if len(merge) > 0 {
			for _, v := range merge {
				if err = s.flushDmSegCache(context.TODO(), v); err != nil {
					log.Error("action:flushSegproc,data:%+v,error(%v)", v, err)
					continue
				}
			}
			merge = make(map[string]*model.FlushDMSeg)
		}
	}
}

func (s *Service) actionCsmproc() {
	for {
		msg, ok := <-s.actionCsmr.Messages()
		if !ok {
			log.Error("action consumer exit")
			return
		}
		act := &model.Action{}
		err := json.Unmarshal(msg.Value, &act)
		if err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		start := time.Now()
		if err = s.actionAct(context.TODO(), act); err != nil {
			log.Error("action:%s,data:%s,error(%v)", act.Action, act.Data, err)
			continue
		}
		log.Info("partition:%d,offset:%d,key:%s,value:%s costing:%+v", msg.Partition, msg.Offset, msg.Key, msg.Value, time.Since(start))
		if err = msg.Commit(); err != nil {
			log.Error("commit offset(%v) error(%v)", msg, err)
		}
	}
}

// add recent dm
func (s *Service) dmRecentproc(dmChan chan *model.DM) {
	var (
		count int64
		c     = context.TODO()
	)
	for {
		dm, ok := <-dmChan
		if !ok {
			log.Error("recent dm channel is closed")
			return
		}
		sub, err := s.subject(c, dm.Type, dm.Oid)
		if err != nil {
			continue
		}
		if dm.State != model.StateNormal && dm.State != model.StateHide && dm.State != model.StateMonitorAfter {
			if err = s.dao.ZRemRecentDM(c, sub.Mid, dm.ID); err != nil {
				continue
			}
		} else {
			if dm.Content == nil {
				if dm.Content, err = s.dao.Content(c, dm.Oid, dm.ID); err != nil {
					continue
				}
			}
			if dm.Pool == model.PoolSpecial && dm.ContentSpe == nil {
				if dm.ContentSpe, err = s.dao.ContentSpecial(c, dm.ID); err != nil {
					continue
				}
			}
			if count, err = s.dao.AddRecentDM(c, sub.Mid, dm); err != nil {
				continue
			}
			if trimCnt := count - _maxUpRecent; trimCnt > 0 {
				if err = s.dao.TrimRecentDM(c, sub.Mid, trimCnt); err != nil {
					continue
				}
			}
		}
	}
}

func (s *Service) asyncAddRecent(c context.Context, dm *model.DM) {
	select {
	case s.dmRecentChan[dm.Oid%int64(s.routineSize)] <- dm:
	default:
		log.Warn("dm recent channel is full,dm(%+v)", dm)
	}
}

func (s *Service) reportCsmproc() {
	for {
		msg, ok := <-s.reportCsmr.Messages()
		if !ok {
			log.Error("report consumer exit")
			return
		}
		log.Info("partition:%d,offset:%d,value:%s", msg.Partition, msg.Offset, msg.Value)
		act := &model.ReportAction{}
		err := json.Unmarshal(msg.Value, &act)
		if err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		if diff := act.HideTime - time.Now().Unix(); diff > 0 {
			log.Info("action:%+v will be processed after %d seconds", act, diff)
			time.Sleep(time.Duration(diff) * time.Second)
		}
		if _, err = s.dao.DelDMHideState(context.TODO(), 1, act.Cid, act.Did); err != nil {
			log.Error("DelDMHideState(%+v) error(%v)", act, err)
		} else {
			log.Info("DelDMHideState(%+v) success ", act)
		}
		if err = msg.Commit(); err != nil {
			log.Error("commit offset(%v) error(%v)", msg, err)
		}
	}
}

func (s *Service) videoupCsmrproc() {
	var (
		err error
		c   = context.TODO()
	)
	for {
		msg, ok := <-s.videoupCsmr.Messages()
		if !ok {
			log.Error("videoup consumer exit")
			return
		}
		log.Info("partition:%d,offset:%d,key:%s,value:%s", msg.Partition, msg.Offset, msg.Key, msg.Value)
		m := &model.VideoupMsg{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		if m.Route == model.RouteSecondRound || m.Route == model.RouteAutoOpen ||
			m.Route == model.RouteForceSync || m.Route == model.RouteDelayOpen {
			if err = s.trackVideoup(c, m.Aid); err != nil {
				continue
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("commit offset(%v) error(%v)", msg, err)
		}
	}
}

func (s *Service) subtitleAuditProc() {
	var (
		err error
		c   = context.Background()
	)
	for {
		msg, ok := <-s.subtitleAuditCsmr.Messages()
		if !ok {
			log.Error("subtitle_audit consumer exit")
			return
		}
		log.Info("partition:%d,offset:%d,key:%s,value:%s", msg.Partition, msg.Offset, msg.Key, msg.Value)
		m := &model.SubtitleAuditMsg{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		if err = s.SubtitleFilter(c, m.Oid, m.SubtitleID); err != nil {
			log.Error("SubtitleFilter(oid:%v,subtitleID:%v),error(%v)", m.Oid, m.SubtitleID, err)
			continue
		}
		if err = msg.Commit(); err != nil {
			log.Error("commit offset(%v) error(%v)", msg, err)
		}
	}
}
