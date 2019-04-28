package service

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpVideo(t *testing.T) {
	Convey("UpVideo", t, func() {
		err := s.UpVideo(context.TODO(), 1, 1)
		So(err, ShouldNotBeNil)
	})
}

func Test_DelVideo(t *testing.T) {
	Convey("DelVideo", t, func() {
		err := s.DelVideo(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
	})
}

func Test_Description(t *testing.T) {
	Convey("Description", t, func() {
		_, err := s.Description(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_Page3(t *testing.T) {
	Convey("Page3", t, func() {
		ps, err := s.Page3(context.TODO(), 10098500)
		So(err, ShouldBeNil)
		for _, p := range ps {
			Printf("%+v\n\n", p)
			bs, _ := json.Marshal(p)
			Printf("%s\n\n", bs)
		}
	})
}

func Test_View3(t *testing.T) {
	Convey("View3", t, func() {
		v, err := s.View3(context.TODO(), 10098500)
		So(err, ShouldBeNil)
		if v.Pages != nil {
			for _, p := range v.Pages {
				Printf("%+v\n\n", p)
				bs, _ := json.Marshal(p)
				Printf("%s\n\n", bs)
			}
		}

	})
}

func Test_Views3(t *testing.T) {
	Convey("Views3", t, func() {
		as, err := s.Views3(context.TODO(), []int64{10098500, 10097755})
		So(err, ShouldBeNil)
		for _, a := range as {
			for _, p := range a.Pages {
				Printf("%+v\n\n", p)
				bs, _ := json.Marshal(p)
				Printf("%s\n\n", bs)
			}
		}
	})
}

func Test_Video3(t *testing.T) {
	Convey("Video3", t, func() {
		v, err := s.Video3(context.TODO(), 10098500, 10109206)
		Printf("%+v\n\n\n", v)
		So(err, ShouldBeNil)
		bs, _ := json.Marshal(v)
		Printf("%s\n\n\n", bs)
	})
}
