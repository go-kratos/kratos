package space

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DelContributeIDCache(t *testing.T) {
	Convey("DelContributeIDCache", t, func() {
		d.DelContributeIDCache(context.TODO(), 1, 1, "")
	})
}

func Test_DelContributeCache(t *testing.T) {
	Convey("DelContributeCache", t, func() {
		d.DelContributeCache(context.TODO(), 1)
	})
}
