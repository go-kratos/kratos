package archive

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	ajbmdl "go-common/app/job/main/archive/model/databus"
	accapi "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/conf"
	"go-common/app/service/main/archive/model/archive"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"

	"go-common/library/sync/errgroup"
)

const (
	_multiInterval = 200
)

var (
	_emptyPages3 = []*api.Page{}
)

// Dao is archive dao.
type Dao struct {
	c *conf.Config
	// db
	db        *sql.DB
	arcReadDB *sql.DB
	resultDB  *sql.DB
	statDB    *sql.DB
	clickDB   *sql.DB
	// acc rpc
	acc accapi.AccountClient
	// memcache
	mc *memcache.Pool
	// redis
	upRds    *redis.Pool
	upExpire int32
	// region
	rgRds *redis.Pool
	// archive stmt
	rgnArcsStmt      *sql.Stmt
	upCntStmt        *sql.Stmt
	upPasStmt        *sql.Stmt
	reportResultStmt *sql.Stmt
	additStmt        *sql.Stmt
	// video stmt
	vdosStmt *sql.Stmt
	// dede stmt
	tpsStmt *sql.Stmt
	// type cache
	tNamem map[int16]string
	// report_result cache
	aidResult map[int64]string
	// cache chan
	cacheCh  chan func()
	hitProm  *prom.Prom
	missProm *prom.Prom
	errProm  *prom.Prom
	infoProm *prom.Prom
	// player http client
	playerClient *bm.Client
	// player qn map
	playerQn     map[int]struct{}
	playerVipQn  map[int]struct{}
	cacheDatabus *databus.Databus
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// db
		db:        sql.NewMySQL(c.DB.Arc),
		arcReadDB: sql.NewMySQL(c.DB.ArcRead),
		resultDB:  sql.NewMySQL(c.DB.ArcResult),
		statDB:    sql.NewMySQL(c.DB.Stat),
		clickDB:   sql.NewMySQL(c.DB.Click),
		// memcache
		mc: memcache.NewPool(c.Memcache.Archive.Config),
		// redis
		upRds:    redis.NewPool(c.Redis.Archive.Config),
		upExpire: 480 * 60 * 60,
		rgRds:    redis.NewPool(c.Redis.Archive.Config),
		// cache chan
		cacheCh:      make(chan func(), 1024),
		hitProm:      prom.CacheHit,
		missProm:     prom.CacheMiss,
		errProm:      prom.BusinessErrCount,
		infoProm:     prom.BusinessInfoCount,
		playerClient: bm.NewClient(c.PlayerClient),
		playerQn:     make(map[int]struct{}),
		playerVipQn:  make(map[int]struct{}),
		cacheDatabus: databus.New(c.CacheDatabus),
	}
	var err error
	if d.acc, err = accapi.NewClient(c.AccClient); err != nil {
		panic(fmt.Sprintf("account GRPC error(%v)!!!!!!!!!!!!!!!!!!!!!!", err))
	}
	d.rgnArcsStmt = d.resultDB.Prepared(_rgnArcsSQL)
	d.upCntStmt = d.resultDB.Prepared(_upCntSQL)
	d.upPasStmt = d.resultDB.Prepared(_upPasSQL)
	d.tpsStmt = d.arcReadDB.Prepared(_tpsSQL)
	d.additStmt = d.arcReadDB.Prepared(_additSQL)
	// video stmt
	d.vdosStmt = d.resultDB.Prepared(_vdosSQL)
	// archive_report_resutl
	d.reportResultStmt = d.arcReadDB.Prepared(_reportResultSQL)
	// type cache
	for _, pn := range d.c.PlayerQn {
		d.playerQn[pn] = struct{}{}
	}
	for _, pvn := range d.c.PlayerVipQn {
		d.playerVipQn[pvn] = struct{}{}
	}
	d.loadTypes()
	d.loadReportResult()
	for i := 0; i < runtime.NumCPU(); i++ {
		go d.cacheproc()
	}
	go d.loadproc()
	return
}

func (d *Dao) addCache(f func()) {
	select {
	case d.cacheCh <- f:
	default:
		log.Warn("d.cacheCh is full")
	}
}

// Archive3 get archive by aid.
func (d *Dao) Archive3(c context.Context, aid int64) (a *api.Arc, err error) {
	var cached = true
	if a, err = d.archive3Cache(c, aid); err != nil {
		log.Error("d.archivePBCache(%d) error(%v)", aid, err)
		err = nil // NOTE ignore error use db
		cached = false
	}
	if a != nil {
		if st, _ := d.Stat3(c, aid); st != nil {
			a.Stat = *st
			a.FillStat()
		}
		a.ReportResult = d.aidResult[aid]
		return
	}
	if a, err = d.archive3(c, aid); err != nil {
		log.Error("d.archive(%d) error(%v)", aid, err)
		return
	}
	if a == nil {
		err = ecode.NothingFound
		return
	}
	if st, _ := d.Stat3(c, aid); st != nil {
		a.Stat = *st
		a.FillStat()
	}
	d.fillArchive3(c, a, cached)
	a.ReportResult = d.aidResult[aid]
	return
}

// Archives3 get archives by aids.
func (d *Dao) Archives3(c context.Context, aids []int64) (am map[int64]*api.Arc, err error) {
	var eg errgroup.Group
	eg.Go(func() (err error) {
		var (
			missed []int64
			missm  map[int64]*api.Arc
			cached = true
		)
		if am, err = d.archive3Caches(c, aids); err != nil {
			log.Error("%+v", err)
			am = make(map[int64]*api.Arc, len(aids))
			err = nil // NOTE: ignore error
			cached = false
		}
		for _, aid := range aids {
			if _, ok := am[aid]; !ok {
				missed = append(missed, aid)
			}
		}
		if len(missed) == 0 {
			return
		}
		if missm, err = d.archives3(c, missed); err != nil {
			log.Error("d.archives(%v) error(%v)", missed, err)
			return
		}
		if len(missm) == 0 {
			log.Warn("archives miss(%+v)", missed)
			return
		}
		d.fillArchives(c, missm, cached)
		for _, a := range missm {
			am[a.Aid] = a
		}
		return
	})
	var stm map[int64]*api.Stat
	eg.Go(func() (err error) {
		var (
			missed []int64
			missm  map[int64]*api.Stat
			cached = true
		)
		if stm, missed, err = d.statCaches3(c, aids); err != nil {
			log.Error("d.statCaches(%d) error(%v)", aids, err)
			missed = aids
			stm = make(map[int64]*api.Stat, len(aids))
			err = nil // NOTE: ignore error
			cached = false
		}
		d.hitProm.Add("stat3", int64(len(stm)))
		if stm != nil && len(missed) == 0 {
			return
		}
		if missm, err = d.stats3(c, missed); err != nil {
			log.Error("d.stats(%v) error(%v)", missed, err)
			err = nil // NOTE: ignore error
		}
		for aid, st := range missm {
			stm[aid] = st
			if cached {
				var cst = &api.Stat{}
				*cst = *st
				d.addCache(func() {
					d.addStatCache3(context.TODO(), cst)
				})
			}
		}
		d.missProm.Add("stat3", int64(len(missm)))
		return
	})
	if err = eg.Wait(); err != nil {
		log.Error("eg.Wait(%v) error(%v)", aids, err)
		return
	}
	for aid, a := range am {
		if st, ok := stm[aid]; ok {
			a.Stat = *st
			a.FillStat()
		}
		a.ReportResult = d.aidResult[aid]
	}
	return
}

// Videos3 get archive videos by aid.
func (d *Dao) Videos3(c context.Context, aid int64) (ps []*api.Page, err error) {
	var cached = true
	if ps, err = d.pageCache3(c, aid); err != nil {
		log.Error("d.pageCache3(%d) error(%v)", aid, err)
		err = nil // NOTE ignore error use db
		cached = false
	}
	if ps != nil {
		d.hitProm.Add("videos3", 1)
		return
	}
	if ps, err = d.videos3(c, aid); err != nil {
		log.Error("d.videos(%d) error(%v)", aid, err)
		return
	}
	if len(ps) == 0 {
		log.Warn("archive(%d) have not passed video", aid)
		ps = _emptyPages3
	}
	d.missProm.Add("videos3", 1)
	if cached {
		d.addCache(func() {
			d.addPageCache3(context.TODO(), aid, ps)
		})
	}
	return
}

// VideosByAids3 get videos by aids
func (d *Dao) VideosByAids3(c context.Context, aids []int64) (vs map[int64][]*api.Page, err error) {
	var (
		missed []int64
		cached = true
	)
	if vs, missed, err = d.pagesCache3(c, aids); err != nil {
		log.Error("d.pagesCache(%v) error(%v)", aids, err)
		missed = aids
		err = nil
		cached = false
	}
	d.hitProm.Add("videos3", int64(len(vs)))
	if len(missed) == 0 && vs != nil {
		return
	}
	var missVs = make(map[int64][]*api.Page, len(missed))
	if missVs, err = d.videosByAids3(c, missed); err != nil {
		log.Error("d.videosByAids3(%v) error(%v)", missed, err)
		return
	}
	d.missProm.Add("videos3", int64(len(missVs)))
	for aid, v := range missVs {
		vs[aid] = v
		if cached {
			var (
				caid = aid
				cv   = v
			)
			d.addCache(func() {
				d.addPageCache3(context.TODO(), caid, cv)
			})
		}
	}
	return
}

// Video3 get video by aid & cid.
func (d *Dao) Video3(c context.Context, aid, cid int64) (p *api.Page, err error) {
	var cached = true
	if p, err = d.videoCache3(c, aid, cid); err != nil {
		log.Error("d.videoCache3(%d, %d) error(%v)", aid, cid, err)
		err = nil // NOTE ignore error use db
		cached = false
	}
	if p != nil {
		d.hitProm.Add("video3", 1)
		return
	}
	if p, err = d.video3(c, aid, cid); err != nil {
		log.Error("d.video3(%d) error(%v)", aid, cid, err)
		return
	}
	if p == nil {
		log.Warn("archive(%d) cid(%d) no passed video", aid, cid)
		err = ecode.NothingFound
		cached = false
	}
	d.missProm.Add("video3", 1)
	if cached {
		d.addCache(func() {
			d.addVideoCache3(context.TODO(), aid, cid, p)
		})
	}
	return
}

// Description get Description from by aid || aid+cid.
func (d *Dao) Description(c context.Context, aid int64) (desc string, err error) {
	var (
		addit *archive.Addit
		a     *api.Arc
	)
	desc, _ = d.descCache(c, aid)
	if desc != "" {
		return
	}
	if addit, err = d.Addit(c, aid); err != nil {
		log.Error("d.Addit(%d) error(%v)", aid, err)
		err = nil
	}
	if addit != nil && addit.Description != "" {
		desc = addit.Description
		d.addCache(func() {
			d.addDescCache(context.TODO(), aid, desc)
		})
		return
	}
	if a, err = d.archive3(c, aid); err != nil {
		log.Error("d.archive(%d) error(%v)", aid, err)
		return
	}
	if a != nil && a.Desc != "" {
		desc = a.Desc
		d.addCache(func() {
			d.addDescCache(context.TODO(), aid, desc)
		})
	}
	return
}

// UpVideo3 update video by aid & cid.
func (d *Dao) UpVideo3(c context.Context, aid, cid int64) (err error) {
	var p *api.Page
	if p, err = d.video3(c, aid, cid); err != nil {
		log.Error("d.video2(%d) error(%v)", aid, cid, err)
		return
	}
	if p == nil {
		err = ecode.NothingFound
		return
	}
	d.addCache(func() {
		d.addVideoCache3(context.TODO(), aid, cid, p)
	})
	return
}

// UpperCache is
func (d *Dao) UpperCache(c context.Context, mid int64, action string, oldname string, oldface string) {
	var (
		aids      []int64
		err       error
		infoReply *accapi.InfoReply
	)
	if action == "updateByAdmin" {
		if infoReply, err = d.acc.Info3(c, &accapi.MidReq{Mid: mid}); err != nil || infoReply == nil {
			log.Error("d.acc.Info3(%d) error(%v)", mid, err)
			return
		}
		if infoReply.Info.Face == oldface && infoReply.Info.Name == oldname {
			log.Info("account updateByAdmin no need update")
			return
		}
	}
	if aids, _, _, err = d.UpperPassed(c, mid); err != nil {
		log.Error("d.UpperPassed(%d) error(%v)", mid, err)
		return
	}
	for _, aid := range aids {
		d.UpArchiveCache(context.TODO(), aid)
	}
}

// UpArchiveCache update archive cache by aid.
func (d *Dao) UpArchiveCache(c context.Context, aid int64) (err error) {
	var (
		a     *api.Arc
		addit *archive.Addit
		desc  string
	)
	if a, err = d.archive3(c, aid); a != nil && err == nil {
		d.fillArchive3(c, a, true)
		desc = a.Desc
	}
	if addit, err = d.Addit(c, aid); err == nil && addit != nil {
		if addit.Description != "" {
			desc = addit.Description
		}
	}
	d.addCache(func() {
		d.addDescCache(context.TODO(), aid, desc)
	})
	// pages3
	var ps3 []*api.Page
	if ps3, err = d.videos3(c, aid); err != nil {
		log.Error("d.videos3(%d) error(%v)", aid, err)
		return
	}
	d.addCache(func() {
		d.addPageCache3(context.TODO(), aid, ps3)
	})
	return
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		log.Error("d.db error(%v)", err)
		return
	}
	// mc
	mconn := d.mc.Get(c)
	defer mconn.Close()
	if err = mconn.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("mc.Store error(%v)", err)
		return
	}
	// upper redis
	rconn := d.upRds.Get(c)
	if _, err = rconn.Do("SET", "ping", "pong"); err != nil {
		rconn.Close()
		log.Error("rds.Set error(%v)", err)
		return
	}
	rconn.Close()
	// region redis
	rconn = d.rgRds.Get(c)
	if _, err = rconn.Do("SET", "ping", "pong"); err != nil {
		log.Error("rds.Set error(%v)", err)
	}
	rconn.Close()
	return
}

// Close close resource.
func (d *Dao) Close() {
	d.db.Close()
	d.mc.Close()
	d.upRds.Close()
	d.rgRds.Close()
}

func (d *Dao) fillArchives(c context.Context, am map[int64]*api.Arc, cached bool) {
	var (
		mids          []int64
		aids          []int64
		staffs        map[int64][]*api.StaffInfo
		cooperationOK error
		err           error
	)
	for _, a := range am {
		a.Fill()
		a.TypeName = d.tNamem[int16(a.TypeID)]
		if a.Author.Mid > 0 {
			mids = append(mids, a.Author.Mid)
		}
		if a.AttrVal(archive.AttrBitIsCooperation) == archive.AttrYes {
			aids = append(aids, a.Aid)
		}
	}
	mi, err := d.acc.Infos3(c, &accapi.MidsReq{Mids: mids})
	if err != nil || mi == nil {
		log.Error("d.acc.Infos(%v) error(%v)", mids, err)
		mi = new(accapi.InfosReply)
	}
	if len(aids) > 0 {
		staffs, cooperationOK = d.staffs(c, aids)
	}
	for _, a := range am {
		if m, ok := mi.Infos[a.Author.Mid]; ok {
			a.Author.Name = m.Name
			a.Author.Face = m.Face
		}
		if a.AttrVal(archive.AttrBitIsCooperation) == archive.AttrYes {
			if staff, ok := staffs[a.Aid]; ok {
				a.StaffInfo = staff
			}
		}
		if cached {
			var ca = &api.Arc{}
			*ca = *a
			d.addCache(func() {
				d.addArchive3Cache(context.TODO(), ca)
				if ca.Author.Name == "" || ca.Author.Face == "" || cooperationOK != nil {
					log.Error("account empty aid(%d) name(%s) face(%s)", ca.Aid, ca.Author.Name, ca.Author.Face)
					d.sendUpperCache(context.TODO(), ca.Aid)
				}
			})
		}
	}
}

func (d *Dao) fillArchive3(c context.Context, a *api.Arc, cached bool) {
	a.Fill()
	a.TypeName = d.tNamem[int16(a.TypeID)]
	if a.Author.Mid > 0 {
		reply, err := d.acc.Info3(c, &accapi.MidReq{Mid: a.Author.Mid})
		if err != nil || reply == nil {
			log.Error("d.acc.Info(%d) error(%v)", a.Author.Mid, err)
			err = nil
		} else {
			a.Author.Name = reply.Info.Name
			a.Author.Face = reply.Info.Face
		}
	}
	var (
		staffs        []*api.StaffInfo
		cooperationOK error
	)
	// cooperation
	if a.AttrVal(archive.AttrBitIsCooperation) == archive.AttrYes {
		if staffs, cooperationOK = d.staff(c, a.Aid); cooperationOK == nil {
			a.StaffInfo = staffs
		}
	}
	if cached {
		var ca = &api.Arc{}
		*ca = *a
		d.addCache(func() {
			d.addArchive3Cache(context.TODO(), ca)
			if ca.Author.Name == "" || ca.Author.Face == "" || cooperationOK != nil {
				log.Error("account empty aid(%d) name(%s) face(%s)", ca.Aid, ca.Author.Name, ca.Author.Face)
				d.sendUpperCache(context.TODO(), ca.Aid)
			}
		})
	}
}

func (d *Dao) sendUpperCache(c context.Context, aid int64) {
	for i := 0; i < 10; i++ {
		err := d.cacheDatabus.Send(context.Background(), strconv.FormatInt(aid, 10), &ajbmdl.Rebuild{Aid: aid})
		if err != nil {
			log.Error("d.cacheDatabus.Send(%d) error(%v)", aid, err)
			time.Sleep(50 * time.Millisecond)
			continue
		}
		log.Info("sendUpperCache(%d) ok", aid)
		break
	}
}

func (d *Dao) loadReportResult() {
	var (
		reportResults map[int64]*arcmdl.ReportResult
		err           error
		ar            = make(map[int64]string)
	)
	reportResults, err = d.ReportResults(context.TODO())
	if err != nil {
		log.Error("d.ReportResults error(%v)", err)
		return
	}
	for _, rr := range reportResults {
		ar[rr.Aid] = rr.Result
	}
	d.aidResult = ar
}

func (d *Dao) loadTypes() {
	var (
		types map[int16]*archive.ArcType
		nm    = make(map[int16]string)
		err   error
	)
	if types, err = d.Types(context.TODO()); err != nil {
		log.Error("d.Types error(%v)", err)
		return
	}
	for _, t := range types {
		nm[t.ID] = t.Name
	}
	d.tNamem = nm
}

func (d *Dao) cacheproc() {
	for {
		f, ok := <-d.cacheCh
		if !ok {
			return
		}
		f()
	}
}

func (d *Dao) loadproc() {
	for {
		time.Sleep(time.Duration(d.c.Tick))
		d.loadTypes()
		d.loadReportResult()
	}
}
