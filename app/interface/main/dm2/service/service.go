package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/app/interface/main/dm2/conf"
	"go-common/app/interface/main/dm2/dao"
	"go-common/app/interface/main/dm2/model"
	"go-common/app/interface/main/dm2/model/oplog"
	accountCli "go-common/app/service/main/account/api"
	arcCli "go-common/app/service/main/archive/api/gorpc"
	assMdl "go-common/app/service/main/assist/model/assist"
	assCli "go-common/app/service/main/assist/rpc/client"
	coinCli "go-common/app/service/main/coin/api/gorpc"
	figureCli "go-common/app/service/main/figure/rpc/client"
	filterCli "go-common/app/service/main/filter/api/grpc/v1"
	locCli "go-common/app/service/main/location/rpc/client"
	memberCli "go-common/app/service/main/member/api/gorpc"
	relCli "go-common/app/service/main/relation/rpc/client"
	seqMdl "go-common/app/service/main/seq-server/model"
	seqCli "go-common/app/service/main/seq-server/rpc/client"
	spyCli "go-common/app/service/main/spy/rpc/client"
	thumbupApi "go-common/app/service/main/thumbup/api"

	ugcPayCli "go-common/app/service/main/ugcpay/api/grpc/v1"
	seasonCli "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/sync/pipeline/fanout"
	"go-common/library/xstr"

	"golang.org/x/sync/singleflight"
)

func keySubject(tp int32, oid int64) string {
	return fmt.Sprintf("sub_local_%d_%d", tp, oid)
}

func keyXML(tp int32, oid int64) string {
	return fmt.Sprintf("%d_%d", tp, oid)
}

func keySeg(tp int32, oid, cnt, num int64) string {
	return fmt.Sprintf("%d_%d_%d_%d", tp, oid, cnt, num)
}

func keyDuration(tp int32, oid int64) string {
	return fmt.Sprintf("d_%d_%d", tp, oid)
}

type broadcast struct {
	Aid int64
	Rnd int64
	*model.DM
}

// Service dm2 service
type Service struct {
	conf          *conf.Config
	dao           *dao.Dao
	arcRPC        *arcCli.Service2
	accountRPC    accountCli.AccountClient
	assRPC        *assCli.Service
	coinRPC       *coinCli.Service
	relRPC        *relCli.Service
	thumbupRPC    thumbupApi.ThumbupClient
	locationRPC   *locCli.Service
	cache         *fanout.Fanout
	realname      map[int64]int64 // key：分区id，value:cid，即该分区中大于cid的视频开启实名制
	singleGroup   singleflight.Group
	arcTypes      map[int16]int16
	broadcastChan chan *broadcast
	// seq serer
	seqDmArg       *seqMdl.ArgBusiness
	seqSubtitleArg *seqMdl.ArgBusiness
	seqRPC         *seqCli.Service2
	// dm xml and dm seg local cache
	localCache    map[string][]byte
	assistLogChan chan *assMdl.ArgAssistLogAdd
	// send operation log with infoc2
	dmOperationLogSvc *infoc.Infoc
	opsLogCh          chan *oplog.Infoc
	// dm monitor merge proc
	moniOidMap map[int64]struct{}
	oidLock    sync.Mutex
	// subtitle singleFlight
	subtitleSingleGroup singleflight.Group
	subtitleLans        model.SubtitleLans
	// block
	memberRPC *memberCli.Service
	// filter
	filterRPC filterCli.FilterClient
	//figure
	figureRPC *figureCli.Service
	// ugc pay
	ugcPayRPC ugcPayCli.UGCPayClient
	// spy RPC
	spyRPC *spyCli.Service
	// season
	seasonRPC seasonCli.SeasonClient
	// garbageDanmu
	garbageDanmu bool
	// broadcast limit
	broadcastLimit         int
	broadcastlimitInterval int
	// view localcache
	localViewCache map[string]*model.ViewDm
	// bnj shield
	aidSheild  map[int64]struct{}
	midsSheild map[int64]struct{}
}

// New return a service instance.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:                   c,
		dao:                    dao.New(c),
		arcRPC:                 arcCli.New2(c.ArchiveRPC),
		assRPC:                 assCli.New(c.AssistRPC),
		coinRPC:                coinCli.New(c.CoinRPC),
		relRPC:                 relCli.New(c.RelationRPC),
		locationRPC:            locCli.New(c.LocationRPC),
		cache:                  fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		realname:               make(map[int64]int64),
		arcTypes:               make(map[int16]int16),
		broadcastChan:          make(chan *broadcast, 1024),
		localCache:             make(map[string][]byte),
		seqDmArg:               &seqMdl.ArgBusiness{BusinessID: c.Seq.DM.BusinessID, Token: c.Seq.DM.Token},
		seqSubtitleArg:         &seqMdl.ArgBusiness{BusinessID: c.Seq.Subtitle.BusinessID, Token: c.Seq.Subtitle.Token},
		seqRPC:                 seqCli.New2(c.SeqRPC),
		assistLogChan:          make(chan *assMdl.ArgAssistLogAdd, 1024),
		dmOperationLogSvc:      infoc.New(c.Infoc2),
		opsLogCh:               make(chan *oplog.Infoc, 1024),
		moniOidMap:             make(map[int64]struct{}),
		memberRPC:              memberCli.New(c.MemberRPC),
		figureRPC:              figureCli.New(c.FigureRPC),
		spyRPC:                 spyCli.New(c.SpyRPC),
		garbageDanmu:           c.Switch.GarbageDanmu,
		broadcastLimit:         c.BroadcastLimit.Limit,
		broadcastlimitInterval: c.BroadcastLimit.Interval,
		localViewCache:         make(map[string]*model.ViewDm),
		aidSheild:              make(map[int64]struct{}),
		midsSheild:             make(map[int64]struct{}),
	}
	for idStr, cid := range s.conf.Realname.Threshold {
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
	accountRPC, err := accountCli.NewClient(c.AccountRPC)
	if err != nil {
		panic(err)
	}
	s.accountRPC = accountRPC
	ugcPayRPC, err := ugcPayCli.NewClient(c.UgcPayRPC)
	if err != nil {
		panic(err)
	}
	s.ugcPayRPC = ugcPayRPC
	filterRPC, err := filterCli.NewClient(s.conf.FilterRPC)
	if err != nil {
		panic(err)
	}
	s.filterRPC = filterRPC
	seasonRPC, err := seasonCli.NewClient(c.SeasonRPC)
	if err != nil {
		panic(fmt.Sprintf("seasonCli.NewClient.error(%v)", err))
	}
	s.seasonRPC = seasonRPC
	thumbupRPC, err := thumbupApi.NewClient(c.ThumbupRPC)
	if err != nil {
		panic(err)
	}
	s.thumbupRPC = thumbupRPC
	subtitleLans, err := s.dao.SubtitleLans(context.Background())
	if err != nil {
		panic(err)
	}
	s.subtitleLans = model.SubtitleLans(subtitleLans)
	go s.broadcastproc()
	go s.archiveTypeproc()
	go s.localcacheproc()
	go s.assistLogproc()
	go s.oplogproc()
	go s.monitorproc()
	go s.viewProc()
	go s.shieldProc()
	return
}

// Ping service ping.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

func (s *Service) archiveTypeproc() {
	var mu sync.Mutex
	for {
		rmap, err := s.dao.TypeMapping(context.TODO())
		if err != nil {
			log.Error("load archive types error(%v)", err)
		} else {
			mu.Lock()
			s.arcTypes = rmap
			mu.Unlock()
		}
		time.Sleep(5 * time.Minute)
	}
}

func (s *Service) broadcastproc() {
	broadcastFmt := `["%.2f,%d,%d,%d,%d,%d,%d,%s,%d","%s"]`
	for m := range s.broadcastChan {
		if m.Pool == model.PoolSpecial {
			continue
		}
		if m.State != model.StateNormal && m.State != model.StateMonitorAfter {
			continue
		}
		if err := s.dao.BroadcastLimit(context.TODO(), m.Oid, m.Type, s.broadcastLimit, s.broadcastlimitInterval); err != nil {
			if err != ecode.LimitExceed {
				log.Error("dao.BroadcastLimit(oid:%v) error(%v)", m.Oid, err)
			}
			continue
		}
		hash := model.Hash(m.Mid, uint32(m.Content.IP))
		msg := strings.Replace(m.Content.Msg, `\`, `\\`, -1)
		msg = strings.Replace(msg, `"`, `\"`, -1)
		info := fmt.Sprintf(broadcastFmt, float32(m.Progress)/1000.0, m.Content.Mode, m.Content.FontSize,
			m.Content.Color, m.Ctime, m.Rnd, m.Pool, hash, m.ID, msg)
		if err := s.dao.BroadcastInGoim(context.TODO(), m.Oid, m.Aid, []byte(info)); err != nil {
			log.Error("dao.BroadcastInGoim(cid:%d,aid:%d,info:%s) error(%v)", m.Oid, m.Aid, info, err)
		} else {
			log.Info("BroadcastInGoim(%s) succeed", info)
		}
		if err := s.dao.Broadcast(context.Background(), m.Oid, m.Aid, info); err != nil {
			log.Error("dao.Broadcast(cid:%d,aid:%d,info:%s) error(%v)", m.Oid, m.Aid, info, err)
		} else {
			log.Info("broadcast(%s) succeed", info)
		}
	}
}

func (s *Service) localcacheproc() {
	for {
		s.loadLocalcache(s.conf.Localcache.Oids)
		time.Sleep(time.Duration(s.conf.Localcache.Expire))
	}
}

func (s *Service) oplogproc() {
	for opLog := range s.opsLogCh {
		if len(opLog.Subject) == 0 || len(opLog.CurrentVal) == 0 || opLog.Source <= 0 ||
			opLog.Operator <= 0 || opLog.OperatorType <= 0 {
			log.Warn("oplogproc() it is an illegal log, warn(%v, %v, %v)", opLog.Subject, opLog.Subject, opLog.CurrentVal)
			continue
		} else {
			for _, dmid := range opLog.DMIds {
				if dmid > 0 {
					s.dmOperationLogSvc.Info(opLog.Subject, strconv.FormatInt(opLog.Oid, 10), strconv.Itoa(opLog.Type),
						strconv.FormatInt(dmid, 10), opLog.Source.String(), opLog.OriginVal,
						opLog.CurrentVal, strconv.FormatInt(opLog.Operator, 10), opLog.OperatorType.String(),
						opLog.OperationTime, opLog.Remark)
				} else {
					log.Warn("oplogproc() it is an illegal log, for dmid value, warn(%d, %+v)", dmid, opLog)
				}
			}
		}
	}
}

// OpLog put a new infoc format operation log into the channel
func (s *Service) OpLog(c context.Context, cid, operator, OperationTime int64, typ int, dmids []int64, subject, originVal, currentVal, remark string, source oplog.Source, operatorType oplog.OperatorType) (err error) {
	infoLog := new(oplog.Infoc)
	infoLog.Oid = cid
	infoLog.Type = typ
	infoLog.DMIds = dmids
	infoLog.Subject = subject
	infoLog.OriginVal = originVal
	infoLog.CurrentVal = currentVal
	infoLog.OperationTime = strconv.FormatInt(OperationTime, 10)
	infoLog.Source = source
	infoLog.OperatorType = operatorType
	infoLog.Operator = operator
	infoLog.Remark = remark
	select {
	case s.opsLogCh <- infoLog:
	default:
		err = fmt.Errorf("opsLogCh full")
		log.Error("opsLogCh full (%v)", infoLog)
	}
	return
}

func (s *Service) monitorproc() {
	for {
		time.Sleep(3 * time.Second)
		s.oidLock.Lock()
		oidMap := s.moniOidMap
		s.moniOidMap = make(map[int64]struct{})
		s.oidLock.Unlock()
		for oid := range oidMap {
			sub, err := s.dao.Subject(context.TODO(), model.SubTypeVideo, oid)
			if err != nil || sub == nil {
				continue
			}
			if err := s.updateMonitorCnt(context.TODO(), sub); err != nil {
				log.Error("s.updateMonitorCnt(%+v) error(%v)", sub, err)
			}
		}
	}
}
