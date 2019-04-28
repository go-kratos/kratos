package dao

import (
	"testing"

	"go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_MallCard(t *testing.T) {
	Convey("normal should get data", t, func() {
		data := `{"code":0,"message":"success","data":{"pageNum":1,"pageSize":2,"size":2,"startRow":1,"endRow":2,"total":2,"pages":1,"list":[{"itemsId":1,"brief":"韩版帅气","isLastestVersion":1,"name":"短袖T恤男学生新款韩版衬衫12","img":["//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg","//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg","//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg","//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg"],"onSaleTime":null,"offSaleTime":null,"price":0,"maxPrice":0,"sales":0,"frozenStock":null,"stock":null,"needUserinfoCollection":[1,2,3],"presaleStartOrderTime":0,"presaleEndOrderTime":0,"depositPrice":0,"deliveryTemplateId":0,"tianmaImg":"","version":3},{"itemsId":5,"brief":"韩版帅气","isLastestVersion":1,"name":"短袖T恤男学生新款韩版衬衫","img":["//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg","//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg","//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg","//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg"],"onSaleTime":null,"offSaleTime":null,"price":0,"maxPrice":0,"sales":0,"frozenStock":null,"stock":null,"needUserinfoCollection":[1,2,3],"presaleStartOrderTime":0,"presaleEndOrderTime":0,"depositPrice":0,"deliveryTemplateId":13566,"tianmaImg":"","version":2}],"prePage":0,"nextPage":0,"isFirstPage":true,"isLastPage":true,"hasPreviousPage":false,"hasNextPage":false,"navigatePages":8,"navigatepageNums":[1],"navigateFirstPage":1,"navigateLastPage":1,"firstPage":1,"lastPage":1}}`
		httpMock("POST", d.c.Cards.MallURL).Reply(200).JSON(data)
		res, err := d.MallCard(ctx(), []int64{1, 5})
		So(err, ShouldBeNil)
		So(res, ShouldResemble, map[int64]*model.MallCard{
			1: &model.MallCard{
				ID:    1,
				Name:  "短袖T恤男学生新款韩版衬衫12",
				Brief: "韩版帅气",
				Images: []string{
					"//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg",
					"//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg",
					"//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg",
					"//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg",
				},
				Price: 0,
			},
			5: &model.MallCard{
				ID:    5,
				Name:  "短袖T恤男学生新款韩版衬衫",
				Brief: "韩版帅气",
				Images: []string{
					"//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg",
					"//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg",
					"//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg",
					"//img10.360buyimg.com/n0/jfs/t5392/234/1745592889/186437/4e9da0f7/5913b8dfNcc393bff.jpg",
				},
				Price: 0,
			},
		})
	})
	Convey("code !=0 should get error", t, func() {
		data := `{"code":-3,"message":"faild","data":{}}`
		httpMock("POST", d.c.Cards.MallURL).Reply(200).JSON(data)
		_, err := d.MallCard(ctx(), []int64{1, 5})
		So(err, ShouldNotBeNil)
	})
}

func Test_TicketCard(t *testing.T) {
	Convey("normal get data", t, func() {
		data := `{"errno":0,"msg":"","data":{"75":{"id":75,"name":"赵丽颖见面会","status":1,"start_time":1500268460,"end_time":1538284460,"performance_image":"//uat-i1.hdslb.com/bfs/openplatform/201707/imrGbwzlkCYUs.jpeg","is_sale":1,"promo_tags":"1-2","stime":"7/17","etime":"9/30","province_name":"上海市","city_name":"上海市","district_name":"浦东新区","venue_name":"梅赛德斯奔驰文化中心","url":"https://show.bilibili.com/m/platform/detail.html?id=75&from=","price_low":0.01,"price_high":500},"80":{"id":80,"name":"演唱会测试C","status":0,"start_time":1501050922,"end_time":1501137326,"performance_image":"//uat-i0.hdslb.com/bfs/openplatform/201707/imXtcy7Kgllz2.jpeg","is_sale":1,"promo_tags":"1-1","stime":"7/26","etime":"7/27","province_name":"上海市","city_name":"上海市","district_name":"浦东新区","venue_name":"文化中心","url":"https://show.bilibili.com/m/platform/detail.html?id=80&from=","price_low":200,"price_high":500}}}`
		httpMock("get", d.c.Cards.TicketURL).Reply(200).JSON(data)
		res, err := d.TicketCard(ctx(), []int64{75, 80})
		So(err, ShouldBeNil)
		So(res, ShouldResemble, map[int64]*model.TicketCard{
			75: &model.TicketCard{
				ID:        75,
				Name:      "赵丽颖见面会",
				Image:     "//uat-i1.hdslb.com/bfs/openplatform/201707/imrGbwzlkCYUs.jpeg",
				StartTime: 1500268460,
				EndTime:   1538284460,
				Province:  "上海市",
				City:      "上海市",
				District:  "浦东新区",
				Venue:     "梅赛德斯奔驰文化中心",
				PriceLow:  0.01,
				URL:       "https://show.bilibili.com/m/platform/detail.html?id=75&from=",
			},
			80: &model.TicketCard{
				ID:        80,
				Name:      "演唱会测试C",
				Image:     "//uat-i0.hdslb.com/bfs/openplatform/201707/imXtcy7Kgllz2.jpeg",
				StartTime: 1501050922,
				EndTime:   1501137326,
				Province:  "上海市",
				City:      "上海市",
				District:  "浦东新区",
				Venue:     "文化中心",
				PriceLow:  200,
				URL:       "https://show.bilibili.com/m/platform/detail.html?id=80&from=",
			},
		})
	})
	Convey("code != 0 should return error", t, func() {
		data := `{"errno":-1,"msg":"","data":{}}`
		httpMock("get", d.c.Cards.TicketURL).Reply(200).JSON(data)
		res, err := d.TicketCard(ctx(), []int64{75, 80})
		So(err, ShouldNotBeNil)
		So(res, ShouldBeNil)
	})
}

func Test_AudioCard(t *testing.T) {
	Convey("normal get data", t, func() {
		data := `{"code":0,"msg":"success","data":{"75":{"song_id":75,"title":"【Hanser】星电感应","up_mid":26609612,"up_name":"siroccox","play_num":17,"reply_num":0,"cover_url":"http://i0.hdslb.com/bfs/test/80740468b108a4f1b98316caa02dc8dcf5976caf.jpg"}}}`
		httpMock("get", d.c.Cards.AudioURL).Reply(200).JSON(data)
		res, err := d.AudioCard(ctx(), []int64{75})
		So(err, ShouldBeNil)
		So(res, ShouldResemble, map[int64]*model.AudioCard{
			75: &model.AudioCard{
				ID:       75,
				Title:    "【Hanser】星电感应",
				UpMid:    26609612,
				UpName:   "siroccox",
				Play:     17,
				Reply:    0,
				CoverURL: "http://i0.hdslb.com/bfs/test/80740468b108a4f1b98316caa02dc8dcf5976caf.jpg",
			},
		})
	})

	Convey("code != 0 should return error", t, func() {
		data := `{"code":-1,"msg":"fail","data":{}}}`
		httpMock("get", d.c.Cards.AudioURL).Reply(200).JSON(data)
		_, err := d.AudioCard(ctx(), []int64{75})
		So(err, ShouldNotBeNil)
	})
}

func Test_BangumiCard(t *testing.T) {
	exp := map[int64]*model.BangumiCard{
		20031: &model.BangumiCard{
			ID:    20031,
			Image: "http://i0.hdslb.com/bfs/bangumi/77605418c0921578c469201d6384d6a32ed218e9.jpg",
			Title: "地狱少女 宵伽",
			Rating: struct {
				Score float64 `json:"score"`
				Count int64   `json:"count"`
			}{
				Score: 0,
				Count: 0,
			},
			Playable:    true,
			FollowCount: 0,
			PlayCount:   0,
		},
	}
	Convey("seasons", t, func() {
		Convey("normal get data", func() {
			data := `{"code":0,"message":"success","result":{"season_map":{"20031":{"allow_review":1,"cover":"http://i0.hdslb.com/bfs/bangumi/77605418c0921578c469201d6384d6a32ed218e9.jpg","is_finish":1,"is_started":1,"media_id":11,"playable":true,"season_id":20031,"season_type":1,"season_type_name":"番剧","title":"地狱少女 宵伽","total_count":13}}}}`
			httpMock("post", d.c.Cards.BangumiURL).Reply(200).JSON(data)
			res, err := d.BangumiCard(ctx(), []int64{20031}, nil)
			So(err, ShouldBeNil)
			So(res, ShouldResemble, exp)
		})
		Convey("code != 0 should return error", func() {
			data := `{"code":-1,"message":"fail","result":{}}`
			httpMock("post", d.c.Cards.BangumiURL).Reply(200).JSON(data)
			_, err := d.BangumiCard(ctx(), []int64{20031}, nil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("eps", t, func() {
		Convey("normal get data", func() {
			data := `{"code":0,"message":"success","result":{"episode_map":{"20031":{"allow_review":1,"cover":"http://i0.hdslb.com/bfs/bangumi/77605418c0921578c469201d6384d6a32ed218e9.jpg","is_finish":1,"is_started":1,"media_id":11,"playable":true,"season_id":20031,"season_type":1,"season_type_name":"番剧","title":"地狱少女 宵伽","total_count":13}}}}`
			httpMock("post", d.c.Cards.BangumiURL).Reply(200).JSON(data)
			res, err := d.BangumiCard(ctx(), nil, []int64{20031})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, exp)
		})
		Convey("code != 0 should return error", func() {
			data := `{"code":-1,"message":"fail","result":{}}`
			httpMock("post", d.c.Cards.BangumiURL).Reply(200).JSON(data)
			_, err := d.BangumiCard(ctx(), nil, []int64{20031})
			So(err, ShouldNotBeNil)
		})
	})
}
