package block

import (
	"go-common/app/service/main/member/conf"
	"go-common/app/service/main/member/dao/block"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"

	"runtime/debug"
)

// Service is
type Service struct {
	c        *conf.Config
	dao      *block.Dao
	whiteMap map[int64]struct{}
	cache    *fanout.Fanout
	missch   chan func()
}

// New is
func New(conf *conf.Config, dao *block.Dao) *Service {
	s := &Service{
		c:      conf,
		dao:    dao,
		cache:  fanout.New("blockCache", fanout.Worker(1), fanout.Buffer(10240)),
		missch: make(chan func(), 10240),
	}
	s.initWhiteMap()
	go s.missproc()
	return s
}

func (s *Service) missproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.missproc panic(%+v) :\n %s", x, debug.Stack())
			go s.missproc()
		}
	}()
	for {
		f := <-s.missch
		f()
	}
}

func (s *Service) mission(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Error("s.missch full")
	}
}

func (s *Service) initWhiteMap() {
	s.whiteMap = make(map[int64]struct{}, len(s.c.BlockProperty.WhiteList))
	for _, mid := range s.c.BlockProperty.WhiteList {
		s.whiteMap[mid] = struct{}{}
	}
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
