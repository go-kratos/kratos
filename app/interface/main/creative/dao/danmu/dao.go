package danmu

import (
	"go-common/app/interface/main/creative/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao is creative dao.
type Dao struct {
	// config
	c *conf.Config
	// http client
	client *bm.Client
	// assist url
	assistDmBannedURL string

	advDmPurchaseListURL   string
	advDmPurchasePassURL   string
	advDmPurchaseDenyURL   string
	advDmPurchaseCancelURL string

	dmSearchURL                string
	dmEditURL                  string
	dmRecentURL                string
	dmTransferURL              string
	dmPoolURL                  string
	dmDistriURL                string
	dmProtectApplyListURL      string
	dmProtectApplyStatusURL    string
	dmProtectApplyVideoListURL string
	dmReportUpListURL          string
	dmReportUpArchivesURL      string
	dmReportUpEditURL          string
}

// New init api url
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                 c,
		client:            bm.NewClient(c.HTTPClient.Slow),
		assistDmBannedURL: c.Host.API + _setDmBannedURI,

		advDmPurchaseListURL:       c.Host.API + _getDmPurchaseListURI,
		advDmPurchasePassURL:       c.Host.API + _setDmPurchasePassURI,
		advDmPurchaseDenyURL:       c.Host.API + _setDmPurchaseDenyURI,
		advDmPurchaseCancelURL:     c.Host.API + _setDmPurchaseCancelURI,
		dmSearchURL:                c.Host.API + _dmSearchURI,
		dmEditURL:                  c.Host.API + _dmEditURI,
		dmRecentURL:                c.Host.API + _dmRecentURI,
		dmTransferURL:              c.Host.API + _dmTransferURI,
		dmPoolURL:                  c.Host.API + _dmPoolURI,
		dmDistriURL:                c.Host.API + _dmDistriURI,
		dmProtectApplyStatusURL:    c.Host.API + _dmProtectApplyStatusURI,
		dmProtectApplyListURL:      c.Host.API + _dmProtectApplyListURI,
		dmProtectApplyVideoListURL: c.Host.API + _dmProtectApplyVideoListURI,
		dmReportUpEditURL:          c.Host.API + _dmReportUpEditURI,
		dmReportUpListURL:          c.Host.API + _dmReportUpListURI,
		dmReportUpArchivesURL:      c.Host.API + _dmReportUpArchivesURI,
	}
	return
}
