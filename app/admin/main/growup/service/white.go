package service

import (
	"context"
)

// InsertWhite insert white
func (s *Service) InsertWhite(c context.Context, mid int64, typ int) (err error) {
	_, err = s.dao.InsertWhitelist(c, mid, typ)
	return
}
