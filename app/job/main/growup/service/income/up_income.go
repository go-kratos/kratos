package income

import (
	"bytes"
	"context"
	"strconv"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
)

// UpIncomeSvr up income service
type UpIncomeSvr struct {
	batchSize int
	dao       *incomeD.Dao
}

// NewUpIncomeSvr new income service
func NewUpIncomeSvr(dao *incomeD.Dao, batchSize int) (svr *UpIncomeSvr) {
	return &UpIncomeSvr{
		batchSize: batchSize,
		dao:       dao,
	}
}

func (p *UpIncomeSvr) getUpIncomeByDate(c context.Context, upCh chan []*model.UpIncome, date string, limit int) (err error) {
	defer close(upCh)
	var id int64
	for {
		var up []*model.UpIncome
		up, err = p.dao.GetUpIncomeTable(c, "up_income", date, id, limit)
		if err != nil {
			return
		}
		upCh <- up
		if len(up) < limit {
			break
		}
		id = up[len(up)-1].ID
	}
	return
}

// BatchInsertUpIncome insert up_income batch
func (p *UpIncomeSvr) BatchInsertUpIncome(c context.Context, us map[int64]*model.UpIncome) (err error) {
	var (
		buff    = make([]*model.UpIncome, p.batchSize)
		buffEnd = 0
	)

	for _, u := range us {
		buff[buffEnd] = u
		buffEnd++

		if buffEnd >= p.batchSize {
			values := upIncomeValues(buff[:buffEnd])
			buffEnd = 0
			_, err = p.dao.InsertUpIncome(c, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := upIncomeValues(buff[:buffEnd])
		buffEnd = 0
		_, err = p.dao.InsertUpIncome(c, values)
	}
	return
}

func upIncomeValues(us []*model.UpIncome) (values string) {
	var buf bytes.Buffer
	for _, u := range us {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(u.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.PlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AudioIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.TaxMoney, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.Income, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.TotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvBaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvTax, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnBaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnTax, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmBaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmTax, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + u.Date.Time().Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BaseIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.AvTotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.ColumnTotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.BgmTotalIncome, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}

///************************************************************ FOR HISTORY ****************************************************************/
//
//// UpdateBusinessIncomeByDate update business income by date
//func (p *UpIncomeSvr) UpdateBusinessIncomeByDate(c context.Context, date string, ustat map[int64]*model.UpIncomeStat) (err error) {
//	us, err := p.businessTotalIncome(c, ustat, date)
//	if err != nil {
//		return
//	}
//	err = p.batchUpdateUpIncome(c, us)
//	if err != nil {
//		return
//	}
//	return p.batchUpdateUpIncomeStat(c, ustat)
//}
//
//// for history data m: map[mid]map[date]*model.UpIncome
//func (p *UpIncomeSvr) businessTotalIncome(c context.Context, ustat map[int64]*model.UpIncomeStat, date string) (m []*model.UpIncome, err error) {
//	var id int64
//	for {
//		var ups []*model.UpIncome
//		ups, err = p.dao.GetUpIncomeTable(c, "up_income", date, id, 2000)
//		if err != nil {
//			return
//		}
//		for _, up := range ups {
//			ut := ustat[up.MID]
//
//			ut.AvTotalIncome += up.AvIncome
//			up.AvTotalIncome = ut.AvTotalIncome
//
//			ut.ColumnTotalIncome += up.ColumnIncome
//			up.ColumnTotalIncome = ut.ColumnTotalIncome
//
//			ut.BgmTotalIncome += up.BgmIncome
//			up.BgmTotalIncome = ut.BgmTotalIncome
//			m = append(m, up)
//		}
//		if len(ups) < 2000 {
//			break
//		}
//		id = ups[len(ups)-1].ID
//	}
//	return
//}
//
//// BatchUpdateUpIncome insert up_income batch
//func (p *UpIncomeSvr) batchUpdateUpIncome(c context.Context, us []*model.UpIncome) (err error) {
//	var (
//		buff    = make([]*model.UpIncome, p.batchSize)
//		buffEnd = 0
//	)
//
//	for _, u := range us {
//		buff[buffEnd] = u
//		buffEnd++
//
//		if buffEnd >= p.batchSize {
//			values := businessValues(buff[:buffEnd])
//			buffEnd = 0
//			_, err = p.dao.FixInsertUpIncome(c, values)
//			if err != nil {
//				return
//			}
//		}
//	}
//	if buffEnd > 0 {
//		values := businessValues(buff[:buffEnd])
//		buffEnd = 0
//		_, err = p.dao.FixInsertUpIncome(c, values)
//	}
//	return
//}
//
//func businessValues(us []*model.UpIncome) (values string) {
//	var buf bytes.Buffer
//	for _, u := range us {
//		buf.WriteString("(")
//		buf.WriteString(strconv.FormatInt(u.MID, 10))
//		buf.WriteByte(',')
//		buf.WriteString(strconv.FormatInt(u.AvTotalIncome, 10))
//		buf.WriteByte(',')
//		buf.WriteString(strconv.FormatInt(u.ColumnTotalIncome, 10))
//		buf.WriteByte(',')
//		buf.WriteString(strconv.FormatInt(u.BgmTotalIncome, 10))
//		buf.WriteByte(',')
//		buf.WriteString("'" + u.Date.Time().Format(_layout) + "'")
//		buf.WriteString(")")
//		buf.WriteByte(',')
//	}
//	if buf.Len() > 0 {
//		buf.Truncate(buf.Len() - 1)
//	}
//	values = buf.String()
//	buf.Reset()
//	return
//}
//
//// BatchUpdateUpIncome insert up_income batch
//func (p *UpIncomeSvr) batchUpdateUpIncomeStat(c context.Context, us map[int64]*model.UpIncomeStat) (err error) {
//	var (
//		buff    = make([]*model.UpIncomeStat, p.batchSize)
//		buffEnd = 0
//	)
//
//	for _, u := range us {
//		buff[buffEnd] = u
//		buffEnd++
//
//		if buffEnd >= p.batchSize {
//			values := businessStatValues(buff[:buffEnd])
//			buffEnd = 0
//			_, err = p.dao.FixInsertUpIncomeStat(c, values)
//			if err != nil {
//				return
//			}
//		}
//	}
//	if buffEnd > 0 {
//		values := businessStatValues(buff[:buffEnd])
//		buffEnd = 0
//		_, err = p.dao.FixInsertUpIncomeStat(c, values)
//	}
//	return
//}
//
//func businessStatValues(ustat []*model.UpIncomeStat) (values string) {
//	var buf bytes.Buffer
//	for _, u := range ustat {
//		buf.WriteString("(")
//		buf.WriteString(strconv.FormatInt(u.MID, 10))
//		buf.WriteByte(',')
//		buf.WriteString(strconv.FormatInt(u.AvTotalIncome, 10))
//		buf.WriteByte(',')
//		buf.WriteString(strconv.FormatInt(u.ColumnTotalIncome, 10))
//		buf.WriteByte(',')
//		buf.WriteString(strconv.FormatInt(u.BgmTotalIncome, 10))
//		buf.WriteString(")")
//		buf.WriteByte(',')
//	}
//	if buf.Len() > 0 {
//		buf.Truncate(buf.Len() - 1)
//	}
//	values = buf.String()
//	buf.Reset()
//	return
//}
