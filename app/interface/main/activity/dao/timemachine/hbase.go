package timemachine

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/activity/model/timemachine"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

const (
	_hBaseUpItemTableName = "dw:mid_main_profile_y"
)

// reverse for string.
func hbaseRowKey(mid int64) string {
	s := strconv.FormatInt(mid, 10)
	rs := []rune(s)
	l := len(rs)
	for f, t := 0, l-1; f < t; f, t = f+1, t-1 {
		rs[f], rs[t] = rs[t], rs[f]
	}
	ns := string(rs)
	if l < 10 {
		for i := 0; i < 10-l; i++ {
			ns = ns + "0"
		}
	}
	return ns
}

// RawTimemachine .
func (d *Dao) RawTimemachine(c context.Context, mid int64) (data *timemachine.Item, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.Hbase.RegionReadTimeout))
		key         = hbaseRowKey(mid)
		tableName   = _hBaseUpItemTableName
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("UserItem d.hbase.GetStr tableName(%s)|mid(%d)|key(%v)|error(%v)", tableName, mid, key, err)
		return
	}
	if result == nil {
		return
	}
	data = &timemachine.Item{Mid: mid}
	for _, c := range result.Cells {
		if c == nil {
			continue
		}
		if !bytes.Equal(c.Family, []byte("m")) {
			continue
		}
		tmFillFields(data, c)
	}
	return
}

func tmFillFields(data *timemachine.Item, c *hrpc.Cell) {
	var (
		intVal int64
		strVal string
	)
	strVal = string(c.Value[:])
	if v, e := strconv.ParseInt(string(c.Value[:]), 10, 64); e == nil {
		intVal = v
	}
	switch {
	case bytes.Equal(c.Qualifier, []byte("is_up")):
		data.IsUp = intVal
	case bytes.Equal(c.Qualifier, []byte("dh")):
		data.DurationHour = intVal
	case bytes.Equal(c.Qualifier, []byte("adh")):
		data.AvDurationHour = intVal
	case bytes.Equal(c.Qualifier, []byte("a_vv")):
		data.ArchiveVv = intVal
	case bytes.Equal(c.Qualifier, []byte("tag_id")):
		data.LikeTagID = intVal
	case bytes.Equal(c.Qualifier, []byte("ls_vv")):
		data.LikeSubtidVv = intVal
	case bytes.Equal(c.Qualifier, []byte("ugc_avs")):
		data.LikesUgc3Avids = strVal
	case bytes.Equal(c.Qualifier, []byte("pgc_avs")):
		data.LikePgc3Avids = strVal
	case bytes.Equal(c.Qualifier, []byte("up")):
		data.LikeBestUpID = intVal
	case bytes.Equal(c.Qualifier, []byte("up_ad")):
		data.LikeUpAvDuration = intVal
	case bytes.Equal(c.Qualifier, []byte("up_ld")):
		data.LikeUpLiveDuration = intVal
	case bytes.Equal(c.Qualifier, []byte("up_avs")):
		data.LikeUp3Avs = strVal
	case bytes.Equal(c.Qualifier, []byte("up_st")):
		data.LikeLiveUpSubTname = strVal
	case bytes.Equal(c.Qualifier, []byte("cir_tm")):
		data.BrainwashCirTime = strVal
	case bytes.Equal(c.Qualifier, []byte("cir_av")):
		data.BrainwashCirAvid = intVal
	case bytes.Equal(c.Qualifier, []byte("cir_v")):
		data.BrainwashCirVv = intVal
	case bytes.Equal(c.Qualifier, []byte("fs_av")):
		data.FirstSubmitAvid = intVal
	case bytes.Equal(c.Qualifier, []byte("fs_tm")):
		data.FirstSubmitTime = strVal
	case bytes.Equal(c.Qualifier, []byte("fs_ty")):
		data.FirstSubmitType = intVal
	case bytes.Equal(c.Qualifier, []byte("s_av_rd")):
		data.SubmitAvsRds = strVal
	case bytes.Equal(c.Qualifier, []byte("bt_av")):
		data.BestAvid = intVal
	case bytes.Equal(c.Qualifier, []byte("bt_ty")):
		data.BestAvidType = intVal
	case bytes.Equal(c.Qualifier, []byte("bt_av_o")):
		data.BestAvidOld = intVal
	case bytes.Equal(c.Qualifier, []byte("bt_av_ty")):
		data.BestAvidOldType = intVal
	case bytes.Equal(c.Qualifier, []byte("o_vv")):
		data.OldAvVv = intVal
	case bytes.Equal(c.Qualifier, []byte("all_vv")):
		data.AllVv = intVal
	case bytes.Equal(c.Qualifier, []byte("live_d")):
		data.UpLiveDuration = intVal
	case bytes.Equal(c.Qualifier, []byte("is_live")):
		data.IsLiveUp = intVal
	case bytes.Equal(c.Qualifier, []byte("ld")):
		data.ValidLiveDays = intVal
	case bytes.Equal(c.Qualifier, []byte("md")):
		data.MaxCdnNumDate = strVal
	case bytes.Equal(c.Qualifier, []byte("mc")):
		data.MaxCdnNum = intVal
	case bytes.Equal(c.Qualifier, []byte("att")):
		data.Attentions = intVal
	case bytes.Equal(c.Qualifier, []byte("fan_vv")):
		data.UpBestFanVv = strVal
	case bytes.Equal(c.Qualifier, []byte("fan_live")):
		data.UpBestFanLiveMinute = strVal
	case bytes.Equal(c.Qualifier, []byte("like_tid")):
		data.Like2Tids = strVal
	case bytes.Equal(c.Qualifier, []byte("like_st")):
		data.Like2SubTids = strVal
	case bytes.Equal(c.Qualifier, []byte("wr")):
		data.WinRatio = strVal
	case bytes.Equal(c.Qualifier, []byte("pd_hr")):
		data.PlayDurationHourRep = intVal
	case bytes.Equal(c.Qualifier, []byte("lu_adr")):
		data.LikeUpAvDurationRep = intVal
	}
}
