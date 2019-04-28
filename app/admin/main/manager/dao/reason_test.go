package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/manager/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddCateSecExt(t *testing.T) {
	var (
		c = context.TODO()
		e = &model.CateSecExt{
			BusinessID: 1,
			Type:       1,
			Name:       "测试type1",
		}
	)
	convey.Convey("AddCateSecExt", t, func(ctx convey.C) {
		err := d.AddCateSecExt(c, e)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateCateSecExt(t *testing.T) {
	var (
		c = context.TODO()
		e = &model.CateSecExt{
			ID:   75,
			Name: "测试更新下",
		}
	)
	convey.Convey("UpdateCateSecExt", t, func(ctx convey.C) {
		err := d.UpdateCateSecExt(c, e)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBanCateSecExt(t *testing.T) {
	var (
		c = context.TODO()
		e = &model.CateSecExt{
			ID:    75,
			State: 0,
		}
	)
	convey.Convey("BanCateSecExt", t, func(ctx convey.C) {
		err := d.BanCateSecExt(c, e)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddAssociation(t *testing.T) {
	var (
		c = context.TODO()
		e = &model.Association{
			BusinessID: 1,
			RoleID:     121,
			CategoryID: 212,
			SecondIDs:  "3,7",
		}
	)
	convey.Convey("AddAssociation", t, func(ctx convey.C) {
		err := d.AddAssociation(c, e)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateAssociation(t *testing.T) {
	var (
		c = context.TODO()
		e = &model.Association{
			ID:         41,
			RoleID:     1111,
			CategoryID: 2222,
			SecondIDs:  "6,5,4,3",
		}
	)
	convey.Convey("UpdateAssociation", t, func(ctx convey.C) {
		err := d.UpdateAssociation(c, e)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBanAssociation(t *testing.T) {
	var (
		c = context.TODO()
		e = &model.Association{
			ID:    41,
			State: 1,
		}
	)
	convey.Convey("BanAssociation", t, func(ctx convey.C) {
		err := d.BanAssociation(c, e)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddReason(t *testing.T) {
	var (
		c = context.TODO()
		e = &model.Reason{
			BusinessID:  1,
			RoleID:      99,
			CategoryID:  98,
			SecondID:    97,
			Common:      0,
			UID:         1,
			Description: "随便测试",
			Weight:      96,
			Flag:        0,
		}
	)
	convey.Convey("AddReason", t, func(ctx convey.C) {
		err := d.AddReason(c, e)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateReason(t *testing.T) {
	var (
		c = context.TODO()
		e = &model.Reason{
			ID:          18,
			RoleID:      100,
			CategoryID:  99,
			SecondID:    98,
			Common:      1,
			Description: "随便测试试试",
			Weight:      97,
			Flag:        1,
			BusinessID:  2,
			TypeID:      3,
			TagID:       4,
		}
	)
	convey.Convey("UpdateReason", t, func(ctx convey.C) {
		err := d.UpdateReason(c, e)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoReasonList(t *testing.T) {
	var (
		c = context.TODO()
		e = &model.SearchReasonParams{
			BusinessID: 1,
		}
	)
	convey.Convey("ReasonList", t, func(ctx convey.C) {
		res, err := d.ReasonList(c, e)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCateSecByIDs(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{1, 2, 3}
	)
	convey.Convey("CateSecByIDs", t, func(ctx convey.C) {
		res, err := d.CateSecByIDs(c, ids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBatchUpdateReasonState(t *testing.T) {
	var (
		c = context.TODO()
		b = &model.BatchUpdateReasonState{
			IDs: []int64{1, 2, 3},
		}
	)
	convey.Convey("BatchUpdateReasonState", t, func(ctx convey.C) {
		err := d.BatchUpdateReasonState(c, b)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCateSecExtList(t *testing.T) {
	var (
		c = context.TODO()
		e = &model.CateSecExt{
			BusinessID: 1,
			Type:       1,
		}
	)
	convey.Convey("CateSecExtList", t, func(ctx convey.C) {
		res, err := d.CateSecExtList(c, e)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCateSecList(t *testing.T) {
	var (
		c   = context.TODO()
		lid = int64(0)
	)
	convey.Convey("CateSecList", t, func(ctx convey.C) {
		res, err := d.CateSecList(c, lid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAssociationList(t *testing.T) {
	var (
		c     = context.TODO()
		state = int64(0)
		lid   = int64(0)
	)
	convey.Convey("AssociationList", t, func(ctx convey.C) {
		res, err := d.AssociationList(c, state, lid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
