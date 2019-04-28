package service

import (
	"context"

	"go-common/library/log"
)

// setName set user name.
func (s *Service) setName(mid int64) (err error) {
	var (
		name string
	)
	if name, err = s.dao.Name(context.TODO(), mid); err != nil {
		log.Error("s.dao.Name(%d) error(%v)", mid, err)
		return
	}
	if len(name) > 0 {
		if err = s.dao.SetName(context.TODO(), mid, name); err != nil {
			return
		}
	}
	return
}
