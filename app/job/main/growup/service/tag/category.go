package tag

import (
	"context"
)

// 获取所有分区
func (s *Service) getArchiveCategory(c context.Context, ctype int) ([]int64, error) {
	switch ctype {
	case _video:
		return s.getAvCategory(c)
	case _column:
		return s.getCmCategory(c)
	case _bgm:
		return []int64{0}, nil
	}
	return []int64{}, nil
}

func (s *Service) getAvCategory(c context.Context) (ids []int64, err error) {
	ids = make([]int64, 0)
	allCategory, err := s.dao.GetVideoTypes(c)
	if err != nil {
		return
	}
	categoryMap := make(map[int64]struct{})
	for _, id := range allCategory {
		if _, ok := categoryMap[id]; !ok {
			categoryMap[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return
}

func (s *Service) getCmCategory(c context.Context) (ids []int64, err error) {
	ids = make([]int64, 0)
	allCategory, err := s.dao.GetColumnTypes(c)
	if err != nil {
		return
	}
	categoryMap := make(map[int64]struct{})
	for _, id := range allCategory {
		if _, ok := categoryMap[id]; !ok {
			categoryMap[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return
}
