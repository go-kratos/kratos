package model

// type replyCommon struct {
// 	Message string      `json:"message"`
// 	Code    int         `json:"code"`
// 	Data    interface{} `json:"data"`
// }

// func createReply(data interface{}) *replyCommon {
// 	return &replyCommon{
// 		Message: "",
// 		Code:    0,
// 		Data:    data,
// 	}
// }

// func printReply(data interface{}) {
// 	var r = createReply(data)
// 	var result, _ = json.MarshalIndent(r, "", "    ")
// 	fmt.Printf(string(result) + "\n")
// }

// var now = xtime.Time(time.Now().Unix())

// func TestMcnGetMcnFansReply(t *testing.T) {
// 	var reply = McnGetMcnFansReply{Result: []*McnGetMcnFansInfo{
// 		{
// 			FansOverview: &dtmdl.DmConMcnFansD{LogDate: now},
// 			FansSex:      &dtmdl.DmConMcnFansSexW{LogDate: now},
// 			FansAge:      &dtmdl.DmConMcnFansAgeW{LogDate: now},
// 			FansPlayWay:  &dtmdl.DmConMcnFansPlayWayW{LogDate: now},
// 			FansArea:     []*dtmdl.DmConMcnFansAreaW{{LogDate: now}},
// 			FansType:     []*dtmdl.DmConMcnFansTypeW{{LogDate: now}},
// 			FansTag:      []*dtmdl.DmConMcnFansTagW{{LogDate: now}},
// 		},
// 	}}
// 	printReply(&reply)
// }
