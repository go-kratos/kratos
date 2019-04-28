package service

import (
	"fmt"
	"strings"
	"time"

	"go-common/app/service/ep/footman/model"

	"github.com/pkg/errors"
)

const _dayMillisecond = 86400000000000

//TestTimeReport test time report
func (s *Service) TestTimeReport(workspaceID, workspaceType, iterationName, excelFilePath, testType string, ips, sps, scps int) (err error) {

	var (
		scis   []*model.StoryChangeByIteration
		iteras []*model.TestTimeByIteration
	)
	pattern := `("field":"status","value_before":")([\w]+)(","value_after":")([\w]+)(")`
	if scis, err = s.storyChangeByIterationWithName(workspaceID, iterationName, pattern, ips, sps, scps); err != nil {
		return errors.Wrap(err, "Error happened when fetching story changes from tapd")
	}

	if iteras, err = s.testTimeByIteration(scis, workspaceType, testType); err != nil {
		return errors.Wrap(err, "Error happened when extracting test time from story changes data")
	}

	return s.testTimeExcel(excelFilePath, iteras, iterationName, testType)
}

//DelayedStoryReport delayed story report
func (s *Service) DelayedStoryReport(workspaceID, workspaceType, iterationName, excelFilePath string, ips, sps, scps int) (err error) {

	var (
		scis       []*model.StoryChangeByIteration
		nameMapRes *model.NameMapResponse
	)
	pattern := `("field":"release_id","value_before":")([\w]*)(","value_after":")([\w]*)(")`
	if scis, err = s.storyChangeByIterationWithName(workspaceID, iterationName, pattern, ips, sps, scps); err != nil {
		return errors.Wrap(err, "Error happened when fetching story changes from tapd")
	}

	nameMapURL := fmt.Sprintf(model.NameMapURL, workspaceID)
	if nameMapRes, err = s.dao.NameMap(nameMapURL); err != nil {
		return errors.Wrap(err, "Error happened when fetching status name mapping from tapd")
	}
	return s.delayedStoryExcel(excelFilePath, scis, workspaceType, iterationName, nameMapRes.Data)
}

//RejectedStoryReport test or experience rejected story report
func (s *Service) RejectedStoryReport(workspaceID, workspaceType, iterationName, excelFilePath, rejectType string, ips, sps, scps int) (err error) {

	var (
		scis   []*model.StoryChangeByIteration
		iteras []*model.RejectedStoryByIteration
	)
	pattern := `("field":"status","value_before":")([\w]+)(","value_after":")([\w]+)(")`
	if scis, err = s.storyChangeByIterationWithName(workspaceID, iterationName, pattern, ips, sps, scps); err != nil {
		return errors.Wrap(err, "Error happened when fetching story changes from tapd")
	}

	if iteras, err = s.rejectedStoryByIteration(scis, workspaceType, rejectType); err != nil {
		return errors.Wrap(err, "Error happened when extracting rejected stories from story change data")
	}

	return s.rejectedStoryExcel(excelFilePath, iteras, iterationName, rejectType)
}

//WaitTimeReport test time report
func (s *Service) WaitTimeReport(workspaceID, iterationName, excelFilePath, waitType string, ips, sps, scps int) (err error) {

	var (
		scis   []*model.StoryChangeByIteration
		iteras []*model.WaitTimeByIteration
	)
	pattern := `("field":"status","value_before":")([\w]+)(","value_after":")([\w]+)(")`
	if scis, err = s.storyChangeByIterationWithName(workspaceID, iterationName, pattern, ips, sps, scps); err != nil {
		return errors.Wrap(err, "Error happened when fetching story changes from tapd")
	}

	if iteras, err = s.waitTimeByIteration(scis, waitType); err != nil {
		return errors.Wrap(err, "Error happened when extracting wait time from story change data")
	}

	return s.waitTimeExcel(excelFilePath, iteras, iterationName, waitType)
}

//StoryWallReport story wall report
func (s *Service) StoryWallReport(workspaceID, workspaceType, iterationName, excelFilePath string, isTime bool, ips, sps, scps, cps int) (err error) {

	var (
		scis        []*model.StoryChangeByIteration
		nameMapRes  *model.NameMapResponse
		categoryMap map[string]string
	)
	pattern := `("field":"status","value_before":")([\w]+)(","value_after":")([\w]+)(")`
	if scis, err = s.storyChangeByIterationWithName(workspaceID, iterationName, pattern, ips, sps, scps); err != nil {
		return errors.Wrap(err, "Error happened when fetching story changes from tapd")
	}

	nameMapURL := fmt.Sprintf(model.NameMapURL, workspaceID)
	if nameMapRes, err = s.dao.NameMap(nameMapURL); err != nil {
		return errors.Wrap(err, "Error happened when fetching status name mapping from tapd")
	}

	categoryURL := fmt.Sprintf(model.CategoryURL, workspaceID, model.CPS) + "&page=%d"
	if categoryMap, err = s.dao.AllCategories(cps, categoryURL); err != nil {
		return errors.Wrap(err, "Error happened when fetching project categories from tapd")
	}
	return s.storyWallExcel(excelFilePath, scis, workspaceID, workspaceType, iterationName, isTime, nameMapRes.Data, categoryMap)
}

// GenStoryReport Gen Story Report.
func (s *Service) GenStoryReport(workspaceType, importExcel, sheetName, exportExcel, iterationName, startDate, endDate string) (err error) {
	var (
		startTime     time.Time
		endTime       time.Time
		dateSort      []string
		excelRets     []map[string]string
		timeLayout    = "2006-01-02"
		tmpExportMap  = make(map[string][]*model.StoryWallTimeModel)
		exportRet     = make(map[string]map[string]int)
		storyFields   []string
		storyColNames map[string]string
	)

	switch workspaceType {
	case "android":
		storyFields = model.AndroidStoryWallFields
		storyColNames = model.AndroidStoryWallColNames

	case "ios":
		storyFields = model.IosStoryWallFields
		storyColNames = model.IosStoryWallColNames
	}

	loc, _ := time.LoadLocation("Local")

	if startTime, err = time.ParseInLocation(timeLayout, startDate, loc); err != nil {
		return
	}
	if endTime, err = time.ParseInLocation(timeLayout, endDate, loc); err != nil {
		return
	}

	if excelRets, err = s.getFromExcel(importExcel, sheetName); err != nil {
		return
	}

	for _, excelRet := range excelRets {
		if strings.Contains(excelRet["iterationname"], iterationName) {
			lastTime := endTime

			for index := len(storyFields); index > 0; index-- {
				//for index, stepField := range model.AndroidStoryWallFields {
				stepField := storyFields[index-1]
				if index == 0 {
					if strings.TrimSpace(excelRet[stepField]) != "" {
						curDate := strings.Split(excelRet[stepField], " ")[0]
						curTime, _ := time.ParseInLocation(timeLayout, curDate, loc)
						if !curTime.After(endTime) {
							tmpExportMap[stepField] = append(tmpExportMap[stepField], &model.StoryWallTimeModel{StepStartTime: curTime, StepEndTime: endTime})
							lastTime = curTime.AddDate(0, 0, -1)
						}
					}
				} else {
					if strings.TrimSpace(excelRet[stepField]) != "" {
						curDate := strings.Split(excelRet[stepField], " ")[0]
						curTime, _ := time.ParseInLocation(timeLayout, curDate, loc)
						tmpExportMap[stepField] = append(tmpExportMap[stepField], &model.StoryWallTimeModel{StepStartTime: curTime, StepEndTime: lastTime})
						lastTime = curTime.AddDate(0, 0, -1)
					}
				}
			}

		}
	}

	durDate := int(endTime.Sub(startTime)/_dayMillisecond + 1)

	for index := 0; index < durDate; index++ {
		curTime := startTime.AddDate(0, 0, index)
		curDate := curTime.Format("2006-01-02")
		dayRet := make(map[string]int)
		for stepKey, stepVals := range tmpExportMap {
			totalCount := 0
			for _, stepVal := range stepVals {
				if !stepVal.StepStartTime.After(curTime) && !curTime.After(stepVal.StepEndTime) {
					totalCount = totalCount + 1
				}
			}
			dayRet[stepKey] = totalCount
		}
		dateSort = append(dateSort, curDate)
		exportRet[curDate] = dayRet
	}

	s.storyWallReportExcel(exportExcel, iterationName, storyFields, dateSort, storyColNames, exportRet)
	return
}
