package dao

import (
	"reflect"
	"testing"

	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/recommend/internal/conf"
)

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
}

func TestOneHotTextEncode(t *testing.T) {
	Convey("oneHotTextEncode", t, func() {
		arr := oneHotTextEncode("", []string{"", ".*人存活", "决赛圈", "正在抽奖", ".*No\\.\\d+", "年度.*主播"})
		So(reflect.DeepEqual(arr, []int64{1, 0, 0, 0, 0, 0}), ShouldBeTrue)
		arr = oneHotTextEncode("23人存活", []string{"", ".*人存活", "决赛圈", "正在抽奖", ".*No\\.\\d+", "年度.*主播"})
		So(reflect.DeepEqual(arr, []int64{0, 1, 0, 0, 0, 0}), ShouldBeTrue)
		arr = oneHotTextEncode("决赛圈", []string{"", ".*人存活", "决赛圈", "正在抽奖", ".*No\\.\\d+", "年度.*主播"})
		So(reflect.DeepEqual(arr, []int64{0, 0, 1, 0, 0, 0}), ShouldBeTrue)
		arr = oneHotTextEncode("正在抽奖", []string{"", ".*人存活", "决赛圈", "正在抽奖", ".*No\\.\\d+", "年度.*主播"})
		So(reflect.DeepEqual(arr, []int64{0, 0, 0, 1, 0, 0}), ShouldBeTrue)
		arr = oneHotTextEncode("上小时电台No.1", []string{"", ".*人存活", "决赛圈", "正在抽奖", ".*No\\.\\d+", "年度.*主播"})
		So(reflect.DeepEqual(arr, []int64{0, 0, 0, 0, 1, 0}), ShouldBeTrue)
		arr = oneHotTextEncode("年度五强主播", []string{"", ".*人存活", "决赛圈", "正在抽奖", ".*No\\.\\d+", "年度.*主播"})
		So(reflect.DeepEqual(arr, []int64{0, 0, 0, 0, 0, 1}), ShouldBeTrue)
	})
}

func TestOneHotEncode(t *testing.T) {
	Convey("oneHotEncode", t, func() {
		arr := oneHotEncode(78, []int64{23, 54, 100, 120})
		So(reflect.DeepEqual(arr, []int64{0, 0, 1, 0, 0}), ShouldBeTrue)
		arr = oneHotEncode(7, []int64{23, 54, 100, 120})
		So(reflect.DeepEqual(arr, []int64{1, 0, 0, 0, 0}), ShouldBeTrue)
		arr = oneHotEncode(200, []int64{23, 54, 100, 120})
		So(reflect.DeepEqual(arr, []int64{0, 0, 0, 0, 1}), ShouldBeTrue)
	})
}

func TestSliceArray(t *testing.T) {
	Convey("sliceArray", t, func() {
		arr := sliceArray([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 4)
		So(reflect.DeepEqual(arr[0], []int64{1, 2, 3, 4}), ShouldBeTrue)
		So(reflect.DeepEqual(arr[1], []int64{5, 6, 7, 8}), ShouldBeTrue)
		So(reflect.DeepEqual(arr[2], []int64{9}), ShouldBeTrue)
	})
}

func TestCreateRoomFeature(t *testing.T) {
	Convey("createRoomFeature", t, func() {
		c := conf.Conf
		arr := createFeature(c, 21, "决赛圈", 2000, 1000)
		So(reflect.DeepEqual(arr, []int64{21, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}), ShouldBeTrue)
	})
}
