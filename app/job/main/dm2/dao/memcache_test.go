package dao

import (
	"context"
	"go-common/app/job/main/dm2/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyXML(t *testing.T) {
	convey.Convey("keyXML", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyXML(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeySubject(t *testing.T) {
	convey.Convey("keySubject", t, func(ctx convey.C) {
		var (
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keySubject(tp, oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyAjax(t *testing.T) {
	convey.Convey("keyAjax", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyAjax(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyTransferLock(t *testing.T) {
	convey.Convey("keyTransferLock", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyTransferLock()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelXMLCache(t *testing.T) {
	convey.Convey("DelXMLCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.DelXMLCache(c, oid)
		})
	})
}

func TestDaoAddXMLCache(t *testing.T) {
	convey.Convey("AddXMLCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			value = []byte("")
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := testDao.AddXMLCache(c, oid, value)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoXMLCache(t *testing.T) {
	convey.Convey("XMLCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.XMLCache(c, oid)
		})
	})
}

func TestDaoSubjectCache(t *testing.T) {
	convey.Convey("SubjectCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.SubjectCache(c, tp, oid)
		})
	})
}

func TestDaoSubjectsCache(t *testing.T) {
	convey.Convey("SubjectsCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			oids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.SubjectsCache(c, tp, oids)
		})
	})
}

func TestDaoAddSubjectCache(t *testing.T) {
	convey.Convey("AddSubjectCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sub = &model.Subject{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := testDao.AddSubjectCache(c, sub)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelSubjectCache(t *testing.T) {
	convey.Convey("DelSubjectCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.DelSubjectCache(c, tp, oid)
		})
	})
}

func TestDaoAddTransferLock(t *testing.T) {
	convey.Convey("AddTransferLock", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			succeed := testDao.AddTransferLock(c)
			ctx.Convey("Then succeed should not be nil.", func(ctx convey.C) {
				ctx.So(succeed, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelTransferLock(t *testing.T) {
	convey.Convey("DelTransferLock", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.DelTransferLock(c)
		})
	})
}

func TestDaoDelAjaxDMCache(t *testing.T) {
	convey.Convey("DelAjaxDMCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.DelAjaxDMCache(c, oid)
		})
	})
}
