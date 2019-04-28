package service

import (
	"testing"

	"go-common/app/service/main/ugcpay/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func TestIncomeAssetMonthly(t *testing.T) {
	Convey("", t, func() {
		list := []*model.AggrIncomeUserAsset{
			{
				OID:      2333,
				OType:    "archive",
				Currency: "bp",
			},
			{
				OID:      2333,
				OType:    "archive",
				Currency: "bp",
			},
			{
				OID:      2333,
				OType:    "archive",
				Currency: "bp",
			},
		}
		l, page := s.pageIncomeUseAsset(list, 4, 20)
		t.Log(l, page)
	})
}
