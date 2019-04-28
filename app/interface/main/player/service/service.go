package service

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	htemplate "html/template"
	"strconv"
	"strings"
	"text/template"
	"time"

	dmrpc "go-common/app/interface/main/dm2/rpc/client"
	hisrpc "go-common/app/interface/main/history/rpc/client"
	"go-common/app/interface/main/player/conf"
	"go-common/app/interface/main/player/dao"
	"go-common/app/interface/main/player/model"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	accclient "go-common/app/service/main/account/api"
	arcclient "go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	assrpc "go-common/app/service/main/assist/rpc/client"
	locrpc "go-common/app/service/main/location/rpc/client"
	resmdl "go-common/app/service/main/resource/model"
	resrpc "go-common/app/service/main/resource/rpc/client"
	ugcclient "go-common/app/service/main/ugcpay/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_resourceID = 2319
	_bgColor    = "#000000"
)

// Service is a service.
type Service struct {
	// config
	c *conf.Config
	// dao
	dao *dao.Dao
	// rpc
	arc *arcrpc.Service2
	his *hisrpc.Service
	ass *assrpc.Service
	res *resrpc.Service
	dm2 *dmrpc.Service
	loc *locrpc.Service
	tag *tagrpc.Service
	// memory cache
	caItems []*model.Item
	params  string
	icon    htemplate.HTML
	// template
	tWithU *template.Template
	tNoU   *template.Template
	// broadcast
	BrBegin time.Time
	BrEnd   time.Time
	// 拜年祭相关
	matOn    bool
	matTime  time.Time
	pastView *model.View
	matView  *model.View
	// vipQn
	vipQn map[int]int
	// grpc client
	accClient    accclient.AccountClient
	arcClient    arcclient.ArchiveClient
	ugcPayClient ugcclient.UGCPayClient
	// playurl paths
	playURL   string
	playURLV3 string
	h5PlayURL string
	highQaURL string
	// bnj2019 view map
	bnj2019ViewMap map[int64]*arcclient.ViewReply
}

// New  new and return service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
		arc: arcrpc.New2(c.ArchiveRPC),
		his: hisrpc.New(c.HistoryRPC),
		ass: assrpc.New(c.AssistRPC),
		res: resrpc.New(c.ResourceRPC),
		dm2: dmrpc.New(c.Dm2RPC),
		loc: locrpc.New(c.LocRPC),
		tag: tagrpc.New2(c.TagRPC),
	}
	s.playURL = c.Host.PlayurlCo + _playurlURI
	s.playURLV3 = c.Host.PlayurlCo + _playurlURIV3
	s.h5PlayURL = c.Host.H5Playurl + _h5PlayURI
	s.highQaURL = c.Host.HighPlayurl + _highQaURI
	var err error
	if s.accClient, err = accclient.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	if s.arcClient, err = arcclient.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	if s.ugcPayClient, err = ugcclient.NewClient(c.UGCPayClient); err != nil {
		panic(err)
	}
	s.initVipQn()
	s.caItems = []*model.Item{}
	// 模板
	s.tWithU, _ = template.New("with_user").Parse(model.TpWithUinfo)
	s.tNoU, _ = template.New("no_user").Parse(model.TpWithNoUinfo)
	// broadcast
	s.BrBegin, _ = time.Parse("2006-01-02 15:04:05", c.Broadcast.Begin)
	s.BrEnd, _ = time.Parse("2006-01-02 15:04:05", c.Broadcast.End)
	// 初始化播放器灰度配置项目
	go s.policyproc()
	go s.resourceproc()
	go s.paramproc()
	go s.iconproc()
	// 拜年祭
	s.matTime, _ = time.Parse(time.RFC3339, c.Matsuri.MatTime)
	s.matOn = true
	go s.matProc()
	go s.bnj2019Viewproc()
	return
}

// Ping check service health
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.Ping(c); err != nil {
		log.Error("s.dao.Ping() error(%v)", err)
	}
	return
}

func (s *Service) matProc() {
	var (
		ctx = context.Background()
		err error
	)
	for {
		var tmp *arcclient.ViewReply
		if tmp, err = s.arcClient.View(ctx, &arcclient.ViewRequest{Aid: s.c.Matsuri.PastID}); err != nil || tmp == nil {
			dao.PromError("View接口错误", "s.arcClientView3(%v) tmp(%v) error(%v) ", s.c.Matsuri.PastID, tmp, err)
			time.Sleep(time.Second)
			continue
		}
		s.pastView = &model.View{Arc: tmp.Arc, Pages: tmp.Pages}
		if tmp, err = s.arcClient.View(ctx, &arcclient.ViewRequest{Aid: s.c.Matsuri.MatID}); err != nil || tmp == nil {
			dao.PromError("View接口错误", "s.arcClientView3(%v) tmp(%v) error(%v)", s.c.Matsuri.MatID, tmp, err)
			time.Sleep(time.Second)
			continue
		}
		s.matView = &model.View{Arc: tmp.Arc, Pages: tmp.Pages}
		time.Sleep(time.Duration(s.c.Matsuri.Tick))
	}
}

func (s *Service) resourceproc() {
	for {
		var (
			items []*model.Item
			err   error
			res   *resmdl.Resource
		)
		if res, err = s.res.Resource(context.Background(), &resmdl.ArgRes{ResID: _resourceID}); err != nil {
			dao.PromError("Resource接口错误", "s.res.Resource(%d) error(%v)", _resourceID, err)
			time.Sleep(time.Second)
			continue
		}
		if res == nil {
			dao.PromError("Resource接口数据为空", "s.res.Resource(%d) is nil", _resourceID)
			time.Sleep(time.Second)
			continue
		}
		if len(res.Assignments) == 0 {
			dao.PromError("Resource接口Assignments数据为空", "s.res.Resource(%d) assignments is nil", _resourceID)
			time.Sleep(time.Second)
			continue
		}
		for _, v := range res.Assignments {
			item := &model.Item{
				Bgcolor:    _bgColor,
				ResourceID: strconv.Itoa(_resourceID),
				SrcID:      strconv.Itoa(v.ResID),
				ID:         strconv.Itoa(v.ID),
			}
			if catalog, ok := model.Catalog[v.PlayerCategory]; ok {
				item.Catalog = catalog
			}
			item.Content = string(rune(10)) + string(rune(13)) + fmt.Sprintf(_content, v.URL, v.Name) + string(rune(10)) + string(rune(13))
			items = append(items, item)
		}
		s.caItems = items
		time.Sleep(time.Duration(s.c.Tick.CarouselTick))
	}
}

func (s *Service) paramproc() {
	for {
		var (
			params []*model.Param
			items  []string
			err    error
		)
		c := context.Background()
		if params, err = s.dao.Param(c); err != nil {
			log.Error("s.dao.Param() error(%v)", err)
			time.Sleep(time.Second)
			continue
		}
		if len(params) == 0 {
			time.Sleep(time.Duration(s.c.Tick.ParamTick))
			continue
		}
		for _, pa := range params {
			nameBy := bytes.NewBuffer(nil)
			valueBy := bytes.NewBuffer(nil)
			if err = xml.EscapeText(nameBy, []byte(pa.Name)); err != nil {
				log.Error("xml.EscapeText(%s) error(%v)", pa.Name, err)
				continue
			} else {
				pa.Name = nameBy.String()
			}
			if err = xml.EscapeText(valueBy, []byte(pa.Value)); err != nil {
				log.Error("xml.EscapeText(%s) error(%v)", pa.Value, err)
				continue
			} else {
				pa.Value = valueBy.String()
			}
			item := "<" + pa.Name + ">" + pa.Value + "</" + pa.Name + ">"
			items = append(items, item)
		}
		if len(items) > 0 {
			s.params = strings.Join(items, "\n")
		}
		time.Sleep(time.Duration(s.c.Tick.ParamTick))
	}
}

func (s *Service) policyproc() {
	s.c.Policy.StartTime, _ = time.Parse("2006-01-02 15:04:05", s.c.Policy.Start)
	s.c.Policy.EndTime, _ = time.Parse("2006-01-02 15:04:05", s.c.Policy.End)
	s.c.Policy.MtimeTime, _ = time.Parse("2006-01-02 15:04:05", s.c.Policy.Mtime)
}

func (s *Service) iconproc() {
	for {
		icon, err := s.res.PlayerIcon(context.Background())
		if err != nil || icon == nil {
			log.Error("iconproc s.res.PlayerIcon error(%v) icon(%v)", err, icon)
			if ecode.Cause(err) == ecode.NothingFound {
				s.icon = ""
			}
			time.Sleep(time.Duration(s.c.Tick.IconTick))
			continue
		}
		icon.URL1 = strings.Replace(icon.URL1, "http://", "//", 1)
		icon.URL2 = strings.Replace(icon.URL2, "http://", "//", 1)
		bs, err := json.Marshal(icon)
		if err != nil {
			log.Error("iconproc json.Marshal(%v) error(%v)", icon, err)
			time.Sleep(time.Second)
			continue
		}
		s.icon = htemplate.HTML(bs)
		time.Sleep(time.Duration(s.c.Tick.IconTick))
	}
}

func (s *Service) initVipQn() {
	tmp := make(map[int]int, len(s.c.Rule.VipQn))
	for _, qn := range s.c.Rule.VipQn {
		tmp[qn] = qn
	}
	s.vipQn = tmp
}

func (s *Service) bnj2019Viewproc() {
	for {
		time.Sleep(time.Duration(s.c.Bnj2019.BnjTick))
		aids := append(s.c.Bnj2019.BnjListAids, s.c.Bnj2019.BnjMainAid)
		if len(aids) == 0 {
			continue
		}
		if views, err := s.arcClient.Views(context.Background(), &arcclient.ViewsRequest{Aids: aids}); err != nil || views == nil {
			log.Error("bnj2019Viewproc s.arcClient.Views(%v) error(%v)", aids, err)
			continue
		} else {
			tmp := make(map[int64]*arcclient.ViewReply, len(aids))
			for _, aid := range aids {
				if view, ok := views.Views[aid]; ok && view.Arc.IsNormal() {
					tmp[aid] = view
				}
			}
			s.bnj2019ViewMap = tmp
		}
	}
}
