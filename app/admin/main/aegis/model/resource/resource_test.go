package resource

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	s = &Result{
		AttributeList: map[string]int8{"one": 0},
	}
	ss = &Res{
		AttributeList: map[string]int8{"one": 0},
	}
)

func TestRes_AttrParse(t *testing.T) {
	ss.AttrParse(nil)
}

func TestRes_AttrSet(t *testing.T) {
	ss.MetaParse()
}

func TestResult_AttrParse(t *testing.T) {
	s.AttrParse(nil)
}

func TestResult_AttrSet(t *testing.T) {
	s.AttrSet(map[string]uint{"one": 1})
	convey.Convey("Result_AttrSet", t, func() {
		convey.So(s.AttributeList["one"], convey.ShouldEqual, 1)
	})
}
