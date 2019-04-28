package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_hbaseMd5Key(t *testing.T) {
	mid := int64(908085)
	convey.Convey("test article stat", t, func() {
		res := hbaseMd5Key(mid)
		convey.Printf("%s", string(res))
	})
}
