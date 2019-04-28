package conf

import (
	"net/http"
	"testing"
)

func TestConf_client(t *testing.T) {
	c := initConf()
	testClientValue(t, c)
	testCheckVersion(t, c)
	testUpdate(t, c)
	testDownload(t, c)
	testGetConfig(t, c)
}

func TestClientNew(t *testing.T) {
	initConf()
	if _, err := New(); err != nil {
		t.Errorf("client.New() error(%v)", err)
		t.FailNow()
	}
}

func testClientValue(t *testing.T, c *Client) {
	key := "breaker"
	testUpdate(t, c)
	test1, ok := c.Value(key)
	if !ok {
		t.Errorf("client.Value() error")
		t.FailNow()
	}
	t.Logf("get the result test1(%s)", test1)
}

func testCheckVersion(t *testing.T, c *Client) {
	ver, err := c.checkVersion(_unknownVersion)
	if err != nil && ver == _unknownVersion {
		t.Errorf("client.checkVersion() error(%v) ver(%d)", err, ver)
		t.FailNow()
	}
}

func testDownload(t *testing.T, c *Client) {
	ver := int64(102)
	if err := c.download(ver); err != nil {
		t.Errorf("client.downloda() error(%v) ", err)
		t.FailNow()
	}

}

func testUpdate(t *testing.T, c *Client) {
	data := &data{
		Version: 199,
		Content: "{\"\":{\"name\":\"\",\"data\":{\"breaker\":\"fuck778\",\"degrade\":\"shit233333\"}},\"redis\":{\"name\":\"redis\",\"data\":{\"444\":\"555\",\"address\":\"172.123.0\",\"array\":\"4,12,test,4\",\"float\":\"3.123\",\"router\":\"test=1,fuck=shit,abc=test\",\"switch\":\"true\",\"timeout\":\"30s\"}}}",
		Md5:     "0843192c43148cbbf43aabb24e3e6442",
	}
	if err := c.update(data); err != nil {
		t.Errorf("client.update() error(%v)", err)
		t.FailNow()
	}
}

func testGetConfig(t *testing.T, c *Client) {
	ver := int64(102)
	data, err := c.getConfig(ver)
	if err != nil {
		t.Errorf("client.getconfiig() error(%v)", err)
		t.FailNow()
	}
	t.Logf("get the result data(%v)", data)
}

func initConf() (c *Client) {
	conf.Addr = "172.16.33.134:9011"
	conf.Host = "testHost"
	conf.Path = "./"
	conf.Svr = "config_test"
	conf.Ver = "shsb-docker-1"
	conf.Env = "10"
	conf.Token = "qmVUPwNXnNfcSpuyqbiIBb0H4GcbSZFV"
	//conf.Appoint = "88"
	c = &Client{
		httpCli: &http.Client{Timeout: _httpTimeout},
		event:   make(chan string, 10),
	}
	c.data.Store(make(map[string]*Namespace))
	return
}
