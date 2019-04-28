package conf

import (
	"net/http"
	"testing"
)

func TestConf_client2(t *testing.T) {
	c := initConf2()
	testClientValue2(t, c)
	testCheckVersion2(t, c)
	testDownload2(t, c)
	testGetConfig2(t, c)
}

func testClientValue2(t *testing.T, c *Client) {
	key := "test.toml"
	testDownload2(t, c)
	test1, ok := c.Value2(key)
	if !ok {
		t.Errorf("client.Value() error")
		t.FailNow()
	}
	t.Logf("get the result test1(%s)", test1)
}

func testCheckVersion2(t *testing.T, c *Client) {
	unknow := &ver{Version: _unknownVersion}
	ver, err := c.checkVersion2(unknow)
	if err != nil {
		t.Errorf("client.checkVersion() error(%v) ver(%d)", err, ver)
		t.FailNow()
	}
}

func testDownload2(t *testing.T, c *Client) {
	ver := &ver{Version: 13}
	if err := c.download2(ver, true); err != nil {
		t.Errorf("client.downloda() error(%v) ", err)
		t.FailNow()
	}
}

func testGetConfig2(t *testing.T, c *Client) {
	ver := &ver{Version: 13}
	data, err := c.getConfig2(ver)
	if err != nil {
		t.Errorf("client.getconfiig() error(%v)", err)
		t.FailNow()
	}
	t.Logf("get the result data(%v)", data)
}

func TestClient_Create(t *testing.T) {
	c := initConf2()
	if err := c.Create("zjx11.toml", "test comment", "zjx", "mark"); err != nil {
		t.Errorf("client.Create() error(%v)", err)
		t.FailNow()
	}
}

func TestClient_Update(t *testing.T) {
	c := initConf2()
	if err := c.Update(21, "test comment11", "zjx", "mark"); err != nil {
		t.Errorf("client.Create() error(%v)", err)
		t.FailNow()
	}
}

func TestClient_ConfIng(t *testing.T) {
	c := initConf2()
	if val, err := c.ConfIng("zjx1.toml"); err != nil {
		t.Errorf("client.Create() error(%v)", err)
		t.FailNow()
	} else {
		t.Logf("%v", val)
	}
}

func initConf2() (c *Client) {
	conf.Addr = "172.16.33.134:9011"
	conf.Host = "testHost"
	conf.Path = "./"
	conf.AppID = "main.common-arch.msm-service"
	conf.Svr = "msm-service"
	conf.Ver = "server-1"
	conf.DeployEnv = "dev"
	conf.Zone = "sh001"
	conf.Token = "45338e440bdc11e880ce02420a0a0204"
	conf.TreeID = "2888"
	c = &Client{
		httpCli: &http.Client{Timeout: _httpTimeout},
		event:   make(chan string, 10),
	}
	return
}
