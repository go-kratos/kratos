package reply

import (
	"context"
	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/database/sql"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewSubjectDao(t *testing.T) {
	convey.Convey("NewSubjectDao", t, func(ctx convey.C) {
		var (
			db = &sql.DB{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dao := NewSubjectDao(db)
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(dao, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplySubjecthit(t *testing.T) {
	convey.Convey("hit", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.Subject.hit(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplySet(t *testing.T) {
	convey.Convey("Set", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sub = &reply.Subject{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.Subject.Set(c, sub)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyInsert(t *testing.T) {
	convey.Convey("Insert", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sub = &reply.Subject{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.Subject.Insert(c, sub)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyUpState(t *testing.T) {
	convey.Convey("UpState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			tp    = int8(0)
			state = int8(0)
			now   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.Subject.UpState(c, oid, tp, state, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyUpMid(t *testing.T) {
	convey.Convey("UpMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			oid = int64(0)
			tp  = int8(0)
			now = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.Subject.UpMid(c, mid, oid, tp, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGet(t *testing.T) {
	convey.Convey("Get", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sub, err := d.Subject.Get(c, oid, tp)
			ctx.Convey("Then err should be nil.sub should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sub, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGets(t *testing.T) {
	convey.Convey("Gets", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oids = []int64{}
			tp   = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Subject.Gets(c, oids, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
