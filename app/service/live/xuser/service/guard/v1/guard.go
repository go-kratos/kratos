package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/live/xuser/dao/account"
	xhttp "net/http"
	"strconv"
	"strings"

	roomApi "go-common/app/service/live/room/api/liverpc"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	v1pb "go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/app/service/live/xuser/conf"
	"go-common/app/service/live/xuser/dao/guard"
	"go-common/app/service/live/xuser/model"
	// account "go-common/app/service/main/account/rpc/client"
	"go-common/library/cache"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

const (
	_entryEffectBusinessId = 1
)

// GuardService guard service
type GuardService struct {
	c                  *conf.Config
	dao                *guard.GuardDao
	cli                *bm.Client
	liveEntryEffectPub *databus.Databus
	roomAPI            *roomApi.Client
	accRPC             *account.Dao
	async              *cache.Cache
	asyncMulti         *cache.Cache
}

// New new guard service
func New(c *conf.Config) *GuardService {
	return &GuardService{
		c:                  c,
		dao:                guard.NewGuardDao(c),
		cli:                bm.NewClient(c.BMClient),
		roomAPI:            roomApi.New(c.LiveRPC["room"]),
		accRPC:             account.New(c),
		liveEntryEffectPub: databus.New(c.LiveEntryEffectPub),
		async:              cache.New(1, 1024),
		asyncMulti:         cache.New(2, 10240),
	}
}

// Buy buy guard
// grpc wrapper
func (s *GuardService) Buy(ctx context.Context, req *v1pb.GuardBuyReq) (reply *v1pb.GuardBuyReply, err error) {
	var status int
	reqDao := &model.GuardBuy{
		OrderId:    req.OrderId,
		Uid:        req.Uid,
		Ruid:       req.Ruid,
		GuardLevel: req.GuardLevel,
		Num:        req.Num,
		Platform:   req.Platform,
		Source:     req.Source,
	}
	status, err = s.saveGuard(ctx, reqDao)
	if err != nil {
		log.Error("[service.v1.guard|Buy] buy guard error(%+v), params(%+v)", err, reqDao)
	}
	reply = &v1pb.GuardBuyReply{
		Status: status,
	}
	return
}

func (s *GuardService) saveGuard(ctx context.Context, req *model.GuardBuy) (status int, err error) {
	ok, err := s.dao.LockOrder(ctx, req.OrderId)
	if !ok {
		log.Info("[service.v1.guard|saveGuard] LockOrder success, req(%+v)", req)
		return 1, nil
	}
	if err != nil {
		log.Warn("[service.v1.guard|saveGuard] LockOrder failed(%+v), req(%+v)", err, req)
		return 2, err
	}

	status = 1

	info, err := s.dao.GetGuardByUIDRuid(ctx, req.Uid, req.Ruid)
	if err != nil {
		status = 2
		return
	}

	s.dao.ClearCache(ctx, req.Uid, req.Ruid)

	buyType := 1 // 1: 开通大航海广播, 2: 续费大航海广播
	if len(info) == 0 {
		// 无守护
		err = s.dao.AddGuard(ctx, req)
		if err != nil {
			log.Warn("[service.v1.guard|saveGuard] AddGuard failed(%+v), req(%+v)", err, req)
			status = 2
			goto END
		}
	} else {
		if info[0].PrivilegeType == req.GuardLevel {
			buyType = 2
			// 购买的等级没有当前佩戴的守护等级高
			s.dao.UpdateGuard(ctx, req, ">=") // @CHECKED
		} else if info[0].PrivilegeType < req.GuardLevel {
			hasEqual := 0
			// 越前面的等级越高
			expiredTime := info[0].ExpiredTime
			// 购买的新守护，只影响等级比(购买的守护等级)更低的守护的过期时间
			for _, v := range info {
				if v.PrivilegeType == req.GuardLevel {
					buyType = 2
					hasEqual = 1
				}
				// 在高级别的最晚有效期上处理
				if v.PrivilegeType < req.GuardLevel && v.ExpiredTime.Sub(expiredTime) > 0 {
					expiredTime = v.ExpiredTime
				}
			}

			if hasEqual == 0 {
				// 新购买, 过期时间在当前佩戴的大航海的过期时间基础上累加
				s.dao.UpsertGuard(ctx, req, expiredTime.AddDate(0, 0, req.Num*30).Format("2006-01-02 15:04:05")) // @CHECKED insert or update
				s.dao.UpdateGuard(ctx, req, ">")                                                                 // @CHECKED
			} else {
				s.dao.UpdateGuard(ctx, req, ">=") // @CHECKED
			}
		} else if info[0].PrivilegeType > req.GuardLevel {
			// 购买的等级比当前拥有的守护等级高
			err = s.dao.AddGuard(ctx, req)
			if err != nil {
				log.Warn("[service.v1.guard|saveGuard] AddGuard failed(%+v), req(%+v)", err, req)
				status = 2
				goto END
			}
			// 更新低级别的过期时间
			s.dao.UpdateGuard(ctx, req, ">") // @CHECKED
		}
	}

END:
	if status == 1 {
		// 发送房间广播
		s.async.Save(func() {
			s.sendRoomMsg(context.Background(), req.Uid, req.Ruid, req.GuardLevel, buyType)
		})
	} else {
		s.dao.UnlockOrder(ctx, req.OrderId)
		return
	}

	s.dao.ClearCache(ctx, req.Uid, req.Ruid)
	// 投递进场特效给rewardcenter落库
	s.async.Save(func() {
		s.sendEntryEffect(context.Background(), req.Uid, req.Ruid)
	})
	return
}

func getEffectIdByLevel(level int) int {
	switch level {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 4
	}
	return 0
}

func (s *GuardService) sendEntryEffect(ctx context.Context, uid int64, ruid int64) (err error) {
	if !s.c.GuardCfg.OpenEntryEffectDatabus {
		return
	}

	info, err := s.dao.GetGuardByUIDRuid(ctx, uid, ruid)
	if err != nil {
		log.Error("[service.guard.v1.guard|sendEntryEffect] GetGuardByUIDRuid error(%v), uid(%d), ruid(%d)", err, uid, ruid)
		return
	}
	if len(info) == 0 {
		log.Error("[service.guard.v1.guard|sendEntryEffect] GetGuardByUIDRuid error(%v), uid(%d), ruid(%d) length==0", err, uid, ruid)
		return
	}

	var data model.GuardEntryEffects
	data.Business = _entryEffectBusinessId
	data.Data = make([]model.GuardEntryEffect, len(info))
	for index := 0; index < len(info); index++ {
		data.Data[index].EffectId = getEffectIdByLevel(info[index].PrivilegeType)
		data.Data[index].Uid = info[index].Uid
		data.Data[index].TargetId = info[index].TargetId
		data.Data[index].EndTime = info[index].ExpiredTime.Format("2006-01-02 15:04:05")
	}

	if err = s.liveEntryEffectPub.Send(ctx, strconv.FormatInt(uid, 10), data); err != nil {
		log.Error("[service.guard.v1.guard|sendEntryEffect] send error(%v), data(%v)", err, data)
	}
	return
}

func (s *GuardService) sendRoomMsg(ctx context.Context, uid int64, ruid int64, level int, buyType int) (err error) {
	if !s.c.GuardCfg.EnableGuardBroadcast {
		return
	}

	usrInfo, err := s.accRPC.GetUserInfo(ctx, uid)
	if err != nil || usrInfo == nil {
		log.Error("[service.v1.guard|sendRoomMsg] s.accRPC.Info3 error(%+v), params(%+v)", err, usrInfo)
		return
	}

	var roomReq roomV2.RoomRoomIdByUidMultiReq
	roomReq.Uids = make([]int64, 1)
	roomReq.Uids[0] = ruid
	roomRsp, err := s.roomAPI.V2Room.RoomIdByUidMulti(ctx, &roomReq)
	if err != nil || roomRsp.Code != 0 {
		log.Error("[service.v1.guard|sendRoomMsg] RoomIdByUidMulti error(%+v), params(%+v)", err, roomRsp)
		return
	}

	roomIdStr, ok := roomRsp.Data[strconv.FormatInt(ruid, 10)]
	if !ok {
		log.Error("[service.v1.guard|sendRoomMsg] RoomIdByUidMulti error(%+v), params(%+v)", err, roomRsp)
		return
	}

	var msg struct {
		Cmd  string `json:"cmd"`
		Data struct {
			OpType     int    `json:"op_type"`
			Uid        int64  `json:"uid"`
			UserName   string `json:"username"`
			GuardLevel int    `json:"guard_level"`
			IsShow     int    `json:"is_show"`
		} `json:"data"`
	}

	msg.Cmd = "USER_TOAST_MSG"
	msg.Data.OpType = buyType
	msg.Data.Uid = uid
	msg.Data.UserName = usrInfo.Name
	msg.Data.GuardLevel = level
	msg.Data.IsShow = 0

	bs, err := json.Marshal(msg)
	if err != nil {
		return
	}

	var resp struct {
		Ret int `json:"ret"`
	}

	req, err := xhttp.NewRequest(xhttp.MethodPost, fmt.Sprintf("%s/dm/1/push?ensure=1&cid=%s", s.c.GuardCfg.DanmuHost, roomIdStr), strings.NewReader(string(bs)))
	if err != nil {
		log.Error("[service.v1.guard|sendRoomMsg] sendRoomMsg error(%+v), params(%+v)", err, msg)
		return
	}

	err = s.cli.Do(ctx, req, &resp)
	if resp.Ret != 1 {
		log.Error("[service.v1.guard|sendRoomMsg] sendRoomMsg error(%+v), params(%+v), resp(%+v)", err, msg, resp)
	}
	return
}

// Close close guard service
func (s *GuardService) Close() {
	s.dao.Close()
}
