package dao

import (
	"context"
	"net/http"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJtYWluIiwicGxhdGZvcm1faWQiOiIyRjNiOGZEVkdsTW5qOGFDRGxNYVciLCJleHAiOjE1NDA2NDE3MjQsImlzcyI6Im1haW4ifQ.SxKceDl69N1su3kDO70c4P1NuPhXOxp5rXpM3n8Jyig"

func TestToken(t *testing.T) {
	d.client.SetTransport(gock.DefaultTransport)
	convey.Convey("get token", t, func(ctx convey.C) {
		ctx.Convey("When http response code != 0", func(ctx convey.C) {
			httpMock("POST", "http://easyst.bilibili.co/v1/token").Reply(200).JSON(`{
				"code": 90000,
				"data": {
					"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJtYWluIiwicGxhdGZvcm1faWQiOiIyRjNiOGZEVkdsTW5qOGFDRGxNYVciLCJleHAiOjE1NDA2NDE3MjQsImlzcyI6Im1haW4ifQ.SxKceDl69N1su3kDO70c4P1NuPhXOxp5rXpM3n8Jyig",
					"user_name": "main",
					"secret": "",
					"expired": 1540641724
				},
				"message": "success",
				"status": 200
			}`)
			data, err := d.Token(context.Background(), "")
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			gock.Off()
			d.client.SetTransport(http.DefaultClient.Transport)
		})
	})
}

func TestAuth(t *testing.T) {
	d.client.SetTransport(gock.DefaultTransport)
	convey.Convey("get auth", t, func(ctx convey.C) {
		ctx.Convey("When http response code != 0", func(ctx convey.C) {
			httpMock("GET", "http://easyst.bilibili.co/v1/auth").Reply(200).JSON(`{"code":90000,"message":"success"}`)
			_, err := d.Auth(context.Background(), "sven-apm=afab110001a20fa3a1c9b5bd67564ccd71c6253406725184bfa5d88c7ece9d17; Path=/; Domain=bilibili.co; Expires=Tue, 23 Oct 2018 07:25:02 GMT; Max-Age=1800; HttpOnly")
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// ctx.So(data, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			gock.Off()
			d.client.SetTransport(http.DefaultClient.Transport)
		})
	})
}

func TestRole(t *testing.T) {
	d.client.SetTransport(gock.DefaultTransport)
	convey.Convey("TestRole", t, func(ctx convey.C) {
		ctx.Convey("When http response code != 0", func(ctx convey.C) {
			httpMock("GET", "http://easyst.bilibili.co/v1/node/role/app").Reply(200).JSON(`{"code":90000,"message":"success"}`)
			data, err := d.Role(context.Background(), token)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			gock.Off()
			d.client.SetTransport(http.DefaultClient.Transport)
		})
	})
}

func TestNodeTree(t *testing.T) {
	d.client.SetTransport(gock.DefaultTransport)
	convey.Convey("NodeTree", t, func(ctx convey.C) {
		ctx.Convey("When http response code != 0", func(ctx convey.C) {
			httpMock("GET", "http://easyst.bilibili.co/v1/node/bilibili.main.common-arch").Reply(200).JSON(`{
				"status": 200,
				"message": "success",
				"code": 90000,
				"data": {
					"count": 0,
					"results": 20,
					"page": 1,
					"data": null
				}
			}`)
			_, err := d.NodeTree(context.Background(), token, "main.common-arch", "")
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// ctx.So(data, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			gock.Off()
			d.client.SetTransport(http.DefaultClient.Transport)
		})
	})
}

func TestTree(t *testing.T) {
	d.client.SetTransport(gock.DefaultTransport)
	convey.Convey("TestTree", t, func(ctx convey.C) {
		ctx.Convey("When http response code != 0", func(ctx convey.C) {
			httpMock("GET", "http://easyst.bilibili.co/v1/node/apptree").Reply(200).JSON(`{
				"code": 90000,
				"data": {
					"bilibili": {
						"id": 0,
						"name": "bilibili",
						"alias": "哔哩哔哩",
						"type": 1,
						"path": "bilibili",
						"tags": {},
						"children": null
					}
				},
				"message": "success",
				"status": 200
			}`)
			_, err := d.Tree(context.Background(), token)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// ctx.So(data, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			gock.Off()
			d.client.SetTransport(http.DefaultClient.Transport)
		})
	})
}
