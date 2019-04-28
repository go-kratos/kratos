package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func Test_MutliSendSysMsg(t *testing.T) {
	Convey("return someting", t, func() {
		defer gock.OffAll()
		httpMock("POST", d.msgURL).Reply(200).JSON(`{"code":0}`)
		err := d.MutliSendSysMsg(context.Background(), []int64{22, 11}, "dasda", "dsadasd", "127.0.0.1")
		So(err, ShouldBeNil)
	})
}

func Test_SendSysMsg(t *testing.T) {
	Convey("return someting", t, func() {
		defer gock.OffAll()
		httpMock("POST", d.msgURL).Reply(200).JSON(`{"code":0}`)
		err := d.SendSysMsg(context.Background(), []int64{22, 11}, "dasda", "dsadasd", "127.0.0.1")
		So(err, ShouldBeNil)
	})
}
