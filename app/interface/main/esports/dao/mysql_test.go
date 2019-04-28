package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaologoURL(t *testing.T) {
	convey.Convey("logoURL", t, func(ctx convey.C) {
		var (
			uri = "/bfs/archive/ad22dba8f05bb7dee6889492d3bb544413ee42c1.jpg"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			logo := logoURL(uri)
			ctx.Convey("Then logo should not be nil.", func(ctx convey.C) {
				ctx.So(logo, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMatchs(t *testing.T) {
	convey.Convey("Matchs", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Matchs(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGames(t *testing.T) {
	convey.Convey("Games", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Games(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTeams(t *testing.T) {
	convey.Convey("Teams", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Teams(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTags(t *testing.T) {
	convey.Convey("Tags", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Tags(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoYears(t *testing.T) {
	convey.Convey("Years", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Years(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCalendar(t *testing.T) {
	convey.Convey("Calendar", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			stime = int64(0)
			etime = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Calendar(c, stime, etime)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSeason(t *testing.T) {
	convey.Convey("Season", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Season(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoModule(t *testing.T) {
	convey.Convey("Module", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mmid = int64(69)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mod, err := d.Module(c, mmid)
			ctx.Convey("Then err should be nil.mod should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mod, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoModules(t *testing.T) {
	convey.Convey("Modules", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mods, err := d.Modules(c, aid)
			ctx.Convey("Then err should be nil.mods should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mods, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTrees(t *testing.T) {
	convey.Convey("Trees", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			madID = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mods, err := d.Trees(c, madID)
			ctx.Convey("Then err should be nil.mods should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mods, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoActive(t *testing.T) {
	convey.Convey("Active", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mod, err := d.Active(c, aid)
			ctx.Convey("Then err should be nil.mod should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mod, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPActDetail(t *testing.T) {
	convey.Convey("PActDetail", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.PActDetail(c, id)
			ctx.Convey("Then err should be nil.mod should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoActData(t *testing.T) {
	convey.Convey("ActDetail", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			actData, err := d.ActDetail(c, aid)
			ctx.Convey("Then err should be nil.actData should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(actData, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAppSeason(t *testing.T) {
	convey.Convey("AppSeason", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.AppSeason(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSeasonAll(t *testing.T) {
	convey.Convey("SeasonAll", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.SeasonAll(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoContest(t *testing.T) {
	convey.Convey("Contest", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			cid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Contest(c, cid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoContestRecent(t *testing.T) {
	convey.Convey("ContestRecent", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			homeid    = int64(1)
			awayid    = int64(2)
			contestid = int64(1)
			ps        = int64(2)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.ContestRecent(c, homeid, awayid, contestid, ps)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoContestData(t *testing.T) {
	convey.Convey("ContestData", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			cid = int64(458)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.ContestData(c, cid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawEpSeasons(t *testing.T) {
	convey.Convey("RawEpSeasons", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			sids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawEpSeasons(c, sids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoActTeams(t *testing.T) {
	convey.Convey("ActTeams", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ActTeams(c, tids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoKDetails(t *testing.T) {
	convey.Convey("KDetails", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.KDetails(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawEpContests(t *testing.T) {
	convey.Convey("RawEpContests", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			cids = []int64{1, 2, 3, 4, 5, 6}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.RawEpContests(c, cids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
