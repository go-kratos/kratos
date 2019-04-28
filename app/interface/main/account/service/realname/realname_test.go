package realname

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s   *Service
	ctx = context.Background()

	rsaPub, rsaPriv = ``, ``
	alipayPub       = ``
	alipayOwnPriv   = ``
)

func TestMain(m *testing.M) {
	var err error
	flag.Set("conf", "../../cmd/account-interface-example.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf, rsaPub, rsaPriv, alipayPub, alipayOwnPriv)

	m.Run()
	os.Exit(0)
}

func TestEncrypt(t *testing.T) {
	var (
		raw []byte
		err error
	)
	if raw, err = ioutil.ReadFile("./test.jpg"); err != nil {
		panic(err)
	}
	Convey("encrypt", t, func() {
		var (
			encrpted []byte
			outputs  []byte
		)
		if encrpted, err = s.mainCryptor.IMGEncrypt(raw); err != nil {
			panic(err)
		}
		if err = ioutil.WriteFile("./encrypted", encrpted, os.ModePerm); err != nil {
			panic(err)
		}

		if outputs, err = s.mainCryptor.IMGDecrypt(encrpted); err != nil {
			panic(err)
		}
		if err = ioutil.WriteFile("./output.jpg", outputs, os.ModePerm); err != nil {
			panic(err)
		}
	})
}

func TestRealname(t *testing.T) {
	Convey("realname", t, func() {
		var (
			mid int64
			err error
		)
		_, err = s.TelInfo(ctx, mid)
		So(err, ShouldBeNil)
		err = s.TelCapture(ctx, mid)
		So(err, ShouldBeNil)
		s.CountryList(ctx)
		s.CardTypes(ctx, "pc", "", "", 1234)
	})
}

func TestRealnameMask(t *testing.T) {
	Convey("realname mask", t, func() {
		var (
			realname1 = ""
			realname2 = "托"
			realname3 = "托米"
			realname4 = "托了一地的米"

			card1 = ""
			card2 = "t"
			card3 = "tom"
			card4 = "tommmmmy233"
		)
		r, c := maskRealnameInfo(realname1, card1)
		So(r, ShouldEqual, "")
		So(c, ShouldEqual, "")

		r, c = maskRealnameInfo(realname2, card2)
		So(r, ShouldEqual, "*")
		So(c, ShouldEqual, "*")

		r, c = maskRealnameInfo(realname3, card3)
		So(r, ShouldEqual, "*米")
		So(c, ShouldEqual, "t**")

		r, c = maskRealnameInfo(realname4, card4)
		So(r, ShouldEqual, "*了一地的米")
		So(c, ShouldEqual, "t********33")
	})
}

func TestAlipay(t *testing.T) {
	Convey("alipay", t, func() {
		// bizno, err := s.alipayInit(context.Background(), "沐阳", "340702199110120011")
		// So(err, ShouldBeNil)
		// t.Log(bizno)

		// url, err := s.alipayCertifyURL(context.Background(), conf.Conf.Realname.Alipay.AppID, bizno)
		// So(err, ShouldBeNil)
		// t.Log(url)

		bizno := "ZM201808153000000757500643766753"

		pass, reason, err := s.alipayQuery(context.Background(), bizno)
		So(err, ShouldBeNil)
		t.Log(pass, reason)
	})
}

func TestGeetest(t *testing.T) {
	Convey("gt", t, func() {
		var (
			mid int64 = 111001723
		)
		// urlStr, err := s.CaptchaGTRegister(ctx, mid, "127.0.0.1", "h5")
		// So(err, ShouldBeNil)
		// t.Log(urlStr)

		// challenge, gt, success, err := s.CaptchaGTRefresh(ctx, mid, "127.0.0.1", "h5", "fake hash")
		// So(err, ShouldBeNil)
		// t.Log(challenge, gt, success)

		var (
			validateArg = &model.ParamRealnameCaptchaGTCheck{
				Remote:    1,
				Challenge: "58fb29d758ef977717b8fd98d5d1371b9c",
				Validate:  "6171ad7e1243a50c8ca9b124cb6ae411",
				Seccode:   "6171ad7e1243a50c8ca9b124cb6ae411|jordan",
			}
		)
		res, err := s.CaptchaGTValidate(ctx, mid, "127.0.0.1", "h5", validateArg)
		So(err, ShouldBeNil)
		t.Log(res)
	})
}

func BenchmarkEncrypt(b *testing.B) {
	var (
		err error
		raw []byte
	)
	if raw, err = ioutil.ReadFile("./test.jpg"); err != nil {
		panic(err)
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err = s.mainCryptor.IMGEncrypt(raw); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkDecrypt(b *testing.B) {
	var (
		err error
		raw []byte
	)
	if raw, err = ioutil.ReadFile("./encrypted"); err != nil {
		panic(err)
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err = s.mainCryptor.IMGDecrypt(raw); err != nil {
				panic(err)
			}
		}
	})
}
