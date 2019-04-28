package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPgcCond(t *testing.T) {
	var (
		c      = context.Background()
		snType = int32(1)
	)
	convey.Convey("PgcCond", t, func(ctx convey.C) {
		httpMock("GET", d.c.Cfg.RefLabel.PgcAPI).Reply(200).JSON(`{"code":0,"message":"success","result":{"filter":[{"id":"area","name":"地区","value":[{"id":"-1","name":"全部"},{"id":"1","name":"中国大陆"},{"id":"6,7","name":"中国港台"},{"id":"3","name":"美国"},{"id":"2","name":"日本"},{"id":"8","name":"韩国"},{"id":"9","name":"法国"},{"id":"4","name":"英国"},{"id":"15","name":"德国"},{"id":"10","name":"泰国"},{"id":"35","name":"意大利"},{"id":"13","name":"西班牙"},{"id":"5,11,12,14,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52","name":"其他国家"}]},{"id":"style_id","name":"风格","value":[{"id":"-1","name":"全部"},{"id":"460","name":"剧情"},{"id":"480","name":"喜剧"},{"id":"490","name":"爱情"},{"id":"500","name":"动作"},{"id":"510","name":"恐怖"},{"id":"520","name":"科幻"},{"id":"530","name":"犯罪"},{"id":"540","name":"惊悚"},{"id":"550","name":"悬疑"},{"id":"560","name":"奇幻"},{"id":"600","name":"战争"},{"id":"610","name":"动画"},{"id":"620","name":"传记"},{"id":"630","name":"家庭"},{"id":"640","name":"歌舞"},{"id":"650","name":"历史"},{"id":"730","name":"漫画改"}]},{"id":"year","name":"年份","value":[{"id":"-1","name":"全部"},{"id":"2018","name":"2018"},{"id":"2017","name":"2017"},{"id":"2016","name":"2016"},{"id":"2015","name":"2015"},{"id":"2014","name":"2014"},{"id":"2013-2010","name":"2013-2010"},{"id":"2009-2005","name":"2009-2005"},{"id":"2004-2000","name":"2004-2000"},{"id":"90年代","name":"90年代"},{"id":"80年代","name":"80年代"},{"id":"更早","name":"更早"}]},{"id":"season_status","name":"付费","value":[{"id":"-1","name":"全部"},{"id":"1","name":"免费"},{"id":"2,6","name":"付费"},{"id":"4,6","name":"大会员"}]}],"order":{"name":"排序","value":[{"id":"2","name":"播放数量","sort":"0,1"},{"id":"0","name":"更新时间","sort":"0,1"},{"id":"6","name":"上映时间","sort":"0,1"}]}}}`)
		result, err := d.PgcCond(c, snType)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}
