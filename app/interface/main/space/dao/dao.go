package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/space/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"

	"go-common/library/database/hbase.v2"
)

// Dao dao struct.
type Dao struct {
	// config
	c *conf.Config
	// db
	db *sql.DB
	// hbase
	hbase *hbase.Client
	// stmt
	channelStmt       []*sql.Stmt
	channelListStmt   []*sql.Stmt
	channelCntStmt    []*sql.Stmt
	channelArcCntStmt []*sql.Stmt
	// redis
	redis *redis.Pool
	// mc
	mc *memcache.Pool
	// http client
	httpR    *bm.Client
	httpW    *bm.Client
	httpGame *bm.Client
	// api URL
	bangumiURL          string
	bangumiConcernURL   string
	bangumiUnConcernURL string
	favFolderURL        string
	favArcURL           string
	favAlbumURL         string
	favMovieURL         string
	shopURL             string
	shopLinkURL         string
	albumCountURL       string
	albumListURL        string
	tagSubURL           string
	tagCancelSubURL     string
	tagSubListURL       string
	accTagsURL          string
	accTagsSetURL       string
	isAnsweredURL       string
	lastPlayGameURL     string
	appPlayedGameURL    string
	arcSearchURL        string
	webTopPhotoURL      string
	topPhotoURL         string
	liveMetalURL        string
	liveURL             string
	medalStatusURL      string
	groupsCountURL      string
	elecURL             string
	audioCardURL        string
	audioUpperCertURL   string
	audioCntURL         string
	dynamicListURL      string
	dynamicURL          string
	dynamicCntURL       string
	// expire
	clExpire        int32
	upArtExpire     int32
	upArcExpire     int32
	mcSettingExpire int32
	mcNoticeExpire  int32
	mcTopArcExpire  int32
	mcMpExpire      int32
	mcThemeExpire   int32
	mcTopDyExpire   int32
	// cache
	cache *fanout.Fanout
}

// New new dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// config
		c:                   c,
		db:                  sql.NewMySQL(c.Mysql),
		hbase:               hbase.NewClient(c.HBase.Config),
		redis:               redis.NewPool(c.Redis.Config),
		mc:                  memcache.NewPool(c.Memcache.Config),
		httpR:               bm.NewClient(c.HTTPClient.Read),
		httpW:               bm.NewClient(c.HTTPClient.Write),
		httpGame:            bm.NewClient(c.HTTPClient.Game),
		bangumiURL:          c.Host.Bangumi + _bangumiURI,
		bangumiConcernURL:   c.Host.Bangumi + _bangumiConcernURI,
		bangumiUnConcernURL: c.Host.Bangumi + _bangumiUnConcernURI,
		favFolderURL:        c.Host.API + _favFolderURI,
		favArcURL:           c.Host.API + _favArchiveURI,
		favAlbumURL:         c.Host.APILive + _favAlbumURI,
		favMovieURL:         c.Host.Bangumi + _favMovieURI,
		shopURL:             c.Host.Mall + _shopURI,
		shopLinkURL:         c.Host.Mall + _shopLinkURI,
		albumCountURL:       c.Host.APIVc + _albumCountURI,
		albumListURL:        c.Host.APIVc + _albumListURI,
		tagSubURL:           c.Host.API + _tagSubURI,
		tagCancelSubURL:     c.Host.API + _tagCancelSubURI,
		tagSubListURL:       c.Host.API + _subTagListURI,
		accTagsURL:          c.Host.Acc + _accTagsURI,
		accTagsSetURL:       c.Host.Acc + _accTagsSetURI,
		isAnsweredURL:       c.Host.API + _isAnsweredURI,
		lastPlayGameURL:     c.Host.Game + _lastPlayGameURI,
		appPlayedGameURL:    c.Host.AppGame + _appPlayedGameURI,
		arcSearchURL:        c.Host.Search + _arcSearchURI,
		webTopPhotoURL:      c.Host.Space + _webTopPhotoURI,
		topPhotoURL:         c.Host.Space + _topPhotoURI,
		liveMetalURL:        c.Host.APILive + _liveMetalURI,
		liveURL:             c.Host.APILive + _liveURI,
		medalStatusURL:      c.Host.APILive + _medalStatusURI,
		groupsCountURL:      c.Host.APIVc + _groupsCountURI,
		elecURL:             c.Host.Elec + _elecURI,
		audioCardURL:        c.Host.API + _audioCardURI,
		audioUpperCertURL:   c.Host.API + _audioUpperCertURI,
		audioCntURL:         c.Host.API + _audioCntURI,
		dynamicListURL:      c.Host.APIVc + _dynamicListURI,
		dynamicURL:          c.Host.APIVc + _dynamicURI,
		dynamicCntURL:       c.Host.APIVc + _dynamicCntURI,
		// expire
		clExpire:        int32(time.Duration(c.Redis.ClExpire) / time.Second),
		upArtExpire:     int32(time.Duration(c.Redis.UpArtExpire) / time.Second),
		upArcExpire:     int32(time.Duration(c.Redis.UpArcExpire) / time.Second),
		mcSettingExpire: int32(time.Duration(c.Memcache.SettingExpire) / time.Second),
		mcNoticeExpire:  int32(time.Duration(c.Memcache.NoticeExpire) / time.Second),
		mcTopArcExpire:  int32(time.Duration(c.Memcache.TopArcExpire) / time.Second),
		mcMpExpire:      int32(time.Duration(c.Memcache.MpExpire) / time.Second),
		mcThemeExpire:   int32(time.Duration(c.Memcache.ThemeExpire) / time.Second),
		mcTopDyExpire:   int32(time.Duration(c.Memcache.TopDyExpire) / time.Second),
		// cache
		cache: fanout.New("cache"),
	}
	d.channelStmt = make([]*sql.Stmt, _chSub)
	d.channelListStmt = make([]*sql.Stmt, _chSub)
	d.channelCntStmt = make([]*sql.Stmt, _chSub)
	d.channelArcCntStmt = make([]*sql.Stmt, _chSub)
	for i := 0; i < _chSub; i++ {
		d.channelStmt[i] = d.db.Prepared(fmt.Sprintf(_chSQL, i))
		d.channelListStmt[i] = d.db.Prepared(fmt.Sprintf(_chListSQL, i))
		d.channelCntStmt[i] = d.db.Prepared(fmt.Sprintf(_chCntSQL, i))
		d.channelArcCntStmt[i] = d.db.Prepared(fmt.Sprintf(_chArcCntSQL, i))
	}
	return
}

// Ping ping dao
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		return
	}
	err = d.pingRedis(c)
	return
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
