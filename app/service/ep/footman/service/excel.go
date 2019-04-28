package service

import (
	"fmt"
	"strconv"
	"strings"

	"go-common/app/service/ep/footman/model"
	"go-common/library/log"

	"github.com/360EntSecGroup-Skylar/excelize"
)

//testTimeExcel test time excel
func (s *Service) testTimeExcel(excelFilePath string, iteras []*model.TestTimeByIteration, iterationName, testType string) (err error) {
	log.Info("Start generating test time excel.")
	var xlsx *excelize.File
	if xlsx, err = excelize.OpenFile(excelFilePath); err != nil {
		xlsx = excelize.NewFile()
	}
	index := xlsx.NewSheet(iterationName)
	var tableHeader []string
	if testType == model.Test {
		tableHeader = []string{"Item", "需求", "迭代", "规模", "预估工时", "测试耗时(小时)"}
	}
	if testType == model.Experience {
		tableHeader = []string{"Item", "需求", "迭代", "规模", "预估工时", "产品验收耗时(小时)"}
	}

	for k, v := range tableHeader {
		xlsx.SetCellValue(iterationName, s.cellLocation(k, 1), v)
	}
	storyCount := 0
	for _, itera := range iteras {
		for _, story := range itera.TimeByStroy {
			storyCount++
			xlsx.SetCellValue(iterationName, s.cellLocation(0, storyCount+1), storyCount)
			xlsx.SetCellValue(iterationName, s.cellLocation(1, storyCount+1), story.StoryName)
			xlsx.SetCellValue(iterationName, s.cellLocation(2, storyCount+1), itera.IterationName)
			xlsx.SetCellValue(iterationName, s.cellLocation(3, storyCount+1), story.StorySize)
			xlsx.SetCellValue(iterationName, s.cellLocation(4, storyCount+1), story.StoryEffort)
			xlsx.SetCellValue(iterationName, s.cellLocation(5, storyCount+1), story.TestTime)
		}
	}
	xlsx.SetActiveSheet(index)
	log.Info("End generating test time excel.")
	return xlsx.SaveAs(excelFilePath)
}

//delayedStoryExcel delayed story list
func (s *Service) delayedStoryExcel(excelFilePath string, iteras []*model.StoryChangeByIteration, workspaceType, iterationName string, nameMap map[string]string) (err error) {
	log.Info("Start generating delayed story excel.")
	var xlsx *excelize.File
	if xlsx, err = excelize.OpenFile(excelFilePath); err != nil {
		xlsx = excelize.NewFile()
	}
	index := xlsx.NewSheet(iterationName)
	tableHeader := []string{"Item", "需求", "迭代", "规模", "预估工时", "状态", "修改前发布计划", "修改后发布计划", "修改人", "修改时间"}
	for k, v := range tableHeader {
		xlsx.SetCellValue(iterationName, s.cellLocation(k, 1), v)
	}
	storyCount := 0
	rn := 0
	var releases map[string]string
	if workspaceType == model.IOS {
		releases = model.IOSRelease
	}
	if workspaceType == model.Android {
		releases = model.AndroidRelease
	}

	for _, itera := range iteras {
		for _, story := range itera.StoryChangeList {
			for k, v := range story.StatusChanges {
				if _, ok := releases[v.ValueBefore]; ok {
					story.StatusChanges[k].ValueBefore = releases[v.ValueBefore]
				}
				if _, ok := releases[v.ValueAfter]; ok {
					story.StatusChanges[k].ValueAfter = releases[v.ValueAfter]
				}
				if v.ValueBefore == "0" {
					story.StatusChanges[k].ValueBefore = "--"
				}
				if v.ValueAfter == "" {
					story.StatusChanges[k].ValueAfter = "--"
				}
			}
			if len(story.StatusChanges) == 0 || (len(story.StatusChanges) == 1 && story.StatusChanges[0].ValueBefore == "--" && strings.Contains(story.StatusChanges[0].ValueAfter, iterationName)) {
				continue
			}
			if len(story.StatusChanges) > 1 {
				toSkip := 1
				for _, v := range story.StatusChanges {
					if v.ValueAfter != "--" && !strings.Contains(v.ValueAfter, iterationName) {
						toSkip = 0
						break
					}
				}
				if toSkip == 1 {
					continue
				}
			}
			storyCount++
			rn++
			xlsx.SetCellValue(iterationName, s.cellLocation(0, rn+1), storyCount)
			xlsx.SetCellValue(iterationName, s.cellLocation(1, rn+1), story.Story.Name)
			xlsx.SetCellValue(iterationName, s.cellLocation(2, rn+1), itera.IterationName)
			xlsx.SetCellValue(iterationName, s.cellLocation(3, rn+1), story.Story.Size)
			xlsx.SetCellValue(iterationName, s.cellLocation(4, rn+1), story.Story.Effort)
			xlsx.SetCellValue(iterationName, s.cellLocation(5, rn+1), nameMap[story.Story.Status])

			for _, c := range story.StatusChanges {
				xlsx.SetCellValue(iterationName, s.cellLocation(6, rn+1), c.ValueBefore)
				xlsx.SetCellValue(iterationName, s.cellLocation(7, rn+1), c.ValueAfter)
				xlsx.SetCellValue(iterationName, s.cellLocation(8, rn+1), c.Creator)
				xlsx.SetCellValue(iterationName, s.cellLocation(9, rn+1), c.Created)
				rn++
			}
			len := len(story.StatusChanges)
			xlsx.MergeCell(iterationName, s.cellLocation(0, rn-len+1), s.cellLocation(0, rn))
			xlsx.MergeCell(iterationName, s.cellLocation(1, rn-len+1), s.cellLocation(1, rn))
			xlsx.MergeCell(iterationName, s.cellLocation(2, rn-len+1), s.cellLocation(2, rn))
			xlsx.MergeCell(iterationName, s.cellLocation(3, rn-len+1), s.cellLocation(3, rn))
			xlsx.MergeCell(iterationName, s.cellLocation(4, rn-len+1), s.cellLocation(4, rn))
			xlsx.MergeCell(iterationName, s.cellLocation(5, rn-len+1), s.cellLocation(5, rn))
			rn--
		}
	}
	xlsx.SetActiveSheet(index)
	log.Info("End generating delayed story excel.")
	return xlsx.SaveAs(excelFilePath)
}

//rejectedStoryExcel rejected story list
func (s *Service) rejectedStoryExcel(excelFilePath string, iteras []*model.RejectedStoryByIteration, iterationName, rejectType string) (err error) {
	log.Info("Start generating rejected story excel.")
	var xlsx *excelize.File
	if xlsx, err = excelize.OpenFile(excelFilePath); err != nil {
		xlsx = excelize.NewFile()
	}
	index := xlsx.NewSheet(iterationName)
	var tableHeader []string
	if rejectType == model.Test {
		tableHeader = []string{"Item", "迭代", "测试打回需求数量", "测试打回需求列表"}
	}
	if rejectType == model.Experience {
		tableHeader = []string{"Item", "迭代", "产品打回需求数量", "产品打回需求列表"}
	}
	for k, v := range tableHeader {
		xlsx.SetCellValue(iterationName, s.cellLocation(k, 1), v)
	}
	iteraCount := 0
	rn := 0

	for _, itera := range iteras {

		if itera.RejectedStoryCount == 0 {
			continue
		}
		iteraCount++
		rn++
		xlsx.SetCellValue(iterationName, s.cellLocation(0, rn+1), iteraCount)
		xlsx.SetCellValue(iterationName, s.cellLocation(1, rn+1), itera.IterationName)
		xlsx.SetCellValue(iterationName, s.cellLocation(2, rn+1), itera.RejectedStoryCount)

		for _, story := range itera.RejectedStoryList {
			xlsx.SetCellValue(iterationName, s.cellLocation(3, rn+1), story)
			rn++
		}
		len := itera.RejectedStoryCount
		xlsx.MergeCell(iterationName, s.cellLocation(0, rn-len+1), s.cellLocation(0, rn))
		xlsx.MergeCell(iterationName, s.cellLocation(1, rn-len+1), s.cellLocation(1, rn))
		xlsx.MergeCell(iterationName, s.cellLocation(2, rn-len+1), s.cellLocation(2, rn))
		rn--
	}
	xlsx.SetActiveSheet(index)
	log.Info("End generating rejected story excel.")
	return xlsx.SaveAs(excelFilePath)
}

//waitTimeExcel wait time excel
func (s *Service) waitTimeExcel(excelFilePath string, iteras []*model.WaitTimeByIteration, iterationName, waitType string) (err error) {
	log.Info("Start generating waiting time excel.")
	var xlsx *excelize.File
	if xlsx, err = excelize.OpenFile(excelFilePath); err != nil {
		xlsx = excelize.NewFile()
	}
	index := xlsx.NewSheet(iterationName)
	var tableHeader []string
	if waitType == model.Test {
		tableHeader = []string{"Item", "需求", "迭代", "规模", "预估工时", "等待测试时长(小时)"}
	}

	for k, v := range tableHeader {
		xlsx.SetCellValue(iterationName, s.cellLocation(k, 1), v)
	}
	storyCount := 0
	for _, itera := range iteras {
		for _, story := range itera.TimeByStroy {
			storyCount++
			xlsx.SetCellValue(iterationName, s.cellLocation(0, storyCount+1), storyCount)
			xlsx.SetCellValue(iterationName, s.cellLocation(1, storyCount+1), story.StoryName)
			xlsx.SetCellValue(iterationName, s.cellLocation(2, storyCount+1), itera.IterationName)
			xlsx.SetCellValue(iterationName, s.cellLocation(3, storyCount+1), story.StorySize)
			xlsx.SetCellValue(iterationName, s.cellLocation(4, storyCount+1), story.StoryEffort)
			xlsx.SetCellValue(iterationName, s.cellLocation(5, storyCount+1), story.WaitTime)
		}
	}
	xlsx.SetActiveSheet(index)
	log.Info("End generating waiting time excel.")
	return xlsx.SaveAs(excelFilePath)
}

//storyWallExcel story wall excel
func (s *Service) storyWallExcel(excelFilePath string, iteras []*model.StoryChangeByIteration, workspaceID, workspaceType, iterationName string, isTime bool, nameMap map[string]string, categoryMap map[string]string) (err error) {
	if iterationName == "" {
		iterationName = "all"
	}

	log.Info("Start generating story wall excel.")
	var (
		xlsx                 *excelize.File
		workflow             []string
		releases             map[string]string
		additionalFields     map[string]string
		additionalFieldsList []string
	)
	if xlsx, err = excelize.OpenFile(excelFilePath); err != nil {
		xlsx = excelize.NewFile()
	}
	index := xlsx.NewSheet(iterationName)

	//init table header map
	xlsx.SetCellValue(iterationName, s.cellLocation(0, 1), "Item")
	xlsx.SetCellValue(iterationName, s.cellLocation(1, 1), "需求")
	xlsx.SetCellValue(iterationName, s.cellLocation(2, 1), "迭代")
	if workspaceType == model.Android {
		workflow = model.AndroidWorkflow
		releases = model.AndroidRelease
		additionalFields = model.AndroidFields
		additionalFieldsList = model.AndroidFieldsList
	}
	if workspaceType == model.IOS {
		workflow = model.IOSWorkflow
		releases = model.IOSRelease
		additionalFields = model.IOSFields
		additionalFieldsList = model.IOSFieldsList
	}

	if workspaceType == model.Live {
		workflow = model.LiveWorkflow
		//releases
		//additionalFields
		//additionalFieldsList

	}

	if workspaceType == model.BPlus {
		workflow = model.BPlusWorkflow
	}

	tableHeader := make(map[string]int)
	order := 3
	for _, v := range workflow {
		tableHeader[v] = order
		xlsx.SetCellValue(iterationName, s.cellLocation(order, 1), v)
		order++
	}
	for _, v := range model.BaseFieldsList {
		tableHeader[v] = order
		xlsx.SetCellValue(iterationName, s.cellLocation(order, 1), v)
		order++
	}

	if workspaceType == model.Android || workspaceType == model.IOS {
		for _, v := range additionalFieldsList {
			tableHeader[v] = order
			xlsx.SetCellValue(iterationName, s.cellLocation(order, 1), v)
			order++
		}
	}

	storyCount := 0
	for _, itera := range iteras {
		for _, story := range itera.StoryChangeList {
			storyCount++
			xlsx.SetCellValue(iterationName, s.cellLocation(0, storyCount+1), storyCount)
			xlsx.SetCellValue(iterationName, s.cellLocation(1, storyCount+1), story.Story.Name)
			xlsx.SetCellValue(iterationName, s.cellLocation(2, storyCount+1), itera.IterationName)

			var createdDate, completedDate string
			if createdDate, err = s.timeToDate(story.Story.Created); err != nil {
				return
			}
			if completedDate, err = s.timeToDate(story.Story.Completed); err != nil {
				completedDate = story.Story.Completed
			}

			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["id"]], storyCount+1), story.Story.ID)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["status"]], storyCount+1), nameMap[story.Story.Status])
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["created"]], storyCount+1), createdDate)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["priority"]], storyCount+1), story.Story.Priority)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["size"]], storyCount+1), story.Story.Size)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["category_id"]], storyCount+1), categoryMap[story.Story.CategoryID])

			if story.Story.ParentID == "0" {
				xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["parent_id"]], storyCount+1), "")
			} else {
				specStoryURL := fmt.Sprintf(model.SpecStoryURL, workspaceID, story.Story.ParentID)
				var specStoryRes *model.SpecStoryResponse
				if specStoryRes, err = s.dao.SpecStory(specStoryURL); err != nil {
					return
				}
				xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["parent_id"]], storyCount+1), specStoryRes.Data.Story.Name)
			}

			var releaseStr string
			if _, ok := releases[story.Story.ReleaseID]; ok {
				releaseStr = releases[story.Story.ReleaseID]
			} else {
				releaseStr = story.Story.ReleaseID
			}

			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["release_id"]], storyCount+1), releaseStr)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["owner"]], storyCount+1), story.Story.Owner)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["developer"]], storyCount+1), story.Story.Developer)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["creator"]], storyCount+1), story.Story.Creator)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["begin"]], storyCount+1), story.Story.Begin)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["due"]], storyCount+1), story.Story.Due)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["completed"]], storyCount+1), completedDate)
			xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[model.BaseFields["effort"]], storyCount+1), story.Story.Effort)

			if workspaceType == model.Android || workspaceType == model.IOS {
				xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[additionalFields["custom_field_99"]], storyCount+1), story.Story.CustomField99)
				xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[additionalFields["custom_field_97"]], storyCount+1), story.Story.CustomField97)
				xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[additionalFields["custom_field_93"]], storyCount+1), story.Story.CustomField93)

				if workspaceType == model.IOS {
					xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[additionalFields["custom_field_92"]], storyCount+1), story.Story.CustomField92)
				}
			}

			for _, change := range story.StatusChanges {
				var changeDate string
				if isTime {
					changeDate = change.Created
				} else if changeDate, err = s.timeToDate(change.Created); err != nil {
					return
				}
				if _, ok := tableHeader[nameMap[change.ValueAfter]]; !ok {
					continue
				}
				if xlsx.GetCellValue(iterationName, s.cellLocation(tableHeader[nameMap[change.ValueAfter]], storyCount+1)) == "" {
					xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[nameMap[change.ValueAfter]], storyCount+1), changeDate)
				}
			}

			if xlsx.GetCellValue(iterationName, s.cellLocation(tableHeader[workflow[0]], storyCount+1)) == "" {
				xlsx.SetCellValue(iterationName, s.cellLocation(tableHeader[workflow[0]], storyCount+1), createdDate)
			}
		}
	}
	xlsx.SetActiveSheet(index)
	log.Info("End generating story wall excel.")
	return xlsx.SaveAs(excelFilePath)
}

func (s *Service) storyWallReportExcel(excelFilePath, iterationName string, storyWallColNames, datesSort []string, storyWallColNameMap map[string]string, excelRet map[string]map[string]int) (err error) {
	log.Info("Start generating story report wall excel.")
	var (
		xlsx *excelize.File
	)
	if xlsx, err = excelize.OpenFile(excelFilePath); err != nil {
		xlsx = excelize.NewFile()
	}

	index := xlsx.NewSheet(iterationName)

	xlsx.SetCellValue(iterationName, s.cellLocation(0, 1), "Date")
	for k, v := range storyWallColNames {
		xlsx.SetCellValue(iterationName, s.cellLocation(k+1, 1), storyWallColNameMap[v])
	}

	storyCount := 0
	for _, dateSort := range datesSort {
		rowDate := excelRet[dateSort]
		if rowDate == nil {
			continue
		}
		storyCount = storyCount + 1
		xlsx.SetCellValue(iterationName, s.cellLocation(0, storyCount+1), dateSort)
		for k, v := range storyWallColNames {
			xlsx.SetCellValue(iterationName, s.cellLocation(k+1, storyCount+1), rowDate[v])
		}
	}

	xlsx.SetActiveSheet(index)
	log.Info("End generating story report wall excel.")
	return xlsx.SaveAs(excelFilePath)
}

//cellLocation get cell location with given column and row
func (s *Service) cellLocation(c int, r int) (cellLocation string) {
	sc := 'A'
	if c < 26 {
		cellLocation = string(int(sc)+c) + strconv.Itoa(r)
	} else {
		cellLocation = "A" + string(int(sc)+c-26) + strconv.Itoa(r)
	}
	return
}

func (s *Service) getFromExcel(excelName, sheetName string) (ret []map[string]string, err error) {
	var (
		xlsx      *excelize.File
		sheetRows [][]string
	)

	if xlsx, err = excelize.OpenFile(excelName); err != nil {
		return
	}

	sheetRows = xlsx.GetRows(sheetName)
	for index := 1; index < len(sheetRows); index++ {
		tmpMap := make(map[string]string)
		for innerIndex := 0; innerIndex < len(sheetRows[0]); innerIndex++ {
			key := sheetRows[0][innerIndex]
			val := sheetRows[index][innerIndex]
			tmpMap[key] = val
		}
		ret = append(ret, tmpMap)
	}
	return
}
