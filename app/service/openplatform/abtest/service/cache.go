package service

import (
	"context"
	"time"

	"go-common/app/service/openplatform/abtest/model"
)

//setGroupCache 使用src重置分组缓存
func (s *Service) setGroupCache(c context.Context, group int, src []*model.AB, ver int64) (err error) {
	m := make(map[int]*model.AB)
	for _, ab := range src {
		m[ab.ID] = ab
	}
	// s.mutex.Lock()
	// s.abCache[group] = m
	// s._versionID[group] = ver
	// s.mutex.Unlock()
	s.abCache.Store(group, m)
	s._versionID.Store(group, ver)
	return
}

//readGroupCache 获取分组缓存；不存在则ok返回false
func (s *Service) readGroupCache(c context.Context, group int) (res map[int]*model.AB, ok bool) {
	// s.mutex.RLock()
	// res, ok = s.abCache[group]
	// s.mutex.RUnlock()
	var v interface{}
	if v, ok = s.abCache.Load(group); ok {
		res = v.(map[int]*model.AB)
	}
	return
}

//VersionIDListCache 获取缓存版本列表
func (s *Service) VersionIDListCache(c context.Context) (res map[int]int64, err error) {
	res = make(map[int]int64)
	// s.mutex.RLock()
	// for k, v := range s._versionID {
	// 	res[k] = v
	// }
	// s.mutex.RUnlock()
	s._versionID.Range(func(key, value interface{}) bool {
		res[key.(int)] = value.(int64)
		return true
	})
	return
}

//versionID get A/B test version ID by group
func (s *Service) versionID(c context.Context, group int) (ver int64) {
	v, ok := s._versionID.Load(group)
	// s.mutex.RLock()
	// ver, ok = s._versionID[group]
	// s.mutex.RUnlock()
	if ok {
		ver = v.(int64)
		return
	}
	ver = time.Now().Unix()
	return
}
