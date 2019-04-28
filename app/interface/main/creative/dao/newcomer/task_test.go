package newcomer

import (
	"context"
	"testing"

	"go-common/app/interface/main/creative/model/newcomer"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewcomergetTableName(t *testing.T) {
	convey.Convey("getTableName", t, func(ctx convey.C) {
		var (
			mid = int64(27515308)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := getTableName(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerUserTaskBind(t *testing.T) {
	convey.Convey("UserTaskBind", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UserTaskBind(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerIsRewardReceived(t *testing.T) {
	convey.Convey("IsRewardReceived", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(27515308)
			rid        = int64(1)
			rewardType = int8(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.IsRewardReceived(c, mid, rid, rewardType)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerRewardReceivedGroup(t *testing.T) {
	convey.Convey("RewardReceivedGroup", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
			ids = []int64{1, 2, 3}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RewardReceivedGroup(c, mid, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerGiftRewards(t *testing.T) {
	convey.Convey("GiftRewards", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			taskType = int8(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GiftRewards(c, taskType)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerTaskGroupRewards(t *testing.T) {
	convey.Convey("TaskGroupRewards", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			groupID = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.TaskGroupRewards(c, groupID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerRewardCompleteState(t *testing.T) {
	convey.Convey("RewardCompleteState", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(27515308)
			tids = []int64{1, 2, 3}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RewardCompleteState(c, mid, tids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerRewardReceive(t *testing.T) {
	convey.Convey("RewardReceive", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			places = "(?, ?, ?, ?, ? ,?)"
			args   = []interface{}{1, 1, 1, 1, 1, 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.RewardReceive(c, places, args)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerRewardActivate(t *testing.T) {
	convey.Convey("RewardActivate", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
			id  = int64(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.RewardActivate(c, mid, id)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerRewardReceives(t *testing.T) {
	convey.Convey("RewardReceives", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RewardReceives(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerRewards(t *testing.T) {
	convey.Convey("Rewards", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Rewards(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerGiftRewardCount(t *testing.T) {
	convey.Convey("GiftRewardCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
			ids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GiftRewardCount(c, mid, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerBaseRewardCount(t *testing.T) {
	convey.Convey("BaseRewardCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
			ids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.BaseRewardCount(c, mid, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerTasks(t *testing.T) {
	convey.Convey("Tasks", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Tasks(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerUserTasksByMID(t *testing.T) {
	convey.Convey("UserTasksByMID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UserTasksByMID(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerUpUserTask(t *testing.T) {
	convey.Convey("UpUserTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
			tid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.UpUserTask(c, mid, tid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerUserTaskType(t *testing.T) {
	convey.Convey("UserTaskType", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UserTaskType(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerAllTaskGroupRewards(t *testing.T) {
	convey.Convey("AllTaskGroupRewards", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.AllTaskGroupRewards(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerAllGiftRewards(t *testing.T) {
	convey.Convey("AllGiftRewards", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.AllGiftRewards(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerUserTasks(t *testing.T) {
	convey.Convey("UserTasks", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UserTasks(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerRewardReceiveByID(t *testing.T) {
	convey.Convey("RewardReceiveByID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
			rid = int64(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RewardReceiveByID(c, mid, rid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerTaskGroups(t *testing.T) {
	convey.Convey("TaskGroups", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.TaskGroups(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerTaskRewards(t *testing.T) {
	convey.Convey("TaskRewards", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.TaskRewards(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerRewardReceive2(t *testing.T) {
	convey.Convey("RewardReceive2", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rrs = make([]*newcomer.RewardReceive2, 0)
			mid = int64(27515308)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RewardReceive2(c, mid, rrs)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerRewardReceiveByOldInfo(t *testing.T) {
	convey.Convey("RewardReceiveByOldInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
			r   = &newcomer.RewardReceive{
				ID:          2,
				MID:         mid,
				TaskGiftID:  0,
				TaskGroupID: 2,
				RewardID:    7,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RewardReceiveByOldInfo(c, r)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewcomerRewardActivate2(t *testing.T) {
	convey.Convey("RewardActivate2", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
			oid = int64(2)
			nid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RewardActivate2(c, mid, oid, nid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
