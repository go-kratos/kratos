package service

import (
	"testing"

	artmdl "go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

var a = map[int64]*artmdl.Category{
	1: &artmdl.Category{ID: 1, ParentID: 0, Position: 1},
	2: &artmdl.Category{ID: 2, ParentID: 0, Position: 0},
	3: &artmdl.Category{ID: 3, ParentID: 1, Position: 3},
	4: &artmdl.Category{ID: 4, ParentID: 3, Position: 4},
	5: &artmdl.Category{ID: 5, ParentID: 2, Position: 5},
}

func Test_transformCategory(t *testing.T) {
	Convey("should generate children", t, func() {
		expt := artmdl.Categories{
			&artmdl.Category{ID: 2, ParentID: 0, Position: 0, Children: []*artmdl.Category{
				&artmdl.Category{ID: 5, ParentID: 2, Position: 5},
			},
			},
			&artmdl.Category{ID: 1, ParentID: 0, Position: 1, Children: []*artmdl.Category{
				&artmdl.Category{ID: 3, ParentID: 1, Position: 3, Children: []*artmdl.Category{
					&artmdl.Category{ID: 4, ParentID: 3, Position: 4},
				},
				},
			}},
		}
		res := transformCategory(a)
		So(res, ShouldResemble, expt)
	})
}

func Test_transformReverseCategory(t *testing.T) {
	Convey("should generate all children category under a category", t, func() {
		expt1 := map[int64][]*artmdl.Category{
			1: []*artmdl.Category{a[3], a[4]},
			2: []*artmdl.Category{a[5]},
			3: []*artmdl.Category{a[4]},
		}
		expt2 := map[int64][]*artmdl.Category{
			1: []*artmdl.Category{a[4], a[3]},
			2: []*artmdl.Category{a[5]},
			3: []*artmdl.Category{a[4]},
		}
		res := transformReverseCategory(a)
		So(res, ShouldBeIn, expt1, expt2)
	})
}

func Test_transformCategoryParents(t *testing.T) {
	Convey("should generate all parent category", t, func() {
		expt := map[int64][]*artmdl.Category{
			1: []*artmdl.Category{a[1]},
			2: []*artmdl.Category{a[2]},
			3: []*artmdl.Category{a[1], a[3]},
			4: []*artmdl.Category{a[1], a[3], a[4]},
			5: []*artmdl.Category{a[2], a[5]},
		}
		res := transformCategoryParents(a)
		So(res, ShouldResemble, expt)
	})
}
