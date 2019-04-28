package service

import (
	"context"

	"go-common/library/log"
)

// add a new resource ID into the to push list
func (s *Service) newPush(ctx context.Context, resID int) (err error) {
	if err = s.dao.ZAddPush(ctx, resID); err != nil {
		log.Error("NewPush Redis for ResID: %d, Error: %v", resID, err)
	}
	return
}
