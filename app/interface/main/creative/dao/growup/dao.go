package growup

import (
	httpx "go-common/library/net/http/blademaster"

	"go-common/app/interface/main/creative/conf"
)

const (
	//up check
	_upStatus = "/allowance/api/x/internal/growup/up/status"
	_upInfo   = "/allowance/api/x/internal/growup/up/info"
	_join     = "/allowance/api/x/internal/growup/up/add"
	_quit     = "/allowance/api/x/internal/growup/up/quit"
	//up income
	_summary = "/up-openapi/api/open_api/v1/income/summary"
	_stat    = "/up-openapi/api/open_api/v1/income/statis"
	_arc     = "/up-openapi/api/open_api/v1/income/archive"
	_breach  = "/up-openapi/api/open_api/v1/income/breach"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *httpx.Client
	// up check uri
	upStatusURL string
	upInfoURL   string
	joinURL     string
	quitURL     string
	// up income uri
	summaryURL string
	statURL    string
	arcURL     string
	breachURL  string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		client: httpx.NewClient(c.HTTPClient.UpMng),
		//up check
		upStatusURL: c.Host.Growup + _upStatus,
		upInfoURL:   c.Host.Growup + _upInfo,
		joinURL:     c.Host.Growup + _join,
		quitURL:     c.Host.Growup + _quit,
		//up check
		summaryURL: c.Host.UpMng + _summary,
		statURL:    c.Host.UpMng + _stat,
		arcURL:     c.Host.UpMng + _arc,
		breachURL:  c.Host.UpMng + _breach,
	}
	return
}
