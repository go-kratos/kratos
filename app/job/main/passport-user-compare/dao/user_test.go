package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/passport-user-compare/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoQueryUserBase(t *testing.T) {
	convey.Convey("QueryUserBase", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryUserBase(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateUserBase(t *testing.T) {
	convey.Convey("UpdateUserBase", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserBase{
				Mid:    1111111,
				UserID: "test_0000_0001",
				Pwd:    []byte{},
				Salt:   "",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.UpdateUserBase(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUserBase(t *testing.T) {
	convey.Convey("InsertUserBase", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserBase{
				Mid:    1111111111,
				UserID: "test_0000_0002",
				Pwd:    []byte{},
				Salt:   "",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.InsertUserBase(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryUserTel(t *testing.T) {
	convey.Convey("QueryUserTel", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1111111111)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryUserTel(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateUserTel(t *testing.T) {
	convey.Convey("UpdateUserTel", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserTel{
				Mid: 1111111111,
				Tel: []byte{},
				Cid: "",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.UpdateUserTel(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUserTel(t *testing.T) {
	convey.Convey("InsertUserTel", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserTel{
				Mid: 1111111111,
				Tel: []byte{},
				Cid: "",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.InsertUserTel(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryUserMail(t *testing.T) {
	convey.Convey("QueryUserMail", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1111111111)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryUserMail(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateUserMail(t *testing.T) {
	convey.Convey("UpdateUserMail", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserEmail{
				Mid:   1111111111,
				Email: []byte{},
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.UpdateUserMail(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateUserMailVerified(t *testing.T) {
	convey.Convey("UpdateUserMailVerified", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserEmail{
				Mid:   1111111111,
				Email: []byte{},
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.UpdateUserMailVerified(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUserEmail(t *testing.T) {
	convey.Convey("InsertUserEmail", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserEmail{
				Mid:   1111111111,
				Email: []byte{},
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.InsertUserEmail(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryUserSafeQuestion(t *testing.T) {
	convey.Convey("QueryUserSafeQuestion", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1111111111)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryUserSafeQuestion(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateUserSafeQuestion(t *testing.T) {
	convey.Convey("UpdateUserSafeQuestion", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserSafeQuestion{
				Mid: 1111111111,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.UpdateUserSafeQuestion(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUserSafeQuestion(t *testing.T) {
	convey.Convey("InsertUserSafeQuestion", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserSafeQuestion{
				Mid: 1111111111,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.InsertUserSafeQuestion(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryUserThirdBind(t *testing.T) {
	convey.Convey("QueryUserThirdBind", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(1111111111)
			platform = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryUserThirdBind(c, mid, platform)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateUserThirdBind(t *testing.T) {
	convey.Convey("UpdateUserThirdBind", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserThirdBind{
				Mid: 1111111111,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.UpdateUserThirdBind(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUserThirdBind(t *testing.T) {
	convey.Convey("InsertUserThirdBind", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserThirdBind{
				Mid: 1111111111,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.InsertUserThirdBind(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryCountryCode(t *testing.T) {
	convey.Convey("QueryCountryCode", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryCountryCode(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetMidByTel(t *testing.T) {
	convey.Convey("GetMidByTel", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserTel{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mid, err := d.GetMidByTel(c, a)
			ctx.Convey("Then err should be nil.mid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(mid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetMidByEmail(t *testing.T) {
	convey.Convey("GetMidByEmail", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserEmail{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mid, err := d.GetMidByEmail(c, a)
			ctx.Convey("Then err should be nil.mid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(mid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUnverifiedEmail(t *testing.T) {
	convey.Convey("GetUnverifiedEmail", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetUnverifiedEmail(c, start)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUserRegOriginByMid(t *testing.T) {
	convey.Convey("GetUserRegOriginByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetUserRegOriginByMid(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUpdateUserRegOriginType(t *testing.T) {
	convey.Convey("InsertUpdateUserRegOriginType", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.UserRegOrigin{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.InsertUpdateUserRegOriginType(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaotableIndex(t *testing.T) {
	convey.Convey("tableIndex", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := tableIndex(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
