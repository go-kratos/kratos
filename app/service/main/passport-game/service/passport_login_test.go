package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-common/library/log"

	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPwdMatches(t *testing.T) {
	Convey("check pwd match when", t, func() {
		Convey("password is incorrect", func() {
			Convey("salt is empty", func() {
				plain := "wrongPassword"
				salt := ""
				cloudRsaPwd := "b602e55952dfc534e222577eee8f469e"

				ok := pwdMatches(plain, salt, cloudRsaPwd)
				So(ok, ShouldBeFalse)
			})

			Convey("salt is not empty", func() {
				plain := "wrongPassword"
				salt := "DC2q1ju5"
				cloudRsaPwd := "3cb9ffbfd219772793d860c97a815e4b"

				ok := pwdMatches(plain, salt, cloudRsaPwd)
				So(ok, ShouldBeFalse)
			})
		})
		Convey("password is correct", func() {
			Convey("salt is empty", func() {
				plain := "123456"
				salt := ""
				cloudRsaPwd := "b602e55952dfc534e222577eee8f469e"

				ok := pwdMatches(plain, salt, cloudRsaPwd)
				So(ok, ShouldBeTrue)
			})

			Convey("salt is not empty", func() {
				plain := "123456"
				salt := "DC2q1ju5"
				cloudRsaPwd := "3cb9ffbfd219772793d860c97a815e4b"

				ok := pwdMatches(plain, salt, cloudRsaPwd)
				So(ok, ShouldBeTrue)
			})
		})
	})
}

// TestService_Login_OK login ok.
func TestService_Login_OK(t *testing.T) {
	once.Do(startService)
	u := "game0033"
	pwdPlain := "123456"
	expectMid := int64(110000139)
	ts := time.Now().Unix()
	rsaPwd, err := s.rsaEncrypt(TsSeconds2Hash(ts), pwdPlain)
	if err != nil {
		t.Errorf("faied to RSA encrypt pwd and ts, rsaEncrypt(%s, %d) error(%v)", pwdPlain, ts, err)
		t.FailNow()
	}
	if res, err := s.Login(context.TODO(), s.appMap[_gameAppKey], 0, u, rsaPwd); err != nil {
		t.Errorf("service.Login() error(%v)", err)
		t.FailNow()
	} else if res == nil || res.Mid != expectMid {
		t.Errorf("res is not correct, expected login token with mid %d but got %v", expectMid, res)
		t.FailNow()
	} else {
		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	}
}

// TestService_Login_BackOriginWhenCloudFail when cloud login not ok,
// should back origin.
func TestService_Login_BackOriginWhenCloudFail(t *testing.T) {
	once.Do(startService)
	u := "test_ss"
	pwdPlain := "123456"
	ts := time.Now().Unix()
	rsaPwd, err := s.rsaEncrypt(TsSeconds2Hash(ts), pwdPlain)
	if err != nil {
		t.Errorf("faied to RSA encrypt pwd and ts, rsaEncrypt(%s, %d) error(%v)", pwdPlain, ts, err)
		t.FailNow()
	}
	if _, err := s.Login(context.TODO(), s.appMap[_gameAppKey], 0, u, rsaPwd); err != nil {
		t.Errorf("back origin falied, service.Login() error(%v)", err)
		t.FailNow()
	}
}

func TestService_LoginOrigin(t *testing.T) {
	once.Do(startService)
	u := "game0033"
	pwdPlain := "123456"
	expectMid := int64(110000139)
	tsHash := TsSeconds2Hash(time.Now().Unix())
	if res, err := s.loginOrigin(context.TODO(), u, pwdPlain, tsHash); err != nil {
		t.Errorf("s.loginOrigin(%s, %s, %s) error(%v)", u, pwdPlain, tsHash, err)
		t.FailNow()
	} else if res == nil || res.Mid != expectMid {
		t.Errorf("res is not correct, expected login token with mid %d but got %v", expectMid, res)
		t.FailNow()
	} else {
		str, _ := json.Marshal(res)
		t.Logf("res: %v", string(str))
	}
}

func TestService_RSADecrypt(t *testing.T) {
	once.Do(startService)
	expect := "123456"
	rsaPwd, err := s.rsaEncrypt(TsSeconds2Hash(time.Now().Unix()), expect)
	if err != nil {
		t.FailNow()
	}
	res, err := s.rsaDecrypt(rsaPwd)
	if err != nil {
		t.Errorf("failed to decrypt, error(%v)", err)
		t.FailNow()
	}
	ts, _ := Hash2TsSeconds(res[:16])
	plain := res[16:]
	if plain != expect {
		t.Errorf("expect plain %s but got %s", expect, plain)
		t.FailNow()
	}
	t.Logf("ts hash: %s, ts: %d, plain: %s", res[:16], ts, res[16:])
}

func BenchmarkService_RSADecrypt(b *testing.B) {
	once.Do(startService)
	minLen := 3
	maxLen := 16
	num := maxLen - minLen + 1
	plains := make([]string, num)
	cipers := make([]string, num)
	for i := 0; i < num; i++ {
		u := strings.Replace(uuid.NewV4().String(), "-", "", 4)
		plains[i] = u[:i+minLen]
		cb, err := rsaEncryptPKCS8(s.cloudRSAKey.pub, []byte(plains[i]))
		if err != nil {
			b.FailNow()
		}
		cipers[i] = base64.StdEncoding.EncodeToString(cb)
	}
	atomicCnt := int64(0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			idx := atomic.AddInt64(&atomicCnt, 1) % int64(num)
			res, err := s.rsaDecrypt(cipers[idx])
			if err != nil {
				b.Errorf("s.rsaDecrypt(%s) error(%v)", cipers[idx], err)
				b.FailNow()
			}
			if res != plains[idx] {
				b.Errorf("expect %s but got %s", plains[idx], res)
				b.FailNow()
			}
		}
	})
	b.ReportAllocs()
}

func (s *Service) rsaEncrypt(timeHash, pwd string) (res string, err error) {
	d, err := rsaEncryptPKCS8(s.cloudRSAKey.pub, []byte(timeHash+pwd))
	if err != nil {
		log.Error("failed to encrypt pwd, error(%v)", err)
		return
	}
	res = base64.StdEncoding.EncodeToString(d)
	return
}
