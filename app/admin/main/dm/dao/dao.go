package dao

import (
	"context"

	"go-common/app/admin/main/dm/conf"
	"go-common/app/admin/main/dm/model"
	"go-common/library/cache/memcache"
	"go-common/library/database/bfs"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

// Dao dao layer.
type Dao struct {
	actionPub *databus.Databus
	// http
	httpSearch *bm.Client
	httpCli    *bm.Client
	// mysql
	dmMetaWriter *sql.DB
	dmMetaReader *sql.DB
	biliDM       *sql.DB
	// memcache
	filterMC *memcache.Pool
	// subtitle mc
	subtitleMC *memcache.Pool
	// elastic client
	esCli *elastic.Elastic
	// bfs client
	bfsCli *bfs.BFS
	// http uri
	sendNotifyURI   string
	addMoralURI     string
	blockUserURI    string
	blockInfoAddURI string
	sendJudgeURI    string
	viewsURI        string
	typesURI        string
	seasonURI       string
	maskURI         string
	workFlowURI     string
	berserkerURI    string
}

// New new a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		actionPub: databus.New(c.ActionPub),
		// mysql
		dmMetaWriter: sql.NewMySQL(c.DB.DMMetaWriter),
		dmMetaReader: sql.NewMySQL(c.DB.DMMetaReader),
		biliDM:       sql.NewMySQL(c.DB.DM),
		// memcache
		filterMC:   memcache.NewPool(c.Memcache.Filter.Config),
		subtitleMC: memcache.NewPool(c.Memcache.Subtitle.Config),
		// elastic client
		esCli: elastic.NewElastic(c.Elastic),
		// http client
		sendNotifyURI:   c.Host.Message + _sendNotify,
		addMoralURI:     c.Host.Account + _addMoral,
		blockUserURI:    c.Host.API + _blockUser,
		blockInfoAddURI: c.Host.API + _blockInfoAdd,
		sendJudgeURI:    c.Host.API + _sendJudge,
		viewsURI:        c.Host.Videoup + _views,
		typesURI:        c.Host.Videoup + _types,
		seasonURI:       c.Host.Season + _season,
		maskURI:         c.Host.Mask + _mask,
		berserkerURI:    c.Host.Berserker,
		workFlowURI:     c.Host.API + _workFlowAppealDelete,
		httpCli:         bm.NewClient(c.HTTPClient.ClientConfig),
		httpSearch:      bm.NewClient(c.HTTPSearch.ClientConfig),
		bfsCli:          bfs.New(c.BFS),
	}
	return
}

// SendAction send action to job.
func (d *Dao) SendAction(c context.Context, k string, action *model.Action) (err error) {
	if err = d.actionPub.Send(c, k, action); err != nil {
		log.Error("actionPub.Send(action:%s,data:%s) error(%v)", action.Action, action.Data, err)
	} else {
		log.Info("actionPub.Send(action:%s,data:%s) success", action.Action, action.Data)
	}
	return
}

//BeginBiliDMTrans begin a transsaction of biliDM
func (d *Dao) BeginBiliDMTrans(c context.Context) (*sql.Tx, error) {
	return d.biliDM.Begin(c)
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.dmMetaWriter.Ping(c); err != nil {
		log.Error("d.dmMetaWriter error(%v)", err)
		return
	}
	if err = d.dmMetaReader.Ping(c); err != nil {
		log.Error("d.dmMetaReader error(%v)", err)
		return
	}
	if err = d.biliDM.Ping(c); err != nil {
		log.Error("d.biliDM error(%v)", err)
	}
	// mc
	filterMC := d.filterMC.Get(c)
	defer filterMC.Close()
	if err = filterMC.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("filterMC.Set error(%v)", err)
		return
	}
	return
}
