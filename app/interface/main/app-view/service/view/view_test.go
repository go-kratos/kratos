package view

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_View(t *testing.T) {
	Convey("View", t, func() {
		v, err := s.View(context.TODO(), 1, 10111556, 0, 0, 10000, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", time.Now())
		Println(v, err)
	})
}

func Test_ViewPage(t *testing.T) {
	Convey("ViewPage", t, func() {
		v, err := s.ViewPage(context.TODO(), 1, 1, 0, 0, 0, "", "", "", "", true, time.Now())
		Println(v, err)
	})
}

func Test_AddShare(t *testing.T) {
	Convey("AddShare", t, func() {
		_, _, _, err := s.AddShare(context.TODO(), 1, 1684013, "127.0.0.1")
		So(err, ShouldBeNil)
	})
}

func Test_Shot(t *testing.T) {
	Convey("Shot", t, func() {
		shot, _ := s.Shot(context.TODO(), 10106351, 10126396)
		fmt.Printf("===%+v===", shot)
	})
}

func Test_Like(t *testing.T) {
	Convey("Like", t, func() {
		s.Like(context.TODO(), 1, 1, 0)
	})
}

func Test_AddCoin(t *testing.T) {
	Convey("AddCoin", t, func() {
		s.AddCoin(context.TODO(), 1, 1684013, 2, 0, 1, "", 0)
	})
}

func Test_AddFav(t *testing.T) {
	Convey("AddFav", t, func() {
		_, err := s.AddFav(context.TODO(), 2, 1684013, []int64{1}, 1, "")
		So(err, ShouldBeNil)
	})
}

func Test_Paster(t *testing.T) {
	Convey("Paster", t, func() {
		s.Paster(context.TODO(), 1, 1, "1", "1", "")
	})
}

func Test_VipPlayURL(t *testing.T) {
	Convey("VipPlayURL", t, func() {
		s.VipPlayURL(context.TODO(), 1, 1, 1684013)
	})
}
