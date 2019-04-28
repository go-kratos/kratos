package service

import (
	"context"
	"database/sql"
	"strings"
	"unicode/utf8"

	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/model"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

func (s *Service) checkGroupData(arg *model.AddGroupArg) error {
	if strings.Trim(arg.Name, " ") == "" || utf8.RuneCountInString(arg.Name) > 15 {
		return errors.New("分组名称需在1-15个字符之间")
	} else if arg.Tag == "" || utf8.RuneCountInString(arg.Tag) > 10 {
		return errors.New("标签名称需在1-10个字符之间")
	} else if utf8.RuneCountInString(arg.ShortTag) > 2 {
		return errors.New("标签简称需在1-2个字符之间")
	} else if utf8.RuneCountInString(arg.Remark) > 50 {
		return errors.New("备注不能超过50字符")
	}
	return nil

}

//AddGroup add group
func (s *Service) AddGroup(c *blademaster.Context, arg *model.AddGroupArg) (result sql.Result, err error) {
	if err = s.checkGroupData(arg); err != nil {
		log.Error("add group error, %v", err)
		return
	}

	var exist bool
	if exist, err = s.mng.CheckGroupExist(c, arg, 0); err != nil {
		log.Error("check group exist fail, err=%v, arg=%+v", err, arg)
		return
	}

	if exist {
		err = errors.New("该分组、标签或标签简称已存在")
		log.Error("group with same name, tag or short_tag exist, arg=%+v", arg)
		return
	}

	result, err = s.mng.AddGroup(c, arg)
	if err != nil {
		log.Error("add group db error, %v", err)
		return
	}

	return
}

//UpdateGroup update group
func (s *Service) UpdateGroup(c *blademaster.Context, arg *model.EditGroupArg) (result sql.Result, err error) {
	if err = s.checkGroupData(arg.AddArg); err != nil {
		log.Error("update group error, %v", err)
		return
	}

	var exist bool
	if exist, err = s.mng.CheckGroupExist(c, arg.AddArg, arg.ID); err != nil {
		log.Error("check group exist fail, err=%v, arg=%+v", err, arg)
		return
	}

	if exist {
		err = errors.New("该分组、标签或标签简称已存在")
		log.Error("group with same name, tag or short_tag exist, arg=%+v", arg)
		return
	}

	result, err = s.mng.UpdateGroup(c, arg)
	if err != nil {
		log.Error("update group db error, %v", err)
		return
	}

	return
}

//RemoveGroup remove group
func (s *Service) RemoveGroup(c *blademaster.Context, arg *model.RemoveGroupArg) (result sql.Result, err error) {
	result, err = s.mng.RemoveGroup(c, arg)
	if err != nil {
		log.Error("remove group db error, %v", err)
		return
	}
	return
}

//GetGroup get group
func (s *Service) GetGroup(c *blademaster.Context, arg *model.GetGroupArg) (res []*model.UpGroup, err error) {
	return s.mng.GetGroup(c, arg)
}

// getGroupCache get group from cache
func (s *Service) getGroupCache(groupID int64) (group *model.UpGroup) {
	if g, ok := s.spGroupsCache[groupID]; ok {
		return &model.UpGroup{ID: g.ID, Name: g.Name, Tag: g.Tag, ShortTag: g.ShortTag, Remark: g.Note,
			State: 1, FontColor: g.FontColor, BgColor: g.BgColor}
	}
	return
}

// UpGroups .
func (s *Service) UpGroups(c context.Context, req *upgrpc.NoArgReq) (res *upgrpc.UpGroupsReply, err error) {
	res = new(upgrpc.UpGroupsReply)
	if len(s.spGroupsCache) > 0 {
		res.UpGroups = s.spGroupsCache
	}
	return
}

func (s *Service) loadSpGroupsMids() {
	var (
		err    error
		gids   []int64
		lastID int64
		c      = context.Background()
		gmap   map[int64][]int64
	)
	for _, v := range s.spGroupsCache {
		gids = append(gids, v.ID)
	}
	gmap = make(map[int64][]int64)
	defer func() {
		if err == nil {
			s.spGroupsMidsCache = gmap
		}
	}()
	for {
		var (
			lid   int64
			ps    = 10000
			gmids map[int64][]int64
		)
		if lid, gmids, err = s.mng.UpGroupsMids(c, gids, lastID, ps); err != nil {
			log.Error("s.mng.UpGroupsMids(%+v,%d,%d)", gids, lastID, ps)
			return
		}
		for _, gid := range gids {
			if mids, ok := gmids[gid]; ok {
				gmap[gid] = append(gmap[gid], mids...)
			}
		}
		if lid == 0 {
			return
		}
		lastID = lid
	}
}
