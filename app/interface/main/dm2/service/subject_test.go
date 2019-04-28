package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubjectInfos(t *testing.T) {
	var (
		c          = context.TODO()
		tp   int32 = 1
		oids       = []int64{1221, 1231, 2386052}
	)
	Convey("get dm subject info", t, func() {
		res, err := svr.SubjectInfos(c, tp, model.MaskPlatMbl, oids)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		for _, hit := range res {
			t.Logf(":%+v", hit)
		}
	})
}

func TestSubjects(t *testing.T) {
	var (
		c          = context.TODO()
		tp   int32 = 1
		oids       = []int64{1221, 1231, 49859595}
	)
	Convey("get dm subject info", t, func() {
		res, err := svr.subjects(c, tp, oids)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		for _, hit := range res {
			t.Logf(":%+v", hit)
		}
	})
}
