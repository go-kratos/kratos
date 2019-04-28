package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const _defaultPrivacy = 1

// SettingInfo get setting info.
func (s *Service) SettingInfo(c context.Context, mid int64) (data *model.Setting, err error) {
	data = new(model.Setting)
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		data.Privacy = s.privacy(errCtx, mid)
		return nil
	})
	group.Go(func() error {
		data.IndexOrder = s.indexOrder(errCtx, mid)
		return nil
	})
	group.Wait()
	return
}

// PrivacyModify privacy modify.
func (s *Service) PrivacyModify(c context.Context, mid int64, field string, value int) (err error) {
	privacy := s.privacy(c, mid)
	for k, v := range privacy {
		if field == k && value == v {
			err = ecode.NotModified
			return
		}
	}
	if err = s.dao.PrivacyModify(c, mid, field, value); err == nil {
		s.dao.DelPrivacyCache(c, mid)
	}
	return
}

// PrivacyBatchModify privacy batch modify.
func (s *Service) PrivacyBatchModify(c context.Context, mid int64, data map[string]int) (err error) {
	group, errCtx := errgroup.WithContext(c)
	for k, v := range data {
		field := k
		value := v
		group.Go(func() error {
			if e := s.PrivacyModify(errCtx, mid, field, value); e != nil {
				log.Warn("PrivacyBatchModify mid(%d) filed(%s) value(%d) error(%v)", mid, field, value, e)
			}
			return nil
		})
	}
	group.Wait()
	return
}

// IndexOrderModify index order modify
func (s *Service) IndexOrderModify(c context.Context, mid int64, orderNum []string) (err error) {
	var orderStr []byte
	if orderStr, err = json.Marshal(orderNum); err != nil {
		log.Error("index order modify json.Marshal(%v) error(%v)", orderNum, err)
		err = ecode.RequestErr
		return
	}
	if err = s.dao.IndexOrderModify(c, mid, string(orderStr)); err == nil {
		s.cache.Do(c, func(c context.Context) {
			var cacheData []*model.IndexOrder
			for _, v := range orderNum {
				i, _ := strconv.Atoi(v)
				cacheData = append(cacheData, &model.IndexOrder{ID: i, Name: model.IndexOrderMap[i]})
			}
			s.dao.SetIndexOrderCache(c, mid, cacheData)
		})
	}
	return
}

func (s *Service) privacy(c context.Context, mid int64) (data map[string]int) {
	var (
		privacy  map[string]int
		err      error
		addCache = true
	)
	if data, err = s.dao.PrivacyCache(c, mid); err != nil {
		addCache = false
	} else if data != nil {
		return
	}
	if privacy, err = s.dao.Privacy(c, mid); err != nil || len(privacy) == 0 {
		data = model.DefaultPrivacy
	} else {
		data = fmtPrivacy(privacy)
	}
	if addCache {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetPrivacyCache(c, mid, data)
		})
	}
	return
}

func (s *Service) indexOrder(c context.Context, mid int64) (data []*model.IndexOrder) {
	var (
		indexOrderStr string
		err           error
		addCache      = true
	)
	if data, err = s.dao.IndexOrderCache(c, mid); err != nil {
		addCache = false
	} else if len(data) != 0 {
		return
	}
	if indexOrderStr, err = s.dao.IndexOrder(c, mid); err != nil || indexOrderStr == "" {
		data = model.DefaultIndexOrder
	} else {
		orderNum := make([]string, 0)
		if err = json.Unmarshal([]byte(indexOrderStr), &orderNum); err != nil {
			log.Error("indexOrder mid(%d) json.Unmarshal(%s) error(%v)", mid, indexOrderStr, err)
			addCache = false
			s.cache.Do(c, func(c context.Context) {
				s.fixIndexOrder(c, mid, indexOrderStr)
			})
			data = model.DefaultIndexOrder
		} else {
			extraOrder := make(map[int]string)
			for _, v := range orderNum {
				if index, err := strconv.Atoi(v); err != nil {
					continue
				} else if name, ok := model.IndexOrderMap[index]; ok {
					data = append(data, &model.IndexOrder{ID: index, Name: name})
					extraOrder[index] = name
				}
			}
			for i, v := range model.IndexOrderMap {
				if _, ok := extraOrder[i]; !ok {
					data = append(data, &model.IndexOrder{ID: i, Name: v})
				}
			}
		}
	}
	if addCache {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetIndexOrderCache(c, mid, data)
		})
	}
	return
}

func fmtPrivacy(privacy map[string]int) (data map[string]int) {
	data = make(map[string]int, len(model.PrivacyFields))
	for _, v := range model.PrivacyFields {
		if value, ok := privacy[v]; ok {
			data[v] = value
		} else {
			data[v] = _defaultPrivacy
		}
	}
	return
}

func (s *Service) fixIndexOrder(c context.Context, mid int64, indexOrderStr string) {
	fixStr := strings.Replace(strings.TrimRight(strings.TrimLeft(indexOrderStr, "["), "]"), "\"", "", -1)
	fixArr := strings.Split(fixStr, ",")
	fixByte, err := json.Marshal(fixArr)
	if err != nil {
		log.Error("fixIndexOrder mid(%d) indexOrder(%s) error(%v)", mid, indexOrderStr, err)
		return
	}
	if err := s.dao.IndexOrderModify(c, mid, string(fixByte)); err == nil {
		s.dao.DelIndexOrderCache(c, mid)
	}
}
