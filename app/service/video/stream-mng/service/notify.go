package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go-common/app/service/video/stream-mng/common"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"net/http"
	"net/url"
	"time"
)

// StreamingNotify 开关播回调
func (s *Service) StreamingNotify(ctx context.Context, p *model.StreamingNotifyParam, open bool) error {
	// 流鉴权
	if open {
		val := &ValidateParams{
			StreamName: p.StreamName,
			Type:       p.Type,
			Src:        p.SRC,
		}

		_, err := s.CheckStreamValidate(ctx, val, true)
		if err != nil {
			return err
		}
	}

	// 校验
	rid, src, bs, streamInfo, err := s.checkNotifyParams(ctx, p, open)
	if err != nil {
		return err
	}

	if bs == nil { // 如果不是备用流
		// 写库
		if p.Type.String() == "1" {
			if open {
				err = s.dao.SetOriginStreamingStatus(ctx, rid, src, _originUpRankNothing, _originUpRankForwardStreaming)
			} else {
				err = s.dao.SetOriginStreamingStatus(ctx, rid, src, _originUpRankForwardStreaming, _originUpRankNothing)
			}

			if err != nil {
				return err
			}
		}

		// 更新redis
		var forwardVendor int64
		var sname string
		var origin int64
		var postInfo *model.StreamFullInfo
		var options int64
		var newoptions int64
		infoLen := len(streamInfo.List)

		if infoLen > 1 {
			postInfo = streamInfo
		}

		if streamInfo != nil && infoLen > 0 {
			for _, v := range streamInfo.List {
				if v.Type == 1 {
					origin = v.Origin
					sname = v.StreamName
					options = v.Options
					for _, k := range v.Forward {
						forwardVendor += k
					}
					break
				}
			}
		}
		// 获取数据失败或者其他情况，直接删除缓存
		if sname == "" {
			s.dao.DeleteStreamByRIDFromCache(ctx, rid)
		} else {
			// 过渡接口 后续和main-stream保存一致
			vendor := common.SrcMapBitwise[int8(src)]
			if open {
				//检查options第二位是否是1，是的话通知AI
				newoptions = options
				if 2&options == 2 {
					//-------------->此处通知AI<---------------
				}

				// 主推
				if p.Type.String() == "0" {
					s.dao.UpdateStreamStatusCache(ctx, &model.StreamStatus{
						RoomID:       rid,
						StreamName:   sname,
						OriginChange: true,
						Origin:       vendor,
					})

					// 解决并发问题
					go func(ctx context.Context, rid int64) {
						time.Sleep(time.Second * 30)
						s.dao.DeleteStreamByRIDFromCache(ctx, rid)
					}(metadata.WithContext(ctx), rid)

					if postInfo != nil {
						for _, v := range postInfo.List {
							if v.StreamName == sname {
								v.Origin = vendor
								break
							}
						}
					}
				} else {
					s.dao.UpdateStreamStatusCache(ctx, &model.StreamStatus{
						RoomID:        rid,
						StreamName:    sname,
						Forward:       (forwardVendor | vendor),
						ForwardChange: true,
					})

					if postInfo != nil {
						for _, v := range postInfo.List {
							if v.StreamName == sname {
								v.Forward = append(v.Forward, vendor)
								break
							}
						}
					}
				}

			} else {
				//关播回调 wmask mmask 第三位第四位都要清零 12 = 00001100
				newoptions = options &^ 4
				newoptions = newoptions &^ 8
				if p.Type.String() == "0" {
					// 关播需要判断是否为当前的流才可以更新
					if origin == common.SrcMapBitwise[int8(src)] {
						s.dao.UpdateStreamStatusCache(ctx, &model.StreamStatus{
							RoomID:        rid,
							StreamName:    sname,
							Origin:        0,
							OriginChange:  true,
							Forward:       0,
							ForwardChange: true,
							Options:       newoptions,
							OptionsChange: true,
						})

						if postInfo != nil {
							for _, v := range postInfo.List {
								if v.StreamName == sname {
									v.Origin = 0
									v.Forward = []int64{}
									break
								}
							}
						}
					}
				} else {
					s.dao.UpdateStreamStatusCache(ctx, &model.StreamStatus{
						RoomID:        rid,
						StreamName:    sname,
						Forward:       (forwardVendor &^ vendor),
						ForwardChange: true,
					})

					if postInfo != nil {
						for _, v := range postInfo.List {
							if v.StreamName == sname {
								// 去除
								forwards := []int64{}
								for _, f := range v.Forward {
									if f != vendor {
										forwards = append(forwards, f)
									}
								}
								v.Forward = forwards
								break
							}
						}
					}
				}
			}
		}

		//log.Warn("%v", postInfo)

		// 同步数据
		go func(ctx context.Context, roomID int64, legacySrc int8, isOpen bool, isOrigin bool, postInfo *model.StreamFullInfo, options int64, newoptions int64) {
			if postInfo != nil {
				s.updateLiveUpStream(ctx, postInfo)
			}

			s.syncMainStream(ctx, roomID, "")
			if vendor, ok := common.SrcMapBitwise[legacySrc]; ok {
				s.dao.MainStreamNotify(ctx, roomID, vendor, isOpen, isOrigin, options, newoptions)
			}
		}(metadata.WithContext(ctx), rid, int8(src), open, p.Type.String() == "0", postInfo, options, newoptions)
	} else { // 如果是备用流
		bs, err = s.dao.SetBackupStreamStreamingStatus(ctx, p, bs, open)
		if err != nil {
			return err
		}

		// 开播&主推
		if open && p.Type.String() == "0" {
			// 解决并发问题
			go func(ctx context.Context, rid int64) {
				time.Sleep(time.Second * 30)
				s.dao.DeleteStreamByRIDFromCache(ctx, rid)
			}(metadata.WithContext(ctx), rid)
		}

		// 已经封装好的bs
		s.dao.UpdateStreamStatusCache(ctx, &model.StreamStatus{
			RoomID:        bs.RoomID,
			StreamName:    bs.StreamName,
			ForwardChange: true,
			Forward:       bs.Streaming &^ bs.OriginUpstream,
			OriginChange:  true,
			Origin:        bs.OriginUpstream,
		})

		m, err := s.GetStreamInfoByRIDMapSrcFromDB(ctx, bs.RoomID)
		if err == nil {
			go s.updateLiveUpStream(metadata.WithContext(ctx), m)
		}
	}

	// 广播（如果有), 无论主流还是备用流均发广播
	if open && p.Type.String() == "0" {
		body := []byte(fmt.Sprintf(`{"cmd":"LIVE", "roomid":"%d"}`, rid))
		log.Info("%+v", string(body))

		q := make(url.Values)
		q.Set("ensure", "0")
		q.Set("cid", fmt.Sprintf("%d", rid))
		err := s.dao.NewRequst(ctx, http.MethodPost, "http://live-dm.bilibili.co/dm/1/push", q, body, nil, nil)
		log.Infov(ctx, log.KV("body", string(body)), log.KV("type", "notify DM"))
		if err != nil {
			log.Errorv(ctx, log.KV("err", err), log.KV("type", "notify DM"))
		}
	}
	return nil
}

func (s *Service) checkNotifyParams(ctx context.Context, p *model.StreamingNotifyParam, open bool) (int64, int, *model.BackupStream, *model.StreamFullInfo, error) {
	ts, _ := p.TS.Int64()
	if p == nil || p.StreamName == "" || p.SRC == "" || ts == 0 {
		return 0, 0, nil, nil, errors.New("invalid params")
	}

	t, err := p.Type.Int64()
	if err != nil {
		return 0, 0, nil, nil, errors.New("invalid typeof type ")
	}

	if time.Now().Sub(time.Unix(ts, 0)) > time.Minute*30 {
		return 0, 0, nil, nil, errors.New("ts expired")
	}
	// sign
	if salt, ok := CDNSalt[p.SRC]; ok {
		uri := "/live_stream/v1/StreamThird/close_notify"
		if open {
			uri = "/live_stream/v1/StreamThird/open_notify"
		}

		h := md5.New()
		h.Write([]byte(fmt.Sprintf("%s%s%s", salt, uri, p.TS.String())))

		log.Warn("sign = %v", hex.EncodeToString(h.Sum(nil)))
		if p.Sign != hex.EncodeToString(h.Sum(nil)) {
			return 0, 0, nil, nil, errors.New("invalid sign")
		}
	} else {
		return 0, 0, nil, nil, errors.New("invalid src")
	}

	if bitwise, ok := common.CdnBitwiseMap[p.SRC]; ok {
		log.Infov(ctx, log.KV("bitwise", bitwise), log.KV("t", t))
		info, _ := s.dao.StreamFullInfo(ctx, 0, p.StreamName)
		if info != nil && info.List != nil {
			for _, row := range info.List {
				if row.StreamName != p.StreamName {
					continue
				}

				// 只有主流需要校验推流逻辑,
				if row.Type == 1 && ((t == 0 && row.DefaultUpStream == bitwise) || (t == 1 && row.DefaultUpStream != bitwise)) {
					return info.RoomID, int(common.BitwiseMapSrc[bitwise]), nil, info, nil
				}

				if row.Type == 2 {
					bs, err := s.dao.GetBackupStreamByStreamName(ctx, p.StreamName)
					if err != nil {
						log.Infov(ctx, log.KV("query_backup_stream_with_err", err))
						return 0, 0, nil, nil, errors.New("invalid stream name")
					}
					if bs != nil && bs.RoomID != 0 {
						return bs.RoomID, 0, bs, nil, nil
					}
				}
			}
		}
	}
	return 0, 0, nil, nil, errors.New("invalid type")
}

// updateLiveUpStream 更新playurl缓存
func (s *Service) updateLiveUpStream(ctx context.Context, m *model.StreamFullInfo) {
	b, err := json.Marshal(m)
	if err == nil {
		h := map[string]string{"Content-Type": "application/json"}
		log.Infov(ctx, log.KV("body", string(b)), log.KV("type", "notify playurl"))
		err := s.dao.NewRequst(ctx, http.MethodPost, "http://live-upstream.bilibili.co/live_stream/v1/Dispatch/set_streaminfo", nil, b, h, nil)
		if err != nil {
			log.Errorv(ctx, log.KV("err", err), log.KV("type", "notify playurl"))
		}
	}
}
