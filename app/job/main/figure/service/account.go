package service

import (
	"context"
	"time"

	"go-common/app/job/main/figure/model"
	"go-common/library/log"
)

// AccountExp handle user exp chenage message
func (s *Service) AccountExp(c context.Context, mid, exp int64) (err error) {
	s.figureDao.UpdateAccountExp(c, mid, exp)
	return
}

// AccountReg handle user register message, init user figure info
func (s *Service) AccountReg(c context.Context, mid int64) (err error) {
	var id int64
	if id, err = s.figureDao.ExistFigure(c, mid); err != nil || id != 0 {
		log.Info("user(%d) already init", mid)
		return
	}
	f := &model.Figure{
		Mid:             mid,
		LawfulScore:     s.c.Figure.Lawful,
		WideScore:       s.c.Figure.Wide,
		FriendlyScore:   s.c.Figure.Friendly,
		BountyScore:     s.c.Figure.Bounty,
		CreativityScore: s.c.Figure.Creativity,
		Ver:             1,
		Ctime:           time.Now(),
		Mtime:           time.Now(),
	}
	if _, err = s.figureDao.SaveFigure(c, f); err != nil {
		return
	}
	s.figureDao.AddFigureInfoCache(c, f)
	return
}

// AccountViewVideo handle user view video message
func (s *Service) AccountViewVideo(c context.Context, mid int64) (err error) {
	s.figureDao.IncArchiveViews(c, mid)
	return
}
