package monitor

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Send(t *testing.T) {
	Convey("Send", t, func() {
		err := d.Send(context.TODO(), "报警短信test")
		So(err, ShouldBeNil)
	})
}
