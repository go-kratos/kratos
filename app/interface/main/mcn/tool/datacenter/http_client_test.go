package datacenter

import (
	"context"
	"crypto/md5"
	"fmt"
	"go-common/library/log"
	"net/url"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Init(nil)
	os.Exit(m.Run())
}

func TestSign(t *testing.T) {
	var params = url.Values{
		"timestamp": {"2018-11-19 18:50:28"},
		"version":   {"1.0"},
	}
	var result = "appKey12345timestamp2018-11-19 18:50:28version1.0"
	var md5hash = md5.New()
	var conf = ClientConfig{
		Key: "12345", Secret: "56473",
	}
	md5hash.Write([]byte(conf.Secret + result + conf.Secret))
	var signString = fmt.Sprintf("%X", md5hash.Sum(nil))
	t.Logf("sign string=%s", signString)
	c := New(&conf)
	s, e := c.sign(params)
	if e != nil {
		t.Errorf("err happend, %+v", e)
		t.FailNow()
		return
	}
	if s != signString {
		t.Logf("fail, expect=%s, get=%s", signString, s)
		t.FailNow()
		return
	}
}

func TestAPI(t *testing.T) {
	var q = &Query{}
	q.Select("id, day").
		Where(
			ConditionMapType{"day": ConditionLte("2018-10-28")}).
		Limit(20, 0).Order("day")
	var params = url.Values{
		//"timestamp": {"2018-11-19 18:50:28"},
		//"version": {"1.0"},
		"query": {q.String()},
	}
	var conf = ClientConfig{
		Key: "b9739fc84d087c4b3c1aa297d01999e6", Secret: "5018928e83a23c0cc9773f2571de01e5",
	}
	c := New(&conf)
	var res = &Response{}

	var realResult struct {
		Result []struct {
			ID  int    `json:"id"`
			Day string `json:"day"`
		}
	}
	var middle interface{} = &realResult.Result
	res.Result = middle
	c.Debug = true
	var e = c.Get(context.Background(), "http://berserker.bilibili.co/avenger/api/151/query", params, res)
	if e != nil {
		t.Errorf("err=%+v", e)
		t.FailNow()
	}

	t.Logf("result=%+v", res)
	for k, v := range realResult.Result {
		t.Logf("k=%v, v=%+v", k, v)
	}
}
