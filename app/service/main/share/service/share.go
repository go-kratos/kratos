package service

import (
	"context"

	"go-common/app/service/main/share/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Add add share
func (s *Service) Add(ctx context.Context, p *model.ShareParams) (shared int64, err error) {
	typ, ok := s.allowType[p.TP]
	if !ok || typ == "" {
		log.Error("share type(%d) not support or typ is empty typ(%s)", p.TP, typ)
		err = ecode.RequestErr
		return
	}
	if err = s.dao.PubShare(ctx, p); err != nil {
		log.Error("s.dao.PubShare oid(%d) mid(%d) error(%v)", p.OID, p.MID, err)
		err = nil
	}
	log.Info("pubshare success: oid(%d) mid(%d) type(%d) ip(%s)", p.OID, p.MID, p.TP, p.IP)
	if shared, err = s.dao.Add(ctx, p); err != nil {
		log.Error("s.dao.Add mid(%d) oid(%d) tp(%d) ip(%s) error(%v)", p.MID, p.OID, p.TP, p.IP, err)
		return
	}
	if err = s.dao.PubStatShare(ctx, typ, p.OID, shared); err != nil {
		log.Error("s.dao.PubArchiveShare oid(%d) error(%v)", p.OID, err)
		err = nil
	}
	log.Info("oid-%d mid-%d type-%d count-%d", p.OID, p.MID, p.TP, shared)
	return
}

// Stat return oid shared count
func (s *Service) Stat(ctx context.Context, oid int64, tp int) (shared int64, err error) {
	if _, ok := s.allowType[tp]; !ok {
		err = ecode.RequestErr
		return
	}
	return s.dao.ShareCount(ctx, oid, tp)
}

// Stats return oids shares
func (s *Service) Stats(ctx context.Context, oids []int64, tp int) (shares map[int64]int64, err error) {
	if _, ok := s.allowType[tp]; !ok {
		err = ecode.RequestErr
		return
	}
	return s.dao.Shares(ctx, oids, tp)
}
