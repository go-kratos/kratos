package order

import (
	"go-common/app/interface/main/creative/conf"
	bm "go-common/library/net/http/blademaster"
)

const (
	// --- from chaodian v2
	_executeOrders = "/api/open_api/v2/execute_orders"
	_ups           = "/api/open_api/v2/ups"
	_getOrderByAid = "/api/open_api/v2/execute_orders/by_av_id"
	_archiveStatus = "/api/open_api/v2/execute_orders/video/status"
	_oasis         = "/api/open_api/v2/ups/up_execute_order_statistics" //绿洲计划
	_launchtime    = "/api/open_api/v2/execute_orders/launch_time"
	// ----
	_upValidate   = "/meet/api/openApi/v1/up/validate"
	_accountState = "/allowance/api/x/admin/growup/up/account/state"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client   *bm.Client
	chaodian *bm.Client
	// uri
	executeOrdersURI string
	upsURI           string
	getOrderByAidURI string
	archiveStatusURI string
	oasisURI         string
	launchTimeURI    string
	upValidateURI    string
	accountStateURI  string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                c,
		client:           bm.NewClient(c.HTTPClient.UpMng),
		chaodian:         bm.NewClient(c.HTTPClient.Chaodian),
		executeOrdersURI: c.Host.Chaodian + _executeOrders,
		upsURI:           c.Host.Chaodian + _ups,
		getOrderByAidURI: c.Host.Chaodian + _getOrderByAid,
		archiveStatusURI: c.Host.Chaodian + _archiveStatus,
		oasisURI:         c.Host.Chaodian + _oasis,
		launchTimeURI:    c.Host.Chaodian + _launchtime,
		upValidateURI:    c.Host.UpMng + _upValidate,
		accountStateURI:  c.Host.Profit + _accountState,
	}
	return
}
