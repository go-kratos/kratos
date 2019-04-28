package service

import "context"

// ExpireUpRatingCache ...
func (s *Service) ExpireUpRatingCache(c context.Context, mid int64) error {
	return s.dao.ExpireUpRatingCache(c, mid)
}
