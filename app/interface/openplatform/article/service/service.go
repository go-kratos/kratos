package service

import (
	"context"
	"math/rand"
	"time"

	hisrpc "go-common/app/interface/main/history/rpc/client"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	"go-common/app/interface/openplatform/article/conf"
	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"
	search "go-common/app/interface/openplatform/article/model/search"
	account "go-common/app/service/main/account/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	arcmdl "go-common/app/service/main/archive/model/archive"
	coinrpc "go-common/app/service/main/coin/api/gorpc"
	favrpc "go-common/app/service/main/favorite/api/gorpc"
	filterrpc "go-common/app/service/main/filter/rpc/client"
	resrpc "go-common/app/service/main/resource/rpc/client"
	thumbuprpc "go-common/app/service/main/thumbup/rpc/client"
	xcache "go-common/library/cache"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/log/anticheat"
	"go-common/library/net/rpc/warden"
)

const _segmentAddr = "127.0.01:6755"

// cache proc
var cache *xcache.Cache

// Service service
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	// rpc
	accountRPC           *accrpc.Service3
	tagRPC               *tagrpc.Service
	favRPC               *favrpc.Service
	thumbupRPC           thumbuprpc.ThumbupRPC
	arcRPC               *arcrpc.Service2
	coinRPC              coinrpc.RPC
	resRPC               *resrpc.Service
	filterRPC            *filterrpc.Service
	hisRPC               *hisrpc.Service
	searchRPC            search.TagboxServiceClient
	categoriesMap        map[int64]*artmdl.Category
	categoriesReverseMap map[int64][]*artmdl.Category
	categoryParents      map[int64][]*artmdl.Category
	primaryCategories    []*artmdl.Category
	Categories           artmdl.Categories
	RecommendsMap        map[int64][][]*artmdl.Recommend
	RecommendsGroups     map[int64]map[int64]bool
	recommendChan        chan [2]int64
	recommendAids        map[int64][]int64
	setting              *artmdl.Setting
	activities           map[int64]*artmdl.Activity
	// infoc
	logCh chan interface{}
	//banner
	bannersMap    map[int8][]*artmdl.Banner
	actBannersMap map[int8][]*artmdl.Banner
	// rank
	ranksMap      map[int64]bool
	sortLimitTime int64
	notices       []*artmdl.Notice
	CheatInfoc    *anticheat.AntiCheat
}

// New new
func New(c *conf.Config) *Service {
	rand.Seed(time.Now().Unix())
	s := &Service{
		c:                    c,
		dao:                  dao.New(c),
		accountRPC:           accrpc.New3(c.AccountRPC),
		tagRPC:               tagrpc.New2(c.TagRPC),
		favRPC:               favrpc.New2(c.FavRPC),
		arcRPC:               arcrpc.New2(c.ArcRPC),
		coinRPC:              coinrpc.New(c.CoinRPC),
		resRPC:               resrpc.New(c.ResRPC),
		thumbupRPC:           thumbuprpc.New(c.ThumbupRPC),
		filterRPC:            filterrpc.New(c.FilterRPC),
		hisRPC:               hisrpc.New(c.HistoryRPC),
		searchRPC:            searchRPC(c.SearchRPC),
		categoriesMap:        make(map[int64]*artmdl.Category),
		categoriesReverseMap: make(map[int64][]*artmdl.Category),
		categoryParents:      make(map[int64][]*artmdl.Category),
		logCh:                make(chan interface{}, 1024),
		recommendChan:        make(chan [2]int64, 10240),
		recommendAids:        make(map[int64][]int64),
		ranksMap:             make(map[int64]bool),
		sortLimitTime:        int64(time.Duration(c.Article.SortLimitTime) / time.Second),
		CheatInfoc:           anticheat.New(c.CheatInfoc),
		RecommendsGroups:     make(map[int64]map[int64]bool),
	}
	s.loadCategories()
	s.loadSettings()
	s.loadRanks()
	go s.loadCategoriesproc()
	go s.loadSettingsproc()
	go s.loadNoticeproc()
	go s.infocproc()
	go s.loadRecommendsproc()
	go s.loadBannersproc()
	go s.loadActBannersproc()
	go s.deleteRecommendproc()
	go s.loadActivityproc()
	return s
}

func (s *Service) loadRecommendsproc() {
	for {
		now := time.Now().Unix()
		c := context.TODO()
		if (s.RecommendsMap == nil) || (now%s.dao.UpdateRecommendsInterval == 0) {
			err := s.UpdateRecommends(c)
			if err != nil {
				dao.PromError("service:更新推荐数据")
				time.Sleep(time.Second)
				continue
			}
			if err = s.groupRecommend(c); err != nil {
				log.Error("s.groupRecommend error(%+v)", err)
			}
		}
		// 这里不是每秒钟一更新
		time.Sleep(time.Second)
	}
}

// Close close dao.
func (s *Service) Close() {
	s.dao.Close()
}

// Ping check connection success.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}

// UserDisabled check user is disabled
func (s *Service) UserDisabled(c context.Context, mid int64) (res bool, level int, err error) {
	arg := account.ArgMid{Mid: mid}
	card, err := s.accountRPC.Card3(c, &arg)
	if (err == ecode.UserNotExist) || (err == ecode.MemberNotExist) {
		return false, 0, nil
	}
	if err != nil {
		dao.PromError("service:用户封禁状态")
		log.Error("s.accountRPC.Card2(%+v) err: %+v", arg, err)
		return
	}
	if card.Silence == 1 {
		res = true
	}
	level = int(card.Level)
	return
}

func (s *Service) isUpper(c context.Context, mid int64) (res bool, err error) {
	arg := &arcmdl.ArgUpCount2{Mid: mid}
	var count int
	if count, err = s.arcRPC.UpCount2(c, arg); err != nil {
		dao.PromError("service:up主投稿")
		log.Error("s.arcRPC.UpCount2(%v) err: %+v", mid, err)
		return
	}
	if count > 0 {
		res = true
	}
	return
}

func (s *Service) loadActivityproc() {
	for {
		if acts, err := s.dao.Activity(context.TODO()); err == nil {
			s.activities = acts
		}
		time.Sleep(time.Minute)
	}
}

func init() {
	cache = xcache.New(1, 1024)
}

func searchRPC(cfg *warden.ClientConfig) search.TagboxServiceClient {
	cc, err := warden.NewClient(cfg).Dial(context.Background(), "discovery://default/search.tagbox")
	if err != nil {
		panic(err)
	}
	return search.NewTagboxServiceClient(cc)
}
