package model

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBatchParam_Format(t *testing.T) {
	Convey("BatchParam_Format", t, func() {
		f := func(p *BatchParam, t *testing.T, args ...interface{}) {
			t.Logf("args(%+v)", args...)
			p.Handler(&p.Params, args...)
			t.Logf("params(%+v)\r\n", p)
		}

		var uuid interface{}
		m := map[string]interface{}{
			"h1":      10,
			"archive": &Archive{ID: int64(12)},
		}
		p1 := NewBatchParam(m, BaseParamHandler)
		f(p1, t)
		_, exist := p1.Params["uuid"]
		So(exist, ShouldEqual, false)

		p2 := NewBatchParam(m, PushParamHandler)
		f(p2, t)
		_, exist = p2.Params["uuid"]
		So(exist, ShouldEqual, true)
		uuid = p2.Params["uuid"]

		f(p2, t, []int{1, 2})
		_, exist = p2.Params["uuid"]
		So(exist, ShouldEqual, true)
		So(p2.Params["uuid"], ShouldNotEqual, uuid)

		f(p2, t, []int64{1, 2})
		_, exist = p2.Params["uuid"]
		So(exist, ShouldEqual, true)

		var h ParamHandler
		p2.Params["h1"] = 1
		t.Logf("params(%+v), %+v, %v", p2.Params, h, h == nil)
	})
}
