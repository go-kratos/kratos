package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/answer/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceBatchUpdateState(t *testing.T) {
	convey.Convey("BatchUpdateState", t, func() {
		err := s.BatchUpdateState(context.TODO(), []int64{1, 2, 3}, 0, "")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestServiceTypes(t *testing.T) {
	convey.Convey("Types", t, func() {
		res, err := s.Types(context.TODO())
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceQuestionAdd(t *testing.T) {
	q := &model.QuestionDB{ID: 1, Question: "test2333"}
	convey.Convey("QuestionAdd", t, func() {
		err := s.QuestionAdd(context.TODO(), q)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("UpdateStatus", t, func() {
		err := s.UpdateStatus(context.TODO(), 1, 0, "")
		convey.So(err, convey.ShouldBeNil)
	})
	arg := &model.ArgQue{State: 1}
	convey.Convey("QuestionList", t, func() {
		res, err := s.QuestionList(context.TODO(), arg)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceloadtypes(t *testing.T) {
	convey.Convey("loadtypes", t, func() {
		t, err := s.loadtypes(context.TODO())
		convey.So(err, convey.ShouldBeNil)
		convey.So(t, convey.ShouldNotBeNil)
	})
}
