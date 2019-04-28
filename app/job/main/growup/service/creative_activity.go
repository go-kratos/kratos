package service

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/growup/model"

	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_like  = "1"
	_share = "2"
	_play  = "3"
	_reply = "4"
	_dm    = "5"
)

// CreativeActivity creative activity job
func (s *Service) CreativeActivity(c context.Context, date time.Time) (err error) {
	activities, err := s.dao.GetCActivities(c)
	if err != nil {
		log.Error("s.dao.GetCActivities error(%v)", err)
		return
	}
	// get signed up
	ups, err := s.signed(c, int64(_dbLimit))
	if err != nil {
		log.Error("s.signed error(%v)", err)
		return
	}
	now := xtime.Time(date.Unix())
	// calculate activity within the statistical period
	for _, ac := range activities {
		if now >= ac.StatisticsStart && now <= ac.StatisticsEnd {
			log.Info("calculate ac: %d", ac.ID)
			err = s.handleActivity(c, ac, ups, date)
			if err != nil {
				log.Error("s.handleActivity error(%v)", err)
				return
			}
		}
	}
	return
}

func (s *Service) handleActivity(c context.Context, ac *model.CActivity, upInfo map[int64]*model.UpInfoVideo, date time.Time) (err error) {
	// 获取已报名并且在签约时间内的up主，
	signedUps, err := s.getSignedUps(c, ac, upInfo)
	if err != nil {
		log.Error("s.getSignedUps error(%v)", err)
		return
	}
	// 获取活动稿件点赞、分享、播放、评论、弹幕
	archiveInfo, err := s.getArchiveInfo(c, ac.ID)
	if err != nil {
		log.Error("s.getArchiveInfo error(%v)", err)
		return
	}
	if len(archiveInfo) == 0 {
		log.Info("activity(%d) get 0 archiveInfo", ac.ID)
		return
	}
	log.Info("get %d archiveInfo", len(archiveInfo))
	fmt.Printf("get %d archiveInfo\n", len(archiveInfo))
	// 获取在投稿时间内的稿件
	avs, err := s.getAvsUpload(c, signedUps)
	if err != nil {
		log.Error("s.getAvsByMID error(%v)", err)
		return
	}
	if len(avs) == 0 {
		log.Info("activity(%d) get 0 av", ac.ID)
		return
	}
	log.Info("get %d archives", len(avs))
	fmt.Printf("get %d archives\n", len(avs))

	// 筛选在签约时间内的稿件
	handleActivityAvs(signedUps, avs, ac)
	for _, u := range signedUps {
		// 获取up主对应稿件要求的和
		err = s.getUpStatisData(c, u, ac, date, archiveInfo)
		if err != nil {
			log.Error("s.getUpData error(%v)", err)
			return
		}
	}
	// 计算中奖up主
	upBonus, err := s.calUpWin(c, signedUps, ac)
	if err != nil {
		log.Error("s.calUpBonus error(%v)", err)
		return
	}

	err = s.updateUpActivity(c, ac.ID, upBonus)
	if err != nil {
		log.Error("s.updateUpActivity error(%v)", err)
	}
	return
}

func (s *Service) calUpWin(c context.Context, ups map[int64]*model.UpActivity, ac *model.CActivity) (bonusUp []*model.UpActivity, err error) {
	activityBonus, err := s.dao.GetActivityBonus(c, ac.ID)
	if err != nil {
		log.Error("s.dao.GetActivityBonus error(%v)", err)
		return
	}
	bonusUp = make([]*model.UpActivity, 0)
	var (
		upList = make([]*model.UpActivity, len(ups))
		index  int
	)
	for _, up := range ups {
		upList[index] = up
		index++
	}
	sort.Slice(upList, func(i, j int) bool {
		return upList[i].ItemVal > upList[j].ItemVal
	})

	// 1:达标型 2:排序型
	if ac.WinType == 1 {
		for _, up := range upList {
			// uplist经过ItemVal,所以如果有ItemVal<RequireValue,可以直接break
			if up.ItemVal < ac.RequireValue {
				break
			}
			bonusUp = append(bonusUp, up)
		}
	} else if ac.WinType == 2 {
		// 去掉item为0的up
		for i := 0; i < len(upList); i++ {
			if upList[i].ItemVal == 0 {
				upList = upList[:i]
				break
			}
		}
		if len(upList) > int(ac.RequireValue) {
			upList = upList[:int(ac.RequireValue)]
		}
		bonusUp = upList
	}
	if len(bonusUp) == 0 {
		return
	}
	// 计算up主奖金
	calUpBonus(bonusUp, activityBonus, ac.BonusType, ac.WinType)
	return
}

// cal ups bonus
func calUpBonus(bonusUp []*model.UpActivity, activityBonus map[int64]int64, bonusType, winType int) (err error) {
	// 中奖类型 1:达标型 2:排序型
	// 奖金类型 1:平分 2:各得
	if winType == 1 {
		money, ok := activityBonus[0]
		if !ok {
			err = fmt.Errorf("活动奖金设置错误:达标型未设置金额")
			return
		}
		if bonusType == 1 {
			money = money / int64(len(bonusUp))
		}
		for _, up := range bonusUp {
			up.Bonus = money
			up.Rank = 0
			if up.State < 2 {
				// 中奖
				up.State = 2
				up.SuccessTime = xtime.Time(time.Now().Unix())
			}
		}
	} else if winType == 2 {
		other := len(activityBonus)
		otherMoney, ok := activityBonus[int64(other)]
		if !ok {
			err = fmt.Errorf("活动奖金设置错误:排序型没有其他金额")
			return
		}
		for i := 0; i < len(bonusUp); i++ {
			var (
				rank  = i + 1
				money int64
				ok    bool
			)
			if rank >= other {
				money = otherMoney
			} else {
				money, ok = activityBonus[int64(rank)]
				if !ok {
					err = fmt.Errorf("活动奖金设置错误:没有名次金额")
					return
				}
			}
			if money == 0 {
				continue
			}
			bonusUp[i].Bonus = money
			bonusUp[i].Rank = i + 1
			if bonusUp[i].State < 2 {
				// 中奖
				bonusUp[i].State = 2
				bonusUp[i].SuccessTime = xtime.Time(time.Now().Unix())
			}
		}
	}
	return
}

func (s *Service) getUpStatisData(c context.Context, up *model.UpActivity, ac *model.CActivity, date time.Time, archiveInfo map[int64]map[int]*model.ArchiveStat) (err error) {
	if len(up.AIDs) == 0 {
		return
	}
	avItem := make([]*model.AvItem, 0)
	aIDs := up.AIDs
	for _, avID := range aIDs {
		// 稿件统计数据: 计算当天的和 - 统计开始前一天的和
		var itemSumStart, itemSumEnd int64
		itemSumEnd, err = s.getAvStatisState(c, avID, ac, archiveInfo, 2)
		if err != nil {
			log.Error("s.getAvStatisState error(%v)", err)
			return
		}
		itemSumStart, err = s.getAvStatisState(c, avID, ac, archiveInfo, 1)
		if err != nil {
			log.Error("s.getAvStatisState error(%v)", err)
			return
		}
		itemSum := itemSumEnd - itemSumStart
		if itemSum > 0 {
			avItem = append(avItem, &model.AvItem{
				AvID:  avID,
				Value: itemSum,
			})
		}
	}
	if len(avItem) == 0 {
		return
	}
	sort.Slice(avItem, func(i, j int) bool {
		return avItem[i].Value > avItem[j].Value
	})

	// 1:uid 2 avid
	if ac.Object == 1 {
		up.AIDs = make([]int64, 0)
		for _, av := range avItem {
			if len(up.AIDs) < 5 {
				up.AIDs = append(up.AIDs, av.AvID)
			}
			up.ItemVal += av.Value
		}
		up.AIDNum = int64(len(avItem))
	} else if ac.Object == 2 {
		up.AIDs = []int64{avItem[0].AvID}
		up.ItemVal = avItem[0].Value
		up.AIDNum = 1
	}
	return
}

func (s *Service) getAvStatisState(c context.Context, avID int64, ac *model.CActivity, archiveInfo map[int64]map[int]*model.ArchiveStat, state int) (sum int64, err error) {
	if _, ok := archiveInfo[avID]; !ok {
		return
	}
	stat, ok := archiveInfo[avID][state]
	if !ok {
		return
	}
	requireItems := strings.Split(ac.RequireItems, ",")
	for _, item := range requireItems {
		switch item {
		case _like:
			sum += stat.Like
		case _share:
			sum += stat.Share
		case _play:
			sum += stat.Play
		case _reply:
			sum += stat.Reply
		case _dm:
			sum += stat.Dm
		}
	}
	return
}

func handleActivityAvs(ups map[int64]*model.UpActivity, avs []*model.AvUpload, ac *model.CActivity) {
	for _, av := range avs {
		if !(av.UploadTime >= ac.UploadStart && av.UploadTime <= ac.UploadEnd) {
			continue
		}
		if _, ok := ups[av.MID]; ok {
			if len(ups[av.MID].AIDs) == 0 {
				ups[av.MID].AIDs = make([]int64, 0)
			}
			ups[av.MID].AIDs = append(ups[av.MID].AIDs, av.AvID)
		}
	}
}

func (s *Service) getAvsUpload(c context.Context, ups map[int64]*model.UpActivity) (avs []*model.AvUpload, err error) {
	avs = make([]*model.AvUpload, 0)
	var id int64
	for {
		var av []*model.AvUpload
		av, err = s.dao.GetAvUploadByMID(c, id, _dbLimit)
		if err != nil {
			return
		}
		avs = append(avs, av...)
		if len(av) < _dbLimit {
			break
		}
		id = av[len(av)-1].ID
	}
	return
}

func (s *Service) getArchiveInfo(c context.Context, activityID int64) (archives map[int64]map[int]*model.ArchiveStat, err error) {
	archives = make(map[int64]map[int]*model.ArchiveStat) // map[av_id][state]*model.ArchiveStat
	var id int64
	for {
		var arch []*model.ArchiveStat
		arch, err = s.dao.GetArchiveInfo(c, activityID, id, _dbLimit)
		if err != nil {
			return
		}
		for _, a := range arch {
			if _, ok := archives[a.AvID]; !ok {
				archives[a.AvID] = make(map[int]*model.ArchiveStat)
			}
			archives[a.AvID][a.State] = a
		}
		if len(arch) < _dbLimit {
			break
		}
		id = arch[len(arch)-1].ID
	}
	return
}

func (s *Service) getSignedUps(c context.Context, ac *model.CActivity, upInfo map[int64]*model.UpInfoVideo) (signedUps map[int64]*model.UpActivity, err error) {
	signedUps = make(map[int64]*model.UpActivity)
	// if need sign up
	if ac.SignUp == 1 {
		var ups []*model.UpActivity
		ups, err = s.dao.ListUpActivity(c, ac.ID)
		if err != nil {
			return
		}
		for _, up := range ups {
			if info, ok := upInfo[up.MID]; ok && info.SignedAt >= ac.SignedStart && info.SignedAt <= ac.SignedEnd {
				signedUps[up.MID] = up
			}
		}
	} else {
		for _, info := range upInfo {
			if info.SignedAt >= ac.SignedStart && info.SignedAt <= ac.SignedEnd {
				signedUps[info.MID] = &model.UpActivity{MID: info.MID, ActivityID: ac.ID, State: 1, Nickname: info.Nickname}
			}
		}
	}
	return
}

func (s *Service) updateUpActivity(c context.Context, id int64, upBonus []*model.UpActivity) (err error) {
	// 更新所有之前获奖up主为已报名
	// update state 2->1 将所有状态设置为已报名
	_, err = s.dao.UpdateUpActivityState(c, id, 2, 1)
	if err != nil {
		log.Error("s.dao.UpdateUpActivityState error(%v)", err)
		return
	}
	if len(upBonus) > 0 {
		_, err = s.insertUpActivityBatch(c, upBonus)
		if err != nil {
			log.Error("s.insertUpActivity error(%v)", err)
		}
	}
	return
}
func (s *Service) insertUpActivityBatch(c context.Context, ups []*model.UpActivity) (rows int64, err error) {
	insert := make([]*model.UpActivity, _dbBatchSize)
	insertIndex := 0
	for _, up := range ups {
		insert[insertIndex] = up
		insertIndex++
		if insertIndex >= _dbBatchSize {
			_, err = s.insertUpActivity(c, insert[:insertIndex])
			if err != nil {
				return
			}
			insertIndex = 0
		}
	}
	if insertIndex > 0 {
		_, err = s.insertUpActivity(c, insert[:insertIndex])
	}
	return
}

func assembleUpActivity(ups []*model.UpActivity) (vals string) {
	var buf bytes.Buffer
	for _, row := range ups {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.MID, 10))
		buf.WriteByte(',')
		buf.WriteString("\"" + row.Nickname + "\"")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.ActivityID, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + xstr.JoinInts(row.AIDs) + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.AIDNum, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(row.Rank))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Bonus, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(row.State))
		buf.WriteByte(',')
		buf.WriteString("'" + row.SuccessTime.Time().Format("2006-01-02 15:04:05") + "'")
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals = buf.String()
	buf.Reset()
	return
}

func (s *Service) insertUpActivity(c context.Context, ups []*model.UpActivity) (rows int64, err error) {
	vals := assembleUpActivity(ups)
	rows, err = s.dao.InsertUpActivityBatch(c, vals)
	return
}
