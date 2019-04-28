package service

import (
	"context"
	"errors"
	"fmt"
	"go-common/library/log"
)

func (s *Service) ChangeMaskStreamByRoomID(ctx context.Context, realRoomID int64, streamname string, mask int64) (t string, err error) {
	//记录一下
	log.Errorv(ctx, log.KV("log", fmt.Sprintf("mask is run = rid->%v sname->%v mask->%v", realRoomID, streamname, mask)))
	retrylimit := 2
retry:
	infos, err := s.dao.StreamFullInfo(ctx, realRoomID, streamname)
	//没有查到结果
	if err != nil {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("select_maskbyroomid_error = %v", err)))
		return "", err
	}
	var newmask int64
	if infos != nil && infos.List != nil {
		//当rid = infos.RoomID
		realRoomID = infos.RoomID
		for _, v := range infos.List {
			if v.Type == 1 {
				//做位运算
				if mask == 0 {
					//mask=0关上蒙版
					newmask = v.Options &^ 2
					newmask = newmask &^ 4
					newmask = newmask &^ 8
				} else if mask == 1 {
					//mask==1打开蒙版
					newmask = v.Options | 2
				} else if mask == 2 {
					//wmask蒙版可调度playurl
					newmask = v.Options | 4
				} else if mask == 3 {
					//wmask蒙版停止调度playurl
					newmask = v.Options &^ 4
				} else if mask == 4 {
					//mmask蒙版可调度playurl
					newmask = v.Options | 8
				} else if mask == 5 {
					//mmask蒙版停止调度playurl
					newmask = v.Options &^ 8
				} else {
					return "", errors.New("mask value is err")
				}

				//修改数据库 options字段
				err := s.dao.ChangeMainStreamOptions(ctx, realRoomID, newmask, v.Options)
				retrylimit--
				if err != nil {
					if retrylimit <= 0 {
						//更新数据库重试2次
						goto retry
					}
					return "", err
				} else {
					//更新缓存
					err := s.dao.UpdateRoomOptionsCache(ctx, realRoomID, v.StreamName, newmask)
					//重试一次，如果缓存更新失败会导致缓存存在期间数据库插入失败
					if err != nil {
						err := s.dao.UpdateRoomOptionsCache(ctx, realRoomID, v.StreamName, newmask)
						if err != nil {
							//再失败则删除缓存
							s.dao.DeleteStreamByRIDFromCache(ctx, realRoomID)
						}
					}
					return "success", nil
				}
			}
		}
	}
	return "", errors.New("select mask options by roomid result is nil")
}
