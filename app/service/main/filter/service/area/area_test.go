package area

import (
	"context"
	"testing"

	"go-common/app/service/main/filter/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	area *Area
	ctx  = context.TODO()
)

func TestMain(m *testing.M) {
	var err error
	area = New()
	var loader = func(context.Context) (list []*model.Area, err error) {
		list = append(list, &model.Area{
			ID:         1,
			GroupID:    1,
			Name:       "message",
			ShowName:   "消息业务",
			CommonFlag: false,
		})
		list = append(list, &model.Area{
			ID:         2,
			GroupID:    1,
			Name:       "danmu",
			ShowName:   "弹幕",
			CommonFlag: true,
		})
		list = append(list, &model.Area{
			ID:         3,
			GroupID:    1,
			Name:       "reply",
			ShowName:   "评论",
			CommonFlag: true,
		})
		list = append(list, &model.Area{
			ID:         4,
			GroupID:    1,
			Name:       "test",
			ShowName:   "测试业务",
			CommonFlag: false,
		})
		list = append(list, &model.Area{
			ID:         5,
			GroupID:    2,
			Name:       "common",
			ShowName:   "基础库",
			CommonFlag: true,
		})
		return
	}
	if err = area.Load(ctx, loader); err != nil {
		panic(err)
	}

	m.Run()
}

func TestArea(t *testing.T) {
	Convey("area", t, func() {
	})
}

func TestStrategy(t *testing.T) {
	Convey("test filter strategy", t, func() {
		areaName := "common"
		s := area.Area(areaName)
		So(s, ShouldNotBeNil)
		So(s.IsFilter(), ShouldBeTrue)
		So(s.IsFilterCommon(), ShouldBeTrue)
		So(s.IsFilterKey([]string{}), ShouldBeFalse)
		So(s.IsFilterRubbish(1), ShouldBeFalse)

		areaName = "message"
		s = area.Area(areaName)
		So(s, ShouldNotBeNil)
		So(s.IsFilter(), ShouldBeFalse)
		So(s.IsFilterCommon(), ShouldBeFalse)
		So(s.IsFilterKey([]string{}), ShouldBeFalse)
		So(s.IsFilterRubbish(1), ShouldBeTrue)
		So(s.IsFilterRubbish(0), ShouldBeTrue)

		areaName = "danmu"
		s = area.Area(areaName)
		So(s, ShouldNotBeNil)
		So(s.IsFilter(), ShouldBeTrue)
		So(s.IsFilterCommon(), ShouldBeTrue)
		So(s.IsFilterKey([]string{}), ShouldBeTrue)
		So(s.IsFilterKey([]string{"aid:2333"}), ShouldBeTrue)
		So(s.IsFilterRubbish(1), ShouldBeTrue)
		So(s.IsFilterRubbish(0), ShouldBeFalse)

		areaName = "reply"
		s = area.Area(areaName)
		So(s, ShouldNotBeNil)
		So(s.IsFilter(), ShouldBeTrue)
		So(s.IsFilterCommon(), ShouldBeTrue)
		So(s.IsFilterKey([]string{}), ShouldBeFalse)
		So(s.IsFilterKey([]string{"test"}), ShouldBeTrue)
		So(s.IsFilterRubbish(0), ShouldBeFalse)
		So(s.IsFilterRubbish(1), ShouldBeTrue)

		areaName = "test"
		s = area.Area(areaName)
		So(s, ShouldNotBeNil)
		So(s.IsFilter(), ShouldBeTrue)
		So(s.IsFilterCommon(), ShouldBeFalse)
		So(s.IsFilterKey([]string{}), ShouldBeFalse)
		So(s.IsFilterKey([]string{"test"}), ShouldBeTrue)
		So(s.IsFilterRubbish(1), ShouldBeFalse)

		areaName = "foo"
		s = area.Area(areaName)
		So(s, ShouldBeNil)
	})
}
