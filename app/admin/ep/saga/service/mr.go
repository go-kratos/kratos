package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go-common/app/admin/ep/saga/model"
	"go-common/app/admin/ep/saga/service/utils"
	"go-common/library/log"

	"github.com/xanzy/go-gitlab"
)

const _specialMrID = 10199

// QueryProjectMr query project commit info according to project id.
func (s *Service) QueryProjectMr(c context.Context, req *model.ProjectDataReq) (resp *model.ProjectDataResp, err error) {

	if resp, err = s.QueryProject(c, model.ObjectMR, req); err != nil {
		return
	}
	return
}

// QueryTeamMr query team commit info according to department and business
func (s *Service) QueryTeamMr(c context.Context, req *model.TeamDataRequest) (resp *model.TeamDataResp, err error) {

	if resp, err = s.QueryTeam(c, model.ObjectMR, req); err != nil {
		return
	}
	return
}

// QueryProjectMrReport query mr review
func (s *Service) QueryProjectMrReport(c context.Context, req *model.ProjectMrReportReq) (resp *model.ProjectMrReportResp, err error) {
	var (
		info               []*model.MrInfo
		changeAdd          int
		changeDel          int
		mrCount            int
		stateCount         int
		discussionCount    int
		discussionResolved int
		mrTime             int
		spentTime          time.Duration
		avarageTime        time.Duration
		reviewers          []string
		reviews            []string
		reviewChangeAdd    int
		reviewChangeDel    int
		reviewTime         int
		reviewTotalTime    time.Duration
	)
	if info, err = s.QueryAllMergeRequestInfo(c, req.ProjectID); err != nil {
		return
	}

	for _, i := range info {
		for _, r := range i.Reviewers {
			if r.Name == req.Member {
				reviews = append(reviews, i.Author)
				reviewChangeAdd += i.ChangeAdd
				reviewChangeDel += i.ChangeDel
				reviewTime += i.SpentTime
			}
		}
		if i.Author == req.Member {
			mrCount++
			if i.State == model.StatusMerged {
				stateCount++
				mrTime += i.SpentTime
				changeAdd += i.ChangeAdd
				changeDel += i.ChangeDel
				discussionCount += i.TotalDiscussion
				discussionResolved += i.SolvedDiscussion
				for _, r := range i.Reviewers {
					reviewers = append(reviewers, r.Name)
				}
			}
		}
	}
	// 判断mr为零的情况
	if mrCount == 0 {
		resp = &model.ProjectMrReportResp{}
		return
	}

	if spentTime, err = time.ParseDuration(fmt.Sprintf("%ds", mrTime)); err != nil {
		return
	}
	if avarageTime, err = time.ParseDuration(fmt.Sprintf("%ds", mrTime/mrCount)); err != nil {
		return
	}
	if reviewTotalTime, err = time.ParseDuration(fmt.Sprintf("%ds", reviewTime)); err != nil {
		return
	}

	resp = &model.ProjectMrReportResp{
		ChangeAdd:       changeAdd,
		ChangeDel:       changeDel,
		MrCount:         mrCount,
		SpentTime:       spentTime.String(),
		StateCount:      stateCount,
		AverageMerge:    avarageTime.String(),
		Reviewers:       reviewers,
		Discussion:      discussionCount,
		Resolve:         discussionResolved,
		ReviewerOther:   reviews,
		ReviewChangeAdd: reviewChangeAdd,
		ReviewChangeDel: reviewChangeDel,
		ReviewTotalTime: reviewTotalTime.String(),
	}

	return
}

// QueryAllMergeRequestInfo ...
func (s *Service) QueryAllMergeRequestInfo(c context.Context, projID int) (info []*model.MrInfo, err error) {
	var (
		until = time.Now()
		since = until.AddDate(0, -1, 0)
		mrs   []*gitlab.MergeRequest
		resp  *gitlab.Response
	)

	for page := 1; ; page++ {
		if mrs, resp, err = s.gitlab.ListProjectMergeRequests(projID, &since, &until, page); err != nil {
			return
		}
		for _, m := range mrs {

			mr := &model.MrInfo{ProjectID: projID, MrID: m.IID, Author: m.Author.Name, State: m.State}

			if m.State == model.StatusMerged {
				spent := m.UpdatedAt.Sub(*m.CreatedAt)
				mr.SpentTime = int(spent.Seconds())
			}

			if mr.ChangeAdd, mr.ChangeDel, err = s.QueryMergeRequestDiff(c, projID, m.IID); err != nil {
				return
			}

			if mr.Reviewers, err = s.QueryMergeRequestReview(c, projID, m.IID); err != nil {
				return
			}

			if mr.TotalDiscussion, mr.SolvedDiscussion, err = s.QueryMergeRequestDiscussion(c, projID, m.IID); err != nil {
				return
			}

			info = append(info, mr)
		}
		if resp.NextPage == 0 {
			break
		}
	}

	return
}

// QueryMergeRequestReview 查询mr reviewer信息
func (s *Service) QueryMergeRequestReview(c context.Context, projectID, mrIID int) (reviewers []*model.MrReviewer, err error) {
	var (
		notes  []*model.StatisticsNotes
		emojis []*model.StatisticsMRAwardEmojis
		owners []string
		r      *regexp.Regexp
	)

	//query note
	if notes, err = s.dao.NoteByMRIID(c, projectID, mrIID); err != nil {
		return
	}

	// query emoji
	if emojis, err = s.dao.AwardEmojiByMRIID(c, projectID, mrIID); err != nil {
		return
	}

	//评论中解析获取owner
	if len(notes) == 0 {
		return
	}
	if r, err = regexp.Compile("OWNER:.*?@(.*)(；|)"); err == nil {
		matchResult := r.FindStringSubmatch(notes[len(notes)-1].Body)
		if len(matchResult) > 0 {
			owners = strings.Split(matchResult[1], " 或 @")
		}
	}

	mrCreatedAt := notes[len(notes)-1].CreatedAt

	// 从评论和表情中解析 reviewers
	for _, note := range notes {
		if note.Body == "+1" {
			reviewers = append(reviewers, &model.MrReviewer{Name: note.AuthorName, FinishedAt: note.CreatedAt, SpentTime: int(note.CreatedAt.Sub(*mrCreatedAt))})
		}
	}
	for _, emoji := range emojis {
		if emoji.Name == "thumbsup" {
			reviewers = append(reviewers, &model.MrReviewer{Name: emoji.UserName, FinishedAt: emoji.CreatedAt, SpentTime: int(emoji.CreatedAt.Sub(*mrCreatedAt))})
		}
	}

	// 判断reviewer类型
	for _, r := range reviewers {
		for _, owner := range owners {
			if owner == r.Name {
				r.UserType = "owner"
				break
			}
		}
		if r.UserType == "" {
			r.UserType = "other"
		}
	}

	return
}

// QueryMergeRequestDiscussion 查询获取mr的discussion
func (s *Service) QueryMergeRequestDiscussion(c context.Context, projectID, mrIID int) (total, solved int, err error) {
	var discussions []*model.StatisticsDiscussions
	if discussions, err = s.dao.DiscussionsByMRIID(c, projectID, mrIID); err != nil {
		return
	}
	for _, d := range discussions {
		var notesArray []int
		if err = json.Unmarshal([]byte(d.Notes), &notesArray); err != nil {
			return
		}
		for _, n := range notesArray {
			var note *model.StatisticsNotes
			if note, err = s.dao.NoteByID(c, projectID, mrIID, n); err != nil {
				return
			}
			if note.Resolvable {
				total++
				if note.Resolved {
					solved++
				}
			}
		}
	}
	return
}

// QueryMergeRequestDiff 查询获取到mr修改文件的总行数, 参数 p project_ID m MR_ID
func (s *Service) QueryMergeRequestDiff(c context.Context, p, m int) (ChangeAdd, ChangeDel int, err error) {
	var mr *gitlab.MergeRequest

	// 此MR数据量太大，曾导致过gitlab服务器崩溃，访问会返回502服务器错误，因此特殊处理。
	if m == _specialMrID {
		return 0, 0, nil
	}

	if mr, _, err = s.gitlab.GetMergeRequestDiff(p, m); err != nil {
		return
	}
	for _, change := range mr.Changes {
		rows := strings.Split(change.Diff, "\n")
		if len(rows) < 3 {
			// 处理diff为空
			continue
		}

		for _, row := range rows[3:] {
			if strings.HasPrefix(row, "+") {
				ChangeAdd++
			} else if strings.HasPrefix(row, "-") {
				ChangeDel++
			}
		}
	}

	return
}

/*-------------------------------------- sync MR ----------------------------------------*/

// SyncProjectMR ...
func (s *Service) SyncProjectMR(c context.Context, projectID int) (result *model.SyncResult, err error) {
	var (
		//syncAllTime = conf.Conf.Property.SyncData.SyncAllTime
		syncAllTime = false
		mrs         []*gitlab.MergeRequest
		resp        *gitlab.Response
		since       *time.Time
		until       *time.Time
		projectInfo *model.ProjectInfo
	)
	result = &model.SyncResult{}

	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		return
	}

	if !syncAllTime {
		since, until = utils.CalSyncTime()
	}
	log.Info("sync project(%d) MR time since: %v, until: %v", projectID, since, until)
	for page := 1; ; page++ {
		result.TotalPage++
		if mrs, resp, err = s.gitlab.ListProjectMergeRequests(projectID, since, until, page); err != nil {
			return
		}

		for _, mr := range mrs {
			if err = s.structureDatabaseMR(c, projectID, projectInfo.Name, mr); err != nil {
				log.Error("mr Save Database err: projectID(%d), MRIID(%d)", projectID, mr.IID)
				err = nil

				errData := &model.FailData{
					ChildID: mr.IID,
				}
				result.FailData = append(result.FailData, errData)
				continue
			}
			result.TotalNum++
		}

		if resp.NextPage == 0 {
			break
		}
	}
	return
}

// structureDatabaseMR ...
func (s *Service) structureDatabaseMR(c context.Context, projectID int, projectName string, mr *gitlab.MergeRequest) (err error) {
	var (
		milestoneID   int
		mrLables      string
		mrChanges     string
		mrLablesByte  []byte
		mrChangesByte []byte
	)

	if mr.Milestone != nil {
		milestoneID = mr.Milestone.ID
	}
	if mrLablesByte, err = json.Marshal(mr.Labels); err != nil {
		mrLables = model.JsonMarshalErrorText
	} else {
		mrLables = string(mrLablesByte)
	}

	if mrChangesByte, err = json.Marshal(mr.Changes); err != nil {
		mrChanges = model.JsonMarshalErrorText
	} else {
		mrChanges = string(mrChangesByte)
	}

	mrDB := &model.StatisticsMrs{
		MRID:                         mr.ID,
		MRIID:                        mr.IID,
		TargetBranch:                 mr.TargetBranch,
		SourceBranch:                 mr.SourceBranch,
		ProjectID:                    mr.ProjectID,
		ProjectName:                  projectName,
		Title:                        mr.Title,
		State:                        mr.State,
		CreatedAt:                    mr.CreatedAt,
		UpdatedAt:                    mr.UpdatedAt,
		Upvotes:                      mr.Upvotes,
		Downvotes:                    mr.Downvotes,
		AuthorID:                     mr.Author.ID,
		AuthorName:                   mr.Author.Name,
		AssigneeID:                   mr.Assignee.ID,
		AssigneeName:                 mr.Assignee.Name,
		SourceProjectID:              mr.SourceProjectID,
		TargetProjectID:              mr.TargetProjectID,
		Labels:                       mrLables,
		Description:                  mr.Description,
		WorkInProgress:               mr.WorkInProgress,
		MilestoneID:                  milestoneID,
		MergeWhenPipelineSucceeds:    mr.MergeWhenPipelineSucceeds,
		MergeStatus:                  mr.MergeStatus,
		MergedByID:                   mr.MergedBy.ID,
		MergedByName:                 mr.MergedBy.Name,
		MergedAt:                     mr.MergedAt,
		ClosedByID:                   mr.ClosedBy.ID,
		ClosedAt:                     mr.ClosedAt,
		Subscribed:                   mr.Subscribed,
		SHA:                          mr.SHA,
		MergeCommitSHA:               mr.MergeCommitSHA,
		UserNotesCount:               mr.UserNotesCount,
		ChangesCount:                 mr.ChangesCount,
		ShouldRemoveSourceBranch:     mr.ShouldRemoveSourceBranch,
		ForceRemoveSourceBranch:      mr.ForceRemoveSourceBranch,
		WebURL:                       mr.WebURL,
		DiscussionLocked:             mr.DiscussionLocked,
		Changes:                      mrChanges,
		TimeStatsHumanTimeEstimate:   mr.TimeStats.HumanTimeEstimate,
		TimeStatsHumanTotalTimeSpent: mr.TimeStats.HumanTotalTimeSpent,
		TimeStatsTimeEstimate:        mr.TimeStats.TimeEstimate,
		TimeStatsTotalTimeSpent:      mr.TimeStats.TotalTimeSpent,
		Squash:                       mr.Squash,
		PipelineID:                   mr.Pipeline.ID,
	}
	if len(mrDB.Labels) > model.MessageMaxLen {
		mrDB.Labels = mrDB.Labels[0 : model.MessageMaxLen-1]
	}
	if len(mrDB.Description) > model.MessageMaxLen {
		mrDB.Description = mrDB.Description[0 : model.MessageMaxLen-1]
	}
	return s.SaveDatabaseMR(c, mrDB)
}

// SaveDatabaseMR ...
func (s *Service) SaveDatabaseMR(c context.Context, mrDB *model.StatisticsMrs) (err error) {
	var total int

	if total, err = s.dao.HasMR(c, mrDB.ProjectID, mrDB.MRIID); err != nil {
		log.Error("SaveDatabaseMR HasMR(%+v)", err)
		return
	}

	// found only one, so update
	if total == 1 {
		return s.dao.UpdateMR(c, mrDB.ProjectID, mrDB.MRIID, mrDB)
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabaseMR mr has more rows(%d)", total)
		return
	}

	// insert row now
	return s.dao.CreateMR(c, mrDB)
}

/*-------------------------------------- agg MR ----------------------------------------*/

// AggregateProjectMR ...
func (s *Service) AggregateProjectMR(c context.Context, projectID int) (err error) {
	var (
		//syncAllTime = conf.Conf.Property.SyncData.SyncAllTime
		syncAllTime = false
		mrs         []*model.StatisticsMrs
		projectInfo *model.ProjectInfo
		since       *time.Time
		until       *time.Time
	)
	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		return
	}

	if !syncAllTime {
		since, until = utils.CalSyncTime()
	}

	if mrs, err = s.dao.MRByProjectID(c, projectID, since, until); err != nil {
		return
	}

	for _, mr := range mrs {
		if err = s.MRAddedInfoDB(c, projectID, mr); err != nil {
			log.Error("MRAddedInfoDB Save Database err: projectID(%d), MRIID(%d)", projectID, mr.MRIID)
			err = nil
		}
		if err = s.MRReviewerDB(c, projectID, mr.MRIID, projectInfo.Name, mr); err != nil {
			log.Error("MRReviewerDB Save Database err: projectID(%d), MRIID(%d)", projectID, mr.MRIID)
			err = nil
		}
	}
	return
}

// MRAddedInfoDB ...
func (s *Service) MRAddedInfoDB(c context.Context, projectID int, mr *model.StatisticsMrs) (err error) {
	if mr.ChangeAdd, mr.ChangeDel, err = s.QueryMergeRequestDiff(c, projectID, mr.MRIID); err != nil {
		return
	}
	if mr.TotalDiscussion, mr.SolvedDiscussion, err = s.QueryMergeRequestDiscussion(c, projectID, mr.MRIID); err != nil {
		return
	}
	return s.dao.UpdateMR(c, projectID, mr.MRIID, mr)
}

// MRReviewerDB 查询mr reviewer信息
func (s *Service) MRReviewerDB(c context.Context, projectID, mrIID int, projectName string, mr *model.StatisticsMrs) (err error) {
	var (
		notes     []*model.StatisticsNotes
		emojis    []*model.StatisticsMRAwardEmojis
		owners    []string
		r         *regexp.Regexp
		reviewers []*model.AggregateMrReviewer
	)

	//query note
	if notes, err = s.dao.NoteByMRIID(c, projectID, mrIID); err != nil {
		return
	}

	// query emoji
	if emojis, err = s.dao.AwardEmojiByMRIID(c, projectID, mrIID); err != nil {
		return
	}

	//评论中解析获取owner
	if len(notes) == 0 {
		return
	}
	if r, err = regexp.Compile("OWNER:.*?@(.*)(；|)"); err == nil {
		matchResult := r.FindStringSubmatch(notes[len(notes)-1].Body)
		if len(matchResult) > 0 {
			owners = strings.Split(matchResult[1], " 或 @")
		}
	}

	mrCreatedAt := notes[len(notes)-1].CreatedAt

	// 从评论和表情中解析 reviewers
	for _, note := range notes {
		if note.Body == "+1" {
			reviewer := &model.AggregateMrReviewer{
				ReviewerID:    note.AuthorID,
				ReviewerName:  note.AuthorName,
				ReviewType:    "note",
				ReviewID:      note.NoteID,
				ReviewCommand: "+1",
				CreatedAt:     note.CreatedAt,
				ApproveTime:   int(note.CreatedAt.Sub(*mrCreatedAt).Seconds()),
				MergeTime:     int(mr.UpdatedAt.Sub(*note.CreatedAt).Seconds()),
			}
			reviewers = append(reviewers, reviewer)
		}
	}
	for _, emoji := range emojis {
		if emoji.Name == "thumbsup" {
			reviewer := &model.AggregateMrReviewer{
				ReviewerID:    emoji.UserID,
				ReviewerName:  emoji.UserName,
				ReviewType:    "emoji",
				ReviewID:      emoji.AwardEmojiID,
				ReviewCommand: "thumbsup",
				CreatedAt:     emoji.CreatedAt,
				ApproveTime:   int(emoji.CreatedAt.Sub(*mrCreatedAt).Seconds()),
				MergeTime:     int(mr.UpdatedAt.Sub(*emoji.CreatedAt).Seconds()),
			}
			reviewers = append(reviewers, reviewer)
		}
	}

	// 判断reviewer类型
	for _, r := range reviewers {
		r.ProjectID = projectID
		r.ProjectName = projectName
		r.MrIID = mrIID
		r.Title = mr.Title
		r.WebUrl = mr.WebURL
		r.AuthorName = mr.AuthorName
		if utils.InSlice(r.ReviewerName, owners) {
			r.UserType = "owner"
		} else {
			r.UserType = "other"
		}
		if err = s.SaveDatabaseAggMR(c, r); err != nil {
			log.Error("mrReviewer 存数据库报错(%+V)", err)
			continue
		}
	}
	return
}

// SaveDatabaseAggMR ...
func (s *Service) SaveDatabaseAggMR(c context.Context, mrDB *model.AggregateMrReviewer) (err error) {
	var total int

	if total, err = s.dao.HasAggregateReviewer(c, mrDB.ProjectID, mrDB.MrIID, mrDB.ReviewerID, mrDB.ReviewID); err != nil {
		log.Error("SaveDatabaseAggMR HasAggMR(%+v)", err)
		return
	}

	// found only one, so update
	if total == 1 {
		return s.dao.UpdateAggregateReviewer(c, mrDB.ProjectID, mrDB.MrIID, mrDB.ReviewerID, mrDB.ReviewID, mrDB)
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabaseAggMR aggMR has more rows(%d)", total)
		return
	}

	// insert row now
	return s.dao.CreateAggregateReviewer(c, mrDB)
}
