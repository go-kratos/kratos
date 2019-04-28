package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_TagNameMap(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		_, err := d.TagNameMap(context.Background(), []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}
