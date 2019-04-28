package client

import (
	"context"
	"os"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

var svr *Service

func TestMain(m *testing.M) {
	// new rpc client
	svr = New(nil)
	os.Exit(m.Run())
}

func TestSubjectInfos(t *testing.T) {
	var (
		tp   int32 = 1
		oids       = []int64{1221, 1491, 1352, 1391, 1291, 10109227}
	)
	Convey("test dm subject info", t, func() {
		arg := &model.ArgOids{Type: tp, Oids: oids, Plat: 1}
		res, err := svr.SubjectInfos(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		for oid, v := range res {
			t.Logf("oid:%d,%+v", oid, v)
		}
	})
}
