package service

import (
	"sync"
	"testing"

	"go-common/app/service/main/identify/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	once sync.Once
	s    *Service
)

func startService() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}

func TestService_ReadCookiesVal(t *testing.T) {
	Convey("test read cookies val", t, func() {
		vals0 := readCookiesVal("SESSDATA=123456", "SESSDATA")
		So(len(vals0), ShouldEqual, 1)
		So(vals0[0], ShouldEqual, "123456")
		vals1 := readCookiesVal("SESSDATA=123456;sid=456", "SESSDATA")
		So(len(vals1), ShouldEqual, 1)
		So(vals1[0], ShouldEqual, "123456")
		vals2 := readCookiesVal("SESSDATA=123456;sid=456", "NOTEXIST")
		So(len(vals2), ShouldEqual, 1)
		So(vals2[0], ShouldEqual, "")
		vals3 := readCookiesVal("fts=1531043495; im_notify_type_9119166=0; rpdid=iwqspsxomwdoskowppkww; UM_distinctid=1647a1d280b9b7-0a08c4306c3a24-16386952-232800-1647a1d280c12b; LIVE_BUVID=51d9a8b51edd9c9fd23a3c894f447327; LIVE_BUVID__ckMd5=9f17f524c5a8688a; sid=bzfd5icm; pgv_pvi=3323167744; pos=44; im_local_unread_1=0; pgv_si=s4746899456; buvid3=1B00778D-FC70-43F5-8BC9-3C1020636B026695infoc; CURRENT_QUALITY=80; finger=14bc3c4e; bp_t_offset_18478831=149485967370985408; DedeUserID=9119166; DedeUserID__ckMd5=c11479f809369580; SESSDATA=22f15449%2C1536288045%2Ca05af0ad; bili_jct=4801fa7ff4587a8a6efb588be1b1a915; _ga=GA1.2.689403130.1533818051; stardustvideo=1; BANGUMI_SS_6360_REC=243804; CURRENT_FNVAL=8; im_local_unread_9119166=0; im_seqno_9119166=24530; bp_t_offset_9119166=156097786055309344; msource=PCbanner; _dfcaptcha=5f9c3f275dcdf853b0ff78d902bbfdd9", "SESSDATA")
		So(len(vals3), ShouldEqual, 1)
		So(vals3[0], ShouldEqual, "22f15449%2C1536288045%2Ca05af0ad")
	})
}

func TestService_isIntranetIP(t *testing.T) {
	once.Do(startService)
	Convey("", t, func() {
		ok := false
		ok = s.isIntranetIP("10.255.255.255")
		So(ok, ShouldBeTrue)
		ok = s.isIntranetIP("172.31.255.255")
		So(ok, ShouldBeTrue)
		ok = s.isIntranetIP("172.32.255.255")
		So(ok, ShouldBeFalse)
		ok = s.isIntranetIP("192.168.255.255")
		So(ok, ShouldBeTrue)
		ok = s.isIntranetIP("192.169.0.1")
		So(ok, ShouldBeFalse)
	})
}
