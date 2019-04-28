package activity

import (
	"context"
	"flag"
	"fmt"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/activity"
	"os"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.creative")
		flag.Set("conf_token", "96b6a6c10bb311e894c14a552f48fef8")
		flag.Set("tree_id", "2305")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/creative.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.client.SetTransport(gock.DefaultTransport)
	return r
}
func Test_Likes(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		missionID = int64(10364)
		likeCnt   int
	)
	convey.Convey("Likes", t, func(ctx convey.C) {
		defer gock.OffAll()
		url := fmt.Sprintf(d.ActLikeURI, 10364)
		httpMock("GET", url).Reply(200).JSON(`{"code":0,"data":{"id":53,"mid":27515310,"uname":"1qs314567","state":2,"type":2,"position":4,"url":"http://i0.hdslb.com/bfs/article/578a61e7caf47b5deaa80940b2806f1d9ce53dde.png","md5":"79d02f08c2b7b0d2b30a6b6a4f61c97e","info":"{\"width\":325,\"height\":50}","ctime":1499764110,"mtime":1499828335},"message":"","ttl":1}`)
		likeCnt, err = d.Likes(c, missionID)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(likeCnt, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}
func Test_Protocol(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		missionID = int64(10364)
		p         *activity.Protocol
	)
	convey.Convey("Protocol", t, func(ctx convey.C) {
		defer gock.OffAll()
		url := fmt.Sprintf(d.ActProtocolURI, 10364)
		httpMock("GET", url).Reply(200).JSON(`{"code":0,"data":{"id":"231","sid":"10364","protocol":"说明不超过201234567890123","types":"","tags":"","hot":"0","bgm_id":"0","paster_id":"0","oids":"","screen_set":"1"}}`)
		p, err = d.Protocol(c, missionID)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p, convey.ShouldNotBeNil)
		})
	})
}

func Test_Subject(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		missionID = int64(10364)
		act       *activity.Activity
	)
	convey.Convey("Subject", t, func(ctx convey.C) {
		defer gock.OffAll()
		url := fmt.Sprintf(d.ActSubjectURI, 10364)
		httpMock("GET", url).Reply(200).JSON(`{"code":0,"data":{"id":"231","sid":"10364","Subject":"说明不超过201234567890123","types":"","tags":"","hot":"0","bgm_id":"0","paster_id":"0","oids":"","screen_set":"1"}}`)
		act, err = d.Subject(c, missionID)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(act, convey.ShouldNotBeNil)
		})
	})
}

func Test_Activities(t *testing.T) {
	var (
		c    = context.TODO()
		err  error
		acts []*activity.Activity
	)
	convey.Convey("Activities", t, func(ctx convey.C) {
		defer gock.OffAll()
		url := d.ActAllListURI
		httpMock("GET", url).Reply(200).JSON(`{"code":0,"data":[{"id":10329,"oid":0,"type":4,"state":1,"stime":"2018-10-14 14:32:00","etime":"2018-12-17 00:00:00","ctime":"2018-08-30 18:03:30","mtime":"2018-10-17 18:27:44","name":"这是一个标题，123kkkk有点长有点长，可能展示不下，只能展示13个字","author":"jinchenchen","act_url":"","lstime":"2018-09-08 00:00:00","letime":"2018-12-01 00:00:00","cover":"//uat-i0.hdslb.com/bfs/test/static/20181017/fb8f33d1a41042b9a1ebb515fdc19d94/nBljdwyCo.jpg","dic":"sdf","flag":"33","uetime":"0000-00-00 00:00:00","ustime":"0000-00-00 00:00:00","level":"0","h5_cover":"","rank":"123","like_limit":"1","android_url":"","ios_url":"","fan_limit_max":"0","fan_limit_min":"0","tags":"","hot":0,"bgm_id":123,"paster_id":123,"oids":"10110549|10110164|10110536|10110546","screen_set":1,"protocol":"活动说明，123MMMMM展示展示，只能展示20个，展示不下，显示"}]}`)
		acts, err = d.Activities(c)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(acts, convey.ShouldNotBeNil)
		})
	})
}

func Test_MissionOnlineByTid(t *testing.T) {
	var (
		c    = context.TODO()
		err  error
		tid  = int16(160)
		plat = int16(1)
		mm   []*activity.ActWithTP
	)
	convey.Convey("MissionOnlineByTid", t, func(ctx convey.C) {
		defer gock.OffAll()
		url := d.ActOnlineByTypeURL
		httpMock("GET", url).Reply(200).JSON(`
		{
			"code": 0,
			"data": [
			  {
				"id": 10329,
				"oid": 0,
				"type": 4,
				"state": 1,
				"stime": "2018-10-14 14:32:00",
				"etime": "2018-12-17 00:00:00",
				"ctime": "2018-08-30 18:03:30",
				"mtime": "2018-10-17 18:27:44",
				"name": "这是一个",
				"author": "jinchenchen",
				"act_url": "",
				"lstime": "2018-09-08 00:00:00",
				"letime": "2018-12-01 00:00:00",
				"cover": "//uat-i0.hdslb.com/bfs/test/static/20181017/fb8f33d1a41042b9a1ebb515fdc19d94/nBljdwyCo.jpg",
				"dic": "sdf",
				"flag": "33",
				"uetime": "0000-00-00 00:00:00",
				"ustime": "0000-00-00 00:00:00",
				"level": "0",
				"h5_cover": "",
				"rank": "123",
				"like_limit": "1",
				"android_url": "",
				"ios_url": "",
				"fan_limit_max": "0",
				"fan_limit_min": "0",
				"tags": "",
				"types": "",
				"hot": 0,
				"bgm_id": 123,
				"paster_id": 123,
				"oids": "10110549|10110164|10110536|10110546",
				"screen_set": 1,
				"protocol": "活动说明"
			  }
			]
		  }
		`)
		mm, err = d.MissionOnlineByTid(c, tid, plat)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(mm, convey.ShouldNotBeNil)
		})
	})
}

func Test_Unbind(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		missionID = int64(10364)
		aid       = int64(10110788)
		ip        = "127.0.0.1"
	)
	convey.Convey("Unbind", t, func(ctx convey.C) {
		defer gock.OffAll()
		url := fmt.Sprintf(d.ActUpdateURI, 10364)
		httpMock("POST", url).Reply(200).JSON(`{"code":0,"data":""}`)
		err = d.Unbind(c, aid, missionID, ip)
		ctx.Convey("Then err should be nil.has should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
