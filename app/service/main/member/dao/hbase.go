package dao

var (
// tableMoralLog = "ugc:MoralLog"
// familyDetail = "detail"
// familyDetailB = []byte(familyDetail)
// columnIP      = "ip"
// columnTs = "ts"
// columnMid     = "mid"
// columnLogID   = "log_id"
// _logDuration = int64(7 * 24 * 3600)
)

// func rowKey(mid, ts int64) string {
// 	return fmt.Sprintf("%d%d_%d", mid%10, mid, math.MaxInt64-ts)
// }

// func moralRowKey(mid, ts int64, tid uint64) string {
// 	return fmt.Sprintf("%d%d_%d_%d", mid%10, mid, math.MaxInt64-ts, tid)
// }

// AddMoralLog add moral modify log.
// func (d *Dao) AddMoralLog(c context.Context, mid, ts int64, content map[string][]byte) (err error) {
// 	var (
// 		mutate      *hrpc.Mutate
// 		key         = moralRowKey(mid, ts, genID())
// 		tsB         = make([]byte, 8)
// 		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.WriteTimeout))
// 	)
// 	defer cancel()
// 	binary.BigEndian.PutUint64(tsB, uint64(ts))
// 	content[columnTs] = tsB
// 	values := map[string]map[string][]byte{familyDetail: content}
// 	if mutate, err = hrpc.NewPutStr(ctx, tableMoralLog, key, values); err != nil {
// 		log.Error("hrpc.NewPutStr(%s, %s, %v) error(%v)", tableMoralLog, key, values, err)
// 		return
// 	}
// 	if _, err = d.hbase.Put(c, mutate); err != nil {
// 		log.Error("hbase.Put mid %d,error(%v) value(%v)", mid, err, values)
// 		return
// 	}
// 	log.Info("hbase.Put moral log success mid: %d,value(%v)", mid, values)
// 	return
// }

// func genID() uint64 {
// 	i := [16]byte(uuid.NewV1())
// 	return farm.Hash64(i[:]) % math.MaxInt64
// }
