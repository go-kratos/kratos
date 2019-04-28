package dao

import (
	"context"
	"errors"
	"fmt"
	"go-common/app/service/video/stream-mng/common"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

// RawStreamFullInfo 直接从数据库中查询流信息，可传入流名， 也可传入rid
func (d *Dao) RawStreamFullInfo(c context.Context, id int64, sname string) (res *model.StreamFullInfo, err error) {
	var (
		official   []*model.OfficialStream
		backup     []*model.StreamBase
		mainStream *model.MainStream
	)

	if sname != "" {
		official, err = d.GetOfficialStreamByName(c, sname)

		// 可以从原表中查询到
		if err == nil && official != nil && len(official) > 0 {
			id = official[0].RoomID
			goto END
		}

		var backUpInfo *model.BackupStream
		// 原表中查询不到
		backUpInfo, err = d.GetBackupStreamByStreamName(c, sname)
		if err != nil {
			log.Errorv(c, log.KV("log", fmt.Sprintf("sql backup_stream err = %v", err)))
			return
		}

		if backUpInfo == nil {
			err = fmt.Errorf("can not find any info by %s", sname)
			return
		}

		id = backUpInfo.RoomID
	}

END:
	// todo 这里用老的errgroup， 新errgroup2 暂时未有人用,bug未知
	group, errCtx := errgroup.WithContext(c)
	// 如果还未查sv_ls_stream则需要查询
	if id > 0 && len(official) == 0 {
		group.Go(func() (err error) {
			log.Warn("group offical")
			if official, err = d.GetOfficialStreamByRoomID(errCtx, id); err != nil {
				log.Errorv(errCtx, log.KV("log", fmt.Sprintf("group offical err=%v", err)))
			}
			return nil
		})
	}

	if id > 0 {
		group.Go(func() (err error) {
			log.Warn("group main")
			if mainStream, err = d.GetMainStreamFromDB(errCtx, id, ""); err != nil {
				log.Errorv(errCtx, log.KV("log", fmt.Sprintf("group main err=%v", err)))
			}
			return nil
		})

		group.Go(func() (err error) {
			log.Warn("group back")
			back, err := d.GetBackupStreamByRoomID(errCtx, id)
			if err != nil {
				log.Errorv(errCtx, log.KV("log", fmt.Sprintf("group backup err=%v", err)))
			} else {
				backup = d.formatBackup2BaseInfo(c, back)
			}
			return nil
		})
	}

	err = group.Wait()

	if err != nil {
		return
	}

	if len(official) == 0 {
		err = fmt.Errorf("can not find any info by room_id=%d", id)
		return
	}

	return d.formatStreamFullInfo(c, official, backup, mainStream)
}

// RawStreamRIDByName 查询rid
func (d *Dao) RawStreamRIDByName(c context.Context, sname string) (res *model.StreamFullInfo, err error) {
	return d.RawStreamFullInfo(c, 0, sname)
}

// RawMultiStreamInfo 批量查询流信息
func (d *Dao) RawMultiStreamInfo(c context.Context, rids []int64) (res map[int64]*model.StreamFullInfo, err error) {
	var (
		official   []*model.OfficialStream
		backup     []*model.BackupStream
		mainStream []*model.MainStream
	)

	group, errCtx := errgroup.WithContext(c)
	group.Go(func() (err error) {
		if official, err = d.GetMultiOfficalStreamByRID(errCtx, rids); err != nil {
			log.Errorv(errCtx, log.KV("log", fmt.Sprintf("group offical err=%v", err)))
		}
		return nil
	})

	group.Go(func() (err error) {
		if backup, err = d.GetMultiBackupStreamByRID(errCtx, rids); err != nil {
			log.Errorv(errCtx, log.KV("log", fmt.Sprintf("group back err=%v", err)))
		}
		return nil
	})

	group.Go(func() (err error) {
		if mainStream, err = d.GetMultiMainStreamFromDB(errCtx, rids); err != nil {
			log.Errorv(errCtx, log.KV("log", fmt.Sprintf("group back err=%v", err)))
		}
		return nil
	})

	err = group.Wait()

	if err != nil {
		return
	}

	// 把rid相同的放为一组
	ridMapOfficial := map[int64][]*model.OfficialStream{}
	for _, v := range official {
		ridMapOfficial[v.RoomID] = append(ridMapOfficial[v.RoomID], v)
	}

	ridMapBackup := map[int64][]*model.BackupStream{}
	for _, v := range backup {
		ridMapBackup[v.RoomID] = append(ridMapBackup[v.RoomID], v)
	}

	ridMapBackupBase := map[int64][]*model.StreamBase{}
	for id, v := range ridMapBackup {
		ridMapBackupBase[id] = d.formatBackup2BaseInfo(c, v)
	}

	ridMapMain := map[int64]*model.MainStream{}
	for _, v := range mainStream {
		ridMapMain[v.RoomID] = v
	}

	infos := map[int64]*model.StreamFullInfo{}

	flag := false
	for id, v := range ridMapOfficial {
		flag = true
		infos[id], _ = d.formatStreamFullInfo(c, v, ridMapBackupBase[id], ridMapMain[id])
	}

	if flag {
		return infos, nil
	}

	log.Errorv(c, log.KV("log", fmt.Errorf("can not find any info by room_ids=%d", rids)))
	return nil, nil
}

// formatStreamFullInfo 格式化流信息
func (d *Dao) formatStreamFullInfo(c context.Context, official []*model.OfficialStream, backup []*model.StreamBase, main *model.MainStream) (*model.StreamFullInfo, error) {
	resp := &model.StreamFullInfo{}
	resp.List = []*model.StreamBase{}

	var roomID int64

	roomID = official[0].RoomID
	resp.RoomID = official[0].RoomID

	base := &model.StreamBase{}
	base.StreamName = official[0].Name
	base.Type = 1
	base.Key = official[0].Key

	if main != nil {
		base.Options = main.Options
		if 4&base.Options == 4 {
			base.Wmask = true
		}
		if 8&base.Options == 8 {
			base.Mmask = true
		}

	}

	for _, item := range official {
		if item.UpRank == 1 {
			if val, ok := common.SrcMapBitwise[item.Src]; ok {
				// todo origin为main-stream取
				if main != nil {
					base.Origin = main.OriginUpstream
				} else {
					// 做个兜底逻辑， main-stream中没有这个数据，但是sv_ls_stream确实在播
					base.Origin = val
				}
				base.DefaultUpStream = val
			} else {
				// 如果上行不在现在的任意一家， 则重新设置上行
				if err := d.UpdateOfficialStreamStatus(c, roomID, common.BVCSrc); err == nil {
					if main != nil {
						base.Origin = main.OriginUpstream
					} else {
						base.Origin = common.BitWiseBVC
					}
					base.DefaultUpStream = common.BitWiseBVC

					go func(c context.Context, rid int64, fromOrigin int8, toOrigin int64, sname string) {
						d.UpdateStreamStatusCache(c, &model.StreamStatus{
							RoomID:          rid,
							StreamName:      sname,
							DefaultChange:   true,
							DefaultUpStream: toOrigin,
						})
						// 插入日志
						d.InsertChangeLog(c, &model.StreamChangeLog{
							RoomID:      rid,
							FromOrigin:  int64(fromOrigin),
							ToOrigin:    toOrigin,
							Reason:      fmt.Sprintf("上行不在五家CDN,old origin=%d", fromOrigin),
							OperateName: "auto_change",
							Source:      "background",
						})
					}(metadata.WithContext(c), roomID, item.Src, common.BitWiseBVC, item.Name)
				}
			}
		} else if item.UpRank == 2 {
			if val, ok := common.SrcMapBitwise[item.Src]; ok {
				base.Forward = append(base.Forward, val)
			}
		}
	}

	resp.List = append(resp.List, base)

	if len(backup) > 0 {
		for _, v := range backup {
			resp.List = append(resp.List, v)
		}
	}

	d.liveAside.Do(c, func(ctx context.Context) {
		d.diffStreamInfo(ctx, resp, main)
	})

	return resp, nil
}

// formatBackup2Base backup 格式化为base
func (d *Dao) formatBackup2BaseInfo(c context.Context, back []*model.BackupStream) (resp []*model.StreamBase) {
	if len(back) > 0 {
		for _, b := range back {
			bs := &model.StreamBase{}
			bs.StreamName = b.StreamName
			bs.Type = 2
			bs.Key = b.Key

			// 原始上行
			bs.Origin = b.OriginUpstream
			bs.DefaultUpStream = b.DefaultVendor
			bs.Options = b.Options

			// 位运算:可满足9家cdn
			var n int64
			for n = 256; n > 0; n /= 2 {
				if (b.Streaming&n) == n && n != bs.Origin {
					bs.Forward = append(bs.Forward, n)
				}
			}

			resp = append(resp, bs)
		}
	}
	return
}

// 比较新表和老表
func (d *Dao) diffStreamInfo(c context.Context, info *model.StreamFullInfo, mainStream *model.MainStream) {
	if info != nil && info.RoomID != 0 && len(info.List) > 0 {
		if mainStream == nil {
			d.syncMainStream(c, info.RoomID, "")
			log.Infov(c, log.KV("log", fmt.Sprintf("diff_err:can find any info, room_id=%d", info.RoomID)))
			return
		}

		offical := info.List[0]
		if mainStream.StreamName != offical.StreamName {
			log.Infov(c, log.KV("log", fmt.Sprintf("diff_err:stream name is different，room_id=%d", info.RoomID)))
			return
		}

		if mainStream.Key != offical.Key {
			log.Infov(c, log.KV("log", fmt.Sprintf("diff_err:key is different，room_id=%d", info.RoomID)))
			return
		}

		if mainStream.DefaultVendor != offical.DefaultUpStream {
			log.Infov(c, log.KV("log", fmt.Sprintf("diff_err:DefaultVendor is different，room_id=%d,main=%d,offical=%d", info.RoomID, mainStream.DefaultVendor, offical.DefaultUpStream)))
			return
		}

		if mainStream.OriginUpstream != 0 && (mainStream.OriginUpstream != mainStream.DefaultVendor) {
			log.Infov(c, log.KV("log", fmt.Sprintf("diff_err:OriginUpstream is different，room_id=%d, main origin=%d, main default=%d", info.RoomID, mainStream.OriginUpstream, mainStream.DefaultVendor)))
			return
		}

		streaming := offical.DefaultUpStream
		for _, v := range offical.Forward {
			streaming += v
		}
		if mainStream.Streaming != streaming {
			log.Infov(c, log.KV("log", fmt.Sprintf("diff_err:Streaming is different，room_id=%d, main=%d, offical=%d", info.RoomID, mainStream.Streaming, streaming)))
			return
		}
	}
}

func (d *Dao) syncMainStream(c context.Context, roomID int64, streamName string) error {
	if roomID <= 0 && streamName == "" {
		return errors.New("invalid params")
	}

	var err error
	exists, err := d.GetMainStreamFromDB(c, roomID, streamName)
	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Errorv(c, log.KV("log", fmt.Sprintf("sync_stream_data_error = %v", err)))
		return err
	}
	if exists != nil && (exists.RoomID == roomID || exists.StreamName == streamName) {
		return nil
	}

	var full *model.StreamFullInfo
	if roomID > 0 && streamName == "" {
		full, err = d.StreamFullInfo(c, roomID, "")
	} else if roomID <= 0 && streamName != "" {
		full, err = d.StreamFullInfo(c, 0, streamName)
	}

	if err != nil {
		return err
	}
	if full == nil {
		return errors.New("unknow response")
	}

	for _, ss := range full.List {
		if ss.Type == 1 {
			ms := &model.MainStream{
				RoomID:        full.RoomID,
				StreamName:    ss.StreamName,
				Key:           ss.Key,
				DefaultVendor: ss.DefaultUpStream,
				Status:        1,
			}

			if ms.DefaultVendor == 0 {
				ms.DefaultVendor = 1
			}

			_, err := d.CreateNewStream(c, ms)
			if err != nil {
				log.Errorv(c, log.KV("log", fmt.Sprintf("sync_stream_data_error = %v", err)))
			}
			break
		}
	}

	return nil
}
