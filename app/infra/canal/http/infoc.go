package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/http"

	"go-common/app/infra/canal/conf"
	"go-common/app/infra/canal/infoc"
	config "go-common/library/conf"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"github.com/BurntSushi/toml"
	"github.com/siddontang/go-mysql/canal"
)

const (
	_heartHeat   = 60
	_readTimeout = 90
	_flavor      = "mysql"
	_updateUser  = "canal"
	_updateMark  = "infoc"
)

// InfocConf .
type infocConf struct {
	Addr     string     `json:"db_addr"`
	User     string     `json:"user"`
	Pass     string     `json:"pass"`
	InfocDBs []*infocDB `json:"databases"`
}

// InfocDB .
type infocDB struct {
	Schema           string       `json:"schema"`
	Tables           []*infoTable `json:"tables"`
	LancerAddr       string       `json:"lancer_addr"`
	LancerTaskID     string       `json:"lancer_task_id"`
	LancerReportAddr string       `json:"lancer_report_addr"`
	Proto            string       `json:"proto"`
}

// InfoTable .
type infoTable struct {
	Name       string   `json:"name"`
	OmitFlied  []string `json:"omit_field"`
	OmitAction []string `json:"omit_action"`
}

func infocPost(c *bm.Context) {
	var (
		ics []*infocConf
		bs  []byte
		err error
		buf *bytes.Buffer
	)
	content := make(map[string]string)
	if bs, err = ioutil.ReadAll(c.Request.Body); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(bs, &ics); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, ifc := range ics {
		databases := make([]*conf.Database, len(ifc.InfocDBs))
		for idx, infocDB := range ifc.InfocDBs {
			tables := make([]*conf.CTable, len(infocDB.Tables))
			for ix, table := range infocDB.Tables {
				tables[ix] = &conf.CTable{
					Name:       table.Name,
					OmitAction: table.OmitAction,
					OmitField:  table.OmitFlied,
				}
			}
			databases[idx] = &conf.Database{
				Schema: infocDB.Schema,
				Infoc: &infoc.Config{
					TaskID:       infocDB.LancerTaskID,
					Addr:         infocDB.LancerAddr,
					ReporterAddr: infocDB.LancerReportAddr,
					Proto:        infocDB.Proto,
				},
				CTables: tables,
			}
		}
		ic := &conf.InsConf{
			Databases: databases,
			Config: &canal.Config{
				Addr:            ifc.Addr,
				User:            ifc.User,
				Password:        ifc.Pass,
				ServerID:        crc32.ChecksumIEEE([]byte(ifc.Addr)),
				Flavor:          _flavor,
				HeartbeatPeriod: _heartHeat,
				ReadTimeout:     _readTimeout,
			},
		}
		var isc = &struct {
			InsConf *conf.InsConf `toml:"instance"`
		}{
			InsConf: ic,
		}
		buf = new(bytes.Buffer)
		if err = toml.NewEncoder(buf).Encode(isc); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		content[fmt.Sprintf("%v.toml", ifc.Addr)] = buf.String()
	}
	for cn, cv := range content {
		value, err := conf.ConfClient.ConfIng(cn)
		if err == nil {
			err = conf.ConfClient.Update(value.CID, cv, _updateUser, _updateMark)
		} else if err == ecode.NothingFound {
			err = conf.ConfClient.Create(cn, cv, _updateUser, _updateMark)
		}
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

func infocCurrent(c *bm.Context) {
	var (
		ok     bool
		result []*config.Value
	)
	if result, ok = conf.ConfClient.Configs(); !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
	ics := make([]*infocConf, 0, len(result))
	for _, ns := range result {
		var ic struct {
			InsConf *conf.InsConf `toml:"instance"`
		}
		if _, err := toml.Decode(ns.Config, &ic); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if ic.InsConf == nil {
			continue
		}
		icf := &infocConf{
			Addr: ic.InsConf.Addr,
			User: ic.InsConf.User,
			Pass: ic.InsConf.Password,
		}
		for _, icdb := range ic.InsConf.Databases {
			if icdb.Infoc == nil {
				continue
			}
			tables := make([]*infoTable, len(icdb.CTables))
			for idx, ctable := range icdb.CTables {
				tables[idx] = &infoTable{
					Name:       ctable.Name,
					OmitFlied:  ctable.OmitField,
					OmitAction: ctable.OmitAction,
				}
			}
			icf.InfocDBs = append(icf.InfocDBs, &infocDB{
				Schema:       icdb.Schema,
				Tables:       tables,
				LancerAddr:   icdb.Infoc.Addr,
				LancerTaskID: icdb.Infoc.TaskID,
			})
		}
		ics = append(ics, icf)
	}
	c.JSON(ics, nil)
}
