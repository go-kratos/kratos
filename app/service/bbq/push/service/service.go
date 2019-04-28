package service

import (
	"time"

	"go-common/app/service/bbq/push/conf"
	"go-common/app/service/bbq/push/dao"
	"go-common/library/log/infoc"

	"github.com/Dai0522/workpool"
)

var (
	// Infoc lancer data collect
	Infoc *infoc.Infoc
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	wp  *workpool.Pool
}

// New init
func New(c *conf.Config) (s *Service) {
	workpoolConf := &workpool.PoolConfig{
		MaxWorkers:     c.WorkPool.MaxWorkers,
		MaxIdleWorkers: c.WorkPool.MaxIdleWorkers,
		MinIdleWorkers: c.WorkPool.MinIdleWorkers,
		KeepAlive:      time.Duration(c.WorkPool.KeepAlive),
	}
	wp, err := workpool.NewWorkerPool(c.WorkPool.Capacity, workpoolConf)
	if err != nil {
		panic(err)
	}

	Infoc = infoc.New(c.Infoc)

	return &Service{
		c:   c,
		dao: dao.New(c),
		wp:  wp,
	}
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
