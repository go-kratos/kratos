package tag

import (
	"context"
)

// GetTagUpInfoMID get tag_up_info mid map
func (s *Service) GetTagUpInfoMID(c context.Context, tags []int64) (tagMID map[int64][]int64, err error) {
	tagMID = make(map[int64][]int64)
	from, limit := 0, 2000
	for {
		count := 0
		count, err = s.dao.GetTagUpInfoByTag(c, tags, from, limit, tagMID)
		if err != nil {
			return
		}
		if count < limit {
			break
		}
		from += limit
	}
	return
}
