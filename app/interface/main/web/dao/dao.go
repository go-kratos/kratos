package dao

import (
	"context"
	xhttp "net/http"
	"time"

	"go-common/app/interface/main/web/conf"
	"go-common/library/cache/redis"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

const (
	_rankURL          = "/data/rank/"
	_feedbackURL      = "/x/internal/feedback/ugc/add"
	_spaceTopPhotoURL = "/api/member/getTopPhoto"
	_coinAddURL       = "/x/coin/add"
	_coinExpURL       = "/x/coin/today/exp"
	_elecShowURL      = "/api/v2/rank/query/av"
	_arcReportURL     = "/videoup/archive/report"
	_arcAppealURL     = "/x/internal/workflow/appeal/add"
	_appealTagsURL    = "/x/internal/workflow/tag/list"
	_relatedURL       = "/recsys/related"
	_onlineURL        = "/x/internal/chat/server/ol"
	_liveOnlineURL    = "/room/v1/Online/get_total_online"
	_helpListURL      = "/kb/getQuestionTypeListByParentIdBilibili/4"
	_helpSearchURL    = "/kb/searchInerDocListBilibili/4"
	_onlineListURL    = "/x/internal/chat/num/top/aid"
	_searchURL        = "/main/search"
	_searchRecURL     = "/search/recommend"
	_searchDefaultURL = "/query/recommend"
	_searchUpRecURL   = "/main/recommend"
	_searchEggURI     = "/x/admin/feed/eggSearchWeb"
	_payWalletURL     = "/wallet-int/wallet/getUserWalletInfo"
	_payWalletOldURL  = "/wallet/api/v1/info"
	_special          = "/x/internal/uper/special"
)

// Dao dao
type Dao struct {
	// config
	c *conf.Config
	// http client
	httpR       *bm.Client
	httpW       *bm.Client
	httpBigData *bm.Client
	httpHelp    *bm.Client
	httpSearch  *bm.Client
	httpPay     *bm.Client
	bfsClient   *xhttp.Client
	// redis
	redis                  *redis.Pool
	redisBak               *redis.Pool
	redisNlExpire          int32
	redisNlBakExpire       int32
	redisRkExpire          int32
	redisRkBakExpire       int32
	redisDynamicBakExpire  int32
	redisArchiveBakExpire  int32
	redisTagBakExpire      int32
	redisCardBakExpire     int32
	redisRcExpire          int32
	redisRcBakExpire       int32
	redisArtBakExpire      int32
	redisIxIconExpire      int32
	redisIxIconBakExpire   int32
	redisHelpBakExpire     int32
	redisOlListBakExpire   int32
	redisWxHotExpire       int32
	redisWxHotBakExpire    int32
	redisAppealLimitExpire int32
	// bigdata url
	rankURL          string
	rankIndexURL     string
	rankRegionURL    string
	rankRecURL       string
	rankTagURL       string
	feedbackURL      string
	spaceTopPhotoURL string
	coinAddURL       string
	coinExpURL       string
	customURL        string
	elecShowURL      string
	arcReportURL     string
	appealTagsURL    string
	arcAppealURL     string
	relatedURL       string
	onlineURL        string
	liveOnlineURL    string
	helpListURL      string
	helpSearchURL    string
	onlineListURL    string
	shopURL          string
	replyHotURL      string
	searchURL        string
	searchRecURL     string
	searchDefaultURL string
	searchUpRecURL   string
	searchEggURL     string
	walletURL        string
	walletOldURL     string
	abServerURL      string
	wxHotURL         string
	bnjConfURL       string
	special          string
	// cache Prom
	cacheProm *prom.Prom
}

// New dao new
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// config
		c: c,
		// http read client
		httpR:       bm.NewClient(c.HTTPClient.Read),
		httpW:       bm.NewClient(c.HTTPClient.Write),
		httpBigData: bm.NewClient(c.HTTPClient.BigData),
		httpHelp:    bm.NewClient(c.HTTPClient.Help),
		httpSearch:  bm.NewClient(c.HTTPClient.Search),
		httpPay:     bm.NewClient(c.HTTPClient.Pay),
		// init bfs http client
		bfsClient: xhttp.DefaultClient,
		// redis
		redis:                  redis.NewPool(c.Redis.LocalRedis.Config),
		redisBak:               redis.NewPool(c.Redis.BakRedis.Config),
		redisNlExpire:          int32(time.Duration(c.Redis.LocalRedis.NewlistExpire) / time.Second),
		redisNlBakExpire:       int32(time.Duration(c.Redis.BakRedis.NewlistExpire) / time.Second),
		redisRkExpire:          int32(time.Duration(c.Redis.LocalRedis.RankingExpire) / time.Second),
		redisRkBakExpire:       int32(time.Duration(c.Redis.BakRedis.RankingExpire) / time.Second),
		redisDynamicBakExpire:  int32(time.Duration(c.Redis.BakRedis.RegionExpire) / time.Second),
		redisArchiveBakExpire:  int32(time.Duration(c.Redis.BakRedis.ArchiveExpire) / time.Second),
		redisTagBakExpire:      int32(time.Duration(c.Redis.BakRedis.TagExpire) / time.Second),
		redisCardBakExpire:     int32(time.Duration(c.Redis.BakRedis.CardExpire) / time.Second),
		redisRcExpire:          int32(time.Duration(c.Redis.LocalRedis.RcExpire) / time.Second),
		redisRcBakExpire:       int32(time.Duration(c.Redis.BakRedis.RcExpire) / time.Second),
		redisArtBakExpire:      int32(time.Duration(c.Redis.BakRedis.ArtUpExpire) / time.Second),
		redisIxIconExpire:      int32(time.Duration(c.Redis.LocalRedis.IndexIconExpire) / time.Second),
		redisIxIconBakExpire:   int32(time.Duration(c.Redis.BakRedis.IndexIconExpire) / time.Second),
		redisHelpBakExpire:     int32(time.Duration(c.Redis.BakRedis.HelpExpire) / time.Second),
		redisOlListBakExpire:   int32(time.Duration(c.Redis.BakRedis.OlListExpire) / time.Second),
		redisWxHotExpire:       int32(time.Duration(c.Redis.LocalRedis.WxHotExpire) / time.Second),
		redisWxHotBakExpire:    int32(time.Duration(c.Redis.BakRedis.WxHotExpire) / time.Second),
		redisAppealLimitExpire: int32(time.Duration(c.Redis.BakRedis.AppealLimitExpire) / time.Second),
		// remote source urls
		rankURL:          c.Host.Rank + _rankURL + _rankURI,
		rankIndexURL:     c.Host.Rank + _rankURL + _rankIndexURI,
		rankRegionURL:    c.Host.Rank + _rankURL + _rankRegionURI,
		rankRecURL:       c.Host.Rank + _rankURL + _rankRecURI,
		wxHotURL:         c.Host.Rank + _rankURL + _wxHotURI,
		rankTagURL:       c.Host.Rank + _rankTagURI,
		feedbackURL:      c.Host.API + _feedbackURL,
		spaceTopPhotoURL: c.Host.Space + _spaceTopPhotoURL,
		coinAddURL:       c.Host.API + _coinAddURL,
		coinExpURL:       c.Host.API + _coinExpURL,
		customURL:        c.Host.Rank + _rankURL + _customURI,
		elecShowURL:      c.Host.Elec + _elecShowURL,
		arcReportURL:     c.Host.ArcAPI + _arcReportURL,
		appealTagsURL:    c.Host.API + _appealTagsURL,
		arcAppealURL:     c.Host.API + _arcAppealURL,
		relatedURL:       c.Host.Data + _relatedURL,
		onlineURL:        c.Host.API + _onlineURL,
		liveOnlineURL:    c.Host.LiveAPI + _liveOnlineURL,
		helpListURL:      c.Host.HelpAPI + _helpListURL,
		helpSearchURL:    c.Host.HelpAPI + _helpSearchURL,
		onlineListURL:    c.Host.API + _onlineListURL,
		shopURL:          c.Host.Mall + _shopURI,
		replyHotURL:      c.Host.API + _hotURI,
		searchURL:        c.Host.Search + _searchURL,
		searchRecURL:     c.Host.Search + _searchRecURL,
		searchDefaultURL: c.Host.Search + _searchDefaultURL,
		searchUpRecURL:   c.Host.Search + _searchUpRecURL,
		searchEggURL:     c.Host.Manager + _searchEggURI,
		walletURL:        c.Host.Pay + _payWalletURL,
		walletOldURL:     c.Host.Pay + _payWalletOldURL,
		abServerURL:      c.Host.AbServer + _abServerURI,
		bnjConfURL:       c.Host.LiveAPI + _bnjConfURI,
		special:          c.Host.API + _special,
		// prom
		cacheProm: prom.CacheHit,
	}
	return d
}

// Ping check connection success.
func (dao *Dao) Ping(c context.Context) (err error) {
	if err = dao.pingRedis(c); err != nil {
		log.Error("dao.pingRedis error(%v)", err)
		return
	}
	if err = dao.pingRedisBak(c); err != nil {
		log.Error("dao.pingRedisBak error(%v)", err)
		return
	}
	return
}

// Close close  resource.
func (dao *Dao) Close() {
	if dao.redis != nil {
		dao.redis.Close()
	}
	if dao.redisBak != nil {
		dao.redisBak.Close()
	}
}

func (dao *Dao) pingRedis(c context.Context) (err error) {
	conn := dao.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

func (dao *Dao) pingRedisBak(c context.Context) (err error) {
	conn := dao.redisBak.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
