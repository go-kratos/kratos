package service

import (
	"context"
	"fmt"
	"go-common/app/service/bbq/user/api"
	"go-common/app/service/bbq/user/internal/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
)

// Login 登录
func (s *Service) Login(c context.Context, in *api.UserBase) (res *api.UserBase, err error) {
	res = &api.UserBase{}
	mid := in.Mid
	newTag := in.NewTag
	//默认参数
	defaultUserBase := api.UserBase{
		Mid:            mid,
		Uname:          "Qing_" + strconv.FormatInt(mid, 10),
		Face:           "http://i0.hdslb.com/bfs/bbq/video-image/userface/1558868601542006937.png",
		Sex:            model.SexAnimal,
		Signature:      "",
		UserType:       model.UserTypeNew,
		CompleteDegree: model.DegreeComp,
	}
	addUserBase := &defaultUserBase
	*res = defaultUserBase
	var rawUserBases map[int64]*api.UserBase
	if rawUserBases, err = s.dao.RawUserBase(c, []int64{mid}); err != nil {
		log.Warnw(c, "log", "get raw user base fail", "mid", mid)
		return
	}

	// 是否在bbq数据库中存在
	if rawUserBase, ok := rawUserBases[mid]; ok {
		*res = *rawUserBase
		if rawUserBase.CompleteDegree == model.DegreeUncomp {
			addUserBase = rawUserBase
			addUserBase.CompleteDegree = model.DegreeComp
			s.dao.UpdateUserBase(c, mid, addUserBase)
		}
		return
	}

	// 获取主站的账号信息
	var userCard *model.UserCard
	userCard, err = s.dao.RawUserCard(c, mid)
	// 获取失败则插入默认
	if err != nil || userCard == nil {
		log.V(10).Infow(c, "log", "get raw user card from main fail")
		if _, err = s.dao.AddUserBase(c, addUserBase); err != nil {
			log.Errorv(c, log.KV("log", fmt.Sprintf("Login inset default info err,err:%v", err)))
			return
		}
		*res = defaultUserBase
	} else {
		// 变换主站信息到bbq，并插入
		if userCard.Sex == "男" {
			addUserBase.Sex = model.SexMan
		} else if userCard.Sex == "女" {
			addUserBase.Sex = model.SexWoman
		} else {
			addUserBase.Sex = model.SexAnimal
		}
		addUserBase.Face = userCard.Face
		// 这里默认用生成的Qing_{{mid}}存到mysql，不使用主站昵称
		addUserBase.CompleteDegree = model.DegreeComp
		addUserBase.UserType = model.UserTypeBili
		if _, err = s.dao.AddUserBase(c, addUserBase); err != nil {
			log.Errorv(c, log.KV("log", fmt.Sprintf("Login inset bili info err,err:%v", err)))
		}

		*res = *addUserBase
		// 返回的时候把主站的数据吐给前端展示，把名字是否重复、屏蔽词等等逻辑都放给edit接口去判断
		res.Uname = userCard.Name
	}

	// 返回一次提示
	res.CompleteDegree = model.DegreeUncomp

	//返回客户端同学要求
	if newTag > 0 {
		defaultUserBase.NewTag = newTag
	}

	return
}

// ListUserInfo 用户信息列表
func (s *Service) ListUserInfo(c context.Context, in *api.ListUserInfoReq) (res *api.ListUserInfoReply, err error) {
	res = new(api.ListUserInfoReply)
	upMIDs := in.UpMid
	if len(upMIDs) == 0 {
		return
	}
	if len(upMIDs) > model.BatchUserLen {
		err = ecode.BatchUserTooLong
		return
	}

	userInfos, err := s.batchUserInfo(c, in.Mid, upMIDs, &api.ListUserInfoConf{NeedDesc: in.NeedDesc, NeedStat: in.NeedStat, NeedFollowState: in.NeedFollowState})
	if err != nil {
		log.Warnv(c, log.KV("log", "batch user info fail"))
		return
	}

	for _, mid := range upMIDs {
		userInfo, exists := userInfos[mid]
		if !exists {
			log.Warnv(c, log.KV("log", fmt.Sprintf("get user info fail: mid=%d", mid)))
			continue
		}
		res.List = append(res.List, userInfo)
	}

	return
}

// UserEdit 完善用户信息
//              该请求需要保证请求的mid已经存在，如果不存在，该接口会返回失败
func (s *Service) UserEdit(c context.Context, in *api.UserBase) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	// 0. 参数正确性校验
	if err = s.dao.CheckUname(c, in.Mid, in.Uname); err != nil {
		return
	}

	// 1. 更新 编辑接口不关心用户信息完成度，直接更新为完善
	in.CompleteDegree = model.DegreeComp
	if _, err = s.dao.UpdateUserBase(c, in.Mid, in); err != nil {
		err = ecode.EditUserBaseErr
		return nil, err
	}
	return
}

// UserCmsTagEdit 修改cms_tag
func (s *Service) UserCmsTagEdit(c context.Context, in *api.CmsTagRequest) (res *empty.Empty, err error) {
	// TODO: 这里就不要用事务了，后面改掉，同时UserFieldEdit也改掉
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("edit user field begin tran err(%v) mid(%d)", err, in.Mid)))
		err = ecode.BBQSystemErr
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if _, err = s.dao.UpdateUserField(c, tx, in.Mid, "cms_tag", in.CmsTag); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("update user cms_tag fail: mid=%d, cms_tag=%d", in.Mid, in.CmsTag)))
		return
	}
	return
}

// UserFieldEdit 仅提供给审核使用
func (s *Service) UserFieldEdit(c context.Context, in *api.UserBase) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("edit user field begin tran err(%v) mid(%d)", err, in.Mid)))
		err = ecode.BBQSystemErr
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if len(in.Uname) != 0 {
		if _, err = s.dao.UpdateUserField(c, tx, in.Mid, "uname", in.Uname); err != nil {
			log.Errorv(c, log.KV("log", fmt.Sprintf("update user uname fail: mid=%d, uname=%s", in.Mid, in.Uname)))
			return
		}
	}
	if len(in.Face) != 0 {
		if _, err = s.dao.UpdateUserField(c, tx, in.Mid, "face", in.Face); err != nil {
			log.Errorv(c, log.KV("log", fmt.Sprintf("update user uname fail: mid=%d, uname=%s", in.Mid, in.Uname)))
			return
		}
	}

	return
}

func (s *Service) batchUserInfo(c context.Context, visitorMID int64, upMIDs []int64, conf *api.ListUserInfoConf) (res map[int64]*api.UserInfo, err error) {
	res = make(map[int64]*api.UserInfo)
	if len(upMIDs) == 0 {
		return
	}
	// 1. 请求user base信息
	userBases, err := s.dao.UserBase(c, upMIDs)
	if err != nil {
		log.Errorv(c, log.KV("event", "get_user_base"), log.KV("error", err))
		return
	}
	log.V(1).Infov(c, log.KV("event", "get_user_base"), log.KV("req_size", len(upMIDs)),
		log.KV("rsp_size", len(userBases)))

	// 2. 请求user statistics信息
	var userStatistics map[int64]*api.UserStat
	if conf.NeedStat {
		userStatistics, err = s.dao.RawBatchUserStatistics(c, upMIDs)
		if err != nil {
			log.Errorv(c, log.KV("event", "user_statistics"))
			err = nil
			userStatistics = make(map[int64]*api.UserStat)
		}
	}

	// 3. 请求社交关系，是否关注，是否粉丝
	var visitorFollowedMID map[int64]bool
	var visitorFanMID map[int64]bool
	if conf.NeedFollowState && visitorMID != 0 {
		visitorFollowedMID = s.dao.IsFollow(c, visitorMID, upMIDs)
		visitorFanMID = s.dao.IsFan(c, visitorMID, upMIDs)
	}

	// 4. 组装回包
	for mid, userBase := range userBases {
		userInfo := new(api.UserInfo)
		res[mid] = userInfo
		userInfo.UserBase = userBase
		if len(userBase.Face) == 0 {
			userBase.Face = "http://i0.hdslb.com/bfs/bbq/video-image/userface/1558868601542006937.png"
		}

		if conf.NeedDesc {
			userInfo.UserBase.UserDesc, userInfo.UserBase.RegionName = s.genUserDesc(c, userBase)
		}

		if conf.NeedStat && userStatistics != nil {
			userInfo.UserStat = userStatistics[mid]
		}

		if conf.NeedFollowState && visitorMID != 0 {
			// mid和visitorMID的关系
			_, follow := visitorFollowedMID[mid]
			_, fan := visitorFanMID[mid]
			if follow {
				userInfo.FollowState |= 1
			}
			if fan {
				userInfo.FollowState |= 2
			}
		}
	}
	return
}

// genUserDesc form tags from birthday, sex and region
func (s *Service) genUserDesc(c context.Context, base *api.UserBase) (userDesc []string, regionName string) {
	var (
		birthday string
		sex      string
	)
	if len(base.Birthday) != 8 {
		birthday = "其他"
	} else {
		if base.Birthday > "20100000" {
			birthday = "10后"
		} else if base.Birthday > "20000000" {
			birthday = "00后"
		} else if base.Birthday > "19900000" {
			birthday = "90后"
		} else if base.Birthday > "19800000" {
			birthday = "80后"
		} else if base.Birthday > "19700000" {
			birthday = "70后"
		} else {
			birthday = "其他"
		}
	}
	userDesc = append(userDesc, birthday)

	switch base.Sex {
	case 1:
		sex = "男"
	case 2:
		sex = "女"
	default:
		sex = "不明生物"
	}
	userDesc = append(userDesc, sex)

	if base.Region == 0 {
		return
	}
	region := int32(base.Region / 100 * 100)
	if location, err := s.dao.GetLocation(c, region); err != nil {
		log.Errorv(c, log.KV("event", "get_location"))
	} else if location == nil {
		log.Errorv(c, log.KV("event", "get_location"), log.KV("reason", "no_location"), log.KV("loc_id", base.Region))
	} else {
		userDesc = append(userDesc, location.Name)
		regionName = location.Name
	}
	return
}

//PhoneCheck ..
func (s *Service) PhoneCheck(c context.Context, in *api.PhoneCheckReq) (res *api.PhoneCheckReply, err error) {
	res = &api.PhoneCheckReply{}
	userProfile, err := s.dao.GetUserBProfile(c, in)
	if err != nil || userProfile == nil {
		log.Warn("PhoneCheck get userb profile err:%v,res:%v", err, userProfile)
	}
	res.TelStatus = userProfile.Profile.GetTelStatus()
	return
}

//ForbidUser ...
func (s *Service) ForbidUser(c context.Context, req *api.ForbidRequest) (res *empty.Empty, err error) {
	if err = s.dao.ForbidUser(c, req.MID, req.ExpireTime); err != nil {
		log.Warnv(c, log.KV("event", "Service ForbidUser"))
		return
	}
	return
}

//ReleaseUser ...
func (s *Service) ReleaseUser(c context.Context, req *api.ReleaseRequest) (res *empty.Empty, err error) {
	if err = s.dao.ReleaseUser(c, req.MID); err != nil {
		log.Warnv(c, log.KV("event", "Service ForbidUser"))
		return
	}
	return
}

// UpdateUserVideoView .
func (s *Service) UpdateUserVideoView(c context.Context, req *api.UserVideoView) (res *empty.Empty, err error) {
	err = s.dao.UpdateUserVideoView(c, req.Mid, req.Views)
	return
}
