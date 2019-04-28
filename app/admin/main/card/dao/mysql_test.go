package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/admin/main/card/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddGroup(t *testing.T) {
	convey.Convey("AddGroup", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.AddGroup{
				Name:     fmt.Sprintf("%v", time.Now().Unix()),
				State:    0,
				Operator: "superman",
				OrderNum: time.Now().Unix(),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddGroup(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateGroup(t *testing.T) {
	convey.Convey("UpdateGroup", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.UpdateGroup{
				ID:       1,
				Name:     fmt.Sprintf("%v", time.Now().UnixNano()),
				State:    0,
				Operator: "superman",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateGroup(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCards(t *testing.T) {
	convey.Convey("Cards", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Cards(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGroupByName(t *testing.T) {
	convey.Convey("GroupByName", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GroupByName("te123156544664st")
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCardByName(t *testing.T) {
	convey.Convey("CardByName", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.CardByName("tes...t")
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCardsByGid(t *testing.T) {
	convey.Convey("CardsByGid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			gid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CardsByGid(c, gid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCardsByIds(t *testing.T) {
	convey.Convey("CardsByIds", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1, 2, 3, 4, 5}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CardsByIds(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGroupsByIds(t *testing.T) {
	convey.Convey("GroupsByIds", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1, 2, 3, 4, 5}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GroupsByIds(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCard(t *testing.T) {
	convey.Convey("AddCard", t, func(ctx convey.C) {
		var (
			arg = &model.AddCard{
				Name:       fmt.Sprintf("%v", time.Now().Unix()),
				State:      0,
				Operator:   "superman",
				OrderNum:   time.Now().Unix(),
				CardURL:    "http://www.baidu.com",
				BigCradURL: "http://www.baidu.com",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCard(arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateCard(t *testing.T) {
	convey.Convey("UpdateCard", t, func(ctx convey.C) {
		var (
			req = &model.UpdateCard{
				ID:       1,
				Name:     fmt.Sprintf("%v", time.Now().UnixNano()),
				State:    0,
				Operator: "superman",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateCard(req)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateCardState(t *testing.T) {
	convey.Convey("UpdateCardState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			state = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateCardState(c, id, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDeleteCard(t *testing.T) {
	convey.Convey("DeleteCard", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DeleteCard(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDeleteGroup(t *testing.T) {
	convey.Convey("DeleteGroup", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DeleteGroup(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateGroupState(t *testing.T) {
	convey.Convey("UpdateGroupState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			state = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateGroupState(c, id, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoMaxCardOrder(t *testing.T) {
	convey.Convey("MaxCardOrder", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			max, err := d.MaxCardOrder()
			ctx.Convey("Then err should be nil.max should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(max, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMaxGroupOrder(t *testing.T) {
	convey.Convey("MaxGroupOrder", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			max, err := d.MaxGroupOrder()
			ctx.Convey("Then err should be nil.max should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(max, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBatchUpdateCardOrder(t *testing.T) {
	convey.Convey("BatchUpdateCardOrder", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			cs = []*model.Card{{
				ID: 1, OrderNum: 1},
				{ID: 2, OrderNum: 2}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.BatchUpdateCardOrder(c, cs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchUpdateCardGroupOrder(t *testing.T) {
	convey.Convey("BatchUpdateCardGroupOrder", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			cs = []*model.CardGroup{{ID: 2, OrderNum: 2}, {ID: 1, OrderNum: 1}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.BatchUpdateCardGroupOrder(c, cs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGroups(t *testing.T) {
	convey.Convey("Groups", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgQueryGroup{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Groups(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
