package service

import (
	"context"
	"strings"
	"time"

	"go-common/app/service/main/sms/conf"
	"go-common/app/service/main/sms/dao"
	"go-common/app/service/main/sms/model"
)

// Service struct of service.
type Service struct {
	dao      *dao.Dao
	c        *conf.Config
	template map[string]*model.ModelTemplate
	missch   chan func()
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		dao:      dao.New(c),
		missch:   make(chan func(), 10240),
		template: make(map[string]*model.ModelTemplate),
	}
	s.loadConf(context.Background())
	go s.loadConfproc()
	go s.loadTaskproc()
	return
}

func (s *Service) loadConfproc() {
	for {
		time.Sleep(time.Minute)
		s.loadConf(context.Background())
	}
}

func (s *Service) loadConf(ctx context.Context) {
	var (
		err error
		res []*model.ModelTemplate
	)
	if res, err = s.dao.TemplateByStatus(ctx, model.TemplateStatusApprovel); err != nil {
		return
	}
	tpl := make(map[string]*model.ModelTemplate, len(res))
	for _, v := range res {
		v.Param = parseTemplateParam(v.Template)
		tpl[v.Code] = v
	}
	s.template = tpl
}

func parseTemplateParam(p string) (param []string) {
	mp := make(map[string]struct{})
	ss := strings.SplitAfter(p, "#[")
	for i, v := range ss {
		if i == 0 {
			continue
		}
		k := v[0:strings.Index(v, "]")]
		mp[k] = struct{}{}
	}
	for k := range mp {
		param = append(param, k)
	}
	return
}

// Ping check server ok.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close dao.
func (s *Service) Close() {
	s.dao.Close()
}
