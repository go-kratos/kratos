package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/growup/dao"
	"go-common/app/admin/main/growup/dao/resource"
	"go-common/app/admin/main/growup/model"

	"go-common/library/log"
)

const (
	_datetimeLayout = "2006-01-02 15:04:05"
	_dateLayout     = "2006-01-02"
)

func (s *Service) tradeDao() dao.TradeDao {
	return s.dao
}

// SyncGoods . secret portal...
func (s *Service) SyncGoods(c context.Context, gt int) (eff int64, err error) {
	if gt != int(model.GoodsVIP) {
		err = fmt.Errorf("unsupported goodsType(%v)", gt)
		return
	}
	// query
	vips, err := resource.VipProducts(c)
	if err != nil {
		return
	}
	existing := make(map[string]*model.GoodsInfo)
	allGoods, err := s.tradeDao().GoodsList(c, fmt.Sprintf("is_deleted=0 AND goods_type=%d", gt), 0, 200)
	if err != nil {
		return
	}
	for _, v := range allGoods {
		existing[v.ProductID] = v
	}
	// hardcoded <vip_product_duration_in_month, external_resource_id> k-v pairs
	m := map[int32]int64{
		1:  16,
		3:  17,
		12: 18,
	}
	// add incrementally
	newGoods := make([]string, 0)
	for k, v := range vips {
		if _, ok := existing[k]; ok {
			continue
		}
		rid, ok := m[v.Month]
		if !ok {
			err = fmt.Errorf("unkonwn vip goods, month(%d)", v.Month)
		}
		newGoods = append(newGoods, fmt.Sprintf("('%s', %d, %d, %d, %d)", v.ProductID, rid, gt, model.DisplayOff, 100))
	}
	if len(newGoods) == 0 {
		return
	}
	fields := "ex_product_id, ex_resource_id, goods_type, is_display, discount"
	return s.tradeDao().AddGoods(c, fields, strings.Join(newGoods, ","))
}

// GoodsList .
func (s *Service) GoodsList(c context.Context, from, limit int) (total int64, res []*model.GoodsInfo, err error) {
	if total, err = s.tradeDao().GoodsCount(c, "is_deleted=0"); err != nil {
		return
	}
	if total == 0 {
		res = make([]*model.GoodsInfo, 0)
		return
	}
	if res, err = s.tradeDao().GoodsList(c, "is_deleted=0", from, limit); err != nil || len(res) == 0 {
		return
	}
	// external information of vip goods
	var vips map[string]*model.GoodsInfo
	if vips, err = resource.VipProducts(c); err != nil {
		return
	}
	for _, target := range res {
		if src, ok := vips[target.ProductID]; ok {
			model.MergeExternal(target, src)
		}
		target.GoodsTypeDesc = target.GoodsType.Desc()
	}
	return
}

// UpdateGoodsInfo by ID
func (s *Service) UpdateGoodsInfo(c context.Context, discount int, ID int64) (int64, error) {
	// select and diff before update?
	set := fmt.Sprintf("discount=%d", discount)
	return s.tradeDao().UpdateGoods(c, set, "", []int64{ID})
}

// UpdateGoodsDisplay by IDs
func (s *Service) UpdateGoodsDisplay(c context.Context, status model.DisplayStatus, IDs []int64) (eff int64, err error) {
	switch status {
	case model.DisplayOn:
		eff, err = s.onlineGoods(c, IDs)
	case model.DisplayOff:
		eff, err = s.offlineGoods(c, IDs)
	default:
		err = fmt.Errorf("illegal display status(%v)", status)
	}
	return
}

// onlineGoods by IDs
func (s *Service) onlineGoods(c context.Context, IDs []int64) (int64, error) {
	now := time.Now().Format(_datetimeLayout)
	set := fmt.Sprintf("is_display=%d, display_on_time='%s'", model.DisplayOn, now)
	return s.tradeDao().UpdateGoods(c, set, "is_deleted=0", IDs)
}

// offlineGoods by IDs
func (s *Service) offlineGoods(c context.Context, IDs []int64) (int64, error) {
	now := time.Now().Format(_datetimeLayout)
	set := fmt.Sprintf("is_display=%d, display_off_time='%s'", model.DisplayOff, now)
	return s.tradeDao().UpdateGoods(c, set, "is_deleted=0", IDs)
}

// OrderStatistics .
func (s *Service) OrderStatistics(c context.Context, arg *model.OrderQueryArg) (data interface{}, err error) {
	if pass := preprocess(c, arg); !pass {
		return
	}
	where := orderQueryStr(arg)
	var orders []*model.OrderInfo
	if orders, err = s.orderAll(c, where); err != nil {
		return
	}
	for _, v := range orders {
		v.GenDerived()
	}
	data = orderStatistics(orders, arg.StartTime, arg.EndTime, arg.TimeType)
	return
}

// orderAll . be careful using this
func (s *Service) orderAll(c context.Context, where string) (orders []*model.OrderInfo, err error) {
	offset, size := 0, 2000
	for {
		list, err := s.tradeDao().OrderList(c, where, offset, size)
		if err != nil {
			return nil, err
		}
		orders = append(orders, list...)
		if len(list) < size {
			break
		}
		offset += len(list)
	}
	return
}

// orderStatistics . ugly...
func orderStatistics(orders []*model.OrderInfo, start, end time.Time, timeType model.TimeType) interface{} {
	type orderStatUnit struct {
		orderNum   int64
		totalPrice int64
		totalCost  int64
	}

	m := make(map[time.Time]*orderStatUnit)
	for _, v := range orders {
		date := timeType.RangeStart(v.OrderTime)
		if _, ok := m[date]; !ok {
			m[date] = &orderStatUnit{}
		}
		m[date].orderNum += v.GoodsNum
		m[date].totalCost += v.TotalCost
		m[date].totalPrice += v.TotalPrice
	}

	dates, orderNums, totalCost, totalPrice := make([]string, 0), make([]int64, 0), make([]int64, 0), make([]int64, 0)
	for start.Before(end) {
		next := timeType.Next()(start)
		dates = append(dates, timeType.RangeDesc(start, next))
		if v, ok := m[start]; ok {
			orderNums = append(orderNums, v.orderNum)
			totalCost = append(totalCost, v.totalCost)
			totalPrice = append(totalPrice, v.totalPrice)
		} else {
			orderNums = append(orderNums, 0)
			totalCost = append(totalCost, 0)
			totalPrice = append(totalPrice, 0)
		}
		start = next
	}
	// result
	data := map[string]interface{}{
		"xAxis":        dates,
		"order_num":    orderNums,
		"total_cost":   totalCost,
		"total_income": totalPrice,
	}
	return data
}

// OrderExport .
func (s *Service) OrderExport(c context.Context, arg *model.OrderQueryArg, from, limit int) (res []byte, err error) {
	_, orders, err := s.OrderList(c, arg, from, limit)
	if err != nil {
		return
	}
	records := make([][]string, 0, len(orders)+1)
	records = append(records, model.OrderExportFields())
	for _, v := range orders {
		records = append(records, v.ExportStrings())
	}
	if res, err = FormatCSV(records); err != nil {
		log.Error("FormatCSV error(%v)", err)
	}
	return
}

// OrderList .
func (s *Service) OrderList(c context.Context, arg *model.OrderQueryArg, from, limit int) (total int64, list []*model.OrderInfo, err error) {
	list = make([]*model.OrderInfo, 0)
	if pass := preprocess(c, arg); !pass {
		return
	}
	where := orderQueryStr(arg)
	if total, err = s.tradeDao().OrderCount(c, where); err != nil || total == 0 {
		return
	}
	if list, err = s.tradeDao().OrderList(c, where, from, limit); err != nil || len(list) == 0 {
		return
	}
	// fetch names
	mids := make([]int64, 0)
	for _, v := range list {
		mids = append(mids, v.MID)
	}
	m, err := resource.NamesByMIDs(c, mids)
	if err != nil {
		return
	}
	// generate & merge
	for _, v := range list {
		v.GenDerived().GenDesc()
		if name, ok := m[v.MID]; ok {
			v.Nickname = name
		}
	}
	return
}

func preprocess(c context.Context, arg *model.OrderQueryArg) bool {
	if arg.Nickname != "" {
		mid, err := resource.MidByNickname(c, arg.Nickname)
		if err != nil || mid == 0 {
			return false
		}
		log.Info("resource.MidByNickname name(%s) mid(%d)", arg.Nickname, arg.MID)
		if arg.MID == 0 {
			arg.MID = mid
		}

		if arg.MID != mid {
			log.Error("illegal mid(%d) and nickname(%s) pair", arg.MID, arg.Nickname)
			return false
		}
	}
	arg.StartTime = arg.TimeType.RangeStart(time.Unix(arg.FromTime, 0))
	arg.EndTime = arg.TimeType.RangeEnd(time.Unix(arg.ToTime, 0))
	return true
}

func orderQueryStr(arg *model.OrderQueryArg) string {
	var where []string
	where = append(where, "is_deleted=0")
	{
		// 左开右闭
		where = append(where, fmt.Sprintf("order_time >= '%s'", arg.StartTime.Format(_dateLayout)))
		where = append(where, fmt.Sprintf("order_time < '%s'", arg.EndTime.Format(_dateLayout)))
	}
	if arg.GoodsType > 0 {
		where = append(where, fmt.Sprintf("goods_type=%d", arg.GoodsType))
	}
	if arg.GoodsID != "" {
		where = append(where, fmt.Sprintf("goods_id='%s'", arg.GoodsID))
	}
	if arg.GoodsName != "" {
		where = append(where, fmt.Sprintf("goods_name='%s'", arg.GoodsName))
	}
	if arg.OrderNO != "" {
		where = append(where, fmt.Sprintf("order_no='%s'", arg.OrderNO))
	}
	if arg.MID > 0 {
		where = append(where, fmt.Sprintf("mid=%d", arg.MID))
	}
	return strings.Join(where, " AND ")
}
