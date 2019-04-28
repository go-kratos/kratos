package dao

import (
	"testing"

	"go-common/app/admin/ep/marthe/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	buglyCookie = &model.BuglyCookie{
		QQAccount:  246802468,
		Cookie:     "eas_sid=6105G3k2D0Q428f5u0c6B527B0; pgv_pvi=95820800; RK=kfyEYkBgS6; btcu_id=c64abeb0e4f6485712c0bb79bf16c19d5b6aa5326be74; vc=vc-01c6b914-d8f9-4449-adcb-abba3ebe137c; vc.sig=b8t0PNSdSX1m4wOr-B12whzwaf75BlJjOjz5Jy7YTkE; _ga=GA1.2.1955813367.1533715741; pgv_si=s229603328; _qpsvr_localtk=0.0198901072762101; o_cookie=972360526; pac_uid=1_972360526; csrfToken=gekJ-F5QdGgVTZqLC0NiBEOA; ptui_loginuin=1211712225; ptisp=ctc; ptcz=f45f877d04ce6b659e432a158d35cbc9dea2c565d17eb6ee23640a0c7f82aaf9; uin=o1211712225; skey=@FGItvXrQ6; pt2gguin=o1211712225; IED_LOG_INFO2=userUin%3D1211712225%26nickName%3D%2525E5%2525B0%25258F%2525E7%2525BE%25258E+%26userLoginTime%3D1545802497; midas_openid=1211712225; midas_openkey=@FGItvXrQ6; pgv_info=ssid=s2529006206&pgvReferrer=; pgv_pvid=7939868100; NODINX_SESS=7XIt-RXcFpUaAwKwVFHFbIsssGiryDAw_dF_oP1uVFP2V5vV95jh92eADSSJIq0v; token-skey=771060e7-cd54-f0b3-960c-c8fe485fde10; token-lifeTime=1545828885; bugly_session=eyJpdiI6IldLOHM2V2lhNXFyemdMV1d6YXQ0SHc9PSIsInZhbHVlIjoiekt5UllBZWU4OEltSDVzTzJOeHRESjdQMWY5Y1wveEpYbUlDNmxrV25XTHR3ME5RMkRUdk9VaGlKbGFrQ0cxc2xoUzBOVXdCM0hzVWZIemlFR1BLZXJnPT0iLCJtYWMiOiJhZTI5ZmVjNmVjNzZjMWI2MTMyM2U4NWE5MGZiNWMxMjQzZmEzMWEyMGZhMTcxZjg1N2FiOTY4OTgxNWZjMDExIn0%3D; referrer=eyJpdiI6Im9FZ00yMHdsS2hIeHp3UERSaFVhWlE9PSIsInZhbHVlIjoiZXN6dmZFWmJ4V3R6UmordnowVXZkMXdhbm8zN3QrNzVcL2NSc1I0eWw1ZUVYbVFvTnlwdDB2QWVoaXp4VmZxY2tFV2VSdDIrWG40bEpqb3hvWTZmaVAwXC9vR1JqNEE5NG1MQnlkR1dvV1dkWitSakV6RjV1dWF4dEtzbGpXRFhsNW10SEhrSDVrZk1tRE9EXC9zUEVBRGxwSzhoTHRzSHhuTktFV1g1ckpOTEo0PSIsIm1hYyI6IjAzZWJiMjQ0YjkyNmUyYTk2MDRmNTdjYjY2OWYwNzIzZjZjMmNiMzU0NWRhZmExZWFhYWUzMGFiMTI2MDI4NzIifQ%3D%3F",
		Token:      "1768129694",
		UsageCount: 0,
		Status:     model.BuglyCookieStatusEnable,
	}

	queryBuglyCookiesRequest = &model.QueryBuglyCookiesRequest{
		Pagination: model.Pagination{
			PageSize: 10,
			PageNum:  1,
		},
		QQAccount: buglyCookie.QQAccount,
	}
)

func Test_Bugly_cookie(t *testing.T) {
	Convey("test insert bugly cookie", t, func() {
		err := d.InsertCookie(buglyCookie)
		So(err, ShouldBeNil)
	})

	Convey("test Update Cookie Status", t, func() {
		err := d.UpdateCookieStatus(buglyCookie.ID, model.BuglyCookieStatusDisable)
		So(err, ShouldBeNil)
	})

	Convey("test Update Cookie Usage Count", t, func() {
		err := d.UpdateCookieUsageCount(buglyCookie.ID, 5)
		So(err, ShouldBeNil)
	})

	Convey("test Find Cookies", t, func() {
		total, buglyCookies, err := d.FindCookies(queryBuglyCookiesRequest)
		So(err, ShouldBeNil)
		So(total, ShouldBeGreaterThan, 0)
		So(len(buglyCookies), ShouldBeGreaterThan, 0)
	})
}
