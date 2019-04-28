package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/app/admin/main/aegis/model/resource"
	"go-common/app/admin/main/aegis/model/task"
)

func Test_AuditLog(t *testing.T) {
	opt := &model.SubmitOptions{
		Binds: []int64{1, 2},
		EngineOption: model.EngineOption{
			BaseOptions: common.BaseOptions{
				RID: 1,
			},
			Result: &resource.Result{},
		},
	}
	res := &net.TriggerResult{
		SubmitToken: &net.TokenPackage{
			Values: map[string]interface{}{
				"state":  1,
				"forbid": 2,
			},
			TokenIDList: []int64{3, 4},
		},
		ResultToken: &net.TokenPackage{
			Values: map[string]interface{}{
				"state":  1,
				"forbid": 2,
			},
			TokenIDList: []int64{3, 4},
		},
	}

	err := s.sendAuditLog(context.TODO(), "submit", opt, res, model.LogTypeAuditSubmit)
	fmt.Printf("err(%v)\n", err)
	t.Fail()
}

func Test_TaskLog(t *testing.T) {
	Content := map[string]interface{}{
		"task": &task.Task{},
	}

	bs, _ := json.Marshal(Content)
	fmt.Printf("bs(%s)", string(bs))
	t.Fail()
}

func Test_ResourceLog(t *testing.T) {
	Content := map[string]interface{}{
		"opt": &model.AddOption{},
		"res": &net.TriggerResult{
			SubmitToken: &net.TokenPackage{
				Values: map[string]interface{}{
					"state":  1,
					"forbid": 2,
				},
				TokenIDList: []int64{3, 4},
			},
			ResultToken: &net.TokenPackage{
				Values: map[string]interface{}{
					"state":  1,
					"forbid": 2,
				},
				TokenIDList: []int64{3, 4},
			},
		},
		"oids": []int64{1, 2, 3},
		"err":  "dqfgug",
	}

	bs, _ := json.Marshal(Content)
	fmt.Printf("bs(%s)", string(bs))
	t.Fail()
}

func TestService_SearchAuditLog1(t *testing.T) {
	convey.Convey("SearchAuditLog", t, func(ctx convey.C) {
		pm := &model.SearchAuditLogParam{
			OID:        []string{"196962673299031503"},
			BusinessID: 1,
			Ps:         20,
			Pn:         1,
			Username:   []string{"业务方"},
			CtimeFrom:  "2018-11-01 00:00:00",
			CtimeTo:    "",
			TaskID:     []int64{0, 660},
			State:      "0",
		}
		data, pager, err := s.SearchAuditLog(cntx, pm)
		t.Logf("pager(%+v)", pager)
		for i, item := range data {
			t.Logf("data  i=%d, %+v", i, item)
		}
		ctx.So(err, convey.ShouldBeNil)
	})
}

func TestService_SearchAuditLogCSV(t *testing.T) {
	convey.Convey("SearchAuditLogCSV", t, func(ctx convey.C) {
		pm := &model.SearchAuditLogParam{
			OID:        []string{"196962673299031503"},
			BusinessID: 1,
			Ps:         20,
			Pn:         1,
			//Username:   "chenxuefeng",
			CtimeFrom: "2018-11-01 00:00:00",
			CtimeTo:   "",
			//TaskID:    []int64{0, 55},
			//State: "1",
		}
		data, err := s.SearchAuditLogCSV(cntx, pm)
		for i, item := range data {
			t.Logf("data  i=%d, %+v", i, item)
		}
		ctx.So(err, convey.ShouldBeNil)
	})
}

func TestService_TrackResource(t *testing.T) {
	convey.Convey("TrackResource", t, func(ctx convey.C) {
		pm := &model.TrackParam{
			OID:          "186464454672655532112",
			BusinessID:   1,
			Pn:           2,
			Ps:           2,
			LastPageTime: "2018-12-06 14:01:15",
		}

		data, p, err := s.TrackResource(cntx, pm)
		ctx.So(err, convey.ShouldBeNil)
		t.Logf("data(%+v) pager(%+v), params(%+v)", data, p, pm)
		for i, item := range data.Add {
			t.Logf("data.add  i=%d, %+v", i, item)
		}
		for i, item := range data.Audit {
			t.Logf("data.audit  i=%d, %+v", i, item)
		}
	})
}

func TestService_searchConsumerLog(t *testing.T) {
	convey.Convey("searchConsumerLog", t, func(ctx convey.C) {
		res, err := s.searchConsumerLog(cntx, 1, 1, []string{"on"}, []int64{1148}, 10)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(res, convey.ShouldNotBeNil)
		t.Logf("res(%+v)", res)
	})
}
