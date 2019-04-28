package service

import (
	"context"

	"go-common/app/infra/notify/model"
	"go-common/app/infra/notify/notify"
	"go-common/library/ecode"
)

// Pub pub message.
func (s *Service) Pub(c context.Context, arg *model.ArgPub) (err error) {
	pc, ok := s.pubConfs[key(arg.Group, arg.Topic)]
	if !ok {
		err = ecode.AccessDenied
		return
	}
	s.plock.RLock()
	pub, ok := s.pubs[key(arg.Group, arg.Topic)]
	s.plock.RUnlock()
	if !ok {
		pub, err = notify.NewPub(pc, s.c)
		if err != nil {
			return
		}
		s.plock.Lock()
		s.pubs[key(arg.Group, arg.Topic)] = pub
		s.plock.Unlock()
	}
	if !pub.Auth(arg.AppSecret) {
		err = ecode.AccessDenied
		return
	}
	err = pub.Send([]byte(arg.Key), []byte(arg.Msg))
	return
}
