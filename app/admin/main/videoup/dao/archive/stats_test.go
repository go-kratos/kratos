package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test_StatsPoints(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		now := time.Now()
		_, err := d.StatsPoints(context.Background(), now.Add(-time.Hour), now, 1)
		So(err, ShouldBeNil)
	}))
}
