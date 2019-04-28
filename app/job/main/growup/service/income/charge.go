package income

import (
	incomeD "go-common/app/job/main/growup/dao/income"
)

// AvChargeSvr av charge service
type AvChargeSvr struct {
	dao *incomeD.Dao
}

// NewAvChargeSvr new av charge service
func NewAvChargeSvr(dao *incomeD.Dao) (svr *AvChargeSvr) {
	return &AvChargeSvr{
		dao: dao,
	}
}
