package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/service/main/spy/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoService(t *testing.T) {
	convey.Convey("Service", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			serviceName = "account-service"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			service, err := d.Service(c, serviceName)
			ctx.Convey("Then err should be nil.service should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(service, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoEvent(t *testing.T) {
	convey.Convey("Event", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			eventName = "init_user_info"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Event(c, eventName)
			ctx.Convey("Then err should be nil.event should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddService(t *testing.T) {
	convey.Convey("AddService", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			service = &model.Service{
				Name:     fmt.Sprintf("ut_%d", time.Now().Unix()),
				NickName: "test",
				Status:   0,
				CTime:    xtime.Time(time.Now().Unix()),
				MTime:    xtime.Time(time.Now().Unix()),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.AddService(c, service)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddEvent(t *testing.T) {
	convey.Convey("AddEvent", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			event = &model.Event{
				Name:     fmt.Sprintf("ut_%d", time.Now().Unix()),
				NickName: "test",
				Status:   1,
				CTime:    xtime.Time(time.Now().Unix()),
				MTime:    xtime.Time(time.Now().Unix()),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.AddEvent(c, event)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFactor(t *testing.T) {
	convey.Convey("Factor", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			serviceID = int64(1)
			eventID   = int64(1)
			riskLevel = int8(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			factor, err := d.Factor(c, serviceID, eventID, riskLevel)
			ctx.Convey("Then err should be nil.factor should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(factor, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFactorGroup(t *testing.T) {
	convey.Convey("FactorGroup", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			groupName = "基础资料分值"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.FactorGroup(c, groupName)
			ctx.Convey("Then err should be nil.factorGroup should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaohitInfo(t *testing.T) {
	convey.Convey("hitInfo", t, func(ctx convey.C) {
		var (
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := hitInfo(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitHistory(t *testing.T) {
	convey.Convey("hitHistory", t, func(ctx convey.C) {
		var (
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := hitHistory(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBeginTran(t *testing.T) {
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.BeginTran(c)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserInfo(t *testing.T) {
	convey.Convey("UserInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UserInfo(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateInfo(t *testing.T) {
	convey.Convey("TxUpdateInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			info = &model.UserInfo{
				Mid: 1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.TxUpdateInfo(c, tx, info)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxAddInfo(t *testing.T) {
	convey.Convey("TxAddInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			info = &model.UserInfo{
				Mid: time.Now().Unix(),
			}
			id int64
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			id, err = d.TxAddInfo(c, tx, info)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxAddEventHistory(t *testing.T) {
	convey.Convey("TxAddEventHistory", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ueh = &model.UserEventHistory{
				Mid:    1,
				Remark: "un_test",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.TxAddEventHistory(c, tx, ueh)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxAddPunishment(t *testing.T) {
	convey.Convey("TxAddPunishment", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(1)
			no     = int8(0)
			reason = "unit test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.TxAddPunishment(c, tx, mid, no, reason)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxAddPunishmentQueue(t *testing.T) {
	convey.Convey("TxAddPunishmentQueue", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(1)
			blockNo = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.TxAddPunishmentQueue(c, tx, mid, blockNo)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddPunishmentQueue(t *testing.T) {
	convey.Convey("AddPunishmentQueue", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(1)
			blockNo = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddPunishmentQueue(c, mid, blockNo)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxUpdateEventScore(t *testing.T) {
	convey.Convey("TxUpdateEventScore", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(0)
			escore = int8(0)
			score  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.TxUpdateEventScore(c, tx, mid, escore, score)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxUpdateBaseScore(t *testing.T) {
	convey.Convey("TxUpdateBaseScore", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ui = &model.UserInfo{Mid: 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.TxUpdateBaseScore(c, tx, ui)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxClearReliveTimes(t *testing.T) {
	convey.Convey("TxClearReliveTimes", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ui = &model.UserInfo{Mid: 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.TxClearReliveTimes(c, tx, ui)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoConfigs(t *testing.T) {
	convey.Convey("Configs", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Configs(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoHistoryList(t *testing.T) {
	convey.Convey("HistoryList", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(1)
			size = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.HistoryList(c, mid, size)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateEventScoreReLive(t *testing.T) {
	convey.Convey("TxUpdateEventScoreReLive", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(1)
			escore = int8(0)
			score  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.TxUpdateEventScoreReLive(c, tx, mid, escore, score)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoStatListByMid(t *testing.T) {
	convey.Convey("StatListByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			list, err := d.StatListByMid(c, mid)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStatListByIDAndMid(t *testing.T) {
	convey.Convey("StatListByIDAndMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
			id  = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			list, err := d.StatListByIDAndMid(c, mid, id)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStatListByID(t *testing.T) {
	convey.Convey("StatListByID", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			list, err := d.StatListByID(c, id)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAllEvent(t *testing.T) {
	convey.Convey("AllEvent", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			list, err := d.AllEvent(c)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTelLevel(t *testing.T) {
	convey.Convey("TelLevel", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.TelLevel(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddTelLevelInfo(t *testing.T) {
	convey.Convey("AddTelLevelInfo", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			no = &model.TelRiskLevel{
				Mid:    time.Now().Unix(),
				Level:  1,
				Origin: 1,
				Ctime:  time.Now(),
				Mtime:  time.Now(),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.AddTelLevelInfo(c, no)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}
