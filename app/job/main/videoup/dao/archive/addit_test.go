package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Addit(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("Test_Addit", t, WithDao(func(d *Dao) {
		_, err = d.Addit(c, 23333)
		So(err, ShouldBeNil)
	}))
}
