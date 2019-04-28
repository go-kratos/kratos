package assist

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
	assistLogsURL      string
	assistListURL      string
	assistInfoURL      string
	assistLogInfoURL   string
	assistAddURL       string
	assistDelURL       string
	assistLogAddURL    string
	assistLogRevocURL  string
	assistStatURL      string
	assistLogObjURL    string
	liveStatusURL      string
	liveAddAssistURL   string
	liveDelAssistURL   string
	liveRevocBannedURL string
	liveAssistsURL     string
	liveCheckAssURL    string
}

// New init api url
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                  c,
		client:             bm.NewClient(c.HTTPClient.Normal),
		assistAddURL:       c.Host.API + _addAssistURI,
		assistDelURL:       c.Host.API + _delAssistURI,
		assistInfoURL:      c.Host.API + _getAssistInfoURI,
		assistLogInfoURL:   c.Host.API + _getAssistLogInfoURI,
		assistLogsURL:      c.Host.API + _getAssistLogsURI,
		assistLogAddURL:    c.Host.API + _addAssistLogURI,
		assistListURL:      c.Host.API + _getAssistURI,
		assistLogRevocURL:  c.Host.API + _revocAssistLogURI,
		assistStatURL:      c.Host.API + _getAssistStatURI,
		assistLogObjURL:    c.Host.API + _getAssistLogObjURI,
		liveStatusURL:      c.Host.Live + _liveStatus,
		liveAddAssistURL:   c.Host.Live + _liveAddAssist,
		liveDelAssistURL:   c.Host.Live + _liveDelAssist,
		liveRevocBannedURL: c.Host.Live + _liveRevocBanned,
		liveAssistsURL:     c.Host.Live + _liveAssists,
		liveCheckAssURL:    c.Host.Live + _liveCheckAssist,
	}
	return
}
