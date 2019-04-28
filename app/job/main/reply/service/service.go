package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	"go-common/app/job/main/reply/conf"
	"go-common/app/job/main/reply/dao/message"
	"go-common/app/job/main/reply/dao/notice"
	"go-common/app/job/main/reply/dao/reply"
	"go-common/app/job/main/reply/dao/search"
	"go-common/app/job/main/reply/dao/spam"
	"go-common/app/job/main/reply/dao/stat"
	model "go-common/app/job/main/reply/model/reply"
	accrpc "go-common/app/service/main/account/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	assrpc "go-common/app/service/main/assist/rpc/client"
	eprpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	es "go-common/library/database/elastic"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"
)

const (
	_chLen = 2048
)

var (
	_rpChs   []chan *databus.Message
	_likeChs []chan *databus.Message
)

// action the message struct of kafka
type consumerMsg struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type searchFlush struct {
	OldState int8
	Reply    *model.Reply
	Report   *model.Report
}

func (s *searchFlush) Key() (key string) {
	if s.Report != nil {
		return fmt.Sprintf("%d%d", s.Report.RpID, s.Report.ID)
	}
	return fmt.Sprintf("%d", s.Reply.RpID)
}

// Service is reply-job service
type Service struct {
	c            *conf.Config
	waiter       *sync.WaitGroup
	dataConsumer *databus.Databus
	likeConsumer *databus.Databus
	searchChan   chan *searchFlush

	// rpc client
	accSrv     accrpc.AccountClient
	arcSrv     *arcrpc.Service2
	articleSrv *artrpc.Service
	assistSrv  *assrpc.Service
	bangumiSrv eprpc.EpisodeClient
	// depend
	messageDao *message.Dao
	// notice
	noticeDao *notice.Dao
	// stat
	statDao *stat.Dao
	// reply
	dao *reply.Dao
	// spam
	spam *spam.Cache
	// search
	searchDao   *search.Dao
	batchNumber int
	es          *es.Elastic

	notify *fanout.Fanout

	typeMapping  map[int32]string
	aliasMapping map[string]int32
	marker       *fanout.Fanout
}

// New return new service
func New(c *conf.Config) (s *Service) {
	if c.Job.BatchNumber <= 0 {
		c.Job.BatchNumber = 2000
	}
	searchHTTPClient = xhttp.NewClient(c.HTTPClient)
	wardenClient := warden.DefaultClient()
	cc, err := wardenClient.Dial(context.Background(), "discovery://default/season.service")
	if err != nil {
		panic(err)
	}
	bangumiClient := eprpc.NewEpisodeClient(cc)
	s = &Service{
		c:            c,
		bangumiSrv:   bangumiClient,
		waiter:       new(sync.WaitGroup),
		searchChan:   make(chan *searchFlush, 1024),
		dataConsumer: databus.New(c.Databus.Consumer),
		likeConsumer: databus.New(c.Databus.Like),
		//rpc
		arcSrv:     arcrpc.New2(c.RPCClient2.Archive),
		articleSrv: artrpc.New(c.RPCClient2.Article),
		assistSrv:  assrpc.New(c.RPCClient2.Assist),
		messageDao: message.NewMessageDao(c),
		searchDao:  search.New(c),
		noticeDao:  notice.New(c),
		// stat
		statDao: stat.New(c),
		// init reply dao
		dao: reply.New(c),
		// init spam cache
		batchNumber:  c.Job.BatchNumber,
		spam:         spam.NewCache(c.Redis.Config),
		notify:       fanout.New("cache", fanout.Worker(1), fanout.Buffer(2048)),
		typeMapping:  make(map[int32]string),
		aliasMapping: make(map[string]int32),
		es:           es.NewElastic(c.Es),
		marker:       fanout.New("marker", fanout.Worker(1), fanout.Buffer(1024)),
	}
	accSvc, err := accrpc.NewClient(c.AccountClient)
	if err != nil {
		panic(err)
	}
	s.accSrv = accSvc
	time.Sleep(time.Second)
	_rpChs = make([]chan *databus.Message, c.Job.Proc)
	_likeChs = make([]chan *databus.Message, c.Job.Proc)
	for i := 0; i < c.Job.Proc; i++ {
		_rpChs[i] = make(chan *databus.Message, _chLen)
		_likeChs[i] = make(chan *databus.Message, _chLen)
		s.waiter.Add(1)
		go s.consumeproc(i)
		s.waiter.Add(1)
		go s.consumelikeproc(i)
	}
	s.waiter.Add(1)
	go s.likeConsume()
	s.waiter.Add(1)
	go s.dataConsume()

	go s.searchproc()
	go s.mappingproc()
	return
}

func (s *Service) addSearchUp(c context.Context, oldState int8, rp *model.Reply, rpt *model.Report) {
	select {
	case s.searchChan <- &searchFlush{OldState: oldState, Reply: rp, Report: rpt}:
	default:
		log.Error("addSearchUp chan full, type:%d oid:%d rpID:%d", rp.Type, rp.Oid, rp.RpID)
	}
}

func (s *Service) searchproc() {
	var (
		m      *searchFlush
		merge  = make(map[string]*searchFlush)
		num    = s.c.Job.SearchNum
		ticker = time.NewTicker(time.Duration(s.c.Job.SearchFlush))
	)
	for {
		select {
		case m = <-s.searchChan:
			merge[m.Key()] = m
			if len(merge) < num {
				continue
			}
		case <-ticker.C:
		}
		if len(merge) > 0 {
			s.callSearchUp(context.Background(), merge)
			merge = make(map[string]*searchFlush)
		}
	}
}

func (s *Service) likeConsume() {
	defer func() {
		s.waiter.Done()
		for i := 0; i < s.c.Job.Proc; i++ {
			close(_rpChs[i])
		}
	}()
	msgs := s.likeConsumer.Messages()
	for {
		msg, ok := <-msgs
		if !ok {
			log.Warn("[service.dataConsume|reply] dataConsumer has been closed.")
			return
		}
		if msg.Topic != s.c.Databus.Like.Topic {
			continue
		}
		rpid, err := strconv.ParseInt(string(msg.Key), 10, 64)
		if err != nil {
			continue
		}
		_likeChs[rpid%int64(s.c.Job.Proc)] <- msg
	}
}

func (s *Service) dataConsume() {
	defer func() {
		s.waiter.Done()
		for i := 0; i < s.c.Job.Proc; i++ {
			close(_rpChs[i])
		}
	}()
	msgs := s.dataConsumer.Messages()
	for {
		msg, ok := <-msgs
		if !ok {
			log.Warn("[service.dataConsume|reply] dataConsumer has been closed.")
			return
		}
		if msg.Topic != s.c.Databus.Consumer.Topic {
			continue
		}
		oid, err := strconv.ParseInt(string(msg.Key), 10, 64)
		if err != nil {
			continue
		}
		_rpChs[oid%int64(s.c.Job.Proc)] <- msg
	}
}

// StatMsg stat msg.
type StatMsg struct {
	Type         string `json:"type,omitempty"`
	ID           int64  `json:"id,omitempty"`
	Count        int    `json:"count,omitempty"`
	Oid          int64  `json:"origin_id,omitempty"`
	DislikeCount int    `json:"dislike_count,omitempty"`
	Timestamp    int64  `json:"timestamp,omitempty"`
	Mid          int64  `json:"mid,omitempty"`
}

func (s *Service) consumelikeproc(i int) {
	defer s.waiter.Done()
	for {
		msg, ok := <-_likeChs[i]
		if !ok {
			log.Info("consumeproc exit")
			return
		}
		cmsg := &StatMsg{}
		if err := json.Unmarshal(msg.Value, cmsg); err != nil {
			log.Error("json.Unmarshal() error(%v)", err)
			continue
		}
		if cmsg.Type != "reply" {
			continue
		}
		s.setLike(context.Background(), cmsg)
		msg.Commit()
		log.Info("consumer topic:%s, partitionId:%d, offset:%d, Key:%s, Value:%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
	}
}

func (s *Service) consumeproc(i int) {
	defer s.waiter.Done()
	for {
		msg, ok := <-_rpChs[i]
		if !ok {
			log.Info("consumeproc exit")
			return
		}
		cmsg := &consumerMsg{}
		if err := json.Unmarshal(msg.Value, cmsg); err != nil {
			log.Error("json.Unmarshal() error(%v)", err)
			continue
		}
		switch cmsg.Action {
		case "add":
			s.actionAdd(context.Background(), cmsg)
		case "add_top":
			s.addTopCache(context.Background(), cmsg)
		case "rpt":
			s.actionRpt(context.Background(), cmsg)
		case "act":
			//s.actionAct(context.Background(), cmsg)
			s.recAct(context.Background(), cmsg)
		case "re_idx":
			s.actionRecoverIndex(context.Background(), cmsg)
		case "idx_floor":
			s.acionRecoverFloorIdx(context.Background(), cmsg)
		case "re_rt_idx":
			s.actionRecoverRootIndex(context.Background(), cmsg)
		case "idx_dialog":
			s.actionRecoverDialog(context.Background(), cmsg)
		case "fix_dialog":
			s.actionRecoverFixDialog(context.Background(), cmsg)
		case "re_act":
			// s.actionRecoverAction(context.Background(),cmsg)
		case "up":
			s.actionUp(context.Background(), cmsg)
		case "admin":
			s.actionAdmin(context.Background(), cmsg)
		case "spam":
			s.addRecReply(context.Background(), cmsg)
			s.addDailyReply(context.Background(), cmsg)
		case "folder":
			s.folderHanlder(context.Background(), cmsg)
		default:
			log.Error("invalid action %s, cmsg is %v", cmsg.Action, cmsg)
		}
		msg.Commit()
		log.Info("consumer topic:%s, partitionId:%d, offset:%d, Key:%s, Value:%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
	}
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

func (s *Service) mappingproc() {
	for {
		if business, err := s.ListBusiness(context.Background()); err != nil {
			log.Error("s.ListBusiness error(%v)", err)
		} else {
			for _, b := range business {
				s.typeMapping[b.Type] = b.Alias
				s.aliasMapping[b.Alias] = b.Type
			}
		}
		time.Sleep(time.Duration(time.Minute * 5))
	}
}

// Close close service
func (s *Service) Close() error {
	return s.dataConsumer.Close()
}

// Wait wait all chan close
func (s *Service) Wait() {
	s.waiter.Wait()
}

// Ping check service health
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}
