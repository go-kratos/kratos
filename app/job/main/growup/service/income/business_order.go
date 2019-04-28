package income

import "context"

// GetBusinessOrders get business orders
func (s *Service) GetBusinessOrders(c context.Context, limit int64) (result map[int64]bool, err error) {
	var id int64
	result = make(map[int64]bool)
	for {
		var m map[int64]bool
		id, m, err = s.dao.BusinessOrders(c, id, limit)
		if err != nil {
			return
		}
		if len(m) == 0 {
			break
		}
		for k, v := range m {
			result[k] = v
		}
	}
	return
}
