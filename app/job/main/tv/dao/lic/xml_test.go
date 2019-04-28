package lic

import (
	"encoding/json"
	"testing"

	model "go-common/app/job/main/tv/model/pgc"

	"github.com/smartystreets/goconvey/convey"
)

func TestLicPrepareXML(t *testing.T) {
	licStr := `{"TId":"cd-IqonHtgOndMAywoPLezvXwYmUxcHJ","InputTime":"20180904","Sign":"timer-import_BILIBILI","XMLData":{"Service":{"ID":"dataSync","Head":{"TradeID":"cd-IqonHtgOndMAywoPLezvXwYmUxcHJ","Date":"2018-09-04","Count":0},"Body":{"ProgramSetList":{"ProgramSet":[{"ProgramSetID":"xds296","ProgramSetName":"我家浴缸的二三事-Test","ProgramSetClass":"基腐,日常,泡面,治愈","ProgramSetType":"电影","ProgramSetPoster":"http://i0.hdslb.com/bfs/bangumi/23fbb5ece1d3adb8700988c02e8e97f30bbfbf33.jpg","Portrait":"","Producer":"","PublishDate":"2014-10-06","Copyright":"bilibili","ProgramCount":13,"CREndData":"1970-01-01","DefinitionType":"SD","CpCode":"BILIBILI","PayStatus":0,"PrimitiveName":"オレん家のフロ事情","Alias":"我家浴室的二三事,我家浴室的现况,オレん家のフロ事情","Zone":"日本","LeadingRole":"若狭：梅原裕一郎\n龙己：岛崎信长\n鹰巢：铃木达央\n真木：津田健次郎\n三国：花江夏树\n霞：木户衣吹\n阿比留：川原庆久","ProgramSetDesc":"一个人无忧无虑独自生活的男子高中生龙己，某日救下了一位倒在河边的美青年，没想到这位青年竟然是人鱼！于是这位人鱼先生似乎很中意龙己的浴缸，变住了下来！基友和卖萌人鱼同居的故事开始上演！","Staff":"原作：いときち\n监督：青井小夜\n演出：青井小夜\n脚本：绫奈由仁子\n系列构成：绫奈由仁子\n角色设计：羽田浩二\n色彩设计：小鹿绘里\n摄影监督：藤坂めぐみ\n美术监督：永吉幸树\n音响监督：小泉纪介\n编辑：斋藤朱里\n动画制作：旭PRODUCTION","ProgramList":{"Program":null}}]}}}}}`
	license := &model.License{}
	json.Unmarshal([]byte(licStr), &license)
	convey.Convey("PrepareXML", t, func(ctx convey.C) {
		body := PrepareXML(license)
		ctx.Convey("Then body should not be nil.", func(ctx convey.C) {
			ctx.So(body, convey.ShouldNotBeNil)
		})
	})
}
