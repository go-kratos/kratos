package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-common/app/interface/main/tag/conf"
	"go-common/app/interface/main/tag/dao"
	"go-common/app/interface/main/tag/model"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	figrpc "go-common/app/service/main/figure/rpc/client"
	rpcModel "go-common/app/service/main/tag/model"
	tagrpc "go-common/app/service/main/tag/rpc/client"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"

	"golang.org/x/sync/singleflight"
)

const (
	_hotTagKey  = "%d_%d"
	_similarKey = "%d_%d"

	_videoTagUpURL  = "/videoup/tag/up"
	_filterURL      = "/x/internal/filter"
	_mFilterURL     = "/x/internal/filter/multi"
	_aiRecommandURL = "/recommand"
)

func (s *Service) hotTagKey(cid, hotType int64) string {
	return fmt.Sprintf(_hotTagKey, cid, hotType)
}

func (s *Service) similarKey(rid, tid int64) string {
	return fmt.Sprintf(_similarKey, rid, tid)
}

// Service service.
type Service struct {
	c            *conf.Config
	dao          *dao.Dao
	arcRPC       *arcrpc.Service2
	tagRPC       *tagrpc.Service
	figRPC       *figrpc.Service
	cacheCh      *fanout.Fanout
	invalidArcCh *fanout.Fanout
	rids         []int64
	ridMap       map[int64][]int64 // rid map
	pridMap      map[int64]int64   // prid map
	hotTag       map[string]*model.HotTags
	simlars      map[string][]*model.SimilarTag

	limitArc    map[int64]int8
	whiteUser   map[int64]struct{}
	tagListUser map[int64]struct{}
	sLock       sync.RWMutex
	hLock       sync.RWMutex
	channelLock sync.RWMutex
	limitLock   sync.RWMutex

	//
	videoTagUpURL   string
	client          *bm.Client
	similarURL      string
	siClient        *bm.Client // http similar client
	filterURL       string
	mFilterURL      string
	aiRecommandlURL string

	// channel.
	channelMap        map[int64]*model.Channel
	channelCategories []*model.ChannelCategory
	channelTypeMap    map[int64][]*model.Channel
	channelRecommand  []*model.Channel
	channelRule       map[int64]*model.ChannelRuleClassifier

	singleGroup singleflight.Group
}

// New new a service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:            c,
		dao:          dao.New(c),
		arcRPC:       arcrpc.New2(c.ArchiveRPC),
		tagRPC:       tagrpc.New(c.TagDisRPC),
		figRPC:       figrpc.New(c.FigureRPC),
		cacheCh:      fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		invalidArcCh: fanout.New("invalid_arc", fanout.Worker(1), fanout.Buffer(1024)),
		ridMap:       make(map[int64][]int64),
		pridMap:      make(map[int64]int64),
		hotTag:       make(map[string]*model.HotTags),
		simlars:      make(map[string][]*model.SimilarTag),
		limitArc:     make(map[int64]int8),
		whiteUser:    make(map[int64]struct{}),
		tagListUser:  make(map[int64]struct{}),

		videoTagUpURL:   c.Host.Archive + _videoTagUpURL,
		client:          bm.NewClient(c.HTTPClient),
		siClient:        bm.NewClient(c.HTTPSimilar),
		filterURL:       c.Host.APICo + _filterURL,
		mFilterURL:      c.Host.APICo + _mFilterURL,
		aiRecommandlURL: c.Host.AI + _aiRecommandURL,
		similarURL:      c.Host.BigDataURL,

		channelMap:        make(map[int64]*model.Channel),
		channelCategories: make([]*model.ChannelCategory, 0),
		channelTypeMap:    make(map[int64][]*model.Channel),
		channelRecommand:  make([]*model.Channel, 0),
		channelRule:       make(map[int64]*model.ChannelRuleClassifier),
	}
	time.Sleep(time.Second * 3)
	if err := s.channelCaches(); err != nil {
		panic(err)
	}
	if len(s.c.Tag.White) > 0 {
		for _, v := range s.c.Tag.White {
			s.tagListUser[v] = struct{}{}
		}
	}
	go s.crontabProc()
	go s.loadLimitArc()
	go s.channelproc()
	time.Sleep(time.Second * 3)
	return
}

func (s *Service) crontabProc() {
	for {
		var err error
		s.rids, s.pridMap, s.ridMap, err = s.ridsService(context.Background())
		if err != nil {
			log.Error("s.ridsService error(%v)", err)
			time.Sleep(time.Second)
			continue
		}
		s.initHotTag()
		s.simlars = make(map[string][]*model.SimilarTag)
		time.Sleep(time.Minute * 5)
	}
}

func (s *Service) initHotTag() {
	var (
		err      error
		hots     []*model.HotTag
		hotTypes = []int64{0, 1}
		hotTmp   = make(map[string]*model.HotTags)
	)
	for _, rid := range s.rids {
		for _, hotType := range hotTypes {
			if hots, err = s.hots(context.Background(), rid, hotType); err != nil {
				continue
			}
			hotTmp[s.hotTagKey(rid, hotType)] = &model.HotTags{Rid: rid, Tags: hots}
		}
	}
	s.hLock.Lock()
	s.hotTag = hotTmp
	s.hLock.Unlock()
}

func (s *Service) loadLimitArc() {
	for {
		var (
			lrs      []*rpcModel.ResourceLimit
			limitArc = make(map[int64]int8)
			midm     = make(map[int64]struct{})
			ctx      = context.Background()
		)
		lrs, _ = s.limitResource(ctx)
		for _, lr := range lrs {
			limitArc[lr.Oid] = int8(lr.Attr)
		}
		midm, _ = s.whiteUserService(ctx)
		s.limitLock.Lock()
		s.limitArc = limitArc
		s.whiteUser = midm
		s.limitLock.Unlock()
		time.Sleep(time.Minute * 10)
	}
}

// Ping Ping.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.PingRe(c); err != nil {
		log.Error("s.dao.PingRe err(%v)", err)
	}
	return
}

// Close .
func (s *Service) Close() {
	s.dao.Close()
	s.cacheCh.Close()
	s.invalidArcCh.Close()
}
