package service

import (
	"context"
	"go-common/app/service/main/archive/api"
)

// Archive3 get a archive by aid.
func (s *Service) Archive3(c context.Context, aid int64) (a *api.Arc, err error) {
	a, err = s.arc.Archive3(c, aid)
	return
}

// Archives3 multi get archives.
func (s *Service) Archives3(c context.Context, aids []int64) (as map[int64]*api.Arc, err error) {
	if len(aids) == 0 {
		as = _emptyArchives3
		return
	}
	as, err = s.arc.Archives3(c, aids)
	return
}
