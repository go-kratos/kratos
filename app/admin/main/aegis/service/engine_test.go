package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/resource"
)

func TestServiceListBizFlow(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("ListBizFlow", t, func(ctx convey.C) {
		r, err := s.ListBizFlow(c, 1, []int64{1}, []int64{1, 2})
		for _, item := range r {
			t.Logf("bizflow(%+v)", item)
		}
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

//conf = biz=1, mid="1,2",oid="oid1,oid3", oid="oid6"&extra2="3,4"
var grayTest = []struct {
	AddOpt *model.AddOption
	Expect bool
}{
	{AddOpt: &model.AddOption{Resource: resource.Resource{
		BusinessID: 100,
		MID:        1,
		OID:        "oid1",
	}},
		Expect: true,
	},
	{AddOpt: &model.AddOption{Resource: resource.Resource{
		BusinessID: 100,
		MID:        2,
		OID:        "oid2",
	}},
		Expect: true,
	},
	{AddOpt: &model.AddOption{Resource: resource.Resource{
		BusinessID: 100,
		MID:        3,
		OID:        "oid3",
	}},
		Expect: true,
	},
	{AddOpt: &model.AddOption{Resource: resource.Resource{
		BusinessID: 100,
		MID:        4,
		OID:        "oid4",
	}},
		Expect: false,
	},
	{AddOpt: &model.AddOption{Resource: resource.Resource{
		BusinessID: 200,
		MID:        5,
		OID:        "oid5",
	}},
		Expect: true,
	},
	{AddOpt: &model.AddOption{Resource: resource.Resource{
		BusinessID: 100,
		MID:        6,
		OID:        "oid6",
	}},
		Expect: false,
	},
	{AddOpt: &model.AddOption{Resource: resource.Resource{
		BusinessID: 100,
		MID:        7,
		OID:        "oid6",
		Extra2:     5,
	}},
		Expect: false,
	},
	{AddOpt: &model.AddOption{Resource: resource.Resource{
		BusinessID: 100,
		MID:        7,
		OID:        "oid6",
		Extra2:     3,
	}},
		Expect: true,
	},
	{AddOpt: &model.AddOption{Resource: resource.Resource{
		BusinessID: 200, //没配置灰度
		MID:        7,
		OID:        "oid6",
		Extra2:     10,
	}},
		Expect: true,
	},
}

func TestServiceGray(t *testing.T) {
	convey.Convey("Gray", t, func(ctx convey.C) {
		optid := 0
		for biz, options := range s.gray {
			for _, fields := range options {
				optid++
				t.Logf("gray biz(%d) options.id(%d), options.fields.len=%d, field(%+v)\r\n", biz, optid, len(fields), fields)
			}
		}

		ctx.Convey("No return values", func(ctx convey.C) {
			for _, item := range grayTest {
				reality := s.Gray(item.AddOpt)
				t.Logf("gray opt(%+v) reality(%v) expect(%v)", item.AddOpt, reality, item.Expect)
				convey.So(reality, convey.ShouldEqual, item.Expect)
			}
		})
	})
}

func TestService_Auth(t *testing.T) {
	convey.Convey("Auth", t, func(ctx convey.C) {
		auth, err := s.Auth(cntx, 421)
		ctx.Convey("auth ok", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			t.Logf("auth(%+v)", auth)
		})
	})
}
