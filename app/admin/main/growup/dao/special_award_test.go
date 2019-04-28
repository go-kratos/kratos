package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

var awardID = time.Now().Unix()

func TestDaoAddAward(t *testing.T) {
	convey.Convey("AddAward", t, func(ctx convey.C) {
		var (
			tx, _         = d.BeginTran(context.TODO())
			awardName     = string(time.Now().Unix())
			cycleStart    = time.Now().Format("2006-01-02 15:04:05")
			cycleEnd      = time.Now().Format("2006-01-02 15:04:05")
			announceDate  = time.Now().Format("2006-01-02")
			openTime      = time.Now().Format("2006-01-02 15:04:00")
			displayStatus = int(1)
			totalWinner   = int(0)
			totalBonus    = int(0)
			createdBy     = "test"
		)

		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			p1, err := AddAward(tx, awardID, awardName, cycleStart, cycleEnd, announceDate, openTime, displayStatus, totalWinner, totalBonus, createdBy)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAward(t *testing.T) {
	convey.Convey("Award", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.Award(c, awardID)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSelectAwardForUpdate(t *testing.T) {
	convey.Convey("SelectAwardForUpdate", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			award, err := SelectAwardForUpdate(tx, awardID)
			ctx.Convey("Then err should be nil.award should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(award, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateAward(t *testing.T) {
	convey.Convey("UpdateAward", t, func(ctx convey.C) {
		var (
			tx, _  = d.BeginTran(context.TODO())
			values = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			p1, err := UpdateAward(tx, awardID, values)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxAward(t *testing.T) {
	convey.Convey("Award", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			data, err := Award(tx, awardID)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListAwardsDivision(t *testing.T) {
	convey.Convey("ListAwardsDivision", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
			where = fmt.Sprintf("award_id=%d", awardID)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			res, err := ListAwardsDivision(tx, where)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListDivision(t *testing.T) {
	convey.Convey("ListDivision", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			res, err := ListDivision(tx, awardID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxDivisionInfo(t *testing.T) {
	convey.Convey("DivisionInfo", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			res, err := DivisionInfo(tx, awardID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDivisionInfo(t *testing.T) {
	convey.Convey("DivisionInfo", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.AwardDivisionInfo(c, awardID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListPrize(t *testing.T) {
	convey.Convey("ListPrize", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			res, err := ListPrize(tx, awardID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListResource(t *testing.T) {
	convey.Convey("ListResource", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			res, err := ListResource(tx, awardID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPrizeInfo(t *testing.T) {
	convey.Convey("PrizeInfo", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			res, err := PrizeInfo(tx, awardID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountAward(t *testing.T) {
	convey.Convey("CountAward", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			total, err := CountAward(tx)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountAwardWinner(t *testing.T) {
	convey.Convey("CountAwardWinner", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
			where = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			total, err := CountAwardWinner(tx, awardID, where)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGroupCountAwardWinner(t *testing.T) {
	convey.Convey("GroupCountAwardWinner", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
			where = fmt.Sprintf("award_id=%d", awardID)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			res, err := GroupCountAwardWinner(tx, where)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAwardWinnerAll(t *testing.T) {
	convey.Convey("AwardWinnerAll", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			res, err := AwardWinnerAll(tx, awardID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListAwardRecord(t *testing.T) {
	convey.Convey("ListAwardRecord", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			where = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.ListAwardRecord(c, awardID, where)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryAwardWinner(t *testing.T) {
	convey.Convey("QueryAwardWinner", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
			where = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			res, err := QueryAwardWinner(tx, awardID, where)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListAward(t *testing.T) {
	convey.Convey("ListAward", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.ListAward(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxListAward(t *testing.T) {
	convey.Convey("ListAward", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
			from  = int(1)
			limit = int(20)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			res, err := ListAward(tx, from, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelDivisionAll(t *testing.T) {
	convey.Convey("DelDivisionAll", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			p1, err := DelDivisionAll(tx, awardID)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelWinner(t *testing.T) {
	convey.Convey("DelWinner", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
			where = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := DelWinner(tx, awardID, where)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelWinnerAll(t *testing.T) {
	convey.Convey("DelWinnerAll", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			p1, err := DelWinnerAll(tx, awardID)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelDivisionsExclude(t *testing.T) {
	convey.Convey("DelDivisionsExclude", t, func(ctx convey.C) {
		var (
			tx, _       = d.BeginTran(context.TODO())
			divisionIDs = []int64{1, 2}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			p1, err := DelDivisionsExclude(tx, awardID, divisionIDs)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelPrizeAll(t *testing.T) {
	convey.Convey("DelPrizeAll", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			p1, err := DelPrizeAll(tx, awardID)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelPrizesExclude(t *testing.T) {
	convey.Convey("DelPrizesExclude", t, func(ctx convey.C) {
		var (
			tx, _    = d.BeginTran(context.TODO())
			prizeIDs = []int64{1, 2}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			p1, err := DelPrizesExclude(tx, awardID, prizeIDs)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelResources(t *testing.T) {
	convey.Convey("DelResources", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.TODO())
			where = fmt.Sprintf("award_id=%d", awardID)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			p1, err := DelResources(tx, where)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSaveWinners(t *testing.T) {
	convey.Convey("SaveWinners", t, func(ctx convey.C) {
		var (
			tx, _  = d.BeginTran(context.TODO())
			fields = "award_id,mid,division_id,prize_id,tag_id"
			values = fmt.Sprintf("(%d,1,1,1,1)", awardID)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := SaveWinners(tx, fields, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSaveDivisions(t *testing.T) {
	convey.Convey("SaveDivisions", t, func(ctx convey.C) {
		var (
			tx, _  = d.BeginTran(context.TODO())
			fields = "award_id,division_id,division_name,tag_id"
			values = fmt.Sprintf("(%d,1,'test-division-name',1)", awardID)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := SaveDivisions(tx, fields, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSaveResource(t *testing.T) {
	convey.Convey("SaveResource", t, func(ctx convey.C) {
		var (
			tx, _  = d.BeginTran(context.TODO())
			fields = "award_id,resource_type,resource_index,content"
			values = fmt.Sprintf("(%d,1,1,'test-content')", awardID)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := SaveResource(tx, fields, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSavePrizes(t *testing.T) {
	convey.Convey("SavePrizes", t, func(ctx convey.C) {
		var (
			tx, _  = d.BeginTran(context.TODO())
			fields = "award_id,prize_id,bonus,quota"
			values = fmt.Sprintf("(%d,1,100,100)", awardID)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := SavePrizes(tx, fields, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
