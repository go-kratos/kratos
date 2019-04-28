package dao

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"

	"github.com/gogo/protobuf/types"
)

//订单相关常量
const (
	cacheTimeout        = 10
	DefaultOrderPSize   = 20        //默认每页数量
	DefaultOrderOrderBy = "id desc" //默认排序
	//sqlGetUserOrders 查询用户订单列表
	sqlGetUserOrders = "SELECT %s FROM order_main WHERE uid=? AND is_deleted = 0 ORDER BY %s LIMIT %d,%d"
	//sqlGetUserItemOrders 查询用户已购买的某个项目订单
	sqlGetUserItemOrders = "SELECT %s FROM order_main WHERE uid=? AND item_id=? AND status IN (%s)"
	//sqlCountUserOrders 查询用户订单数量
	sqlCountUserOrders = "SELECT COUNT(*) FROM order_main WHERE uid=? AND is_deleted=0"
	//sqlGetOrders 查询订单列表
	sqlGetOrders = "SELECT %s FROM order_main WHERE order_id IN (%s)"
	//sqlInsertOrderMains 批量写入order_main
	sqlInsertOrderMains = "INSERT INTO order_main(%s) VALUES %s"
	//sqlGetOrderDetails 查询订单详情
	sqlGetOrderDetails = "SELECT %s FROM order_detail WHERE order_id IN (%s)"
	//sqlInsertOrderDetails 批量写入order_detail
	sqlInsertOrderDetails = "INSERT INTO order_detail(%s) VALUES %s"
	//sqlGetOrderSkus 查询order_sku
	sqlGetOrderSkus = "SELECT %s FROM order_sku WHERE order_id IN (%s)"
	//sqlInsertOrderSkus 批量写入order_sku
	sqlInsertOrderSkus = "INSERT INTO order_sku(%s) VALUES %s"
	//sqlGetBoughtSkus 获取已购买的sku
	sqlGetBoughtSkus = "SELECT `count` FROM order_sku WHERE order_id IN (%s) AND sku_id IN (%s)"
	//sqlGetOrderPayChs 获取支付流水
	sqlGetOrderPayChs = "SELECT %s FROM order_pay_charge WHERE order_id IN (%s) AND paid=1"

	//获取结算对账订单(0元单)
	sqlGetSettleCompareOrders = "SELECT id,order_id FROM order_main WHERE ctime>=? and ctime<? AND id>? AND status IN (%s) AND pay_money=0 ORDER BY ctime,id LIMIT ?"
	//获取结算对帐退款单
	sqlGetSettleCompareRefunds = "SELECT id,order_id,refund_apply_time FROM order_refund WHERE ctime>=? AND ctime<? AND id>? AND status IN (%s) AND refund_money=0 ORDER BY ctime,id LIMIT ?"
)

//RawOrders 从db查询用户订单信息，包括以下情况
//* 按照订单号查询
//* 按照uid查分页
//* 按照uid+商品id+状态查询列表
func (d *Dao) RawOrders(ctx context.Context, req *model.OrderMainQuerier) (orders []*model.OrderMain, err error) {
	defer LogX(ctx, req, orders, err)
	o := new(model.OrderMain)
	sqls, args := pickOrderSQL(req, o.GetFields(nil))
	q := sqls[0]
	if q == "" {
		return
	}
	r, err := d.db.Query(ctx, q, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	defer r.Close()
	for r.Next() {
		vptrs := make(map[int]interface{})
		ptrs := o.GetPtrs(nil, vptrs)
		if err = r.Scan(ptrs...); err != nil {
			return
		}
		for k, v := range vptrs {
			json.Unmarshal([]byte(*ptrs[k].(*string)), v)
		}
		orders = append(orders, o)
		o = new(model.OrderMain)
	}
	return
}

//RawOrderCount 从db获取用户订单数目，当查询条件仅有uid时生效，否则返回0
func (d *Dao) RawOrderCount(ctx context.Context, req *model.OrderMainQuerier) (cnt int64, err error) {
	defer LogX(ctx, req, cnt, err)
	sqls, args := pickOrderSQL(req, nil)
	q := sqls[1]
	if q == "" {
		return
	}
	r := d.db.QueryRow(ctx, q, args...)
	err = r.Scan(&cnt)
	//缓存回源逻辑0会被当成没命中缓存
	if cnt == 0 {
		cnt = -1
	}
	return
}

//RawOrderDetails 从db获取订单详细
func (d *Dao) RawOrderDetails(ctx context.Context, oids []int64) (orders map[int64]*model.OrderDetail, err error) {
	defer LogX(ctx, oids, orders, err)
	lo := len(oids)
	if lo == 0 {
		return
	}
	a := make([]interface{}, lo)
	for k, v := range oids {
		a[k] = v
	}
	o := new(model.OrderDetail)
	f := o.GetFields(nil)
	q := fmt.Sprintf(sqlGetOrderDetails, "`"+strings.Join(f, "`,`")+"`", strings.Repeat(",?", lo)[1:])
	r, err := d.db.Query(ctx, q, a...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	defer r.Close()
	orders = make(map[int64]*model.OrderDetail, lo)
	for r.Next() {
		vptrs := make(map[int]interface{})
		ptrs := o.GetPtrs(nil, vptrs)
		if err = r.Scan(ptrs...); err != nil {
			return
		}
		for k, v := range vptrs {
			json.Unmarshal([]byte(*ptrs[k].(*string)), v)
		}
		o.Decrypt(d.c.Encrypt)
		orders[o.OrderID] = o
		o = new(model.OrderDetail)
	}
	return
}

//RawOrderSKUs 从db获取订单的sku
func (d *Dao) RawOrderSKUs(ctx context.Context, oids []int64) (skus map[int64][]*model.OrderSKU, err error) {
	defer LogX(ctx, oids, skus, err)
	lo := len(oids)
	if lo == 0 {
		return
	}
	a := make([]interface{}, lo)
	for k, v := range oids {
		a[k] = v
	}
	o := new(model.OrderSKU)
	f := o.GetFields(nil)
	q := fmt.Sprintf(sqlGetOrderSkus, "`"+strings.Join(f, "`,`")+"`", strings.Repeat(",?", lo)[1:])
	r, err := d.db.Query(ctx, q, a...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	defer r.Close()
	skus = make(map[int64][]*model.OrderSKU, lo)
	for r.Next() {
		vptrs := make(map[int]interface{})
		ptrs := o.GetPtrs(nil, vptrs)
		if err = r.Scan(ptrs...); err != nil {
			return
		}
		for k, v := range vptrs {
			json.Unmarshal([]byte(*ptrs[k].(*string)), v)
		}
		skus[o.OrderID] = append(skus[o.OrderID], o)
		o = new(model.OrderSKU)
	}
	return
}

//RawBoughtCount 按照票种获取用户已购的订单票数
func (d *Dao) RawBoughtCount(ctx context.Context, uid string, itemID int64, skuIDs []int64) (cnt int64, err error) {
	query := &model.OrderMainQuerier{
		UID:    uid,
		ItemID: itemID,
		Status: []int16{consts.OrderStatusUnpaid, consts.OrderStatusPaid},
	}
	orders, err := d.RawOrders(ctx, query)
	if err != nil {
		return
	}
	//存在部分退款，要减去已退款票数
	pt := false
	lo := len(orders)
	oids := make([]int64, lo)
	for k, v := range orders {
		if v.RefundStatus == consts.RefundStatusPtRefunded {
			pt = true
		}
		cnt += v.Count
		oids[k] = v.OrderID
	}
	ls := len(skuIDs)
	//查具体sku的，要再次从order_sku表统计
	if ls > 0 && lo > 0 {
		cnt, err = d.rawBoughtSkusCnt(ctx, oids, skuIDs)
		if err != nil {
			return
		}
	}
	//部分退款的减去已退款张数
	if pt {
		var rCnt int64
		rCnt, err = d.rawRefundTicketCnt(ctx, oids)
		if err != nil {
			return
		}
		cnt -= rCnt
	}
	return
}

//rawBoughtSkusCnt 按照skuID统计用户购买数
func (d *Dao) rawBoughtSkusCnt(ctx context.Context, oids []int64, skuIDs []int64) (cnt int64, err error) {
	lo := len(oids)
	ls := len(skuIDs)
	if lo == 0 || ls == 0 {
		err = ecode.RequestErr
		return
	}
	q := fmt.Sprintf(sqlGetBoughtSkus, strings.Repeat(",?", lo)[1:], strings.Repeat(",?", ls)[1:])
	a := make([]interface{}, lo+ls)
	for k, v := range oids {
		a[k] = v
	}
	for k, v := range skuIDs {
		a[k+lo] = v
	}
	r, err := d.db.Query(ctx, q, a...)
	if err != nil {
		return
	}
	defer r.Close()
	var c int64
	for r.Next() {
		if err = r.Scan(&c); err != nil {
			return
		}
		cnt += c
	}
	return
}

//RawOrderPayCharges 从db获取订单流水
func (d *Dao) RawOrderPayCharges(ctx context.Context, oids []int64) (chs map[int64]*model.OrderPayCharge, err error) {
	defer LogX(ctx, oids, chs, err)
	lo := len(oids)
	if lo == 0 {
		return
	}
	a := make([]interface{}, lo)
	for k, v := range oids {
		a[k] = v
	}
	o := new(model.OrderPayCharge)
	f := o.GetFields(nil)
	q := fmt.Sprintf(sqlGetOrderPayChs, "`"+strings.Join(f, "`,`")+"`", strings.Repeat(",?", lo)[1:])
	r, err := d.db.Query(ctx, q, a...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	defer r.Close()
	chs = make(map[int64]*model.OrderPayCharge, lo)
	for r.Next() {
		ptrs := o.GetPtrs(&types.FieldMask{Paths: f}, nil)
		if err = r.Scan(ptrs...); err != nil {
			return
		}
		chs[o.OrderID] = o
		o = new(model.OrderPayCharge)
	}
	return
}

//CacheOrders 从缓存获取订单基础信息
func (d *Dao) CacheOrders(ctx context.Context, req *model.OrderMainQuerier) (orders []*model.OrderMain, err error) {
	pool := d.redis.Get(ctx)
	defer func() {
		pool.Close()
		LogX(ctx, req, orders, err)
	}()
	var keys []interface{}
	if l := len(req.OrderID); l > 0 {
		keys = make([]interface{}, l)
		for k, v := range req.OrderID {
			keys[k] = fmt.Sprintf("%s:%d", model.CacheKeyOrderMn, v)
		}
	} else if key := oidCacheKey(req); key != "" {
		//如果查的是列表，先查出列表缓存的orderId，再根据orderId查订单信息
		var b []byte
		key = model.CacheKeyOrderList + ":" + key
		if b, err = redis.Bytes(pool.Do("GET", key)); err != nil {
			if err == redis.ErrNil {
				err = nil
			}
			return
		}
		s := string(b)
		LogX(ctx, []string{"GET", key}, s, err)
		oids := strings.Split(s, ",")
		keys = make([]interface{}, len(oids))
		for k, v := range oids {
			keys[k] = model.CacheKeyOrderMn + ":" + v
		}
	}
	if len(keys) == 0 {
		return
	}
	var data [][]byte
	if data, err = redis.ByteSlices(pool.Do("MGET", keys...)); err != nil {
		return
	}
	LogX(ctx, append([]interface{}{"MGET"}, keys...), data, err)
	for _, v := range data {
		if v != nil {
			o := &model.OrderMain{}
			orders = append(orders, o)
			json.Unmarshal(v, o)
		}
	}
	return
}

//CacheOrderCount 从缓存获取订单数目
func (d *Dao) CacheOrderCount(ctx context.Context, req *model.OrderMainQuerier) (cnt int64, err error) {
	pool := d.redis.Get(ctx)
	defer func() {
		pool.Close()
		LogX(ctx, req, cnt, err)
	}()
	if key := oidCacheKey(req); key != "" {
		if cnt, err = redis.Int64(pool.Do("GET", model.CacheKeyOrderCnt+":"+key)); err == redis.ErrNil {
			err = nil
		}
	}
	return
}

//CacheOrderDetails 从缓存获取订单详细
func (d *Dao) CacheOrderDetails(ctx context.Context, oids []int64) (orders map[int64]*model.OrderDetail, err error) {
	pool := d.redis.Get(ctx)
	defer func() {
		pool.Close()
		LogX(ctx, oids, orders, err)
	}()
	lo := len(oids)
	keys := make([]interface{}, lo)
	for k, v := range oids {
		keys[k] = fmt.Sprintf("%s:%d", model.CacheKeyOrderDt, v)
	}
	var data [][]byte
	if data, err = redis.ByteSlices(pool.Do("MGET", keys...)); err != nil {
		return
	}
	LogX(ctx, append([]interface{}{"MGET"}, keys...), data, nil)
	orders = make(map[int64]*model.OrderDetail, lo)
	for k, v := range oids {
		if data[k] == nil {
			continue
		}
		o := &model.OrderDetail{}
		if err = json.Unmarshal(data[k], o); err != nil {
			err = nil
		} else {
			orders[v] = o
		}
	}
	return
}

//CacheOrderSKUs 获取订单sku缓存
func (d *Dao) CacheOrderSKUs(ctx context.Context, oids []int64) (skus map[int64][]*model.OrderSKU, err error) {
	pool := d.redis.Get(ctx)
	defer func() {
		pool.Close()
		LogX(ctx, oids, skus, err)
	}()
	lo := len(oids)
	keys := make([]interface{}, lo)
	for k, v := range oids {
		keys[k] = fmt.Sprintf("%s:%d", model.CacheKeyOrderSKU, v)
	}
	var data [][]byte
	if data, err = redis.ByteSlices(pool.Do("MGET", keys...)); err != nil {
		return
	}
	LogX(ctx, append([]interface{}{"MGET"}, keys...), data, nil)
	skus = make(map[int64][]*model.OrderSKU, lo)
	for k, v := range oids {
		if data[k] == nil {
			continue
		}
		o := []*model.OrderSKU{}
		if err = json.Unmarshal(data[k], &o); err != nil {
			err = nil
		} else {
			skus[v] = o
		}
	}
	return
}

//CacheOrderPayCharges 从缓存获取订单流水
func (d *Dao) CacheOrderPayCharges(ctx context.Context, oids []int64) (chs map[int64]*model.OrderPayCharge, err error) {
	pool := d.redis.Get(ctx)
	defer func() {
		pool.Close()
		LogX(ctx, oids, chs, err)
	}()
	lo := len(oids)
	if lo == 0 {
		return
	}
	keys := make([]interface{}, lo)
	for k, v := range oids {
		keys[k] = fmt.Sprintf("%s:%d", model.CacheKeyOrderPayCh, v)
	}
	var data [][]byte
	if data, err = redis.ByteSlices(pool.Do("MGET", keys...)); err != nil {
		return
	}
	LogX(ctx, append([]interface{}{"MGET"}, keys...), data, nil)
	chs = make(map[int64]*model.OrderPayCharge, lo)
	for k, v := range oids {
		if data[k] == nil {
			continue
		}
		ch := &model.OrderPayCharge{}
		if err = json.Unmarshal(data[k], ch); err != nil {
			err = nil
		} else {
			chs[v] = ch
		}
	}
	return
}

//AddCacheOrders 设置订单基础信息缓存
func (d *Dao) AddCacheOrders(ctx context.Context, req *model.OrderMainQuerier, res []*model.OrderMain) (err error) {
	pool := d.redis.Get(ctx)
	defer func() {
		pool.Flush()
		pool.Close()
		LogX(ctx, []interface{}{req, res}, nil, err)
	}()
	data := make([]interface{}, len(res)*2)
	var sOids string
	for k, v := range res {
		if sOids == "" {
			sOids = fmt.Sprintf("%d", v.OrderID)
		} else {
			sOids += fmt.Sprintf(",%d", v.OrderID)
		}
		b, _ := json.Marshal(v)
		data[k*2] = fmt.Sprintf("%s:%d", model.CacheKeyOrderMn, v.OrderID)
		data[k*2+1] = b
	}
	//设置列表orderID缓存
	if key := oidCacheKey(req); key != "" {
		arg := []interface{}{model.CacheKeyOrderList + ":" + key, cacheTimeout, sOids}
		LogX(ctx, append([]interface{}{"SETEX"}, arg...), nil, nil)
		if err = pool.Send("SETEX", arg...); err != nil {
			return
		}
	}
	LogX(ctx, append([]interface{}{"MSET"}, data...), nil, nil)
	if err = pool.Send("MSET", data...); err != nil {
		return
	}
	for i := 0; i < len(data); i += 2 {
		pool.Send("EXPIRE", data[i], cacheTimeout)
	}
	return
}

//AddCacheOrderCount 设置订单数目缓存
func (d *Dao) AddCacheOrderCount(ctx context.Context, req *model.OrderMainQuerier, res int64) (err error) {
	pool := d.redis.Get(ctx)
	defer func() {
		pool.Close()
		LogX(ctx, []interface{}{req, res}, nil, err)
	}()
	if key := oidCacheKey(req); key != "" {
		_, err = pool.Do("SETEX", model.CacheKeyOrderCnt+":"+key, cacheTimeout, res)
	}
	return
}

//AddCacheOrderDetails 增加订单详细缓存
func (d *Dao) AddCacheOrderDetails(ctx context.Context, orders map[int64]*model.OrderDetail) (err error) {
	pool := d.redis.Get(ctx)
	defer func() {
		pool.Flush()
		pool.Close()
		LogX(ctx, orders, nil, err)
	}()
	data := make([]interface{}, len(orders)*2)
	i := 0
	for k, v := range orders {
		key := fmt.Sprintf("%s:%d", model.CacheKeyOrderDt, k)
		var b []byte
		b, _ = json.Marshal(v)
		data[i] = key
		data[i+1] = b
		i += 2
	}
	LogX(ctx, append([]interface{}{"MSET"}, data...), nil, nil)
	if err = pool.Send("MSET", data...); err != nil {
		return
	}
	for i := 0; i < len(data); i += 2 {
		pool.Send("EXPIRE", data[i], cacheTimeout)
	}
	return
}

//AddCacheOrderSKUs 增加订单sku缓存
func (d *Dao) AddCacheOrderSKUs(ctx context.Context, skus map[int64][]*model.OrderSKU) (err error) {
	pool := d.redis.Get(ctx)
	defer func() {
		pool.Flush()
		pool.Close()
		LogX(ctx, skus, nil, err)
	}()
	data := make([]interface{}, len(skus)*2)
	i := 0
	for k, v := range skus {
		key := fmt.Sprintf("%s:%d", model.CacheKeyOrderSKU, k)
		var b []byte
		b, err = json.Marshal(v)
		data[i] = key
		data[i+1] = b
		i += 2
	}
	LogX(ctx, append([]interface{}{"MSET"}, data...), nil, nil)
	if err = pool.Send("MSET", data...); err != nil {
		return
	}
	for i := 0; i < len(data); i += 2 {
		pool.Send("EXPIRE", data[i], cacheTimeout)
	}
	return
}

//AddCacheOrderPayCharges 增加订单流水缓存
func (d *Dao) AddCacheOrderPayCharges(ctx context.Context, chs map[int64]*model.OrderPayCharge) (err error) {
	pool := d.redis.Get(ctx)
	defer func() {
		pool.Flush()
		pool.Close()
		LogX(ctx, chs, nil, err)
	}()
	data := make([]interface{}, len(chs)*2)
	i := 0
	for k, v := range chs {
		key := fmt.Sprintf("%s:%d", model.CacheKeyOrderPayCh, k)
		var b []byte
		b, _ = json.Marshal(v)
		data[i] = key
		data[i+1] = b
		i += 2
	}
	LogX(ctx, append([]interface{}{"MSET"}, data...), nil, nil)
	if err = pool.Send("MSET", data...); err != nil {
		return
	}
	for i := 0; i < len(data); i += 2 {
		pool.Send("EXPIRE", data[i], cacheTimeout)
	}
	return
}

//DelCacheOrders 删除订单相关缓存，如果是新增或删除订单，要删除这个uid的订单列表缓存
func (d *Dao) DelCacheOrders(ctx context.Context, req *model.OrderMainQuerier) {
	pool := d.redis.Get(ctx)
	var ret interface{}
	defer func() {
		pool.Close()
		LogX(ctx, req, ret, nil)
	}()
	for _, v := range req.OrderID {
		ret, _ = pool.Do("DEL",
			fmt.Sprintf("%s:%d", model.CacheKeyOrderMn, v),
			fmt.Sprintf("%s:%d", model.CacheKeyOrderDt, v),
			fmt.Sprintf("%s:%d", model.CacheKeyOrderSKU, v),
		)
	}
	if key := oidCacheKey(req); key != "" {
		ret, _ = pool.Do("DEL", model.CacheKeyOrderList+":"+key, model.CacheKeyOrderCnt+":"+key)
	}
}

//TxInsertOrders 插入订单表，返回成功行数
func (d *Dao) TxInsertOrders(tx *xsql.Tx, orders []*model.OrderMain) (cnt int64, err error) {
	lo := len(orders)
	if lo == 0 {
		return
	}
	f := orders[0].GetFields(&types.FieldMask{Paths: []string{"ctime", "mtime"}})
	lf := len(f)
	hlds := model.InsPlHlds(lf, lo)
	a := make([]interface{}, lf*lo)
	i := 0
	for _, o := range orders {
		vals := o.GetVals(&types.FieldMask{Paths: f}, true)
		copy(a[i:], vals)
		i += lf
	}
	r, err := tx.Exec(fmt.Sprintf(sqlInsertOrderMains, "`"+strings.Join(f, "`,`")+"`", hlds), a...)
	if err != nil {
		return
	}
	cnt, err = r.RowsAffected()
	return
}

//TxInsertOrderDetails 写入orderDetail表
func (d *Dao) TxInsertOrderDetails(tx *xsql.Tx, orders []*model.OrderDetail) (cnt int64, err error) {
	for _, v := range orders {
		v.Encrypt(d.c.Encrypt)
	}
	lo := len(orders)
	if lo == 0 {
		return
	}
	f := orders[0].GetFields(&types.FieldMask{Paths: []string{"ctime", "mtime"}})
	lf := len(f)
	hlds := model.InsPlHlds(lf, lo)
	a := make([]interface{}, lf*lo)
	i := 0
	for _, o := range orders {
		vals := o.GetVals(&types.FieldMask{Paths: f}, true)
		copy(a[i:], vals)
		i += lf
	}
	r, err := tx.Exec(fmt.Sprintf(sqlInsertOrderDetails, "`"+strings.Join(f, "`,`")+"`", hlds), a...)
	if err != nil {
		cnt = 0
		return
	}
	cnt, err = r.RowsAffected()
	return
}

//TxInsertOrderSKUs 增加订单Sku
func (d *Dao) TxInsertOrderSKUs(tx *xsql.Tx, orders []*model.OrderSKU) (cnt int64, err error) {
	lo := len(orders)
	if lo == 0 {
		return
	}
	f := orders[0].GetFields(&types.FieldMask{Paths: []string{"ctime", "mtime"}})
	lf := len(f)
	hlds := model.InsPlHlds(lf, lo)
	a := make([]interface{}, lf*lo)
	i := 0
	for _, o := range orders {
		vals := o.GetVals(&types.FieldMask{Paths: f}, true)
		copy(a[i:], vals)
		i += lf
	}
	r, err := tx.Exec(fmt.Sprintf(sqlInsertOrderSkus, "`"+strings.Join(f, "`,`")+"`", hlds), a...)
	if err != nil {
		cnt = 0
		return
	}
	cnt, err = r.RowsAffected()
	return
}

//OrderID 获取订单号的方法
func (d *Dao) OrderID(ctx context.Context, n int) ([]int64, error) {
	if n <= 0 {
		return nil, nil
	}
	if _, ok := d.c.URLs["basecenter"]; !ok {
		return nil, errors.New("miss basecenter's url conf")
	}
	u := d.c.URLs["basecenter"]
	var res struct {
		Errno int32   `json:"errno"`
		Data  []int64 `json:"data"`
	}
	uv := url.Values{}
	uv.Set("app_id", d.c.BaseCenter.AppID)
	uv.Set("app_token", d.c.BaseCenter.Token)
	uv.Set("count", fmt.Sprintf("%d", n))
	if err := d.httpClientR.Get(ctx, u+"/orderid/get", "", uv, &res); err != nil {
		return nil, err
	}
	if len(res.Data) == 0 {
		return nil, ecode.TicketGetOidFail
	}
	return res.Data, nil
}

//pickOrderSQL 查订单的sql，返回sqls:{0:查列表的sql, 1:查数量的sql}, args:需要代入的变量
func pickOrderSQL(q *model.OrderMainQuerier, fields []string) (sqls []string, args []interface{}) {
	sqls = make([]string, 2)
	var f string
	if len(fields) > 0 {
		f = "`" + strings.Join(fields, "`,`") + "`"
	}
	if l := len(q.OrderID); l > 0 {
		if f != "" {
			sqls[0] = fmt.Sprintf(sqlGetOrders, f, strings.Repeat(",?", l)[1:])
		}
		args = make([]interface{}, len(q.OrderID))
		for k, v := range q.OrderID {
			args[k] = v
		}
		return
	}
	if l := len(q.Status); l > 0 && q.ItemID > 0 {
		if f != "" {
			sqls[0] = fmt.Sprintf(sqlGetUserItemOrders, f, strings.Repeat(",?", l)[1:])
		}
		args = make([]interface{}, 2+l)
		args[0] = q.UID
		args[1] = q.ItemID
		for k, v := range q.Status {
			args[k+2] = v
		}
		return
	}
	if q.OrderBy == "" {
		q.OrderBy = DefaultOrderOrderBy
	}
	if q.Limit == 0 {
		q.Limit = DefaultOrderPSize
	}
	if f != "" {
		sqls[0] = fmt.Sprintf(sqlGetUserOrders, f, q.OrderBy, q.Offset, q.Limit)
	}
	sqls[1] = sqlCountUserOrders
	args = []interface{}{q.UID}
	return
}

//oidCacheKey orderId列表的缓存key, 返回空不缓存列表
func oidCacheKey(q *model.OrderMainQuerier) string {
	s := ""
	if len(q.OrderID) > 0 {
		return s
	}
	if q.OrderBy == "" {
		q.OrderBy = DefaultOrderOrderBy
	}
	if q.Limit == 0 {
		q.Limit = DefaultOrderPSize
	}
	//仅缓存第一页
	if q.Offset > 0 || q.Limit != DefaultOrderPSize || strings.ToLower(q.OrderBy) != DefaultOrderOrderBy {
		return ""
	}
	if q.UID != "" {
		s = "&uid=" + q.UID
	}
	if q.ItemID > 0 {
		s += fmt.Sprintf("&item_id=%d", q.ItemID)
	}
	fs := map[string][]int16{
		"status":        q.Status,
		"sub_status":    q.SubStatus,
		"refund_status": q.RefundStatus,
	}
	for k, v := range fs {
		if len(v) > 0 {
			s += fmt.Sprintf("&%s=%s", k, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v)), ","), "[]"))
		}
	}
	ls := len(s)
	if ls > 32 {
		return fmt.Sprintf("%x", md5.Sum([]byte(s[1:])))
	}
	if ls > 2 {
		return s[1:]
	}
	//空缓存key
	return "nil"
}

//RawGetSettleOrders 获取待结算订单
func (d *Dao) RawGetSettleOrders(ctx context.Context, bt time.Time, et time.Time, id int64, size int) (res *model.SettleOrders, offset int64, err error) {
	q := fmt.Sprintf(sqlGetSettleCompareOrders, "?,?")
	if size > 200 || size <= 0 {
		size = 200
	}
	r, err := d.db.Query(ctx, q, bt.Format("2006-01-02"), et.Format("2006-01-02"), id, consts.OrderStatusPaid, consts.OrderStatusRefunded, size+1)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	defer r.Close()
	res = &model.SettleOrders{}
	res.Data = make([]*model.SettleOrder, 0, size+1)
	for r.Next() {
		o := &model.SettleOrder{}
		if err = r.Scan(&o.ID, &o.OrderID); err != nil {
			return
		}
		res.Data = append(res.Data, o)
	}
	l := len(res.Data)
	if l > size {
		res.Data = res.Data[:l-1]
		offset = res.Data[l-2].ID
	}
	return
}

//RawGetSettleRefunds 获取待结算退款订单
func (d *Dao) RawGetSettleRefunds(ctx context.Context, bt time.Time, et time.Time, id int64, size int) (res *model.SettleOrders, offset int64, err error) {
	q := fmt.Sprintf(sqlGetSettleCompareRefunds, "?")
	if size > 200 || size <= 0 {
		size = 200
	}
	r, err := d.db.Query(ctx, q, bt.Format("2006-01-02"), et.Format("2006-01-02"), id, consts.RefundTxStatusSucc, size+1)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	defer r.Close()
	res = &model.SettleOrders{}
	res.Data = make([]*model.SettleOrder, 0)
	for r.Next() {
		o := &model.SettleOrder{}
		if err = r.Scan(&o.RefID, &o.OrderID, &o.RefundApplyTime); err != nil {
			return
		}
		res.Data = append(res.Data, o)
	}
	l := len(res.Data)
	if l > size {
		res.Data = res.Data[:l-1]
		offset = res.Data[l-2].RefID
	}
	return
}
