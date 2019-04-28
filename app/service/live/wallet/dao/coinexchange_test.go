package dao

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"testing"
	"time"
)

func TestDao_NewCoinExchangeRecord(t *testing.T) {
	Convey("New CoinExchange", t, func() {
		once.Do(startService)
		record := new(model.CoinExchangeRecord)
		record.ExchangeTime = time.Now().Unix()
		record.Status = 1
		record.DestNum = 1
		record.DestType = 1
		record.SrcType = 2
		record.SrcNum = 1
		record.TransactionId = "abcdef"
		record.Uid = 1

		affect, err := d.NewCoinExchangeRecord(ctx, record)
		So(affect, ShouldEqual, 1)
		So(err, ShouldBeNil)
	})
}
