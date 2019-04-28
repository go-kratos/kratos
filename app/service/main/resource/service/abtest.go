package service

import (
	"context"
	"strings"
	"time"

	"go-common/app/service/main/resource/model"
	"go-common/library/log"
)

// AbTest get abtest by group name
func (s *Service) AbTest(c context.Context, names, ipaddr string) (res map[string]*model.AbTest) {
	var (
		now = time.Now().Unix()
		mis []string
	)
	res = make(map[string]*model.AbTest)
	ns := strings.Split(names, ",")
	s.abTestLock.Lock()
	for _, n := range ns {
		if ab, ok := s.abTestCache[n]; ok && (now-ab.UTime <= 300) {
			r := &model.AbTest{}
			*r = *ab
			res[n] = r
		} else {
			mis = append(mis, n)
		}
	}
	s.abTestLock.Unlock()
	if len(mis) > 0 {
		nabs, err := s.abtest.AbTest(c, strings.Join(mis, ","), ipaddr)
		if err != nil {
			log.Error("AbTest(%v, %v) error(%v)", mis, ipaddr, err)
			return
		}
		for _, nab := range nabs {
			// add in res
			r := &model.AbTest{}
			*r = *nab
			res[r.Name] = r
		}
		s.abTestLock.Lock()
		s.UpdateAbTestCache(nabs)
		s.abTestLock.Unlock()
	}
	return
}

// UpdateAbTestCache update abtest
func (s *Service) UpdateAbTestCache(nabs []*model.AbTest) {
	var now = time.Now().Unix()
	for _, nab := range nabs {
		nab.UTime = now
		s.abTestCache[nab.Name] = nab
	}
}
