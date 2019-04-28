package service

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/ep/footman/model"
	"go-common/library/log"
)

//storyChangeByIteration get stories and story changes for input iterations
func (s *Service) storyChangeByIteration(iterationRes *model.IterationResponse, workspaceID, pattern string, sps, scps int) (iterations []*model.StoryChangeByIteration, err error) {
	var (
		storyRes       *model.StoryResponse
		storyChangeRes *model.StoryChangeResponse
	)
	iterations = make([]*model.StoryChangeByIteration, 0, 5)

	for _, itera := range iterationRes.Data {
		sci := &model.StoryChangeByIteration{}
		sci.IterationName = itera.Iteration.Name

		storyURL := fmt.Sprintf(model.StoryURL, workspaceID, itera.Iteration.ID, sps) + "&page=%d"
		log.Info("Story:" + fmt.Sprintf(model.StoryURL, workspaceID, itera.Iteration.ID, sps))
		if storyRes, err = s.dao.AllStories(sps, storyURL); err != nil {
			return
		}

		if storyRes == nil {
			continue
		}

		storyChanges := make([]*model.TargetStoryChange, 0, 5)

		if len(storyRes.Data) == 0 {
			continue
		}
		for _, story := range storyRes.Data {
			change := &model.TargetStoryChange{}
			change.Story = story.Story

			storyChangeURL := fmt.Sprintf(model.StoryChangeURL, workspaceID, change.Story.ID, scps) + "&page=%d"
			if storyChangeRes, err = s.dao.AllStoryChanges(scps, storyChangeURL); err != nil {
				return
			}

			if storyChangeRes == nil {
				continue
			}
			statusChanges := make([]*model.StatusChange, 0, 5)
			for _, change := range storyChangeRes.Data {
				reg := regexp.MustCompile(pattern)
				fields := reg.FindStringSubmatch(change.WorkitemChange.Changes)
				if len(fields) >= 1 {
					statusChange := &model.StatusChange{}
					statusChange.Created = change.WorkitemChange.Created
					statusChange.Creator = change.WorkitemChange.Creator
					statusChange.ValueBefore = fields[2]
					statusChange.ValueAfter = fields[4]
					statusChanges = append(statusChanges, statusChange)
				}
			}
			change.StatusChanges = statusChanges
			storyChanges = append(storyChanges, change)
		}

		sci.StoryChangeList = storyChanges
		sci.StoryCount = len(storyChanges)
		iterations = append(iterations, sci)
	}

	return
}

//storyChangeByIterationWithName get stories and story changes for input iteration name
func (s *Service) storyChangeByIterationWithName(workspaceID, iterationName, pattern string, ips, sps, scps int) (iterations []*model.StoryChangeByIteration, err error) {
	var (
		iterationRes *model.IterationResponse
	)
	iterationURL := fmt.Sprintf(model.IterationURL, workspaceID, iterationName, ips) + "&page=%d"
	log.Info("Iteration:" + fmt.Sprintf(model.IterationURL, workspaceID, iterationName, ips))
	if iterationRes, err = s.dao.AllIterations(ips, iterationURL); err != nil || len(iterationRes.Data) == 0 {
		return
	}
	return s.storyChangeByIteration(iterationRes, workspaceID, pattern, sps, scps)
}

//testTimeByIteration get test time data from story changes
func (s *Service) testTimeByIteration(scis []*model.StoryChangeByIteration, workspaceType, testType string) (iterations []*model.TestTimeByIteration, err error) {
	iterations = make([]*model.TestTimeByIteration, 0, 5)
	for _, sci := range scis {
		itera := &model.TestTimeByIteration{}
		itera.IterationName = sci.IterationName
		stories := make([]*model.TestTimeByStory, 0, 5)
		for _, story := range sci.StoryChangeList {
			ts := &model.TestTimeByStory{}
			ts.StoryName = story.Story.Name
			ts.StorySize = story.Story.Size
			ts.StoryEffort = story.Story.Effort

			found := 0
			var t1, t2 *time.Time
			var starttime, endtime time.Time
			total := 0.0
			for _, change := range story.StatusChanges {
				var condition1, condition2 bool
				if workspaceType == model.Android && testType == model.Test {
					condition1 = change.ValueBefore == "testing" && (change.ValueAfter == "developing" || change.ValueAfter == "status_6" || change.ValueAfter == "product_experience")
					condition2 = change.ValueBefore == "status_5" && change.ValueAfter == "testing" && found == 1
				}
				if workspaceType == model.IOS && testType == model.Test {
					condition1 = change.ValueBefore == "testing" && (change.ValueAfter == "developing" || change.ValueAfter == "status_7" || change.ValueAfter == "product_experience")
					condition2 = change.ValueBefore == "status_5" && change.ValueAfter == "testing" && found == 1
				}
				if testType == model.Experience {
					condition1 = change.ValueBefore == "product_experience" && (change.ValueAfter == "developing" || change.ValueAfter == "status_5")
					condition2 = change.ValueBefore == "developing" && change.ValueAfter == "product_experience" && found == 1
				}
				if condition1 {
					if starttime, err = s.stringToTime(change.Created); err != nil {
						return
					}
					found = 1
					t1 = &starttime
				} else if condition2 {
					if endtime, err = s.stringToTime(change.Created); err != nil {
						return
					}
					found = 2
					t2 = &endtime
				} else {
					found = 0
					t1 = nil
				}
				if t1 != nil && t2 != nil {
					total = total + t1.Sub(*t2).Hours() - float64(s.weekendDays(t1, t2)*24)
					t1 = nil
					t2 = nil
				}
			}
			if ts.TestTime, err = strconv.ParseFloat(strconv.FormatFloat(total, 'f', 2, 64), 64); err != nil {
				return
			}
			stories = append(stories, ts)
		}
		itera.TimeByStroy = stories
		itera.StoryCount = len(stories)
		iterations = append(iterations, itera)
	}
	return
}

//rejectedStoryByIteration get rejected stories
func (s *Service) rejectedStoryByIteration(scis []*model.StoryChangeByIteration, workspaceType, rejectType string) (iterations []*model.RejectedStoryByIteration, err error) {
	iterations = make([]*model.RejectedStoryByIteration, 0, 5)
	for _, sci := range scis {
		itera := &model.RejectedStoryByIteration{}
		itera.IterationName = sci.IterationName
		stories := make([]string, 0, 5)
		for _, story := range sci.StoryChangeList {
			for _, change := range story.StatusChanges {
				var condition bool
				if (workspaceType == model.Android || workspaceType == model.IOS) && rejectType == model.Test {
					condition = change.ValueBefore == "testing" && (change.ValueAfter == "developing" || change.ValueAfter == "product_experience")
				}
				if (workspaceType == model.Android || workspaceType == model.IOS) && rejectType == model.Experience {
					condition = change.ValueBefore == "product_experience" && change.ValueAfter == "developing"
				}
				if condition {
					stories = append(stories, story.Story.Name)
					break
				}
			}
		}
		itera.RejectedStoryList = stories
		itera.RejectedStoryCount = len(stories)
		iterations = append(iterations, itera)
	}
	return
}

//waitTimeByIteration get wait time data from story changes
func (s *Service) waitTimeByIteration(scis []*model.StoryChangeByIteration, waitType string) (iterations []*model.WaitTimeByIteration, err error) {
	iterations = make([]*model.WaitTimeByIteration, 0, 5)
	for _, sci := range scis {
		itera := &model.WaitTimeByIteration{}
		itera.IterationName = sci.IterationName
		stories := make([]*model.WaitTimeByStory, 0, 5)
		for _, story := range sci.StoryChangeList {
			ws := &model.WaitTimeByStory{}
			ws.StoryName = story.Story.Name
			ws.StorySize = story.Story.Size
			ws.StoryEffort = story.Story.Effort

			found := 0
			var t1, t2 *time.Time
			var starttime, endtime time.Time
			total := 0.0
			for _, change := range story.StatusChanges {
				var condition1, condition2 bool
				if waitType == model.Test {
					condition1 = change.ValueBefore == "status_5" && change.ValueAfter == "testing"
					condition2 = change.ValueBefore == "product_experience" && change.ValueAfter == "status_5" && found == 1
				}
				if condition1 {
					if starttime, err = s.stringToTime(change.Created); err != nil {
						return
					}
					found = 1
					t1 = &starttime
				} else if condition2 {
					if endtime, err = s.stringToTime(change.Created); err != nil {
						return
					}
					found = 2
					t2 = &endtime
				} else {
					found = 0
					t1 = nil
				}
				if t1 != nil && t2 != nil {
					total = t1.Sub(*t2).Hours() - float64(s.weekendDays(t1, t2)*24)
					break
				}
			}
			if ws.WaitTime, err = strconv.ParseFloat(strconv.FormatFloat(total, 'f', 2, 64), 64); err != nil {
				return
			}
			stories = append(stories, ws)
		}
		itera.TimeByStroy = stories
		itera.StoryCount = len(stories)
		iterations = append(iterations, itera)
	}
	return
}

//storyChangeItems get story change items
func (s *Service) storyChangeItems(workspaceID, storyID, pattern string, scps int) (changes []*model.StoryChangeItem, err error) {
	storyChangeURL := fmt.Sprintf(model.StoryChangeURL, workspaceID, storyID, scps) + "&page=%d"
	var storyChangeRes *model.StoryChangeResponse
	if storyChangeRes, err = s.dao.AllStoryChanges(scps, storyChangeURL); err != nil || storyChangeRes == nil {
		return
	}
	for k, v := range storyChangeRes.Data {
		reg := regexp.MustCompile(pattern)
		changeStrList := strings.Split(s.unicode2Chinese(v.WorkitemChange.Changes), "},{")
		for _, changeStr := range changeStrList {
			fields := reg.FindStringSubmatch(changeStr)
			if len(fields) >= 1 {
				change := &model.StoryChangeItem{}
				change.ID = v.WorkitemChange.ID
				change.ChangeSummay = v.WorkitemChange.ChangeSummay
				change.Comment = v.WorkitemChange.Comment
				change.EntityType = v.WorkitemChange.EntityType
				change.Number = strconv.Itoa(k)
				change.StoryID = v.WorkitemChange.StoryID
				change.WorkspaceID = v.WorkitemChange.WorkspaceID
				change.Created = v.WorkitemChange.Created
				change.Creator = v.WorkitemChange.Creator
				change.Field = fields[2]
				if fields[2] != "description" {
					change.ValueBefore = fields[4]
					change.ValueAfter = fields[6]
				}

				changes = append(changes, change)
			}
		}
	}
	return
}

func (s *Service) unicode2Chinese(str string) string {
	buf := bytes.NewBuffer(nil)
	i, j := 0, len(str)
	for i < j {
		x := i + 6
		if x > j {
			buf.WriteString(str[i:])
			break
		}
		if str[i] == '\\' && str[i+1] == 'u' {
			hex := str[i+2 : x]
			r, err := strconv.ParseUint(hex, 16, 64)
			if err == nil {
				buf.WriteRune(rune(r))
			} else {
				buf.WriteString(str[i:x])
			}
			i = x
		} else {
			buf.WriteByte(str[i])
			i++
		}
	}
	return buf.String()
}
