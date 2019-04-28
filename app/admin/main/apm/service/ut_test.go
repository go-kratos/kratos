package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/apm/model/ut"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceUtDetailList(t *testing.T) {
	convey.Convey("UtDetailList", t, func(ctx convey.C) {
		req := &ut.DetailReq{
			CommitID: "4f5ca7122e023c605555ed4bddf274709635d019",
			PKG:      "go-common/app/admin/main/apm",
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			utpkgs, err := svr.UtDetailList(context.Background(), req)
			ctx.Convey("Than utpkgs should not be nil. err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(utpkgs, convey.ShouldNotBeNil)
			})
		})
	})
}
func TestServiceUtHistoryCommit(t *testing.T) {
	convey.Convey("UtHistoryCommit", t, func(ctx convey.C) {
		dt := &ut.HistoryCommitReq{
			MergeID: 1,
			Ps:      20,
			Pn:      1,
		}
		ctx.Convey("When everything goes possitive", func(ctx convey.C) {
			utcmts, count, err := svr.UtHistoryCommit(context.Background(), dt)
			ctx.Convey("Then err should be nil.count shoulde be greater than 0.utcmts should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldBeGreaterThan, 0)
				ctx.So(utcmts, convey.ShouldNotBeNil)
				t.Logf("history count: %+v", count)
				for _, commit := range utcmts {
					for _, pkg := range commit.PkgAnls {
						t.Logf("%s %s cov:%.2f\n", pkg.CommitID, pkg.PKG, pkg.Coverage)
					}
				}
			})
		})
	})
}

//
//func TestService_CreateUtMerge(t *testing.T) {
//	convey.Convey("CreateUtMerge", t, func() {
//		res := &ut.UploadRes{
//			CommitID: "eaefwa",
//			UserName: "ChengXing",
//			MergeID:  1000,
//		}
//		err := svr.AddUTMerge(context.Background(), res)
//		convey.So(err, convey.ShouldBeNil)
//	})
//}
//
//func TestService_CreateUtCommit(t *testing.T) {
//	convey.Convey("CreateUtCommit", t, func() {
//		res := &ut.UploadRes{
//			CommitID: "eaefwa",
//			UserName: "ChengXing",
//			MergeID:  1000,
//		}
//		err := svr.AddUTCommit(context.Background(), res)
//		convey.So(err, convey.ShouldBeNil)
//	})
//}
//
//func TestService_CreateUtPKG(t *testing.T) {
//	convey.Convey("CreateUtPKG", t, func() {
//		pkg := &ut.PkgAnls{
//			CommitID:   "eawfsdfeasdfe",
//			PassRate:   64.12,
//			PKG:        "go-common/app/server/main/coin/dao",
//			Assertions: 10,
//			Passed:     8,
//		}
//		res := &ut.UploadRes{
//			CommitID: "eaefwa",
//			UserName: "ChengXing",
//			MergeID:  1000,
//		}
//		err := svr.AddUTPKG(context.Background(), pkg, res, "", "")
//		convey.So(err, convey.ShouldBeNil)
//	})
//}

// func TestService_SAGAReport(t *testing.T) {
// 	convey.Convey("SAGAReport", t, func() {
// 		saga, err := svr.SAGAReport(context.Background(), "8c9d27c106caf6325c1155d127fcc952d1f4a441")
// 		convey.So(saga, convey.ShouldNotBeEmpty)
// 		convey.So(err, convey.ShouldBeNil)
// 		for _, s := range saga {
// 			t.Logf("Coverage:(%.02f) PKG:(%s)", s.Coverage, s.PKG)
// 		}
// 	})
// }

func TestServiceQATrend(t *testing.T) {
	convey.Convey("QATrend", t, func(ctx convey.C) {
		req := &ut.QATrendReq{
			User:      "fengshanshan",
			LastTime:  30,
			Period:    "hour",
			StartTime: 1540465325,
			EndTime:   1540551725,
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			trend, err := svr.QATrend(context.Background(), req)
			t.Logf("trend:%+v\n", trend)
			ctx.Convey("Than trend should not be nil. err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(trend, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUTGernalCommit(t *testing.T) {
	convey.Convey("UTGernalCommit", t, func(ctx convey.C) {
		req := "ae1377033a11ca85a19bca365af32a5b0ebea31c,324a7609f2be27ccbedf7540f970982317d6cd6b,a1e94b169c728dab26ff583f2619b35d40519752"
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			commits, err := svr.UTGernalCommit(context.Background(), req)
			for _, commit := range commits {
				t.Logf("commit:%+v\n", commit)
				t.Logf("commit:%+v\n", commit.GitlabCommit)
			}
			ctx.Convey("Than commitInfo should not be nil. err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(commits, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestService_CommitHistory(t *testing.T) {
	convey.Convey("CommitHistory", t, func() {
		data, err := svr.CommitHistory(context.Background(), "zhaobingqing", 10)
		convey.So(data, convey.ShouldNotBeEmpty)
		convey.So(err, convey.ShouldBeNil)
		for _, d := range data {
			t.Logf("d.MergeID=%d, d.CommitID=%s, d.PKG=%s, d.Coverage=%0.2f, d.PassRate=%0.2f", d.MergeID, d.CommitID, d.PKG, d.Coverage, d.PassRate)
		}
	})
}
