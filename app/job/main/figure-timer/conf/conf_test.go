package conf

import (
	"flag"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	flag.Set("conf", "../cmd/figure-timer-job-test.toml")
}

func TestInit(t *testing.T) {
	Convey("TEST conf", t, func() {
		err := Init()
		So(err, ShouldBeNil)
		So(Conf.Property, ShouldNotBeNil)
		So(Conf.Property.Calc, ShouldNotBeNil)
		fmt.Printf("%+v", Conf.Property.Calc)
	})
}
