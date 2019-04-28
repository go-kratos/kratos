package service

import (
	"context"
	"time"

	"go-common/app/service/main/relation/model"
)

// Prompt incr user prompt count and return if prompt.
func (s *Service) Prompt(c context.Context, m *model.ArgPrompt) (b bool, err error) {
	r, err := s.Relation(c, m.Mid, m.Fid)
	if err != nil {
		return
	}
	if r != nil && r.Following() {
		return
	}
	ucount, bcount, err := s.dao.IncrPromptCount(c, m.Mid, m.Fid, time.Now().Unix(), m.Btype)
	if err != nil {
		return
	}
	if s.c.Relation.Bcount > bcount && s.c.Relation.Ucount > ucount {
		b = true
	}
	return
}

// ClosePrompt close prompt.
func (s *Service) ClosePrompt(c context.Context, m *model.ArgPrompt) (err error) {
	return s.dao.ClosePrompt(c, m.Mid, m.Fid, time.Now().Unix(), m.Btype)
}
