package service

import (
	"context"
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestService_SDConvert(t *testing.T) {
	Convey("test convert sd", t, func() {
		Convey("new sd encode and decode", func() {
			// decode
			sd := "25ded96af42eb61677730d0a74eb4c51"
			sdb, err := decodeSD(sd)
			So(err, ShouldBeNil)
			So(len(sdb), ShouldEqual, _newTokenBinByteLen)

			// encode
			s := encodeSD(sdb)
			So(s, ShouldEqual, sd)
		})

		Convey("old sd encode and decode", func() {
			// decode
			sd := "396c38bb,1519380539,d73804f2"
			sdb, err := decodeSD(sd)
			So(err, ShouldBeNil)
			So(len(sdb), ShouldEqual, 16+2+10)

			// encode
			s := encodeSD(sdb)
			So(s, ShouldEqual, sd)
		})
	})
}

func TestService_CookieInfo(t *testing.T) {
	once.Do(startService)
	Convey("Test Query Token", t, func() {
		sd := "396c38bb,1519380539,d73804f2"
		res, err := s.cookieInfo(context.TODO(), sd)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		So(res.Session, ShouldEqual, sd)

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}
