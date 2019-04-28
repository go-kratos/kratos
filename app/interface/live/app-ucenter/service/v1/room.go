package v1

import (
	"context"

	v1pb "go-common/app/interface/live/app-ucenter/api/http/v1"
	"go-common/app/interface/live/app-ucenter/conf"
	"go-common/app/interface/live/app-ucenter/dao"
	"go-common/library/ecode"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

// RoomService struct
type RoomService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewRoomService init
func NewRoomService(c *conf.Config) (s *RoomService) {
	s = &RoomService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

// GetInfo implementation
// 获取房间基本信息
// `method:"GET" midware:"auth"`
func (s *RoomService) GetInfo(ctx context.Context, req *v1pb.GetRoomInfoReq) (resp *v1pb.GetRoomInfoResp, err error) {
	uid, ok := metadata.Value(ctx, metadata.Mid).(int64)
	if uid <= 0 || !ok {
		err = ecode.UidError
		return
	}
	uidArr := make([]int64, 0)
	uidArr = append(uidArr, uid)
	resp = &v1pb.GetRoomInfoResp{}
	resp.Uid = uid
	resp.FullText = "V等级5级或UP等级10级才能开通粉丝勋章哦~加油！"
	resp.OpenMedalLevel = dao.OpenFansMealLevel
	resp.MaxLevel = 40
	eg := errgroup.Group{}
	eg.Go(func() (err error) {
		room, err := s.dao.GetRoomInfosByUids(ctx, uidArr)
		if err != nil {
			err = ecode.CallRoomError
			return
		}
		roomInfo := room[uid]
		if roomInfo == nil {
			return
		}
		LockStatus := 0
		//if roomInfo.LockTill >= dao.PERMANENT_LOCK_TIME {
		//	LockStatus = 1
		//}
		resp.RoomId = roomInfo.RoomId
		resp.Title = roomInfo.Title
		resp.LiveStatus = roomInfo.LiveStatus
		resp.AreaV2Id = roomInfo.AreaV2Id
		resp.AreaV2Name = roomInfo.AreaV2Name
		resp.LockTime = roomInfo.LockTill
		resp.LockStatus = int64(LockStatus)
		resp.ParentId = roomInfo.AreaV2ParentId
		resp.ParentName = roomInfo.AreaV2ParentName
		return
	})
	eg.Go(func() (err error) {
		user, err := s.dao.GetUserInfo(ctx, uidArr)
		if err != nil {
			err = ecode.CallUserError
			return
		}
		if user == nil {
			err = ecode.UserNotFound
			return
		}
		userInfo, ok := user[uid]
		if !ok {
			err = ecode.UserNotFound
			return
		}
		resp.Face = userInfo.GetInfo().GetFace()
		resp.Uname = userInfo.GetInfo().GetUname()
		if userInfo.Exp != nil {
			resp.MasterScore = userInfo.Exp.Rcost / 100
			if userInfo.Exp.MasterLevel != nil {
				masterLevel := userInfo.Exp.MasterLevel.Level
				masterNextLevel := masterLevel + 1
				if masterNextLevel > 40 {
					masterNextLevel = 40
				}
				resp.MasterLevel = masterLevel
				resp.MasterLevelColor = userInfo.Exp.MasterLevel.Color
				resp.MasterNextLevel = masterNextLevel
				if len(userInfo.Exp.MasterLevel.Next) >= 2 {
					resp.MasterNextLevelScore = userInfo.Exp.MasterLevel.Next[1]
				}
			}
		}

		return
	})
	eg.Go(func() (err error) {
		relation, err := s.dao.GetUserFc(ctx, uid)
		if err != nil {
			err = ecode.CallRelationError
			return
		}
		resp.FcNum = relation.Fc
		return
	})
	eg.Go(func() (err error) {
		fansMedal, err := s.dao.GetFansMedalInfo(ctx, uid)
		if err != nil {
			err = ecode.CallFansMedalError
			return
		}
		if fansMedal == nil {
			return
		}
		resp.IsMedal = fansMedal.MasterStatus
		resp.MedalName = fansMedal.MedalName
		resp.MedalRenameStatus = fansMedal.RenameStatus
		resp.MedalStatus = fansMedal.Status
		return
	})
	eg.Go(func() (err error) {
		identifyStatus, err := s.dao.GetIdentityStatus(ctx, uid)
		if err != nil {
			err = ecode.CallMainMemberError
			return
		}
		resp.IdentifyStatus = int64(identifyStatus)
		return
	})
	err = eg.Wait()
	return
}

// Create implementation
// 创建房间
// `method:"POST" midware:"auth"`
func (s *RoomService) Create(ctx context.Context, req *v1pb.CreateReq) (resp *v1pb.CreateResp, err error) {
	uid, ok := metadata.Value(ctx, metadata.Mid).(int64)
	if uid <= 0 || !ok {
		err = ecode.UidError
		return
	}
	resp = &v1pb.CreateResp{}
	room, err := s.dao.CreateRoom(ctx, uid)
	if err != nil {
		err = ecode.CallRoomError
		return
	}
	resp.RoomId = room.Roomid
	return
}
