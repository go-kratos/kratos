package service

import (
	"context"
	"go-common/app/interface/main/videoup/model/archive"
)

func (s *Service) dealSubtitle(c context.Context, sub *archive.Subtitle, aid, mid int64, ip string) (err error) {
	if sub != nil {
		var open bool
		if sub.Open == 1 {
			open = true
		}
		s.sub.Update(c, aid, open, sub.Lan)
	}
	return
}
