package service

import (
	"context"
	"math/rand"
	"sync"
	"time"

	dm2rpc "go-common/app/interface/main/dm2/rpc/client"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	"go-common/app/interface/main/web/conf"
	"go-common/app/interface/main/web/dao"
	"go-common/app/interface/main/web/model"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	accclient "go-common/app/service/main/account/api"
	arcclient "go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	brdrpc "go-common/app/service/main/broadcast/api/grpc/v1"
	coinclient "go-common/app/service/main/coin/api"
	couprpc "go-common/app/service/main/coupon/rpc/client"
	dyrpc "go-common/app/service/main/dynamic/rpc/client"
	favrpc "go-common/app/service/main/favorite/api/gorpc"
	locrpc "go-common/app/service/main/location/rpc/client"
	relrpc "go-common/app/service/main/relation/rpc/client"
	resrpc "go-common/app/service/main/resource/rpc/client"
	shareclient "go-common/app/service/main/share/api"
	thumbrpc "go-common/app/service/main/thumbup/rpc/client"
	ugcclient "go-common/app/service/main/ugcpay/api/grpc/v1"
	"go-common/library/log"
	"go-common/library/log/anticheat"
	"go-common/library/sync/pipeline/fanout"
)

// Service service
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	// rpc
	arc      *arcrpc.Service2
	dy       *dyrpc.Service
	tag      *tagrpc.Service
	loc      *locrpc.Service
	art      *artrpc.Service
	res      *resrpc.Service
	relation *relrpc.Service
	thumbup  *thumbrpc.Service
	coupon   *couprpc.Service
	dm2      *dm2rpc.Service
	fav      *favrpc.Service
	// chans
	rids map[int32]struct{}
	// cache proc
	cache                                   *fanout.Fanout
	regionCount                             map[int16]int
	onlineArcs                              []*model.OnlineArc
	allArchivesCount, playOnline, webOnline int64
	// elec show typeids
	elecShowTypeIds map[int32]struct{}
	// no related data aids
	noRelAids map[int64]struct{}
	// rand source
	r         *rand.Rand
	indexIcon *model.IndexIcon
	// infoc2
	CheatInfoc *anticheat.AntiCheat
	// TypeNames
	typeNames map[int32]*arcclient.Tp
	// Broadcast grpc client
	broadcastClient brdrpc.ZergClient
	coinClient      coinclient.CoinClient
	arcClient       arcclient.ArchiveClient
	accClient       accclient.AccountClient
	shareClient     shareclient.ShareClient
	ugcPayClient    ugcclient.UGCPayClient
	// searchEggs
	searchEggs map[int64]*model.SearchEggRes
	// bnj
	bnj2019View    *model.Bnj2019View
	BnjElecInfo    *model.ElecShow
	bnj2019List    []*model.Bnj2019Related
	bnj2019LiveArc *arcclient.ArcReply
	bnjGrayUids    map[int64]struct{}
	specialMids    map[int64]struct{}
}

// New new
func New(c *conf.Config) *Service {
	s := &Service{
		c:           c,
		dao:         dao.New(c),
		arc:         arcrpc.New2(c.ArchiveRPC),
		dy:          dyrpc.New(c.DynamicRPC),
		tag:         tagrpc.New2(c.TagRPC),
		loc:         locrpc.New(c.LocationRPC),
		art:         artrpc.New(c.ArticleRPC),
		res:         resrpc.New(c.ResourceRPC),
		relation:    relrpc.New(c.RelationRPC),
		thumbup:     thumbrpc.New(c.ThumbupRPC),
		coupon:      couprpc.New(c.CouponRPC),
		dm2:         dm2rpc.New(c.Dm2RPC),
		fav:         favrpc.New2(c.FavRPC),
		cache:       fanout.New("cache"),
		regionCount: make(map[int16]int),
		r:           rand.New(rand.NewSource(time.Now().Unix())),
		specialMids: map[int64]struct{}{},
	}
	var err error
	if s.broadcastClient, err = brdrpc.NewClient(c.BroadcastClient); err != nil {
		panic(err)
	}
	if s.coinClient, err = coinclient.NewClient(c.CoinClient); err != nil {
		panic(err)
	}
	if s.arcClient, err = arcclient.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	if s.accClient, err = accclient.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	if s.shareClient, err = shareclient.NewClient(c.ShareClient); err != nil {
		panic(err)
	}
	if s.ugcPayClient, err = ugcclient.NewClient(c.UGCClient); err != nil {
		panic(err)
	}
	s.initRules()
	go s.newCountproc()
	go s.onlineCountproc()
	go s.onlineListproc()
	go s.indexIconproc()
	go s.typeNameproc()
	go s.searchEggproc()
	go s.bnj2019proc()
	go s.loadManager()
	// init infoc
	if c.Infoc2 != nil {
		s.CheatInfoc = anticheat.New(c.Infoc2)
	}
	return s
}

func (s *Service) initRules() {
	tmpRids := make(map[int32]struct{}, len(s.c.Rule.Rids))
	for _, v := range s.c.Rule.Rids {
		tmpRids[v] = struct{}{}
	}
	s.rids = tmpRids
	tmpElec := make(map[int32]struct{}, len(s.c.Rule.ElecShowTypeIDs))
	for _, id := range s.c.Rule.ElecShowTypeIDs {
		tmpElec[id] = struct{}{}
	}
	s.elecShowTypeIds = tmpElec
	tmpNoRel := make(map[int64]struct{}, len(s.c.Rule.NoRelAids))
	for _, id := range s.c.Rule.NoRelAids {
		tmpNoRel[id] = struct{}{}
	}
	s.noRelAids = tmpNoRel
}

func (s *Service) typeNameproc() {
	for {
		if typesReply, err := s.arcClient.Types(context.Background(), &arcclient.NoArgRequest{}); err != nil {
			log.Error("s.arc.Types2 error(%v)", err)
			time.Sleep(time.Second)
			continue
		} else {
			s.typeNames = typesReply.Types
		}
		time.Sleep(time.Duration(s.c.WEB.PullOnlineInterval))
	}
}

func (s *Service) searchEggproc() {
	var mutex = sync.Mutex{}
	for {
		if eggs, err := s.dao.SearchEgg(context.Background()); err != nil {
			log.Error("s.dao.SearchEgg error(%v)", err)
			time.Sleep(5 * time.Second)
			continue
		} else {
			data := make(map[int64]*model.SearchEggRes, len(eggs))
			for _, v := range eggs {
				if source, ok := v.Plat[_searchEggWebPlat]; ok {
					for _, egg := range source {
						if _, isSet := data[egg.EggID]; !isSet {
							data[egg.EggID] = &model.SearchEggRes{
								EggID:     egg.EggID,
								ShowCount: v.ShowCount,
							}
						}
						source := &model.SearchEggSource{URL: egg.URL, MD5: egg.MD5, Size: egg.Size}
						data[egg.EggID].Source = append(data[egg.EggID].Source, source)
					}
				}
			}
			mutex.Lock()
			s.searchEggs = data
			mutex.Unlock()
		}
		time.Sleep(time.Duration(s.c.WEB.SearchEggInterval))
	}
}

// Ping check connection success.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}

// Close close resource.
func (s *Service) Close() {
	s.dao.Close()
}

func archivesArgLog(name string, aids []int64) {
	if aidLen := len(aids); aidLen >= 50 {
		log.Info("s.arc.Archives3 func(%s) len(%d), arg(%v)", name, aidLen, aids)
	}
}
