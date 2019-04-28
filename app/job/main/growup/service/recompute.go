package service

import (
	"context"

	"go-common/app/job/main/growup/model"
)

// AvIncomes av_income
func (s *Service) AvIncomes(c context.Context, mid int64, date string) (result map[int64]*model.Patch, err error) {
	avs, err := s.dao.GetAvs(c, date, mid)
	if err != nil {
		return
	}

	var avIds []int64
	for avID := range avs {
		avIds = append(avIds, avID)
	}

	charges, err := s.dao.GetAvCharges(c, avIds, date)
	if err != nil {
		return
	}
	result = avIncomes(charges, avs)
	return
}

// result key: av_id, value: income
func avIncomes(avCharges map[int64]int64, avs map[int64]*model.Av) (result map[int64]*model.Patch) {
	var totalCharge int64
	for _, charge := range avCharges {
		totalCharge += charge
	}
	tax := int64(Round(Tax(Div(float64(totalCharge), 100))*100, 0))
	netIncome := totalCharge - tax
	percent := Div(float64(netIncome), float64(totalCharge))

	result = make(map[int64]*model.Patch)
	for avID, charge := range avCharges {
		avIncome := int64(float64(charge) * percent)
		avTax := int64(Round(Mul(float64(tax), Div(float64(avIncome), float64(netIncome))), 0))
		result[avID] = &model.Patch{
			Tax:       avTax,
			Income:    avIncome,
			OldTax:    avs[avID].TaxMoney,
			OldIncome: avs[avID].Income,
			MID:       avs[avID].MID,
			TagID:     avs[avID].TagID,
		}
	}
	return
}
