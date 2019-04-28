package dsn

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseDSN(t *testing.T) {
	Convey("test parsedsn", t, func() {
		var (
			s = "key:secret@group/topic=1&role=2&color=red"
			t = &DSN{
				Key:    "key",
				Secret: "secret",
				Group:  "group",
				Topic:  "1",
				Role:   "2",
				Color:  "red",
			}
		)
		d, err := ParseDSN(s)
		So(err, ShouldBeNil)
		So(d, ShouldResemble, t)
		s = "key:secret@group/top:ic=1&role=2"
		_, err = ParseDSN(s)
		So(err, ShouldNotBeNil)
	})
}
