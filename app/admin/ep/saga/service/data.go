package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/ep/saga/model"
	"go-common/app/admin/ep/saga/service/utils"
	"go-common/library/log"

	"github.com/xanzy/go-gitlab"
)

/*-------------------------------------- sync issue ----------------------------------------*/

// SyncAllIssues ...
func (s *Service) SyncAllIssues(projectID int) (totalPage, totalNum int, err error) {
	var (
		//syncAllTime = conf.Conf.Property.SyncData.SyncAllTime
		syncAllTime = false
		issues      []*gitlab.Issue
		resp        *gitlab.Response
		since       *time.Time
		until       *time.Time
	)
	if !syncAllTime {
		since, until = utils.CalSyncTime()
	}

	for page := 1; ; page++ {
		totalPage++
		if issues, resp, err = s.gitlab.ListProjectIssues(projectID, page, since, until); err != nil {
			return
		}
		for _, issue := range issues {
			var (
				issueAssignees string
				issueLabels    string
				issueTimeStats string
				milestoneID    int
			)
			if issue.Milestone != nil {
				milestoneID = issue.Milestone.ID
			}
			issueAssigneesByte, _ := json.Marshal(issue.Assignees)
			issueAssignees = string(issueAssigneesByte)
			issueLabelsByte, _ := json.Marshal(issue.Labels)
			issueLabels = string(issueLabelsByte)
			issueTimeStatsByte, _ := json.Marshal(issue.TimeStats)
			issueTimeStats = string(issueTimeStatsByte)
			issueDB := &model.StatisticsIssues{
				ProjectID:        projectID,
				IssueID:          issue.ID,
				IssueIID:         issue.IID,
				MilestoneID:      milestoneID,
				Description:      issue.Description,
				State:            issue.State,
				Assignees:        issueAssignees,
				Upvotes:          issue.Upvotes,
				Downvotes:        issue.Downvotes,
				Labels:           issueLabels,
				Title:            issue.Title,
				UpdatedAt:        issue.UpdatedAt,
				CreatedAt:        issue.CreatedAt,
				ClosedAt:         issue.ClosedAt,
				Subscribed:       issue.Subscribed,
				UserNotesCount:   issue.UserNotesCount,
				DueDate:          issue.DueDate,
				WebURL:           issue.WebURL,
				TimeStats:        issueTimeStats,
				Confidential:     issue.Confidential,
				Weight:           issue.Weight,
				DiscussionLocked: issue.DiscussionLocked,
				IssueLinkID:      issue.IssueLinkID,
			}

			if issue.Author != nil {
				issueDB.AuthorID = issue.Author.ID
				issueDB.AuthorName = issue.Author.Name
			}
			if issue.Assignee != nil {
				issueDB.AssigneeID = issue.Assignee.ID
				issueDB.AssigneeName = issue.Assignee.Name
			}

			if len(issueDB.Description) > model.MessageMaxLen {
				issueDB.Description = issueDB.Description[0 : model.MessageMaxLen-1]
			}

			if err = s.SaveDatabaseIssue(issueDB); err != nil {
				log.Error("issue Save Database err: projectID(%d), IssueID(%d)", projectID, issue.ID)
				err = nil
				continue
			}
			totalNum++
		}

		if resp.NextPage == 0 {
			break
		}
	}
	return
}

// SaveDatabaseIssue ...
func (s *Service) SaveDatabaseIssue(issueDB *model.StatisticsIssues) (err error) {
	var total int

	if total, err = s.dao.HasIssue(issueDB.ProjectID, issueDB.IssueID); err != nil {
		log.Error("SaveDatabaseIssue HasIssue(%+v)", err)
		return
	}

	// update found row
	if total == 1 {
		if err = s.dao.UpdateIssue(issueDB.ProjectID, issueDB.IssueID, issueDB); err != nil {
			if strings.Contains(err.Error(), model.DatabaseErrorText) {
				issueDB.Title = strconv.QuoteToASCII(issueDB.Title)
				issueDB.Description = strconv.QuoteToASCII(issueDB.Description)
				issueDB.Title = utils.Unicode2Chinese(issueDB.Title)
				issueDB.Description = utils.Unicode2Chinese(issueDB.Description)
			}
			if err = s.dao.UpdateIssue(issueDB.ProjectID, issueDB.IssueID, issueDB); err != nil {
				log.Error("SaveDatabaseIssue UpdateIssue(%+v)", err)
				return
			}
		}
		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabaseIssue issue has more rows(%d)", total)
		return
	}

	// insert row now
	if err = s.dao.CreateIssue(issueDB); err != nil {
		if strings.Contains(err.Error(), model.DatabaseErrorText) {
			issueDB.Title = strconv.QuoteToASCII(issueDB.Title)
			issueDB.Description = strconv.QuoteToASCII(issueDB.Description)
			issueDB.Title = utils.Unicode2Chinese(issueDB.Title)
			issueDB.Description = utils.Unicode2Chinese(issueDB.Description)
		}
		if err = s.dao.CreateIssue(issueDB); err != nil {
			log.Error("SaveDatabaseIssue CreateIssue(%+v)", err)
			return
		}
	}

	return
}

/*-------------------------------------- sync note ----------------------------------------*/

// SyncProjectNotes ...
func (s *Service) SyncProjectNotes(c context.Context, projectID int) (totalPage, totalNum int, err error) {
	var (
		//syncAllTime = conf.Conf.Property.SyncData.SyncAllTime
		syncAllTime = false
		resp        *gitlab.Response
		projectInfo *model.ProjectInfo
		mrs         []*model.StatisticsMrs
		notes       []*gitlab.Note
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
		for page := 1; ; page++ {
			totalPage++
			if notes, resp, err = s.gitlab.ListMRNotes(c, projectID, mr.MRIID, page); err != nil {
				return
			}
			for _, note := range notes {
				var (
					notePosition     string
					notePositionByte []byte
				)
				if notePositionByte, err = json.Marshal(note.Position); err != nil {
					return
				}
				notePosition = string(notePositionByte)
				noteDB := &model.StatisticsNotes{
					ProjectID:      projectID,
					ProjectName:    projectInfo.Name,
					MrIID:          mr.MRIID,
					NoteID:         note.ID,
					Body:           note.Body,
					Attachment:     note.Attachment,
					Title:          note.Title,
					FileName:       note.FileName,
					AuthorID:       note.Author.ID,
					AuthorName:     note.Author.Name,
					System:         note.System,
					ExpiresAt:      note.ExpiresAt,
					UpdatedAt:      note.UpdatedAt,
					CreatedAt:      note.CreatedAt,
					NoteableID:     note.NoteableID,
					NoteableType:   note.NoteableType,
					Position:       notePosition,
					Resolvable:     note.Resolvable,
					Resolved:       note.Resolved,
					ResolvedByID:   note.ResolvedBy.ID,
					ResolvedByName: note.ResolvedBy.Username,
					NoteableIID:    note.NoteableIID,
				}
				if err = s.SaveDatabaseNote(c, noteDB); err != nil {
					log.Error("note Save Database err: projectID(%d), NoteID(%d)", projectID, note.ID)
					err = nil
					continue
				}
				totalNum++
			}
			if resp.NextPage == 0 {
				break
			}
		}
	}
	return
}

// HandleNoteDBError ...
func (s *Service) HandleNoteDBError(c context.Context, noteDB *model.StatisticsNotes, errText string) (note *model.StatisticsNotes) {
	if strings.Contains(errText, model.DatabaseErrorText) {
		noteDB.Title = strconv.QuoteToASCII(noteDB.Title)
		noteDB.Body = strconv.QuoteToASCII(noteDB.Body)
		noteDB.Title = utils.Unicode2Chinese(noteDB.Title)
		noteDB.Body = utils.Unicode2Chinese(noteDB.Body)
		if len(noteDB.Body) > model.MessageMaxLen {
			noteDB.Body = noteDB.Body[0 : model.MessageMaxLen-1]
		}
	} else if strings.Contains(errText, model.DatabaseMaxLenthErrorText) {
		noteDB.Body = noteDB.Body[0 : model.MessageMaxLen-1]
	}
	return noteDB
}

// SaveDatabaseNote ...
func (s *Service) SaveDatabaseNote(c context.Context, noteDB *model.StatisticsNotes) (err error) {
	var total int

	if total, err = s.dao.HasNote(c, noteDB.ProjectID, noteDB.NoteID); err != nil {
		log.Error("SaveDatabaseNote HasNote(%+v)", err)
		return
	}

	// update found row
	if total == 1 {
		if err = s.dao.UpdateNote(c, noteDB.ProjectID, noteDB.NoteID, noteDB); err != nil {
			noteDB = s.HandleNoteDBError(c, noteDB, err.Error())
			if err = s.dao.UpdateNote(c, noteDB.ProjectID, noteDB.NoteID, noteDB); err != nil {
				log.Error("SaveDatabaseNote UpdateNote(%+v)", err)
				return
			}
		}
		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabaseNote Note has more rows(%d)", total)
		return
	}

	// insert row now
	if err = s.dao.CreateNote(c, noteDB); err != nil {
		noteDB = s.HandleNoteDBError(c, noteDB, err.Error())
		if err = s.dao.CreateNote(c, noteDB); err != nil {
			log.Error("SaveDatabaseNote CreateNote(%+v)", err)
			return
		}
	}

	return
}

/*-------------------------------------- sync emoji ----------------------------------------*/

// SyncProjectAwardEmoji ...
func (s *Service) SyncProjectAwardEmoji(c context.Context, projectID int) (totalPage, totalNum int, err error) {

	var (
		//syncAllTime = conf.Conf.Property.SyncData.SyncAllTime
		syncAllTime = false
		resp        *gitlab.Response
		projectInfo *model.ProjectInfo
		mrs         []*model.StatisticsMrs
		awardEmojis []*gitlab.AwardEmoji
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
		for page := 1; ; page++ {
			totalPage++
			if awardEmojis, resp, err = s.gitlab.ListMRAwardEmoji(projectID, mr.MRIID, page); err != nil {
				return
			}
			for _, awardEmoji := range awardEmojis {
				awardEmojiDB := &model.StatisticsMRAwardEmojis{
					ProjectID:     projectID,
					ProjectName:   projectInfo.Name,
					MrIID:         mr.MRIID,
					AwardEmojiID:  awardEmoji.ID,
					Name:          awardEmoji.Name,
					UserID:        awardEmoji.User.ID,
					UserName:      awardEmoji.User.Name,
					CreatedAt:     awardEmoji.CreatedAt,
					UpdatedAt:     awardEmoji.UpdatedAt,
					AwardableID:   awardEmoji.AwardableID,
					AwardableType: awardEmoji.AwardableType,
				}
				if err = s.SaveDatabaseAwardEmoji(c, awardEmojiDB); err != nil {
					log.Error("awardEmoji Save Database err: projectID(%d), AwardEmojiID(%d)", projectID, awardEmoji.ID)
					err = nil
					continue
				}
				totalNum++
			}
			if resp.NextPage == 0 {
				break
			}
		}
	}
	return
}

// SaveDatabaseAwardEmoji ...
func (s *Service) SaveDatabaseAwardEmoji(c context.Context, awardEmojiDB *model.StatisticsMRAwardEmojis) (err error) {
	var total int

	if total, err = s.dao.HasMRAwardEmoji(c, awardEmojiDB.ProjectID, awardEmojiDB.MrIID, awardEmojiDB.AwardEmojiID); err != nil {
		log.Error("SaveDatabaseAwardEmoji HasAwardEmoji(%+v)", err)
		return
	}

	// update found row
	if total == 1 {
		if err = s.dao.UpdateMRAwardEmoji(c, awardEmojiDB.ProjectID, awardEmojiDB.MrIID, awardEmojiDB.AwardEmojiID, awardEmojiDB); err != nil {
			log.Error("SaveDatabaseAwardEmoji UpdateAwardEmoji(%+v)", err)
		}
		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabaseAwardEmoji Note has more rows(%d)", total)
		return
	}

	// insert row now
	if err = s.dao.CreateMRAwardEmoji(c, awardEmojiDB); err != nil {
		log.Error("SaveDatabaseAwardEmoji CreateAwardEmoji(%+v)", err)
		return

	}

	return
}

/*-------------------------------------- sync Discussion ----------------------------------------*/

// SyncProjectDiscussion ...
func (s *Service) SyncProjectDiscussion(c context.Context, projectID int) (totalPage, totalNum int, err error) {

	var (
		//syncAllTime = conf.Conf.Property.SyncData.SyncAllTime
		syncAllTime = false
		resp        *gitlab.Response
		projectInfo *model.ProjectInfo
		mrs         []*model.StatisticsMrs
		discussions []*gitlab.Discussion
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
		for page := 1; ; page++ {
			totalPage++
			if discussions, resp, err = s.gitlab.ListMRDiscussions(projectID, mr.MRIID, page); err != nil {
				return
			}
			for _, discussion := range discussions {
				var (
					notesArray []int
					notesByte  []byte
					notes      string
				)
				if discussion.Notes != nil {
					for _, note := range discussion.Notes {
						notesArray = append(notesArray, note.ID)
					}
				}
				if notesByte, err = json.Marshal(notesArray); err != nil {
					notes = model.JsonMarshalErrorText
				} else {
					notes = string(notesByte)
				}
				discussionDB := &model.StatisticsDiscussions{
					ProjectID:      projectID,
					ProjectName:    projectInfo.Name,
					MrIID:          mr.MRIID,
					DiscussionID:   discussion.ID,
					IndividualNote: discussion.IndividualNote,
					Notes:          notes,
				}
				if err = s.SaveDatabaseDiscussion(c, discussionDB); err != nil {
					log.Error("discussion Save Database err: projectID(%d), DiscussionID(%s)", projectID, discussion.ID)
					err = nil
					continue
				}
				totalNum++
			}
			if resp.NextPage == 0 {
				break
			}
		}
	}
	return
}

// SaveDatabaseDiscussion ...
func (s *Service) SaveDatabaseDiscussion(c context.Context, discussionDB *model.StatisticsDiscussions) (err error) {
	var total int

	if total, err = s.dao.HasDiscussion(c, discussionDB.ProjectID, discussionDB.MrIID, discussionDB.DiscussionID); err != nil {
		log.Error("SaveDatabaseDiscussion HasDiscussion(%+v)", err)
		return
	}

	// update found row
	if total == 1 {
		if err = s.dao.UpdateDiscussion(c, discussionDB.ProjectID, discussionDB.MrIID, discussionDB.DiscussionID, discussionDB); err != nil {
			log.Error("SaveDatabaseDiscussion UpdateDiscussion(%+v)", err)
		}
		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabaseDiscussion Note has more rows(%d)", total)
		return
	}

	// insert row now
	if err = s.dao.CreateDiscussion(c, discussionDB); err != nil {
		log.Error("SaveDatabaseDiscussion CreateDiscussion(%+v)", err)
		return

	}

	return
}
