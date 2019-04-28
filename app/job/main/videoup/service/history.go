package service

import (
	"context"
	"time"

	"go-common/library/log"
)

func (s *Service) delArcEditHistory(limit int64) (delRows int64, err error) {
	var (
		c      = context.TODO()
		mtime  = time.Now().Add(-2 * 30 * 24 * time.Hour)
		before = time.Date(mtime.Year(), mtime.Month(), mtime.Day(), mtime.Hour(), 0, 0, 0, mtime.Location())
	)
	if delRows, err = s.arc.DelArcEditHistoryBefore(c, before, limit); err != nil {
		log.Error("s.arc.TxDelArcEditHistoryBefore(%s) error(%v)", mtime.Format("2006-01-02 15:04:05"), err)
		return
	}
	log.Info("delArchiveHistory before mtime(%s) rows(%d)", before.Format("2006-01-02 15:04:05"), delRows)
	return
}

func (s *Service) delArcVideoEditHistory(limit int64) (delRows int64, err error) {
	var (
		c      = context.TODO()
		mtime  = time.Now().Add(-2 * 30 * 24 * time.Hour)
		before = time.Date(mtime.Year(), mtime.Month(), mtime.Day(), mtime.Hour(), 0, 0, 0, mtime.Location())
	)
	if delRows, err = s.arc.DelArcVideoEditHistoryBefore(c, before, limit); err != nil {
		log.Error("s.arc.TxDelArcVideoEditHistoryBefore(%s) error(%v)", mtime.Format("2006-01-02 15:04:05"), err)
		return
	}
	log.Info("delArcVideoEditHistory before mtime(%s) rows(%d)", before.Format("2006-01-02 15:04:05"), delRows)
	return
}
