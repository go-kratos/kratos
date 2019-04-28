package dao

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/admin/main/filter/conf"
	"go-common/app/admin/main/filter/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	dao *Dao
	ctx = context.TODO()
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/filter-admin-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	dao = New(conf.Conf)
	defer dao.Close()
	os.Exit(m.Run())
}

func Test_Area(t *testing.T) {
	var groupID int64
	Convey("TEST area group", t, func() {
		var (
			areaGroup = &model.AreaGroup{
				Name: "test group",
			}
			checkedAreaGroup *model.Area
		)
		tx, err := dao.BeginTran(ctx)
		So(err, ShouldBeNil)

		newID, err := dao.TxInsertAreaGroup(ctx, tx, areaGroup)
		So(err, ShouldBeNil)
		So(newID, ShouldBeGreaterThan, 0)
		groupID = newID
		areaGroupLog := &model.AreaGroupLog{
			AdID:    2333,
			AdName:  "muyang",
			Comment: "test comment",
			State:   model.LogStateAdd,
		}

		err = dao.TxInsertAreaGroupLog(ctx, tx, newID, areaGroupLog)
		So(err, ShouldBeNil)

		err = tx.Commit()
		So(err, ShouldBeNil)

		checkedAreaGroup, err = dao.AreaByName(ctx, areaGroup.Name)
		So(err, ShouldBeNil)
		So(checkedAreaGroup, ShouldNotBeNil)
		So(checkedAreaGroup.Name, ShouldEqual, areaGroup.Name)
	})
	Convey("TEST area", t, func() {
		So(groupID, ShouldBeGreaterThan, 0)
		var (
			area = &model.Area{
				Name:     "test area",
				ShowName: "test show name",
				GroupID:  int(groupID),
			}
			checkedArea *model.Area
		)
		tx, err := dao.BeginTran(ctx)
		So(err, ShouldBeNil)

		newID, err := dao.TxInsertArea(ctx, tx, area)
		So(err, ShouldBeNil)
		So(newID, ShouldBeGreaterThan, 0)

		areaLog := &model.AreaLog{
			AdID:    2333,
			AdName:  "muyang",
			Comment: "test comment",
			State:   model.LogStateAdd,
		}
		err = dao.TxInsertAreaLog(ctx, tx, newID, areaLog)
		So(err, ShouldBeNil)

		err = tx.Commit()
		So(err, ShouldBeNil)

		checkedArea, err = dao.AreaByName(ctx, area.Name)
		So(err, ShouldBeNil)
		So(checkedArea, ShouldNotBeNil)
		So(checkedArea.GroupID, ShouldEqual, groupID)
		So(checkedArea.Name, ShouldEqual, area.Name)
		So(checkedArea.ShowName, ShouldEqual, area.ShowName)
	})
}

func TestFilter(t *testing.T) {
	Convey("TEST filter", t, func() {
		var (
			f = &model.FilterInfo{
				Mode:    model.RegMode,
				Filter:  "ut.{0,2}测试",
				Level:   20,
				Source:  0,
				Type:    0,
				TpIDs:   []int64{0},
				Stime:   xtime.Time(time.Now().Unix()),
				Etime:   xtime.Time(time.Now().Add(time.Hour).Unix()),
				Comment: "ut test",
				State:   model.FilterStateNormal,
			}
		)
		tx, err := dao.BeginTran(ctx)
		So(err, ShouldBeNil)

		f.ID, err = dao.UpsertRule(ctx, tx, f.Filter, f.Comment, f.Level, f.Mode, f.Source, f.Type, f.Stime.Time(), f.Etime.Time())
		So(err, ShouldBeNil)
		So(f.ID, ShouldBeGreaterThan, 0)

		newID, err := dao.UpsertRule(ctx, tx, f.Filter, f.Comment, 30, f.Mode, f.Source, f.Type, f.Stime.Time(), f.Etime.Time())
		So(err, ShouldBeNil)
		So(newID, ShouldEqual, f.ID)

		_, err = dao.UpdateRuleState(ctx, f.ID, model.FilterStateDeleted)
		So(err, ShouldBeNil)

	})
}

func TestSearch(t *testing.T) {
	Convey("TEST search", t, func() {
		_, err := dao.Search(ctx, "test", "common", "0", "0", "20", 0, 0, 0, 10)
		So(err, ShouldNotBeNil)
		_, err = dao.SearchCount(ctx, "test", "common", "0", "0", "20", 0, 0)
		So(err, ShouldNotBeNil)
		_, err = dao.SearchKey(ctx, "test_key", "test", 0, 20, model.FilterStateNormal)
		So(err, ShouldNotBeNil)
		_, err = dao.SearchWhiteContent(ctx, "test", "common", 0, 10)
		So(err, ShouldNotBeNil)
	})
}

func TestConKeyById(t *testing.T) {
	Convey("TEST ConKeyById", t, func() {
		_, err := dao.ConKeyByID(ctx, 1, "test")
		So(err, ShouldNotBeNil)
	})
}

func Test_WhiteContent(t *testing.T) {
	Convey("TEST WhiteContent", t, func() {
		_, err := dao.WhiteContent(ctx, "test")
		So(err, ShouldNotBeNil)
	})
}

func Test_WhiteInfo(t *testing.T) {
	Convey("TEST WhiteInfo", t, func() {
		_, err := dao.WhiteInfo(ctx, 1)
		So(err, ShouldNotBeNil)
	})
}
func Test_AreaTotal(t *testing.T) {
	Convey("TEST AreaTotal", t, func() {
		_, err := dao.AreaTotal(ctx, 1)
		So(err, ShouldNotBeNil)
	})
}

func Test_KeyArea(t *testing.T) {
	Convey("TEST KeyArea", t, func() {
		_, err := dao.KeyArea(ctx, "test", 1)
		So(err, ShouldNotBeNil)
	})
}

func Test_AreaGroupTotal(t *testing.T) {
	Convey("TEST AreaGroupTotal", t, func() {
		_, err := dao.AreaGroupTotal(ctx)
		So(err, ShouldNotBeNil)
	})
}

func Test_MaxFilterID(t *testing.T) {
	Convey("TEST MaxFilterID", t, func() {
		_, err := dao.MaxFilterID(ctx)
		So(err, ShouldNotBeNil)
	})
}

func Test_CountKey(t *testing.T) {
	Convey("TEST CountKey", t, func() {
		_, err := dao.CountKey(ctx, "test", "comment", 1)
		So(err, ShouldNotBeNil)
	})
}

func Test_AreaGroupByName(t *testing.T) {
	Convey("TEST AreaGroupByName", t, func() {
		_, err := dao.AreaGroupByName(ctx, "test")
		So(err, ShouldNotBeNil)
	})
}
