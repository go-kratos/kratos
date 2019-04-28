package dao

import (
	"context"
	"go-common/app/admin/main/search/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoArchiveVideoScore(t *testing.T) {
	convey.Convey("ArchiveVideoScore", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			sp = &model.QueryParams{
				QueryBody: &model.QueryBody{},
				AppIDConf: &model.QueryConfDetail{
					ESCluster: "",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, debug, err := d.ArchiveVideoScore(c, sp)
			ctx.Convey("Then err should be nil.res,debug should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(debug, convey.ShouldNotBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoArchiveScore(t *testing.T) {
	convey.Convey("ArchiveScore", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			sp = &model.QueryParams{
				QueryBody: &model.QueryBody{},
				AppIDConf: &model.QueryConfDetail{
					ESCluster: "",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, debug, err := d.ArchiveScore(c, sp)
			ctx.Convey("Then err should be nil.res,debug should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(debug, convey.ShouldNotBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTaskQaRandom(t *testing.T) {
	convey.Convey("TaskQaRandom", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			sp = &model.QueryParams{
				QueryBody: &model.QueryBody{},
				AppIDConf: &model.QueryConfDetail{
					ESCluster: "",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, debug, err := d.TaskQaRandom(c, sp)
			ctx.Convey("Then err should be nil.res,debug should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(debug, convey.ShouldNotBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoEsportsContestsDate(t *testing.T) {
	convey.Convey("EsportsContestsDate", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			sp = &model.QueryParams{
				QueryBody: &model.QueryBody{},
				AppIDConf: &model.QueryConfDetail{
					ESCluster: "pcie_pub_out01",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, debug, err :=
			d.EsportsContestsDate(c, sp)
			ctx.Convey("Then err should be nil.res,debug should not be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
				//ctx.So(debug, convey.ShouldNotBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCreativeArchiveSearch(t *testing.T) {
	convey.Convey("CreativeArchiveSearch", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			sp = &model.QueryParams{
				QueryBody: &model.QueryBody{
					Where: &model.QueryBodyWhere{
						EQ: map[string]interface{}{"mid": 1},
					},
				},
				AppIDConf: &model.QueryConfDetail{
					ESCluster: "ssd_pub_in01",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, debug, err :=
			d.CreativeArchiveSearch(c, sp)
			ctx.Convey("Then err should be nil.res,debug should not be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
				//ctx.So(debug, convey.ShouldNotBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCreativeArchiveStaff(t *testing.T) {
	convey.Convey("CreativeArchiveStaff", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			sp = &model.QueryParams{
				QueryBody: &model.QueryBody{
					Where: &model.QueryBodyWhere{
						Combo: []model.QueryBodyWhereCombo{
							{
								EQ: []map[string]interface{}{{"mid": 1}},
							},
						},
						Like: []model.QueryBodyWhereLike{{
							KWFields: []string{"title"},
							KW:       []string{"title"},
						}},
					},
				},
				AppIDConf: &model.QueryConfDetail{
					ESCluster: "ssd_pub_in02",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, debug, err :=
			d.CreativeArchiveStaff(c, sp)
			ctx.Convey("Then err should be nil.res,debug should not be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
				//ctx.So(debug, convey.ShouldNotBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCreativeArchiveApply(t *testing.T) {
	convey.Convey("CreativeArchiveApply", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			sp = &model.QueryParams{
				QueryBody: &model.QueryBody{
					Where: &model.QueryBodyWhere{
						EQ: map[string]interface{}{"apply_staff.apply_staff_mid": "1"},
						Like: []model.QueryBodyWhereLike{{
							KWFields: []string{"title"},
							KW:       []string{"title"},
						}},
					},
				},
				AppIDConf: &model.QueryConfDetail{
					ESCluster: "ssd_pub_in02",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, debug, err :=
			d.CreativeArchiveApply(c, sp)
			ctx.Convey("Then err should be nil.res,debug should not be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
				//ctx.So(debug, convey.ShouldNotBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
