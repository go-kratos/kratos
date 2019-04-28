package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/job/main/credit/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DelOrigin(t *testing.T) {
	var (
		c  = context.TODO()
		mr = &model.Case{}
	)
	mr.OriginType = int64(model.OriginTag)
	mr.RelationID = "57-4052400"
	Convey("should return err be nil", t, func() {
		s.DelOrigin(c, mr)
	})
}

func Test_DelOrigin2(t *testing.T) {
	var (
		c  = context.TODO()
		mr = &model.Case{}
	)
	mr.OriginType = int64(model.OriginReply)
	mr.RelationID = "111002543-3-400"
	Convey("should return err be nil", t, func() {
		s.DelOrigin(c, mr)
	})
}

func TestServiceDelJuryInfoCache(t *testing.T) {
	var (
		c  = context.TODO()
		mr = &model.Jury{}
	)
	mr.ID = 304
	mr.Status = model.CaseStatusDealing
	mr.Mid = 1
	Convey("should return err be nil", t, func() {
		b, _ := json.Marshal(&mr)
		err := s.DelJuryInfoCache(c, b)
		So(err, ShouldBeNil)
	})
}
