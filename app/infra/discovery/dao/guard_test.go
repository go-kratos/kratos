package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIncrExp(t *testing.T) {
	Convey("test IncrExp", t, func() {
		re := new(Guard)
		re.incrExp()
		So(re.expPerMin, ShouldResemble, int64(2))
	})
}

func TestDecrExp(t *testing.T) {
	Convey("test DecrExp", t, func() {
		re := new(Guard)
		re.incrExp()
		re.decrExp()
		So(re.expPerMin, ShouldResemble, int64(0))
	})
}

func TestSetExp(t *testing.T) {
	Convey("test SetExp", t, func() {
		re := new(Guard)
		re.setExp(10)
		So(re.expPerMin, ShouldResemble, int64(20))
		So(re.expThreshold, ShouldResemble, int64(17))
	})
}

func TestUpdateFac(t *testing.T) {
	Convey("test UpdateFac", t, func() {
		re := new(Guard)
		re.incrFac()
		re.updateFac()
		So(re.facLastMin, ShouldResemble, int64(1))
	})
}

func TestIncrFac(t *testing.T) {
	Convey("test IncrFac", t, func() {
		re := new(Guard)
		re.incrFac()
		So(re.facInMin, ShouldResemble, int64(1))
	})
}

func TestIsProtected(t *testing.T) {
	Convey("test IncrFac", t, func() {
		re := new(Guard)
		re.incrExp()
		re.incrExp()
		re.incrFac()
		re.updateFac()
		So(re.ok(), ShouldBeTrue)
		re = new(Guard)
		re.incrExp()
		re.incrFac()
		re.updateFac()
		So(re.ok(), ShouldBeFalse)
	})
}
