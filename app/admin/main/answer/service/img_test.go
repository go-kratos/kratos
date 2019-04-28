package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/answer/model"
	"go-common/library/xstr"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceImg(t *testing.T) {
	Convey("TestServiceImg", t, func() {
		ids, err := xstr.SplitInts("35")
		s.CreateBFSImg(context.TODO(), ids)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestAddImg
func TestAddImg(t *testing.T) {
	q := &model.QuestionDB{ID: 1, State: 1, Question: "test2333"}
	Convey("TestAddImgQuestionAdd", t, func() {
		err := s.QuestionAdd(context.TODO(), q)
		So(err, ShouldBeNil)
	})
	Convey("TestAddImg", t, func() {
		s.GenerateImage(context.TODO())
	})
}
