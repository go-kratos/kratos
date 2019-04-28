package operation

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/web-show/conf"
	"go-common/app/interface/main/web-show/dao/operation"
	opdml "go-common/app/interface/main/web-show/model/operation"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/library/log"
)

const (
	_rankCacheLen = 20
)

// Service struct
type Service struct {
	dao    *operation.Dao
	arcRPC *arcrpc.Service2
	cache  map[string]map[int][]*opdml.Operation
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		cache: make(map[string]map[int][]*opdml.Operation, len(opdml.Types)),
	}
	s.arcRPC = arcrpc.New2(c.RPCClient2.Archive)
	s.dao = operation.New(c)

	s.reload()
	go s.loadproc()
	return
}

// Notice return notice info
func (s *Service) operation(tp string, rank, num int) (res map[string][]*opdml.Operation) {
	res = make(map[string][]*opdml.Operation)
	tmp, ok := s.cache[tp]
	if ok {
		if rank != 0 {
			if ns := tmp[rank]; ns != nil {
				if len(ns) < num || num < 0 {
					num = len(ns)
				}
				ns = ns[:num]
				res[strconv.FormatInt(int64(rank), 10)] = ns
			}
		} else {
			for rk, ns := range tmp {
				if ns != nil {
					if len(ns) < num || num < 0 {
						num = len(ns)
					}
					ns = ns[:num]
					res[strconv.FormatInt(int64(rk), 10)] = ns
				}
			}
		}

	}
	return
}

// reload Service
func (s *Service) reload() {
	var (
		tmpT = make(map[string]map[int][]*opdml.Operation)
	)
	ops, err := s.dao.Operation(context.Background())
	if err != nil {
		log.Error("s.reloadNotice error(%v)", err)
		return
	}
	for _, op := range ops {
		tmp, ok := tmpT[op.Type]
		if !ok {
			tmp = make(map[int][]*opdml.Operation)
		}
		if len(tmp) > _rankCacheLen {
			continue
		}
		tmp[op.Rank] = append(tmp[op.Rank], op)
		tmpT[op.Type] = tmp
	}
	s.cache = tmpT
}

// loadproc Service
func (s *Service) loadproc() {
	for {
		s.reload()
		time.Sleep(time.Duration(conf.Conf.Reload.Notice))
	}
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}
