package cms

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"

	"go-common/app/interface/main/tv/conf"
	"go-common/library/database/sql"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d   *Dao
	ctx = context.TODO()
)

const (
	_pickSids  = 1
	_pickEpids = 2
	_pickAids  = 3
	_pickCids  = 4
)

func init() {
	// dir, _ := filepath.Abs("../../cmd/tv-interface.toml")
	// flag.Set("conf", dir)
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.web-svr.tv-interface")
		flag.Set("conf_token", "07c1826c1f39df02a1411cdd6f455879")
		flag.Set("tree_id", "15326")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}

func pickIDs(db *sql.DB, ttype int) (testSids []int64, err error) {
	var (
		querySQL string
		rows     *sql.Rows
	)
	switch ttype {
	case _pickSids:
		querySQL = "SELECT DISTINCT(id) FROM tv_ep_season WHERE is_deleted = 0 ORDER BY id DESC LIMIT 20"
	case _pickEpids:
		querySQL = "SELECT DISTINCT(epid) FROM tv_content WHERE is_deleted = 0 ORDER BY id DESC LIMIT 20"
	case _pickAids:
		querySQL = "SELECT DISTINCT(aid) FROM ugc_archive WHERE deleted = 0 ORDER BY id DESC LIMIT 20"
	case _pickCids:
		querySQL = "SELECT DISTINCT(cid) FROM ugc_video WHERE deleted = 0 ORDER BY id DESC LIMIT 20"
	}
	rows, err = db.Query(ctx, querySQL)
	if err != nil {
		fmt.Println("Query Err ", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			fmt.Printf("Scan Id %d, Err %v", id, err)
			return
		}
		testSids = append(testSids, id)
	}
	return
}

func TestDao_MixedFilter(t *testing.T) {
	Convey("TestDao_MixedFilter", t, WithDao(func(d *Dao) {
		sids, errDB := pickIDs(d.db, _pickSids)
		if errDB != nil {
			fmt.Println("PickSids Err ", errDB)
			return
		}
		aids, errDB2 := pickIDs(d.db, _pickAids)
		if errDB != nil {
			fmt.Println("PickAids Err ", errDB2)
			return
		}
		okSids, okAids := d.MixedFilter(ctx, sids, aids)
		fmt.Println(okSids)
		fmt.Println(okAids)
		So(len(okSids), ShouldBeLessThanOrEqualTo, len(sids))
		So(len(okAids), ShouldBeLessThanOrEqualTo, len(aids))
	}))
}

func TestDao_LoadVideosMeta(t *testing.T) {
	Convey("TestDao_MixedFilter", t, WithDao(func(d *Dao) {
		sids, errDB := pickIDs(d.db, _pickCids)
		if errDB != nil {
			fmt.Println("PickSids Err ", errDB)
			return
		}
		res, err := d.LoadVideosMeta(ctx, sids)
		So(err, ShouldBeNil)
		data, _ := json.Marshal(res)
		fmt.Println(string(data))
	}))
}
