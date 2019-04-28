package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPubDatabus(t *testing.T) {
	var (
		tp       int32 = 1
		aid      int64 = 10097265
		oid      int64 = 1508
		cnt      int64 = 26
		num      int64 = 1
		duration int64 = 9031000
		c              = context.TODO()
	)
	Convey("flush segment dm xml", t, func() {
		err := testDao.PubDatabus(c, tp, aid, oid, cnt, num, duration)
		So(err, ShouldBeNil)
	})
}

func TestSendAction(t *testing.T) {
	var (
		c     = context.TODO()
		flush = &model.Flush{Oid: 1221, Type: 1}
	)
	Convey("flush xml", t, func() {
		data, err := json.Marshal(flush)
		So(err, ShouldBeNil)
		act := &model.Action{Action: model.ActionFlush, Data: data}
		err = testDao.SendAction(c, fmt.Sprint(flush.Oid), act)
		So(err, ShouldBeNil)
	})
}
