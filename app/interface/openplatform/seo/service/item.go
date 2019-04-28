package service

import (
	"context"
)

// GetItem get item detail
func (s *Service) GetItem(c context.Context, id int, bot bool) (res []byte, err error) {
	res, err = s.dao.GetItem(c, id, bot)
	return
}
