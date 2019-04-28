package dao

import (
	"context"
	"time"

	"go-common/app/service/main/push/conf"
	"go-common/app/service/main/push/dao/apns2"
	"go-common/app/service/main/push/dao/fcm"
	"go-common/app/service/main/push/dao/huawei"
	"go-common/app/service/main/push/dao/jpush"
	"go-common/app/service/main/push/dao/mi"
	"go-common/app/service/main/push/dao/oppo"
	"go-common/app/service/main/push/model"
	"go-common/library/cache/memcache"
	xredis "go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
)

const (
	_retry = 3
)

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	//mc: -key=tokenKey -type=get
	TokenCache(c context.Context, key string) (*model.Report, error)

	//mc: -key=tokenKey -expire=d.mcReportExpire
	AddTokenCache(c context.Context, key string, value *model.Report) error
	//mc: -key=tokenKey -expire=d.mcReportExpire
	AddTokensCache(c context.Context, values map[string]*model.Report) error

	//mc: -key=tokenKey
	DelTokenCache(c context.Context, key string) error
}

// Dao .
type Dao struct {
	c                      *conf.Config
	db                     *xsql.DB
	mc                     *memcache.Pool
	redis                  *xredis.Pool
	reportPub              *databus.Databus
	callbackPub            *databus.Databus
	clientsIPhone          map[int64][]*apns2.Client
	clientsIPad            map[int64][]*apns2.Client
	clientsMi              map[int64][]*mi.Client
	clientMiByMids         map[int64]*mi.Client
	clientsHuawei          map[int64][]*huawei.Client
	clientsOppo            map[int64][]*oppo.Client
	clientsJpush           map[int64][]*jpush.Client
	clientsFCM             map[int64][]*fcm.Client
	clientsLen             map[string]int
	clientsIndex           map[string]*uint32
	huaweiAuth             map[int64]*huawei.Access
	oppoAuth               map[int64]*oppo.Auth
	addTaskStmt            *xsql.Stmt
	updateTaskStatusStmt   *xsql.Stmt
	updateTaskProgressStmt *xsql.Stmt
	taskStmt               *xsql.Stmt
	businessesStmt         *xsql.Stmt
	settingStmt            *xsql.Stmt
	setSettingStmt         *xsql.Stmt
	authsStmt              *xsql.Stmt
	addReportStmt          *xsql.Stmt
	updateReportStmt       *xsql.Stmt
	reportStmt             *xsql.Stmt
	reportByIDStmt         *xsql.Stmt
	delReportStmt          *xsql.Stmt
	reportsByMidStmt       *xsql.Stmt
	lastReportIDStmt       *xsql.Stmt
	addCallbackStmt        *xsql.Stmt
	redisTokenExpire       int32
	redisLaterExpire       int32
	redisMidsExpire        int32
	mcReportExpire         int32
	mcSettingExpire        int32
	mcUUIDExpire           int32
}

var (
	errorsCount = prom.BusinessErrCount
	infosCount  = prom.BusinessInfoCount
	missedCount = prom.CacheMiss
	cachedCount = prom.CacheHit
)

// New creates a push-service DAO instance.
func New(c *conf.Config) *Dao {
	d := &Dao{
		c:                c,
		db:               xsql.NewMySQL(c.MySQL),
		mc:               memcache.NewPool(c.Memcache.Config),
		redis:            xredis.NewPool(c.Redis.Config),
		reportPub:        databus.New(c.ReportPub),
		callbackPub:      databus.New(c.CallbackPub),
		clientsIPhone:    make(map[int64][]*apns2.Client),
		clientsIPad:      make(map[int64][]*apns2.Client),
		clientsMi:        make(map[int64][]*mi.Client),
		clientMiByMids:   make(map[int64]*mi.Client),
		clientsHuawei:    make(map[int64][]*huawei.Client),
		clientsOppo:      make(map[int64][]*oppo.Client),
		clientsJpush:     make(map[int64][]*jpush.Client),
		clientsFCM:       make(map[int64][]*fcm.Client),
		clientsLen:       make(map[string]int),
		clientsIndex:     make(map[string]*uint32),
		huaweiAuth:       make(map[int64]*huawei.Access),
		oppoAuth:         make(map[int64]*oppo.Auth),
		redisTokenExpire: int32(time.Duration(c.Redis.TokenExpire) / time.Second),
		redisLaterExpire: int32(time.Duration(c.Redis.LaterExpire) / time.Second),
		redisMidsExpire:  int32(time.Duration(c.Redis.MidsExpire) / time.Second),
		mcReportExpire:   int32(time.Duration(c.Memcache.ReportExpire) / time.Second),
		mcSettingExpire:  int32(time.Duration(c.Memcache.SettingExpire) / time.Second),
		mcUUIDExpire:     int32(time.Duration(c.Memcache.UUIDExpire) / time.Second),
	}
	d.addTaskStmt = d.db.Prepared(_addTaskSQL)
	d.updateTaskStatusStmt = d.db.Prepared(_upadteTaskStatusSQL)
	d.updateTaskProgressStmt = d.db.Prepared(_upadteTaskProgressSQL)
	d.taskStmt = d.db.Prepared(_taskByIDSQL)
	d.businessesStmt = d.db.Prepared(_businessesSQL)
	d.settingStmt = d.db.Prepared(_settingSQL)
	d.setSettingStmt = d.db.Prepared(_setSettingSQL)
	d.authsStmt = d.db.Prepared(_authsSQL)
	d.addReportStmt = d.db.Prepared(_addReportSQL)
	d.updateReportStmt = d.db.Prepared(_updateReportSQL)
	d.addCallbackStmt = d.db.Prepared(_addCallbackSQL)
	d.reportStmt = d.db.Prepared(_reportSQL)
	d.reportByIDStmt = d.db.Prepared(_reportByIDSQL)
	d.delReportStmt = d.db.Prepared(_delReportSQL)
	d.reportsByMidStmt = d.db.Prepared(_reportsByMidSQL)
	d.lastReportIDStmt = d.db.Prepared(_lastReportIDSQL)

	go d.refreshAuthproc()
	time.Sleep(time.Second)
	d.loadClients()
	return d
}

func (d *Dao) refreshAuthproc() {
	for {
		auths, err := d.auths(context.Background())
		if err != nil {
			log.Error("d.auths() error(%v)", err)
			time.Sleep(time.Second)
			continue
		}
		for _, a := range auths {
			d.refreshAuth(a)
		}
		time.Sleep(1 * time.Minute)
	}
}

func (d *Dao) refreshAuth(a *model.Auth) {
	i := fmtRoundIndex(a.APPID, a.PlatformID)
	switch a.PlatformID {
	case model.PlatformOppo:
		if d.clientsLen[i] == 0 || d.oppoAuth[a.APPID] == nil || d.oppoAuth[a.APPID].IsExpired() {
			auth, err := oppo.NewAuth(a.Key, a.Value)
			if err != nil {
				log.Error("new oppo auth failed, key(%s) secret(%s) error(%v)", a.Key, a.Value, err)
				return
			}
			log.Info("oppo refresh auth app(%d) auth(%+v)", a.APPID, auth)
			if d.oppoAuth[a.APPID] == nil {
				d.oppoAuth[a.APPID] = new(oppo.Auth)
			}
			*d.oppoAuth[a.APPID] = *auth
			if d.clientsLen[i] == 0 {
				cs := d.newOppoClients(a.APPID, a.BundleID)
				if len(cs) > 0 {
					d.clientsOppo[a.APPID] = cs
					d.clientsLen[i] = len(d.clientsOppo)
					log.Info("oppo renew push clients app(%d)", a.APPID)
				}
			}
		}
	case model.PlatformHuawei:
		if d.clientsLen[i] == 0 || d.huaweiAuth[a.APPID] == nil || d.huaweiAuth[a.APPID].IsExpired() {
			ac, err := huawei.NewAccess(a.Key, a.Value)
			if err != nil {
				log.Error("new huawei access failed, id(%s) secret(%s) error(%v)", a.Key, a.Value, err)
				return
			}
			log.Info("huawei refresh auth app(%d) auth(%+v)", a.APPID, ac)
			if d.huaweiAuth[a.APPID] == nil {
				d.huaweiAuth[a.APPID] = new(huawei.Access)
			}
			*d.huaweiAuth[a.APPID] = *ac
			if d.clientsLen[i] == 0 {
				cs := d.newHuaweiClients(a.APPID, a.BundleID)
				if len(cs) > 0 {
					d.clientsHuawei[a.APPID] = cs
					d.clientsLen[i] = len(d.clientsHuawei)
					log.Info("huawei renew push clients app(%d)", a.APPID)
				}
			}
		}
	}
}

// PromError prom error
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromInfo add prom info
func PromInfo(name string) {
	infosCount.Incr(name)
}

// PromChanLen channel length
func PromChanLen(name string, length int64) {
	infosCount.State(name, length)
}

// BeginTx begin transaction.
func (d *Dao) BeginTx(c context.Context) (*xsql.Tx, error) {
	return d.db.Begin(c)
}

// Close dao.
func (d *Dao) Close() (err error) {
	if err = d.db.Close(); err != nil {
		log.Error("d.db.Close() error(%v)", err)
		PromError("db:close")
	}
	if err = d.redis.Close(); err != nil {
		log.Error("d.redis.Close() error(%v)", err)
		PromError("redis:close")
	}
	if err = d.mc.Close(); err != nil {
		log.Error("d.mc.Close() error(%v)", err)
		PromError("mc:close")
	}
	return
}

// Ping check connection status.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		PromError("redis:Ping")
		log.Error("d.pingRedis error(%v)", err)
		return
	}
	if err = d.pingMC(c); err != nil {
		PromError("mc:Ping")
		log.Error("d.pingMC error(%v)", err)
		return
	}
	if err = d.db.Ping(c); err != nil {
		PromError("mysql:Ping")
		log.Error("d.db.Ping error(%v)", err)
	}
	return
}
