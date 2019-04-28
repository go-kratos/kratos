package service

import (
	"context"
	"time"
)

func (s *Service) recommendAuthorproc() {
	for {
		s.calRecommendAuthor()
		time.Sleep(time.Hour)
	}
}

func (s *Service) calRecommendAuthor() {
	var (
		mids []int64
		err  error
		m    = make(map[int64][]int64)
	)
	for {
		if mids, err = s.dao.MidsByPublishTime(context.TODO(), time.Now().Unix()-86400*s.c.Job.StatDays); err != nil {
			time.Sleep(time.Minute)
			continue
		}
		for _, mid := range mids {
			if categoris, ok := s.statAuthorInfo(mid); ok {
				for _, category := range categoris {
					m[category] = append(m[category], mid)
				}
				s.dao.AddAuthorMostCategories(context.TODO(), mid, categoris)
			}
		}
		s.dao.AddcategoriesAuthors(context.TODO(), m)
		return
	}
}

func (s *Service) statAuthorInfo(mid int64) (categories []int64, res bool) {
	var (
		m     map[int64][2]int64
		err   error
		words int64
		cm    = make(map[int64]int64)
		max   int64
	)

	if m, err = s.dao.StatByMid(context.TODO(), mid); err != nil {
		return
	}
	for cid, data := range m {
		if cate, ok := s.categoriesMap[cid]; ok {
			if cate.ParentID != 0 {
				cid = cate.ParentID
			}
			words += data[1]
			if _, ok := cm[cid]; ok {
				cm[cid] += data[0]
			} else {
				cm[cid] = data[0]
			}
		}
	}
	if words < s.c.Job.Words {
		return
	}
	res = true
	for k, v := range cm {
		if v > max {
			max = v
			categories = []int64{k}
		} else if v == max {
			categories = append(categories, k)
		}
	}
	return
}
