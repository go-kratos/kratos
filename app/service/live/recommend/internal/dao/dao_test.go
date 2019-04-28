package dao

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeArrWithOrder(t *testing.T) {
	Convey("mergeArrWithOrder", t, func() {
		a := []int64{1, 2, 3, 4}
		b := []int64{4, 5, 6, 7}
		c := mergeArrWithOrder(a, b, len(a)+len(b))
		So(reflect.DeepEqual(c, []int64{1, 2, 3, 4, 5, 6, 7}), ShouldBeTrue)
		c = mergeArrWithOrder(a, b, 2)
		So(reflect.DeepEqual(c, []int64{1, 2, 3, 4}), ShouldBeTrue)
		c = mergeArrWithOrder(a, b, 6)
		So(reflect.DeepEqual(c, []int64{1, 2, 3, 4, 5, 6}), ShouldBeTrue)
	})
}

func TestMergeArr(t *testing.T) {
	Convey("mergeArr", t, func() {
		a := []int64{1, 2, 3, 4}
		b := []int64{4, 5, 6, 7}
		c := mergeArr(a, b)
		So(len(c) == 7, ShouldBeTrue)
	})
}
