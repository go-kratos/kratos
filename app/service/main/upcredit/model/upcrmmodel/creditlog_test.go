package upcrmmodel

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/library/time"
	"testing"
	xtime "time"
)

func Test_Orm(t *testing.T) {
	var (
		arg = ArgCreditLogAdd{
			BusinessType: 1,
			Type:         1,
			OpType:       2,
			Reason:       101,
			Mid:          12345,
			Oid:          13021,
			UID:          1,
			Content:      "稿件打回",
			CTime:        time.Time(xtime.Now().Unix()),
			Extra:        []byte("{ \"key\" : \"value\"}"),
		}
	)
	Convey("orm", t, func() {
		Convey("connect", func() {
			var js, _ = json.Marshal(arg)
			t.Log(string(js))
		})
	})
}
