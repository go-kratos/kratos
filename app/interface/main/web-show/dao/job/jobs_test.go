package job

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_Job(t *testing.T) {
	Convey("test job", t, WithDao(func(d *Dao) {
		data, err := d.Jobs(context.TODO())
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}
