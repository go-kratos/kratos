package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_statsCache(t *testing.T) {
	Convey("should get statsCache", t, func() {
		d.statsCache(context.Background(), []int64{1, 2, 3, 4})
	})
}

func Test_arcsCache(t *testing.T) {
	Convey("should get arcsCache", t, func() {
		d.arcsCache(context.Background(), []int64{1, 2, 3, 4})
	})
}
