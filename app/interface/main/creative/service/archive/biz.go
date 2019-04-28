package archive

import (
	"context"
	"time"

	model "go-common/app/interface/main/creative/model/archive"
	"go-common/library/log"
)

// BIZsByTime list bizs by time and type
func (s *Service) BIZsByTime(c context.Context, start, end *time.Time, tp int8) (bizs []*model.BIZ, err error) {
	if bizs, err = s.arc.BIZsByTime(c, start, end, tp); err != nil {
		log.Error("s.arc.BIZsByTime error(%+v)", err)
	}
	return
}
