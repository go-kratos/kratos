package service

import (
	"context"
)

// GetPro get project detail
func (s *Service) GetPro(c context.Context, id int, bot bool) (res []byte, err error) {
	res, err = s.dao.GetPro(c, id, bot)
	return
}
