package dao

// var (
// 	familyDetail = "detail"
// 	columnTs     = "ts"
// )

// func rowKey(mid, ts int64, op string) string {
// 	return fmt.Sprintf("%d%d_%d_%s", mid%10, mid, math.MaxInt64-ts, op)
// }

// // AddLog coin modify log.
// func (d *Dao) AddLog(c context.Context, mid, ts int64, content map[string][]byte, table string) (err error) {
// 	var (
// 		mutate      *hrpc.Mutate
// 		operator    = string(content["operater"])
// 		key         = rowKey(mid, ts, operator)
// 		tsB         = make([]byte, 8)
// 		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.WriteTimeout))
// 	)
// 	defer cancel()
// 	binary.BigEndian.PutUint64(tsB, uint64(ts))
// 	content[columnTs] = tsB
// 	values := map[string]map[string][]byte{familyDetail: content}
// 	if mutate, err = hrpc.NewPutStr(ctx, table, key, values); err != nil {
// 		log.Error("hrpc.NewPutStr(%s, %s, %v) error(%v)", table, key, values, err)
// 		return
// 	}
// 	if _, err = d.hbase.Put(c, mutate); err != nil {
// 		log.Error("hbase.Put mid %d,error(%v)", mid, err)
// 		return
// 	}
// 	return
// }
