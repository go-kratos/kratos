package search

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/feed/conf"
	searchModel "go-common/app/admin/main/feed/model/search"
	"go-common/library/log"

	"github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c = context.Background()
)

func init() {
	var (
		err error
	)
	dir, _ := filepath.Abs("../../cmd/feed-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	if s = New(conf.Conf); err != nil {
		log.Error("bgmdao.New error(%v)", err)
		return
	}
}

func TestIsTodayAutoPubHot(t *testing.T) {
	convey.Convey("isTodayAutoPubHot", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			err error
			res bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err = s.isTodayAutoPubHot(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIsTodayAutoPubDark(t *testing.T) {
	convey.Convey("isTodayAutoPubDark", t, func(ctx convey.C) {
		var (
			err error
			res bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err = s.isTodayAutoPubDark(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestParseTime(t *testing.T) {
	convey.Convey("parseTime", t, func(ctx convey.C) {
		var (
			err error
			res time.Time
		)
		timeTwelve := time.Now().Format("2006-01-02 ") + "12:00:00"
		layout := "2006-01-02 15:04:05"
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err = s.parseTime(timeTwelve, layout)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestHotwordFromDB(t *testing.T) {
	convey.Convey("HotwordFromDB", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Intervene
		)

		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, _, err = s.HotwordFromDB("2018-09-05")
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestDarkwordFromDB(t *testing.T) {
	convey.Convey("HotwordFromDB", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Dark
		)

		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, _, err = s.DarkwordFromDB(time.Now().Format("2006-01-02"))
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestOpenHotList(t *testing.T) {
	convey.Convey("HotwordFromDB", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Intervene
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//c = context.Background()
			//res, err = s.OpenHotList(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestHotList(t *testing.T) {
	convey.Convey("HotList", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Intervene
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, err = s.HotList(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestDarkList(t *testing.T) {
	convey.Convey("DarkList", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Intervene
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, err = s.DarkList(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestBlackList(t *testing.T) {
	convey.Convey("BlackList", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Black
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err = s.BlackList()
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestDelBlack(t *testing.T) {
	convey.Convey("DelBlack", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Black
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, err = s.DelBlack()
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestAddBlack(t *testing.T) {
	convey.Convey("AddBlack", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Black
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//err = s.AddBlack()
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestCheckBlack(t *testing.T) {
	convey.Convey("checkBlack", t, func(ctx convey.C) {
		var (
			err error
			res bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err = s.checkBlack("test")
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestCheckInter(t *testing.T) {
	convey.Convey("checkInter", t, func(ctx convey.C) {
		var (
			err error
			res bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err = s.checkInter("test", 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestCheckTimeConflict(t *testing.T) {
	convey.Convey("checkTimeConflict", t, func(ctx convey.C) {
		var (
			err   error
			res   bool
			model = searchModel.InterveneAdd{
				Rank:  10,
				Stime: 1536134791,
				Etime: 1536134791,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err = s.checkTimeConflict(model, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestAddInter(t *testing.T) {
	convey.Convey("checkTimeConflict", t, func(ctx convey.C) {
		var (
			err error
			res bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, err = s.AddInter(model, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestUpdateInter(t *testing.T) {
	convey.Convey("UpdateInter", t, func(ctx convey.C) {
		var (
			err error
			res bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, err = s.AddInter(model, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestUpdateSearch(t *testing.T) {
	convey.Convey("UpdateSearch", t, func(ctx convey.C) {
		var (
			err error
			res bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, err = s.AddInter(model, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestDeleteHot(t *testing.T) {
	convey.Convey("DeleteHot", t, func(ctx convey.C) {
		var (
			err error
			res bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err = s.DeleteHot(c, 10, 2, "quguolin", 100)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestDeleteDark(t *testing.T) {
	convey.Convey("DeleteDark", t, func(ctx convey.C) {
		var (
			err error
			res bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err = s.DeleteDark(c, 10, "quguolin", 100)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestOpenAddDarkword(t *testing.T) {
	convey.Convey("OpenAddDarkword", t, func(ctx convey.C) {
		var (
			err   error
			res   bool
			value searchModel.OpenDark
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err = s.OpenAddDarkword(c, value)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestOpenAddHotword(t *testing.T) {
	convey.Convey("OpenAddHotword", t, func(ctx convey.C) {
		var (
			err error
			res bool
			//value searchModel.OpenHot
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//err = s.OpenAddHotword(c, value)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestGetHotPub(t *testing.T) {
	convey.Convey("GetHotPub", t, func(ctx convey.C) {
		var (
			err error
			res bool
			//value searchModel.OpenHot
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//err = s.OpenAddHotword(c, value)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestGetDarkPub(t *testing.T) {
	convey.Convey("GetDarkPub", t, func(ctx convey.C) {
		var (
			err error
			res bool
			//value searchModel.OpenHot
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//err = s.GetDarkPub(c, value)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestSetHotPub(t *testing.T) {
	convey.Convey("SetHotPub", t, func(ctx convey.C) {
		var (
			err error
			res bool
			//value searchModel.OpenHot
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err = s.SetHotPub(c, "quguolin", 100)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestHotPubLog(t *testing.T) {
	convey.Convey("HotPubLog", t, func(ctx convey.C) {
		var (
			err   error
			res   bool
			value []searchModel.Intervene
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err = s.HotPubLog(value)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestGetHotPubLog(t *testing.T) {
	convey.Convey("GetHotPubLog", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Intervene
			pub bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, pub, err = s.GetHotPubLog("2018-09-03")
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res, pub)
		})
	})
}

func TestGetDarkPubLog(t *testing.T) {
	convey.Convey("GetDarkPubLog", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Dark
			pub bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, pub, err = s.GetDarkPubLog("2018-09-03")
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res, pub)
		})
	})
}

func TestSetDarkPub(t *testing.T) {
	convey.Convey("SetDarkPub", t, func(ctx convey.C) {
		var (
			err error
			res []searchModel.Intervene
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.SetDarkPub(c, "quguolin", 10)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
			fmt.Println(res)
		})
	})
}

func TestDarkPubLog(t *testing.T) {
	convey.Convey("DarkPubLog", t, func(ctx convey.C) {
		var (
			err  error
			dark []searchModel.Dark
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.DarkPubLog(dark)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
