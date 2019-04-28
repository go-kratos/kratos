package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_TypeMapping(t *testing.T) {
	Convey("TypeMapping", t, func() {
		_, err := d.TypeMapping(context.TODO())
		So(err, ShouldBeNil)
	})
}
