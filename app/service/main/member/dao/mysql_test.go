package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/member/model"
	"go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohit(t *testing.T) {
	var (
		id = int64(0)
	)
	convey.Convey("hit", t, func(ctx convey.C) {
		p1 := hit(id)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBeginTran(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		tx, err := d.BeginTran(c)
		defer tx.Commit()
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("tx should not be nil", func(ctx convey.C) {
			ctx.So(tx, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBaseInfo(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("BaseInfo", t, func(ctx convey.C) {
		r, err := d.BaseInfo(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("r should not be nil", func(ctx convey.C) {
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetBase(t *testing.T) {
	var (
		c    = context.Background()
		base = &model.BaseInfo{
			Mid: 0,
		}
	)
	convey.Convey("SetBase", t, func(ctx convey.C) {
		err := d.SetBase(c, base)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetSign(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(0)
		sign = "test"
	)
	convey.Convey("SetSign", t, func(ctx convey.C) {
		err := d.SetSign(c, mid, sign)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetName(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(0)
		name = "test"
	)
	convey.Convey("SetName", t, func(ctx convey.C) {
		err := d.SetName(c, mid, name)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetRank(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(0)
		rank = int64(100)
	)
	convey.Convey("SetRank", t, func(ctx convey.C) {
		err := d.SetRank(c, mid, rank)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetSex(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
		sex = int64(1)
	)
	convey.Convey("SetSex", t, func(ctx convey.C) {
		err := d.SetSex(c, mid, sex)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetBirthday(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(0)
		birthday = time.Time(946656000)
	)
	convey.Convey("SetBirthday", t, func(ctx convey.C) {
		err := d.SetBirthday(c, mid, birthday)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetFace(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(0)
		face = "test"
	)
	convey.Convey("SetFace", t, func(ctx convey.C) {
		err := d.SetFace(c, mid, face)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoExpDB(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("ExpDB", t, func(ctx convey.C) {
		count, err := d.ExpDB(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("count should not be nil", func(ctx convey.C) {
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetExp(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(0)
		count = int64(10)
	)
	convey.Convey("SetExp", t, func(ctx convey.C) {
		affect, err := d.SetExp(c, mid, count)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("affect should not be nil", func(ctx convey.C) {
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateExp(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(0)
		delta = int64(6)
	)
	convey.Convey("UpdateExp", t, func(ctx convey.C) {
		affect, err := d.UpdateExp(c, mid, delta)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("affect should not be nil", func(ctx convey.C) {
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUserAttrDB(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(0)
		attr uint
	)
	convey.Convey("UserAttrDB", t, func(ctx convey.C) {
		hasAttr, err := d.UserAttrDB(c, mid, attr)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("hasAttr should not be nil", func(ctx convey.C) {
			ctx.So(hasAttr, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetUserAttr(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(0)
		attr uint
	)
	convey.Convey("SetUserAttr", t, func(ctx convey.C) {
		err := d.SetUserAttr(c, mid, attr)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOfficials(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("Officials", t, func(ctx convey.C) {
		om, err := d.Officials(c)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("om should not be nil", func(ctx convey.C) {
			ctx.So(om, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoOfficial(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("Official", t, func(ctx convey.C) {
		p1, p2 := d.Official(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p2, convey.ShouldBeNil)
		})
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoMoralDB(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(4780461)
	)
	convey.Convey("MoralDB", t, func(ctx convey.C) {
		moral, err := d.MoralDB(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("moral should not be nil", func(ctx convey.C) {
			ctx.So(moral, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxMoralDB(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.db.Begin(c)
		mid   = int64(8)
	)
	defer tx.Commit()
	convey.Convey("TxMoralDB", t, func(ctx convey.C) {
		moral, err := d.TxMoralDB(tx, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("moral should not be nil", func(ctx convey.C) {
			ctx.So(moral, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxUpdateMoral(t *testing.T) {
	var (
		c        = context.Background()
		tx, _    = d.db.Begin(c)
		mid      = int64(0)
		moral    = int64(0)
		added    = int64(0)
		deducted = int64(0)
	)
	defer tx.Commit()
	convey.Convey("TxUpdateMoral", t, func(ctx convey.C) {
		err := d.TxUpdateMoral(tx, mid, moral, added, deducted)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxUpdateMoralRecoverDate(t *testing.T) {
	var (
		c           = context.Background()
		tx, _       = d.db.Begin(c)
		mid         = int64(0)
		recoverDate = time.Time(946656000)
	)
	defer tx.Commit()
	convey.Convey("TxUpdateMoralRecoverDate", t, func(ctx convey.C) {
		err := d.TxUpdateMoralRecoverDate(tx, mid, recoverDate)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxInitMoral(t *testing.T) {
	var (
		c               = context.Background()
		tx, _           = d.db.Begin(c)
		mid             = int64(10)
		moral           = int64(0)
		added           = int64(0)
		deducted        = int64(0)
		lastRecoverDate = time.Time(946656000)
	)
	defer tx.Commit()
	convey.Convey("TxInitMoral", t, func(ctx convey.C) {
		err := d.TxInitMoral(tx, mid, moral, added, deducted, lastRecoverDate)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetOfficialDoc(t *testing.T) {
	var (
		c  = context.Background()
		od = &model.OfficialDoc{
			Mid: 4780461,
		}
	)
	convey.Convey("SetOfficialDoc", t, func(ctx convey.C) {
		err := d.SetOfficialDoc(c, od)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOfficialDoc(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(4780461)
	)
	convey.Convey("OfficialDoc", t, func(ctx convey.C) {
		p1, p2 := d.OfficialDoc(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(p2, convey.ShouldBeNil)
		})
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRealnameInfo(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(46333)
	)
	convey.Convey("TestDaoRealnameInfo", t, func(ctx convey.C) {
		info, err := d.RealnameInfo(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("info should not be nil", func(ctx convey.C) {
			t.Logf("info:%+v", info)
			ctx.So(info, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRealnameInfoByCard(t *testing.T) {
	var (
		c       = context.Background()
		cardMD5 = "0088bdb8af58d25c1d9864e568a3cfb8"
	)
	convey.Convey("TestDaoRealnameInfoByCard", t, func(ctx convey.C) {
		info, err := d.RealnameInfoByCard(c, cardMD5)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("info should not be nil", func(ctx convey.C) {
			t.Logf("info:%+v", info)
			ctx.So(info, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRealnameApply(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(4780461)
	)
	convey.Convey("RealnameApply", t, func(ctx convey.C) {
		apply, err := d.RealnameApply(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("apply should not be nil", func(ctx convey.C) {
			ctx.So(apply, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInsertRealnameApply(t *testing.T) {
	var (
		c    = context.Background()
		data = &model.RealnameApply{
			MID: 4780461,
		}
	)
	convey.Convey("InsertRealnameApply", t, func(ctx convey.C) {
		err := d.InsertRealnameApply(c, data)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoInsertRealnameApplyImg(t *testing.T) {
	var (
		c    = context.Background()
		data = &model.RealnameApplyImage{
			ID:      1,
			IMGData: "1234",
		}
	)
	convey.Convey("InsertRealnameApplyImg", t, func(ctx convey.C) {
		id, err := d.InsertRealnameApplyImg(c, data)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("id should not be nil", func(ctx convey.C) {
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRealnameApplyIMG(t *testing.T) {
	var (
		c  = context.Background()
		id = int(1)
	)
	convey.Convey("RealnameApplyIMG", t, func(ctx convey.C) {
		img, err := d.RealnameApplyIMG(c, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("img should not be nil", func(ctx convey.C) {
			ctx.So(img, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInsertOldRealnameApply(t *testing.T) {
	var (
		c    = context.Background()
		data = &model.RealnameApply{
			ID: 1,
		}
	)
	convey.Convey("InsertOldRealnameApply", t, func(ctx convey.C) {
		id, err := d.InsertOldRealnameApply(c, data)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("id should not be nil", func(ctx convey.C) {
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInsertOldRealnameApplyImg(t *testing.T) {
	var (
		c    = context.Background()
		data = &model.RealnameApplyImage{
			ID:      1,
			IMGData: "1234",
		}
	)
	convey.Convey("InsertOldRealnameApplyImg", t, func(ctx convey.C) {
		id, err := d.InsertOldRealnameApplyImg(c, data)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("id should not be nil", func(ctx convey.C) {
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}
