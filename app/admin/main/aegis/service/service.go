package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/aegis/conf"
	daoOrm "go-common/app/admin/main/aegis/dao/gorm"
	daoHttp "go-common/app/admin/main/aegis/dao/http"
	daoMc "go-common/app/admin/main/aegis/dao/mc"
	daoMysql "go-common/app/admin/main/aegis/dao/mysql"
	daoRedis "go-common/app/admin/main/aegis/dao/redis"
	daoRpc "go-common/app/admin/main/aegis/dao/rpc"
	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/middleware"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"
)

// Service struct
type Service struct {
	c *conf.Config

	gorm     *daoOrm.Dao
	http     *daoHttp.Dao
	mc       *daoMc.Dao
	mysql    *daoMysql.Dao
	redis    *daoRedis.Dao
	rpc      *daoRpc.Dao
	aegisPub *databus.Databus

	// cache
	groupCache         map[int64]*common.Group
	bizCfgCache        map[string]string
	bizRoleCache       map[int64]map[string]int64  //business:cfg_name:role_id/manager_id映射
	taskRoleCache      map[int64]map[int64][]int64 //business:manager_id:flow_id映射
	netCacheCh         chan map[string]string
	bizMiddlewareCache map[int64][]*middleware.Aggregate

	netCache        map[int64]*net.Net
	tokenCache      map[int64]*net.Token
	flowCache       map[int64]*net.Flow
	transitionCache map[int64]*net.Transition
	bindCache       map[int64]*net.TokenBind

	//async
	async *fanout.Fanout

	//gray
	gray map[int64][][]common.GrayField
}

// Cache .
func (s *Service) Cache() map[string]interface{} {
	return map[string]interface{}{
		"groupCache":         s.groupCache,
		"bizCfgCache":        s.bizCfgCache,
		"bizRoleCache":       s.bizRoleCache,
		"taskRoleCache":      s.taskRoleCache,
		"bizMiddlewareCache": s.bizMiddlewareCache,
	}
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		// dao
		gorm:     daoOrm.New(c),
		http:     daoHttp.New(c),
		mc:       daoMc.New(c),
		mysql:    daoMysql.New(c),
		redis:    daoRedis.New(c),
		rpc:      daoRpc.New(c),
		aegisPub: databus.New(c.AegisPub),

		// cache
		netCache:        map[int64]*net.Net{},
		tokenCache:      map[int64]*net.Token{},
		flowCache:       map[int64]*net.Flow{},
		transitionCache: map[int64]*net.Transition{},
		bindCache:       map[int64]*net.TokenBind{},

		//gray
		gray:       loadGray(c),
		netCacheCh: make(chan map[string]string, 10240),

		//async
		async: fanout.New("async", fanout.Worker(10), fanout.Buffer(10240)),
	}

	go s.cacheProc()
	go s.setNetCache()
	return s
}

// Debug .
func (s *Service) Debug() string {
	return s.c.Debug
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.mysql.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.async.Close()
	s.gorm.Close()
	s.mysql.Close()
	s.mc.Close()
	s.redis.Close()
	s.aegisPub.Close()
}

func (s *Service) cacheProc() {
	for {
		s.syncUpCache(context.Background())
		s.syncBizCache(context.Background())

		time.Sleep(10 * time.Minute)
	}
}

//IsAdmin uid是否为管理员
func (s *Service) IsAdmin(uid int64) bool {
	return len(s.c.Admin) > 0 && strings.Contains(","+s.c.Admin+",", fmt.Sprintf(",%d,", uid))
}

func loadGray(c *conf.Config) (gray map[int64][][]common.GrayField) {
	if c.Gray == nil {
		gray = make(map[int64][][]common.GrayField)
		return
	}

	gray = make(map[int64][][]common.GrayField)
	for _, biz := range c.Gray.Biz {
		ones := [][]common.GrayField{}
		for _, opts := range biz.Options {
			options := []common.GrayField{}
			for _, field := range opts.Fields {
				field.Name = strings.TrimSpace(field.Name)
				field.Value = strings.TrimSpace(field.Value)
				if field.Value == "" || field.Name == "" {
					continue
				}
				options = append(options, common.GrayField{
					Name:  field.Name,
					Value: fmt.Sprintf(",%s,", field.Value),
				})
			}
			ones = append(ones, options)
		}
		gray[biz.BusinessID] = ones
	}
	return
}

//GetMiddlewareCache get cache by bizid
func (s *Service) GetMiddlewareCache(bizid int64) []*middleware.Aggregate {
	return s.bizMiddlewareCache[bizid]
}
