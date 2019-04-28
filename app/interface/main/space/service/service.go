package service

import (
	"context"
	"time"

	"go-common/app/interface/main/space/conf"
	"go-common/app/interface/main/space/dao"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	accclient "go-common/app/service/main/account/api"
	accwar "go-common/app/service/main/account/api"
	accmdl "go-common/app/service/main/account/model"
	arcclient "go-common/app/service/main/archive/api"
	assrpc "go-common/app/service/main/assist/rpc/client"
	coinclient "go-common/app/service/main/coin/api"
	favrpc "go-common/app/service/main/favorite/api/gorpc"
	fltrpc "go-common/app/service/main/filter/rpc/client"
	member "go-common/app/service/main/member/api/gorpc"
	"go-common/app/service/main/relation/rpc/client"
	thumbup "go-common/app/service/main/thumbup/rpc/client"
	upclient "go-common/app/service/main/up/api/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"
)

// Service service struct.
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	// rpc
	art      *artrpc.Service
	ass      *assrpc.Service
	tag      *tagrpc.Service
	filter   *fltrpc.Service
	fav      *favrpc.Service
	thumbup  *thumbup.Service
	relation *relation.Service
	member   *member.Service
	// grpc
	accClient  accwar.AccountClient
	arcClient  arcclient.ArchiveClient
	coinClient coinclient.CoinClient
	upClient   upclient.UpClient
	// cache proc
	cache *fanout.Fanout
	// noNoticeMids
	noNoticeMids   map[int64]struct{}
	BlacklistValue map[int64]struct{}
}

// New new service.
func New(c *conf.Config) *Service {
	s := &Service{
		c:        c,
		dao:      dao.New(c),
		art:      artrpc.New(c.ArticleRPC),
		ass:      assrpc.New(c.AssistRPC),
		tag:      tagrpc.New2(c.TagRPC),
		fav:      favrpc.New2(c.FavoriteRPC),
		filter:   fltrpc.New(c.FilterRPC),
		thumbup:  thumbup.New(c.ThumbupRPC),
		relation: relation.New(c.RelationRPC),
		member:   member.New(c.MemberRPC),
		cache:    fanout.New("cache"),
	}
	var err error
	if s.accClient, err = accclient.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	if s.arcClient, err = arcclient.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	if s.coinClient, err = coinclient.NewClient(c.CoinClient); err != nil {
		panic(err)
	}
	if s.upClient, err = upclient.NewClient(c.UpClient); err != nil {
		panic(err)
	}
	s.initMids()
	go s.loadBlacklist()
	return s
}

// Ping ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.Ping(c); err != nil {
		log.Error("s.dao.Ping error(%v)", err)
	}
	return
}

func (s *Service) initMids() {
	tmp := make(map[int64]struct{}, len(s.c.Rule.NoNoticeMids))
	for _, id := range s.c.Rule.NoNoticeMids {
		tmp[id] = struct{}{}
	}
	s.noNoticeMids = tmp
}

func (s *Service) realName(c context.Context, mid int64) (profile *accmdl.Profile, err error) {
	var reply *accwar.ProfileReply
	if reply, err = s.accClient.Profile3(c, &accwar.MidReq{Mid: mid}); err != nil || reply == nil {
		log.Error("s.accClient.Profile3(%d) error(%v)", mid, err)
		return
	}
	profile = reply.Profile
	if !s.c.Rule.RealNameOn {
		return
	}
	if profile.Identification == 0 && profile.TelStatus == 0 {
		err = ecode.UserCheckNoPhone
		return
	}
	if profile.Identification == 0 && profile.TelStatus == 2 {
		err = ecode.UserCheckInvalidPhone
		return
	}
	return
}

func (s *Service) privacyCheck(c context.Context, vmid int64, field string) (err error) {
	privacy := s.privacy(c, vmid)
	if value, ok := privacy[field]; !ok || value != _defaultPrivacy {
		err = ecode.SpaceNoPrivacy
		return
	}
	return
}

// loadBlacklist load spack blacklist
func (s *Service) loadBlacklist() {
	for {
		time.Sleep(time.Duration(conf.Conf.Rule.BlackFre))
		s.Blacklist(context.Background())
	}
}
