package geetest

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"go-common/app/interface/main/answer/conf"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account-law.answer")
		flag.Set("conf_appid", "main.account-law.answer")
		flag.Set("conf_token", "ba3ee255695e8d7b46782268ddc9c8a3")
		flag.Set("tree_id", "25260")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_env", "10")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/answer-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func TestDaoPreProcess(t *testing.T) {
	var (
		c          = context.Background()
		mid        = int64(2205)
		ip         = "127.0.0,1"
		geeType    = "1222"
		gc         = conf.GeetestConfig{CaptchaID: "22"}
		newCaptcha = 1
	)
	convey.Convey("PreProcess", t, func(ctx convey.C) {
		ctx.Convey("req, err = http.NewRequest;err!=nil", func(ctx convey.C) {
			monkey.Patch(http.NewRequest, func(_ string, _ string, _ io.Reader) (_ *http.Request, _ error) {
				return nil, fmt.Errorf("Error")
			})
			_, err := d.PreProcess(c, mid, ip, geeType, gc, newCaptcha)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("req, err = http.NewRequest;err==nil", func(ctx convey.C) {
			ctx.Convey("res, err = d.client.Do; err != nil", func(ctx convey.C) {
				monkey.PatchInstanceMethod(reflect.TypeOf(d.client), "Do", func(_ *http.Client, _ *http.Request) (_ *http.Response, _ error) {
					return nil, fmt.Errorf("Error")
				})
				_, err := d.PreProcess(c, mid, ip, geeType, gc, newCaptcha)
				ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldNotBeNil)
				})
			})
			ctx.Convey("res, err = d.client.Do; err == nil", func(ctx convey.C) {
				params := url.Values{}
				params.Set("user_id", strconv.FormatInt(mid, 10))
				params.Set("new_captcha", strconv.Itoa(newCaptcha))
				params.Set("ip_address", ip)
				params.Set("client_type", geeType)
				params.Set("gt", gc.CaptchaID)
				d.client.Transport = gock.DefaultTransport
				ctx.Convey("res.StatusCode >= http.StatusInternalServerError", func(ctx convey.C) {
					httpMock("GET", d.registerURI+"?"+params.Encode()).Reply(501)
					_, err := d.PreProcess(c, mid, ip, geeType, gc, newCaptcha)
					ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
						ctx.So(err, convey.ShouldNotBeNil)
					})
				})
				ctx.Convey("res.StatusCode < http.StatusInternalServerError", func(ctx convey.C) {
					httpMock("GET", d.registerURI+"?"+params.Encode()).Reply(200).SetHeaders(map[string]string{
						"StatusCode": "200",
					})

					ctx.Convey("bs, err = ioutil.ReadAll(res.Body); err != nil ", func(ctx convey.C) {
						monkey.Patch(ioutil.ReadAll, func(_ io.Reader) (_ []byte, _ error) {
							return nil, fmt.Errorf("Error")
						})
						_, err := d.PreProcess(c, mid, ip, geeType, gc, newCaptcha)
						ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
							ctx.So(err, convey.ShouldNotBeNil)
						})
					})
					ctx.Convey("bs, err = ioutil.ReadAll(res.Body); err == nil ", func(ctx convey.C) {
						ctx.Convey("len(bs) != 32 ", func(ctx convey.C) {
							challenge, err := d.PreProcess(c, mid, ip, geeType, gc, newCaptcha)
							ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
								ctx.So(challenge, convey.ShouldBeEmpty)
								ctx.So(err, convey.ShouldBeNil)
							})
						})
						ctx.Convey("len(bs) == 32 ", func(ctx convey.C) {
							var (
								str = "testeeeeeeeeeeyyyyyyyyyyyyrrrrrr"
								bs  = []byte(str)
							)
							monkey.Patch(ioutil.ReadAll, func(_ io.Reader) (_ []byte, _ error) {
								return bs, nil
							})
							challenge, err := d.PreProcess(c, mid, ip, geeType, gc, newCaptcha)
							ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
								ctx.So(challenge, convey.ShouldNotBeNil)
								ctx.So(err, convey.ShouldBeNil)
							})
						})
					})
				})
			})
		})
		ctx.Reset(func() {
			gock.OffAll()
			d.client.Transport = http.DefaultClient.Transport
			monkey.UnpatchAll()
		})
	})
}

func TestDaoValidate(t *testing.T) {
	var (
		c          = context.Background()
		challenge  = "1"
		seccode    = "127.0.0,1"
		clientType = ""
		ip         = "1222"
		captchaID  = "22"
		mid        = int64(14771787)
	)
	convey.Convey("Validate", t, func(ctx convey.C) {
		ctx.Convey("req, err = http.NewRequest;err!=nil", func(ctx convey.C) {
			monkey.Patch(http.NewRequest, func(_ string, _ string, _ io.Reader) (_ *http.Request, _ error) {
				return nil, fmt.Errorf("Error")
			})
			_, err := d.Validate(c, challenge, seccode, clientType, ip, captchaID, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("req, err = http.NewRequest;err==nil", func(ctx convey.C) {
			ctx.Convey("res, err = d.client.Do; err != nil", func(ctx convey.C) {
				monkey.PatchInstanceMethod(reflect.TypeOf(d.client), "Do", func(_ *http.Client, _ *http.Request) (_ *http.Response, _ error) {
					return nil, fmt.Errorf("Error")
				})
				_, err := d.Validate(c, challenge, seccode, clientType, ip, captchaID, mid)
				ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldNotBeNil)
				})
			})
			ctx.Convey("res, err = d.client.Do; err == nil", func(ctx convey.C) {
				params := url.Values{}
				params.Set("seccode", seccode)
				params.Set("challenge", challenge)
				params.Set("captchaid", captchaID)
				params.Set("client_type", clientType)
				params.Set("ip_address", ip)
				params.Set("json_format", "1")
				params.Set("sdk", "golang_3.0.0")
				params.Set("user_id", strconv.FormatInt(mid, 10))
				params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
				// *bm.Client
				d.clientx.SetTransport(gock.DefaultTransport)
				httpMock("POST", d.validateURI).Reply(200).JSON(`{"code":0}`)
				_, err := d.Validate(c, challenge, seccode, clientType, ip, captchaID, mid)
				ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
				})
			})
		})
		ctx.Reset(func() {
			gock.OffAll()
			d.clientx.SetTransport(http.DefaultClient.Transport)
			monkey.UnpatchAll()
		})
	})
}

func TestDaoGeeConfig(t *testing.T) {
	var gtc = &conf.Geetest{PC: conf.GeetestConfig{CaptchaID: "22", PrivateKEY: "123"}, H5: conf.GeetestConfig{CaptchaID: "22", PrivateKEY: "456"}}
	convey.Convey("GeeConfi", t, func(ctx convey.C) {
		ctx.Convey("t=pc", func(ctx convey.C) {
			var t = "pc"
			gc, geetype := d.GeeConfig(t, gtc)
			ctx.Convey("gc=gtc.PC,geetype =web", func(ctx convey.C) {
				ctx.So(gc, convey.ShouldResemble, gtc.PC)
				ctx.So(geetype, convey.ShouldResemble, "web")
			})
		})
		ctx.Convey("t=h5", func(ctx convey.C) {
			var t = "h5"
			gc, geetype := d.GeeConfig(t, gtc)
			ctx.Convey("gc=gtc.H5,geetype =web", func(ctx convey.C) {
				ctx.So(gc, convey.ShouldResemble, gtc.H5)
				ctx.So(geetype, convey.ShouldResemble, "web")
			})
		})
		ctx.Convey("t=", func(ctx convey.C) {
			var t = ""
			gc, geetype := d.GeeConfig(t, gtc)
			ctx.Convey("gc=gtc.PC,geetype =web", func(ctx convey.C) {
				ctx.So(gc, convey.ShouldResemble, gtc.PC)
				ctx.So(geetype, convey.ShouldResemble, "web")
			})
		})
	})
}
