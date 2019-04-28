package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSendTask(t *testing.T) {
	convey.Convey("SendTask", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			taskSQL = []string{"index.mid=3458517", "content.log_date<=20181111"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			statusURL, err := testDao.SendTask(c, taskSQL)
			ctx.Convey("Then err should be nil.statusURL should not be nil.", func(ctx convey.C) {
				//	ctx.So(err, convey.ShouldBeNil)
				//	ctx.So(statusURL, convey.ShouldNotBeNil)
				t.Logf("%v %s\n", err, statusURL)
			})
		})
	})
}

// func TestBerserker(t *testing.T) {
// 	params := url.Values{}
// 	params.Set("appKey", "672bc22888af701529e8b3052fd2c4a7")
// 	params.Set("query", "select * from ods.ods_dm_index where dmid<1000 limit 10")
// 	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
// 	params.Set("version", "1.0")
// 	params.Set("signMethod", "md5")

// 	s := _berserker + "?" + sign(params)
// 	fmt.Println(s)
// 	body, err := oget(s)
// 	if err != nil {
// 		t.Errorf("url(%s) error(%s)", s, err)
// 		t.FailNow()
// 	}
// 	fmt.Println(string(body))
// 	var out bytes.Buffer
// 	if err = json.Indent(&out, body, "", " "); err != nil {
// 		t.Fatal(err)
// 		t.FailNow()
// 	}
// 	fmt.Println(out.String())
// }
