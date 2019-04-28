package zlimit

import (
	"context"
	"time"

	"go-common/app/service/main/location/model"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xip "go-common/library/net/ip"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_prefixBlackList = "zl_"
)

// Config default struct
type Config struct {
	DB        *sql.Config
	Redis     *Redis
	IPFile    string
	FlushTime xtime.Duration
}

// Redis redis struct
type Redis struct {
	*redis.Config
	Expire xtime.Duration
}

// Service zlimit service struct
type Service struct {
	// mysql
	db                 *sql.DB
	getPolicyStmt      *sql.Stmt
	getRelationStmt    *sql.Stmt
	getGroupPolicyStmt *sql.Stmt
	// redis
	redis     *redis.Pool
	expire    int32
	flushTime time.Duration
	// cache
	policy      map[int64]map[int64]int64
	groupPolicy map[int64][]int64
	missch      chan interface{}
	// xip
	list *xip.List
}

// New new zlimit service
func New(c *Config) (s *Service) {
	var err error
	s = &Service{
		db:          sql.NewMySQL(c.DB),
		redis:       redis.NewPool(c.Redis.Config),
		expire:      int32(time.Duration(c.Redis.Expire) / time.Second),
		missch:      make(chan interface{}, 1024),
		policy:      make(map[int64]map[int64]int64),
		groupPolicy: make(map[int64][]int64),
		flushTime:   time.Duration(c.FlushTime),
	}
	s.getPolicyStmt = s.db.Prepared(_getPolicySQL)
	s.getRelationStmt = s.db.Prepared(_getRelationSQL)
	s.getGroupPolicyStmt = s.db.Prepared(_getGolbalPolicySQL)
	s.load()
	s.list, err = xip.New(c.IPFile)
	if err != nil {
		log.Error("xip.New(%s) error(%v)", c.IPFile, err)
		panic(err)
	}
	go s.reloadproc()
	go s.cacheproc()
	return
}

func (s *Service) load() {
	var (
		tmpPolicy      map[int64]map[int64]int64
		tmpGroupPolicy map[int64][]int64
		err            error
	)
	if tmpPolicy, err = s.policies(context.TODO()); err != nil {
		log.Error("s.policies error(%v)", err)
	} else if len(tmpPolicy) > 0 {
		s.policy = tmpPolicy
	}
	if tmpGroupPolicy, err = s.groupPolicies(context.TODO()); err != nil {
		log.Error("s.groupPolicies error(%v)", err)
	} else if len(tmpGroupPolicy) > 0 {
		s.groupPolicy = tmpGroupPolicy
	}
}

// reloadproc reload data from db
func (s *Service) reloadproc() {
	for {
		s.load()
		time.Sleep(s.flushTime)
	}
}

// Find redio rule by aid and ipaddr
func (s *Service) Find(c context.Context, aid int64, ipaddr, cdnip string) (ret, retdown int64, err error) {
	var (
		ok                  bool
		auth, pid, zid, gid int64
		rules, pids         []int64
		zids                map[int64]int64
		ipInfo              *xip.Zone
	)
	ipInfo = s.list.Zone(ipaddr)
	if (ipInfo != nil) && (ipInfo.Province == "共享地址" || ipInfo.City == "共享地址") && cdnip != "" {
		ipInfo = s.list.Zone(cdnip)
	}
	if ipInfo == nil {
		ret = model.Allow
		retdown = model.AllowDown
		return
	}
	uz := s.zoneids(ipInfo) // country, state, city
	if ok, err = s.existsRule(c, aid); err != nil {
		log.Error("s.existsRule error(%v)", err)
		err = nil
	} else if ok {
		if rules, err = s.rule(c, aid, uz); err != nil {
			log.Error("s.rule(%d) error(%v) ", aid, err)
			err = nil
		} else {
			for _, auth = range rules {
				retdown = 0xff & auth
				ret = auth >> 8
				if ret != 0 {
					break
				}
			}
			if ret == 0 {
				ret = model.Allow
				retdown = model.AllowDown
			}
			return
		}
	}
	if gid, err = s.groupid(c, aid); err != nil {
		return
	} else if gid != 0 {
		if pids, ok = s.groupPolicy[gid]; ok {
			for _, pid = range pids {
				if zids, ok = s.policy[pid]; !ok {
					continue
				}
				if ret == 0 {
					//  ret already set skip check
					for _, zid = range uz {
						if auth, ok = zids[zid]; ok {
							if ret == 0 {
								retdown = 0xff & auth
								ret = auth >> 8 // ret must not be zero
								break
							}
						}
					}
				}
				tmpZids := map[int64]map[int64]int64{
					aid: zids,
				}
				s.addCache(tmpZids)
			}
			if ret == 0 {
				ret = model.Allow
				retdown = model.AllowDown
			}
			return
		}
	}
	ret = model.Allow
	retdown = model.AllowDown
	zids = make(map[int64]int64)
	zids[0] = ret<<8 | retdown
	tmpZids := map[int64]map[int64]int64{
		aid: zids,
	}
	s.addCache(tmpZids)
	return
}

// Forbid check ip is forbid or not.
func (s *Service) Forbid(c context.Context, pstr string, ipaddr string) (err error) {
	if pstr == "" {
		return
	}
	var (
		ret  int64
		pids []int64
	)
	if pids, err = xstr.SplitInts(pstr); err != nil {
		log.Error("xstr.SplitInts(%s) error(%v)", pstr, err)
		return
	}
	if ret, _ = s.FindByPid(c, pids, ipaddr); ret == model.Forbidden {
		err = ecode.ZlimitForbidden
	}
	return
}

// FindByPid redio rule by policy_id and ipaddr
func (s *Service) FindByPid(c context.Context, pids []int64, ipaddr string) (ret, retdown int64) {
	var (
		ok      bool
		auth    int64
		zoneids []int64
	)
	ret = model.Allow
	retdown = model.AllowDown
	ipInfo := s.list.Zone(ipaddr)
	if ipInfo == nil {
		return
	}
	zoneids = s.zoneids(ipInfo)
	for _, pid := range pids {
		if _, ok = s.policy[pid]; !ok {
			continue
		}
		for _, zoneid := range zoneids {
			if auth, ok = s.policy[pid][zoneid]; ok {
				retdown = 0xff & auth
				ret = auth >> 8
				break
			}
		}
	}
	return
}

// FindByGid redio rule by group_id and ipaddr(or cdnip)
func (s *Service) FindByGid(c context.Context, gid int64, ipaddr, cdnip string) (ret, retdown int64) {
	var ipInfo *xip.Zone
	ret = model.Allow
	retdown = model.AllowDown
	ipInfo = s.list.Zone(ipaddr)
	if (ipInfo != nil) && (ipInfo.Province == "共享地址" || ipInfo.City == "共享地址") && cdnip != "" {
		ipInfo = s.list.Zone(cdnip)
	}
	if ipInfo == nil {
		return
	}
	zoneids := s.zoneids(ipInfo)
	if pids, ok := s.groupPolicy[gid]; ok {
		for _, pid := range pids {
			if _, ok := s.policy[pid]; !ok {
				continue
			}
			for _, zoneid := range zoneids {
				if auth, ok := s.policy[pid][zoneid]; ok {
					retdown = 0xff & auth
					ret = auth >> 8
					break
				}
			}
		}
	}
	return
}

// zoneids make zoneids
func (s *Service) zoneids(ipinfos *xip.Zone) []int64 {
	cZid := xip.ZoneID(ipinfos.Country, "", "")
	cpZid := xip.ZoneID(ipinfos.Country, ipinfos.Province, "")
	cpcZid := xip.ZoneID(ipinfos.Country, ipinfos.Province, ipinfos.City)
	zoneids := []int64{0, cZid, cpZid, cpcZid}
	return zoneids
}

func (s *Service) addCache(d interface{}) {
	// asynchronous add rules to redis
	select {
	case s.missch <- d:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for add rules into redis.
func (s *Service) cacheproc() {
	for {
		d := <-s.missch
		switch d.(type) {
		case map[int64]map[int64]int64:
			v := d.(map[int64]map[int64]int64)
			if err := s.addRule(context.TODO(), v); err != nil {
				log.Error("s.addRule error(%v) error(%v)", v, err)
			}
		default:
			log.Warn("cacheproc can't process the type")
		}
	}
}
