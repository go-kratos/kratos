package service

import (
	"context"

	"go-common/app/service/bbq/recsys/conf"
	"go-common/app/service/bbq/recsys/dao"
	postProcess "go-common/app/service/bbq/recsys/service/postprocess"
	"go-common/app/service/bbq/recsys/service/rank"
	"go-common/app/service/bbq/recsys/service/retrieve"
	"go-common/library/log/infoc"
	"go-common/library/stat/prom"
)

// Service struct
type Service struct {
	c                *conf.Config
	dao              *dao.Dao
	infoc            *infoc.Infoc
	retrieverManager *retrieve.RetrieverManager
	recallManager    *retrieve.RecallManager
	rankManager      *RankManager
	rankModelManager *rank.RankModelManager
	filterManager    *FilterManager
	postProcessor    *postProcess.PostProcessor

	//monitor
	businessInfoCount *prom.Prom
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                 c,
		dao:               dao.New(c),
		infoc:             infoc.New(c.Infoc),
		retrieverManager:  retrieve.NewRetrieverManager(),
		recallManager:     retrieve.NewRecallManager(),
		rankManager:       NewRankManager(),
		rankModelManager:  rank.NewRankModelManager(),
		filterManager:     NewFilterManager(),
		postProcessor:     postProcess.NewPostProcessor(),
		businessInfoCount: prom.BusinessInfoCount,
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
