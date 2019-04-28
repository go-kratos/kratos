package service

import "context"

// AddAuthority .
func (s *Service) AddAuthority(c context.Context, mids []int64) (int64, error) {
	return s.dao.AddAuthority(c, mids)
}

// RemoveAuthority .
func (s *Service) RemoveAuthority(c context.Context, mids []int64) (int64, error) {
	return s.dao.RmAuthority(c, mids)
}
