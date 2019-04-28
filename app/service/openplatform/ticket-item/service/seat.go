package service

import (
	"context"
	"encoding/json"
	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"strings"

	"github.com/jinzhu/gorm"
)

// SeatInfo 添加/修改座位信息
func (s *ItemService) SeatInfo(c context.Context, info *item.SeatInfoRequest) (res *item.SeatInfoReply, err error) {
	var (
		oriArea           *model.Area
		oriAreaSeatmap    *model.AreaSeatmap
		tx                *gorm.DB
		oriAreaSeatsMap   map[int32]map[int32]int64
		oriAreaSeatsArr   []*model.AreaSeats
		recoveryAreaSeats []int64
		insertAreaSeats   []*model.AreaSeats
	)
	oriArea = &model.Area{
		ID:        info.Area,
		SeatsNum:  info.SeatsNum,
		Width:     info.Width,
		Height:    info.Height,
		RowList:   info.RowList,
		SeatStart: info.SeatStart,
	}

	if err = v.Struct(info); err != nil {
		err = ecode.RequestErr
		return
	}
	// 创建事务
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = ecode.NotModified
		return
	}
	// 更新area表信息
	if err = s.dao.TxUpdateSeat(c, tx, oriArea); err != nil {
		tx.Rollback()
		return
	}
	// 若有seat_map字段，则更新area_seatmap表
	if info.SeatMap != "" {
		oriAreaSeatmap = &model.AreaSeatmap{
			ID:      info.Area,
			SeatMap: info.SeatMap,
		}
		if err = s.dao.TxSaveAreaSeatmap(c, tx, oriAreaSeatmap); err != nil {
			tx.Rollback()
			return
		}
	}
	// 查询并记录现有area_seats表信息到oriAreaSeatsMap
	if oriAreaSeatsArr, err = s.dao.TxGetAreaSeats(c, tx, info.Area); err != nil {
		tx.Rollback()
		return
	}
	oriAreaSeatsMap = make(map[int32]map[int32]int64)
	for _, oriAreaSeats := range oriAreaSeatsArr {
		if _, ok := oriAreaSeatsMap[oriAreaSeats.X]; !ok {
			oriAreaSeatsMap[oriAreaSeats.X] = make(map[int32]int64)
		}
		oriAreaSeatsMap[oriAreaSeats.X][oriAreaSeats.Y] = oriAreaSeats.ID
	}
	// 软删除现有区域下area_seats表记录
	if err = s.dao.TxBatchDeleteAreaSeats(c, tx, info.Area); err != nil {
		tx.Rollback()
		return
	}
	// 比对新旧座位信息，若重合则恢复，若不存在则添加
	recoveryAreaSeats = make([]int64, 0, 100)
	insertAreaSeats = make([]*model.AreaSeats, 0, 100)
	for _, oriSeat := range info.Seats {
		if yMap, ok := oriAreaSeatsMap[oriSeat.X]; ok {
			if id, ok := yMap[oriSeat.Y]; ok {
				recoveryAreaSeats = append(recoveryAreaSeats, id)
				continue
			}
		}
		insertAreaSeats = append(insertAreaSeats, &model.AreaSeats{
			X:       oriSeat.X,
			Y:       oriSeat.Y,
			Label:   oriSeat.Label,
			Bgcolor: oriSeat.Bgcolor,
			Area:    info.Area,
		})
	}
	// 批量添加待添加座位到area_seats表
	if len(insertAreaSeats) > 0 {
		if err = s.dao.TxBatchAddAreaSeats(c, tx, insertAreaSeats); err != nil {
			tx.Rollback()
			return
		}
	}
	// 批量恢复area_seats表的存在记录
	if len(recoveryAreaSeats) > 0 {
		if err = s.dao.TxBatchRecoverAreaSeats(c, tx, recoveryAreaSeats); err != nil {
			tx.Rollback()
			return
		}
	}
	// 提交事务
	if err = s.dao.CommitTran(c, tx); err != nil {
		err = ecode.NotModified
		return
	}
	res = &item.SeatInfoReply{Success: true}
	return
}

// SeatStock 设置座位库存
// 确保座位均未售出，从seat_set表获取座位图（不存在则取area_setmap），
// 遍历请求的SeatInfo数组，补足未定义的票价（从票种复制），生成seat_order表，插入或更新seat_set座位图
func (s *ItemService) SeatStock(c context.Context, info *item.SeatStockRequest) (res *item.SeatStockReply, err error) {
	var (
		seatOrders      []*model.SeatOrder
		seatSet         *model.SeatSet
		seatChartRows   []string
		ticketPrices    []*model.TicketPrice
		ticketPrice     *model.TicketPrice
		priceSymbol     map[int64]string
		symbolPrice     map[string]*model.TicketPrice
		insertPrice     []*model.TicketPrice
		insertSeatOrder []*model.SeatOrder
		seatmap         *model.AreaSeatmap
		seatChart       []byte
		tx              *gorm.DB
	)

	if err = v.Struct(info); err != nil {
		err = ecode.RequestErr
		return
	}
	// 创建事务
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = ecode.NotModified
		return
	}
	// 获取区域不可售座位
	if seatOrders, err = s.dao.TxGetUnsaleableSeatOrders(c, tx, info.Screen, info.Area); err != nil {
		tx.Rollback()
		err = ecode.NotModified
		return
	}
	// 已有不可售座位则不可修改
	if len(seatOrders) != 0 {
		tx.Rollback()
		log.Error("区域（ID:%d）下已有不可售座位，请勿修改！", info.Area)
		err = ecode.NotModified
		return
	}
	// 查找票价设置图
	if seatSet, err = s.dao.TxGetSeatChart(c, tx, info.Screen, info.Area); err != nil {
		tx.Rollback()
		err = ecode.NotModified
		return
	}
	if seatSet.ID == 0 {
		seatSet.AreaID = info.Area
		seatSet.ScreenID = info.Screen
	}
	if seatSet.SeatChart == "" {
		if seatmap, err = s.dao.TxRawAreaSeatmap(c, tx, info.Area); err != nil {
			tx.Rollback()
			err = ecode.NotModified
			return
		}
		seatSet.SeatChart = seatmap.SeatMap
	}
	// 没有设置图则无法生成库存
	if seatSet.SeatChart == "" {
		tx.Rollback()
		log.Error("场次（ID:%d）下没有找到区域（ID:%d）的票价设置图！", info.Screen, info.Area)
		err = ecode.NotModified
		return
	}
	// 反序列化出座位数据
	if err = json.Unmarshal([]byte(seatSet.SeatChart), &seatChartRows); err != nil {
		tx.Rollback()
		log.Error("场次（ID:%d）和区域（ID:%d）下票价设置图数据不可用:%s", info.Screen, info.Area, err)
		err = ecode.NotModified
		return
	}

	// 获取场次下的票价，建立票种和标志映射，标志和票价映射
	if ticketPrices, err = s.dao.TxGetPriceSymbols(c, tx, info.Screen); err != nil {
		tx.Rollback()
		err = ecode.NotModified
		return
	}
	priceSymbol = make(map[int64]string)
	symbolPrice = make(map[string]*model.TicketPrice)
	for _, tp := range ticketPrices {
		priceSymbol[tp.ParentID] = tp.Symbol
		symbolPrice[tp.Symbol] = tp
	}

	// 遍历请求的SeatInfo数组，补足未定义的票价，修改座位图，组装价格图
	insertPrice = make([]*model.TicketPrice, 0)
	insertSeatOrder = make([]*model.SeatOrder, 0)
	alphabetTable := model.AlphabetTable()
	for _, sp := range info.SeatInfo {
		var _symbol string
		if sp.Price == 0 {
			_symbol = "#"
		} else {
			var ok bool
			if _symbol, ok = priceSymbol[sp.Price]; !ok {
				// 枚举_symbol从alphabetTable[len(priceSymbol)]逆序到alphabetTable[0]，若priceSymbol没有用过这个_symbol，则将其作为这个票价的symbol
				for i := len(priceSymbol); i >= 0; i-- {
					_symbol = alphabetTable[i]
					found := false
					for _, symbol := range priceSymbol {
						if symbol == _symbol {
							found = true
							break
						}
					}
					if !found {
						var genID int64
						if genID, err = model.GetTicketIDFromBase(); err != nil {
							tx.Rollback()
							err = ecode.NotModified
							return
						}
						// 取票种
						if ticketPrice, err = s.dao.TxGetParentTicketPrice(c, tx, sp.Price); err != nil {
							tx.Rollback()
							err = ecode.NotModified
							return
						}
						ticketPrice.ID = genID
						ticketPrice.ScreenID = info.Screen
						ticketPrice.Symbol = _symbol
						ticketPrice.ParentID = sp.Price
						ticketPrice.OriginPrice = -1
						ticketPrice.MarketPrice = -1
						ticketPrice.IsSale = 0
						ticketPrice.IsVisible = 0
						ticketPrice.IsRefund = 10
						// 待插入票价
						insertPrice = append(insertPrice, ticketPrice)
						priceSymbol[sp.Price] = _symbol
						symbolPrice[_symbol] = ticketPrice
						break
					}
				}
			}
		}
		// 修改座位图
		bs := []byte(seatChartRows[sp.X])
		bs[sp.Y] = _symbol[0]
		seatChartRows[sp.X] = string(bs)
		// 待插入座位订单
		insertSeatOrder = append(insertSeatOrder, &model.SeatOrder{
			Row:      sp.X,
			Col:      sp.Y,
			PriceID:  symbolPrice[_symbol].ID,
			Price:    symbolPrice[_symbol].Price,
			ScreenID: info.Screen,
			AreaID:   info.Area,
		})
	}
	// 插入票价
	if err = s.dao.TxBatchAddTicketPrice(c, tx, insertPrice); err != nil {
		tx.Rollback()
		err = ecode.NotModified
		return
	}
	// 插入座位订单
	if err = s.dao.TxBatchAddSeatOrder(c, tx, insertSeatOrder); err != nil {
		tx.Rollback()
		err = ecode.NotModified
		return
	}
	// 更新座位图PriceSet的seat_chart
	if seatChart, err = json.Marshal(seatChartRows); err != nil {
		tx.Rollback()
		log.Error("座位图seat_chart序列化失败:%s", err)
		err = ecode.NotModified
		return
	}
	seatSet.SeatChart = string(seatChart)
	if seatSet.ID == 0 {
		if err = s.dao.TxAddSeatChart(c, tx, seatSet); err != nil {
			tx.Rollback()
			err = ecode.NotModified
			return
		}
	} else {
		if err = s.dao.TxUpdateSeatChart(c, tx, seatSet.ID, seatSet.SeatChart); err != nil {
			tx.Rollback()
			err = ecode.NotModified
			return
		}
	}

	// 提交事务
	if err = s.dao.CommitTran(c, tx); err != nil {
		err = ecode.NotModified
		return
	}
	res = &item.SeatStockReply{Success: true}
	return
}

// RemoveSeatOrders 删除坐票票价下所有座位（包括seat_order和seat_set中的记录）
func (s *ItemService) RemoveSeatOrders(c context.Context, info *item.RemoveSeatOrdersRequest) (res *item.RemoveSeatOrdersReply, err error) {
	var (
		tp         *model.TicketPrice
		seatOrders []*model.SeatOrder
		seatSets   []*model.SeatSet
		delIDs     []int64
		areas      []int64
		tx         *gorm.DB
	)
	if err = v.Struct(info); err != nil {
		err = ecode.RequestErr
		return
	}
	// 创建事务
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = ecode.NotModified
		return
	}

	// 查询票价
	if tp, err = s.dao.TxGetTicketPrice(c, tx, info.Price); err != nil {
		tx.Rollback()
		err = ecode.NotModified
		return
	}
	// 确保票价的价格标志存在
	if tp.Symbol == "" {
		tx.Rollback()
		log.Error("票价（ID:%d）不存在", info.Price)
		err = ecode.NotModified
		return
	}
	// 根据场次和票价ID获取可销售座位的ID和区域ID（可能场次冗余？）
	if seatOrders, err = s.dao.TxGetSaleableSeatOrders(c, tx, tp.ScreenID, info.Price); err != nil {
		tx.Rollback()
		err = ecode.NotModified
		return
	}
	// 临时map，key为区域ID，用于去重
	areaMap := make(map[int64]bool)
	// 统计要删除的座位订单ID和区域ID
	delIDs = make([]int64, 0)
	areas = make([]int64, 0)
	for _, seatOrder := range seatOrders {
		delIDs = append(delIDs, seatOrder.ID)
		areaMap[seatOrder.AreaID] = true
	}
	for area := range areaMap {
		areas = append(areas, area)
	}
	// 批量删除座位订单
	if err = s.dao.TxBatchDeleteSeatOrder(c, tx, delIDs); err != nil {
		tx.Rollback()
		err = ecode.NotModified
		return
	}
	// 依次更新区域的座位订单
	if seatSets, err = s.dao.TxGetSeatCharts(c, tx, tp.ScreenID, areas); err != nil {
		tx.Rollback()
		err = ecode.NotModified
		return
	}
	for _, seatSet := range seatSets {
		seatSet.SeatChart = strings.Replace(seatSet.SeatChart, tp.Symbol, "#", -1)
		if err = s.dao.TxUpdateSeatChart(c, tx, seatSet.ID, seatSet.SeatChart); err != nil {
			tx.Rollback()
			err = ecode.NotModified
			return
		}
	}
	// 提交事务
	if err = s.dao.CommitTran(c, tx); err != nil {
		err = ecode.NotModified
		return
	}
	res = &item.RemoveSeatOrdersReply{Areas: areas}
	return
}
