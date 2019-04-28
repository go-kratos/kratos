package service

import "go-common/app/interface/main/answer/model"

// RankShire struct .
type RankShire struct {
	Share    *model.CoolShare
	ViewMore string
	VideoArr []*model.CoolVideo
}

var (
	_rankShire = map[int]*RankShire{
		122: {
			Share: &model.CoolShare{Content: "老司机中的老司机，哲学王中的哲学王。",
				ShortContent: "老司机中的老司机，哲学王中的哲学王。"},
			ViewMore: "http://www.bilibili.com/ranking",
			VideoArr: []*model.CoolVideo{
				{
					URL:      "http://www.bilibili.com/video/av221107/",
					Name:     "乾杯 - ( ゜- ゜)つロ",
					Img:      "https://static-s.bilibili.com/account/img/answer/221107.jpg",
					WatchNum: "85.8万",
					UpNum:    "87515",
				}, {
					URL:      "http://www.bilibili.com/video/av296938/",
					Name:     "我们的BILIBILI",
					Img:      "https://static-s.bilibili.com/account/img/answer/296938.jpg",
					WatchNum: "83.4万",
					UpNum:    "94643",
				},
			},
		},
		124: {
			Share: &model.CoolShare{Content: "在无限大的梦想后面，<br/>无论喜欢动画的你去哪，<br/>那只幻化的蝴蝶，永远会陪伴在你身边。",
				ShortContent: "无限大的梦想后面，那只蝴蝶会永远陪伴着你。"},
			ViewMore: "http://www.bilibili.com/video/bangumi.html",
			VideoArr: []*model.CoolVideo{
				{
					URL:      "http://www.bilibili.com/video/av19781/",
					Name:     "【催泪向MAD】butter-fly (主MIX) 【数码宝贝】",
					Img:      "https://static-s.bilibili.com/account/img/answer/19781.jpg",
					WatchNum: "104.1万",
					UpNum:    "173642",
				}, {
					URL:      "http://www.bilibili.com/video/av2216404/",
					Name:     "海尔兄弟【bilibili正版】",
					Img:      "https://static-s.bilibili.com/account/img/answer/2216404.jpg",
					WatchNum: "40.8万",
					UpNum:    "39693",
				}, {
					URL:      "http://www.bilibili.com/video/av491271/",
					Name:     "【综漫MAD】致一直陪伴着我的二次元——谢谢",
					Img:      "https://static-s.bilibili.com/account/img/answer/491271.jpg",
					WatchNum: "29.8万",
					UpNum:    "19541",
				}, {
					URL:      "http://www.bilibili.com/video/av101457/",
					Name:     "【MAD】因为我是活在二次元的女孩",
					Img:      "https://static-s.bilibili.com/account/img/answer/101457.jpg",
					WatchNum: "53.7万",
					UpNum:    "19034",
				}, {
					URL:      "http://www.bilibili.com/video/av853594/",
					Name:     "致二次元，谢谢你给了我整个世界",
					Img:      "https://static-s.bilibili.com/account/img/answer/853594.png",
					WatchNum: "17.0万",
					UpNum:    "16830",
				},
			},
		},
		127: {
			Share: &model.CoolShare{Content: "在无限大的梦想后面，<br/>无论喜欢动画的你去哪，<br/>那只幻化的蝴蝶，永远会陪伴在你身边。",
				ShortContent: "无限大的梦想后面，那只蝴蝶会永远陪伴着你。"},
			ViewMore: "http://www.bilibili.com/video/bangumi.html",
			VideoArr: []*model.CoolVideo{
				{
					URL:      "http://www.bilibili.com/video/av19781/",
					Name:     "【催泪向MAD】butter-fly (主MIX) 【数码宝贝】",
					Img:      "https://static-s.bilibili.com/account/img/answer/19781.jpg",
					WatchNum: "104.1万",
					UpNum:    "173642",
				},
				{
					URL:      "http://www.bilibili.com/video/av2216404/",
					Name:     "海尔兄弟【bilibili正版】",
					Img:      "https://static-s.bilibili.com/account/img/answer/2216404.jpg",
					WatchNum: "40.8万",
					UpNum:    "39693",
				},
				{
					URL:      "http://www.bilibili.com/video/av491271/",
					Name:     "【综漫MAD】致一直陪伴着我的二次元——谢谢",
					Img:      "https://static-s.bilibili.com/account/img/answer/491271.jpg",
					WatchNum: "29.8万",
					UpNum:    "19541",
				},
				{
					URL:      "http://www.bilibili.com/video/av101457/",
					Name:     "【MAD】因为我是活在二次元的女孩",
					Img:      "https://static-s.bilibili.com/account/img/answer/101457.jpg",
					WatchNum: "53.7万",
					UpNum:    "19034",
				},
				{
					URL:      "http://www.bilibili.com/video/av853594/",
					Name:     "致二次元，谢谢你给了我整个世界",
					Img:      "https://static-s.bilibili.com/account/img/answer/853594.png",
					WatchNum: "17.0万",
					UpNum:    "16830",
				},
			},
		},
		126: {
			Share: &model.CoolShare{Content: "给你最想听的音乐，让耳朵来一场旅行，我是旋律的向导。",
				ShortContent: "给你最想听的音乐，让耳朵来一场旅行，我是旋律的向导。"},
			ViewMore: "http://www.bilibili.com/video/music-vocaloid-1.html",
			VideoArr: []*model.CoolVideo{
				{
					URL:      "http://www.bilibili.com/video/av2129461/",
					Name:     "洛天依，言和原创《普通DISCO》",
					Img:      "https://static-s.bilibili.com/account/img/answer/2129461.jpg",
					WatchNum: "173.7万",
					UpNum:    "53752",
				},
				{
					URL:      "http://www.bilibili.com/video/av644136/",
					Name:     "【洛天依翻唱】跳蛋的性福理论【手书PV】【JumPingEgG】",
					Img:      "https://static-s.bilibili.com/account/img/answer/644136.jpg",
					WatchNum: "64.6万",
					UpNum:    "25293",
				},
				{
					URL:      "http://www.bilibili.com/video/av2075941/",
					Name:     "【洛天依古风原创曲】权御天下【原创PV付】",
					Img:      "https://static-s.bilibili.com/account/img/answer/2075941.jpg",
					WatchNum: "128.7万",
					UpNum:    "27589",
				},
				{
					URL:      "http://www.bilibili.com/video/av482844/",
					Name:     "【洛天依原创】一半一半",
					Img:      "https://static-s.bilibili.com/account/img/answer/482844.jpg",
					WatchNum: "50.5万",
					UpNum:    "22569",
				},
				{
					URL:      "http://www.bilibili.com/video/av556019/",
					Name:     "【niconico超会议2現場版】威风堂々【Vmoe字幕组】",
					Img:      "https://static-s.bilibili.com/account/img/answer/556019.jpg",
					WatchNum: "67.4万",
					UpNum:    "19833",
				},
			},
		},
		123: {
			Share: &model.CoolShare{Content: "给你最想听的音乐，让耳朵来一场旅行，我是旋律的向导。",
				ShortContent: "给你最想听的音乐，让耳朵来一场旅行，我是旋律的向导。"},
			ViewMore: "http://www.bilibili.com/video/music.html",
			VideoArr: []*model.CoolVideo{
				{
					URL:      "http://www.bilibili.com/video/av643809/",
					Name:     "茶太--团子大家族 现场版",
					Img:      "https://static-s.bilibili.com/account/img/answer/643809.jpg",
					WatchNum: "14.8万",
					UpNum:    "7947",
				},
				{
					URL:      "http://www.bilibili.com/video/av345249/",
					Name:     "【AMV、時代的眼泪、泪腺崩坏、燃烧殆尽！】最后的Butterfly重制版",
					Img:      "https://static-s.bilibili.com/account/img/answer/345249.png",
					WatchNum: "31.7万",
					UpNum:    "39937",
				},
				{
					URL:      "http://www.bilibili.com/video/av736852/",
					Name:     "看到鼓手时我跪下尿了一地！",
					Img:      "https://static-s.bilibili.com/account/img/answer/736852.jpg",
					WatchNum: "91.2万",
					UpNum:    "11302",
				},
				{
					URL:      "http://www.bilibili.com/video/av1393947/",
					Name:     "电二胡的咆哮 【致童年—数码宝贝】butterfly！！",
					Img:      "https://static-s.bilibili.com/account/img/answer/1393947.jpg",
					WatchNum: "35.6万",
					UpNum:    "14199",
				},
				{
					URL:      "http://www.bilibili.com/video/av1507163/",
					Name:     "你们要的小苹果交响版",
					Img:      "https://static-s.bilibili.com/account/img/answer/1507163.jpg",
					WatchNum: "30.9万",
					UpNum:    "7092",
				},
			},
		},
		121: {
			Share: &model.CoolShare{Content: "常存好奇之心，<br/>新技能 get 是他们的座右铭。",
				ShortContent: "常存好奇之心，新技能 get 是他们的座右铭。"},
			ViewMore: "http://www.bilibili.com/video/tech-future-military-1.html",
			VideoArr: []*model.CoolVideo{
				{
					URL:      "http://www.bilibili.com/video/av2075767/",
					Name:     "巴雷特狙击枪后座力有多强？",
					Img:      "https://static-s.bilibili.com/account/img/answer/2075767.jpg",
					WatchNum: "33.7万",
					UpNum:    "2002",
				},
				{
					URL:      "http://www.bilibili.com/video/av1952604/",
					Name:     "【张召忠】 印度史上最搞笑的大阅兵 铁血军情20150201",
					Img:      "https://static-s.bilibili.com/account/img/answer/1952604.jpg",
					WatchNum: "28.3万",
					UpNum:    "12463",
				},
				{
					URL:      "http://www.bilibili.com/video/av927165/",
					Name:     "【军武次位面】第九期：十大战列舰",
					Img:      "https://static-s.bilibili.com/account/img/answer/927165.png",
					WatchNum: "32.6万",
					UpNum:    "10880",
				},
			},
		},

		125: {
			Share: &model.CoolShare{Content: "你就是来自那个世界的勇士吗？<br/>果然你在这个虚拟世界有着强大的力量啊！<br/>准备好攻略这场战斗了么？那么…Link Start!",
				ShortContent: "准备好攻略这场战斗了么？那么…Link Start!"},
			ViewMore: "http://www.bilibili.com/video/game.html",
			VideoArr: []*model.CoolVideo{
				{
					URL:      "http://www.bilibili.com/video/av1293479/",
					Name:     "【若风噩梦人机】瞬间爆炸的恐怖电脑们！",
					Img:      "https://static-s.bilibili.com/account/img/answer/1293479.jpg",
					WatchNum: "73.1万",
					UpNum:    "25452",
				},
				{
					URL:      "http://www.bilibili.com/video/av1561567/",
					Name:     "【散人】大型励志剧 娱乐圈小助理养成计划（更新至P20 遇龙2杀青）",
					Img:      "https://static-s.bilibili.com/account/img/answer/1561567.jpg",
					WatchNum: "306.5万",
					UpNum:    "581921",
				},
				{
					URL:      "http://www.bilibili.com/video/av862182/",
					Name:     "【文明5】美丽新世界神级实况（7P完结）",
					Img:      "https://static-s.bilibili.com/account/img/answer/862182.jpg",
					WatchNum: "14.4万",
					UpNum:    "9325",
				},
				{
					URL:      "http://www.bilibili.com/video/av885977/",
					Name:     "<Mugen>狂下左右节操全无大会最终章-燃烧热情吧！向着梦想的彼方！",
					Img:      "https://static-s.bilibili.com/account/img/answer/885977.jpg",
					WatchNum: "14.2万",
					UpNum:    "11825",
				},
				{
					URL:      "http://www.bilibili.com/video/av2269587/",
					Name:     "LOL：这也能翻盘？史上最奇葩的翻盘，这竟然是钻石排位",
					Img:      "https://static-s.bilibili.com/account/img/answer/2269587.jpg",
					WatchNum: "99.3万",
					UpNum:    "37457",
				},
			},
		},
		129: {
			Share: &model.CoolShare{Content: "常存好奇之心，<br/>新技能 get 是他们的座右铭。",
				ShortContent: "常存好奇之心，新技能 get 是他们的座右铭。"},
			ViewMore: "http://www.bilibili.com/video/technology.html",
			VideoArr: []*model.CoolVideo{
				{
					URL:      "http://www.bilibili.com/video/av1787919/",
					Name:     "史上最无节操的手机架子鼓教程 ，结尾开大招。",
					Img:      "https://static-s.bilibili.com/account/img/answer/dcc6580b4dcc7e339545018cc312bf76.jpg",
					WatchNum: "29.2万",
					UpNum:    "4744",
				},
				{
					URL:      "http://www.bilibili.com/video/av679319/",
					Name:     "论爱情动作片和真实生活啪啪啪的区别",
					Img:      "https://static-s.bilibili.com/account/img/answer/13752405809fed020ca2372901.jpg",
					WatchNum: "43.3万",
					UpNum:    "3028",
				},
				{
					URL:      "http://www.bilibili.com/video/av2275735/",
					Name:     "Besiege贴吧4月作品精选",
					Img:      "https://static-s.bilibili.com/account/img/answer/d382baaee8f144b924975b05cc592c60.jpg",
					WatchNum: "10.5万",
					UpNum:    "6506",
				},
				{
					URL:      "http://www.bilibili.com/video/av2278660/",
					Name:     "英梨梨&诗羽de绘成方法 五一特别篇【别人君】",
					Img:      "https://static-s.bilibili.com/account/img/answer/fd41bb17a5f512c37f01570ddd994b49.jpg",
					WatchNum: "4.7万",
					UpNum:    "2837",
				},
				{
					URL:      "http://www.bilibili.com/video/av2278660/",
					Name:     "这大概是最好的日语入门教学了吧----五十音学习",
					Img:      "https://static-s.bilibili.com/account/img/answer/8c7436785697cc5c450b33eb93c2353f.jpg",
					WatchNum: "63.6万",
					UpNum:    "104048",
				},
			},
		},
		130: {
			Share: &model.CoolShare{Content: "一天不看片会死星人就是你吗？<br/>大开脑洞YY主人公也大丈夫~<br/>在异次元的空间里跟我们一起做梦吧。",
				ShortContent: "一天不看片会死星人就是你吗？在异次元的空间里跟我们一起做梦吧。"},
			ViewMore: "http://www.bilibili.com/video/tv-presentation-1.html",
			VideoArr: []*model.CoolVideo{
				{
					URL:      "http://www.bilibili.com/video/av2019928/",
					Name:     "泰剧《不一样的美男》中字第一集@天府泰剧",
					Img:      "https://static-s.bilibili.com/account/img/answer/daab100da59667394e70822e3e63d254.jpg",
					WatchNum: "203.5万",
					UpNum:    "219407",
				},
				{
					URL:      "http://www.bilibili.com/video/av639407/",
					Name:     "【熟肉】半泽直树 01【人人字幕】",
					Img:      "https://static-s.bilibili.com/account/img/answer/076d81509a1c93a4569e942b281f46c1.jpg",
					WatchNum: "141.8万",
					UpNum:    "61416",
				},
				{
					URL:      "http://www.bilibili.com/video/av1999475/",
					Name:     "【国产】少年包青天 第一部 2000 40集全集",
					Img:      "https://static-s.bilibili.com/account/img/answer/33a6c9cf7b38bb08c4da7ec5ac3b12f4.jpg",
					WatchNum: "45.8万",
					UpNum:    "167676",
				},
			},
		},
		128: {
			Share: &model.CoolShare{Content: "从未见过如此才华横溢之人，<br/>他们的技术、创意和努力，<br/>使「鬼畜」成为一种艺术。",
				ShortContent: "他们的技术、创意和努力，使「鬼畜」成为一种艺术。"},
			ViewMore: "http://www.bilibili.com/video/kichiku.html",
			VideoArr: []*model.CoolVideo{
				{
					URL:      "http://www.bilibili.com/video/av75179/",
					Name:     "【葛平金曲】循环（完整版）",
					Img:      "https://static-s.bilibili.com/account/img/answer/1301416217-61b.jpg",
					WatchNum: "45.1万",
					UpNum:    "21196",
				},
				{
					URL:      "http://www.bilibili.com/video/av1858893/",
					Name:     "全是猴【白金王司猴】",
					Img:      "https://static-s.bilibili.com/account/img/answer/f1d17e3ce7f9e71e1508ced43d6a8656.jpg",
					WatchNum: "111.4万",
					UpNum:    "14230",
				},
				{
					URL:      "http://www.bilibili.com/video/av1076105/",
					Name:     "妮可 妮可 妮可",
					Img:      "https://static-s.bilibili.com/account/img/answer/ee087327bf0a239ad114633e2806fa79.jpg",
					WatchNum: "94.6万",
					UpNum:    "25738",
				},
				{
					URL:      "http://www.bilibili.com/video/av2271112/",
					Name:     "【循环向】跟着雷总摇起来！Are you OK！",
					Img:      "https://static-s.bilibili.com/account/img/answer/2fc528fee5d0cbfb98b266bb7ec3a1ad.jpg",
					WatchNum: "118.1万",
					UpNum:    "18960",
				},
				{
					URL:      "http://www.bilibili.com/video/av794506/",
					Name:     "【元首葛炮】要金坷垃",
					Img:      "https://static-s.bilibili.com/account/img/answer/1d34e9f58def856acf501289d2cacc56.jpg",
					WatchNum: "64.3万",
					UpNum:    "13795",
				},
			}},
	}
)
