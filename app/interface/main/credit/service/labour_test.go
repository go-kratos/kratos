package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/credit/model"

	. "github.com/smartystreets/goconvey/convey"
)

// TestServiceAddQs .
func TestServiceAddQs(t *testing.T) {
	Convey("TestServiceAddQs", t, func() {
		qs := &model.LabourQs{}
		qs.Question = "test"
		qs.AvID = 666
		qs.Source = 1
		err := s.AddQs(context.TODO(), qs)
		So(err, ShouldBeNil)
	})
}

// TestServiceSetQs
func TestServiceSetQs(t *testing.T) {
	Convey("TestServiceSetQs", t, func() {
		var (
			c            = context.TODO()
			id     int64 = 3
			ans    int64 = 2
			status int64 = 2
		)
		err := s.SetQs(c, id, ans, status)
		So(err, ShouldBeNil)
	})
}

// TestCommit .
func TestCommit(t *testing.T) {
	Convey("TestCommit", t, func() {
		ans := &model.LabourAns{
			ID:     []int64{},
			Answer: []int64{},
		}
		for i := 0; i < 40; i++ {
			ans.ID = append(ans.ID, int64(i))
			ans.Answer = append(ans.Answer, 1)
		}
		res, err := s.CommitQs(context.TODO(), 88895349, "aa", "1", "11", ans)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
