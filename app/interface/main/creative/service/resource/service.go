package resource

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/resource"
	"go-common/app/interface/main/creative/dao/tool"
	model "go-common/app/interface/main/creative/model/resource"
	"go-common/app/interface/main/creative/service"
	resmdl "go-common/app/service/main/resource/model"
	"go-common/library/log"
	"sort"
	"strconv"
	"time"
)

// Android iPhone
const (
	BanneriPhone    = 2417
	BannerAndroid   = 2431
	AcademyiPhone   = 2873
	AcademyAndroid  = 2877
	BannerCooperate = 2893
)

//Service struct
type Service struct {
	c      *conf.Config
	resDao *resource.Dao
	Seed   int64
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:      c,
		resDao: resource.New(c),
	}
	return s
}

// TopBanner fn
func (s *Service) TopBanner(c context.Context, mobiApp, device, network, ipAddr, buvid, adExtra string, build, resID int, plat int8, mid int64, isAd bool) (res []*model.Banner, err error) {
	var bnsm map[int][]*resmdl.Banner
	if resID == 0 {
		if model.IsAndroid(plat) {
			resID = BannerAndroid
		} else if model.IsIPhone(plat) || model.IsIPad(plat) {
			resID = BanneriPhone
			mobiApp = "iphone"
			device = "phone"
			plat = resmdl.PlatIPhone
		}
	}
	if bnsm, err = s.resDao.Banner(c, mobiApp, device, network, "", ipAddr, buvid, adExtra, strconv.Itoa(resID), build, plat, mid, isAd); err != nil {
		log.Error("s.resDao.Banner err(%v)", err)
		return
	}
	for _, rb := range bnsm[resID] {
		b := &model.Banner{}
		b.ChangeBanner(rb)
		if b.ClientIp == "" {
			b.ClientIp = ipAddr
		}
		res = append(res, b)
	}
	topLen := 5
	if len(res) > topLen {
		res = res[:5]
	}
	sort.Slice(res, func(i, j int) bool { return res[i].Rank < res[j].Rank })
	return
}

// AcademyBanner fn
func (s *Service) AcademyBanner(c context.Context, mobiApp, device, network, ipAddr, buvid, adExtra string, build, resID int, plat int8, mid int64, isAd bool) (randomRes []*model.Banner, err error) {
	var (
		bnsm map[int][]*resmdl.Banner
		res  = make([]*model.Banner, 0)
		keys []int
	)
	randomRes = make([]*model.Banner, 0)
	if model.IsAndroid(plat) {
		resID = AcademyAndroid
	} else if model.IsIPhone(plat) {
		resID = AcademyiPhone
	} else if model.IsIPad(plat) {
		return
	}
	if bnsm, err = s.resDao.Banner(c, mobiApp, device, network, "", ipAddr, buvid, adExtra, strconv.Itoa(resID), build, plat, mid, isAd); err != nil {
		log.Error("s.resDao.Banner err(%v)", err)
		return
	}
	for _, rb := range bnsm[resID] {
		b := &model.Banner{}
		b.ChangeBanner(rb)
		if b.ClientIp == "" {
			b.ClientIp = ipAddr
		}
		res = append(res, b)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].Rank < res[j].Rank })
	randLength := 3
	if len(res) > randLength {
		keys = tool.RandomSliceKeys(0, len(res), randLength, time.Now().Unix())
	} else {
		keys = tool.RandomSliceKeys(0, len(res), len(res), time.Now().Unix())
	}
	for _, k := range keys {
		randomRes = append(randomRes, res[k])
	}
	return
}

// CooperateBanner fn
func (s *Service) CooperateBanner(c context.Context, mobiApp, device, network, buvid, adExtra string, build int, plat int8, mid int64, isAd bool) (ass []*resmdl.Assignment, err error) {
	var res *resmdl.Resource
	if res, err = s.resDao.SimpleResource(c, BannerCooperate); err != nil {
		log.Error("Resource SimpleResource (%d) error(%v)", BannerCooperate, err)
		return
	}
	if res != nil {
		ass = res.Assignments
		return
	}
	return
}
