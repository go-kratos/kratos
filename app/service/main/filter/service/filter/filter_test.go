package filter

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/filter/conf"
	"go-common/app/service/main/filter/model"
	"go-common/app/service/main/filter/service/area"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	loader FilterAreas
	filter *Filter
	a      *area.Area
	ctx    = context.TODO()
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../../cmd/filter-service-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	//mock get area
	a = area.New()
	var areaLoader = func(context.Context) (list []*model.Area, err error) {
		area := &model.Area{
			ID:         1,
			GroupID:    1,
			Name:       "area1",
			ShowName:   "测试业务1",
			CommonFlag: true,
		}
		list = append(list, area)
		return
	}
	if err = a.Load(ctx, areaLoader); err != nil {
		panic(err)
	}

	//mock get rules
	loader = func(c context.Context, srcMask int64, area string) (rules []*model.FilterAreaInfo, err error) {
		rules = make([]*model.FilterAreaInfo, 0)
		rule1 := &model.FilterAreaInfo{
			ID:     1,
			TpIDs:  []int64{0, 1, 2},
			Mode:   1,
			Filter: "金三胖",
			Area:   "common",
			Source: 0x0,
			State:  0,
		}
		rule1.SetLevel(20, 30)
		rule2 := &model.FilterAreaInfo{
			ID:     2,
			TpIDs:  []int64{0, 1, 2},
			Mode:   0,
			Filter: "(A|C).*站",
			Area:   "reply",
			Source: 0x2,
			State:  0,
		}
		rule2.SetLevel(20, 20)
		if area == "common" {
			if srcMask == 0x00 {
				rules = append(rules, rule1)
			}
		}
		if area == "reply" {
			if srcMask == 0x02 {
				rules = append(rules, rule2)
			}
		}
		return
	}
	filter = New()
	os.Exit(m.Run())
}

func TestFilter(t *testing.T) {
	Convey("load", t, func() {
		err := filter.Load(ctx, loader, a)
		So(err, ShouldBeNil)
		Convey("filter", func() {
			var (
				filters []*model.Filter
			)
			filters = filter.GetFilters("common", true)
			So(filters, ShouldHaveLength, 1)

			level, rules := filters[0].Matcher.Test("铁拳无敌金三胖")
			So(level, ShouldEqual, 20)
			So(rules, ShouldHaveLength, 1)

			filters = filter.GetFilters("reply", false)
			So(filters, ShouldHaveLength, 1)
			So(filters[0].Regs, ShouldHaveLength, 1)
			matchFlag := filters[0].Regs[0].Reg.MatchString("A站")
			So(matchFlag, ShouldBeTrue)

			filters = filter.GetFilters("reply", true)
			So(filters, ShouldHaveLength, 2)

			filters = filter.GetFilters("reply", false)
			So(filters, ShouldHaveLength, 0)
		})

		Convey("filter by area", func() {
			filters := filter.GetFiltersByArea("common")
			So(filters, ShouldHaveLength, 1)

			filters = filter.GetFiltersByArea("reply")
			So(filters, ShouldHaveLength, 1)

			filters = filter.GetFiltersByArea("23333")
			So(filters, ShouldHaveLength, 0)
		})
	})
}
