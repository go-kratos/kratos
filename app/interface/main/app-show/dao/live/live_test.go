package live

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"

	. "github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-show")
		flag.Set("conf_token", "Pae4IDOeht4cHXCdOkay7sKeQwHxKOLA")
		flag.Set("tree_id", "2687")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func TestFeed(t *testing.T) {
	Convey("Feed", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.live).Reply(200).JSON(`{"code":0,"count":1,"lives":[{"owner":{"face":"xxx","mid":1,"name":"xxxx"}}]}`)
		res, err := d.Feed(ctx(), 1, "", "", time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRecommend(t *testing.T) {
	Convey("Recommend", t, func() {
		d.clientAsyn.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.rec).Reply(200).JSON(`{
			"code": 0,
			"data": {
				"count": 1,
				"lives": {
					"subject": [{
						"owner": {
							"face": "xxx",
							"mid": 1,
							"name": "xxxx"
						}
					}],
					"hot": [{
						"owner": {
							"face": "xxx",
							"mid": 1,
							"name": "xxxx"
						}
					}]
				}
			}
		}`)
		res, err := d.Recommend(time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestTopicHots(t *testing.T) {
	Convey("TopicHots", t, func() {
		d.clientAsyn.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.topic).Reply(200).JSON(`{
			"code": 0,
			"data": {
				"list": [{
					"topic_id": 7279615,
					"topic_name": "CP23",
					"picture": "{\"image_src\":\"https:\\/\\/i0.hdslb.com\\/bfs\\/album\\/b46cda4c7e953764c0fcda49c0f06639e7092792.jpg\",\"image_width\":800,\"image_height\":500}"
				}, {
					"topic_id": 9029281,
					"topic_name": "2018COS总结",
					"picture": "{\"image_src\":\"https:\\/\\/i0.hdslb.com\\/bfs\\/album\\/6e703453de103b5b5b9e00640ff9e9d9156950e6.jpg\",\"image_width\":800,\"image_height\":500}"
				}, {
					"topic_id": 2838293,
					"topic_name": "陪你过冬天",
					"picture": "{\"image_src\":\"https:\\/\\/i0.hdslb.com\\/bfs\\/album\\/716e12155e1f2150418c9f215de6bb3c9b38f516.png\",\"image_width\":800,\"image_height\":500}"
				}, {
					"topic_id": 2525230,
					"topic_name": "故事王StoryMan",
					"picture": "{\"image_src\":\"https:\\/\\/i0.hdslb.com\\/bfs\\/album\\/5adb6388723c321d99c5b04fe9839cdaacf6da0d.jpg\",\"image_width\":800,\"image_height\":500}"
				}, {
					"topic_id": 8977836,
					"topic_name": "萌宠暖宝宝",
					"picture": "{\"image_src\":\"https:\\/\\/i0.hdslb.com\\/bfs\\/album\\/adba2d41863e5bf1d269cdf2d4803fac099f8288.jpg\",\"image_width\":800,\"image_height\":500}"
				}, {
					"topic_id": 105286,
					"topic_name": "国家宝藏",
					"picture": "{\"image_src\":\"https:\\/\\/i0.hdslb.com\\/bfs\\/album\\/9650de6e22d0742b63db0feac5bd891faf050b54.png\",\"image_width\":800,\"image_height\":500}"
				}, {
					"topic_id": 8948501,
					"topic_name": "2018绘画总结",
					"picture": "{\"image_src\":\"https:\\/\\/i0.hdslb.com\\/bfs\\/album\\/94c6edd82202ffdfa62c91d0d9d62dcb5b40a0b8.png\",\"image_width\":800,\"image_height\":500}"
				}, {
					"topic_id": 8977910,
					"topic_name": "冬天喝奶茶",
					"picture": "{\"image_src\":\"https:\\/\\/i0.hdslb.com\\/bfs\\/album\\/c0f2533e71732eb700a6ec5a37b633fd5ae707f3.png\",\"image_width\":800,\"image_height\":500}"
				}, {
					"topic_id": 2872407,
					"topic_name": "冬日必备",
					"picture": "{\"image_src\":\"https:\\/\\/i0.hdslb.com\\/bfs\\/album\\/54e1efc022377e209ef98a680f1fb7777e242ca4.png\",\"image_width\":800,\"image_height\":500}"
				}, {
					"topic_id": 2953953,
					"topic_name": "冬日穿搭",
					"picture": "{\"image_src\":\"https:\\/\\/i0.hdslb.com\\/bfs\\/album\\/6522e801f6882c171622a47eb3087d50af42fe7b.jpg\",\"image_width\":800,\"image_height\":500}"
				}]
			}
		}`)
		res, err := d.TopicHots(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestDynamicHot(t *testing.T) {
	Convey("DynamicHot", t, func() {
		d.clientAsyn.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.dynamichot).Reply(200).JSON(`{
			"code": 0,
			"data": {
				"list": [{
					"dynamic_id": 198340924012294846,
					"audit_status": 0,
					"delete_status": 0,
					"mid": 440290,
					"nick_name": "音乐鱼",
					"face_img": "http://i2.hdslb.com/bfs/face/d5e81352871fbe33e7c1aed68ad237c1e0588db3.jpg",
					"rid_type": 2,
					"rid": 10362627,
					"view_count": 5641,
					"comment_count": 10,
					"rcmd_reason": "",
					"dynamic_text": "#CP23##COSPLAY#这两天腿和腰都要断了，光顾着逛没怎么拍 ，但还是满足了😂昨晚上还吃了老上海味道，现在急需一顿火锅！",
					"img_count": 9,
					"imgs": ["https://i0.hdslb.com/bfs/album/615e9efe438a78f25b4ee399fec1e168eaac8c6e.jpg", "https://i0.hdslb.com/bfs/album/1a31a0e93c383183b7b7b02e777248689cac1368.jpg", "https://i0.hdslb.com/bfs/album/270f4b3f84337fa8d2755cff46f7aef6d6f633de.jpg", "https://i0.hdslb.com/bfs/album/25e2fc36946d3f40e242c7619898e2ad5564fd1c.jpg", "https://i0.hdslb.com/bfs/album/2897ee217c455502cde345666ea66086d9f16fc3.jpg", "https://i0.hdslb.com/bfs/album/7c5948d362e17a5118c9cdf095d04c193e50de37.jpg", "https://i0.hdslb.com/bfs/album/473d201e4756fca5b38bc0b535a1086db44d2173.jpg", "https://i0.hdslb.com/bfs/album/91fc858bba009c338d4e1a13ba1953378bd19b15.jpg", "https://i0.hdslb.com/bfs/album/d34307d2e9006847b7e974351ed89c0fe3ad282a.jpg"]
				}]
			}
		}`)
		res, err := d.DynamicHot(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
