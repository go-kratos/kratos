package service

import (
	"strings"
	"time"
)

//timeToDate convert time to date
func (s *Service) timeToDate(ts string) (d string, err error) {
	var t time.Time
	if t, err = s.stringToTime(ts); err != nil {
		return
	}
	d = t.Format("2006/01/02")
	return
}

//stringToTime convert string to time
func (s *Service) stringToTime(ts string) (t time.Time, err error) {
	timeLayout := "2006-01-02 15:04:05"
	var loc *time.Location
	if loc, err = time.LoadLocation("Local"); err != nil {
		return
	}
	return time.ParseInLocation(timeLayout, ts, loc)
}

//weekendDays calculate weekend days between two time
func (s *Service) weekendDays(t1, t2 *time.Time) (days int) {
	var t *time.Time
	if t1.After(*t2) {
		t = t1
		t1 = t2
		t2 = t
	}
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	if y1 == y2 && m1 == m2 && d1 == d2 {
		return 0
	}

	gaps := t2.Sub(*t1).Hours() / 24
	days = int(gaps) / 7 * 2

	wd1 := int(t1.Weekday())
	wd2 := int(t2.Weekday())
	h1 := t1.Hour()
	mu1 := t1.Minute()
	s1 := t1.Second()
	h2 := t2.Hour()
	mu2 := t2.Minute()
	s2 := t2.Second()

	if wd2 < wd1 || ((wd2 == wd1) && (h2 < h1 || (h2 == h1 && mu2 < mu1) || (h2 == h1 && mu2 == mu1 && s2 < s1))) {
		days += 2
		if wd1 == 6 {
			days -= 1
		}
		if wd2 == 0 {
			days -= 1
		}
	}

	var temp1, temp2, temp3, temp4, temp5, temp6 time.Time
	temp1, _ = s.stringToTime("2018-11-05 00:00:00")
	temp2, _ = s.stringToTime("2018-11-06 23:59:59")

	temp3, _ = s.stringToTime("2018-11-03 00:00:00")
	temp4, _ = s.stringToTime("2018-11-03 23:59:59")
	temp5, _ = s.stringToTime("2018-11-11 00:00:00")
	temp6, _ = s.stringToTime("2018-11-11 23:59:59")

	if t1.Before(temp1) && t2.After(temp2) {
		days += 2
	}

	if t1.Before(temp3) && t2.After(temp4) {
		days -= 1
	}

	if t1.Before(temp5) && t2.After(temp6) {
		days -= 1
	}

	return
}

//removeCRCFInString replace CRCF with space in string
func (s *Service) removeCRCFInString(str string) (rstr string) {
	rstr = strings.Replace(str, "\n", "", -1)
	return strings.Replace(rstr, "\r", "", -1)
}
