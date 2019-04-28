package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// forbidTopTypesForAll fn 175=>ASMR
func (s *Service) forbidTopTypesForAll(tid int16) bool {
	return tid == 175
}

func (s *Service) allowOrderUps(mid int64) (ok bool) {
	if mid <= 0 {
		ok = false
		return
	}
	_, ok = s.orderUps[mid]
	return
}

func (s *Service) allowType(typeid int16) (ok bool) {
	_, ok = s.typeCache[typeid]
	if s.forbidTopTypesForAll(typeid) {
		ok = false
	}
	return
}

func (s *Service) allowCopyright(cp int8) (ok bool) {
	ok = archive.InCopyrights(cp)
	return
}

func (s *Service) allowSource(cp int8, source string) (ok bool) {
	ok = cp == archive.CopyrightOriginal || (cp == archive.CopyrightCopy && len(strings.TrimSpace(source)) > 0)
	return
}

func (s *Service) allowTag(tag string) (ok bool) {
	if len(tag) == 0 {
		return
	}
	for _, reg := range _emptyUnicodeReg {
		if reg.MatchString(tag) {
			return
		}
	}
	tags := strings.Split(tag, ",")
	if len(tags) > 12 {
		return
	}
	for _, t := range tags {
		if utf8.RuneCountInString(t) > 30 {
			return
		}
	}
	ok = true
	return
}

func (s *Service) allowDelayTime(dtime xtime.Time) (ok bool) {
	if dtime == 0 {
		ok = true
		return
	}
	const (
		min = int64(4 * time.Hour / time.Second)
		max = int64(15 * 24 * time.Hour / time.Second)
	)
	diff := int64(dtime) - time.Now().Unix()
	ok = min < diff && diff < max
	return
}

func (s *Service) allowHalfMin(c context.Context, mid int64) (ok bool) {
	// 活动等其他业务方运营需要,接触半分钟的限速
	if _, white := s.exemptHalfMinUps[mid]; white {
		return true
	}
	log.Info("halfMin  start | mid(%d).", mid)
	exist, _, _ := s.acc.HalfMin(c, mid)
	log.Info("halfMin from cache | mid(%d) exist(%v).", mid, exist)
	//先判断缓存,ok取反,如果存在则不允许,继续等冷却时间;如果缓存不存在,则默认继续添加冷却时间窗口
	if ok = !exist; ok {
		log.Info("halfMin not exist | mid(%d)", mid)
		s.acc.AddHalfMin(c, mid)
		log.Info("halfMin add cache | mid(%d).", mid)
	}
	log.Info("halfMin  end | mid(%d).", mid)
	return
}

func (s *Service) allowRepeat(c context.Context, mid int64, title string) (ok bool) {
	log.Info("allowRepeat check start | mid(%d) title(%s).", mid, title)
	exist, _ := s.acc.SubmitCache(c, mid, title)
	log.Info("allowRepeat from cache | mid(%d) title(%s) exist(%d).", mid, title, exist)
	if ok = exist == 0; ok {
		log.Info("allowRepeat not exist | mid(%d) title(%s)", mid, title)
		s.acc.AddSubmitCache(c, mid, title)
		log.Info("allowRepeat add cache | mid(%d) title(%s).", mid, title)
	}
	log.Info("allowRepeat check end | mid(%d) title(%s).", mid, title)
	return
}
