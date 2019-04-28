package service

import (
	"context"
	"time"

	"go-common/app/interface/main/upload/conf"
	"go-common/app/interface/main/upload/dao"
	"go-common/app/interface/main/upload/model"
	"go-common/library/log"
)

// Service .
type Service struct {
	dao         *dao.Dao
	bfs         *dao.Bfs
	c           *conf.Config
	bucketCache map[string]*model.Bucket
}

// New .
func New(c *conf.Config) *Service {
	s := &Service{
		dao:         dao.NewDao(c),
		bfs:         dao.NewBfs(c),
		c:           c,
		bucketCache: make(map[string]*model.Bucket),
	}
	go s.cacheproc()

	return s
}

// Ping .
func (s *Service) Ping(c context.Context) (err error) {
	return
}

func (s *Service) cacheproc() {
	for {
		s.loadBucketCache()
		time.Sleep(5 * time.Minute)
	}
}

func (s *Service) loadBucketCache() {
	var (
		bMap map[string]*model.Bucket
		err  error
	)
	if bMap, err = s.dao.Buckets(); err != nil {
		log.Error("get bucket meta failed! error(%v)", err)
		return
	}

	s.bucketCache = bMap
}

// GetRateLimit return rate limit of bucket and dir
func (s *Service) GetRateLimit(bucket, dir string) (model.DirRateConfig, bool) {
	b, ok := s.bucketCache[bucket]
	if !ok {
		return model.DirRateConfig{}, false
	}
	config, ok := b.DirLimit[dir]
	if config == nil {
		return model.DirRateConfig{}, false
	}
	return config.Rate, ok
}
