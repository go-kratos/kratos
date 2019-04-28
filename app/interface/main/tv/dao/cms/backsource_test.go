package cms

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/interface/main/tv/model"
	"go-common/library/database/sql"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmsVideoMetaDB(t *testing.T) {
	var (
		c   = context.Background()
		cid = int64(0)
	)
	convey.Convey("VideoMetaDB", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.meta should not be nil.", func(ctx convey.C) {
			sids, err := pickIDs(d.db, _pickCids)
			if err != nil || len(sids) == 0 {
				fmt.Println("Empty Sids ", err)
				return
			}
			cid = sids[0]
			meta, err := d.VideoMetaDB(c, cid)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(meta, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err := d.VideoMetaDB(c, cid)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsArcMetaDB(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(0)
	)
	convey.Convey("ArcMetaDB", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.meta should not be nil.", func(ctx convey.C) {
			sids, err := pickIDs(d.db, _pickAids)
			if err != nil || len(sids) == 0 {
				fmt.Println("Empty Sids ", err)
				return
			}
			aid = sids[0]
			meta, err := d.ArcMetaDB(c, aid)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(meta, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err := d.ArcMetaDB(c, aid)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsVideoMetas(t *testing.T) {
	var (
		c    = context.Background()
		cids = []int64{}
		err  error
	)
	convey.Convey("VideoMetas", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.meta should not be nil.", func(ctx convey.C) {
			if cids, err = pickIDs(d.db, _pickCids); err != nil || len(cids) == 0 {
				fmt.Println("Empty Sids ", err)
				return
			}
			meta, err := d.VideoMetas(c, cids)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(meta, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err := d.VideoMetas(c, cids)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsArcMetas(t *testing.T) {
	var (
		c    = context.Background()
		aids = []int64{}
		err  error
	)
	convey.Convey("ArcMetas", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.metas should not be nil.", func(ctx convey.C) {
			if aids, err = pickIDs(d.db, _pickAids); err != nil || len(aids) == 0 {
				fmt.Println("Empty Sids ", err)
				return
			}
			metas, err := d.ArcMetas(c, aids)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(metas, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err := d.ArcMetas(c, aids)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsSeasonMetas(t *testing.T) {
	var (
		c    = context.Background()
		sids = []int64{}
		err  error
	)
	convey.Convey("SeasonMetas", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.metas should not be nil.", func(ctx convey.C) {
			if sids, err = pickIDs(d.db, _pickSids); err != nil || len(sids) == 0 {
				fmt.Println("Empty Sids ", err)
				return
			}
			metas, err := d.SeasonMetas(c, sids)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(metas, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err := d.SeasonMetas(c, sids)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsNewestOrder(t *testing.T) {
	var (
		c   = context.Background()
		sid = int64(0)
	)
	convey.Convey("NewestOrder", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.epid,newestOrder should not be nil.", func(ctx convey.C) {
			epid, newestOrder, err := d.NewestOrder(c, sid)
			if err != nil {
				ctx.So(err, convey.ShouldEqual, sql.ErrNoRows)
			} else {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(newestOrder, convey.ShouldNotBeNil)
				ctx.So(epid, convey.ShouldNotBeNil)
			}
		})
	})
}

func TestCmsEpMetas(t *testing.T) {
	var (
		c     = context.Background()
		epids = []int64{}
		err   error
	)
	convey.Convey("EpMetas", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.metas should not be nil.", func(ctx convey.C) {
			if epids, err = pickIDs(d.db, _pickEpids); err != nil || len(epids) == 0 {
				fmt.Println("Empty epids ", err)
				return
			}
			metas, err := d.EpMetas(c, epids)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(metas, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err := d.EpMetas(c, epids)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsEpAuthDB(t *testing.T) {
	var (
		c    = context.Background()
		epid = int64(0)
	)
	convey.Convey("EpAuthDB", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.ep should not be nil.", func(ctx convey.C) {
			epids, err := pickIDs(d.db, _pickEpids)
			if err != nil || len(epids) == 0 {
				fmt.Println("Empty epids ", err)
				return
			}
			epid = epids[0]
			ep, err := d.EpAuthDB(c, epid)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ep, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err := d.EpAuthDB(c, epid)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsSnAuthDB(t *testing.T) {
	var (
		c    = context.Background()
		sids []int64
		sid  int64
		err  error
	)
	convey.Convey("SnAuthDB", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.s should not be nil.", func(ctx convey.C) {
			if sids, err = pickIDs(d.db, _pickSids); err != nil || len(sids) == 0 {
				fmt.Println("Empty Sids ", err)
				return
			}
			sid = sids[0]
			s, err := d.SnAuthDB(c, sid)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(s, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err := d.SnAuthDB(c, sid)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsSnsAuthDB(t *testing.T) {
	var (
		c       = context.Background()
		sids    []int64
		err     error
		snsAuth map[int64]*model.SnAuth
	)
	convey.Convey("SnsAuthDB", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.snsAuth should not be nil.", func(ctx convey.C) {
			if sids, err = pickIDs(d.db, _pickSids); err != nil || len(sids) == 0 {
				fmt.Println("Empty Sids ", err)
				return
			}
			snsAuth, err = d.SnsAuthDB(c, sids)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(snsAuth, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err = d.SnsAuthDB(c, sids)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsEpsAuthDB(t *testing.T) {
	var (
		c       = context.Background()
		epids   []int64
		err     error
		epsAuth map[int64]*model.EpAuth
	)
	convey.Convey("EpsAuthDB", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.epsAuth should not be nil.", func(ctx convey.C) {
			epids, err = pickIDs(d.db, _pickEpids)
			if err != nil || len(epids) == 0 {
				fmt.Println("Empty epids ", err)
				return
			}
			epsAuth, err = d.EpsAuthDB(c, epids)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(epsAuth, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err = d.EpsAuthDB(c, epids)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsSeasonCMS(t *testing.T) {
	var (
		c      = context.Background()
		sids   []int64
		sid    = int64(0)
		err    error
		season *model.SeasonCMS
	)
	convey.Convey("SeasonCMS", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.season should not be nil.", func(ctx convey.C) {
			sids, err = pickIDs(d.db, _pickSids)
			if err != nil || len(sids) == 0 {
				fmt.Println("Empty Sids ", err)
				return
			}
			sid = sids[0]
			season, err = d.SeasonCMS(c, sid)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(season, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err = d.SeasonCMS(c, sid)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsEpCMS(t *testing.T) {
	var (
		c    = context.Background()
		epid = int64(0)
	)
	convey.Convey("EpCMS", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.ep should not be nil.", func(ctx convey.C) {
			epids, err := pickIDs(d.db, _pickEpids)
			if err != nil || len(epids) == 0 {
				fmt.Println("Empty epids ", err)
				return
			}
			epid = epids[0]
			ep, err := d.EpCMS(c, epid)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ep, convey.ShouldNotBeNil)
		})
		ctx.Convey("Db Error", func(ctx convey.C) {
			d.db.Close()
			_, err := d.EpCMS(c, epid)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}
