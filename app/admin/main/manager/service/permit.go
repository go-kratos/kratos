package service

import (
	"context"

	"go-common/library/log"
	"go-common/library/net/http/blademaster/middleware/permit"
)

const (
	_sessUnKey  = "username"
	_sessUIDKey = "uid"
)

// Login .
func (s *Service) Login(ctx context.Context, mngsid, dsbsid string) (sid, uname string, err error) {
	si := s.session(ctx, mngsid)
	var username string
	if si.Get(_sessUnKey) == nil {
		if username, err = s.dao.VerifyDsb(ctx, dsbsid); err != nil {
			log.Error("s.dao.VerifyDsb error(%v)", err)
			return
		}
		si.Set(_sessUnKey, username)
		si.Set(_sessUIDKey, s.userIds[username])
		if err = s.dao.SetSession(ctx, si); err != nil {
			log.Error("s.dao.SetSession(%v) error(%v)", si, err)
			err = nil
		}
	} else {
		username = si.Get(_sessUnKey).(string)
	}
	sid = si.Sid
	uname = username
	return
}

// session .
func (s *Service) session(ctx context.Context, sid string) (res *permit.Session) {
	if res, _ = s.dao.Session(ctx, sid); res == nil {
		res = s.dao.NewSession(ctx)
	}
	return
}
