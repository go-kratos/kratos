package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/ep/melloi/model"
)

var (
	label = model.Label{
		ID:          1,
		Name:        "name",
		Description: "description",
		Color:       "#602230",
		Active:      1,
	}

	lr = model.LabelRelation{
		ID:          1,
		LabelID:     1,
		LabelName:   "name",
		Color:       "#602230",
		Description: "description",
		TargetID:    1,
		Type:        1,
		Active:      1,
	}
	ids = []int64{1, 2, 3, 4, 5}
)

func Test_Label(t *testing.T) {

	Convey("add label", t, func() {
		err := s.AddLabel(&label)
		So(err, ShouldBeNil)
	})
	Convey("query label", t, func() {
		_, err := s.QueryLabel(c)
		So(err, ShouldBeNil)
	})
	Convey("delete label", t, func() {
		err := s.DeleteLabel(label.ID)
		So(err, ShouldBeNil)
	})

	Convey("delete label relation", t, func() {
		err := s.AddLabelRelation(&lr)
		So(err, ShouldBeNil)
	})

	Convey("delete label relation", t, func() {
		err := s.DeleteLabelRelation(lr.ID)
		So(err, ShouldBeNil)
	})

	Convey("query label relation", t, func() {
		_, err := s.QueryLabelRelation(&lr)
		So(err, ShouldBeNil)
	})
	Convey("query label relation by ids", t, func() {
		_, err := s.QueryLabelRelationByIDs(ids)
		So(err, ShouldBeNil)
	})

}
