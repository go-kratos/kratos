package service_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	regexpTaobao    = regexp.MustCompile(`￥([\w\s]+)￥`)
	regexpURL       = regexp.MustCompile(`(?:http|https|www)(?:[\s\.:\/\/]{1,})([\w%+:\s\/\.?=]{1,})`)
	regexpWhitelist = regexp.MustCompile(`((acg|im9|bili|gov).*(com|html|cn|tv)|(av\d{8,}|AV\d{8,}))`)
	regexpQQ        = regexp.MustCompile(`(?:[加qQ企鹅号码\s]{2,}|[群号]{1,})(?:[\x{4e00}-\x{9eff}]*)(?:[:，：]?)([\d\s]{6,})`)
	regexpWechat    = regexp.MustCompile(`(?:[加+微＋+➕薇？vV威卫星♥❤姓xX信]{2,}|weixin|weix)(?:[，❤️.\s]?)(?:[\x{4e00}-\x{9eff}]?)(?:[:，：]?)([\w\s]{6,})`)
)

func TestRegexp(t *testing.T) {
	cases := []struct {
		name   string
		regexp *regexp.Regexp
		inputs []struct {
			content         string
			expectedKeyword string
		}
	}{
		{
			"wechat",
			regexpWechat,
			[]struct {
				content         string
				expectedKeyword string
			}{
				{
					"加微微：Leslie9999990",
					"Leslie9999990",
				},
				{
					"开车开车 想看片加微 18250182201",
					"18250182201",
				},
				{
					"要gv的➕威信 RiverLeee",
					"RiverLeee",
				},
				{
					"未删减版＋威信luijixiang",
					"luijixiang",
				},
			},
		},
		{
			"url",
			regexpURL,
			[]struct {
				content         string
				expectedKeyword string
			}{
				{
					`http://fuli94.com/portal.php?x=611649loli资源`,
					"fuli94.com/portal.php?x=611649loli",
				},
				{
					`老司机开车了 懂得上车：http://zh.cilex.com.cn/http://www.xxoo.jp/?x=156053 萝莉福利：http://zh.cilex.com.cn/http://www.xxoo.jp/?x=156053 福利来了，请叫我雷峰：http://zh.cilex.co`,
					"zh.cilex.com.cn/http://www.xxoo.jp/?x=156053 ",
				},
				{
					`http://flba90.com/forum.php?x=671250 http://flba90.com/forum.php?x=671250`,
					"flba90.com/forum.php?x=671250 http://flba90.com/forum.php?x=671250",
				},
			},
		},
		{
			"taobao",
			regexpTaobao,
			[]struct {
				content         string
				expectedKeyword string
			}{
				{
					`AA网的泥膜，专卖店有10元券。65元入手啊。。。 复制这条信息，打开淘宝￥NEVzZCbjQze￥`,
					"NEVzZCbjQze",
				},
				{
					`这家店还有佳雪神鲜水，这玩意没假货吧，过段日子就难说啦。 ----------------- 复制这条信息，￥jhpeZCWwpvq￥ ，打开【手机淘宝】即可查看`,
					"jhpeZCWwpvq",
				},
				{
					`佳雪神鲜水肌底菁华液神仙水面部精华液补水保湿提亮肤色小样试用【包邮】 【在售价】158.00元 【券后价】153.00元 【下单链接】http://c.b1yt.com/h.jJhBrx?cv=CQfzZCeZauK ----------------- 复制这条信息，￥CQfzZCeZauK￥ ，打`,
					"CQfzZCeZauK",
				},
			},
		},
		{
			"qq",
			regexpQQ,
			[]struct {
				content         string
				expectedKeyword string
			}{
				{
					`HMMD各种类型 老司机开车(^・ω・^ ) 加qq2 6 4 8141670`,
					"2 6 4 8141670",
				},
				{
					`欢迎加入新日暮里，企鹅群：450809463`,
					"450809463",
				},
				{
					`想看完整版加企鹅: 2418046299 高清未删减✔ 各种资源应有尽有。 想看完整版加企鹅: 2418046299 高清未删减✔ 各种资源应有尽有。 想看完整版加企鹅: 2418046299 高清未删减✔ 各种资源应有尽有。`,
					" 2418046299 ",
				},
			},
		},
		{
			"whitelist",
			regexpWhitelist,
			[]struct {
				content         string
				expectedKeyword string
			}{
				{
					`http://big.bilibili.com/site/big.html`,
					`bilibili.com/site/big.html`,
				},
				{
					`live.bilibili.com`,
					`bilibili.com`,
				},
				{
					`http://www.gov.cn/`,
					`gov.cn`,
				},
			},
		},
	}
	for _, c := range cases {
		for _, input := range c.inputs {
			t.Run(c.name, func(t *testing.T) {
				assert := assert.New(t)
				k := c.regexp.FindStringSubmatch(input.content)[1]
				assert.Equal(k, input.expectedKeyword, "")
			})
		}
	}
}
