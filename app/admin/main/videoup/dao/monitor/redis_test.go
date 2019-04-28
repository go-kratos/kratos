package monitor

import (
	"context"
	"go-common/app/admin/main/videoup/model/monitor"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_GetRules(t *testing.T) {
	Convey("GetRules", t, func() {
		rules, err := d.GetRules(context.TODO(), 1, 1, true)
		So(err, ShouldBeNil)
		So(rules, ShouldNotBeNil)
	})
}

func TestDao_SetRule(t *testing.T) {
	Convey("SetRule", t, func() {
		rule := &monitor.Rule{
			ID:       1,
			Type:     1,
			Business: 1,
			Name:     "一审阶段",
			State:    1,
			RuleConf: &monitor.RuleConf{
				Name: "一审长耗时",
				MoniCdt: map[string]struct {
					Comp  string `json:"comparison"`
					Value int64  `json:"value"`
				}{
					"state": {
						Comp:  "=",
						Value: -1,
					},
				},
				NotifyCdt: map[string]struct {
					Comp  string `json:"comparison"`
					Value int64  `json:"value"`
				}{
					"count": {
						Comp:  ">",
						Value: 10,
					},
					"time": {
						Comp:  ">",
						Value: 10,
					},
				},
				Notify: struct {
					Way    int8     `json:"way"`
					Member []string `json:"member"`
				}{
					Way:    monitor.NotifyTypeEmail,
					Member: []string{"liusiming@bilibili.com"},
				},
			},
		}
		/*rule := &monitor.Rule{
			ID:       6,
			Type:     1,
			Business: 2,
			Name:     "二审阶段",
			State:    1,
			RuleConf: &monitor.RuleConf{
				Name: "二审长耗时",
				MoniCdt: map[string]struct {
					Comp  string `json:"comparison"`
					Value int64  `json:"value"`
				}{
					"state": {
						Comp:  "=",
						Value: -1,
					},
					"round": {
						Comp:  "=",
						Value: 10,
					},
				},
				NotifyCdt: map[string]struct {
					Comp  string `json:"comparison"`
					Value int64  `json:"value"`
				}{
					"count": {
						Comp:  ">",
						Value: 10,
					},
					"time": {
						Comp:  ">",
						Value: 10,
					},
				},
				Notify: struct {
					Way    int8     `json:"way"`
					Member []string `json:"member"`
				}{
					Way:    monitor.NotifyTypeEmail,
					Member: []string{"liusiming@bilibili.com"},
				},
			},
		}*/
		err := d.SetRule(context.TODO(), rule)
		So(err, ShouldBeNil)
	})
}

func TestDao_SetRuleState(t *testing.T) {
	Convey("SetRuleState", t, func() {
		err := d.SetRuleState(context.TODO(), 1, 1, 1, monitor.RuleStateOK)
		So(err, ShouldBeNil)
	})
}

func TestDao_BusKeys(t *testing.T) {
	Convey("BusKeys", t, func() {
		_, keys, err := d.BusStatsKeys(context.TODO(), 1)
		So(err, ShouldBeNil)
		So(keys, ShouldNotBeNil)
	})
}
func TestDao_GetAllRules(t *testing.T) {
	Convey("BusKeys", t, func() {
		_, err := d.GetAllRules(context.Background(), true)
		So(err, ShouldBeNil)
	})
}
