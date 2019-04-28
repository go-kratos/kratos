package charge

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	model "go-common/app/job/main/growup/model/charge"

	xtime "go-common/library/time"
)

func (s *Service) calUpCharge(c context.Context, t time.Time, avc chan []*model.AvCharge) (daily, weekly, monthly map[int64]*model.UpCharge, err error) {
	od, err := s.getUpDailyCharge(c, t)
	if err != nil {
		return
	}
	daily = calUpDailyCharge(od, avc)

	weekly, err = s.getUpWeeklyCharge(c)
	if err != nil {
		return
	}
	// group by weekly
	upChargeStat(startWeeklyDate, weekly, daily)

	monthly, err = s.getUpMonthlyCharge(c)
	if err != nil {
		return
	}
	// group by monthly
	upChargeStat(startMonthlyDate, monthly, daily)

	return
}

// 1 od: old daily up charges
func (s *Service) getUpDailyCharge(c context.Context, t time.Time) (od map[int64]*model.UpCharge, err error) {
	// old: yesterday
	old := t.AddDate(0, 0, -1)
	return s.getUpCharges(c, "up_daily_charge", old.Format("2006-01-02"))
}

// 2 ow: old weekly up charges
func (s *Service) getUpWeeklyCharge(c context.Context) (ow map[int64]*model.UpCharge, err error) {
	return s.getUpCharges(c, "up_weekly_charge", startWeeklyDate.Format("2006-01-02"))
}

// 3 om: old monthly up charges
func (s *Service) getUpMonthlyCharge(c context.Context) (om map[int64]*model.UpCharge, err error) {
	return s.getUpCharges(c, "up_monthly_charge", startMonthlyDate.Format("2006-01-02"))
}

func calUpDailyCharge(od map[int64]*model.UpCharge, avc chan []*model.AvCharge) (mu map[int64]*model.UpCharge) {
	mu = make(map[int64]*model.UpCharge)
	for charges := range avc {
		for _, charge := range charges {
			if charge.IncCharge <= 0 {
				continue
			}
			// udc: up daily charge
			if udc, ok := mu[charge.MID]; ok {
				udc.IncCharge += charge.IncCharge
				udc.TotalCharge += charge.IncCharge
			} else {
				var total int64
				if o, ok := od[charge.MID]; ok {
					total = o.TotalCharge
				}
				mu[charge.MID] = &model.UpCharge{
					MID:         charge.MID,
					IncCharge:   charge.IncCharge,
					TotalCharge: total + charge.IncCharge,
					Date:        charge.Date,
				}
			}
		}
	}
	return
}

// os: old weekly/month charge chan, empty maybe
func upChargeStat(t time.Time, os, daily map[int64]*model.UpCharge) {
	for mid, ucd := range daily {
		if charge, ok := os[mid]; ok {
			// update
			charge.IncCharge += ucd.IncCharge
			charge.TotalCharge += ucd.IncCharge
			charge.Date = xtime.Time(t.Unix())
		} else {
			// new
			os[mid] = &model.UpCharge{
				MID:         mid,
				IncCharge:   ucd.IncCharge,
				TotalCharge: ucd.TotalCharge,
				Date:        xtime.Time(t.Unix()),
			}
		}
	}
}

// get up charges by date
func (s *Service) getUpCharges(c context.Context, table, date string) (m map[int64]*model.UpCharge, err error) {
	var id int64
	m = make(map[int64]*model.UpCharge)
	for {
		var charges map[int64]*model.UpCharge
		id, charges, err = s.dao.GetUpCharges(c, table, date, id, 2000)
		if err != nil {
			return
		}
		if len(charges) == 0 {
			break
		}
		for k, v := range charges {
			m[k] = v
		}
	}
	return
}

// BatchInsertUpCharge batch insert up charge
func (s *Service) BatchInsertUpCharge(c context.Context, table string, us map[int64]*model.UpCharge) (err error) {
	var (
		buff    = make([]*model.UpCharge, _batchSize)
		buffEnd = 0
	)
	for _, u := range us {
		buff[buffEnd] = u
		buffEnd++
		if buffEnd >= _batchSize {
			values := upChargeValues(buff[:buffEnd])
			buffEnd = 0
			_, err = s.dao.InsertUpCharge(c, table, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := upChargeValues(buff[:buffEnd])
		buffEnd = 0
		_, err = s.dao.InsertUpCharge(c, table, values)
	}
	return
}

func upChargeValues(us []*model.UpCharge) (values string) {
	var buf bytes.Buffer
	for _, charge := range us {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(charge.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(charge.IncCharge, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(charge.TotalCharge, 10))
		buf.WriteByte(',')
		buf.WriteString(fmt.Sprintf("'%s'", charge.Date.Time().Format("2006-01-02")))
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
