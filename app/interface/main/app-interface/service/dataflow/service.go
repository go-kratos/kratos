package dataflow

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/library/log"
	"go-common/library/log/infoc"
)

// Service is search service
type Service struct {
	c     *conf.Config
	infoc *infoc.Infoc
}

// New is search service initial func
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		infoc: infoc.New(c.Infoc),
	}
	return
}

func (s *Service) Report(c context.Context, eventID, eventType, buvid, fts, messageInfo string, now time.Time) (err error) {
	if err = s.infoc.Info(strconv.FormatInt(now.Unix(), 10), eventID, eventType, buvid, fts, messageInfo); err != nil {
		log.Error("s.infoc2.Info(%v,%v,%v,%v,%v,%v) error(%v)", strconv.FormatInt(now.Unix(), 10), eventID, eventType, buvid, fts, messageInfo, err)
	}
	return
}
