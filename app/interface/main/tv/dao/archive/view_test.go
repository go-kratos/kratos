package archive

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_View3(t *testing.T) {
	Convey("view3", t, func() {
		aid, errGet := getPassAid(d.db)
		if errGet != nil {
			fmt.Println(errGet)
			return
		}
		vp, err := d.view3(context.Background(), aid)
		fmt.Println(vp, " ", aid)
		So(err, ShouldBeNil)
		So(vp, ShouldNotBeNil)
	})
}

func Test_ViewCache(t *testing.T) {
	Convey("viewCache", t, func() {
		reply, err := d.viewCache(context.TODO(), 123)
		So(reply, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestDao_Progress(t *testing.T) {
	Convey("TestDao_Progress", t, func() {
		aid, errGet := getPassAid(d.db)
		if errGet != nil {
			fmt.Println(errGet)
			return
		}
		h, err := d.Progress(context.TODO(), aid, 88895137)
		So(err, ShouldBeNil)
		fmt.Println(h)
	})
}

func TestDao_LoadViews(t *testing.T) {
	Convey("TestDao_LoadViews", t, func() {
		resMetas := d.LoadViews(context.Background(), []int64{10110328, 10099960})
		data, _ := json.Marshal(resMetas)
		fmt.Println(string(data))
	})
}

func TestDao_Archives(t *testing.T) {
	Convey("TestDao_Archives", t, WithDao(func(d *Dao) {
		aid, errGet := getPassAid(d.db)
		if errGet != nil {
			fmt.Println(errGet)
			return
		}
		as, err := d.Archives(context.Background(), []int64{aid})
		So(err, ShouldBeNil)
		So(len(as), ShouldBeGreaterThan, 0)
		fmt.Println(as)
	}))
}

func TestDao_GetView(t *testing.T) {
	Convey("TestDao_Archives", t, WithDao(func(d *Dao) {
		aid, errGet := getPassAid(d.db)
		if errGet != nil {
			fmt.Println(errGet)
			return
		}
		vp, err := d.GetView(context.Background(), aid)
		So(err, ShouldBeNil)
		So(vp, ShouldNotBeNil)
		fmt.Println(vp)
	}))
}

func TestDao_Archive3(t *testing.T) {
	Convey("TestDao_Archives", t, WithDao(func(d *Dao) {
		aid, errGet := getPassAid(d.db)
		if errGet != nil {
			fmt.Println(errGet)
			return
		}
		vp, err := d.Archive3(context.Background(), aid)
		So(err, ShouldBeNil)
		So(vp, ShouldNotBeNil)
		fmt.Println(vp)
	}))
}
