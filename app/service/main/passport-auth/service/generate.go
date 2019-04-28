package service

import (
	"strconv"
	"time"
)

func calcMonDelta(ak string, now time.Time) (res int, err error) {
	var parsedMon int
	if parsedMon, err = parseMonth(ak); err != nil {
		return
	}

	curMon := int(now.Month())
	// check if cur mon
	if curMon == parsedMon {
		return 0, nil
	}

	delta := curMon - parsedMon
	if delta < 0 {
		delta += 12
	}
	return delta, nil
}

func monDiff(t time.Time, delta int) time.Time {
	if delta == 0 {
		return t
	}
	year, month, _ := t.Date()
	thisMonthFirstDay := time.Date(year, month, 1, 1, 1, 1, 1, t.Location())
	return thisMonthFirstDay.AddDate(0, delta, 0)
}

func parseMonth(ak string) (int, error) {
	n, err := strconv.ParseInt(ak[30:31], 16, 64)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}
