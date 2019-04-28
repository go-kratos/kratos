package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_config(t *testing.T) {
	initd()
	Convey("Test config get", t, func() {
		cf, err := d.GetPushConfig(context.TODO())
		t.Logf("the result included(%v) err(%v)", cf, err)
		So(err, ShouldEqual, nil)
	})
}

func Test_GetPushInterval(t *testing.T) {
	initd()
	Convey("Test GetPushInterval", t, func() {
		cf, err := d.GetPushInterval(context.TODO())
		t.Logf("the result included(%v) err(%v)", cf, err)
		So(err, ShouldEqual, nil)
	})
}
