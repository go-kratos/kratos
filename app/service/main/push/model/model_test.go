package model

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFuncs(t *testing.T) {
	Convey("int string functions", t, func() {
		Convey("SplitInts", func() {
			res := SplitInts("1,2,3")
			So(res, ShouldResemble, []int{1, 2, 3})
		})

		Convey("JoinInts", func() {
			ints := []int{1, 2, 3}
			res := JoinInts(ints)
			So(res, ShouldEqual, "1,2,3")
		})

		Convey("existsInt", func() {
			exists := ExistsInt([]int{}, 4)
			So(exists, ShouldBeFalse)
			ints := []int{1, 2, 3}
			exists = ExistsInt(ints, 1)
			So(exists, ShouldBeTrue)
			exists = ExistsInt(ints, 4)
			So(exists, ShouldBeFalse)
		})

		Convey("gen temp task id", func() {
			id := TempTaskID()
			So(len(id), ShouldEqual, 9)
		})

		Convey("gen job name", func() {
			name := JobName(time.Now().UnixNano(), "123", "456", "g")
			t.Logf("job name is: %d", name)
		})
	})

	Convey("ParseBuild", t, func() {
		buildString := `{"2":{"Build":100,"Condition":"gt"}}`
		build := ParseBuild(buildString)
		So(build, ShouldResemble, map[int]*Build{2: {Build: 100, Condition: "gt"}})
	})

	Convey("platform", t, func() {
		plat := Platform("iphone", PushSDKApns)
		So(plat, ShouldEqual, PlatformIPhone)
		plat = Platform("ipad", PushSDKApns)
		So(plat, ShouldEqual, PlatformIPad)
		plat = Platform("whatever", PushSDKXiaomi)
		So(plat, ShouldEqual, PlatformXiaomi)
	})

	Convey("parse silent time", t, func() {
		st := ParseSilentTime("22:30-06:00")
		So(st, ShouldResemble, BusinessSilentTime{
			BeginHour:   22,
			EndHour:     6,
			BeginMinute: 30,
			EndMinute:   0,
		})
	})
}

func TestValidateBuild(t *testing.T) {
	builds := map[int]*Build{
		1: {Build: 520000, Condition: "eq"},
		2: {Build: 123456, Condition: "gt"},
	}
	Convey("ValidateBuild", t, func() {
		b := ValidateBuild(2, 123455, builds)
		So(b, ShouldBeFalse)
		b = ValidateBuild(2, 123457, builds)
		So(b, ShouldBeTrue)
		b = ValidateBuild(4, 520001, builds)
		So(b, ShouldBeFalse)
		b = ValidateBuild(4, 519999, builds)
		So(b, ShouldBeFalse)
		b = ValidateBuild(4, 520000, builds)
		So(b, ShouldBeTrue)
	})
}

func TestScheme(t *testing.T) {
	Convey("Scheme()", t, func() {
		scheme := Scheme(LinkTypeLive, "1,0", PlatformAndroid, 5300000)
		So(scheme, ShouldEqual, "bilibili://live/1?broadcast_type=0")
		scheme = Scheme(LinkTypeLive, "1", PlatformAndroid, 5280000)
		So(scheme, ShouldEqual, "bili:///?type=bililive&roomid=1")
		scheme = Scheme(LinkTypeLive, "1,1", PlatformIPhone, 5300000)
		So(scheme, ShouldEqual, "bilibili://live/1?broadcast_type=1")
		scheme = Scheme(LinkTypeLive, "1,0", PlatformIPhone, 5280000)
		So(scheme, ShouldEqual, "bilibili://live/1?broadcast_type=0")
		scheme = Scheme(LinkTypeCustom, "custom_scheme", PlatformIPhone, 68)
		So(scheme, ShouldEqual, "custom_scheme")
	})
}
