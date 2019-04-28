package income

import (
	"context"
)

// GetBubbleMeta .
func (s *Service) GetBubbleMeta(c context.Context) (data map[int64][]int, err error) {
	var id int64
	data = make(map[int64][]int)
	for {
		var meta map[int64][]int
		meta, id, err = s.dao.GetBubbleMeta(c, id, int64(_limitSize))
		if err != nil {
			return
		}
		if len(meta) == 0 {
			break
		}
		for avID, bTypes := range meta {
			data[avID] = append(data[avID], bTypes...)
		}
	}
	return
}
