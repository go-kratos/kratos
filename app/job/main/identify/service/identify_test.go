package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServie_Ping(t *testing.T) {
	Convey("test identify consumer proc", t, func() {
		once.Do(startService)

		s.Ping(context.Background())
	})
}

func TestService_identify(t *testing.T) {
	Convey("test identify consumer proc", t, func() {
		once.Do(startService)
		var (
			bmsg interface{}
			err  error
			res  int
		)
		msg := <-s.identifySub.Messages()
		bmsg, err = s.identifyNew(msg)
		So(err, ShouldBeNil)
		So(bmsg, ShouldNotBeNil)
		res = s.identifySplit(msg, bmsg)
		So(res, ShouldNotBeNil)

		bmsgs := []interface{}{bmsg}
		s.processIdentifyInfo(bmsgs)
	})
}

func TestService_auth(t *testing.T) {
	Convey("test identify consumer proc", t, func() {
		once.Do(startService)
		var (
			m   interface{}
			err error
			res int
		)
		msg := <-s.authDataBus.Messages()
		m, err = s.new(msg)
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)
		res = s.spilt(msg, m)
		So(res, ShouldNotBeNil)

		ms := []interface{}{m}
		s.processIdentifyInfo(ms)
	})
}
