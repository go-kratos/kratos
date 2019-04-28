package archive

import (
	"context"
	"time"

	"fmt"
	"go-common/app/interface/main/creative/conf"
	arcMdl "go-common/app/interface/main/creative/model/archive"
	archive "go-common/app/service/main/archive/api/gorpc"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"
)

// Dao is archive dao.
type Dao struct {
	// config
	c *conf.Config
	// rpc
	arc *archive.Service2
	// select
	client *bm.Client
	// db
	db *sql.DB
	// mc
	mc       *memcache.Pool
	mcExpire int32
	//cache tool
	cache *fanout.Fanout
	// redis
	redis       *redis.Pool
	redisExpire int32
	// uri
	view          string
	views         string
	del           string
	video         string
	hList         string
	hView         string
	flow          string
	upArchives    string
	descFormat    string
	nsMd5         string
	simpleVideos  string
	simpleArchive string
	videoJam      string
	upSpecialURL  string
	flowJudge     string
	staffApplies  string
	staffApply    string
	staffCheck    string
}

// New init api url
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:   c,
		arc: archive.New2(c.ArchiveRPC),
		// http client
		client: bm.NewClient(c.HTTPClient.Normal),
		// db
		db:    sql.NewMySQL(c.DB.Archive),
		cache: fanout.New("dao_archive", fanout.Worker(5), fanout.Buffer(10240)),
		// mc
		mc:       memcache.NewPool(c.Memcache.Archive.Config),
		mcExpire: 600,
		//Fav redis cache
		redis:       redis.NewPool(c.Redis.Cover.Config),
		redisExpire: int32(time.Duration(c.Redis.Cover.Expire) / time.Second),
		// uri
		view:          c.Host.Videoup + _view,
		views:         c.Host.Videoup + _views,
		video:         c.Host.Videoup + _video,
		del:           c.Host.Videoup + _del,
		hList:         c.Host.Videoup + _hList,
		hView:         c.Host.Videoup + _hView,
		flow:          c.Host.Videoup + _flow,
		upArchives:    c.Host.Videoup + _archives,
		descFormat:    c.Host.Videoup + _descFormat,
		nsMd5:         c.Host.Videoup + _nsMd5,
		simpleVideos:  c.Host.Videoup + _simpleVideos,
		simpleArchive: c.Host.Videoup + _simpleArchive,
		videoJam:      c.Host.Videoup + _videoJam,
		upSpecialURL:  c.Host.API + _upSpecial,
		flowJudge:     c.Host.Videoup + _flowjudge,
		staffApplies:  c.Host.Videoup + _staffApplies,
		staffApply:    c.Host.Videoup + _staffApply,
		staffCheck:    c.Host.Videoup + _staffCheck,
	}
	return
}

// Ping fn
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		return
	}
	return
}

// Close fn
func (d *Dao) Close() (err error) {
	if d.redis != nil {
		d.redis.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
	return d.db.Close()
}

func staffKey(aid int64) string {
	return fmt.Sprintf("staff_aid_%d", aid)
}

func (d *Dao) cacheSFStaffData(aid int64) string {
	return fmt.Sprintf("staff_aid_sf_%d", aid)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -singleflight=true -nullcache=[]*arcMdl.Staff{{ID:-1}} -check_null_code=len($)==1&&$[0].ID==-1
	StaffData(c context.Context, aid int64) ([]*arcMdl.Staff, error)
	ViewPoint(c context.Context, aid int64, cid int64) (vp *arcMdl.ViewPointRow, err error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// mc: -key=staffKey
	CacheStaffData(c context.Context, key int64) ([]*arcMdl.Staff, error)
	// 这里也支持自定义注释 会替换默认的注释
	// mc: -key=staffKey -expire=3 -encode=json|gzip
	AddCacheStaffData(c context.Context, key int64, value []*arcMdl.Staff) error
	// mc: -key=staffKey
	DelCacheStaffData(c context.Context, key int64) error
	//mc: -key=viewPointCacheKey -expire=_viewPointExp -encode=json
	AddCacheViewPoint(c context.Context, aid int64, vp *arcMdl.ViewPointRow, cid int64) (err error)
	//mc: -key=viewPointCacheKey
	CacheViewPoint(c context.Context, aid int64, cid int64) (vp *arcMdl.ViewPointRow, err error)
}
