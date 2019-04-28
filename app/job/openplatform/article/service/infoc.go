package service

import (
	"go-common/library/log"
	binfoc "go-common/library/log/infoc"
	"go-common/library/stat/prom"
)

// 用户阅读专栏时长上报
type readInfo struct {
	aid      int64
	mid      int64
	buvid    string
	ip       string
	duration int64
	from     string
}

// ReadInfoc .
func (s *Service) ReadInfoc(aid int64, mid int64, buvid string, ip string, duration int64, from string) {
	s.infoc(readInfo{
		aid:      aid,
		mid:      mid,
		buvid:    buvid,
		ip:       ip,
		duration: duration,
		from:     from,
	})
}

func (s *Service) infoc(i interface{}) {
	select {
	case s.logCh <- i:
	default:
		log.Warn("infocproc chan full")
	}
}

// writeInfoc
func (s *Service) infocproc() {
	var (
		readInfoc = binfoc.New(s.c.ReadInfoc)
	)
	for {
		i, ok := <-s.logCh
		if !ok {
			log.Warn("infoc proc exit")
			return
		}
		prom.BusinessInfoCount.State("infoc_channel", int64(len(s.logCh)))
		switch l := i.(type) {
		case readInfo:
			readInfoc.Info(l.aid, l.mid, l.buvid, l.ip, l.duration, l.from)
			log.Info("infocproc readInfoc param(%+v)", l)
		}
	}
}
