package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unicode"

	"go-common/app/service/ep/footman/model"
	"go-common/library/log"
)

const (
	_successTile   = "Success Notice"
	_successContxt = "Success to write %s Data To File,current retry time number: %d"
	_failedTile    = "Failed Notice"
	_failedContxt  = "Failed to write %s Data To File,current retry time number: %d, error:(+%v)"
)

//SaveFiles  save story, story change, iteration info to files
func (s *Service) SaveFiles(c context.Context) (err error) {
	return s.cache.Save(func() {
		s.saveFiles(context.Background())
	})
}

// SaveFilesForTask Save Files For Task.
func (s *Service) SaveFilesForTask() {
	s.cache.Save(func() {
		s.saveFiles(context.Background())
	})
}

//saveFiles  save story, story change, iteration info to files
func (s *Service) saveFiles(c context.Context) (err error) {

	for index := 0; index < s.c.Tapd.RetryTime; index++ {
		if err = s.writeBugDataToFile(c); err == nil {
			s.TapdMailNotice(_successTile, fmt.Sprintf(_successContxt, "bug", index+1))
			break
		} else {
			s.TapdMailNotice(_failedTile, fmt.Sprintf(_failedContxt, "bug", index+1, err))
		}
		time.Sleep(time.Duration(s.c.Tapd.WaitTime))
	}

	for index := 0; index < s.c.Tapd.RetryTime; index++ {
		if err = s.writeStoryDataToFile(c); err == nil {
			s.TapdMailNotice(_successTile, fmt.Sprintf(_successContxt, "story", index+1))
			break
		} else {
			s.TapdMailNotice(_failedTile, fmt.Sprintf(_failedContxt, "story", index+1, err))
		}
		time.Sleep(time.Duration(s.c.Tapd.WaitTime))
	}

	for index := 0; index < s.c.Tapd.RetryTime; index++ {
		if err = s.writeIterationDataToFile(c); err == nil {
			s.TapdMailNotice(_successTile, fmt.Sprintf(_successContxt, "iteration", index+1))
			break
		} else {
			s.TapdMailNotice(_failedTile, fmt.Sprintf(_failedContxt, "iteration", index+1, err))
		}
		time.Sleep(time.Duration(s.c.Tapd.WaitTime))
	}

	return
}

//DownloadStoryFile return the story file content
func (s *Service) DownloadStoryFile(c context.Context) (content []byte, err error) {
	return ioutil.ReadFile(s.c.Tapd.StoryFilePath)
}

//DownloadChangeFile return the story change file content
func (s *Service) DownloadChangeFile(c context.Context) (content []byte, err error) {
	return ioutil.ReadFile(s.c.Tapd.ChangeFilePath)
}

//DownloadIterationFile return the iteration file content
func (s *Service) DownloadIterationFile(c context.Context) (content []byte, err error) {
	return ioutil.ReadFile(s.c.Tapd.IterationFilePath)
}

//DownBugFile return the bug file content
func (s *Service) DownBugFile(c context.Context) (content []byte, err error) {
	return ioutil.ReadFile(s.c.Tapd.BugFilePath)
}

//writeStoryDataToFile write story and story change info to file
func (s *Service) writeStoryDataToFile(c context.Context) (err error) {
	var (
		storyRes              *model.StoryResponse
		storyFile, changeFile *os.File
		storyChangeItems      []*model.StoryChangeItem
	)

	log.Info("Writing story and story change info to file: start")

	writeTime := time.Now().Format("2006-01-02 15:04:05")

	if _, fileErr := os.Stat(s.c.Tapd.StoryFilePath); fileErr == nil {
		if err = os.Remove(s.c.Tapd.StoryFilePath); err != nil {
			log.Error("Error happened when removing story file. Error:(%v)", err)
			return
		}
	}
	if storyFile, err = os.Create(s.c.Tapd.StoryFilePath); err != nil {
		log.Error("Error happened when creating story file. Error:(%v)", err)
		return
	}
	defer storyFile.Close()

	if _, fileErr := os.Stat(s.c.Tapd.ChangeFilePath); fileErr == nil {
		if err = os.Remove(s.c.Tapd.ChangeFilePath); err != nil {
			log.Error("Error happened when removing story change file. Error:(%v)", err)
			return
		}
	}
	if changeFile, err = os.Create(s.c.Tapd.ChangeFilePath); err != nil {
		log.Error("Error happened when creating story change file. Error:(%v)", err)
		return
	}
	defer changeFile.Close()

	for _, workspaceID := range s.c.Tapd.StoryWorkspaceIDs {
		storyURL := fmt.Sprintf(model.StoryURLWithoutItera, workspaceID, s.c.Tapd.SPS) + "&page=%d"
		if storyRes, err = s.dao.AllStories(s.c.Tapd.SPS, storyURL); err != nil {
			log.Error("Error happened when fetching all stories for workspace %s. Error: %v", workspaceID, err)
			continue
		}
		for _, story := range storyRes.Data {
			//releaseId change to release name
			releaseName := s.handleReleaseName(story.Story.WorkspaceID, story.Story.ReleaseID)
			categoryName := s.handleCategoryPreName(story.Story.WorkspaceID, story.Story.CategoryID)

			customStoryFields := make(map[string]string)
			customStoryFields["CustomFieldOne"] = s.removeCRCFInString(story.Story.CustomFieldOne)
			customStoryFields["CustomFieldTwo"] = s.removeCRCFInString(story.Story.CustomFieldTwo)
			customStoryFields["CustomFieldThree"] = s.removeCRCFInString(story.Story.CustomFieldThree)
			customStoryFields["CustomFieldFour"] = s.removeCRCFInString(story.Story.CustomFieldFour)
			customStoryFields["CustomFieldFive"] = s.removeCRCFInString(story.Story.CustomFieldFive)
			customStoryFields["CustomFieldSix"] = s.removeCRCFInString(story.Story.CustomFieldSix)
			customStoryFields["CustomFieldSeven"] = s.removeCRCFInString(story.Story.CustomFieldSeven)
			customStoryFields["CustomFieldEight"] = s.removeCRCFInString(story.Story.CustomFieldEight)
			customStoryFields["CustomField9"] = s.removeCRCFInString(story.Story.CustomField9)
			customStoryFields["CustomField10"] = s.removeCRCFInString(story.Story.CustomField10)
			customStoryFields["CustomField11"] = s.removeCRCFInString(story.Story.CustomField11)
			customStoryFields["CustomField12"] = s.removeCRCFInString(story.Story.CustomField12)
			customStoryFields["CustomField13"] = s.removeCRCFInString(story.Story.CustomField13)
			customStoryFields["CustomField14"] = s.removeCRCFInString(story.Story.CustomField14)
			customStoryFields["CustomField15"] = s.removeCRCFInString(story.Story.CustomField15)
			customStoryFields["CustomField16"] = s.removeCRCFInString(story.Story.CustomField16)
			customStoryFields["CustomField17"] = s.removeCRCFInString(story.Story.CustomField17)
			customStoryFields["CustomField18"] = s.removeCRCFInString(story.Story.CustomField18)
			customStoryFields["CustomField19"] = s.removeCRCFInString(story.Story.CustomField19)
			customStoryFields["CustomField20"] = s.removeCRCFInString(story.Story.CustomField20)
			customStoryFields["CustomField21"] = s.removeCRCFInString(story.Story.CustomField21)
			customStoryFields["CustomField22"] = s.removeCRCFInString(story.Story.CustomField22)
			customStoryFields["CustomField23"] = s.removeCRCFInString(story.Story.CustomField23)
			customStoryFields["CustomField24"] = s.removeCRCFInString(story.Story.CustomField24)
			customStoryFields["CustomField25"] = s.removeCRCFInString(story.Story.CustomField25)
			customStoryFields["CustomField26"] = s.removeCRCFInString(story.Story.CustomField26)
			customStoryFields["CustomField27"] = s.removeCRCFInString(story.Story.CustomField27)
			customStoryFields["CustomField28"] = s.removeCRCFInString(story.Story.CustomField28)
			customStoryFields["CustomField29"] = s.removeCRCFInString(story.Story.CustomField29)
			customStoryFields["CustomField30"] = s.removeCRCFInString(story.Story.CustomField30)
			customStoryFields["CustomField31"] = s.removeCRCFInString(story.Story.CustomField31)
			customStoryFields["CustomField32"] = s.removeCRCFInString(story.Story.CustomField32)
			customStoryFields["CustomField33"] = s.removeCRCFInString(story.Story.CustomField33)
			customStoryFields["CustomField34"] = s.removeCRCFInString(story.Story.CustomField34)
			customStoryFields["CustomField35"] = s.removeCRCFInString(story.Story.CustomField35)
			customStoryFields["CustomField36"] = s.removeCRCFInString(story.Story.CustomField36)
			customStoryFields["CustomField37"] = s.removeCRCFInString(story.Story.CustomField37)
			customStoryFields["CustomField38"] = s.removeCRCFInString(story.Story.CustomField38)
			customStoryFields["CustomField39"] = s.removeCRCFInString(story.Story.CustomField39)
			customStoryFields["CustomField40"] = s.removeCRCFInString(story.Story.CustomField40)
			customStoryFields["CustomField41"] = s.removeCRCFInString(story.Story.CustomField41)
			customStoryFields["CustomField42"] = s.removeCRCFInString(story.Story.CustomField42)
			customStoryFields["CustomField43"] = s.removeCRCFInString(story.Story.CustomField43)
			customStoryFields["CustomField44"] = s.removeCRCFInString(story.Story.CustomField44)
			customStoryFields["CustomField45"] = s.removeCRCFInString(story.Story.CustomField45)
			customStoryFields["CustomField46"] = s.removeCRCFInString(story.Story.CustomField46)
			customStoryFields["CustomField47"] = s.removeCRCFInString(story.Story.CustomField47)
			customStoryFields["CustomField48"] = s.removeCRCFInString(story.Story.CustomField48)
			customStoryFields["CustomField49"] = s.removeCRCFInString(story.Story.CustomField49)
			customStoryFields["CustomField50"] = s.removeCRCFInString(story.Story.CustomField50)
			customStoryFields["CustomField51"] = s.removeCRCFInString(story.Story.CustomField51)
			customStoryFields["CustomField52"] = s.removeCRCFInString(story.Story.CustomField52)
			customStoryFields["CustomField53"] = s.removeCRCFInString(story.Story.CustomField53)
			customStoryFields["CustomField54"] = s.removeCRCFInString(story.Story.CustomField54)
			customStoryFields["CustomField55"] = s.removeCRCFInString(story.Story.CustomField55)
			customStoryFields["CustomField56"] = s.removeCRCFInString(story.Story.CustomField56)
			customStoryFields["CustomField57"] = s.removeCRCFInString(story.Story.CustomField57)
			customStoryFields["CustomField58"] = s.removeCRCFInString(story.Story.CustomField58)
			customStoryFields["CustomField59"] = s.removeCRCFInString(story.Story.CustomField59)
			customStoryFields["CustomField60"] = s.removeCRCFInString(story.Story.CustomField60)
			customStoryFields["CustomField61"] = s.removeCRCFInString(story.Story.CustomField61)
			customStoryFields["CustomField62"] = s.removeCRCFInString(story.Story.CustomField62)
			customStoryFields["CustomField63"] = s.removeCRCFInString(story.Story.CustomField63)
			customStoryFields["CustomField64"] = s.removeCRCFInString(story.Story.CustomField64)
			customStoryFields["CustomField65"] = s.removeCRCFInString(story.Story.CustomField65)
			customStoryFields["CustomField66"] = s.removeCRCFInString(story.Story.CustomField66)
			customStoryFields["CustomField67"] = s.removeCRCFInString(story.Story.CustomField67)
			customStoryFields["CustomField68"] = s.removeCRCFInString(story.Story.CustomField68)
			customStoryFields["CustomField69"] = s.removeCRCFInString(story.Story.CustomField69)
			customStoryFields["CustomField70"] = s.removeCRCFInString(story.Story.CustomField70)
			customStoryFields["CustomField71"] = s.removeCRCFInString(story.Story.CustomField71)
			customStoryFields["CustomField72"] = s.removeCRCFInString(story.Story.CustomField72)
			customStoryFields["CustomField73"] = s.removeCRCFInString(story.Story.CustomField73)
			customStoryFields["CustomField74"] = s.removeCRCFInString(story.Story.CustomField74)
			customStoryFields["CustomField75"] = s.removeCRCFInString(story.Story.CustomField75)
			customStoryFields["CustomField76"] = s.removeCRCFInString(story.Story.CustomField76)
			customStoryFields["CustomField77"] = s.removeCRCFInString(story.Story.CustomField77)
			customStoryFields["CustomField78"] = s.removeCRCFInString(story.Story.CustomField78)
			customStoryFields["CustomField79"] = s.removeCRCFInString(story.Story.CustomField79)
			customStoryFields["CustomField80"] = s.removeCRCFInString(story.Story.CustomField80)
			customStoryFields["CustomField81"] = s.removeCRCFInString(story.Story.CustomField81)
			customStoryFields["CustomField82"] = s.removeCRCFInString(story.Story.CustomField82)
			customStoryFields["CustomField83"] = s.removeCRCFInString(story.Story.CustomField83)
			customStoryFields["CustomField84"] = s.removeCRCFInString(story.Story.CustomField84)
			customStoryFields["CustomField85"] = s.removeCRCFInString(story.Story.CustomField85)
			customStoryFields["CustomField86"] = s.removeCRCFInString(story.Story.CustomField86)
			customStoryFields["CustomField87"] = s.removeCRCFInString(story.Story.CustomField87)
			customStoryFields["CustomField88"] = s.removeCRCFInString(story.Story.CustomField88)
			customStoryFields["CustomField89"] = s.removeCRCFInString(story.Story.CustomField89)
			customStoryFields["CustomField90"] = s.removeCRCFInString(story.Story.CustomField90)
			customStoryFields["CustomField91"] = s.removeCRCFInString(story.Story.CustomField91)
			customStoryFields["CustomField94"] = s.removeCRCFInString(story.Story.CustomField94)
			customStoryFields["CustomField95"] = s.removeCRCFInString(story.Story.CustomField95)
			customStoryFields["CustomField96"] = s.removeCRCFInString(story.Story.CustomField96)
			customStoryFields["CustomField98"] = s.removeCRCFInString(story.Story.CustomField98)
			customStoryFields["CustomField100"] = s.removeCRCFInString(story.Story.CustomField100)

			var jsonStr string
			if mJSON, jsonErr := json.Marshal(customStoryFields); jsonErr == nil {
				jsonStr = string(mJSON)
			}

			storyStr := s.removeCRCFInString(story.Story.ID+"\x01"+story.Story.Name+"\x01"+story.Story.WorkspaceID+"\x01"+
				story.Story.Creator+"\x01"+story.Story.Created+"\x01"+story.Story.Modified+"\x01"+
				story.Story.Status+"\x01"+story.Story.Owner+"\x01"+story.Story.Cc+"\x01"+
				story.Story.Begin+"\x01"+story.Story.Due+"\x01"+story.Story.Size+"\x01"+
				story.Story.Priority+"\x01"+story.Story.Developer+"\x01"+story.Story.IterationID+"\x01"+
				story.Story.TestFocus+"\x01"+story.Story.Type+"\x01"+story.Story.Source+"\x01"+
				story.Story.Module+"\x01"+story.Story.Version+"\x01"+story.Story.Completed+"\x01"+
				categoryName+"\x01"+story.Story.ParentID+"\x01"+story.Story.ChildrenID+"\x01"+
				story.Story.AncestorID+"\x01"+story.Story.BusinessValue+"\x01"+story.Story.Effort+"\x01"+
				story.Story.EffortCompleted+"\x01"+story.Story.Exceed+"\x01"+story.Story.Remain+"\x01"+
				releaseName+"\x01"+story.Story.CustomField92+"\x01"+story.Story.CustomField93+"\x01"+
				story.Story.CustomField97+"\x01"+story.Story.CustomField99+"\x01"+writeTime+"\x01"+jsonStr) + "\n"

			if _, err = io.WriteString(storyFile, storyStr); err != nil {
				log.Error("Error happened when writing story info to file. Error:(%v)", err)
				return
			}
			pattern := `("field":")(.*)(","value_before":")(.*)(","value_after":")(.*)(")`
			if storyChangeItems, err = s.storyChangeItems(workspaceID, story.Story.ID, pattern, s.c.Tapd.SCPS); err != nil {
				log.Error("Error happened when fetching story change items for story %s. Error:%v", story.Story.ID, err)
				return
			}
			for _, item := range storyChangeItems {
				itemBeforeVal := item.ValueBefore
				if s.isStringAsDigit(item.ValueBefore) {
					itemBeforeVal = s.handleReleaseNameForStoryChange(workspaceID, item.ValueBefore)
				}

				itemAfterVal := item.ValueAfter
				if s.isStringAsDigit(item.ValueAfter) {
					itemAfterVal = s.handleReleaseNameForStoryChange(workspaceID, item.ValueAfter)
				}

				itemStr := s.removeCRCFInString(item.ID+"\x01"+item.WorkspaceID+"\x01"+item.StoryID+"\x01"+
					item.Number+"\x01"+item.Field+"\x01"+item.Creator+"\x01"+
					item.Created+"\x01"+itemBeforeVal+"\x01"+itemAfterVal+"\x01"+
					item.ChangeSummay+"\x01"+item.Comment+"\x01"+item.EntityType+"\x01"+writeTime) + "\n"

				if _, err = io.WriteString(changeFile, itemStr); err != nil {
					log.Error("Error happened when writing story change item to file. Error:(%v)", err)
					return
				}
			}
		}
	}
	log.Info("Writing story and story change info to file: done")
	return
}

//writeIterationDataToFile write iteration info to file
func (s *Service) writeIterationDataToFile(c context.Context) (err error) {
	var (
		iterationRes  *model.IterationResponse
		iterationFile *os.File
	)

	log.Info("Writing iteration info to file: start")

	writeTime := time.Now().Format("2006-01-02 15:04:05")

	if _, fileErr := os.Stat(s.c.Tapd.IterationFilePath); fileErr == nil {
		if err = os.Remove(s.c.Tapd.IterationFilePath); err != nil {
			return
		}
	}
	if _, err = os.Create(s.c.Tapd.IterationFilePath); err != nil {
		return
	}

	if iterationFile, err = os.OpenFile(s.c.Tapd.IterationFilePath, os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		return
	}

	for _, workspaceID := range s.c.Tapd.IterationWorkspaceIDs {
		iterationURL := fmt.Sprintf(model.AllIterationURL, workspaceID, s.c.Tapd.IPS) + "&page=%d"
		if iterationRes, err = s.dao.AllIterations(s.c.Tapd.IPS, iterationURL); err != nil {
			log.Error("Error happened when fetching iterations from tapd. Error:(%dv)", err)
			continue
		}
		for _, iteration := range iterationRes.Data {
			customFields := make(map[string]string)
			customFields["CustomField1"] = s.removeCRCFInString(iteration.Iteration.CustomField1)
			customFields["CustomField2"] = s.removeCRCFInString(iteration.Iteration.CustomField2)
			customFields["CustomField3"] = s.removeCRCFInString(iteration.Iteration.CustomField3)
			customFields["CustomField4"] = s.removeCRCFInString(iteration.Iteration.CustomField4)
			customFields["CustomField5"] = s.removeCRCFInString(iteration.Iteration.CustomField5)
			customFields["CustomField6"] = s.removeCRCFInString(iteration.Iteration.CustomField6)
			customFields["CustomField7"] = s.removeCRCFInString(iteration.Iteration.CustomField7)
			customFields["CustomField8"] = s.removeCRCFInString(iteration.Iteration.CustomField8)
			customFields["CustomField9"] = s.removeCRCFInString(iteration.Iteration.CustomField9)
			customFields["CustomField10"] = s.removeCRCFInString(iteration.Iteration.CustomField10)
			customFields["CustomField11"] = s.removeCRCFInString(iteration.Iteration.CustomField11)
			customFields["CustomField12"] = s.removeCRCFInString(iteration.Iteration.CustomField12)
			customFields["CustomField13"] = s.removeCRCFInString(iteration.Iteration.CustomField13)
			customFields["CustomField14"] = s.removeCRCFInString(iteration.Iteration.CustomField14)
			customFields["CustomField15"] = s.removeCRCFInString(iteration.Iteration.CustomField15)
			customFields["CustomField16"] = s.removeCRCFInString(iteration.Iteration.CustomField16)
			customFields["CustomField17"] = s.removeCRCFInString(iteration.Iteration.CustomField17)
			customFields["CustomField18"] = s.removeCRCFInString(iteration.Iteration.CustomField18)
			customFields["CustomField19"] = s.removeCRCFInString(iteration.Iteration.CustomField19)
			customFields["CustomField20"] = s.removeCRCFInString(iteration.Iteration.CustomField20)
			customFields["CustomField21"] = s.removeCRCFInString(iteration.Iteration.CustomField21)
			customFields["CustomField22"] = s.removeCRCFInString(iteration.Iteration.CustomField22)
			customFields["CustomField23"] = s.removeCRCFInString(iteration.Iteration.CustomField23)
			customFields["CustomField24"] = s.removeCRCFInString(iteration.Iteration.CustomField24)
			customFields["CustomField25"] = s.removeCRCFInString(iteration.Iteration.CustomField25)
			customFields["CustomField26"] = s.removeCRCFInString(iteration.Iteration.CustomField26)
			customFields["CustomField27"] = s.removeCRCFInString(iteration.Iteration.CustomField27)
			customFields["CustomField28"] = s.removeCRCFInString(iteration.Iteration.CustomField28)
			customFields["CustomField29"] = s.removeCRCFInString(iteration.Iteration.CustomField29)
			customFields["CustomField30"] = s.removeCRCFInString(iteration.Iteration.CustomField30)
			customFields["CustomField31"] = s.removeCRCFInString(iteration.Iteration.CustomField31)
			customFields["CustomField32"] = s.removeCRCFInString(iteration.Iteration.CustomField32)
			customFields["CustomField33"] = s.removeCRCFInString(iteration.Iteration.CustomField33)
			customFields["CustomField34"] = s.removeCRCFInString(iteration.Iteration.CustomField34)
			customFields["CustomField35"] = s.removeCRCFInString(iteration.Iteration.CustomField35)
			customFields["CustomField36"] = s.removeCRCFInString(iteration.Iteration.CustomField36)
			customFields["CustomField37"] = s.removeCRCFInString(iteration.Iteration.CustomField37)
			customFields["CustomField38"] = s.removeCRCFInString(iteration.Iteration.CustomField38)
			customFields["CustomField39"] = s.removeCRCFInString(iteration.Iteration.CustomField39)
			customFields["CustomField40"] = s.removeCRCFInString(iteration.Iteration.CustomField40)
			customFields["CustomField41"] = s.removeCRCFInString(iteration.Iteration.CustomField41)
			customFields["CustomField42"] = s.removeCRCFInString(iteration.Iteration.CustomField42)
			customFields["CustomField43"] = s.removeCRCFInString(iteration.Iteration.CustomField43)
			customFields["CustomField44"] = s.removeCRCFInString(iteration.Iteration.CustomField44)
			customFields["CustomField45"] = s.removeCRCFInString(iteration.Iteration.CustomField45)
			customFields["CustomField46"] = s.removeCRCFInString(iteration.Iteration.CustomField46)
			customFields["CustomField47"] = s.removeCRCFInString(iteration.Iteration.CustomField47)
			customFields["CustomField48"] = s.removeCRCFInString(iteration.Iteration.CustomField48)
			customFields["CustomField49"] = s.removeCRCFInString(iteration.Iteration.CustomField49)
			customFields["CustomField50"] = s.removeCRCFInString(iteration.Iteration.CustomField50)

			var jsonStr string
			if mJSON, jsonErr := json.Marshal(customFields); jsonErr == nil {
				jsonStr = string(mJSON)
			}

			iterationStr := s.removeCRCFInString(iteration.Iteration.ID+"\x01"+
				iteration.Iteration.Name+"\x01"+
				iteration.Iteration.WorkspaceID+"\x01"+
				iteration.Iteration.Startdate+"\x01"+
				iteration.Iteration.Enddate+"\x01"+
				iteration.Iteration.Status+"\x01"+
				iteration.Iteration.ReleaseID+"\x01"+
				iteration.Iteration.Description+"\x01"+
				iteration.Iteration.Creator+"\x01"+
				iteration.Iteration.Created+"\x01"+
				iteration.Iteration.Modified+"\x01"+
				iteration.Iteration.Completed+"\x01"+
				writeTime+"\x01"+jsonStr) + "\n"
			if _, err = io.WriteString(iterationFile, iterationStr); err != nil {
				log.Error("Error happened when writing iteration into to file. Error:(%v)", err)
				return
			}
		}
	}

	log.Info("Writing iteration info to file: done")
	return
}

// writeBugDataToFile write Bug Data To File.
func (s *Service) writeBugDataToFile(c context.Context) (err error) {
	var (
		bugRes    *model.BugResponse
		bugFile   *os.File
		writeTime = time.Now().Format("2006-01-02 15:04:05")
	)

	log.Info("Writing bug info to file: start")

	if _, fileErr := os.Stat(s.c.Tapd.BugFilePath); fileErr == nil {
		if err = os.Remove(s.c.Tapd.BugFilePath); err != nil {
			log.Error("Error happened when removing bug file. Error:(%v)", err)
			return
		}
	}

	if bugFile, err = os.Create(s.c.Tapd.BugFilePath); err != nil {
		log.Error("Error happened when creating bug file. Error:(%v)", err)
		return
	}
	defer bugFile.Close()

	for _, workspaceID := range s.c.Tapd.BugWorkspaceIDs {
		bugURL := fmt.Sprintf(model.BugURL, workspaceID, s.c.Tapd.SPS) + "&page=%d"
		if bugRes, err = s.dao.AllBugs(s.c.Tapd.SPS, bugURL); err != nil {
			log.Error("Error happened when fetching all bugs for workspace %s. Error: %v", workspaceID, err)
			continue
		}

		for _, bug := range bugRes.Data {
			//releaseId change to release name
			releaseName := s.handleReleaseName(bug.Bug.WorkspaceID, bug.Bug.ReleaseID)
			customBugFields := make(map[string]string)
			customBugFields["CustomField10"] = s.removeCRCFInString(bug.Bug.CustomField10)
			customBugFields["CustomField11"] = s.removeCRCFInString(bug.Bug.CustomField11)
			customBugFields["CustomField12"] = s.removeCRCFInString(bug.Bug.CustomField12)
			customBugFields["CustomField13"] = s.removeCRCFInString(bug.Bug.CustomField13)
			customBugFields["CustomField14"] = s.removeCRCFInString(bug.Bug.CustomField14)
			customBugFields["CustomField15"] = s.removeCRCFInString(bug.Bug.CustomField15)
			customBugFields["CustomField16"] = s.removeCRCFInString(bug.Bug.CustomField16)
			customBugFields["CustomField17"] = s.removeCRCFInString(bug.Bug.CustomField17)
			customBugFields["CustomField18"] = s.removeCRCFInString(bug.Bug.CustomField18)
			customBugFields["CustomField19"] = s.removeCRCFInString(bug.Bug.CustomField19)
			customBugFields["CustomField20"] = s.removeCRCFInString(bug.Bug.CustomField20)
			customBugFields["CustomField21"] = s.removeCRCFInString(bug.Bug.CustomField21)
			customBugFields["CustomField22"] = s.removeCRCFInString(bug.Bug.CustomField22)
			customBugFields["CustomField23"] = s.removeCRCFInString(bug.Bug.CustomField23)
			customBugFields["CustomField24"] = s.removeCRCFInString(bug.Bug.CustomField24)
			customBugFields["CustomField25"] = s.removeCRCFInString(bug.Bug.CustomField25)
			customBugFields["CustomField26"] = s.removeCRCFInString(bug.Bug.CustomField26)
			customBugFields["CustomField27"] = s.removeCRCFInString(bug.Bug.CustomField27)
			customBugFields["CustomField28"] = s.removeCRCFInString(bug.Bug.CustomField28)
			customBugFields["CustomField29"] = s.removeCRCFInString(bug.Bug.CustomField29)
			customBugFields["CustomField30"] = s.removeCRCFInString(bug.Bug.CustomField30)
			customBugFields["CustomField31"] = s.removeCRCFInString(bug.Bug.CustomField31)
			customBugFields["CustomField32"] = s.removeCRCFInString(bug.Bug.CustomField32)
			customBugFields["CustomField33"] = s.removeCRCFInString(bug.Bug.CustomField33)
			customBugFields["CustomField34"] = s.removeCRCFInString(bug.Bug.CustomField34)
			customBugFields["CustomField35"] = s.removeCRCFInString(bug.Bug.CustomField35)
			customBugFields["CustomField36"] = s.removeCRCFInString(bug.Bug.CustomField36)
			customBugFields["CustomField37"] = s.removeCRCFInString(bug.Bug.CustomField37)
			customBugFields["CustomField38"] = s.removeCRCFInString(bug.Bug.CustomField38)
			customBugFields["CustomField39"] = s.removeCRCFInString(bug.Bug.CustomField39)
			customBugFields["CustomField40"] = s.removeCRCFInString(bug.Bug.CustomField40)
			customBugFields["CustomField41"] = s.removeCRCFInString(bug.Bug.CustomField41)
			customBugFields["CustomField42"] = s.removeCRCFInString(bug.Bug.CustomField42)
			customBugFields["CustomField43"] = s.removeCRCFInString(bug.Bug.CustomField43)
			customBugFields["CustomField44"] = s.removeCRCFInString(bug.Bug.CustomField44)
			customBugFields["CustomField45"] = s.removeCRCFInString(bug.Bug.CustomField45)
			customBugFields["CustomField46"] = s.removeCRCFInString(bug.Bug.CustomField46)
			customBugFields["CustomField47"] = s.removeCRCFInString(bug.Bug.CustomField47)
			customBugFields["CustomField48"] = s.removeCRCFInString(bug.Bug.CustomField48)
			customBugFields["CustomField49"] = s.removeCRCFInString(bug.Bug.CustomField49)
			customBugFields["CustomField50"] = s.removeCRCFInString(bug.Bug.CustomField50)

			var jsonStr string
			if mJSON, jsonErr := json.Marshal(customBugFields); jsonErr == nil {
				jsonStr = string(mJSON)
			}

			bugStr := s.removeCRCFInString(bug.Bug.ID+"\x01"+bug.Bug.Title+"\x01"+bug.Bug.Description+"\x01"+bug.Bug.Priority+"\x01"+
				bug.Bug.Severity+"\x01"+bug.Bug.Module+"\x01"+bug.Bug.Status+"\x01"+bug.Bug.Reporter+"\x01"+
				bug.Bug.Deadline+"\x01"+bug.Bug.Created+"\x01"+bug.Bug.BugType+"\x01"+bug.Bug.Resolved+"\x01"+
				bug.Bug.Closed+"\x01"+bug.Bug.Modified+"\x01"+bug.Bug.LastModify+"\x01"+bug.Bug.Auditer+"\x01"+
				bug.Bug.DE+"\x01"+bug.Bug.VersionTest+"\x01"+bug.Bug.VersionReport+"\x01"+bug.Bug.VersionClose+"\x01"+
				bug.Bug.VersionFix+"\x01"+bug.Bug.BaselineFind+"\x01"+bug.Bug.BaselineJoin+"\x01"+bug.Bug.BaselineClose+"\x01"+
				bug.Bug.BaselineTest+"\x01"+bug.Bug.SourcePhase+"\x01"+bug.Bug.TE+"\x01"+bug.Bug.CurrentOwner+"\x01"+
				bug.Bug.IterationID+"\x01"+bug.Bug.Resolution+"\x01"+bug.Bug.Source+"\x01"+bug.Bug.OriginPhase+"\x01"+
				bug.Bug.Confirmer+"\x01"+bug.Bug.Milestone+"\x01"+bug.Bug.Participator+"\x01"+bug.Bug.Closer+"\x01"+
				bug.Bug.Platform+"\x01"+bug.Bug.OS+"\x01"+bug.Bug.TestType+"\x01"+bug.Bug.TestPhase+"\x01"+
				bug.Bug.Frequency+"\x01"+bug.Bug.CC+"\x01"+bug.Bug.RegressionNumber+"\x01"+bug.Bug.Flows+"\x01"+
				bug.Bug.Feature+"\x01"+bug.Bug.TestMode+"\x01"+bug.Bug.Estimate+"\x01"+bug.Bug.IssueID+"\x01"+
				bug.Bug.CreatedFrom+"\x01"+bug.Bug.InProgressTime+"\x01"+bug.Bug.VerifyTime+"\x01"+bug.Bug.RejectTime+"\x01"+
				bug.Bug.ReopenTime+"\x01"+bug.Bug.AuditTime+"\x01"+bug.Bug.SuspendTime+"\x01"+bug.Bug.Due+"\x01"+
				bug.Bug.Begin+"\x01"+releaseName+"\x01"+bug.Bug.WorkspaceID+"\x01"+bug.Bug.CustomFieldOne+"\x01"+
				bug.Bug.CustomFieldTwo+"\x01"+bug.Bug.CustomFieldThree+"\x01"+bug.Bug.CustomFieldFour+"\x01"+bug.Bug.CustomFieldFive+"\x01"+
				bug.Bug.CustomField6+"\x01"+bug.Bug.CustomField7+"\x01"+bug.Bug.CustomField8+"\x01"+bug.Bug.CustomField9+"\x01"+
				writeTime+"\x01"+jsonStr) + "\n"

			if _, err = io.WriteString(bugFile, bugStr); err != nil {
				log.Error("Error happened when writing bug info to file. Error:(%v)", err)
				return
			}
		}
	}
	log.Info("Writing bug info to file: done")
	return
}

// handleReleaseName handle Release Name
func (s *Service) handleReleaseName(workSpaceID, releaseID string) (releaseName string) {
	if len(releaseID) < 10 {
		return releaseID
	}

	var err error
	if releaseID == "0" || strings.TrimSpace(releaseID) == "" {
		releaseName = ""
		return
	}

	if releaseName, err = s.dao.ReleaseName(workSpaceID, releaseID); err != nil {
		if workSpaceID == "20060791" {
			if releaseName, err = s.dao.ReleaseName("20055921", releaseID); err != nil {
				log.Error("Error happened when fetch release for replace workspace:%s, release id: %s, err %v", "20055921", releaseID, err)
				releaseName = releaseID
			}
		} else {
			log.Error("Error happened when fetch release for workspace:%s, release id: %s, err %v", workSpaceID, releaseID, err)
			releaseName = releaseID
		}
	}
	return
}

// handleReleaseName handle Release Name
func (s *Service) handleReleaseNameForStoryChange(workSpaceID, releaseID string) (releaseName string) {
	if len(releaseID) < 10 {
		return releaseID
	}

	var err error
	if releaseID == "0" || strings.TrimSpace(releaseID) == "" {
		releaseName = ""
		return
	}

	if releaseName, err = s.dao.ReleaseName(workSpaceID, releaseID); err != nil {
		if workSpaceID == "20060791" {
			if releaseName, err = s.dao.ReleaseName("20055921", releaseID); err != nil {
				log.Error("Error happened when fetch release for replace workspace:%s, release id: %s, err %v", "20055921", releaseID, err)
				releaseName = releaseID
			} else {
				releaseName = "发布计划#" + releaseName
			}
		} else {
			log.Error("Error happened when fetch release for workspace:%s, release id: %s, err %v", workSpaceID, releaseID, err)
			releaseName = releaseID
		}
	} else {
		releaseName = "发布计划#" + releaseName
	}
	return
}

// handleCategoryPreName handle CategoryPre Name
func (s *Service) handleCategoryPreName(workSpaceID, categoryID string) (categoryName string) {
	var err error
	if len(categoryID) < 10 {
		return categoryID
	}

	if categoryID == "0" || strings.TrimSpace(categoryID) == "" {
		categoryName = ""
		return
	}

	if categoryName, err = s.dao.CategoryPreName(workSpaceID, categoryID); err != nil {
		log.Error("Error happened when fetch release for workspace:%s, category id: %s, err %v", workSpaceID, categoryID, err)
		categoryName = categoryID
	}
	return
}

// isStringAsDigit is String As Digit.
func (s *Service) isStringAsDigit(str string) bool {
	if len(str) < 10 {
		return false
	}

	for _, c := range str {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
