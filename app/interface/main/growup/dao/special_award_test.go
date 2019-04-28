package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPastAwards(t *testing.T) {
	convey.Convey("PastAwards", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Exec(c, "INSERT INTO special_award(award_id,award_name,open_status) VALUES(666,'test','2') ON DUPLICATE KEY UPDATE award_id=VALUES(award_id)")
			as, err := d.PastAwards(c)
			convCtx.Convey("Then err should be nil.as should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(as, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoJoinedSpecialAwards(t *testing.T) {
	convey.Convey("JoinedSpecialAwards", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			awardIDs = []int64{666}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Exec(c, "INSERT INTO special_award(award_id,award_name,open_status) VALUES(666,'test','2') ON DUPLICATE KEY UPDATE award_id=VALUES(award_id)")
			sas, err := d.JoinedSpecialAwards(c, awardIDs)
			convCtx.Convey("Then err should be nil.sas should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(sas, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAwardSchedule(t *testing.T) {
	convey.Convey("GetAwardSchedule", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			awardID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Exec(c, "INSERT INTO special_award(award_id,award_name,open_status) VALUES(666,'test','2') ON DUPLICATE KEY UPDATE award_id=VALUES(award_id)")
			award, err := d.GetAwardSchedule(c, awardID)
			convCtx.Convey("Then err should be nil.award should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(award, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetResources(t *testing.T) {
	convey.Convey("GetResources", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			awardID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Exec(c, "INSERT INTO special_award_resource(award_id,resource_type,content,resource_index) VALUES(666,1,'test',1) ON DUPLICATE KEY UPDATE award_id=VALUES(award_id)")
			res, err := d.GetResources(c, awardID)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetWinners(t *testing.T) {
	convey.Convey("GetWinners", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			awardID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Exec(c, "INSERT INTO special_award_winner(award_id,mid,division_name) VALUES(666,2,'test') ON DUPLICATE KEY UPDATE award_id=VALUES(award_id)")
			mids, err := d.GetWinners(c, awardID)
			convCtx.Convey("Then err should be nil.mids should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(mids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAwardIDsByWinner(t *testing.T) {
	convey.Convey("AwardIDsByWinner", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			am, err := d.AwardIDsByWinner(c, mid)
			convCtx.Convey("Then err should be nil.am should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(am, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDivisionName(t *testing.T) {
	convey.Convey("DivisionName", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			names, err := d.DivisionName(c, mid)
			convCtx.Convey("Then err should be nil.names should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(names, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoJoinedCount(t *testing.T) {
	convey.Convey("JoinedCount", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(2)
			awardID = int64(666)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Exec(c, "INSERT INTO special_award_record(award_id,mid) VALUES(666,2) ON DUPLICATE KEY UPDATE award_id=VALUES(award_id)")
			count, err := d.JoinedCount(c, mid, awardID)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAwardBonus(t *testing.T) {
	convey.Convey("AwardBonus", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			awardID = int64(66)
			prizeID = int64(2)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Exec(c, "INSERT INTO special_award_prize(award_id,prize_id, bonus) VALUES(66,2,100) ON DUPLICATE KEY UPDATE award_id=VALUES(award_id)")
			bonus, err := d.AwardBonus(c, awardID, prizeID)
			convCtx.Convey("Then err should be nil.bonus should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(bonus, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddToAwardRecord(t *testing.T) {
	convey.Convey("AddToAwardRecord", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(2)
			awardID = int64(666)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Exec(c, "DELETE FROM special_award_record WHERE mid=2")
			rows, err := d.AddToAwardRecord(c, mid, awardID)
			convCtx.Convey("Then err should be nil.rows should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetSpecialAwards(t *testing.T) {
	convey.Convey("GetSpecialAwards", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			awards, err := d.GetSpecialAwards(c)
			convCtx.Convey("Then err should be nil.awards should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(awards, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetSpecialAwardDivision(t *testing.T) {
	convey.Convey("GetSpecialAwardDivision", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			awardID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Exec(c, "INSERT INTO special_award_division(award_id,division_name) VALUES(666,'test') ON DUPLICATE KEY UPDATE award_id=VALUES(award_id)")
			divisions, err := d.GetSpecialAwardDivision(c, awardID)
			convCtx.Convey("Then err should be nil.divisions should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(divisions, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAwardWinRecord(t *testing.T) {
	convey.Convey("GetAwardWinRecord", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			awardIDs, err := d.GetAwardWinRecord(c, mid)
			convCtx.Convey("Then err should be nil.awardIDs should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(awardIDs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAwardJoinRecord(t *testing.T) {
	convey.Convey("GetAwardJoinRecord", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			awardIDs, err := d.GetAwardJoinRecord(c, mid)
			convCtx.Convey("Then err should be nil.awardIDs should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(awardIDs, convey.ShouldNotBeNil)
			})
		})
	})
}
