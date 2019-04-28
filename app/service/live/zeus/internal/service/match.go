package service

import (
	"errors"
	"go-common/app/service/live/zeus/expr"
	"go-common/app/service/live/zeus/internal/model"
	"go-common/library/log"
	"sync/atomic"
)

type MatchService struct {
	matcher atomic.Value
}

func NewMatchService(config string) (*MatchService, error) {
	s := &MatchService{}
	if err := s.Reload(config); err != nil {
		return nil, err
	}
	return s, nil
}

func (m *MatchService) Reload(config string) error {
	matcher, err := model.NewMatcher(config)
	if err != nil {
		log.Error("Match Service reload config failed, config:%s", config)
		return err
	}
	m.matcher.Store(matcher)
	log.Info("Match Service reload config success, config:%s", config)
	return nil
}

func (m *MatchService) Match(group string, env expr.Env) (bool, string, error) {
	matcher := m.matcher.Load().(*model.Matcher)
	if matcher == nil {
		return false, "", errors.New("matcher not initialized")
	}
	return matcher.Match(group, env)
}
