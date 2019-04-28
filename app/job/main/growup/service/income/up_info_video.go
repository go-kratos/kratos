package income

import (
	"context"

	model "go-common/app/job/main/growup/model/income"
)

// Signed signed up by business
func (s *Service) Signed(c context.Context, business string, limit int64) (m map[int64]*model.Signed, err error) {
	var offset int64
	m = make(map[int64]*model.Signed)
	for {
		var us map[int64]*model.Signed
		offset, us, err = s.dao.Ups(c, business, offset, limit)
		if err != nil {
			return
		}
		if len(us) == 0 {
			break
		}
		for k, v := range us {
			if v.AccountState == 3 || v.AccountState == 5 || v.AccountState == 6 {
				m[k] = v
			}
		}
	}
	return
}
