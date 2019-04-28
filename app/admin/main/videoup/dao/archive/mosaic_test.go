package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Mosaic(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		m, err := d.Mosaic(context.Background(), 10116994)
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)
	}))
}
