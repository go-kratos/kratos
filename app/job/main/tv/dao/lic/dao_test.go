package lic

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/tv/conf"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx     = context.Background()
	d       *Dao
	xmlBody = `inputTime=20180606&sign=timer-import_BILIBILI&tId=UHZFmufgweRpWhqAzToFYMWtYuZhMKCU&xmlData=<?xmlversion="1.0"encoding="UTF-8"?><Serviceid="dataSync"><Head><TradeId>UHZFmufgweRpWhqAzToFYMWtYuZhMKCU</TradeId><Date>2018-06-06</Date><Count>1</Count></Head><Body><programSetList><programSet><programSetId>ugc10100044</programSetId><programSetName>drm</programSetName><programSetClass></programSetClass><programSetType></programSetType><programSetPoster>http://i1.hdslb.com/bfs/archive/diuren.png</programSetPoster><publishDate>2018-01-30</publishDate><copyright>bilibili</copyright><programCount>1</programCount><cREndDate>1970-01-01</cREndDate><definitionType>SD</definitionType><cpCode>BILIBILI</cpCode><payStatus>0</payStatus><primitiveName></primitiveName><alias></alias><zone></zone><leadingRole></leadingRole><programSetDesc>drm</programSetDesc><Staff></Staff><programList><program><programId>ugc10114149</programId><programName>1</programName><programPoster></programPoster><programLength>1448</programLength><publishDate>1970-01-01</publishDate><ifPreview>0</ifPreview><number>1</number><definitionType>SD</definitionType><playCount>0</playCount><drm>0</drm><programMediaList><programMedia><mediaId>ugc10114149</mediaId><playUrl>http://upos-hz-tvshenhe.acgvideo.com/upgcxcode/87/75/41057587/41057587-1-6.mp4</playUrl><definition>SD</definition><htmlUrl>http://upos-hz-tvshenhe.acgvideo.com/upgcxcode/87/75/41057587/41057587-1-6.mp4</htmlUrl></programMedia></programMediaList></program></programList></programSet></programSetList></Body></Service>`
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../../cmd/tv-job-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		if d == nil {
			d = New(conf.Conf)
		}
		f(d)
	}
}

func TestDelEpLic(t *testing.T) {
	Convey("TestDao_CallRetry", t, WithDao(func(d *Dao) {
		res := DelEpLic("ugc", "timer-import_BILIBILI", []int{10109083, 10109084})
		So(len(res), ShouldBeGreaterThan, 0)
		fmt.Println(res)
	}))
}

func TestDao_CallRetry(t *testing.T) {
	Convey("TestDao_CallRetry", t, WithDao(func(d *Dao) {
		res, err := d.CallRetry(ctx, d.conf.Sync.API.AddURL, xmlBody)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestDao_CallLic(t *testing.T) {
	Convey("TestDao_CallLic", t, WithDao(func(d *Dao) {
		result, err := d.callLic(ctx, d.conf.Sync.API.AddURL, xmlBody)
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
	}))
}
