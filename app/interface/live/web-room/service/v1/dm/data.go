package v1

import (
	"encoding/json"
	v1pb "go-common/app/interface/live/web-room/api/http/v1"
	"go-common/app/interface/live/web-room/model"
)

//HistoryData 历史数据处理
func HistoryData(data *v1pb.HistoryResp) map[string][]*model.History {
	var result = make(map[string][]*model.History)
	result["admin"] = make([]*model.History, 0, 10)
	result["room"] = make([]*model.History, 0, 10)

	for i := 0; i < len(data.Admin); i++ {
		var h = &model.History{}
		err := json.Unmarshal([]byte(data.Admin[i]), h)
		if err != nil {
			break
		}
		result["admin"] = append(result["admin"], h)
	}
	for i := 0; i < len(data.Room); i++ {
		var h = &model.History{}
		err := json.Unmarshal([]byte(data.Room[i]), h)
		if err != nil {
			break
		}
		result["room"] = append(result["room"], h)
	}
	return result
}
