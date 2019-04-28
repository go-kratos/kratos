package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-common/app/admin/main/growup/conf"
	"go-common/app/admin/main/growup/dao"
	"go-common/app/admin/main/growup/dao/shell"
	"go-common/app/admin/main/growup/model/offlineactivity"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"io/ioutil"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	uploadFilePath = "/data/uploadfiles/"
	dateFmt        = "20060102"
	exportMaxCount = 1000
)

func (s *Service) offlineactivityCheckSendDbProc() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("write stat data Runtime error caught, try recover: %+v", r)
			go s.offlineactivityCheckSendDbProc()
		}
	}()
	var timer = time.NewTicker(60 * time.Second)
	for {
		select {
		case <-timer.C:
		case <-s.chanCheckDb:
		}
		s.checkResultNotSendingOrder()
	}
}

func (s *Service) offlineactivityCheckShellOrderProc() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("write stat data Runtime error caught, try recover: %+v", r)
			go s.offlineactivityCheckShellOrderProc()
		}
	}()
	var timer = time.NewTimer(10 * time.Second)
	var timerStart = false
	var resultIds []int64
	for {
		select {
		case orderID := <-s.chanCheckShellOrder:
			var all = []*offlineactivity.OfflineActivityResult{orderID}
			for {
			drain:
				for {
					select {
					case orderID := <-s.chanCheckShellOrder:
						all = append(all, orderID)
						if len(all) >= 100 {
							break drain
						}
					default:
						break drain
					}
				}
				s.checkShellOrder(context.Background(), all)
				if len(all) < 100 {
					break
				}
				all = nil
			}
		case resultID := <-s.chanCheckActivity:
			if !timerStart {
				timer.Reset(10 * time.Second)
				timerStart = true
			}
			resultIds = append(resultIds, resultID)
		case <-timer.C:
			// 去检查对应活动的状态
			timerStart = false
			var ids []int64
			ids = append(ids, resultIds...)
			s.checkActivityState(context.Background(), ids)
			// clear map
			resultIds = nil
		}
	}
}

//PreAddOfflineActivity just test add resutl
func (s *Service) PreAddOfflineActivity(ctx context.Context, arg *offlineactivity.AddActivityArg) (res *offlineactivity.PreAddActivityResult, err error) {
	if arg == nil {
		return
	}
	res = new(offlineactivity.PreAddActivityResult)

	var bonusList = arg.BonusList
	var hasBonus = false
	res.BonusType = arg.BonusType
	for _, bonus := range bonusList {
		if bonus.TotalMoney <= 0 {
			err = fmt.Errorf("bonus money < 0, money=%f", bonus.TotalMoney)
			log.Error(err.Error())
			return
		}
		res.TotalMoney += bonus.TotalMoney
		if bonus.Filename != "" {
			var fullpath = uploadFilePath + bonus.Filename
			var filecontent []byte
			filecontent, err = ioutil.ReadFile(fullpath)
			if err != nil {
				log.Error("read file fail, path=%s", fullpath)
				return
			}
			bonus.Mids = string(filecontent)
			bonus.MidList = dao.ParseMidsFromString(bonus.Mids)
			bonus.MemberCount = int64(len(bonus.MidList))
			if bonus.MemberCount == 0 {
				log.Warn("no mid for this bonus, bonus：money=%d", bonus.TotalMoney)
				continue
			}
		} else {
			bonus.MidList = dao.ParseMidsFromString(bonus.Mids)
			bonus.MemberCount = int64(len(bonus.MidList))
			if bonus.MemberCount == 0 {
				log.Warn("no mid for this bonus, bonus：money=%d", bonus.TotalMoney)
				continue
			}
		}
		res.MemberCount += bonus.MemberCount
		hasBonus = true
	}
	// 如果一个bonus都没有，则什么也不做
	if !hasBonus {
		log.Error("no bonus")
		err = fmt.Errorf("没有获奖人员信息/解析文件失败(请上传csv/txt等文本格式文件)")
		return
	}

	return
}

//AddOfflineActivity add offline activity
func (s *Service) AddOfflineActivity(ctx context.Context, arg *offlineactivity.AddActivityArg) (res *offlineactivity.AddActivityResult, err error) {
	if arg == nil {
		return
	}
	var bonusList = arg.BonusList
	var hasBonus = false
	var bc = ctx.(*blademaster.Context)
	var cookie, _ = bc.Request.Cookie("username")
	if cookie != nil {
		arg.Creator = cookie.Value
	}
	for _, bonus := range bonusList {
		if bonus.TotalMoney <= 0 {
			err = fmt.Errorf("bonus money < 0, money=%f", bonus.TotalMoney)
			log.Error(err.Error())
			return
		}
		if bonus.Filename != "" {
			var fullpath = uploadFilePath + bonus.Filename
			var filecontent []byte
			filecontent, err = ioutil.ReadFile(fullpath)
			if err != nil {
				log.Error("read file fail, path=%s", fullpath)
				return
			}
			bonus.Mids = string(filecontent)
		} else {
			bonus.MidList = dao.ParseMidsFromString(bonus.Mids)
			bonus.MemberCount = int64(len(bonus.MidList))
			if bonus.MemberCount == 0 {
				log.Warn("no mid for this bonus, bonus：money=%d", bonus.TotalMoney)
				continue
			}
		}

		hasBonus = true
	}
	// 如果一个bonus都没有，则什么也不做
	if !hasBonus {
		log.Error("no bonus")
		err = fmt.Errorf("no bonus info")
		return
	}

	go func() {
		err = s.dao.OfflineActivityAddActivity(context.Background(), arg)
		if err != nil {
			log.Error("offline add fail")
			return
		}
		log.Info("offline add ok")
		if len(s.chanCheckDb) == 0 {
			s.chanCheckDb <- 1
		}
	}()

	return
}

//ShellCallback shell callback
func (s *Service) ShellCallback(ctx context.Context, arg *shell.OrderCallbackParam) (err error) {
	if arg == nil {
		log.Error("arg is nil")
		return ecode.RequestErr
	}
	var result = shell.OrderCallbackJSON{}
	if err = json.Unmarshal([]byte(arg.MsgContent), &result); err != nil {
		log.Error("msgid=%s, unmarshal msg content fail, err=%s, msgcontent=%s", err, arg.MsgID, string(arg.MsgContent))
		return
	}
	log.Info("order id=%s, handle shell callback, status=%s", result.ThirdOrderNo, result.Status)
	if result.CustomerID != conf.Conf.ShellConf.CustomID {
		log.Error("order id=%s, customerid not the same, give=%s, expect=%s", result.ThirdOrderNo, result.CustomerID, conf.Conf.ShellConf.CustomID)
		return
	}
	orderInfo, err := s.dao.ShellCallbackUpdate(ctx, &result, arg.MsgID)
	if err == nil {
		s.queueToUpdateActivityState(orderInfo.ResultID)
	}
	return
}

//OfflineActivityQueryActivity query activity
func (s *Service) OfflineActivityQueryActivity(ctx context.Context, arg *offlineactivity.QueryActivityByIDArg) (res *offlineactivity.QueryActivityResult, err error) {
	var db = s.dao.OfflineActivityGetDB()
	// make sure it's valid
	var limit, offset = arg.CheckPageValidation()
	if arg.ExportFormat() != "" {
		limit = exportMaxCount
	}
	var now = time.Now()
	if arg.FromDate == "" {
		arg.FromDate = now.AddDate(0, -1, 0).Format(dateFmt)
	}
	if arg.ToDate == "" {
		arg.ToDate = now.Format(dateFmt)
	}
	// 1.查询 OfflineActivityInfo 与 OfflineActivityBonus 数据，然后聚合
	var activityList []*offlineactivity.OfflineActivityInfo
	var total = 1
	if arg.ID != 0 {
		if err = db.Where("id=?", arg.ID).Find(&activityList).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Error("fail to get from db, err=%s", err)
			return
		}
		total = len(activityList)
	} else {
		if arg.FromDate == "" || arg.ToDate == "" {
			log.Error("request error, fromdate or todate is nill")
			err = ecode.RequestErr
			return
		}

		var todate, e = time.Parse(dateFmt, arg.ToDate)
		err = e
		if err != nil {
			log.Error("todate format err, todate=%s, err=%s", arg.ToDate, err)
			return
		}
		todate = todate.AddDate(0, 0, 1)
		var todatestr = todate.Format(dateFmt)
		if err = db.Table(offlineactivity.TableOfflineActivityInfo).Where("ctime>=? and ctime<=?", arg.FromDate, todatestr).Count(&total).Error; err != nil {
			log.Error("fail to get from db, err=%s", err)
			return
		}

		if err = db.Where("ctime>=? and ctime<=?", arg.FromDate, todatestr).Order("id desc").Offset(offset).Limit(limit).Find(&activityList).Error; err != nil {
			log.Error("fail to get from db, err=%s", err)
			return
		}
	}

	var activityIDs []int64
	for _, v := range activityList {
		activityIDs = append(activityIDs, v.ID)
	}
	if len(activityList) == 0 {
		log.Warn("0 activity list")
		return
	}

	// 查询bonus info
	var bonusList []*offlineactivity.OfflineActivityBonus
	if err = db.Where("activity_id in (?)", activityIDs).Find(&bonusList).Error; err != nil {
		log.Error("fail to get from db, err=%s", err)
		return
	}

	res = new(offlineactivity.QueryActivityResult)
	var activityMap = make(map[int64]*offlineactivity.QueryActivityInfo, len(activityList))
	for _, v := range activityList {
		var info = new(offlineactivity.QueryActivityInfo)
		info.CopyFromActivityDB(v)
		activityMap[info.ID] = info
	}

	for _, v := range bonusList {
		var info, _ = activityMap[v.ActivityID]
		if info == nil {
			continue
		}
		info.TotalMoney += offlineactivity.GetMoneyFromDb(v.TotalMoney)
		info.MemberCount += v.MemberCount
	}

	for _, v := range activityList {
		if info, ok := activityMap[v.ID]; ok {
			res.Result = append(res.Result, info)
		}
	}

	res.PageResult = arg.ToPageResult(total)
	return
}

// OfflineActivityQueryUpBonusSummary query up bonus info
func (s *Service) OfflineActivityQueryUpBonusSummary(ctx context.Context, arg *offlineactivity.QueryUpBonusByMidArg) (res *offlineactivity.QueryUpBonusByMidResult, err error) {
	var limit, offset = arg.CheckPageValidation()
	if arg.ExportFormat() != "" {
		limit = exportMaxCount
	}
	// 查询所有的 result for mid
	// 区分已结算、未结算
	// 已结算= 成功， 未结算= 初始、发送、等待
	upResult, total, err := s.dao.OfflineActivityGetUpBonusResult(ctx, true, limit, offset, "mid=?", arg.Mid)
	if err != nil {
		log.Error("get from up result fail, err=%s", err)
		return
	}
	res = new(offlineactivity.QueryUpBonusByMidResult)
	res.PageResult = arg.ToPageResult(total)

	var bonusMap = make(map[int64]*offlineactivity.UpSummaryBonusInfo)
	for _, v := range upResult {
		var bonusInfo, ok = bonusMap[v.Mid]
		if !ok {
			bonusInfo = &offlineactivity.UpSummaryBonusInfo{
				Mid: v.Mid,
			}
			bonusMap[v.Mid] = bonusInfo
		}
		switch offlineactivity.ActivityState(v.State) {
		case offlineactivity.ActivityStateSucess:
			bonusInfo.BilledMoney += offlineactivity.GetMoneyFromDb(v.BonusMoney)
			if bonusInfo.TmpBillTime < v.MTime {
				bonusInfo.TmpBillTime = v.MTime
			}
		case offlineactivity.ActivityStateInit, offlineactivity.ActivityStateSending, offlineactivity.ActivityStateWaitResult:
			bonusInfo.UnbilledMoney += offlineactivity.GetMoneyFromDb(v.BonusMoney)
		}
	}

	for _, v := range bonusMap {
		v.Finish()
		res.Result = append(res.Result, v)
	}
	return
}

//OfflineActivityQueryUpBonusByActivity get bonus info group by activity
func (s *Service) OfflineActivityQueryUpBonusByActivity(ctx context.Context, arg *offlineactivity.QueryUpBonusByMidArg) (res *offlineactivity.QueryUpBonusByActivityResult, err error) {
	var limit, offset = arg.CheckPageValidation()
	if arg.ExportFormat() != "" {
		limit = exportMaxCount
	}
	// 查询所有的 result for mid
	// 区分已结算、未结算
	// 已结算= 成功， 未结算= 初始、发送、等待
	upResult, total, err := s.dao.OfflineActivityGetUpBonusByActivityResult(ctx, limit, offset, arg.Mid)
	if err != nil {
		log.Error("get from up result fail, err=%s", err)
		return
	}
	res = new(offlineactivity.QueryUpBonusByActivityResult)
	res.PageResult = arg.ToPageResult(total)

	// [mid][activity_id], 按Activity做聚合
	var bonusMap = make(map[int64]map[int64]*offlineactivity.UpSummaryBonusInfo)
	for _, v := range upResult {
		var bonusActivityMap, ok = bonusMap[v.Mid]
		if !ok {
			bonusActivityMap = make(map[int64]*offlineactivity.UpSummaryBonusInfo)
			bonusMap[v.Mid] = bonusActivityMap
		}
		bonusInfo, ok := bonusActivityMap[v.ActivityID]
		if !ok {
			bonusInfo = &offlineactivity.UpSummaryBonusInfo{
				Mid:        v.Mid,
				ActivityID: v.ActivityID,
			}
			bonusActivityMap[v.ActivityID] = bonusInfo
		}
		var bonusMoney = offlineactivity.GetMoneyFromDb(v.BonusMoney)
		switch offlineactivity.ActivityState(v.State) {
		case offlineactivity.ActivityStateSucess:
			bonusInfo.BilledMoney += bonusMoney
			if bonusInfo.TmpBillTime < v.MTime {
				bonusInfo.TmpBillTime = v.MTime
			}
			bonusInfo.TotalBonusMoney += bonusMoney
		case offlineactivity.ActivityStateInit, offlineactivity.ActivityStateSending, offlineactivity.ActivityStateWaitResult:
			bonusInfo.UnbilledMoney += bonusMoney
			bonusInfo.TotalBonusMoney += bonusMoney
		}
	}

	for _, activityMap := range bonusMap {
		for _, v := range activityMap {
			v.Finish()
			res.Result = append(res.Result, v)
		}
	}
	return
}

//OfflineActivityQueryActivityByMonth activity by month
func (s *Service) OfflineActivityQueryActivityByMonth(ctx context.Context, arg *offlineactivity.QueryActvityMonthArg) (res *offlineactivity.QueryActivityMonthResult, err error) {
	var db = s.dao.OfflineActivityGetDB()

	var bonusInfo []*offlineactivity.OfflineActivityBonus
	var limit, offset = 100, 0
	var lastCount = limit
	for limit == lastCount {
		var bonusInfoTmp []*offlineactivity.OfflineActivityBonus
		if err = db.Find(&bonusInfoTmp).Error; err != nil {
			log.Error("get from db fail, err=%s", err)
			return
		}
		bonusInfo = append(bonusInfo, bonusInfoTmp...)
		lastCount = len(bonusInfoTmp)
		offset += lastCount
	}

	var now = time.Now()
	var dateStr = now.Format(dateFmt)
	var monthDataMap = make(map[string]*offlineactivity.ActivityMonthInfo)
	for _, v := range bonusInfo {
		var date = v.CTime.Time().Format("200601")
		var monthData, ok = monthDataMap[date]
		if !ok {
			monthData = &offlineactivity.ActivityMonthInfo{
				CreateTime: date,
			}
			monthDataMap[date] = monthData
			monthData.GenerateDay = dateStr
		}
		monthData.AddBonus(v)
	}

	var monthInfoS []*offlineactivity.ActivityMonthInfo
	for _, v := range monthDataMap {
		monthInfoS = append(monthInfoS, v)
		v.Finish()
	}

	sort.Slice(monthInfoS, func(i, j int) bool {
		return monthInfoS[i].CreateTime > monthInfoS[j].CreateTime
	})

	var lastIndex = len(monthInfoS) - 1
	if lastIndex >= 0 {
		monthInfoS[lastIndex].TotalMoneyAccumulate = monthInfoS[lastIndex].TotalBonusMoneyMonth
	}

	for i := len(monthInfoS) - 1; i > 0; i-- {
		monthInfoS[i-1].TotalMoneyAccumulate = monthInfoS[i].TotalMoneyAccumulate + monthInfoS[i-1].TotalBonusMoneyMonth
	}

	res = new(offlineactivity.QueryActivityMonthResult)
	res.Result = monthInfoS
	return
}

/*
	查询对应的activity_result表state=0 && bonus_type = 1 && ctime <= 86400
	for each item state=0 and activity_id = ? limit 100
		if 没有order_id，
			生成order_id，写入数据库, result表与shellorder表
			写入成功，记录等待发送
		if 有order id，
			// 说明已经经历过上一步，但是发送失败？那么加入到检查order id的队列
			记录等待发送

	批量进行发送-> shell
	if	发送成功
		更新所有的item状态 -> 1
	else if 失败：
		错误码是8002999997->order_id有重复
		更新所有的item order id = ''，等待重试
*/
func (s *Service) checkResultNotSendingOrder() {
	var db, err = gorm.Open("mysql", conf.Conf.ORM.Growup.DSN)
	if err != nil {
		log.Error("open db fail, dsn=%s", conf.Conf.ORM.Growup.DSN)
		return
	}
	defer db.Close()
	db.LogMode(false)
	// 查询对应的activity_result表state=0 && bonus_type = 1 && diff(ctime) <= 86400
	var minTime = time.Now().Add(-time.Hour * 24)
	var limit = 100
	var lastCount = limit
	var lastID int64
	for limit == lastCount {
		var needSendResult = make([]*offlineactivity.OfflineActivityResult, limit)
		if err = db.Select("id, mid, bonus_money, order_id, activity_id, bonus_id").
			Where("state=0 and bonus_type=1 and ctime>=? and id>?", minTime, lastID).Limit(limit).
			Find(&needSendResult).Error; err != nil {
			log.Error("get result fail, err=%s", err)
			return
		}

		lastCount = len(needSendResult)
		var now = time.Now()
		var sendRequestResult []*offlineactivity.OfflineActivityResult
		for _, activityResult := range needSendResult {
			if activityResult.ID > lastID {
				lastID = activityResult.ID
			}
			// if 没有order_id，
			if activityResult.OrderID == "" {
				activityResult.OrderID = generateOrderID(activityResult, now)
				if err = s.offlineActivityUpdateResultForSend(db, activityResult); err != nil {
					log.Warn("fail update result, err=%s, will retry next time", err)
					continue
				}
				log.Info("update result ok, mid=%d, activity_id=%d, order_id=%s", activityResult.Mid, activityResult.ActivityID, activityResult.OrderID)
				sendRequestResult = append(sendRequestResult, activityResult)
			} else {
				// if 有order id，
				go func(res *offlineactivity.OfflineActivityResult) {
					s.chanCheckShellOrder <- res
					log.Info("order has id, need to check, order id=%s, mid=%d, activityid=%d, bonusid=%d",
						res.OrderID, res.Mid, res.ActivityID, res.BonusID)
				}(activityResult)
			}
			if len(sendRequestResult) >= 10 {
				err = s.sendRequestAndUpdate(db, sendRequestResult)
				if err != nil {
					log.Error("send request err, err=%s", err)
				} else {
					log.Info("send to shell ok, length=%d", len(sendRequestResult))
				}

				sendRequestResult = nil
			}
		}
		if len(sendRequestResult) > 0 {
			err = s.sendRequestAndUpdate(db, sendRequestResult)
			if err != nil {
				log.Error("send request err, err=%s", err)
			} else {
				log.Info("send to shell ok, length=%d", len(sendRequestResult))
			}
		}
	}

}

func (s *Service) sendRequestAndUpdate(db *gorm.DB, sendRequestResult []*offlineactivity.OfflineActivityResult) (err error) {
	switch {
	default:
		var res *shell.OrderResponse
		if res, err = s.sendShellOrder(context.Background(), sendRequestResult); err != nil || res == nil {
			if err != nil {
				log.Error("fail to send to shell order, err=%s", err)
			}
			break
		}
		var ids []int64
		for _, v := range sendRequestResult {
			ids = append(ids, v.ID)
		}
		// 需要区分是否返回 错误码是8002999997,表示有重复
		if res.Errno == 8002999997 {
			log.Error("fail to send request, err=%d, msg=%s，重复！", res.Errno, res.Msg)
			// 将order id 更新为""
			if err = db.Table(offlineactivity.TableOfflineActivityResult).Where("id in (?)", ids).
				Update("order_id=''").Error; err != nil {
				log.Error("fail to update order id for duplicate order, id=%v", ids)
				break
			}
			return
		} else if res.Errno != 0 {
			log.Error("fail to send request, err=%d, msg=%s", res.Errno, res.Msg)
			return
		}

		if err = db.Table(offlineactivity.TableOfflineActivityResult).Where("id in (?)", ids).
			Update("state", offlineactivity.ActivityStateWaitResult).Error; err != nil {
			log.Error("fail to update state, id=%v", ids)
			break
		}
		log.Info("send request to shell order")
	}
	return
}

func generateOrderID(result *offlineactivity.OfflineActivityResult, tm time.Time) string {
	var order = fmt.Sprintf("%s%04d%010d%04d", time.Now().Format("20060102150405"), (result.ActivityID*100+result.BonusID)%10000, result.Mid, rand.Int()%10000)
	if len(order) > 32 {
		order = order[:32]
	}
	return order
}

func (s *Service) offlineActivityUpdateResultForSend(db *gorm.DB, result *offlineactivity.OfflineActivityResult) (err error) {
	if result == nil {
		err = fmt.Errorf("nil pointer")
		return
	}
	var tx = db.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	// create
	var query = tx.Select("order_id").Where("order_id=''").Save(result)
	err = db.Error
	if err != nil {
		log.Error("err save offline result, err=%s", err)
		return
	}
	if query.RowsAffected == 0 {
		var msg = fmt.Sprintf("update order fail, it may already be sent, result_id=%d, mid=%d, activity_id=%d", result.ID, result.Mid, result.ActivityID)
		log.Warn(msg)
		err = fmt.Errorf(msg)
		return
	}
	var shellOrder = offlineactivity.OfflineActivityShellOrder{
		OrderID:  result.OrderID,
		ResultID: result.ID,
	}
	if err = tx.Save(&shellOrder).Error; err != nil {
		log.Error("err insert offline shell order, err=%s", err)
		return
	}

	return tx.Commit().Error
}

func (s *Service) sendShellOrder(ctx context.Context, needSendResult []*offlineactivity.OfflineActivityResult) (res *shell.OrderResponse, err error) {
	if len(needSendResult) == 0 {
		log.Warn("no need to send, len=0")
		return
	}
	var nowtime = time.Now()
	var now = time.Now().UnixNano() / int64(time.Millisecond)
	var request = shell.OrderRequest{
		ProductName: "活动奖励",
		NotifyURL:   conf.Conf.ShellConf.CallbackURL,
		Rate:        "1.0",
		SignType:    "MD5",
		Timestamp:   strconv.Itoa(int(beginOfDay(nowtime).UnixNano() / int64(time.Millisecond))),
	}
	for _, item := range needSendResult {
		var money = fmt.Sprintf("%0.2f", float64(offlineactivity.GetMoneyFromDb(item.BonusMoney)))
		var orderInfo = shell.OrderInfo{
			Mid:          item.Mid,
			Brokerage:    money,
			ThirdCoin:    money,
			ThirdOrderNo: item.OrderID,
			ThirdCtime:   strconv.Itoa(int(now)),
		}
		request.Data = append(request.Data, orderInfo)
	}
	//log.Info("request=%+v", request)
	res, err = s.shellClient.SendOrderRequest(ctx, &request)
	if err != nil || res == nil {
		log.Error("fail to send request, err=%s", err)
	} else {
		log.Info("send shell request, msg=%s, errno=%d", res.Msg, res.Errno)
	}
	return
}

/*
2.定单查询 - 查询定单的状态，
	for each state=0 and 存在order id
		去贝壳查询该定单状态
		if 定单不存在
			设置order id='', state=0 // 清除定单，等待重新发送
		else if 定单成功
			设置state=10
		else if 定单失败
			设置state=11
*/
// 最大不要超过100个订单
func (s *Service) checkShellOrder(ctx context.Context, values []*offlineactivity.OfflineActivityResult) (err error) {
	if len(values) == 0 {
		err = fmt.Errorf("values slice 0 length")
		return
	}
	var orderMap = make(map[string]*offlineactivity.OfflineActivityResult, len(values))
	var orderIds []string
	for i, v := range values {
		if i > 99 {
			break
		}
		orderIds = append(orderIds, v.OrderID)
		orderMap[v.OrderID] = v
	}
	var orderIDString = strings.Join(orderIds, ",")
	var orderCheckRequest = shell.OrderCheckRequest{
		Timestamp:     time.Now().UnixNano() / int64(time.Millisecond),
		ThirdOrderNos: orderIDString,
	}
	res, err := s.shellClient.SendCheckOrderRequest(ctx, &orderCheckRequest)
	if err != nil {
		log.Error("fail to check order, order id=%s, err=%s", orderIDString, err)
		return
	}
	// 区分找到的订单和未找到的订单
	// 订单找到，更新状态
	for _, v := range res.Orders {
		log.Info("order find, orderid=%s, current status=%s", v.ThirdOrderNo, v.Status)
		// 删除已找到的订单
		delete(orderMap, v.ThirdOrderNo)
		var resultJSON = shell.OrderCallbackJSON{
			Status:       v.Status,
			ThirdOrderNo: v.ThirdOrderNo,
			Mid:          v.Mid,
		}
		// 更新订单结果
		var orderInfo, e = s.dao.ShellCallbackUpdate(ctx, &resultJSON, "checkorder")
		if e == nil {
			s.queueToUpdateActivityState(orderInfo.ID)
		}
	}
	// 订单未找到,重新发送请求
	var needSendRequest []*offlineactivity.OfflineActivityResult
	for _, v := range orderMap {
		needSendRequest = append(needSendRequest, v)
	}
	err = s.sendRequestAndUpdate(s.dao.OfflineActivityGetDB(), needSendRequest)
	if err != nil {
		log.Error("send order request fail, err=%s", err)
		return
	}
	return
}
func (s *Service) checkActivityState(ctx context.Context, values []int64) (err error) {
	if len(values) == 0 {
		log.Warn("no activity need to check")
		return
	}
	//
	upResult, err := s.dao.OfflineActivityGetUpBonusResultSelect(ctx, "distinct(activity_id) as activity_id", "id in (?)", values)
	for _, v := range upResult {
		var _, e = s.dao.UpdateActivityState(ctx, v.ActivityID)
		if e != nil {
			log.Error("err when update activity state, err=%s", e)
		}
	}
	return
}

func (s *Service) queueToUpdateActivityState(resultID int64) {
	s.chanCheckActivity <- resultID
}

func beginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
