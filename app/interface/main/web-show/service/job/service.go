package job

import (
	"context"
	"time"

	"go-common/app/interface/main/web-show/conf"
	"go-common/app/interface/main/web-show/dao/job"
	jobmdl "go-common/app/interface/main/web-show/model/job"
	"go-common/library/log"
)

var (
	_emptyJobs = make([]*jobmdl.Job, 0)
)

// Service struct
type Service struct {
	dao   *job.Dao
	cache []*jobmdl.Job
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{}
	s.dao = job.New(c)
	s.cache = _emptyJobs
	s.reload()
	go s.loadproc()
	return
}

// jobproc load job infos to cache
func (s *Service) loadproc() {
	for {
		s.reload()
		time.Sleep(time.Duration(conf.Conf.Reload.Jobs))
	}
}

// reload
func (s *Service) reload() {
	js, err := s.dao.Jobs(context.Background())
	if err != nil {
		log.Error("s.job.Jobs error(%v)", err)
		return
	} else if len(js) == 0 {
		s.cache = _emptyJobs
	}
	cates, err := s.dao.Categories(context.Background())
	if err != nil {
		log.Error("job.Categories error(%v)", err)
		return
	}
	cs := make(map[int]string, len(cates))
	for _, cate := range cates {
		cs[cate.ID] = cate.Name
	}
	for _, j := range js {
		j.JobsCla = cs[j.CateID]
		j.Location = cs[j.AddrID]
	}
	s.cache = js
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
