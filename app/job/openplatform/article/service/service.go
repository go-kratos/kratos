package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	tagrpc "go-common/app/interface/main/tag/rpc/client"
	artmdl "go-common/app/interface/openplatform/article/model"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	"go-common/app/job/openplatform/article/conf"
	"go-common/app/job/openplatform/article/dao"
	"go-common/app/job/openplatform/article/model"
	thumdl "go-common/app/service/main/thumbup/model"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
)

const (
	_sharding     = 100 // goroutines for dealing the stat
	_chanSize     = 10240
	_articleTable = "filtered_articles" // article table name
	_authorTable  = "article_authors"   // 作者表
)

type lastTimeStat struct {
	time int64
	stat *artmdl.StatMsg
}

// Service .
type Service struct {
	c      *conf.Config
	dao    *dao.Dao
	waiter *sync.WaitGroup
	closed bool
	// monitor              *monitor.Service
	articleSub           *databus.Databus
	articleStatSub       *databus.Databus
	likeStatSub          *databus.Databus
	replyStatSub         *databus.Databus
	favoriteStatSub      *databus.Databus
	coinStatSub          *databus.Databus
	articleRPC           *artrpc.Service
	tagRPC               *tagrpc.Service
	statCh               [_sharding]chan *artmdl.StatMsg
	statLastTime         [_sharding]map[int64]*lastTimeStat
	categoriesMap        map[int64]*artmdl.Category
	categoriesReverseMap map[int64][]*artmdl.Category
	sitemapMap           map[int64]struct{}
	urlListHead          *urlNode
	urlListTail          *urlNode
	lastURLNode          *urlNode
	sitemapXML           string
	likeCh               chan *thumdl.StatMsg
	updateDbInterval     int64
	updateSortInterval   time.Duration
	cheatInfoc           *infoc.Infoc
	// *Cnt means sum of consumed messages.
	statCnt, binCnt int64
	cheatArts       map[int64]int
	sortLimitTime   int64
	setting         *model.Setting
	// infoc
	logCh chan interface{}
}

// New creates a Service instance.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		waiter: new(sync.WaitGroup),
		// monitor:              monitor.New(),
		articleSub:           databus.New(c.ArticleSub),
		articleStatSub:       databus.New(c.ArticleStatSub),
		replyStatSub:         databus.New(c.ReplyStatSub),
		favoriteStatSub:      databus.New(c.FavoriteStatSub),
		coinStatSub:          databus.New(c.CoinStatSub),
		likeStatSub:          databus.New(c.LikeStatSub),
		articleRPC:           artrpc.New(c.ArticleRPC),
		tagRPC:               tagrpc.New2(c.TagRPC),
		updateDbInterval:     int64(time.Duration(c.Job.UpdateDbInterval) / time.Second),
		updateSortInterval:   time.Duration(c.Job.UpdateSortInterval),
		likeCh:               make(chan *thumdl.StatMsg, 1e4),
		cheatInfoc:           infoc.New(c.CheatInfoc),
		categoriesMap:        make(map[int64]*artmdl.Category),
		categoriesReverseMap: make(map[int64][]*artmdl.Category),
		sitemapMap:           make(map[int64]struct{}),
		sortLimitTime:        int64(time.Duration(c.Job.SortLimitTime) / time.Second),
		logCh:                make(chan interface{}, 1024),
	}
	// s.monitor.SetConfig(c.HTTPClient)
	for i := int64(0); i < _sharding; i++ {
		i := i
		// for stat
		s.statCh[i] = make(chan *artmdl.StatMsg, _chanSize)
		s.statLastTime[i] = make(map[int64]*lastTimeStat)
		s.waiter.Add(1)
		go s.statproc(i)
	}
	s.loadCategories()
	s.loadSettings()
	// 10 go routines with WaitGroup
	s.waiter.Add(10)
	go s.consumeStat()
	go s.consumeCanal()
	go s.checkConsumer()
	go s.retryStat()
	go s.retryReply()
	go s.retryGame()
	go s.retryFlow()
	go s.retryDynamic()
	go s.retryCDN()
	go s.retryCache()
	// other go routines
	go s.updateSortproc()
	go s.activityLikeproc()
	go s.updateReadCountproc()
	go s.updateHotspotsproc()
	go s.updateCheatArtsproc()
	go s.loadCategoriesproc()
	go s.loadSettingsproc()
	go s.recommendAuthorproc()
	go s.checkReadStatusProc()
	go s.infocproc()
	go s.sitemapproc()
	return
}

// consumeStat consumes article's stat.
func (s *Service) consumeStat() {
	defer s.waiter.Done()
	for {
		if s.closed {
			for i := 0; i < _sharding; i++ {
				close(s.statCh[i])
			}
			log.Info("databus: article-job stat consumer exit!")
			return
		}
	L:
		select {
		case msg, ok := <-s.articleStatSub.Messages():
			if !ok {
				break L
			}
			s.statCnt++
			msg.Commit()
			sm := &artmdl.StatMsg{}
			if err := json.Unmarshal(msg.Value, sm); err != nil {
				log.Error("json.Unmarshal(%s) error(%+v)", msg.Value, err)
				dao.PromError("service:解析计数databus消息")
				break L
			}
			if sm.Aid <= 0 {
				log.Warn("aid(%d) <=0 message(%s)", sm.Aid, msg.Value)
				break L
			}
			key := sm.Aid % _sharding
			s.statCh[key] <- sm
			prom.BusinessInfoCount.State(fmt.Sprintf("statChan-%v", key), int64(len(s.statCh[key])))
			log.Info("consumeStat key:%s partition:%d offset:%d msg: %v)", msg.Key, msg.Partition, msg.Offset, sm.String())
		case msg, ok := <-s.likeStatSub.Messages():
			if !ok {
				break L
			}
			msg.Commit()
			like := &thumdl.StatMsg{}
			if err := json.Unmarshal(msg.Value, like); err != nil {
				log.Error("json.Unmarshal(%s) error(%+v)", msg.Value, err)
				dao.PromError("service:解析like计数databus消息")
				break L
			}
			if like.Type != "article" {
				continue
			}
			select {
			case s.likeCh <- like:
			default:
			}
			key := like.ID % _sharding
			s.statCh[key] <- &artmdl.StatMsg{Like: &like.Count, Dislike: &like.DislikeCount, Aid: like.ID}
			prom.BusinessInfoCount.State(fmt.Sprintf("statChan-%v", key), int64(len(s.statCh[key])))
			log.Info("consumeLikeStat key:%s partition:%d offset:%d msg: %+v)", msg.Key, msg.Partition, msg.Offset, like)
		case msg, ok := <-s.replyStatSub.Messages():
			if !ok {
				break L
			}
			msg.Commit()
			m := &thumdl.StatMsg{}
			if err := json.Unmarshal(msg.Value, m); err != nil {
				log.Error("reply json.Unmarshal(%s) error(%+v)", msg.Value, err)
				dao.PromError("service:解析reply计数databus消息")
				break L
			}
			if m.Type != "article" {
				continue
			}
			key := m.ID % _sharding
			s.statCh[key] <- &artmdl.StatMsg{Reply: &m.Count, Aid: m.ID}
			prom.BusinessInfoCount.State(fmt.Sprintf("statChan-%v", key), int64(len(s.statCh[key])))
			log.Info("consumeReplyStat key:%s partition:%d offset:%d msg: %+v)", msg.Key, msg.Partition, msg.Offset, m)
		case msg, ok := <-s.favoriteStatSub.Messages():
			if !ok {
				break L
			}
			msg.Commit()
			m := &thumdl.StatMsg{}
			if err := json.Unmarshal(msg.Value, m); err != nil {
				log.Error("favorite json.Unmarshal(%s) error(%+v)", msg.Value, err)
				dao.PromError("service:解析favorite计数databus消息")
				break L
			}
			if m.Type != "article" {
				continue
			}
			key := m.ID % _sharding
			s.statCh[key] <- &artmdl.StatMsg{Favorite: &m.Count, Aid: m.ID}
			prom.BusinessInfoCount.State(fmt.Sprintf("statChan-%v", key), int64(len(s.statCh[key])))
			log.Info("consumeFavoriteStat key:%s partition:%d offset:%d msg: %+v)", msg.Key, msg.Partition, msg.Offset, m)
		case msg, ok := <-s.coinStatSub.Messages():
			if !ok {
				break L
			}
			msg.Commit()
			m := &thumdl.StatMsg{}
			if err := json.Unmarshal(msg.Value, m); err != nil {
				log.Error("coin json.Unmarshal(%s) error(%+v)", msg.Value, err)
				dao.PromError("service:解析coin计数databus消息")
				break L
			}
			if m.Type != "article" {
				continue
			}
			key := m.ID % _sharding
			s.statCh[key] <- &artmdl.StatMsg{Coin: &m.Count, Aid: m.ID}
			prom.BusinessInfoCount.State(fmt.Sprintf("statChan-%v", key), int64(len(s.statCh[key])))
			log.Info("consumeCoinStat key:%s partition:%d offset:%d msg: %+v)", msg.Key, msg.Partition, msg.Offset, m)
		}
	}
}

// consumeCanal consumes article's binlog databus.
func (s *Service) consumeCanal() {
	defer s.waiter.Done()
	var c = context.TODO()
	for {
		msg, ok := <-s.articleSub.Messages()
		if !ok {
			log.Info("databus: article-job binlog consumer exit!")
			return
		}
		s.binCnt++
		msg.Commit()
		m := &model.Message{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%+v)", msg.Value, err)
			continue
		}
		switch m.Table {
		case _articleTable:
			s.upArticles(c, m.Action, m.New, m.Old)
		case _authorTable:
			s.upAuthors(c, m.Action, m.New, m.Old)
		}
		log.Info("consumeCanal key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
	}
}

func (s *Service) retryStat() {
	defer s.waiter.Done()
	var (
		err error
		bs  []byte
		c   = context.TODO()
	)
	for {
		if s.closed {
			return
		}
		bs, err = s.dao.PopStat(c)
		if err != nil || bs == nil {
			time.Sleep(time.Second)
			continue
		}
		msg := &dao.StatRetry{}
		if err = json.Unmarshal(bs, msg); err != nil {
			log.Error("json.Unmarshal(%s) error(%+v)", bs, err)
			dao.PromError("service:解析计数重试消息")
			continue
		}
		if msg.Count > dao.RetryStatCount {
			continue
		}
		log.Info("retry: %s", bs)
		switch msg.Action {
		case dao.RetryUpdateStatCache:
			if err = s.updateCache(c, msg.Data, msg.Count+1); err != nil {
				time.Sleep(100 * time.Millisecond)
			}
			dao.PromInfo("service:重试计数更新缓存")
		case dao.RetryUpdateStatDB:
			if err = s.updateDB(c, msg.Data, msg.Count+1); err != nil {
				time.Sleep(100 * time.Millisecond)
			}
			dao.PromInfo("service:重试计数更新DB")
		}
	}
}

func (s *Service) retryReply() {
	defer s.waiter.Done()
	var (
		err      error
		aid, mid int64
		c        = context.TODO()
	)
	for {
		if s.closed {
			return
		}
		aid, mid, err = s.dao.PopReply(c)
		if err != nil || aid == 0 || mid == 0 {
			time.Sleep(time.Second)
			continue
		}
		log.Info("retry reply: aid(%d) mid(%d)", aid, mid)
		if err = s.openReply(c, aid, mid); err != nil {
			log.Error("s.openReply(%d,%d) error(%+v)", aid, mid, err)
			dao.PromInfo("service:重试打开评论区")
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Info("s.openReply(%d,%d) retry success", aid, mid)
	}
}

func (s *Service) retryCDN() {
	defer s.waiter.Done()
	var (
		err  error
		file string
		c    = context.TODO()
	)
	for {
		if s.closed {
			return
		}
		file, err = s.dao.PopCDN(c)
		if err != nil || file == "" {
			time.Sleep(time.Second)
			continue
		}
		log.Info("retry CDN: file(%s)", file)
		if err = s.dao.PurgeCDN(c, file); err != nil {
			log.Error("s.dao.PurgeCDN(%s) error(%+v)", file, err)
			dao.PromInfo("service:刷新CDN重试")
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Info("s.dao.PurgeCDN(%s) retry success.", file)
	}
}

func (s *Service) retryCache() {
	defer s.waiter.Done()
	var (
		err error
		bs  []byte
		c   = context.TODO()
	)
	for {
		if s.closed {
			return
		}
		bs, err = s.dao.PopArtCache(c)
		if err != nil || bs == nil {
			time.Sleep(time.Second)
			continue
		}
		msg := &dao.CacheRetry{}
		if err = json.Unmarshal(bs, msg); err != nil {
			log.Error("json.Unmarshal(%s) error(%+v)", bs, err)
			dao.PromError("service:解析文章缓存重试消息")
			continue
		}
		log.Info("retry cache: %s", bs)
		switch msg.Action {
		case dao.RetryAddArtCache:
			if err = s.addArtCache(c, msg.Aid); err != nil {
				time.Sleep(100 * time.Millisecond)
			}
			dao.PromInfo("service:重试添加文章缓存")
		case dao.RetryUpdateArtCache:
			if err = s.updateArtCache(c, msg.Aid, msg.Cid); err != nil {
				time.Sleep(100 * time.Millisecond)
			}
			dao.PromInfo("service:重试更新文章缓存")
		case dao.RetryDeleteArtCache:
			if err = s.deleteArtCache(c, msg.Aid, msg.Mid); err != nil {
				time.Sleep(100 * time.Millisecond)
			}
			dao.PromInfo("service:重试删除文章缓存")
		case dao.RetryDeleteArtRecCache:
			if err = s.deleteArtRecommendCache(c, msg.Aid, msg.Cid); err != nil {
				time.Sleep(100 * time.Millisecond)
			}
			dao.PromInfo("service:重试删除文章推荐缓存")
		}
	}
}

// checkConsumer checks consumer state.
func (s *Service) checkConsumer() {
	defer s.waiter.Done()
	if env.DeployEnv != env.DeployEnvProd {
		return
	}
	var (
		// ctx = context.TODO()
		// c* means sum of consumed messages.
		c1, c2 int64
	)
	for {
		time.Sleep(5 * time.Hour)
		if s.statCnt-c1 == 0 {
			// msg := "databus: article-job stat did not consume within a minute"
			// s.monitor.Sms(ctx, s.c.SMS.Phone, s.c.SMS.Token, msg)
			// log.Warn(msg)
			log.Warn("databus: article-job stat did not consume within a minute")
		}
		c1 = s.statCnt
		if s.binCnt-c2 == 0 {
			// msg := "databus: article-job binlog did not consume within a minute"
			// s.monitor.Sms(ctx, s.c.SMS.Phone, s.c.SMS.Token, msg)
			// log.Warn(msg)
			log.Warn("databus: article-job binlog did not consume within a minute")
		}
		c2 = s.binCnt
	}
}

// Ping reports the heath of services.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close releases resources which owned by the Service instance.
func (s *Service) Close() (err error) {
	defer s.waiter.Wait()
	s.articleSub.Close()
	s.articleStatSub.Close()
	s.closed = true
	log.Info("article-job has been closed.")
	return
}

func (s *Service) retryGame() {
	defer s.waiter.Done()
	var (
		err  error
		info *model.GameCacheRetry
		c    = context.TODO()
	)
	for {
		if s.closed {
			return
		}
		info, err = s.dao.PopGameCache(c)
		if (err != nil) || (info == nil) {
			time.Sleep(time.Second)
			continue
		}
		for {
			log.Info("retry_game: %+v", info)
			dao.PromInfo("service:重试同步游戏")
			if err = s.dao.GameSync(c, info.Action, info.Aid); err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			break
		}
		log.Info("s.GameSync(%s, aid: %d) retry success", info.Action, info.Aid)
	}
}

func (s *Service) activityLikeproc() {
	var c = context.TODO()
	for {
		msg, ok := <-s.likeCh
		if !ok {
			log.Info("activityLikeproc: exit!")
			return
		}
		s.dao.LikeSync(c, msg.ID, msg.Count)
	}
}

func (s *Service) retryFlow() {
	defer s.waiter.Done()
	var (
		err  error
		info *model.FlowCacheRetry
		c    = context.TODO()
	)
	for {
		if s.closed {
			return
		}
		info, err = s.dao.PopFlowCache(c)
		if (err != nil) || (info == nil) {
			time.Sleep(time.Second)
			continue
		}
		for {
			log.Info("retry_flow: %+v", info)
			dao.PromInfo("service:重试同步flow")
			if err = s.dao.FlowSync(c, info.Mid, info.Aid); err != nil {
				time.Sleep(time.Second * 1)
				continue
			}
			break
		}
		log.Info("s.FlowSync(mid:%v, aid: %d) retry success", info.Mid, info.Aid)
	}
}

func (s *Service) retryDynamic() {
	defer s.waiter.Done()
	var (
		err  error
		info *model.DynamicCacheRetry
		c    = context.TODO()
	)
	for {
		if s.closed {
			return
		}
		info, err = s.dao.PopDynamicCache(c)
		if (err != nil) || (info == nil) {
			time.Sleep(time.Second)
			continue
		}
		for {
			log.Info("retry_dynamic: %+v", info)
			dao.PromInfo("service:重试同步dynamic")
			if err = s.dao.PubDynamic(c, info.Mid, info.Aid, info.Show, info.Comment, info.Ts, info.DynamicIntro); err != nil {
				time.Sleep(time.Second * 1)
				continue
			}
			break
		}
		log.Info("s.PubDynamic(mid:%v, aid: %d) retry success %+v", info.Mid, info.Aid, info)
	}
}

func (s *Service) updateReadCountproc() {
	var c = context.TODO()
	for {
		err := s.articleRPC.RebuildAllListReadCount(c)
		if err != nil {
			log.Error("s.updateReadCountproc() err: %+v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		log.Info("s.updateReadCountproc() success")
		dao.PromError("更新文集阅读数")
		time.Sleep(time.Duration(s.c.Job.ListReadCountInterval))
	}
}

func (s *Service) updateHotspotsproc() {
	var c = context.TODO()
	var lastUpdate int64
	for {
		var force bool
		duration := int64(time.Duration(s.c.Job.HotspotForceInterval) / time.Second)
		if (time.Now().Unix() - lastUpdate) > duration {
			force = true
		}
		err := s.articleRPC.UpdateHotspots(c, &artmdl.ArgForce{Force: force})
		if err != nil {
			log.Error("s.UpdateHotspots() err: %+v", err)
			dao.PromError("更新热点运营文章")
			time.Sleep(time.Second * 5)
			continue
		}
		if force {
			lastUpdate = time.Now().Unix()
		}
		log.Info("s.UpdateHotspots() success force:%v", force)
		time.Sleep(time.Duration(s.c.Job.HotspotInterval))
	}
}

func (s *Service) updateCheatArtsproc() {
	var c = context.TODO()
	for {
		arts, err := s.dao.CheatArts(c)
		if err != nil {
			log.Error("s.updateCheatArtsproc() err: %+v", err)
			dao.PromError("更新反作弊文章列表")
			time.Sleep(time.Second * 5)
			continue
		}
		log.Info("s.updateCheatArtsproc() success, len: %v", len(arts))
		s.cheatArts = arts
		time.Sleep(time.Minute)
	}
}
