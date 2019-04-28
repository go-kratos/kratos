package web

import (
	"context"

	"go-common/library/log"
)

// PgcFull pgc full .
func (s *Service) PgcFull(ctx context.Context, tp int, pn, ps int64, source string) (res interface{}, err error) {
	if res, err = s.dao.PgcFull(ctx, tp, pn, ps, source); err != nil {
		log.Error("s.dao.PgcFull error(%v)", err)
	}
	return
}

// PgcIncre pgc incre .
func (s *Service) PgcIncre(ctx context.Context, tp int, pn, ps, start, end int64, source string) (res interface{}, err error) {
	if res, err = s.dao.PgcIncre(ctx, tp, pn, ps, start, end, source); err != nil {
		log.Error("s.dao.PgcIncre error(%s)", err)
	}
	return
}
