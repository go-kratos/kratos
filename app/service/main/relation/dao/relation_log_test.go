package dao

import (
	"context"
	"go-common/app/service/main/relation/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddFollowingLog(t *testing.T) {
	var (
		ctx = context.Background()
		rl  = &model.RelationLog{}
	)
	convey.Convey("AddFollowingLog", t, func(cv convey.C) {
		d.AddFollowingLog(ctx, rl)
		cv.Convey("No return values", func(cv convey.C) {
		})
	})
}

func TestDaoDelFollowingLog(t *testing.T) {
	var (
		ctx = context.Background()
		rl  = &model.RelationLog{}
	)
	convey.Convey("DelFollowingLog", t, func(cv convey.C) {
		d.DelFollowingLog(ctx, rl)
		cv.Convey("No return values", func(cv convey.C) {
		})
	})
}

func TestDaoDelFollowerLog(t *testing.T) {
	var (
		ctx = context.Background()
		rl  = &model.RelationLog{}
	)
	convey.Convey("DelFollowerLog", t, func(cv convey.C) {
		d.DelFollowerLog(ctx, rl)
		cv.Convey("No return values", func(cv convey.C) {
		})
	})
}

func TestDaoAddWhisperLog(t *testing.T) {
	var (
		ctx = context.Background()
		rl  = &model.RelationLog{}
	)
	convey.Convey("AddWhisperLog", t, func(cv convey.C) {
		d.AddWhisperLog(ctx, rl)
		cv.Convey("No return values", func(cv convey.C) {
		})
	})
}

func TestDaoDelWhisperLog(t *testing.T) {
	var (
		ctx = context.Background()
		rl  = &model.RelationLog{}
	)
	convey.Convey("DelWhisperLog", t, func(cv convey.C) {
		d.DelWhisperLog(ctx, rl)
		cv.Convey("No return values", func(cv convey.C) {
		})
	})
}

func TestDaoAddBlackLog(t *testing.T) {
	var (
		ctx = context.Background()
		rl  = &model.RelationLog{}
	)
	convey.Convey("AddBlackLog", t, func(cv convey.C) {
		d.AddBlackLog(ctx, rl)
		cv.Convey("No return values", func(cv convey.C) {
		})
	})
}

func TestDaoDelBlackLog(t *testing.T) {
	var (
		ctx = context.Background()
		rl  = &model.RelationLog{}
	)
	convey.Convey("DelBlackLog", t, func(cv convey.C) {
		d.DelBlackLog(ctx, rl)
		cv.Convey("No return values", func(cv convey.C) {
		})
	})
}

func TestDaoaddLog(t *testing.T) {
	var (
		ctx      = context.Background()
		business = int(0)
		action   = ""
		rl       = &model.RelationLog{}
	)
	convey.Convey("addLog", t, func(cv convey.C) {
		d.addLog(ctx, business, action, rl)
		cv.Convey("No return values", func(cv convey.C) {
		})
	})
}
