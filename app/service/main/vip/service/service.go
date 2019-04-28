package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	v1 "go-common/app/service/main/coupon/api"
	couponrpc "go-common/app/service/main/coupon/rpc/client"
	memrpc "go-common/app/service/main/member/api/gorpc"
	memmdl "go-common/app/service/main/member/model"
	pointrpc "go-common/app/service/main/point/rpc/client"
	"go-common/app/service/main/vip/conf"
	"go-common/app/service/main/vip/dao"
	"go-common/app/service/main/vip/model"
	"go-common/app/service/main/vip/service/price"
	"go-common/library/log"
)

var (
	_retryTimes               = 3
	_iosPaychannel            = 100
	_ymonth             int16 = 12
	_pointRemark              = "积分兑换"
	_maxSizeUsers             = 50
	_defOrderExpire           = 7200
	_autoPayserviceType       = 7
	_daysecond                = 86400
	_remindday                = 8
	_remindtxt                = "%d天之后就要收回特权了哦，快去续费吧！"
	_defpn                    = 1
	_defps                    = 20
	_tiplimit                 = 3
	_annualMonth        int32 = 12
	_millis             int64 = 1000
	_willexpiredays     int64 = 7
	_yyyymmdd                 = "2006-01-02"
	_defround                 = 2
)

//Service vip service
type Service struct {
	dao            *dao.Dao
	c              *conf.Config
	vipConfig      map[string]*model.VipConfig
	missch         chan func()
	missBcoin      chan func()
	months         map[string]*model.Month
	pointpricemap  map[int16]*model.PointExchangePrice
	pointprices    []*model.PointExchangePrice
	pointRPC       *pointrpc.Service
	memRPC         *memrpc.Service
	couRPC         *couponrpc.Service
	r              *rand.Rand
	tips           map[string][]*model.Tips
	vipPriceConf   *price.Price
	vipConfSuitMax map[int64]int8
	vipPriceMap    map[int64]*model.VipPriceConfig
	//map[int8][]*model.Privilege
	vipPrivilege             sync.Map
	vipPrivilegeResourcesMap map[int64]map[int8]*model.PrivilegeResources
	jointlyList              []*model.Jointly
	associateVipMap          map[int8][]*model.AssociateVipResp
	platformConf             map[string]int64
	pLock                    sync.RWMutex
	//map[int64]map[int64]*model.ConfDialog
	dialogMap sync.Map
	// coupon grpc service
	coupongRPC v1.CouponClient
}

//New new service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:             c,
		dao:           dao.New(c),
		missch:        make(chan func(), 10240),
		missBcoin:     make(chan func(), 10240),
		months:        map[string]*model.Month{},
		pointpricemap: map[int16]*model.PointExchangePrice{},
		pointprices:   []*model.PointExchangePrice{},
		pointRPC:      pointrpc.New(c.RPCClient2.Point),
		memRPC:        memrpc.New(c.RPCClient2.Member),
		couRPC:        couponrpc.New(c.RPCClient2.Coupon),
		r:             rand.New(rand.NewSource(time.Now().Unix())),
		tips:          map[string][]*model.Tips{},
		vipPrivilegeResourcesMap: map[int64]map[int8]*model.PrivilegeResources{},
		jointlyList:              []*model.Jointly{},
		associateVipMap:          map[int8][]*model.AssociateVipResp{},
		platformConf:             map[string]int64{},
		vipPriceConf:             price.New(),
	}
	if err := s.loadtips(); err != nil {
		panic(err)
	}
	if err := s.loadjointly(); err != nil {
		panic(err)
	}
	s.loadVipPriceConfig()
	if err := s.loadPrivilege(); err != nil {
		panic(err)
	}
	if err := s.loadAssociateVip(); err != nil {
		panic(err)
	}
	if err := s.loadplatformconf(); err != nil {
		panic(err)
	}
	coupongRPC, err := v1.NewClient(c.CouponClient)
	if err != nil {
		panic(err)
	}
	s.coupongRPC = coupongRPC
	go s.loadassociatevipproc()
	go s.loadvippriceconfigproc()
	go s.loadvipconfigproc()
	go s.bcoinproc()
	go s.cacheproc()
	go s.loadmonthproc()
	go s.loadpointpriceproc()
	go s.loadtipsproc()
	go s.loadprivilegeproc()
	go s.loadjointlysproc()
	go s.loadplatformconfproc()

	s.loadDialog()
	go s.loaddialogproc()
	return
}

func (s *Service) loadvippriceconfigproc() {
	for {
		time.Sleep(1 * time.Second)
		s.loadVipPriceConfig()
	}
}

func (s *Service) loadvipconfigproc() {
	defer func() {
		//if x :=
	}()
	for {
		s.loadVipConfig(context.TODO())
		time.Sleep(time.Second * 1)
	}
}
func (s *Service) loadVipConfig(c context.Context) {
	var (
		res []*model.VipConfig
		err error
	)
	if res, err = s.dao.SelAllConfig(c); err != nil {
		log.Error("s.dao.SelAllConfig err(%v)", err)
		return
	}
	vcm := make(map[string]*model.VipConfig)
	for _, v := range res {
		vcm[v.ConfigKey] = v
	}
	s.vipConfig = vcm

}

//Ping check db live
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

//clean cache
func (s *Service) cache(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}

func (s *Service) asyncBcoin(f func()) {
	select {
	case s.missBcoin <- f:
	default:
		log.Warn("bcoinproc chan full")
	}
}
func (s *Service) bcoinproc() {
	for {
		f := <-s.missBcoin
		f()
	}
}

//Close close DB
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) loadmonth() (err error) {
	var (
		res []*model.Month
	)
	if res, err = s.dao.AllMonthByOrder(context.TODO(), "ASC"); err != nil {
		log.Error("all months by order error(%+v)", err)
		return
	}
	tmp := make(map[string]*model.Month, len(res))
	for _, r := range res {
		tmp[s.monthkey(r.Month, r.MonthType)] = r
	}
	s.months = tmp
	log.Info("loadmonth (%v) load success", tmp)
	return

}

func (s *Service) monthkey(month int16, monthType int8) string {
	return strconv.Itoa(int(month)) + "_" + strconv.Itoa(int(monthType))
}

func (s *Service) loadmonthproc() {
	for {
		s.loadmonth()
		time.Sleep(10 * time.Second)
	}
}

func (s *Service) loadpointpriceproc() {
	for {
		s.loadpointprice()
		time.Sleep(2 * time.Minute)
	}
}

func (s *Service) loadpointprice() {
	var (
		res []*model.PointExchangePrice
		err error
	)
	if res, err = s.dao.AllPointExchangePrice(context.TODO()); err != nil {
		log.Error("load all point price error %+v", err)
		return
	}
	tmp := make(map[int16]*model.PointExchangePrice, len(res))
	for _, v := range res {
		tmp[v.Month] = v
	}
	s.pointpricemap = tmp
	s.pointprices = res
	log.Info("loadpointprice (%v) load success", tmp)
}

func (s *Service) loadtipsproc() {
	for {
		time.Sleep(2 * time.Minute)
		s.loadtips()
	}
}

func (s *Service) loadprivilegeproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.loadprivilegeproc panic(%v)", x)
			go s.loadprivilegeproc()
			log.Info("service.loadprivilegeproc recover")
		}
	}()
	for {
		time.Sleep(2 * time.Minute)
		if err := s.loadPrivilege(); err != nil {
			log.Error("service.loadprivilegeproc err(%+v)", err)
		}
	}
}

func (s *Service) loadtips() (err error) {
	var (
		ts  []*model.Tips
		now = time.Now().Unix()
		key string
	)
	if ts, err = s.dao.AllTips(context.TODO(), now); err != nil {
		log.Error("load all tips error %+v", err)
		return
	}
	tmp := map[string][]*model.Tips{}
	for _, v := range ts {
		key = s.tipsKey(v.Platform, v.Position)
		tmp[key] = append(tmp[key], v)
	}
	s.tips = tmp
	log.Info("load all tips success(%v)", tmp)
	return
}

func (s *Service) tipsKey(platform int64, position int8) string {
	return fmt.Sprintf("%d_%d", platform, position)
}

func (s *Service) loadjointlysproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.loadjointlysproc panic(%v)", x)
			go s.loadjointlysproc()
			log.Info("service.loadjointlysproc recover")
		}
	}()
	for {
		time.Sleep(2 * time.Minute)
		if err := s.loadjointly(); err != nil {
			log.Error("service.loadjointlysproc err(%+v)", err)
		}
	}
}

func (s *Service) loadplatformconfproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.loadplatformconfproc panic(%v)", x)
			go s.loadplatformconfproc()
			log.Info("service.loadplatformconfproc recover")
		}
	}()
	for {
		time.Sleep(30 * time.Second)
		if err := s.loadplatformconf(); err != nil {
			log.Error("service.loadplatformconf err(%+v)", err)
		}
	}
}

func (s *Service) loadassociatevipproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.loadassociatevipproc panic(%v)", x)
			go s.loadassociatevipproc()
			log.Info("service.loadassociatevipproc recover")
		}
	}()
	for {
		time.Sleep(30 * time.Second)
		if err := s.loadAssociateVip(); err != nil {
			log.Error("service.loadassociatevipproc err(%+v)", err)
		}
	}
}

func (s *Service) retryGetMemberInfo(c context.Context, mid int64) (m *memmdl.BaseInfo, err error) {
	for i := 0; i < _retryTimes; i++ {
		if m, err = s.memRPC.Base(c, &memmdl.ArgMemberMid{Mid: mid}); err != nil {
			continue
		}
		break
	}
	return
}
