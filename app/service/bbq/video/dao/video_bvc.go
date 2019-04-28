package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/bbq/common/db/bbq"
	"go-common/app/service/bbq/video/model/grpc"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"strconv"
	"strings"
)

const (
	_querySvPlay = "select `svid`,`path`,`resolution_retio`,`code_rate`,`video_code`,`file_size`,`duration` from %s where svid in (%s) and is_deleted = 0 order by code_rate desc"
)
const (
	_defaultPlatform = "html5"
	_playBcNum       = 1
)

//RawSVBvcKey 批量获取playurl相对地址
func (d *Dao) RawSVBvcKey(c context.Context, svids []int64) (res map[int64][]*bbq.VideoBvc, err error) {
	var (
		tb   map[string][]string
		rows *sql.Rows
	)
	res = make(map[int64][]*bbq.VideoBvc)
	tb = make(map[string][]string)
	tName := "video_bvc_%02d"
	for _, v := range svids {
		if v <= 0 {
			continue
		}
		tbName := fmt.Sprintf(tName, v%100)
		tb[tbName] = append(tb[tbName], strconv.FormatInt(v, 10))
	}
	for k, v := range tb {
		query := fmt.Sprintf(_querySvPlay, k, strings.Join(v, ","))
		if rows, err = d.db.Query(c, query); err != nil {
			log.Errorv(c, log.KV("log", "RawSVBvcKey query sql"), log.KV("err", err))
			continue
		}
		for rows.Next() {
			tmp := bbq.VideoBvc{}
			if err = rows.Scan(&tmp.SVID, &tmp.Path, &tmp.ResolutionRetio, &tmp.CodeRate, &tmp.VideoCode, &tmp.FileSize, &tmp.Duration); err != nil {
				log.Errorv(c, log.KV("log", "RawSVBvcKey scan"), log.KV("err", err))
				continue
			}
			res[tmp.SVID] = append(res[tmp.SVID], &tmp)
		}
	}
	return
}

// RelPlayURLs 相对地址批量获取playurl
func (d *Dao) RelPlayURLs(c context.Context, addrs []string) (res map[string]*grpc.VideoKeyItem, err error) {
	res = make(map[string]*grpc.VideoKeyItem)
	req := &grpc.RequestMsg{
		Keys:     addrs,
		Backup:   uint32(_playBcNum),
		Platform: _defaultPlatform,
		UIP:      metadata.String(c, metadata.RemoteIP),
	}
	_str, _ := json.Marshal(req)
	log.V(5).Infov(c, log.KV("log", fmt.Sprintf("bvc play req (%s)", string(_str))))
	r, err := d.bvcPlayClient.ProtobufPlayurl(c, req)
	_str, _ = json.Marshal(r)
	if err != nil {
		log.Error("bvc play err[%v] ret[%s]", err, string(_str))
		return
	}
	log.V(5).Infov(c, log.KV("log", fmt.Sprintf("bvc play ret (%s)", string(_str))))
	res = r.Data
	return
}
