package service

import (
	"go-common/app/interface/openplatform/monitor-end/model"
	"go-common/library/log"
	"go-common/library/log/infoc"
)

// writeInfoc
func (s *Service) infocproc() {
	var (
		collectInfoc = infoc.New(s.c.CollectInfoc)
	)
	for {
		i, ok := <-s.infoCh
		if !ok {
			log.Warn("infoc proc exit")
			return
		}
		switch l := i.(type) {
		case model.CollectParams:
			collectInfoc.Info(l.Source, l.Product, l.Event, l.SubEvent, l.Code, l.ExtJSON, l.Mid, l.IP, l.Buvid, l.UserAgent)
		}
	}
}
