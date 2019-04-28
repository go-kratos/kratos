package elastic

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Upsert(t *testing.T) {
	Convey("upsert", t, func() {
		e := NewElastic(nil)
		us := e.NewUpdate("articledata").Insert()
		data := map[string]int{"id": 552, "share": 321}
		us.AddData("articledata", data)
		data = map[string]int{"id": 558, "share": 1137}
		us.AddData("articledata", data)
		t.Logf("params(%s)", us.Params())
		err := us.Do(context.Background())
		So(err, ShouldBeNil)
	})
}

func Test_Upsert_ByMod(t *testing.T) {
	Convey("upsert", t, func() {
		data := make([]map[string]interface{}, 0)
		d := map[string]interface{}{"id": 22, "share": 22, "oid": int64(2266)}
		data = append(data, d)
		d = map[string]interface{}{"id": 33, "share": 33, "oid": int64(3388)}
		data = append(data, d)

		e := NewElastic(nil)
		us := e.NewUpdate("reply").Insert()
		for _, v := range data {
			oid, ok := v["oid"]
			if !ok {
				continue
			}
			oidReal, ok := oid.(int64) //must as same as interface type
			if !ok {
				continue
			}
			us.AddData(us.IndexByMod("reply_record", oidReal, 100), v)
		}
		t.Logf("params(%s)", us.Params())
		if us.HasData() {
			err := us.Do(context.Background())
			So(err, ShouldBeNil)
		} else {
			So(us.HasData(), ShouldBeTrue)
		}
	})
}

func Test_Update_Index(t *testing.T) {
	Convey("test update index", t, func() {
		e := NewElastic(nil)
		us := e.NewUpdate("reply").Insert()
		Convey("example by mod 100", func() {
			oid := int64(808)
			index := us.IndexByMod("reply", oid, 100)
			So(index, ShouldEqual, "reply_08")
		})
		Convey("example by mod 1000", func() {
			oid := int64(808)
			index := us.IndexByMod("reply", oid, 1000)
			So(index, ShouldEqual, "reply_808")
		})
		Convey("example by mod 10000", func() {
			oid := int64(808)
			index := us.IndexByMod("reply", oid, 10000)
			So(index, ShouldEqual, "reply_0808")
		})
		Convey("example by mod oid 0", func() {
			oid := int64(0)
			index := us.IndexByMod("reply", oid, 100)
			So(index, ShouldEqual, "reply_00")
		})
		Convey("example by mod oid 20", func() {
			oid := int64(808)
			index := us.IndexByMod("reply", oid, 20)
			So(index, ShouldEqual, "reply_08")
		})
	})
}

func Test_Upsert_ByTime(t *testing.T) {
	Convey("upsert", t, func() {
		data := make([]map[string]interface{}, 0)
		d := map[string]interface{}{"id": 22, "share": 22, "ctime": time.Now().AddDate(-2, 0, 0)}
		data = append(data, d)
		d = map[string]interface{}{"id": 33, "share": 33, "ctime": time.Now().AddDate(-3, 0, 0)}
		data = append(data, d)

		e := NewElastic(nil)
		us := e.NewUpdate("reply_list").Insert()
		for _, v := range data {
			ctime, ok := v["ctime"]
			if !ok {
				continue
			}
			ctimeReal, ok := ctime.(time.Time) //must as same as interface type
			if !ok {
				continue
			}
			indexName := us.IndexByTime("reply_list_hot", IndexTypeYear, ctimeReal)
			v["ctime"] = ctimeReal.Format("2006-01-02 15:04:05")
			us.AddData(indexName, v)
		}
		t.Logf("params(%s)", us.Params())
		if us.HasData() {
			err := us.Do(context.Background())
			So(err, ShouldBeNil)
		} else {
			So(us.HasData(), ShouldBeTrue)
		}
	})
}
