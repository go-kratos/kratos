package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/passport-user/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoubKey(t *testing.T) {
	convey.Convey("ubKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := ubKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoutKey(t *testing.T) {
	convey.Convey("utKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := utKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoueKey(t *testing.T) {
	convey.Convey("ueKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := ueKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaouroKey(t *testing.T) {
	convey.Convey("uroKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := uroKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaousqKey(t *testing.T) {
	convey.Convey("usqKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := usqKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoqqKey(t *testing.T) {
	convey.Convey("qqKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := qqKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosinaKey(t *testing.T) {
	convey.Convey("sinaKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := sinaKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetUserBaseCache(t *testing.T) {
	convey.Convey("SetUserBaseCache", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ub = &model.UserBase{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUserBaseCache(c, ub)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetUserTelCache(t *testing.T) {
	convey.Convey("SetUserTelCache", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ut = &model.UserTel{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUserTelCache(c, ut)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetUserEmailCache(t *testing.T) {
	convey.Convey("SetUserEmailCache", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ue = &model.UserEmail{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUserEmailCache(c, ue)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetUserRegOriginCache(t *testing.T) {
	convey.Convey("SetUserRegOriginCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uro = &model.UserRegOrigin{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUserRegOriginCache(c, uro)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetUserSafeQuestionCache(t *testing.T) {
	convey.Convey("SetUserSafeQuestionCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			usq = &model.UserSafeQuestion{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUserSafeQuestionCache(c, usq)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetUserThirdBindQQCache(t *testing.T) {
	convey.Convey("SetUserThirdBindQQCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			utb = &model.UserThirdBind{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUserThirdBindQQCache(c, utb)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetUserThirdBindSinaCache(t *testing.T) {
	convey.Convey("SetUserThirdBindSinaCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			utb = &model.UserThirdBind{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUserThirdBindSinaCache(c, utb)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelUserBaseCache(t *testing.T) {
	convey.Convey("DelUserBaseCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelUserBaseCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelUserTelCache(t *testing.T) {
	convey.Convey("DelUserTelCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelUserTelCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelUserEmailCache(t *testing.T) {
	convey.Convey("DelUserEmailCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelUserEmailCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaopingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingMC(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
