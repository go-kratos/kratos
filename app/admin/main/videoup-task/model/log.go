package model

const (
	//LogClientVideo 视频business id
	LogClientVideo = int(2)
	//LogClientTypeVideo 视频 type id
	LogClientTypeVideo = int(1)

	//LogClientConsumer 一审任务 business id
	LogClientConsumer = int(131)
	//LogClientTypeConsumer 一审任务type id
	LogClientTypeConsumer = int(1)
)

// SearchLogResult is.
type SearchLogResult struct {
	Code int `json:"code"`
	Data struct {
		Order  string `json:"order"`
		Sort   string `json:"sort"`
		Result []struct {
			UID    int64  `json:"uid"`
			Uname  string `json:"uname"`
			OID    int64  `json:"oid"`
			Type   int8   `json:"type"`
			Action string `json:"action"`
			Str0   string `json:"str_0"`
			Str1   string `json:"str_1"`
			Str2   string `json:"str_2"`
			Int0   int    `json:"int_0"`
			Int1   int    `json:"int_1"`
			Int2   int    `json:"int_2"`
			Ctime  string `json:"ctime"`
			Extra  string `json:"extra_data"`
		} `json:"result"`
		Page struct {
			Num   int `json:"num"`
			Size  int `json:"size"`
			Total int `json:"total"`
		} `json:"page"`
	} `json:"data"`
}
