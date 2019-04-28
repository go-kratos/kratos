package income

import (
	"time"
)

var (
	_layout      = "2006-01-02"
	_layoutMonth = "2006-01"
)

const (
	_upIncome        = "up_income"
	_upIncomeWeekly  = "up_income_weekly"
	_upIncomeMonthly = "up_income_monthly"

	_avDailyStatis   = "av_income_daily_statis"
	_avWeeklyStatis  = "av_income_weekly_statis"
	_avMonthlyStatis = "av_income_monthly_statis"

	_cmDailyStatis   = "column_income_daily_statis"
	_cmWeeklyStatis  = "column_income_weekly_statis"
	_cmMonthlyStatis = "column_income_monthly_statis"

	_bgmDailyStatis   = "bgm_income_daily_statis"
	_bgmWeeklyStatis  = "bgm_income_weekly_statis"
	_bgmMonthlyStatis = "bgm_income_monthly_statis"

	_avChargeDailyStatis   = "av_charge_daily_statis"
	_avChargeWeeklyStatis  = "av_charge_weekly_statis"
	_avChargeMonthlyStatis = "av_charge_monthly_statis"

	_cmChargeDailyStatis   = "column_charge_daily_statis"
	_cmChargeWeeklyStatis  = "column_charge_weekly_statis"
	_cmChargeMonthlyStatis = "column_charge_monthly_statis"

	_bgmChargeDailyStatis   = "bgm_charge_daily_statis"
	_bgmChargeWeeklyStatis  = "bgm_charge_weekly_statis"
	_bgmChargeMonthlyStatis = "bgm_charge_monthly_statis"

	_upIncomeDailyStatis = "up_income_daily_statis"
	_upAvDailyStatis     = "up_av_daily_statis"
	_upColumnDailyStatis = "up_column_daily_statis"
	_upBgmDailyStatis    = "up_bgm_daily_statis"

	_groupDay   = 1
	_groupWeek  = 2
	_groupMonth = 3

	_video = 0
	// _audio  = 1
	_column  = 2
	_bgm     = 3
	_up      = 4
	_lottery = 5 // 动态抽奖

	_leastAvIncome = 30000

	// add blacklist reason
	_avBlack  = 1
	_avBreach = 2
)

func getUpInfoTable(typ int) (table string) {
	switch typ {
	case _video:
		table = "up_info_video"
	case _column:
		table = "up_info_column"
	case _bgm:
		table = "up_info_bgm"
	}
	return
}

func getUpInfoByType(typ int) (table string, incomeType string) {
	switch typ {
	case _video:
		table, incomeType = _upAvDailyStatis, "av_income"
	case _column:
		table, incomeType = _upColumnDailyStatis, "column_income"
	case _bgm:
		table, incomeType = _upBgmDailyStatis, "bgm_income"
	case _up:
		table, incomeType = _upIncomeDailyStatis, "income"
	}
	return
}

func getUpFieldByType(typ int) (field string) {
	switch typ {
	case _video:
		field = "av_income,av_tax,av_base_income,av_total_income"
	case _column:
		field = "column_income,column_tax,column_base_income,column_total_income"
	case _bgm:
		field = "bgm_income,bgm_tax,bgm_base_income,bgm_total_income"
	case _up:
		field = "income,tax_money,base_income,total_income"
	}
	return
}

func setUpTableByGroup(groupType int) string {
	table := _upIncome
	if groupType == _groupWeek {
		table = _upIncomeWeekly
	} else if groupType == _groupMonth {
		table = _upIncomeMonthly
	}
	return table
}

func setArchiveTableByGroup(typ, groupType int) string {
	table := ""
	switch typ {
	case _video:
		table = _avDailyStatis
		if groupType == _groupWeek {
			table = _avWeeklyStatis
		} else if groupType == _groupMonth {
			table = _avMonthlyStatis
		}
	case _column:
		table = _cmDailyStatis
		if groupType == _groupWeek {
			table = _cmWeeklyStatis
		} else if groupType == _groupMonth {
			table = _cmMonthlyStatis
		}
	case _bgm:
		table = _bgmDailyStatis
		if groupType == _groupWeek {
			table = _bgmWeeklyStatis
		} else if groupType == _groupMonth {
			table = _bgmMonthlyStatis
		}
	}
	return table
}

func setChargeTableByGroup(typ, groupType int) string {
	table := ""
	switch typ {
	case _video:
		table = _avChargeDailyStatis
		if groupType == _groupWeek {
			table = _avChargeWeeklyStatis
		} else if groupType == _groupMonth {
			table = _avChargeMonthlyStatis
		}
	case _column:
		table = _cmChargeDailyStatis
		if groupType == _groupWeek {
			table = _cmChargeWeeklyStatis
		} else if groupType == _groupMonth {
			table = _cmChargeMonthlyStatis
		}
	case _bgm:
		table = _bgmChargeDailyStatis
		if groupType == _groupWeek {
			table = _bgmChargeWeeklyStatis
		} else if groupType == _groupMonth {
			table = _bgmChargeMonthlyStatis
		}
	}
	return table
}

func getDateByGroup(groupType int, date time.Time) time.Time {
	if groupType == _groupWeek {
		return getStartWeekDate(date)
	} else if groupType == _groupMonth {
		return getStartMonthDate(date)
	}
	return date
}

func addDayByGroup(groupType int, date time.Time) time.Time {
	if groupType == _groupWeek {
		return date.AddDate(0, 0, 7)
	} else if groupType == _groupMonth {
		return date.AddDate(0, 1, 0)
	}
	return date.AddDate(0, 0, 1)
}

func getStartWeekDate(date time.Time) time.Time {
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	return date
}

func getStartMonthDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
}

func fromYuanToFen(d int64) float64 {
	return float64(d) / float64(100)
}
