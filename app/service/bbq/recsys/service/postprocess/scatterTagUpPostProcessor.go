package postprocess

import (
	"context"
	"go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"strings"
)

//ScatterTagUpProcessor ...
type ScatterTagUpProcessor struct {
	Processor
	lastScreenCount   int
	tagTotalLimit     int
	tagAdjacencyLimit int
	upTotalLimit      int
	upAdjacencyLimit  int
	lastRecordsType   string
	//adjacentQueueTagTotalLimit int
	//adjacentQueueTagAdjacencyLimit int
	//adjacentQueueTagUpTotalLimit int
	//adjacentQueueTagUpAdjacencyLimit int
}

func (p *ScatterTagUpProcessor) name() (name string) {
	name = "ScatterTagUp"
	return
}

func (p *ScatterTagUpProcessor) process(ctx context.Context, request *v1.RecsysRequest, response *v1.RecsysResponse, u *model.UserProfile) (err error) {
	if response.List == nil || len(response.List) == 0 {
		return
	}
	for _, record := range response.List {
		tagNameStr := record.Map[model.TagsName]
		tagTypeStr := record.Map[model.TagsType]
		tagNames := strings.Split(tagNameStr, "|")
		tagTypes := strings.Split(tagTypeStr, "|")
		for i, tagType := range tagTypes {
			if tagType == TagTypeZone2 {
				tagName := tagNames[i]
				//FIXME BBQ HACK
				if tagName == "宅舞" || tagName == "MMD·3D" || tagName == "三次元舞蹈" || tagName == "舞蹈教程" {
					tagName = "舞蹈"
				}
				record.Map[model.ScatterTag] = tagName
			}
		}
	}
	//{key, value[adjacentCnt, totalCnt]}
	//{key, totalCnt}
	tagCntMap := make(map[string]int)
	upCntMap := make(map[string]int)
	adjacentTag := make(map[string]int)
	adjacentUp := make(map[string]int)
	var lastRecords []model.Record4Dup
	if p.lastRecordsType == "lastRecords" {
		lastRecords = u.LastRecords
		//对于主推荐流程,保持4刷中没有重复up主,此处取前三刷,加上本轮的一刷.
		lastRecordsLength := len(lastRecords)
		for k := max(0, lastRecordsLength-3*p.lastScreenCount); k < lastRecordsLength; k++ {
			upMid := lastRecords[int64(k)].MID
			if _, ok := upCntMap[upMid]; ok {
				upCntMap[upMid]++
			} else {
				upCntMap[upMid] = 1
			}

		}
	} else {
		lastRecords = u.LastUpsRecords
	}
	if lastRecords != nil && len(lastRecords) > 0 {
		//记录上页最后一个结果
		lastRecord := lastRecords[len(lastRecords)-1]
		processAdjacentMap(adjacentTag, lastRecord.Tag)
		processAdjacentMap(adjacentUp, lastRecord.MID)
	}

	pageSize := int(request.Limit)
	originRecords := make([]*v1.RecsysRecord, len(response.List))
	copy(originRecords, response.List)
	excludedRecords := make([]*v1.RecsysRecord, 0)
	i := 0
	j := 0
	fromExcluded := false
	isSatisfiedOne := false
	record := originRecords[0]
	for i < len(originRecords) {
		if isSatisfiedOne && len(excludedRecords) > 0 {
			record = excludedRecords[0]
			fromExcluded = true
		} else {
			record = originRecords[i]
			i++
			fromExcluded = false
		}
		isSatisfiedOne = false

		if tagKey, ok := record.Map[model.ScatterTag]; ok {
			if upKey, ok := record.Map[model.UperMid]; ok {
				//是否同时满足
				if isSatisfied(tagKey, p.tagTotalLimit, p.tagAdjacencyLimit, tagCntMap, adjacentTag) {
					if isSatisfied(upKey, p.upTotalLimit, p.upAdjacencyLimit, upCntMap, adjacentUp) {
						response.List[j] = record
						j++
						isSatisfiedOne = true
						//只有添加时,才处理相邻关系的map
						processAdjacentMap(adjacentTag, tagKey)
						processAdjacentMap(adjacentUp, upKey)
						tagCntMap[tagKey]++
						upCntMap[upKey]++
						if fromExcluded {
							excludedRecords = excludedRecords[1:]
						}
					}
				}
			}
		}
		if !isSatisfiedOne {
			record.Map[model.OrderPostProcess] = p.name()
			if !fromExcluded {
				excludedRecords = append(excludedRecords, record)
			}
		}
		//结果个数满足,则跳出循环
		if j >= pageSize {
			break
		}
	}

	for k := 0; k < len(excludedRecords); k++ {
		response.List[j] = excludedRecords[k]
		j++
	}
	for ; i < len(originRecords); i++ {
		response.List[j] = originRecords[i]
		j++
	}
	return
}

func isSatisfied(findKey string, totalLimit int, adjacentLimit int, totalCntMap map[string]int, adjacentCntMap map[string]int) (isOk bool) {
	isOk = false
	if findKey == "" {
		return isOk
	}
	if total, ok := totalCntMap[findKey]; ok {
		if totalLimit <= 0 || total < totalLimit { // limit <=0 意味着没有限制
			if adjacentCnt, ok := adjacentCntMap[findKey]; ok {
				if adjacentLimit <= 0 || adjacentCnt < adjacentLimit { // limit <=0 意味着没有限制
					isOk = true
					return isOk
				}
			} else {
				isOk = true
				return isOk
			}
		}
	} else if adjacentCnt, ok := adjacentCntMap[findKey]; ok {
		if adjacentLimit <= 0 || adjacentCnt < adjacentLimit { // limit <=0 意味着没有限制
			isOk = true
			return isOk
		}
	} else {
		isOk = true
		return isOk
	}
	return isOk
}

func processAdjacentMap(srcMap map[string]int, key string) {
	if srcMap == nil {
		return
	}
	if _, ok := srcMap[key]; ok {
		srcMap[key]++
	} else {
		for k := range srcMap {
			delete(srcMap, k)
		}
		srcMap[key] = 1
	}
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
