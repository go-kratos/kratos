package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceBaseQ(t *testing.T) {
	convey.Convey("BaseQ", t, func() {
		res, err := s.BaseQ(context.Background(), 14771787, "", false)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceBaseQs(t *testing.T) {
	convey.Convey("BaseQs", t, func() {
		rqs, err := s.BaseQs(context.Background(), 14771787, "", false)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rqs, convey.ShouldNotBeNil)
	})
}

func TestServiceConvertExtraQs(t *testing.T) {
	convey.Convey("ConvertExtraQs", t, func() {
		res, err := s.ConvertExtraQs(context.Background(), 14771787, "", false)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceExtraQs(t *testing.T) {
	convey.Convey("ExtraQs", t, func() {
		rqs, err := s.ExtraQs(context.Background(), 14771787, "", false)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rqs, convey.ShouldNotBeNil)
	})
}

func TestServicecheckExtraState(t *testing.T) {
	convey.Convey("checkExtraState", t, func() {
		h, err := s.checkExtraState(context.Background(), 14771787, time.Now())
		convey.So(err, convey.ShouldBeNil)
		convey.So(h, convey.ShouldNotBeNil)
	})
}

func TestServiceProTypes(t *testing.T) {
	convey.Convey("ProTypes", t, func() {
		res, err := s.proTypes(context.Background(), 14771787)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceProType(t *testing.T) {
	convey.Convey("ProType", t, func() {
		res, err := s.ProType(context.Background(), 14771787, "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceConvertProQues(t *testing.T) {
	convey.Convey("ConvertProQues", t, func() {
		res, err := s.ConvertProQues(context.Background(), 14771787, "", "", false)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceProQues(t *testing.T) {
	convey.Convey("ProQues", t, func() {
		rqs, err := s.ProQues(context.Background(), 14771787, "", "", false)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rqs, convey.ShouldNotBeNil)
	})
}

func TestServicecheckBase(t *testing.T) {
	convey.Convey("checkBase", t, func() {
		ah, err := s.checkBase(context.Background(), 0, time.Now())
		convey.So(err, convey.ShouldBeNil)
		convey.So(ah, convey.ShouldNotBeNil)
	})
}

func TestServicecheckTime(t *testing.T) {
	convey.Convey("checkTime", t, func() {
		at, rs := s.checkTime(context.Background(), 0, time.Now())
		convey.So(rs, convey.ShouldNotBeNil)
		convey.So(at, convey.ShouldNotBeNil)
	})
}

func TestServiceconcatData(t *testing.T) {
	convey.Convey("concatData", t, func() {
		rqs, err := s.concatData(context.Background(), 14771787, []int64{}, "", false, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rqs, convey.ShouldNotBeNil)
	})
}

func TestServiceconcatExtraData(t *testing.T) {
	convey.Convey("concatExtraData", t, func() {
		rqs, err := s.concatExtraData(context.Background(), 14771787, []int64{}, []int64{}, []int64{}, "", false, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rqs, convey.ShouldNotBeNil)
	})
}

func TestServiceansHash(t *testing.T) {
	convey.Convey("ansHash", t, func() {
		ansHash := s.ansHash(0, "")
		convey.So(ansHash, convey.ShouldNotBeNil)
	})
}

func TestServiceimgPosition(t *testing.T) {
	convey.Convey("imgPosition", t, func() {
		rq := s.imgPosition(context.Background(), nil, 14771787, "", false)
		convey.So(rq, convey.ShouldNotBeNil)
	})
}

func TestServiceimgExtraPosition(t *testing.T) {
	convey.Convey("imgExtraPosition", t, func() {
		rq := s.imgExtraPosition(context.Background(), nil, 14771787, "", false)
		convey.So(rq, convey.ShouldNotBeNil)
	})
}

func TestServiceloadQidsCache(t *testing.T) {
	convey.Convey("loadQidsCache", t, func() {
		s.loadQidsCache()
	})
}

func TestServiceloadExtraQidsCache(t *testing.T) {
	convey.Convey("loadExtraQidsCache", t, func() {
		s.loadExtraQidsCache()
	})
}

func TestServiceCool(t *testing.T) {
	convey.Convey("Cool", t, func() {
		cool, err := s.Cool(context.Background(), 0, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(cool, convey.ShouldNotBeNil)
	})
}

func TestServiceExtraScore(t *testing.T) {
	convey.Convey("ExtraScore", t, func() {
		score, err := s.ExtraScore(context.Background(), 6383240)
		fmt.Println(score)
		convey.So(err, convey.ShouldBeNil)
		convey.So(score, convey.ShouldBeGreaterThanOrEqualTo, 0)
	})
}
func TestServicehistory(t *testing.T) {
	convey.Convey("history", t, func() {
		ah, err := s.history(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ah, convey.ShouldNotBeNil)
	})
}

func TestServiceanswerDuration(t *testing.T) {
	convey.Convey("answerDuration", t, func() {
		d := s.answerDuration()
		convey.So(d, convey.ShouldNotBeNil)
	})
}

func TestSliceAtoi(t *testing.T) {
	convey.Convey("sliceAtoi", t, func() {
		p1, p2 := sliceAtoi([]string{})
		convey.So(p2, convey.ShouldBeNil)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestServiceextraQueByBigData(t *testing.T) {
	convey.Convey("extraQueByBigData", t, func() {
		ok, passids, npassids := s.extraQueByBigData(context.Background(), 0, "")
		convey.So(npassids, convey.ShouldNotBeNil)
		convey.So(passids, convey.ShouldNotBeNil)
		convey.So(ok, convey.ShouldNotBeNil)
	})
}

func TestServiceloadtypes(t *testing.T) {
	convey.Convey("loadtypes", t, func() {
		t := s.loadtypes()
		convey.So(t, convey.ShouldNotBeNil)
	})
}
