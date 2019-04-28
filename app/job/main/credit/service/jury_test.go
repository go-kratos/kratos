package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/job/main/credit/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_BlockUser(t *testing.T) {
	var (
		c  = context.TODO()
		mr = &model.Case{}
	)
	mr.ID = 304
	mr.JudgeType = model.JudgeTypeViolate
	mr.OriginTitle = "304"
	mr.Mid = 1
	mr.OriginURL = "http:304"
	mr.OriginContent = "304cont"
	mr.OriginType = 4
	mr.Operator = "lgs"
	mr.Status = 4
	mr.ReasonType = 3
	mr.RelationID = "10741-4052410"
	mr.PunishResult = 2
	Convey("should return err be nil", t, func() {
		err := s.BlockUser(c, mr)
		So(err, ShouldBeNil)
	})
}

func Test_Judge(t *testing.T) {
	var (
		c  = context.TODO()
		mr = &model.Case{}
	)
	mr.ID = 304
	mr.Status = model.CaseStatusDealing
	mr.JudgeType = model.JudgeTypeViolate
	mr.OriginTitle = "304"
	mr.Mid = 1
	mr.OriginURL = "http:304"
	mr.OriginContent = "304cont"
	mr.OriginType = 4
	mr.Operator = "lgs"
	mr.Status = 4
	mr.ReasonType = 3
	mr.RelationID = "10741-4052410"
	mr.PunishResult = 2
	Convey("should return err be nil", t, func() {
		b, _ := json.Marshal(&mr)
		err := s.Judge(c, b, b)
		So(err, ShouldBeNil)
	})
}

func Test_dealMoralCase(t *testing.T) {
	var (
		c  = context.TODO()
		mr = &model.Case{}
	)
	mr.ID = 304
	mr.Status = model.CaseStatusDealing
	mr.JudgeType = model.JudgeTypeViolate
	mr.OriginTitle = "304"
	mr.Mid = 1
	mr.OriginURL = "http:304"
	mr.OriginContent = "304cont"
	mr.OriginType = 4
	mr.Operator = "lgs"
	mr.Status = 4
	mr.ReasonType = 3
	mr.RelationID = "10741-4052410"
	mr.PunishResult = 2
	Convey("dealMoralCase", t, func() {
		err := s.dealMoralCase(c, mr)
		So(err, ShouldBeNil)
	})
}
