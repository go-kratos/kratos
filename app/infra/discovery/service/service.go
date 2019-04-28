package service

import (
	"context"
	"os"
	"sync"

	"go-common/app/infra/discovery/conf"
	"go-common/app/infra/discovery/dao"
	bm "go-common/library/net/http/blademaster"
)

// Service discovery main service
type Service struct {
	c        *conf.Config
	client   *bm.Client
	registry *dao.Registry
	nodes    *dao.Nodes
	tLock    sync.RWMutex
	tree     map[int64]string // treeid->appid
	env      *env
}

type env struct {
	Region string
	Zone   string
}

// New get a discovery service
func New(c *conf.Config) (s *Service, cancel context.CancelFunc) {
	s = &Service{
		c:        c,
		client:   bm.NewClient(c.HTTPClient),
		registry: dao.NewRegistry(),
		nodes:    dao.NewNodes(c),
		tree:     make(map[int64]string),
	}
	s.getEnv()
	s.syncUp()
	cancel = s.regSelf()
	go s.nodesproc()
	return
}

func (s *Service) getEnv() {
	s.env = &env{
		Region: os.Getenv("REGION"),
		Zone:   os.Getenv("ZONE"),
	}
}
