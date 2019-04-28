package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDao_WeightLog(t *testing.T) {
	Convey("WeightLog", t, WithDao(func(d *Dao) {
		d.WeightLog(context.TODO(), 4575)
	}))
}
