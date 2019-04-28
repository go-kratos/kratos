package dao

import (
	"testing"

	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"time"
)

func TestGetTable(t *testing.T) {
	Convey("Get Table", t, func() {
		once.Do(startService)

		t := getCoinStreamTable("89689d8c2290d93183b96c7ec83ea3c61114")
		So(t, ShouldEqual, "t_coin_stream_201704")

		t = getCoinStreamTable("31ea4d87162d3ed9c1dfb804eb960c640119")
		So(t, ShouldEqual, "t_coin_stream_201709")

		t = getCoinStreamTable("9937727a8c7fed5bbecd93642b4982561121")
		So(t, ShouldEqual, "t_coin_stream_201801")

		tid := GetTid(1, map[string]string{"a": "a", "b": "b"})
		t = getCoinStreamTable(tid)
		now := time.Now()
		So(t, ShouldEqual, fmt.Sprintf("t_coin_stream_%d%02d", now.Year(), int(now.Month())))

	})
}

func TestDao_NewCoinStreamRecord(t *testing.T) {
	Convey("new coin stream", t, func() {

		once.Do(startService)

		r := model.CoinStreamRecord{
			Uid:           1,
			TransactionId: GetTid(1, map[string]string{"a": "a", "b": "b"}),
			ExtendTid:     "abcdef",
			DeltaCoinNum:  1,
			CoinType:      0,
			OrgCoinNum:    100,
			OpResult:      2,
			OpReason:      0,
			OpType:        1,
			OpTime:        time.Now().Unix(),
			BizCode:       "wtf",
			Area:          0,
			Source:        "",
			MetaData:      "",
			BizSource:     "",
		}

		affect, err := d.NewCoinStreamRecord(ctx, &r)
		So(err, ShouldBeNil)
		So(affect, ShouldEqual, 1)

	})
}
