package service

import (
	"context"
	"testing"

	credit "go-common/app/interface/main/credit/model"

	. "github.com/smartystreets/goconvey/convey"
)

// TestSpJuryCase .
func TestSpJuryCase(t *testing.T) {
	Convey("TestSpJuryCase", t, func() {
		mid := int64(27515259)
		cid := int64(797)
		bc, err := s.SpJuryCase(context.TODO(), mid, cid)
		So(err, ShouldBeNil)
		So(bc, ShouldNotBeNil)
	})
}

// TestVoteInfoMID .
func TestVoteInfoMID(t *testing.T) {
	Convey("TestVoteInfoMID", t, func() {
		mid := int64(27515259)
		cid := int64(797)
		bc, err := s.VoteInfoCache(context.TODO(), mid, cid)
		So(err, ShouldBeNil)
		So(bc, ShouldNotBeNil)
	})
}

// TestObtainCase .
func TestObtainCase(t *testing.T) {
	Convey("TestObtainCase", t, func() {
		mid := int64(88889018)
		cid, err := s.caseVoteID(context.TODO(), mid, 0)
		So(err, ShouldBeNil)
		So(cid, ShouldNotBeNil)
	})
}

// TestServiceApply
func TestServiceApply(t *testing.T) {
	Convey("TestServiceApply", t, func() {
		var (
			mid int64 = 2089809
			c         = context.TODO()
		)
		err := s.Apply(c, mid)
		So(err, ShouldBeNil)
	})
}

// TestServiceRequirement
func TestServiceRequirement(t *testing.T) {
	Convey("TestServiceApply", t, func() {
		var (
			mid int64 = 88889019
			c         = context.TODO()
		)
		jr, err := s.Requirement(c, mid)
		So(err, ShouldBeNil)
		So(jr, ShouldNotBeNil)
	})
}

// TestServiceJury
func TestServiceJury(t *testing.T) {
	Convey("TestServiceJury", t, func() {
		var (
			mid int64 = 27515274
			c         = context.TODO()
		)
		res, err := s.Jury(c, mid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

// TestServiceCaseObtain
func TestServiceCaseObtain(t *testing.T) {
	Convey("TestServiceJury", t, func() {
		var (
			mid int64 = 21432418
			c         = context.TODO()
		)
		cid, err := s.CaseObtain(c, mid, 308)
		So(err, ShouldBeNil)
		So(cid, ShouldNotBeNil)
	})
}

// TestServiceVote .
func TestServiceVote(t *testing.T) {
	Convey("TestServiceVote", t, func() {
		var (
			mid int64 = 9
			cid int64 = 1
			v   int8  = 2
			c         = context.TODO()
		)
		err := s.Vote(c, mid, cid, v, 1, 1, 1, "", []int64{1, 2, 3}, []int64{4, 5, 6})
		So(err, ShouldBeNil)
	})
}

func TestServiceVoteInfo(t *testing.T) {
	Convey("TestServiceVoteInfo", t, func() {
		var (
			mid int64 = 9
			cid int64 = 1
			c         = context.TODO()
		)
		res, err := s.VoteInfo(c, mid, cid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestServiceCaseInfo(t *testing.T) {
	Convey("TestServiceCaseInfo", t, func() {
		var (
			c         = context.TODO()
			cid int64 = 304
		)
		res, err := s.CaseInfo(c, cid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestServiceJuryCase(t *testing.T) {
	Convey("TestServiceJuryCase", t, func() {
		var (
			mid int64 = 27515274
			cid int64 = 708
			c         = context.TODO()
		)
		res, err := s.JuryCase(c, mid, cid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestServiceCaseList(t *testing.T) {
	Convey("TestServiceCaseList", t, func() {
		var (
			mid int64 = 21432418
			pn  int64 = 1
			ps  int64 = 10
			c         = context.TODO()
		)
		res, err := s.CaseList(c, mid, ps, pn)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_KPIList(t *testing.T) {
	Convey("TestService_KPIList", t, func() {
		var (
			err error
			c   = context.TODO()
			res []*credit.KPI
		)
		res, err = s.KPIList(c, 88889017)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestCaseOpinion(t *testing.T) {
	Convey("TestCaseOpinion", t, func() {
		res, _, err := s.CaseOpinion(context.TODO(), 679, 1, 3)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestVoteOpinion(t *testing.T) {
	Convey("TestVoteOpinion", t, func() {
		res, _, err := s.VoteOpinion(context.TODO(), 679, 1, 3, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestAddBlockedCases(t *testing.T) {
	Convey("TestVoteOpinion", t, func() {
		var (
			bc []*credit.ArgJudgeCase
			b  = &credit.ArgJudgeCase{}
			b1 = &credit.ArgJudgeCase{}
		)
		b.MID = 1
		b.Operator = "a"
		b.OContent = "'aaaaa,"
		b.PunishResult = 1
		b.OTitle = "ot"
		b.OType = 1
		b.OURL = "ou"
		b.BlockedDays = 1
		b.ReasonType = 1
		b.RelationID = "1-1"
		b.RPID = 11
		b.Type = 1
		b.OID = 1111
		b.MID = 1

		b1.Operator = "b"
		b1.OContent = ",'bbbbb'''"
		b1.PunishResult = 2
		b1.OTitle = "ot"
		b1.OType = 1
		b1.OURL = "ou"
		b1.BlockedDays = 2
		b1.ReasonType = 2
		b1.RelationID = "1-1"
		b1.RPID = 22
		b1.Type = 2
		b1.OID = 22222
		bc = append(bc, b)
		bc = append(bc, b1)
		err := s.AddBlockedCases(context.TODO(), bc)
		So(err, ShouldBeNil)
	})
}

func TestCaseObtainByID(t *testing.T) {
	Convey("TestVoteOpinion", t, func() {
		var (
			c = context.TODO()
		)
		err := s.CaseObtainByID(c, 88889018, 309)
		So(err, ShouldBeNil)
	})
}
