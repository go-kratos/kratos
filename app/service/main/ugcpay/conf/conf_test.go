package conf

import (
	"flag"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConf(t *testing.T) {
	Convey("", t, func() {
		flag.Set("conf", "../cmd/test.toml")
		err := Init()
		So(err, ShouldBeNil)
	})
}
