package service

import (
	"context"
	"encoding/json"
	"runtime"
	"strconv"
	"time"

	"go-common/app/interface/main/reply/conf"
	"go-common/app/interface/main/reply/dao/bigdata"
	"go-common/app/interface/main/reply/dao/drawyoo"
	"go-common/app/interface/main/reply/dao/fans"
	"go-common/app/interface/main/reply/dao/reply"
	"go-common/app/interface/main/reply/dao/search"
	"go-common/app/interface/main/reply/dao/vip"
	"go-common/app/interface/main/reply/dao/workflow"
	model "go-common/app/interface/main/reply/model/reply"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	accrpc "go-common/app/service/main/account/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	assrpc "go-common/app/service/main/assist/rpc/client"
	figrpc "go-common/app/service/main/figure/rpc/client"
	filgrpc "go-common/app/service/main/filter/api/grpc/v1"
	replyfeed "go-common/app/service/main/reply-feed/api"

	locrpc "go-common/app/service/main/location/rpc/client"
	seqmdl "go-common/app/service/main/seq-server/model"
	seqrpc "go-common/app/service/main/seq-server/rpc/client"
	thumbup "go-common/app/service/main/thumbup/rpc/client"
	ugcpay "go-common/app/service/main/ugcpay/api/grpc/v1"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/sync/pipeline/fanout"
)

var _defMemberJSON = []byte(`{"mid":"0","uname":"","sex":"保密","sign":"没签名","avatar":"http://static.hdslb.com/images/member/noface.gif","rank":"10000","DisplayRank":"0","level_info":{"current_level":0,"current_min":0,"current_exp":0,"next_exp":0},"pendant":{"pid":0,"name":"","image":"","expire":0},"nameplate":{"nid":0,"name":"","image":"","image_small":"","level":"","condition":""},"official_verify":{"type":-1,"desc":""},"vip":{"vipType":0,"vipDueDate":0,"dueRemark":"","accessStatus":1,"vipStatus":0,"vipStatusWarn":""}}`)

// Service is reply service
type Service struct {
	sndDefCnt  int
	oneDaySec  int // If block time from word filter is greater than one day, then banned reply and add moral to account.
	useBigData bool
	// content filter
	bigdata *bigdata.Dao

	drawyoo *drawyoo.Dao

	// fans http get method
	fans *fans.Dao

	//  http request
	search *search.Dao

	// rpc client
	acc        accrpc.AccountClient
	feedClient replyfeed.ReplyFeedClient
	filcli     filgrpc.FilterClient
	location   *locrpc.Service
	assist     *assrpc.Service
	figure     *figrpc.Service
	thumbup    thumbup.ThumbupRPC
	// seq rpc
	seqArg *seqmdl.ArgBusiness
	seqSrv *seqrpc.Service2
	// workflow dao
	workflow   *workflow.Dao
	arcSrv     *arcrpc.Service2
	articleSrv *artrpc.Service
	ugcpay     ugcpay.UGCPayClient
	// reply cache dao
	dao *reply.Dao
	// asynchronized add cache use channel
	replyChan chan replyChan
	topRpChan chan topRpChan
	// vip
	vip     *vip.Dao
	emojisM map[string]int64
	emojis  []*model.EmojiPackage
	// reply notice
	notice map[int8][]*model.Notice
	// default member
	defMember *model.Info
	cache     *fanout.Fanout

	typeMapping  map[int32]string
	aliasMapping map[string]int32
	aidWhiteList []int64

	infoc          *infoc.Infoc
	bnjAid         map[int64]struct{}
	sortByHot      map[int64]int8
	sortByTime     map[int64]int8
	hideFloor      map[int64]int8
	hotReplyConfig map[int8]map[int64]int
}

// New new a service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		cache:          fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		sndDefCnt:      conf.Conf.Reply.SecondDefSize,
		oneDaySec:      24 * 60 * 60,
		useBigData:     conf.Conf.Reply.BigdataFilter,
		replyChan:      make(chan replyChan, _replyChanBuf),
		topRpChan:      make(chan topRpChan, _topRpChanBuf),
		aidWhiteList:   c.Reply.AidWhiteList,
		arcSrv:         arcrpc.New2(c.RPCClient2.Archive),
		articleSrv:     artrpc.New(c.RPCClient2.Article),
		infoc:          infoc.New(c.Infoc),
		bnjAid:         make(map[int64]struct{}),
		sortByTime:     make(map[int64]int8),
		sortByHot:      make(map[int64]int8),
		hideFloor:      make(map[int64]int8),
		hotReplyConfig: make(map[int8]map[int64]int),
	}
	for oidStr, tp := range c.Reply.SortByHotOids {
		var (
			oid int64
			err error
		)
		if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
			panic(err)
		}
		s.sortByHot[oid] = tp
	}
	for oidStr, tp := range c.Reply.SortByTimeOids {
		var (
			oid int64
			err error
		)
		if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
			panic(err)
		}
		s.sortByTime[oid] = tp
	}
	for oidStr, tp := range c.Reply.HideFloorOids {
		var (
			oid int64
			err error
		)
		if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
			panic(err)
		}
		s.hideFloor[oid] = tp
	}
	for tpStr, mapping := range c.Reply.HotReplyConfig {
		var (
			tp  int64
			oid int64
			err error
		)
		if tp, err = strconv.ParseInt(tpStr, 10, 64); err != nil {
			panic(err)
		}
		businessConfig := make(map[int64]int)
		for oidStr, num := range mapping {
			if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
				panic(err)
			}
			businessConfig[oid] = num
		}
		s.hotReplyConfig[int8(tp)] = businessConfig
	}
	for _, oid := range c.Reply.BnjAidList {
		s.bnjAid[oid] = struct{}{}
	}
	acc, err := accrpc.NewClient(c.AccountGRPCClient)
	if err != nil {
		panic(err)
	}
	s.acc = acc
	filcli, err := filgrpc.NewClient(c.FilterGRPCClient)
	if err != nil {
		panic(err)
	}
	s.filcli = filcli
	feedClient, err := replyfeed.NewClient(c.FeedGRPCClient)
	if err != nil {
		panic(err)
	}
	s.feedClient = feedClient
	// init depend dao
	s.fans = fans.New(c)
	s.search = search.New(c)
	s.drawyoo = drawyoo.New(c)
	s.bigdata = bigdata.New(c)
	s.vip = vip.New(c)
	s.emojisM = make(map[string]int64)
	s.emojis = make([]*model.EmojiPackage, 0)
	// init rpc client
	s.location = locrpc.New(c.RPCClient2.Location)
	s.assist = assrpc.New(c.RPCClient2.Assist)
	s.figure = figrpc.New(c.RPCClient2.Figure)
	s.seqSrv = seqrpc.New2(c.RPCClient2.Seq)
	s.seqArg = &seqmdl.ArgBusiness{Token: c.Seq.Token, BusinessID: c.Seq.BusinessID}
	s.thumbup = thumbup.New(c.RPCClient2.Thumbup)
	s.workflow = workflow.New(c)
	s.ugcpay, err = ugcpay.NewClient(nil)
	if err != nil {
		panic(err)
	}
	// init reply cache dao
	s.dao = reply.New(c)
	// default member
	s.defMember = &model.Info{}
	s.typeMapping = make(map[int32]string)
	s.aliasMapping = make(map[string]int32)
	if err := json.Unmarshal(_defMemberJSON, s.defMember); err != nil {
		panic(err)
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		go s.cacheproc()
	}
	s.loadEmoji()
	s.loadRplNotice()
	s.loadBusiness()
	go s.loadproc()
	return
}

// TypeToAlias map type to alias
func (s *Service) TypeToAlias(t int32) (alias string, exists bool) {
	alias, exists = s.typeMapping[t]
	return
}

// AliasToType map alias to type
func (s *Service) AliasToType(alias string) (t int32, exists bool) {
	t, exists = s.aliasMapping[alias]
	return
}

func (s *Service) loadproc() {
	for {
		s.loadEmoji()
		s.loadRplNotice()
		s.loadBusiness()
		time.Sleep(time.Duration(conf.Conf.Reply.EmojiExpire))
	}
}

// Ping check service health
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close close service connection
func (s *Service) Close() {
	s.dao.Close()
}

// ResetFloor ...
func (s *Service) ResetFloor(rps ...*model.Reply) {
	for _, rp := range rps {
		if rp == nil {
			continue
		}
		rp.Floor = 0
		if rp.Replies == nil {
			continue
		}
		for _, childRp := range rp.Replies {
			if childRp == nil {
				continue
			}
			childRp.Floor = 0
		}
	}
}

// hotNum ...
func (s *Service) hotNumWeb(oid int64, tp int8) int {
	if businessConfig, ok := s.hotReplyConfig[tp]; ok {
		if num, exists := businessConfig[oid]; exists {
			return num
		}
	}
	return _hotSizeWeb
}

// hotNum ...
func (s *Service) hotNum(oid int64, tp int8) int {
	if businessConfig, ok := s.hotReplyConfig[tp]; ok {
		if num, exists := businessConfig[oid]; exists {
			return num
		}
	}
	return _hotSize
}

// IsBnj avid...
func (s *Service) IsBnj(oid int64, tp int8) bool {
	if _, ok := s.bnjAid[oid]; ok && tp == model.SubTypeArchive {
		return true
	}
	return false
}

// ShowFloor ...
func (s *Service) ShowFloor(oid int64, tp int8) bool {
	if typ, ok := s.hideFloor[oid]; ok && typ == tp {
		return false
	}
	return true
}

// loadEmoji load emoji
func (s *Service) loadEmoji() (err error) {
	var (
		emojis = make([]*model.EmojiPackage, 0)
		c      = context.Background()
		emjM   = make(map[string]int64)
	)
	emjs, err := s.dao.Emoji.ListEmojiPack(c)
	if err != nil {
		log.Error("service.ListEmojiPack error (%v)", err)
		return err
	}
	for _, d := range emjs {
		d.Emojis, err = s.dao.Emoji.EmojiListByPid(c, d.ID)
		if err != nil {
			log.Error("service.EmojiListByPid err (%v)", err)
			return err
		}
		for _, emo := range d.Emojis {
			emjM[emo.Name] = emo.ID
		}
		if len(d.Emojis) > 0 {
			emojis = append(emojis, d)
		}
	}
	s.emojis = emojis
	s.emojisM = emjM
	return
}
