package shell

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/growup/conf"
	"go-common/library/net/http/blademaster"

	"github.com/smartystreets/goconvey/convey"
)

var (
	client *Client
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "mobile.studio.growup-admin")
		flag.Set("conf_token", "ac1fd397cbc33eb60541e8734844bdd5")
		flag.Set("tree_id", "13583")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/growup-admin.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	client = New(conf.Conf.ShellConf, blademaster.NewClient(conf.Conf.HTTPClient))
	os.Exit(m.Run())
}

func TestShellSetSign(t *testing.T) {
	convey.Convey("SetSign", t, func(ctx convey.C) {
		var (
			sign = "abc"
			o    = OrderRequest{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			o.SetSign(sign)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestShellSetCustomerID(t *testing.T) {
	convey.Convey("SetCustomerID", t, func(ctx convey.C) {
		var (
			customerID = "111"
			o          = OrderRequest{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			o.SetCustomerID(customerID)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestShellSetSignType(t *testing.T) {
	convey.Convey("SetSignType", t, func(ctx convey.C) {
		var (
			signType = "111"
			o        = OrderRequest{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			o.SetSignType(signType)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestShellIsSuccess(t *testing.T) {
	convey.Convey("IsSuccess", t, func(ctx convey.C) {
		var (
			o = OrderCallbackJSON{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := o.IsSuccess()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestShellIsFail(t *testing.T) {
	convey.Convey("IsFail", t, func(ctx convey.C) {
		var (
			o = OrderCallbackJSON{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := o.IsFail()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestShellIsCreate(t *testing.T) {
	convey.Convey("IsCreate", t, func(ctx convey.C) {
		var (
			o = OrderCallbackJSON{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := o.IsCreate()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestShellNew(t *testing.T) {
	convey.Convey("New", t, func(ctx convey.C) {
		var (
			conf       = &conf.ShellConfig{}
			httpClient = &blademaster.Client{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := New(conf, httpClient)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestShellSetDebug(t *testing.T) {
	convey.Convey("SetDebug", t, func(ctx convey.C) {
		var (
			isDebug = true
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			client.SetDebug(isDebug)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestShellSendOrderRequest(t *testing.T) {
	convey.Convey("SendOrderRequest", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &OrderRequest{
				CustomerID:  "1001",
				ProductName: "test",
				NotifyURL:   "test",
				Rate:        "1",
				SignType:    "test",
				Timestamp:   "test",
				Sign:        "test",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := client.SendOrderRequest(c, req)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestShellSendCheckOrderRequest(t *testing.T) {
	convey.Convey("SendCheckOrderRequest", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &OrderCheckRequest{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := client.SendCheckOrderRequest(c, req)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestShellSendShellRequest(t *testing.T) {
	convey.Convey("SendShellRequest", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			url = "localhost:8080"
			req = interface{}(0)
			res = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := client.SendShellRequest(c, url, req, res)
			ctx.Convey("Then err should be not nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
