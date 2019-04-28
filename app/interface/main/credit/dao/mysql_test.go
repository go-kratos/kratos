package dao

import (
	"context"
	"math/rand"
	"testing"
	"time"

	model "go-common/app/interface/main/credit/model"

	"github.com/smartystreets/goconvey/convey"
)

// TestDao_LoadConf .
func TestDao_LoadConf(t *testing.T) {
	var c = context.Background()
	convey.Convey("return someting", t, func(convCtx convey.C) {
		v, err := d.LoadConf(c)
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(v, convey.ShouldNotBeNil)
	})
}

// TestJuryOpinion .
func TestJuryOpinion(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		tx, err := d.BeginTran(c)
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(tx, convey.ShouldNotBeNil)
		d.SetVoteTx(tx, 21432418, 383, 1)
		d.AddOpinionTx(tx, 383, 632, 11, "aaa", 1, 1, 1)
		tx.Commit()
		_, err = d.AddHates(c, []int64{632})
		convCtx.So(err, convey.ShouldBeNil)
		_, err = d.AddLikes(c, []int64{632})
		convCtx.So(err, convey.ShouldBeNil)
		ops, err := d.Opinions(c, []int64{632})
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(ops, convey.ShouldBeNil)
	})
}

// TestDao_IsVote .
func TestDao_IsVote(t *testing.T) {
	var (
		err error
		c   = context.Background()
		res int64
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		res, err = d.IsVote(c, 2528, 631)
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(res, convey.ShouldNotBeNil)
	})
}

// TestDao_CaseInfo .
func TestDao_CaseInfo(t *testing.T) {
	var (
		err error
		c   = context.Background()
		res = &model.BlockedCase{}
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		res, err = d.CaseInfo(c, -1)
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(res, convey.ShouldBeNil)
	})
}

// TestDao_VoteInfo .
func TestDao_VoteInfo(t *testing.T) {
	var (
		err error
		c   = context.Background()
		res = &model.VoteInfo{}
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		res, err = d.VoteInfo(c, 2528, 631)
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(res, convey.ShouldNotBeNil)
	})
}

// TestDao_ApplyJuryInfo .
func TestDao_ApplyJuryInfo(t *testing.T) {
	var (
		err error
		c   = context.Background()
		res = &model.BlockedJury{}
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		res, err = d.JuryInfo(c, 2528)
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(res, convey.ShouldNotBeNil)
	})
}

// TestDao_JuryInfo .
func TestDao_JuryInfo(t *testing.T) {
	var (
		err error
		c   = context.Background()
		res = &model.BlockedJury{}
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		res, err = d.JuryInfo(c, 0)
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(res, convey.ShouldNotBeNil)
	})
}

// TestDao_Setvote .
func TestDao_Setvote(t *testing.T) {
	var (
		err error
		c   = context.Background()
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		err = d.Setvote(c, 2528, 631, 1)
		convCtx.So(err, convey.ShouldBeNil)
	})
}

// TestDao_InsVote .
func TestDao_InsVote(t *testing.T) {
	var (
		err error
		c   = context.Background()
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		err = d.InsVote(c, rand.Int63n(9999999), rand.Int63n(9999999), 1)
		convCtx.So(err, convey.ShouldBeNil)
	})
}

func TestDao_AddUserVoteTotal(t *testing.T) {
	var (
		err error
		c   = context.Background()
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		err = d.AddUserVoteTotal(c, 2528)
		convCtx.So(err, convey.ShouldBeNil)
	})
}

func TestDao_AddCaseVoteTotal(t *testing.T) {
	var (
		err error
		c   = context.Background()
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		err = d.AddCaseVoteTotal(c, "put_total", 631, 1)
		convCtx.So(err, convey.ShouldBeNil)
	})
}

func TestDao_JuryApply(t *testing.T) {
	var (
		err error
		c   = context.Background()
	)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		err = d.JuryApply(c, 2528, time.Now().AddDate(0, 0, 30))
		convCtx.So(err, convey.ShouldBeNil)
	})
}

func TestOpinionIdx(t *testing.T) {
	convey.Convey("return someting", t, func(convCtx convey.C) {
		res, err := d.OpinionIdx(context.Background(), -1)
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(res, convey.ShouldBeNil)
	})
}

func TestDao_addQs(t *testing.T) {
	var qs = &model.LabourQs{}
	qs.Question = "test"
	qs.Ans = 1
	qs.AvID = 123
	qs.Status = 1
	convey.Convey("return someting", t, func(convCtx convey.C) {
		err := d.AddQs(context.Background(), qs)
		convCtx.So(err, convey.ShouldBeNil)
	})
}

func TestDao_setQs(t *testing.T) {
	convey.Convey("return someting", t, func(convCtx convey.C) {
		err := d.SetQs(context.Background(), 1, 2, 2)
		convCtx.So(err, convey.ShouldBeNil)
	})
}

func TestDao_delQs(t *testing.T) {
	convey.Convey("return someting", t, func(convCtx convey.C) {
		err := d.DelQs(context.Background(), 1, 1)
		convCtx.So(err, convey.ShouldBeNil)
	})
}

func TestDao_AddBlockedCases(t *testing.T) {
	var (
		err error
		c   = context.Background()
		bc  []*model.ArgJudgeCase
		b   = &model.ArgJudgeCase{}
	)
	b.MID = 1
	b.Operator = "a"
	b.OContent = "oc"
	b.PunishResult = 1
	b.OTitle = "ot"
	b.OType = 1
	b.OURL = "ou"
	b.BlockedDays = 1
	b.ReasonType = 1
	b.RelationID = "1-1"
	bc = append(bc, b)
	b = &model.ArgJudgeCase{}
	b.MID = 1
	b.Operator = "a"
	b.OContent = "oc"
	b.PunishResult = 1
	b.OTitle = "ot"
	b.OType = 1
	b.OURL = "ou"
	b.BlockedDays = 1
	b.ReasonType = 1
	b.RelationID = "1-1"
	bc = append(bc, b)
	convey.Convey("return someting", t, func(convCtx convey.C) {
		err = d.AddBlockedCases(c, bc)
		convCtx.So(err, convey.ShouldBeNil)
	})
}

func TestDao_NewKPI(t *testing.T) {
	convey.Convey("should return err be nil & map be nil", t, func(convCtx convey.C) {
		rate, err := d.NewKPI(context.Background(), 4)
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(rate, convey.ShouldBeGreaterThanOrEqualTo, 0)
		t.Logf("rate:%d", rate)
	})
}
