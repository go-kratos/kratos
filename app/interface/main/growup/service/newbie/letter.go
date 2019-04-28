package newbie

import (
	"context"
	"go-common/app/interface/main/growup/conf"
	"go-common/app/interface/main/growup/dao/newbiedao"
	"go-common/app/interface/main/growup/model"
	accApi "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup.v2"

	"strconv"
	"time"
)

// Letter newbie letter
func (s *Service) Letter(c context.Context, req *model.NewbieLetterReq) (*model.NewbieLetterRes, error) {
	var (
		group        *errgroup.Group
		recUps       = make(map[int64]*model.RecommendUp)
		NewbieConf   = conf.Conf.Newbie
		category     *model.Category
		recUpMidList []int64
		activities   []*model.Activity
		archive      *model.VideoUpArchive
		accInfo      *accApi.InfoReply

		i   = 0
		ok  bool
		err error
	)
	res := new(model.NewbieLetterRes)

	log.Info("req: %+v", req)
	group = errgroup.WithCancel(c)
	// get up info
	group.Go(func(ctx context.Context) error {
		accInfo, err = s.dao.GetInfo(ctx, req.Mid)
		return err
	})

	// get activities
	group.Go(func(ctx context.Context) error {
		activities, err = s.dao.GetActivities(ctx)
		return err
	})

	// get video up, and set talent
	group.Go(func(ctx context.Context) error {
		archive, err = s.dao.GetVideoUp(ctx, req.Aid)
		return err
	})

	err = group.Wait()
	if err != nil {
		return nil, err
	}

	// data validation, deal with default data
	if req.Mid != archive.Mid {
		log.Error("The archive is not yours, mid: %d, archive.Mid: %v", req.Mid, archive.Mid)
		return nil, ecode.GrowupArchiveNotYours
	}
	if category, ok = newbiedao.Categories[archive.Tid]; !ok {
		log.Error("not found the sub tid, sub tid: %d, Categories: %v", archive.Tid, newbiedao.Categories)
		return nil, ecode.GrowupSubTidNotExist
	}
	if _, ok = newbiedao.Categories[category.Pid]; !ok {
		log.Error("not found the tid, tid: %d, Categories: %v", archive.Tid, newbiedao.Categories)
		return nil, ecode.GrowupTidNotExist
	}
	res.Area = newbiedao.Categories[category.Pid].Name
	log.Info("sub tid: %d, tid: %d", archive.Tid, category.Pid)

	sTid := strconv.FormatInt(archive.Tid, 10)
	if res.Talent, ok = NewbieConf.Talents[sTid]; !ok {
		res.Talent = NewbieConf.DefaultTalent
	}
	for _, activity := range activities {
		if i >= NewbieConf.ActivityCount {
			break
		}
		if activity.Type != NewbieConf.ActivityShotType {
			continue
		}
		if activity.Cover == "" {
			activity.Cover = NewbieConf.DefaultCover
		}
		res.Activities = append(res.Activities, activity)
		i++
	}
	if len(res.Activities) < NewbieConf.ActivityCount {
		log.Error("activity count is not enough %d", NewbieConf.ActivityCount)
		return nil, ecode.GrowupActivityCountNotEnough
	}

	res.UperInfo = new(model.NewbieLetterUpInfo)
	res.UperInfo.Mid = accInfo.Info.Mid
	res.UperInfo.Name = accInfo.Info.Name

	res.Archive = new(model.NewbieLetterArchive)
	res.Archive.Title = archive.Title
	res.Archive.PTime = time.Unix(archive.PTime, 0).Format(model.TimeLayout)
	log.Info("after data validation: data(%+v)", res)

	// get recommend up list
	if _, ok := newbiedao.RecommendUpList[category.Pid]; !ok {
		for _, lists := range newbiedao.RecommendUpList {
			for recUpMid, recUp := range lists {
				recUps[recUpMid] = recUp
				break
			}
		}
		log.Info("Not found recommend up list, system random get them : %+v", recUps)
	} else {
		recUps = newbiedao.RecommendUpList[category.Pid]
		log.Info("found recommend up list : %+v", recUps)
	}

	// get relations
	i = 0
	for recUpMid := range recUps {
		if i >= NewbieConf.RecommendUpPoolCount {
			break
		}
		if recUpMid == req.Mid {
			continue
		}
		recUpMidList = append(recUpMidList, recUpMid)
		i++
	}
	log.Info("recUpMidList: %+v", recUpMidList)

	relations, err := s.dao.GetRelations(c, req.Mid, recUpMidList)
	if err != nil {
		return nil, err
	}
	log.Info("relations: %+v", relations)

	// get ups info
	infosReply, err := s.dao.GetInfos(c, recUpMidList)
	if err != nil {
		err = ecode.GrowupRecommendUpNotExist
		return nil, err
	}
	log.Info("recUpInfos: %+v", infosReply.Infos)

	// select 3 ups
	i = 0
	for recUpMid := range recUps {
		if i >= NewbieConf.RecommendUpCount {
			break
		}
		if _, ok := infosReply.Infos[recUpMid]; !ok {
			continue
		}

		if _, ok := relations[recUpMid]; !ok {
			relations[recUpMid] = &model.Relation{
				Mid:       recUpMid,
				Attribute: -1,
			}
		}
		relations[recUpMid].Face = infosReply.Infos[recUpMid].Face
		relations[recUpMid].Name = infosReply.Infos[recUpMid].Name
		res.Relations = append(res.Relations, relations[recUpMid])

		i++
	}

	log.Info("res.Relations: %+v", res.Relations)
	return res, nil
}
