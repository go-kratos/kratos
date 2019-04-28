package service

import (
	"context"

	coinm "go-common/app/service/main/coin/model"
)

// PutCoinInfo handle user coin chenage message
func (s *Service) PutCoinInfo(c context.Context, msg *coinm.DataBus) (err error) {
	s.figureDao.PutCoinCount(c, msg.Mid)
	return
}
