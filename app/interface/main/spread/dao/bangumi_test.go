package dao

import (
	"context"
	"net/http"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBangumiContent(t *testing.T) {
	httpMock("GET", "http://uat-bangumi.bilibili.co/ext/internal/archive/channel/content").Reply(http.StatusOK).JSON(`
{
  "code": 0,
  "message": "success",
  "result": [
    {
      "akira": "11的期望",
      "alias": "qwdkhj qwknd,qwjdbhqwdjkbqw,qkjcxsa,jcb,sacb,askjcbscajbsakhc",
      "copyright": "bilibili",
      "cover_image": "http://i0.hdslb.com/bfs/bangumi/4f84b91e5b90e99d8b96a336385af7d84c308b48.jpg",
      "display_address": "https://www.bilibili.com/bangumi/play/ss20017?bsource=baidu_os",
      "download_address": "http://app.bilibili.com?bsource=baidu_os",
      "duration": 21,
      "episodes": [
        {
          "cover": "http://i0.hdslb.com/bfs/archive/496ea8899680d4a80d163d2edb401b23.jpg",
          "duration": 0,
          "id": 116664,
          "index": 1,
          "play_url": "https://www.bilibili.com/bangumi/play/ep116664?bsource=baidu_os",
          "pub_real_time": "2018-08-07 00:00:00",
          "title": "第二集"
        },
        {
          "cover": "http://i0.hdslb.com/bfs/archive/496ea8899680d4a80d163d2edb401b23.jpg",
          "duration": 0,
          "id": 116865,
          "index": 2,
          "play_url": "https://www.bilibili.com/bangumi/play/ep116865?bsource=baidu_os",
          "pub_real_time": "2018-09-10 04:00:00",
          "title": "9.10zuixin"
        },
        {
          "cover": "http://i0.hdslb.com/bfs/archive/1fda382339317a7f6c918827b261965c24cac831.jpg",
          "duration": 0,
          "id": 117307,
          "index": 3,
          "play_url": "https://www.bilibili.com/bangumi/play/ep117307?bsource=baidu_os",
          "pub_real_time": "2018-11-07 11:23:00",
          "title": "不可播，就不玩了，找邱穗姬"
        }
      ],
      "intro": "kate_sponsor_谁都不能动dqw qwd ",
      "is_finish": 0,
      "media_id": 2130686907,
      "name": "免费时承包,转付费后随便看的番",
      "play_count": 0,
      "premieredate": "2018",
      "pub_real_time": 1541560980,
      "pub_time": "2018-04-02 00:00:00",
      "season": {
        "id": 20017,
        "index": 1,
        "pay_price": 0.0,
        "paymentstatus": 1,
        "title": "第一季",
        "total_count": 6
      },
      "seasonId": 20017,
      "season_series": [
        {
          "id": 20017,
          "index": 1,
          "title": "免费时承包,转付费后随便看的番"
        },
        {
          "id": 33409,
          "index": 2,
          "title": "介绍姜姜的小店的故事"
        }
      ],
      "staff": {},
      "tag": [],
      "type": 1
    }
  ],
  "total": 13
}
`)
	convey.Convey("BangumiContent", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			pn     = int(1)
			ps     = int(10)
			typ    = int8(1)
			appkey = "douban"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			resp, err := d.BangumiContent(c, pn, ps, typ, appkey)
			ctx.Convey("Then err should be nil.resp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(resp, convey.ShouldNotBeEmpty)
			})
		})
	})
}

func TestDaoBangumiOff(t *testing.T) {
	resp := `
{
  "code": 0,
  "message": "success",
  "ttl": 1,
  "data": [
    {
      "name": "中二病也要谈恋爱！恋",
      "seasonid": 4349,
      "type": 1
    },
    {
      "name": "天空與海洋之間（僅限港澳台地區）",
      "seasonid": 25687,
      "type": 1
    },
    {
      "name": "只要別西卜大小姐喜歡就好（僅限港澳台地區）",
      "seasonid": 25836,
      "type": 1
    },
    {
      "name": "嫁给非人类",
      "seasonid": 25711,
      "type": 1
    },
    {
      "name": "Tokyo Guru: re (Part 2)（僅限港澳台地區）",
      "seasonid": 25727,
      "type": 1
    },
    {
      "name": "產子救世錄（僅限港澳台地區）",
      "seasonid": 25959,
      "type": 1
    },
    {
      "name": "精灵宝可梦 日月",
      "seasonid": 5707,
      "type": 1
    },
    {
      "name": "剧场版「吸血鬼仆人 - Alice in the Garden -」",
      "seasonid": 25951,
      "type": 1
    },
    {
      "name": "",
      "seasonid": 25958,
      "type": 1
    },
    {
      "name": "新战神金刚：传奇的保护神 第七季",
      "seasonid": 25411,
      "type": 1
    },
    {
      "name": "告诉我魔法钟摆～莉露莉露妖精莉露～",
      "seasonid": 24579,
      "type": 1
    },
    {
      "name": "草莓棉花糖 OVA 第1期",
      "seasonid": 4828,
      "type": 1
    },
    {
      "name": "NEKOPARA EXTRA 小猫篇（猫娘乐园）",
      "seasonid": 25152,
      "type": 1
    },
    {
      "name": "新战神金刚：传奇的保护神 第六季",
      "seasonid": 25013,
      "type": 1
    },
    {
      "name": "闪电十一人 第一季 日语",
      "seasonid": 24833,
      "type": 1
    },
    {
      "name": "致命紫罗兰编号044",
      "seasonid": 24779,
      "type": 1
    },
    {
      "name": "灰与幻想的格林姆迦尔 OVA",
      "seasonid": 24745,
      "type": 1
    },
    {
      "name": "tsetfj",
      "seasonid": 24660,
      "type": 1
    },
    {
      "name": "未来卡 神搭档对战",
      "seasonid": 24416,
      "type": 1
    },
    {
      "name": "明日之丈",
      "seasonid": 24332,
      "type": 1
    }
  ]
}
`
	httpMock("GET", "http://uat-bangumi.bilibili.co/ext/internal/archive/channel/content/offshelve").Reply(http.StatusOK).JSON(resp)
	convey.Convey("BangumiOff", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			pn     = int(1)
			ps     = int(10)
			typ    = int8(1)
			appkey = ""
			ts     = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			resp, err := d.BangumiOff(c, pn, ps, typ, appkey, ts)
			ctx.Convey("Then err should be nil.resp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(resp, convey.ShouldNotBeEmpty)
			})
		})
	})
}
