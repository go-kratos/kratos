package feed

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Hots(t *testing.T) {
	Convey("Hots", t, func() {
		d.Hots(context.TODO())
	})
}

func Test_UpRcmdCache(t *testing.T) {
	Convey("UpRcmdCache", t, func() {
		d.UpRcmdCache(context.TODO(), nil)
	})
}
