package model

import (
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/library/log"
)

func TestElecPrepUPRankShift(t *testing.T) {
	convey.Convey("shift\n", t, func() {
		rank := &RankElecPrepUPProto{
			Count: 2,
			UPMID: 1,
			Size_: 100,
		}
		rank.List = []*RankElecPrepElementProto{
			&RankElecPrepElementProto{
				MID:    46333,
				Rank:   1,
				Amount: 100,
			}, &RankElecPrepElementProto{
				MID:    35858,
				Rank:   2,
				Amount: 99,
			}, &RankElecPrepElementProto{
				MID:    233,
				Rank:   3,
				Amount: 98,
			}, &RankElecPrepElementProto{
				MID:    2,
				Rank:   4,
				Amount: 1,
			},
		}
		rank.shift(0, 2)
		log.Info("%s", rank)

		convey.So(rank.List, convey.ShouldHaveLength, 4)
		convey.So(rank.List[0], convey.ShouldBeNil)
		convey.So(rank.List[1].Amount, convey.ShouldEqual, 100)
		convey.So(rank.List[2].Amount, convey.ShouldEqual, 99)
	})
}

func TestElecPrepUPRankUpdate(t *testing.T) {
	convey.Convey("shuffle\n", t, func() {
		rank := &RankElecPrepUPProto{
			Count: 2,
			UPMID: 1,
			Size_: 100,
		}
		rank.List = []*RankElecPrepElementProto{
			&RankElecPrepElementProto{
				MID:    46333,
				Rank:   1,
				Amount: 100,
			}, &RankElecPrepElementProto{
				MID:    35858,
				Rank:   2,
				Amount: 99,
			}, &RankElecPrepElementProto{
				MID:    233,
				Rank:   3,
				Amount: 98,
			}, &RankElecPrepElementProto{
				MID:    2,
				Rank:   4,
				Amount: 1,
			},
		}
		rank.update(&RankElecPrepElementProto{
			MID:    2,
			Rank:   4,
			Amount: 2,
		})
		log.Info("%s", rank)
	})
}

func TestElecPrepUPRankInsert(t *testing.T) {
	convey.Convey("shuffle\n", t, func() {
		rank := &RankElecPrepUPProto{
			Count: 2,
			UPMID: 1,
			Size_: 100,
		}
		rank.List = []*RankElecPrepElementProto{
			&RankElecPrepElementProto{
				MID:    46333,
				Rank:   1,
				Amount: 100,
			}, &RankElecPrepElementProto{
				MID:    35858,
				Rank:   2,
				Amount: 99,
			}, &RankElecPrepElementProto{
				MID:    233,
				Rank:   3,
				Amount: 98,
			}, &RankElecPrepElementProto{
				MID:    2,
				Rank:   4,
				Amount: 1,
			},
		}
		rank.insert(&RankElecPrepElementProto{
			MID:    322,
			Rank:   -1,
			Amount: 101,
		})
		log.Info("%s", rank)
	})
}

func TestElecPrepUPRankCharge(t *testing.T) {
	convey.Convey("charge\n", t, func() {
		rank := &RankElecPrepUPProto{
			Count: 0,
			UPMID: 2,
			Size_: 3,
		}
		rank.Charge(35858, 100, true)
		convey.So(len(rank.List), convey.ShouldEqual, 1)
		convey.So(rank.List[0], convey.ShouldResemble, &RankElecPrepElementProto{
			MID:    35858,
			Rank:   1,
			Amount: 100,
		})

		rank.Charge(46333, 100, true)
		convey.So(len(rank.List), convey.ShouldEqual, 2)
		convey.So(rank.List[1], convey.ShouldResemble, &RankElecPrepElementProto{
			MID:    46333,
			Rank:   2,
			Amount: 100,
		})

		rank.Charge(233, 100, true)
		convey.So(len(rank.List), convey.ShouldEqual, 3)
		convey.So(rank.List[2], convey.ShouldResemble, &RankElecPrepElementProto{
			MID:    233,
			Rank:   3,
			Amount: 100,
		})
		log.Info("%s", rank)

		rank.Charge(46333, 101, false)
		log.Info("%s", rank)
		convey.So(len(rank.List), convey.ShouldEqual, 3)
		convey.So(rank.List[0], convey.ShouldResemble, &RankElecPrepElementProto{
			MID:    46333,
			Rank:   1,
			Amount: 101,
		})
		convey.So(rank.List[1], convey.ShouldResemble, &RankElecPrepElementProto{
			MID:    35858,
			Rank:   2,
			Amount: 100,
		})

		rank.Charge(233, 200, false)
		convey.So(len(rank.List), convey.ShouldEqual, 3)
		convey.So(rank.List[0], convey.ShouldResemble, &RankElecPrepElementProto{
			MID:    233,
			Rank:   1,
			Amount: 200,
		})
		convey.So(rank.List[1], convey.ShouldResemble, &RankElecPrepElementProto{
			MID:    46333,
			Rank:   2,
			Amount: 101,
		})

		log.Info("%s", rank)
	})
}

func TestElecPrepUPRankUpdateMessage(t *testing.T) {
	convey.Convey("charge\n", t, func() {
		rank := &RankElecPrepUPProto{
			Count: 0,
			UPMID: 2,
			Size_: 3,
		}
		rank.Charge(35858, 100, true)
		rank.Charge(46333, 100, true)
		rank.Charge(233, 100, true)

		rank.UpdateMessage(35858, "ut", false)
		rank.UpdateMessage(46333, "ut", true)
		rank.UpdateMessage(233, "ut-hello", false)

		convey.So(rank.Find(35858), convey.ShouldNotBeNil)
		convey.So(rank.Find(35858).Message, convey.ShouldNotBeNil)
		convey.So(rank.Find(35858).Message.Message, convey.ShouldEqual, "ut")
		convey.So(rank.Find(35858).Message.Hidden, convey.ShouldEqual, false)
	})
}

func TestElecPrepUPRankChargeRandom(t *testing.T) {
	rank := &RankElecPrepUPProto{
		Count: 0,
		UPMID: 1,
		Size_: 100,
	}
	for i := 0; i < 100; i++ {
		rank.Charge(randomEle())
	}
	log.Info("%s", rank)
}

func BenchmarkElecPrepUPRankCharge(b *testing.B) {
	rank := &RankElecPrepUPProto{
		Count: 0,
		UPMID: 1,
		Size_: 100,
	}
	for i := 0; i < b.N; i++ {
		rank.Charge(randomEle())
	}
	log.Info("%s", rank)
}

var (
	mids = []int64{2, 46333, 35858, 233}
	i    = 0
)

func randomEle() (mid int64, amount int64, isNew bool) {
	mid = mids[rand.Intn(len(mids))]
	i++
	amount = int64(i)
	isNew = true
	return
}

func TestElecUserSetting(t *testing.T) {
	convey.Convey("TestElecUserSetting", t, func() {
		set := ElecUserSetting(2147483646)
		convey.So(set.ShowMessage(), convey.ShouldBeFalse)
		set = ElecUserSetting(0x7ffffff)
		convey.So(set.ShowMessage(), convey.ShouldBeTrue)
		set = ElecUserSetting(0x1)
		convey.So(set.ShowMessage(), convey.ShouldBeTrue)
	})
}
