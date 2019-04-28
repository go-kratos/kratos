package service

import (
	"testing"

	"go-common/app/interface/openplatform/monitor-end/model"

	. "github.com/smartystreets/goconvey/convey"
)

// TestFELog .
func TestFELog(t *testing.T) {
	var (
		tcs = []TestCase{
			TestCase{
				tag:      "logtype1 normal",
				testData: `[{"level":"info","logtype":1,"url":"https://show.bilibili.com/api/ticket/district/geocoder","status":200,"cost":323,"traceid_end":"2565683971","traceid_svr":"2565683971"}]`,
				expected: 0,
			},
			TestCase{
				tag:      "logtype1 empty data",
				testData: "",
				expected: 1,
			},
			TestCase{
				tag:      "logtype1 valide data2",
				testData: `[{"level":"info","logtype":1,"url":"https://show.bilibili.com/api/ticket/district/geocoder","status":"200","cost":"323","traceid_end":2565683971,"traceid_svr":"2565683971"}]`,
				expected: 0,
			},
			TestCase{
				tag:      "logtype1 valid data2",
				testData: `[{"level":"info","logtype":1,"url":111,"status":"200","cost":"323","traceid_end":2565683971,"traceid_svr":"2565683971"}]`,
				expected: 0,
			},
			TestCase{
				tag:      "logtype1 valid data3",
				testData: `[{"level":"error","logtype":1,"url":111,"status":"200","cost":"323","traceid_end":2565683971,"traceid_svr":"2565683971"}]`,
				expected: 0,
			},
			TestCase{
				tag:      "logtype2 data",
				testData: `[{"level":"info","logtype":2,"url":"https://show.bilibili.com/platform/home.html","navigationStart":0,"redirectStart":0,"redirectEnd":0,"fetchStart":2,"domainLookupStart":2,"domainLookupEnd":2,"connectStart":2,"secureConnectionStart":0,"connectEnd":2,"requestStart":71,"responseStart":73,"responseEnd":89,"domLoading":107,"domInteractive":442,"domContentLoadedEventStart":442,"domContentLoadedEventEnd":471,"domComplete":873,"loadEventStart":873,"loadEventEnd":873,"firstPaint":301,"firstContentfulPaint":515}]`,
				expected: 0,
			},
		}
	)
	for _, tc := range tcs {
		Convey(tc.tag, t, func() {
			err := svr.Report(ctx, &model.LogParams{Source: "aaa", Log: tc.testData}, 111, "1.1.1.1", "", "")
			if tc.expected == 0 {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldNotBeNil)
			}
		})
	}
}

// TestAPPLog .
func TestAPPLog(t *testing.T) {
	var (
		tcs = []TestCase{
			TestCase{
				tag:      "logtype1 normal",
				testData: `{"request_uri":"\/log\/mobile?ios","time_iso":"1539331278223","ip":"163.142.141.67","version":"2","buvid":"6baad8d5278a2982d204cafd403f089e","fts":"1525072552","proid":"1","chid":"AppStore","pid":"11","brand":"Apple","deviceid":"6baad8d5278a2982d204cafd403f089e","model":"iPhone 7","osver":"11.4.1","ctime":"20181012160107","mid":"18260473","ver":"5.32(8170)","net":"1","oid":"","product":"payment","createtime":"20181012-16:01:07.029GMT+08:00","event":"payment_iap","sub_event":"fetch_products","log_type":"16","duration":"0","message":"\u4eceApple\u83b7\u53d6\u6240\u6709products","result":"0","ext_json":"{\"productIds\":\"tv.danmaku.bilianimex68Bcoin,tv.danmaku.bilianimexnewpanel4998Bcoin,tv.danmaku.bilianimex3BigBcoin,tv.danmaku.bilianimex3VIPbf1,tv.danmaku.bilianimexnewpanel158Bcoin,tv.danmaku.bilianimexnewpanel1598Bcoin,tv.danmaku.bilianimex998Bcoin,tv.danmaku.bilianimexnewpanel648Bcoin,tv.danmaku.bilianimexnewpanel68Bcoin,tv.danmaku.bilianimex12VIPbf1,tv.danmaku.bilianimex3VIP,tv.danmaku.bilianimex18Bcoin,tv.danmaku.bilianimex12VIP,tv.danmaku.bilianimexnewpanel388Bcoinbf","traceid":"","desc":"","network":"1"}`,
				expected: 0,
			},
			TestCase{
				tag: "logtype1 normal",
				testData: `	{"request_uri":"\/log\/mobile?android","time_iso":"1539335958928","ip":"39.181.159.21","version":"2","buvid":"HCwYL0pzQSBDIhIrV2VXZVdmVGVcaggxBnoGinfoc","fts":"1538901223","proid":"1","chid":"vivo","pid":"13","brand":"vivo","deviceid":"HCwYL0pzQSBDIhIrV2VXZVdmVGVcaggxBnoGdkEnSiZUIkUnUjFQZ1IjTA","model":"vivo Y75","osver":"7.1.1","ctime":"20181012171917","mid":"169534145","ver":"5.32.0","net":"1","oid":"46002","product":"music","createtime":"1539335957529","event":"network","sub_event":"https:\/\/api.bilibili.com\/audio\/music-service-c\/url","log_type":"10101","duration":"277","message":"{\"code\":0,\"msg\":\"success\",\"data\":{\"sid\":507989,\"type\":2,\"info\":\"\",\"timeout\":10800,\"size\":3463341,\"cdns\":[\"https:\/\/upos-hz-mirrorkodou.acgvideo.com\/ugaxcode\/i180908tx2zju6l33laab3311t7j7gdi-320k.m4a?deadline=1539346757&platform=android&upsig=558186152304d9145b77172de6c8941d\",\"https:\/\/upos-hz-mirrorossu.acgvideo.com\/ugaxcode\/i180908tx2zju6l33laab3311t7j7gdi-320k.m4a?deadline=1539346757","result":"1","ext_json":"{\"code\":\"200\",\"detail\":\"{\\\"t_befSendReq\\\":\\\"153\\\",\\\"t_parse\\\":\\\"7\\\",\\\"t_ttfb\\\":\\\"33\\\"}\",\"respsize\":\"783\"}","traceid":"","desc":"access_key=f71acb8feb40b3b0effec5d84b9503fc&appkey=1d8b6e7d45233436&build=5320000&mid=169534145&mobi_app=android&platform=android&privilege=2&quality=2&songid=507989&ts=1539335957&sign=480ade7c9622d846e6d91499fdd8fc80","network":"1"}`,
				expected: 0,
			},
			TestCase{
				tag: "logtype1 error duration",
				testData: `	{"request_uri":"\/log\/mobile?android","time_iso":"1539335958928","ip":"39.181.159.21","version":"2","buvid":"HCwYL0pzQSBDIhIrV2VXZVdmVGVcaggxBnoGinfoc","fts":"1538901223","proid":"1","chid":"vivo","pid":"13","brand":"vivo","deviceid":"HCwYL0pzQSBDIhIrV2VXZVdmVGVcaggxBnoGdkEnSiZUIkUnUjFQZ1IjTA","model":"vivo Y75","osver":"7.1.1","ctime":"20181012171917","mid":"169534145","ver":"5.32.0","net":"1","oid":"46002","product":"music","createtime":"1539335957529","event":"network","sub_event":"https:\/\/api.bilibili.com\/audio\/music-service-c\/url","log_type":"10101","duration":"277a","message":"{\"code\":0,\"msg\":\"success\",\"data\":{\"sid\":507989,\"type\":2,\"info\":\"\",\"timeout\":10800,\"size\":3463341,\"cdns\":[\"https:\/\/upos-hz-mirrorkodou.acgvideo.com\/ugaxcode\/i180908tx2zju6l33laab3311t7j7gdi-320k.m4a?deadline=1539346757&platform=android&upsig=558186152304d9145b77172de6c8941d\",\"https:\/\/upos-hz-mirrorossu.acgvideo.com\/ugaxcode\/i180908tx2zju6l33laab3311t7j7gdi-320k.m4a?deadline=1539346757","result":"1","ext_json":"{\"code\":\"200\",\"detail\":\"{\\\"t_befSendReq\\\":\\\"153\\\",\\\"t_parse\\\":\\\"7\\\",\\\"t_ttfb\\\":\\\"33\\\"}\",\"respsize\":\"783\"}","traceid":"","desc":"access_key=f71acb8feb40b3b0effec5d84b9503fc&appkey=1d8b6e7d45233436&build=5320000&mid=169534145&mobi_app=android&platform=android&privilege=2&quality=2&songid=507989&ts=1539335957&sign=480ade7c9622d846e6d91499fdd8fc80","network":"1"}`,
				expected: 0,
			},
			TestCase{
				tag: "logtype1 no json",
				testData: `	{request_uri":"\/log\/mobile?android","time_iso":"1539335958928","ip":"39.181.159.21","version":"2","buvid":"HCwYL0pzQSBDIhIrV2VXZVdmVGVcaggxBnoGinfoc","fts":"1538901223","proid":"1","chid":"vivo","pid":"13","brand":"vivo","deviceid":"HCwYL0pzQSBDIhIrV2VXZVdmVGVcaggxBnoGdkEnSiZUIkUnUjFQZ1IjTA","model":"vivo Y75","osver":"7.1.1","ctime":"20181012171917","mid":"169534145","ver":"5.32.0","net":"1","oid":"46002","product":"music","createtime":"1539335957529","event":"network","sub_event":"https:\/\/api.bilibili.com\/audio\/music-service-c\/url","log_type":"10101","duration":"277","message":"{\"code\":0,\"msg\":\"success\",\"data\":{\"sid\":507989,\"type\":2,\"info\":\"\",\"timeout\":10800,\"size\":3463341,\"cdns\":[\"https:\/\/upos-hz-mirrorkodou.acgvideo.com\/ugaxcode\/i180908tx2zju6l33laab3311t7j7gdi-320k.m4a?deadline=1539346757&platform=android&upsig=558186152304d9145b77172de6c8941d\",\"https:\/\/upos-hz-mirrorossu.acgvideo.com\/ugaxcode\/i180908tx2zju6l33laab3311t7j7gdi-320k.m4a?deadline=1539346757","result":"1","ext_json":"{\"code\":\"200\",\"detail\":\"{\\\"t_befSendReq\\\":\\\"153\\\",\\\"t_parse\\\":\\\"7\\\",\\\"t_ttfb\\\":\\\"33\\\"}\",\"respsize\":\"783\"}","traceid":"","desc":"access_key=f71acb8feb40b3b0effec5d84b9503fc&appkey=1d8b6e7d45233436&build=5320000&mid=169534145&mobi_app=android&platform=android&privilege=2&quality=2&songid=507989&ts=1539335957&sign=480ade7c9622d846e6d91499fdd8fc80","network":"1"}`,
				expected: 1,
			},
			TestCase{
				tag: "logtype1 normal",
				testData: `	{"request_uri":"\/log\/mobile?android","time_iso":"1539335958928","ip":"39.181.159.21","version":"2","buvid":"HCwYL0pzQSBDIhIrV2VXZVdmVGVcaggxBnoGinfoc","fts":"1538901223","proid":"1","chid":"vivo","pid":"13","brand":"vivo","deviceid":"HCwYL0pzQSBDIhIrV2VXZVdmVGVcaggxBnoGdkEnSiZUIkUnUjFQZ1IjTA","model":"vivo Y75","osver":"7.1.1","ctime":"20181012171917","mid":"169534145","ver":"5.32.0","net":"1","oid":"46002","product":"music","createtime":"1539335957529","event":"network","sub_event":"https:\/\/api.bilibili.com\/audio\/music-service-c\/url","log_type":"10101","duration":"277","message":"{\"code\":0,\"msg\":\"success\",\"data\":{\"sid\":507989,\"type\":2,\"info\":\"\",\"timeout\":10800,\"size\":3463341,\"cdns\":[\"https:\/\/upos-hz-mirrorkodou.acgvideo.com\/ugaxcode\/i180908tx2zju6l33laab3311t7j7gdi-320k.m4a?deadline=1539346757&platform=android&upsig=558186152304d9145b77172de6c8941d\",\"https:\/\/upos-hz-mirrorossu.acgvideo.com\/ugaxcode\/i180908tx2zju6l33laab3311t7j7gdi-320k.m4a?deadline=1539346757","result":"0","ext_json":"{\"code\":\"200\",\"detail\":\"{\\\"t_befSendReq\\\":\\\"153\\\",\\\"t_parse\\\":\\\"7\\\",\\\"t_ttfb\\\":\\\"33\\\"}\",\"respsize\":\"783\"}","traceid":"","desc":"access_key=f71acb8feb40b3b0effec5d84b9503fc&appkey=1d8b6e7d45233436&build=5320000&mid=169534145&mobi_app=android&platform=android&privilege=2&quality=2&songid=507989&ts=1539335957&sign=480ade7c9622d846e6d91499fdd8fc80","network":"1"}`,
				expected: 0,
			},
		}
	)
	for _, tc := range tcs {
		Convey(tc.tag, t, func() {
			err := svr.Report(ctx, &model.LogParams{Source: "test", Log: tc.testData, IsAPP: 1}, 111, "1.1.1.1", "", "")
			if tc.expected == 0 {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldNotBeNil)
			}
		})
	}
}
