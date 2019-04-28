package service

import (
	"context"
	"time"

	"go-common/app/job/main/ugcpay/conf"
	"go-common/library/log"
)

func dayRange(offset int) (from, to time.Time) {
	tmp := time.Now().AddDate(0, 0, offset)
	from = time.Date(tmp.Year(), tmp.Month(), tmp.Day(), 0, 0, 0, 0, time.Local)
	to = from.Add(24*time.Hour - 1)
	return
}

func monthRange(offset int) (from, to time.Time) {
	tmp := time.Now().AddDate(0, offset, 0)
	from = time.Date(tmp.Year(), tmp.Month(), 1, 0, 0, 0, 0, time.Local)
	to = from.AddDate(0, 1, 0).Add(-1)
	return
}

func dailyBillVer(t time.Time) int64 {
	// 2006-01-02 15:04:05
	return int64(t.Year()*10000 + int(t.Month())*100 + t.Day())
}

func monthlyBillVer(t time.Time) int64 {
	return int64(t.Year()*100 + int(t.Month()))
}

func runCAS(ctx context.Context, fn func(ctx context.Context) (effected bool, err error)) (err error) {
	times := conf.Conf.Biz.RunCASTimes
	if times <= 0 {
		times = 2
	}
	effected := false
	for times > 0 {
		times--
		if effected, err = fn(ctx); err != nil {
			return
		}
		if effected {
			return
		}
	}
	if times <= 0 {
		log.Error("runCAS failed!!!")
	}
	return
}

func calcAssetIncome(fee int64) (userIncome int64, bizIncome int64) {
	if fee <= 0 {
		return 0, 0
	}
	userIncome = int64((1.0 - conf.Conf.Biz.Tax.AssetRate) * float64(fee))
	bizIncome = fee - userIncome
	return
}
