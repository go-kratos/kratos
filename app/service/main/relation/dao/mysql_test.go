package dao

import (
	"context"
	"go-common/app/service/main/relation/model"
	"go-common/app/service/main/relation/model/i64b"
	"math/rand"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestDaohit(t *testing.T) {
	var (
		id = int64(1)
	)
	convey.Convey("hit", t, func(cv convey.C) {
		p1 := hit(id)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaostatHit(t *testing.T) {
	var (
		id = int64(1)
	)
	convey.Convey("statHit", t, func(cv convey.C) {
		p1 := statHit(id)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaotagHit(t *testing.T) {
	var (
		id = int64(1)
	)
	convey.Convey("tagHit", t, func(cv convey.C) {
		p1 := tagHit(id)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaotagUserHit(t *testing.T) {
	var (
		id = int64(1)
	)
	convey.Convey("tagUserHit", t, func(cv convey.C) {
		p1 := tagUserHit(id)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBeginTran(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("BeginTran", t, func(cv convey.C) {
		p1, err := d.BeginTran(c)
		cv.Convey("Then err should be nil.p1 should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(p1, convey.ShouldNotBeNil)
		})
		p1.Commit()
	})
}

func TestDaoFollowings(t *testing.T) {
	var (
		c      = context.Background()
		tx, _  = d.BeginTran(c)
		mid    = int64(1)
		fid    = int64(2)
		mask   = uint32(2)
		source = uint8(1)
		now    = time.Now()
	)
	convey.Convey("Followings", t, func(cv convey.C) {
		affected, err := d.TxAddFollowing(c, tx, mid, fid, mask, source, now)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
		tx.Commit()

		res, err := d.Followings(c, mid)
		cv.Convey("Then err should be nil.res should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoFollowingsIn(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(1)
		fids = []int64{1, 2}
	)
	convey.Convey("FollowingsIn", t, func(cv convey.C) {
		res, err := d.FollowingsIn(c, mid, fids)
		cv.Convey("Then err should be nil.res should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxSetFollowing(t *testing.T) {
	var (
		c         = context.Background()
		tx, _     = d.BeginTran(c)
		mid       = int64(1)
		fid       = int64(2)
		attribute = uint32(2)
		source    = uint8(1)
		status    = int(0)
		now       = time.Now()
	)
	convey.Convey("TxSetFollowing", t, func(cv convey.C) {
		affected, err := d.TxSetFollowing(c, tx, mid, fid, attribute, source, status, now)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
		tx.Commit()
	})
}

func TestDaoFollowers(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("Followers", t, func(cv convey.C) {
		res, err := d.Followers(c, mid)
		cv.Convey("Then err should be nil.res should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxAddFollower(t *testing.T) {
	var (
		c      = context.Background()
		tx, _  = d.BeginTran(c)
		mid    = int64(1)
		fid    = int64(2)
		mask   = uint32(2)
		source = uint8(1)
		now    = time.Now()
	)
	convey.Convey("TxAddFollower", t, func(cv convey.C) {
		affected, err := d.TxAddFollower(c, tx, mid, fid, mask, source, now)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
		tx.Commit()
	})
}

func TestDaoTxSetFollower(t *testing.T) {
	var (
		c         = context.Background()
		tx, _     = d.BeginTran(c)
		mid       = int64(1)
		fid       = int64(2)
		attribute = uint32(2)
		source    = uint8(1)
		status    = int(0)
		now       = time.Now()
	)
	convey.Convey("TxSetFollower", t, func(cv convey.C) {
		affected, err := d.TxSetFollower(c, tx, mid, fid, attribute, source, status, now)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
		tx.Commit()
	})
}

func TestDaoStat(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("Stat", t, func(cv convey.C) {
		stat, err := d.Stat(c, mid)
		cv.Convey("Then err should be nil.stat should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(stat, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxStat(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		mid   = int64(1)
	)
	convey.Convey("TxStat", t, func(cv convey.C) {
		stat, err := d.TxStat(c, tx, mid)
		cv.Convey("Then err should be nil.stat should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(stat, convey.ShouldNotBeNil)
		})
		tx.Commit()
	})
}

func TestDaoAddStat(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(1)
		stat = &model.Stat{
			Mid: 1,
		}
		now = time.Now()
	)
	convey.Convey("AddStat", t, func(cv convey.C) {
		affected, err := d.AddStat(c, mid, stat, now)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxAddStat(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		mid   = int64(1)
		stat  = &model.Stat{
			Mid: 1,
		}
		now = time.Now()
	)
	convey.Convey("TxAddStat", t, func(cv convey.C) {
		affected, err := d.TxAddStat(c, tx, mid, stat, now)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
		tx.Commit()
	})
}

func TestDaoTxSetStat(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		mid   = int64(1)
		stat  = &model.Stat{
			Mid: 1,
		}
		now = time.Now()
	)
	convey.Convey("TxSetStat", t, func(cv convey.C) {
		affected, err := d.TxSetStat(c, tx, mid, stat, now)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
		tx.Commit()
	})
}

func TestDaoRelation(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
		fid = int64(2)
	)
	convey.Convey("Relation", t, func(cv convey.C) {
		attr, err := d.Relation(c, mid, fid)
		cv.Convey("Then err should be nil.attr should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(attr, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoLoadMonitor(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("LoadMonitor", t, func(cv convey.C) {
		affected, err := d.AddMonitor(c, 1, time.Now())
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})

		mids, err := d.LoadMonitor(c)
		cv.Convey("Then err should be nil.mids should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(mids, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelMonitor(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("DelMonitor", t, func(cv convey.C) {
		affected, err := d.DelMonitor(c, mid)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxDelTagUser(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		mid   = int64(1)
		fid   = int64(2)
	)
	convey.Convey("TxDelTagUser", t, func(cv convey.C) {
		affected, err := d.TxDelTagUser(c, tx, mid, fid)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
		tx.Commit()
	})
}

func TestDaoTags(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("Tags", t, func(cv convey.C) {
		res, err := d.Tags(c, mid)
		cv.Convey("Then err should be nil.res should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelTag(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
		id  = int64(2)
	)
	convey.Convey("DelTag", t, func(cv convey.C) {
		affected, err := d.DelTag(c, mid, id)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetTagName(t *testing.T) {
	var (
		c    = context.Background()
		id   = int64(1)
		mid  = int64(2)
		name = "test"
		now  = time.Now()
	)
	convey.Convey("SetTagName", t, func(cv convey.C) {
		affected, err := d.SetTagName(c, id, mid, name, now)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTagUserByMidFid(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
		fid = int64(2)
	)
	convey.Convey("TagUserByMidFid", t, func(cv convey.C) {
		lastID, err := d.AddTag(c, mid, fid, "test"+RandStringRunes(5), time.Now())
		cv.Convey("AddTag; Then err should be nil.lastID should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(lastID, convey.ShouldNotBeNil)
		})

		tids := i64b.Int64Bytes([]int64{lastID})
		affected, err := d.SetTagUser(c, mid, fid, tids, time.Now())
		cv.Convey("SetTagUser; Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})

		affected2, err := d.AddTagUser(c, mid, fid, tids, time.Now())
		cv.Convey("AddTagUser; Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected2, convey.ShouldNotBeNil)
		})

		tag1, err := d.TagUserByMidFid(c, mid, fid)
		cv.Convey("TagUserByMidFid; Then err should be nil.tag1 should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(tag1, convey.ShouldNotBeNil)
		})

		tags2, err := d.UsersTags(c, mid, []int64{fid})
		cv.Convey("UsersTags; Then err should be nil.tags2 should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(tags2, convey.ShouldNotBeNil)
		})

		tags3, err := d.UserTag(c, mid)
		cv.Convey("UserTag; Then err should be nil.tags3 should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(tags3, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoHasReachAchieve(t *testing.T) {
	var (
		c       = context.Background()
		mid     = int64(1)
		achieve model.AchieveFlag
	)
	convey.Convey("HasReachAchieve", t, func(cv convey.C) {
		p1 := d.HasReachAchieve(c, mid, achieve)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoFollowerNotifySetting(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("FollowerNotifySetting", t, func(cv convey.C) {
		p1, err := d.FollowerNotifySetting(c, mid)
		cv.Convey("Then err should be nil.p1 should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoEnableFollowerNotify(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("EnableFollowerNotify", t, func(cv convey.C) {
		affected, err := d.EnableFollowerNotify(c, mid)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDisableFollowerNotify(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("DisableFollowerNotify", t, func(cv convey.C) {
		affected, err := d.DisableFollowerNotify(c, mid)
		cv.Convey("Then err should be nil.affected should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(affected, convey.ShouldNotBeNil)
		})
	})
}
