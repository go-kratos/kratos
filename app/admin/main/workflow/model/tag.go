package model

// TagMeta tag from manager/list
type TagMeta struct {
	ID          int64  `json:"id"`
	Bid         int8   `json:"bid"`
	Tid         int64  `json:"tid"`
	TagID       int64  `json:"tag_id"` //map to old workflow tag id
	TName       string `json:"tname"`
	RID         int8   `json:"rid"`
	RName       string `json:"rname"`
	Name        string `json:"name"`
	Weight      int64  `json:"weight"`
	State       int8   `json:"state"`
	UID         int64  `json:"uid"`
	UName       string `json:"uname"`
	Description string `json:"description"`
	CTime       int    `json:"ctime"`
	MTime       int    `json:"mtime"`
}

// TagListResult .
type TagListResult struct {
	*CommonResponse
	Data struct {
		Tags []*TagMeta `json:"data"`
		Page `json:"page"`
	} `json:"data"`
}
