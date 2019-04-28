package service

import (
	"context"
	"go-common/app/admin/main/apm/model/ut"
	"io/ioutil"

	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/smartystreets/goconvey/convey"
// )

// const content = `panic: runtime error: invalid memory address or nil pointer dereference
// [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x1565c97]

// goroutine 1 [running]:
// go-common/app/interface/main/answer/dao/account.New(0x1a972e0, 0xc42017f4d0)
// 	/Users/bilibili/go/src/go-common/app/interface/main/answer/dao/account/pendant.go:36 +0x37
// go-common/app/interface/main/answer/dao/account.init.0()
// 	/Users/bilibili/go/src/go-common/app/interface/main/answer/dao/account/pendant_test.go:21 +0x68
// FAIL	go-common/app/interface/main/answer/dao/account	0.031s
// panic: runtime error: invalid memory address or nil pointer dereference
// [signal SIGSEGV: segmentation violation code=0x1 addr=0x50 pc=0x1439ed6]
// === RUN   TestDao_MoralLog
// >->->OPEN-JSON->->->
// {
//   "Title": "MoralLog",
//   "File": "/Users/bilibili/go/src/go-common/app/service/main/member/dao/hbase_test.go",
//   "Line": 11,
//   "Depth": 1,
//   "Assertions": [
//     {
//       "File": "/Users/bilibili/go/src/go-common/app/service/main/member/dao/hbase_test.go",
//       "Line": 12,
//       "Expected": "",
//       "Actual": "",
//       "Failure": "",
//       "Error": "runtime error: invalid memory address or nil pointer dereference",",
//       "Skipped": false
//     }
//   ],
//   "Output": ""
// },
// <-<-<-CLOSE-JSON<-<-<
// --- FAIL: TestDao_MoralLog (0.00s)

// `

// func TestService_Upload(t *testing.T) {
// 	convey.Convey("ParserContent", t, func() {
// 		data, err := svr.ParseContent(context.Background(), []byte(content))
// 		convey.So(err, convey.ShouldBeNil)
// 		convey.So(data, convey.ShouldNotBeNil)
// 		t.Logf("after parsercontent: %s", string(data))
// 		convey.So(err, convey.ShouldBeNil)
// 		info, err := svr.CalcCount(context.Background(), data)
// 		convey.So(err, convey.ShouldBeNil)
// 		t.Logf("pass: %d", info.Passed)
// 		t.Logf("fail: %d", info.Failures)
// 		t.Logf("skip: %d", info.Skipped)
// 		t.Logf("panics: %d", info.Panics)
// 		t.Logf("total: %d", info.Assertions)
// 		t.Logf("coverage: %s", info.Coverage)

// 		convey.Convey("Upload", func() {
// 			var (
// 				body = data
// 			)
// 			url, err := svr.Upload(context.Background(), "json", time.Now().Unix(), body)
// 			convey.So(err, convey.ShouldBeNil)
// 			convey.So(url, convey.ShouldNotBeNil)
// 			t.Logf("Location: %s", url)
// 		})
// 	})
// }

func TestServiceCalcCountFiles(t *testing.T) {
	convey.Convey("CalcCountFiles", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			res = &ut.UploadRes{
				CommitID: "somestringhasnothingtodo",
				PKG:      "go-common/app/admin/main/apm/dao",
			}
			filename = "/data/ut1/cover.out"
		)
		body, _ := ioutil.ReadFile(filename)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			utfiles, err := svr.CalcCountFiles(c, res, body)
			t.Logf("\nutfiles:%#v\n", utfiles)
			for i, utfile := range utfiles {
				t.Logf("\nutfiles[%d]:%#v\n", i, utfile)
			}
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
