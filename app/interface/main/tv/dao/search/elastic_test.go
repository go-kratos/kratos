package search

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	seaMdl "go-common/app/interface/main/tv/model/search"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestSearchyearTrans(t *testing.T) {
	var (
		year = "2018-2019"
	)
	convey.Convey("yearTrans", t, func(ctx convey.C) {
		stime, etime, err := yearTrans(year)
		ctx.Convey("Then err should be nil.stime,etime should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(etime, convey.ShouldNotBeNil)
			ctx.So(stime, convey.ShouldNotBeNil)
			fmt.Println(stime, etime)
		})
	})
}

func TestSearchPgcIdx(t *testing.T) {
	var (
		c     = context.Background()
		pnDoc = seaMdl.ReqEsPn{
			Ps:    50,
			Pn:    1,
			Order: 6,
			Sort:  0,
		}
		reqDoc = &seaMdl.ReqPgcIdx{
			SeasonType:    3,
			ProducerID:    31,
			Year:          "",
			StyleID:       -1,
			PubDate:       "",
			SeasonMonth:   0,
			SeasonStatus:  []int{},
			Copyright:     []int{},
			IsFinish:      "",
			Area:          []int{},
			SeasonVersion: 0,
		}
		reqJp = &seaMdl.ReqPgcIdx{
			SeasonType:    1,
			ProducerID:    -1,
			Year:          "",
			StyleID:       136,
			PubDate:       "2018-2018",
			SeasonMonth:   10,
			SeasonStatus:  []int{2, 6}, // need to pay
			Copyright:     []int{0, 1, 2, 4},
			IsFinish:      "0",
			Area:          []int{},
			SeasonVersion: 1,
		}
	)
	convey.Convey("PgcIdxJp", t, func(ctx convey.C) {
		ctx.Convey("PgcIdxJp should not be nil", func(ctx convey.C) {
			reqDoc.ReqEsPn = pnDoc
			dataDoc, err2 := d.PgcIdx(c, reqDoc)
			ctx.So(err2, convey.ShouldBeNil)
			ctx.So(dataDoc, convey.ShouldNotBeNil)
			str, _ := json.Marshal(dataDoc)
			fmt.Println(string(str))
		})
		ctx.Convey("PgcIdxDoc should not be nil", func(ctx convey.C) {
			pnDoc.Sort = 1
			reqJp.ReqEsPn = pnDoc
			dataJp, err := d.PgcIdx(c, reqJp)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(dataJp, convey.ShouldNotBeNil)
			fmt.Println(dataJp)
			str, _ := json.Marshal(dataJp)
			fmt.Println(string(str))
		})
	})
}

func TestSeachUgcIdx(t *testing.T) {
	var (
		c     = context.Background()
		pnDoc = seaMdl.ReqEsPn{
			Ps:    50,
			Pn:    1,
			Order: 6,
			Sort:  0,
		}
		reqNone = &seaMdl.SrvUgcIdx{}
		reqTIDs = &seaMdl.SrvUgcIdx{
			TIDs: []int32{75},
			PubTime: &seaMdl.UgcTime{
				STime: "2017-01-01 00:00:00",
				ETime: "2018-12-31 23:59:59",
			},
		}
	)
	convey.Convey("UgcIdx", t, func(ctx convey.C) {
		ctx.Convey("If TIDs empty, return request error ", func(ctx convey.C) {
			reqTIDs.ReqEsPn = pnDoc
			_, err := d.UgcIdx(c, reqNone)
			ctx.So(err, convey.ShouldEqual, ecode.RequestErr)
		})
		ctx.Convey("Pick Result with type_ids AND pub_time", func(ctx convey.C) {
			dataDoc, err := d.UgcIdx(c, reqTIDs)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(dataDoc, convey.ShouldNotBeNil)
			str, _ := json.Marshal(dataDoc)
			fmt.Println(string(str))
		})
	})
}
