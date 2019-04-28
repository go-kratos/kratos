package ad

import (
	"context"
	"sort"
	"strings"
	"time"

	"go-common/app/interface/main/creative/conf"
	gD "go-common/app/interface/main/creative/dao/game"
	porderD "go-common/app/interface/main/creative/dao/porder"
	gM "go-common/app/interface/main/creative/model/game"
	porderM "go-common/app/interface/main/creative/model/porder"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
)

//Service struct.
type Service struct {
	c            *conf.Config
	game         *gD.Dao
	porder       *porderD.Dao
	industryList []*porderM.Config
	showList     []*porderM.Config
}

//New get service.
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:            c,
		game:         gD.New(c),
		porder:       porderD.New(c),
		industryList: []*porderM.Config{},
		showList:     []*porderM.Config{},
	}
	s.loadPorderConfigList()
	go s.loadproc()
	return s
}

// GameList fn
func (s *Service) GameList(c context.Context, keywordStr, letterStr string, pn, ps int, ip string) (glist *gM.ListWithPager, err error) {
	list, err := s.game.List(context.TODO(), keywordStr, ip)
	if err != nil || len(list) == 0 {
		return
	}
	res := []*gM.ListItem{}
	if len(letterStr) > 0 {
		letterStr = strings.ToUpper(string(letterStr[0]))
		for _, ele := range list {
			if ele.Letter == letterStr {
				res = append(res, ele)
			}
		}
	} else {
		res = list
	}
	sort.Slice(res, func(i, j int) bool { return res[i].GameBaseID > res[j].GameBaseID })
	total := len(res)
	start := (pn - 1) * ps
	end := pn * ps
	glist = &gM.ListWithPager{
		Pn:    pn,
		Ps:    ps,
		Total: total,
	}
	if total <= start {
		glist.List = make([]*gM.ListItem, 0)
	} else if total <= end {
		glist.List = res[start:total]
	} else {
		glist.List = res[start:end]
	}
	return
}

// IndustryList fn
func (s *Service) IndustryList(c context.Context) (list []*porderM.Config, err error) {
	list = s.industryList
	return
}

// ShowList fn
func (s *Service) ShowList(c context.Context) (list []*porderM.Config, err error) {
	list = s.showList
	return
}

func (s *Service) loadproc() {
	//for wait rpc client connect
	time.Sleep(time.Duration(2 * time.Second))
	for {
		s.loadPorderConfigList()
		time.Sleep(time.Duration(2 * time.Minute))
	}
}

func (s *Service) loadPorderConfigList() {
	list, err := s.porder.ListConfig(context.TODO())
	if err != nil || list == nil || len(list) == 0 {
		return
	}
	s.industryList = []*porderM.Config{}
	s.showList = []*porderM.Config{}
	for _, v := range list {
		if v.Tp == porderM.ConfigTypeIndustry {
			s.industryList = append(s.industryList, v)
		}
		if v.Tp == porderM.ConfigTypeShow {
			s.showList = append(s.showList, v)
		}
	}
	log.Info("ListConfig (%d)|(%d)", len(s.showList), len(s.industryList))
}
