package dao

import (
	"context"

	"go-common/app/interface/main/growup/model"
)

const (
	_upBillSQL = "SELECT mid,first_income,max_income,total_income,av_count,av_max_income,av_id,quality_value,defeat_num,title,share_items,first_time,max_time,signed_at,end_at FROM up_bill WHERE mid = ?"
)

// GetUpBill get up bill by mid
func (d *Dao) GetUpBill(c context.Context, mid int64) (up *model.UpBill, err error) {
	up = &model.UpBill{}
	err = d.db.QueryRow(c, _upBillSQL, mid).Scan(&up.MID, &up.FirstIncome, &up.MaxIncome, &up.TotalIncome, &up.AvCount, &up.AvMaxIncome, &up.AvID, &up.QualityValue, &up.DefeatNum, &up.Title, &up.ShareItems, &up.FirstTime, &up.MaxTime, &up.SignedAt, &up.EndAt)
	return
}
