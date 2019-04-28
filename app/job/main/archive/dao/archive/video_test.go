package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Videos(t *testing.T) {
	Convey("Videos", t, func() {
		_, err := d.Videos(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_Videos2(t *testing.T) {
	Convey("Videos2", t, func() {
		vs, err := d.Videos2(context.TODO(), 10098500)
		So(err, ShouldBeNil)
		for _, v := range vs {
			Printf("%+v", v)
		}
	})
}
