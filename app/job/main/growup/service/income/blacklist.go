package income

import (
	"context"
)

// Blacklist map[ctype]map[int64]bool
func (s *Service) Blacklist(c context.Context, limit int64) (m map[int]map[int64]bool, err error) {
	var id int64
	m = make(map[int]map[int64]bool)
	for {
		var ab map[int][]int64
		ab, id, err = s.dao.Blacklist(c, id, limit)
		if err != nil {
			return
		}
		if len(ab) == 0 {
			break
		}
		for ctype, avIDs := range ab {
			for _, avID := range avIDs {
				if _, ok := m[ctype]; ok {
					m[ctype][avID] = true
				} else {
					m[ctype] = map[int64]bool{
						avID: true,
					}
				}
			}
		}
	}
	return
}
