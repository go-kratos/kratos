package recommend

import (
	"context"
	"flag"
	"fmt"
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

func TestHots(t *testing.T) {
	Convey("Hots", t, func() {
		d.clientAsyn.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.hotUrl).Reply(200).JSON(`{
			"note": false,
			"source_date": "2019-01-07",
			"code": 0,
			"num": 500,
			"list": [{
				"aid": 39185037,
				"score": 176
			}, {
				"aid": 39658458,
				"score": 174
			}, {
				"aid": 39532823,
				"score": 168
			}, {
				"aid": 39477161,
				"score": 168
			}, {
				"aid": 39852951,
				"score": 168
			}, {
				"aid": 39672060,
				"score": 168
			}, {
				"aid": 39832577,
				"score": 168
			}, {
				"aid": 39987017,
				"score": 168
			}, {
				"aid": 39700424,
				"score": 163
			}]
		}`)
		res, err := d.Hots(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRegion(t *testing.T) {
	Convey("Region", t, func() {
		d.clientAsyn.SetTransport(gock.DefaultTransport)
		api := fmt.Sprintf(d.regionUrl, "33")
		httpMock("GET", api).Reply(200).JSON(`{
			"code": 0,
			"list": [{
				"aid": "39911001",
				"score": 523
			}, {
				"aid": "39852951",
				"score": 6732
			}, {
				"aid": "39845334",
				"score": 31
			}]
		}`)
		res, err := d.Region(ctx(), "33")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRegionHots(t *testing.T) {
	Convey("RegionHots", t, func() {
		d.clientAsyn.SetTransport(gock.DefaultTransport)
		api := fmt.Sprintf(d.rankRegionAppUrl, 1)
		httpMock("GET", api).Reply(200).JSON(`{
			"note": "统计3日内新投稿的数据综合得分，每二十分钟更新一次。",
			"source_date": "2019-01-07",
			"code": 0,
			"num": 100,
			"list": [{
				"aid": 39894949,
				"mid": 808171,
				"score": 546760
			}, {
				"aid": 39877679,
				"mid": 7487399,
				"score": 516724
			}]
		}`)
		res, err := d.RegionHots(ctx(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRegionList(t *testing.T) {
	Convey("RegionList", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.regionListUrl).Reply(200).JSON(`{
			"code": 0,
			"list": [{
				"aid": 39903065
			}]
		}`)
		res, err := d.RegionList(ctx(), 1, 1, 1, 1, 1, "")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRegionChildHots(t *testing.T) {
	Convey("RegionChildHots", t, func() {
		d.clientAsyn.SetTransport(gock.DefaultTransport)
		api := fmt.Sprintf(d.regionChildHotUrl, 1)
		httpMock("GET", api).Reply(200).JSON(`{
			"code": 0,
			"list": [{
				"aid": 39903065
			}]
		}`)
		res, err := d.RegionChildHots(ctx(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRegionArcList(t *testing.T) {
	Convey("RegionArcList", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.regionArcListUrl).Reply(200).JSON(`{
			"code": 0,
			"data": {
				"archives": [{
					"aid": 39903065
				}]
			}
		}`)
		res, err := d.RegionArcList(ctx(), 1, 1, 1, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRankRegion(t *testing.T) {
	Convey("RankRegion", t, func() {
		d.clientAsyn.SetTransport(gock.DefaultTransport)
		httpMock("GET", fmt.Sprintf(d.rankRegionUrl, "all", 1)).Reply(200).JSON(`{
			"rank": {
				"code": 0,
				"list": [{
					"aid": 39903065
				}]
			}
		}`)
		res, err := d.RankRegion(ctx(), 1, "all")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRankAll(t *testing.T) {
	Convey("RankAll", t, func() {
		d.clientAsyn.SetTransport(gock.DefaultTransport)
		httpMock("GET", fmt.Sprintf(d.rankOriginalUrl, "all")).Reply(200).JSON(`{
			"rank": {
				"code": 0,
				"list": [{
					"aid": 39903065
				}]
			}
		}`)
		res, err := d.RankAll(ctx(), "all")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRankBangumi(t *testing.T) {
	Convey("RankBangumi", t, func() {
		d.clientAsyn.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.rankBangumilUrl).Reply(200).JSON(`{
			"rank": {
				"code": 0,
				"list": [{
					"aid": 39903065
				}]
			}
		}`)
		res, err := d.RankBangumi(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestFeedDynamic(t *testing.T) {
	Convey("FeedDynamic", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.feedDynamicUrl).Reply(200).JSON(`{
			"code": 0,
			"data": [12587337, 1840325, 38132621, 5910308, 26879875, 26308630, 7348036, 1766719, 6374879, 24937721],
			"hot": null,
			"ctop": 12587337,
			"cbottom": 24937721
		}`)
		_, newAids, _, _, err := d.FeedDynamic(ctx(), false, 1, 1, 1, 1, time.Now())
		So(newAids, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRankAppRegion(t *testing.T) {
	Convey("RankAppRegion", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		api := fmt.Sprintf(d.rankRegionAppUrl, 1)
		httpMock("GET", api).Reply(200).JSON(`{
			"note": "统计3日内新投稿的数据综合得分，每二十分钟更新一次。",
			"source_date": "2019-01-07",
			"code": 0,
			"num": 100,
			"list": [{
				"aid": 39954334,
				"mid": 38125899,
				"score": 509800,
				"others": [{
					"aid": 39903065,
					"score": 48222
				}]
			}, {
				"aid": 39953503,
				"mid": 3969839,
				"score": 430381
			}]
		}`)
		res, _, _, err := d.RankAppRegion(ctx(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRankAppOrigin(t *testing.T) {
	Convey("RankAppOrigin", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.rankOriginAppUrl).Reply(200).JSON(`{
			"note": "统计3日内新投稿的数据综合得分，每二十分钟更新一次。",
			"source_date": "2019-01-07",
			"code": 0,
			"num": 100,
			"list": [{
				"aid": 39954334,
				"mid": 38125899,
				"score": 509800,
				"others": [{
					"aid": 39903065,
					"score": 48222
				}]
			}, {
				"aid": 39953503,
				"mid": 3969839,
				"score": 430381
			}]
		}`)
		res, _, _, err := d.RankAppOrigin(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRankAppAll(t *testing.T) {
	Convey("RankAppAll", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.rankAllAppUrl).Reply(200).JSON(`{
			"note": "统计3日内新投稿的数据综合得分，每二十分钟更新一次。",
			"source_date": "2019-01-07",
			"code": 0,
			"num": 100,
			"list": [{
				"aid": 39954334,
				"mid": 38125899,
				"score": 509800,
				"others": [{
					"aid": 39903065,
					"score": 48222
				}]
			}, {
				"aid": 39953503,
				"mid": 3969839,
				"score": 430381
			}]
		}`)
		res, _, _, err := d.RankAppAll(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRankAppBangumi(t *testing.T) {
	Convey("RankAppBangumi", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.rankBangumiAppUrl).Reply(200).JSON(`{
			"note": "统计3日内新投稿的数据综合得分，每二十分钟更新一次。",
			"source_date": "2019-01-07",
			"code": 0,
			"num": 100,
			"list": [{
				"aid": 39954334,
				"mid": 38125899,
				"score": 509800,
				"others": [{
					"aid": 39903065,
					"score": 48222
				}]
			}, {
				"aid": 39953503,
				"mid": 3969839,
				"score": 430381
			}]
		}`)
		res, _, _, err := d.RankAppBangumi(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestHotTab(t *testing.T) {
	Convey("HotTab", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		httpMock("GET", d.hottabURL).Reply(200).JSON(`{
			"note": false,
			"source_date": "2019-01-07",
			"code": 0,
			"num": 100,
			"list": [{
				"aid": 40063426,
				"mid": 837470,
				"score": 764906,
				"desc": "很多人分享",
				"corner_mark": 0
			}, {
				"aid": 39425207,
				"mid": 4870926,
				"score": 690583,
				"desc": "百万播放",
				"corner_mark": 0
			}]
		}`)
		res, err := d.HotTab(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestHotTenTab(t *testing.T) {
	Convey("HotTenTab", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		httpMock("GET", fmt.Sprintf(d.hotHetongURL, 1)).Reply(200).JSON(`{
			"note": false,
			"source_date": "2019-01-07",
			"code": 0,
			"num": 100,
			"list": [{
				"aid": 40063426,
				"mid": 837470,
				"score": 764906,
				"desc": "很多人分享",
				"corner_mark": 0
			}, {
				"aid": 39425207,
				"mid": 4870926,
				"score": 690583,
				"desc": "百万播放",
				"corner_mark": 0
			}]
		}`)
		res, err := d.HotTenTab(ctx(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestHotHeTongTabCard(t *testing.T) {
	Convey("HotHeTongTabCard", t, func() {
		d.client.SetTransport(gock.DefaultTransport)
		httpMock("GET", fmt.Sprintf(d.hotHeTongtabcardURL, 1)).Reply(200).JSON(`{
			"code": 0,
			"list": [{
				"goto": "av",
				"id": 40063426,
				"from_type": "recommend",
				"desc": "8千分享",
				"corner_mark": 0
			}, {
				"goto": "av",
				"id": 39425207,
				"from_type": "recommend",
				"desc": "百万播放",
				"corner_mark": 0
			}, {
				"goto": "av",
				"id": 39920213,
				"from_type": "recommend",
				"desc": "百万播放",
				"corner_mark": 0
			}, {
				"goto": "av",
				"id": 39237975,
				"from_type": "recommend",
				"desc": "百万播放",
				"corner_mark": 0
			}]
		}`)
		res, err := d.HotHeTongTabCard(ctx(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
