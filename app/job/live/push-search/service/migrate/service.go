package migrate

import (
	"context"
	defaultsql "database/sql"
	"fmt"
	"go-common/app/job/live/push-search/conf"
	"go-common/app/job/live/push-search/dao/migrate"
	"go-common/app/job/live/push-search/model"
	"go-common/library/database/sql"
	"strconv"
	"sync"
	"time"
)

const _sql = "select roomid, short_id,uid,uname,area,title,tags, mtime,a.ctime,try_time,user_cover,a.lock_status,hidden_status,attentions,live_time,area_v2_id,area_v2_parent_id,b.name as area_v2_name, virtual,round_status,on_flag,online, cover from ap_room a left join ap_room_area_v2 b on a.area_v2_id=b.id where roomid > %d order by roomid asc limit 100 "

const hbaseTable   = "live:PushSearch"
const hbaseFamily  = "search"

type message struct {
	rowKey string
	values map[string]map[string][]byte
}

type MService struct {
	c                  *conf.Config
	dao                *migrate.Dao
	hChan			   []chan *message
	waiterChan			*sync.WaitGroup
	mainWaiter			*sync.WaitGroup
}


func NewMigrateS(c *conf.Config) (s *MService) {
	s = &MService{
		c:   c,
		dao: migrate.NewMigrate(c),
		hChan: make([]chan *message, c.MigrateNum),
		waiterChan:  new(sync.WaitGroup),
		mainWaiter:  new(sync.WaitGroup),
	}

	//ap room 表 binlog qps 高, hash roomId 并行
	for i := 0; i < c.MigrateNum; i++ {
		ch := make(chan *message, 1024)
		s.hChan[i] = ch
		go s.handle(ch)
	}

	return s
}

func (ms *MService) Migrate (roomid string, isTest string) {

	id, err := strconv.Atoi(roomid)
	if err != nil {
		fmt.Println("roomid error")
	}
	ms.mainWaiter.Add(1)
	defer ms.mainWaiter.Done()
	var rows *sql.Rows
	online := &defaultsql.NullInt64{}
	cover := &defaultsql.NullString{}
	areaV2Name := &defaultsql.NullString{}
	for {
		rows, err := ms.dao.RoomDb.Query(context.TODO(), fmt.Sprintf(_sql, id))
		if err != nil {
			fmt.Println("query error:%+v", err)
			return
		}
		for rows.Next() {
			r := new(model.TableField)
			if err = rows.Scan(&r.RoomId, &r.ShortId, &r.Uid, &r.UName, &r.Area, &r.Title, &r.Tag, &r.MTime, &r.CTime, &r.TryTime, &r.UserCover, &r.LockStatus, &r.HiddenStatus, &r.Attentions, &r.LiveTime, &r.AreaV2Id, &r.AreaV2ParentId, areaV2Name, &r.Virtual, &r.RoundStatus, &r.OnFlag, online, cover); err != nil {
				if !online.Valid {
					r.Online = 0
				}
				if !cover.Valid {
					r.Cover = ""
				}
			}
			r.AreaV2Name = areaV2Name.String
			r.Online = int(online.Int64)
			r.Cover = cover.String

			zijie := ms.generateSearchInfo(r)
			if r.LiveTime != "0000-00-00 00:00:00" {
				fmt.Println(r.RoomId, "jump live room")
				continue
			}
			values := map[string]map[string][]byte{hbaseFamily: zijie}
			rowKey := ms.rowKey(r.RoomId)
			m := &message{
				rowKey: rowKey,
				values: values,
			}
			ms.hChan[r.RoomId % ms.c.MigrateNum] <- m
			fmt.Println(r.RoomId)
			if isTest == "1" {
				return
			}
			id = r.RoomId
		}
	}
	rows.Close()
}

func (ms *MService) handle(c chan *message) {
	ms.waiterChan.Add(1)
	defer ms.waiterChan.Done()
	for {
		msgData, ok := <-c
		if !ok {
			fmt.Println("close chan")
			return
		}
		ms.dao.SearchHBase.PutStr(context.TODO(), hbaseTable, msgData.rowKey, msgData.values)
	}
}

func (ms *MService) Close() {
	ms.dao.Close()
	ms.mainWaiter.Wait()
	for _, ch := range ms.hChan {
		close(ch)
	}
	ms.waiterChan.Wait()
	ms.dao.SearchHBase.Close()
}
func (ms *MService) rowKey(roomId int) string{
	key := fmt.Sprintf("%d_%d", roomId % 10, roomId)
	return key
}
func (ms *MService) generateSearchInfo(new *model.TableField) (retByte map[string][]byte){
	newByteMap := make(map[string][]byte)
	newByteMap["id"] = []byte(strconv.Itoa(new.RoomId))
	newByteMap["short_id"] = []byte(strconv.Itoa(new.ShortId))
	newByteMap["uid"] = []byte(strconv.FormatInt(new.Uid, 10))
	newByteMap["uname"] = []byte(new.UName)
	newByteMap["category"] = []byte(strconv.Itoa(new.Area))
	newByteMap["title"] = []byte(new.Title)
	newByteMap["tag"] = []byte(new.Tag)

	tryTime, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", new.TryTime, time.Local)
	tryTimeStr := tryTime.Format("2006-01-02 15:04:05")
	if tryTimeStr == "0001-01-01 00:00:00" {
		tryTimeStr = "0000-00-00 00:00:00"
		new.TryTime = tryTimeStr
	}
	newByteMap["try_time"] = []byte(tryTimeStr)

	newByteMap["cover"] = []byte(new.Cover)
	newByteMap["user_cover"] = []byte(new.UserCover)

	lockStatus, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", new.LockStatus, time.Local)
	lockStatusStr := lockStatus.Format("2006-01-02 15:04:05")
	if lockStatusStr == "0001-01-01 00:00:00" {
		lockStatusStr = "0000-00-00 00:00:00"
		new.LockStatus = lockStatusStr
	}
	newByteMap["lock_status"] = []byte(strconv.Itoa(ms.getLockStatus(lockStatusStr)))

	hiddenStatus, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", new.HiddenStatus, time.Local)
	hiddenStatusStr := hiddenStatus.Format("2006-01-02 15:04:05")
	if hiddenStatusStr == "0001-01-01 00:00:00" {
		hiddenStatusStr = "0000-00-00 00:00:00"
		new.HiddenStatus = hiddenStatusStr
	}
	newByteMap["hidden_status"] = []byte(strconv.Itoa(ms.getHiddenStatus(hiddenStatusStr)))

	newByteMap["attentions"] = []byte(strconv.Itoa(new.Attentions))
	newByteMap["attention"] = []byte(strconv.Itoa(new.Attentions))
	newByteMap["online"] = []byte(strconv.Itoa(new.Online))

	liveTime, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", new.LiveTime, time.Local)
	liveTimeStr := liveTime.Format("2006-01-02 15:04:05")
	if liveTimeStr == "0001-01-01 00:00:00" {
		liveTimeStr = "0000-00-00 00:00:00"
		new.LiveTime = liveTimeStr
	}
	newByteMap["live_time"] = []byte(liveTimeStr)

	newByteMap["area_v2_id"] = []byte(strconv.Itoa(new.AreaV2Id))
	newByteMap["ord"] = []byte(strconv.Itoa(new.AreaV2ParentId))
	newByteMap["arcrank"] = []byte(strconv.Itoa(new.Virtual))

	cTime, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", new.CTime, time.Local)
	cTimeStr := cTime.Format("2006-01-02 15:04:05")
	if cTimeStr == "0001-01-01 00:00:00" {
		new.CTime = "0000-00-00 00:00:00"
	}else{
		new.CTime = cTimeStr
	}
	mTime, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", new.MTime, time.Local)
	mTimeStr := mTime.Format("2006-01-02 15:04:05")
	if mTimeStr == "0001-01-01 00:00:00" {
		new.MTime = "0000-00-00 00:00:00"
	}else{
		new.MTime = mTimeStr
	}
	newByteMap["lastupdate"] = []byte(ms.getLastUpdate(new))
	newByteMap["is_live"] = []byte(strconv.Itoa(ms.getLiveStatus(new)))
	newByteMap["s_category"] = []byte(new.AreaV2Name)
	return newByteMap
}

//获取直播状态
func (ms *MService) getLiveStatus(roomInfo *model.TableField) int{
	if roomInfo.LiveTime != "0000-00-00 00:00:00" {
		return 1
	}

	if roomInfo.RoundStatus == 1 && roomInfo.OnFlag == 1{
		return 2
	}

	return 0
}

//获取房间最后更新时间
func (ms *MService) getLastUpdate(roomInfo *model.TableField) string{
	if roomInfo.MTime != "0000-00-00 00:00:00" {
		return roomInfo.MTime
	}
	return roomInfo.CTime
}

func (ms *MService) getLockStatus(lockStatus string) int{
	status := 0
	if lockStatus != "0000-00-00 00:00:00" {
		status = 1
	}
	return status
}

func (ms *MService) getHiddenStatus(HiddenStatus string) int{
	status := 0
	if HiddenStatus != "0000-00-00 00:00:00" {
		status = 1
	}
	return status
}

