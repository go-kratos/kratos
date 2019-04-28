package dao

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetnxMaskJob(t *testing.T) {
	Convey("", t, func() {
		ok, err := testDao.SetnxMaskJob(context.Background(), time.Now().String())
		So(err, ShouldBeNil)
		t.Logf("ok:%v", ok)
	})
}

func TestGetMaskJob(t *testing.T) {
	Convey("", t, func() {
		value, err := testDao.GetMaskJob(context.Background())
		So(err, ShouldBeNil)
		t.Logf("ok:%v", value)
	})
}

func TestGetSetMaskJob(t *testing.T) {
	Convey("", t, func() {
		value, err := testDao.GetSetMaskJob(context.Background(), time.Now().String())
		So(err, ShouldBeNil)
		t.Logf("ok:%v", value)
	})
}

func TestDelMaskJob(t *testing.T) {
	Convey("", t, func() {
		err := testDao.DelMaskJob(context.Background())
		So(err, ShouldBeNil)
	})
}
