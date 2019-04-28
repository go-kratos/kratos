package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/push/conf"
	"go-common/app/admin/main/push/dao"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Service biz service def.
type Service struct {
	c          *conf.Config
	dao        *dao.Dao
	DB         *gorm.DB
	httpClient *bm.Client
	dpClient   *bm.Client
	partitions map[int]string
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		dao:        dao.New(c),
		httpClient: bm.NewClient(c.HTTPClient),
		dpClient:   bm.NewClient(c.DPClient),
		partitions: make(map[int]string),
	}
	s.DB = s.dao.DB.Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("dtime = ?", 0)
	})
	s.loadPartitions()
	go s.loadPartitionsproc()
	go s.cleanDiskFilesproc()
	return s
}

func (s *Service) loadPartitions() (err error) {
	m, err := s.dao.Partitions(context.Background())
	if err != nil {
		return
	}
	if len(m) > 0 {
		s.partitions = m
	}
	return
}

func (s *Service) loadPartitionsproc() {
	for {
		if err := s.loadPartitions(); err != nil {
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Minute)
	}
}

func (s *Service) cleanDiskFilesproc() {
	for {
		fs, err := ioutil.ReadDir(conf.Conf.Cfg.MountDir)
		if err != nil {
			log.Error("s.cleanFilesproc() read dir error(%v)", err)
			time.Sleep(time.Minute)
			continue
		}
		divDate := time.Now().Add(-time.Duration(24*conf.Conf.Cfg.DiskFileExpireDay) * time.Hour).Format("20060102")
		div, _ := strconv.ParseInt(divDate, 10, 64)
		for _, f := range fs {
			if !f.IsDir() {
				continue
			}
			d, _ := strconv.ParseInt(f.Name(), 10, 64)
			if d < div {
				dir := fmt.Sprintf("%s/%s", strings.TrimSuffix(conf.Conf.Cfg.MountDir, "/"), f.Name())
				if err = os.RemoveAll(dir); err != nil {
					log.Error("s.cleanFilesproc() remove dir(%s) error(%v)", dir, err)
					time.Sleep(time.Minute)
				}
			}
		}
		time.Sleep(time.Hour)
	}
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Wait wait all closed.
func (s *Service) Wait() {}

// Close close all dao.
func (s *Service) Close() {
	s.dao.Close()
}
