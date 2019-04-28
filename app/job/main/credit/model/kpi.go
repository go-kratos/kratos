package model

// const kpi const
const (
	// account pendPendant
	_accountPendPendantS = 140
	_accountPendPendantA = 139
	_accountPendPendantB = 138
	// _accountPendPendantC = 0

	// KPI level and pendant info
	_kpiLevelS = int8(1)
	_kpiLevelA = int8(2)
	_kpiLevelB = int8(3)
	_kpiLevelC = int8(4)

	_kpiNameplateA = 72
	_kpiNameplateB = 71
	_kpiNameplateC = 70

	_kpiRateTimesA = 12
	_kpiRateTimesB = 6
	_kpiRateTimesC = 3

	_kpiCoinsS       = float64(20)
	_kpiCoinsA       = float64(10)
	_kpiCoinsB       = float64(5)
	_kpiCoinsDefealt = float64(0)

	// kpi defealt send days
	KPIDefealtPendSendDays = 30

	KPICoinsReason = "风纪委员任期奖励"
)

// var kpi var.
var (
	// credit level mapping pendant info.
	_levelPendant = map[int8][]int64{
		_kpiLevelS: {_accountPendPendantS, _accountPendPendantA, _accountPendPendantB},
		_kpiLevelA: {_accountPendPendantA, _accountPendPendantB},
		_kpiLevelB: {_accountPendPendantB},
		_kpiLevelC: {},
	}
)

// LevelPendantByKPI get levelPendant by kpi level.
func LevelPendantByKPI(kpiLevel int8) (lps []int64, ok bool) {
	lps, ok = _levelPendant[kpiLevel]
	return
}

// KpiCoinsRate get coins by rate.
func KpiCoinsRate(rate int8) (coins float64) {
	switch rate {
	case _kpiLevelS:
		coins = _kpiCoinsS
	case _kpiLevelA:
		coins = _kpiCoinsA
	case _kpiLevelB:
		coins = _kpiCoinsB
	default:
		coins = _kpiCoinsDefealt
	}
	return
}

// KpiPlateIDRateTimes get plate_id by rate times.
func KpiPlateIDRateTimes(rateTimes int) (plateID int64) {
	switch rateTimes {
	case _kpiRateTimesA:
		plateID = _kpiNameplateA
	case _kpiRateTimesB:
		plateID = _kpiNameplateB
	case _kpiRateTimesC:
		plateID = _kpiNameplateC
	}
	return
}
