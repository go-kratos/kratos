package income

import (
	"context"
	"time"

	model "go-common/app/job/main/growup/model/income"
)

func (s *Service) columnCharges(c context.Context, date time.Time, ch chan []*model.ColumnCharge) (err error) {
	defer func() {
		close(ch)
	}()
	var id int64
	for {
		var charges []*model.ColumnCharge
		charges, err = s.dao.ColumnDailyCharge(c, date, id, _limitSize)
		if err != nil {
			return
		}
		ch <- charges
		if len(charges) < _limitSize {
			break
		}
		id = charges[len(charges)-1].ID
	}
	return
}
