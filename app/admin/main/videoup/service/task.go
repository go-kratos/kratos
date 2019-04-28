package service

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/manager"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

// TaskTooksByHalfHour get task books by ctime
func (s *Service) TaskTooksByHalfHour(c context.Context, stime, etime time.Time) (tooks []*archive.TaskTook, err error) {
	if tooks, err = s.arc.TaskTooksByHalfHour(c, stime, etime); err != nil {
		log.Error("s.arc.TaskTooksByHalfHour(%v,%v)", stime, etime)
		return
	}
	return
}

//lockVideo Lock specified category videos
func (s *Service) lockVideo() {
	//TODO It's a temporary function to lock videos. When no longer needed, remove it.
	var (
		c              = context.TODO()
		adminID  int64 = 399          //Temporary task admin id for lock video.
		uname          = "videoupjob" //Temporary task admin name for lock video.
		reason         = "版权原因，该视频不予审核通过"
		reasonID int64 = 197
		rCateID  int64 = 76
		tagID    int64 = 7 //版权tag
		note           = "自动锁定分区视频【欧美电影】，【日本电影】，【其他国家】，【港台剧】，【海外剧】"
		ctime          = time.Now()
		mtime          = ctime
		err      error
		tx       *xsql.Tx
	)

	//Note: check if another instance is locking video
	locking, err := s.arc.IsLockingVideo(c)
	if err != nil {
		log.Error("s.lockVideo() s.arc.IsLockingVideo() err(%v)", err)
		return
	}
	if locking {
		log.Info("s.lockVideo() another instance is locking video")
		return
	}
	//Set locking video redis
	if err = s.arc.LockingVideo(c, 1); err != nil {
		log.Error("s.lockVideo() s.arc.LockingVideo() err(%v)", err)
		return
	}
	defer func() {
		//Unlock locking video redis
		if err = s.arc.LockingVideo(c, 0); err != nil {
			log.Error("s.lockVideo() s.arc.LockingVideo() err(%v)", err)
		}
	}()
	if _, err = s.arc.TaskUserCheckIn(c, adminID); err != nil {
		log.Error("s.lockVideo() s.arc.TaskUserCheckIn(%d) error(%v)", adminID, err)
		return
	}
	tasks, err := s.arc.UserUndoneSpecTask(c, adminID)
	if err != nil {
		log.Error("s.lockVideo() error(%v)", err)
		return
	}
	if len(tasks) == 0 {
		log.Info("s.lockVideo() no task.")
		return
	}

	var vps = []*archive.VideoParam{}
	for _, t := range tasks {
		if t.State == archive.TypeFinished {
			continue
		}
		v, err := s.arc.VideoByCID(c, t.Cid)
		if err != nil {
			log.Error("s.lockVideo() s.arc.VideoByCID(%d) error(%v)", t.Cid, err)
			continue
		}
		arc, err := s.arc.Archive(c, t.Aid)
		if err != nil {
			log.Error("s.lockVideo() s.arc.Archive(%d) error(%v)", v.Aid, err)
			continue
		}
		//If archive's mid in white list, release task and continue
		if s.PGCWhite(arc.Mid) {
			log.Info("s.lockVideo() mid in white list, release task(%d)", t.ID)
			//Begin update task state and add task history
			if tx, err = s.arc.BeginTran(c); err != nil {
				log.Error("s.arc.BeginTran error(%v)", err)
				continue
			}

			if _, err = s.arc.TxReleaseByID(tx, t.ID); err != nil {
				log.Error("s.lockVideo() s.arc.TxReleaseByID(%d) error(%v)", t.ID, err)
				tx.Rollback()
				continue
			}
			if _, err = s.arc.TxAddTaskHis(tx, 0, archive.ActionRelease /*action*/, t.ID /*task_id*/, t.Cid /*cid*/, 0, 0, 0, "lockVideo release" /*reason*/); err != nil {
				log.Error("s.lockVideo() s.arc.TxAddTaskHis error(%v)", err)
				tx.Rollback()
				continue
			}

			if err = tx.Commit(); err != nil {
				log.Error("tx.Commit error(%v)", err)
			}
			continue
		}
		//Get archive's top type id
		rTp, err := s.TypeTopParent(arc.TypeID)
		if err != nil {
			log.Error("s.lockVideo() s.arc.TypeTopParent(%d) error(%v)", arc.TypeID, err)
			continue
		}
		//Add video lock reason log
		if _, err = s.mng.AddReasonLog(c, v.Cid, manager.ReasonLogTypeVideo, rCateID, reasonID, adminID, arc.TypeID, ctime, mtime); err != nil {
			log.Error("s.lockVideo() s.arc.AddReasonLog(%d,%d,%d,%d,%d,%d,%v,%v) error(%v)", v.Cid, manager.ReasonLogTypeVideo, rCateID, reasonID, adminID, arc.TypeID, ctime, mtime, err)
		}
		//Begin update task state and add task history
		if tx, err = s.arc.BeginTran(c); err != nil {
			log.Error("s.arc.BeginTran error(%v)", err)
			continue
		}

		if _, err = s.arc.TxUpTaskByID(tx, t.ID, map[string]interface{}{"state": archive.TypeFinished, "utime": 0}); err != nil {
			log.Error("s.lockVideo() s.arc.TxUpTaskByID(%d) error(%v)", t.ID, err)
			tx.Rollback()
			continue
		}
		if _, err = s.arc.TxAddTaskHis(tx, archive.PoolForFirst, archive.TypeFinished, t.ID, t.Cid, adminID, 0, archive.VideoStatusLock, reason); err != nil {
			log.Error("s.lockVideo() s.arc.SubmitTask(%d) error(%v)", t.ID, err)
			tx.Rollback()
			continue
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit error(%v)", err)
			continue
		}
		//Set video param
		vp := &archive.VideoParam{}
		vp.ID = v.ID
		vp.Aid = v.Aid
		vp.Cid = v.Cid
		vp.Filename = v.Filename
		vp.RegionID = rTp.ID
		vp.Status = archive.VideoStatusLock
		vp.UID = adminID
		vp.Oname = uname
		vp.Note = note + " [任务ID]task:" + strconv.Itoa(int(t.ID))
		vp.Reason = reason
		vp.TagID = tagID
		vp.Encoding = 0
		vps = append(vps, vp)
		log.Info("s.lockVideo() add video. cid(%d)", v.Cid)
	}
	if len(vps) == 0 {
		log.Info("s.lockVideo() no belongs to 399 task.")
		return
	}
	//Add videos to batch list
	s.BatchVideo(c, vps, archive.ActionVideoSubmit)
}

/* 批量查询，批量转换
 * list 	[]*struct{}
 * multrans 转化器，根据ID查出其他值
 * ID   	id字段名称，id字段类型必须是int64
 * Names 	查出来的各个字段名称
 */
func (s *Service) mulIDtoName(c context.Context, list interface{}, multrans func(context.Context, []int64) (map[int64][]interface{}, error), ID string, Names ...string) (err error) {
	var (
		lV, itemI, itemIE, idFiled, nameFiled, valueField reflect.Value
		id                                                int64
		ids                                               []int64
		hashIDName                                        = make(map[int64][]interface{})
	)

	if lV = reflect.ValueOf(list); !lV.IsValid() || lV.IsNil() || lV.Kind() != reflect.Slice {
		return fmt.Errorf("invalid list")
	}

	count := lV.Len()
	for i := 0; i < count; i++ {
		if itemI = lV.Index(i); !itemI.IsValid() || itemI.IsNil() || itemI.Kind() != reflect.Ptr {
			return fmt.Errorf("invalid itemI")
		}
		if itemIE = itemI.Elem(); !itemIE.IsValid() || itemIE.Kind() != reflect.Struct {
			return fmt.Errorf("invalid itemIE")
		}
		if idFiled = itemIE.FieldByName(ID); !idFiled.IsValid() || idFiled.Kind() != reflect.Int64 {
			return fmt.Errorf("invalid idFiled")
		}
		for _, name := range Names {
			if nameFiled = itemIE.FieldByName(name); !nameFiled.IsValid() || !nameFiled.CanSet() {
				return fmt.Errorf("invalid nameFiled")
			}
		}
		if id = idFiled.Int(); id != 0 {
			if _, ok := hashIDName[id]; !ok {
				hashIDName[id] = []interface{}{}
				ids = append(ids, id)
			}
		}
	}
	if hashIDName, err = multrans(c, ids); err != nil {
		log.Error("multrans error(%v)", ids)
		return
	}
	for i := 0; i < count; i++ {
		itemIE = lV.Index(i).Elem()
		id = itemIE.FieldByName(ID).Int()
		if names, ok := hashIDName[id]; ok && len(names) == len(Names) {
			for i, name := range names {
				nameFiled = itemIE.FieldByName(Names[i])
				valueField = reflect.ValueOf(name)
				if nameFiled.Kind() != valueField.Kind() {
					log.Error("multrans return %v while need ", ids)
					continue
				}
				itemIE.FieldByName(Names[i]).Set(reflect.ValueOf(name))
			}
		}
	}
	return
}

// 每个ID单独查询 strict严格模式下一次错误，直接返回
func (s *Service) singleIDtoName(c context.Context, list interface{}, singletrans func(context.Context, int64) ([]interface{}, error), strict bool, ID string, Names ...string) (err error) {
	var (
		lV, itemI, itemIE, idFiled, nameFiled, valueField reflect.Value
		id                                                int64
		values                                            []interface{}
	)

	if lV = reflect.ValueOf(list); !lV.IsValid() || lV.IsNil() || lV.Kind() != reflect.Slice {
		return fmt.Errorf("invalid list")
	}

	count := lV.Len()
	for i := 0; i < count; i++ {
		if itemI = lV.Index(i); !itemI.IsValid() || itemI.IsNil() || itemI.Kind() != reflect.Ptr {
			return fmt.Errorf("invalid itemI")
		}
		if itemIE = itemI.Elem(); !itemIE.IsValid() || itemIE.Kind() != reflect.Struct {
			return fmt.Errorf("invalid itemIE")
		}
		if idFiled = itemIE.FieldByName(ID); !idFiled.IsValid() || idFiled.Kind() != reflect.Int64 {
			return fmt.Errorf("invalid idFiled")
		}
		for _, Name := range Names {
			if nameFiled = itemIE.FieldByName(Name); !nameFiled.IsValid() || !nameFiled.CanSet() {
				return fmt.Errorf("invalid nameFiled")
			}
		}

		if id = idFiled.Int(); id != 0 {
			if values, err = singletrans(c, id); err != nil || len(values) != len(Names) {
				log.Error("s.sigleIDtoName error(%v) len(values)=%d len(Names)=%d", err, len(values), len(Names))
				if strict {
					return
				}
				err = nil
				continue
			}
			for i, value := range values {
				nameFiled = itemIE.FieldByName(Names[i])
				valueField = reflect.ValueOf(value)
				if nameFiled.Kind() != valueField.Kind() {
					log.Error("singletrans return %s while need %s", valueField.Kind().String(), nameFiled.Kind().String())
					continue
				}
				nameFiled.Set(valueField)
			}
		}
	}
	return
}

// GetUID 获取uid，有时候cookie没有uid
func (s *Service) GetUID(c context.Context, name string) (uid int64, err error) {
	return s.mng.GetUIDByName(c, name)
}
