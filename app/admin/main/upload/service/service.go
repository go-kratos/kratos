package service

import (
	"context"
	"time"

	"go-common/app/admin/main/upload/conf"
	"go-common/app/admin/main/upload/dao"
	"go-common/app/admin/main/upload/model"
	"go-common/library/database/orm"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// Service struct
type Service struct {
	c           *conf.Config
	orm         *gorm.DB
	bucketCache map[string]*model.Bucket
	dao         *dao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		orm: orm.NewMySQL(c.ORM),
		dao: dao.New(c),
	}
	s.bucketCache = make(map[string]*model.Bucket)
	go s.cacheproc()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return nil
}

// Close Service
func (s *Service) Close() {
	s.orm.Close()
}

func (s *Service) cacheproc() {
	for {
		s.loadcache()
		time.Sleep(5 * time.Minute)
	}
}

func (s *Service) loadcache() {
	var buckets []*model.Bucket
	if err := s.orm.Table("bucket").Order("id desc").Limit(1000).Find(&buckets).Error; err != nil {
		log.Error("read bucket error(%v)", err)
		return
	}
	b := make(map[string]*model.Bucket)
	for _, v := range buckets {
		b[v.BucketName] = v
	}
	s.bucketCache = b
}
