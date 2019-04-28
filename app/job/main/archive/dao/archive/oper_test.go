package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_PassedOper(t *testing.T) {
	Convey("PassedOper", t, func() {
		_, err := d.PassedOper(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}
