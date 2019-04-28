package model

// ArgGetQusBank 题库
type ArgGetQusBank struct {
	QsBId int64 `form:"qb_id" validate:"required"`
	Stat  int64 `form:"idel" validate:"min=0"`
}

// ArgAddQusBank 添加题库
type ArgAddQusBank struct {
	QBName       string `json:"qb_name" validate:"required"`
	CdTime       int64  `json:"cd_time" validate:"required,min=1,gte=1"`
	MaxRetryTime int64  `json:"max_retry_time" validate:"required,min=1,gte=1"`
}

// ArgBaseBank 基本题库
type ArgBaseBank struct {
	QsBId int64 `json:"qb_id" validate:"required"`
}

// ArgUpdateQusBank 更新题库
type ArgUpdateQusBank struct {
	QsBId        int64  `json:"qb_id" validate:"required"`
	QBName       string `json:"qb_name"`
	CdTime       int64  `json:"cd_time" validate:"min=1,gte=1"`
	MaxRetryTime int64  `json:"max_retry_time" validate:"min=1,gte=1"`
}

// ArgPage 分页
type ArgPage struct {
	PageNo   int `form:"page" validate:"required,min=1,gte=1"`
	PageSize int `form:"page_size" validate:"required,min=1,gte=1"`
}

// ArgBankList 题库列表
type ArgBankList struct {
	ArgPage
	Name string `form:"key"`
}

// ArgGetQus 题目
type ArgGetQus struct {
	QsID int64 `form:"qid" validate:"required"`
	Stat int64 `form:"status" validate:"min=0"`
}

// ArgBaseQus 基本题目信息
type ArgBaseQus struct {
	Type   int8     `json:"question_type" validate:"required,min=1,max=4"`
	AnType int8     `json:"answer_type" validate:"required,min=1,max=4"`
	Name   string   `json:"question_name" validate:"required,min=1,gte=1"`
	BId    int64    `json:"qb_id" validate:"required,min=1,gte=1"`
	Dif    int8     `json:"difficulty" validate:"required,min=1,max=3"`
	Answer []Answer `json:"answer" validate:"required"`
}

// ArgAddQus 添加题目
type ArgAddQus struct {
	ArgBaseQus
	QsID string `json:"qid" validate:""`
}

// ArgUpdateQus 更新题目
type ArgUpdateQus struct {
	ArgBaseQus
	QsID int64 `json:"qid" validate:"required"`
}

// ArgQusList 题目列表
type ArgQusList struct {
	ArgPage
	QsBId int64 `form:"qb_id" validate:""`
}

// AddQus 添加题目
type AddQus struct {
	QsID    int64  `json:"qid" validate:""`
	Type    int8   `json:"question_type" validate:"required"`
	AnType  int8   `json:"answer_type" validate:"required,min=1,max=4"`
	Name    string `json:"question_name" validate:"required,min=1,gte=1"`
	BId     int64  `json:"qb_id" validate:"required,min=1,gte=1"`
	Dif     int8   `json:"difficulty" validate:"required,min=1,gte=1"`
	Answers []Answer
}

// ArgQuestionBankBind 关联题库/修改关联
type ArgQuestionBankBind struct {
	QsBId          int64  `json:"qb_id" form:"qb_id" validate:"required"`
	Source         int8   `json:"source" validate:"required"`
	TargetItemType int8   `json:"target_item_type" validate:"required"`
	UseInTime      int64  `json:"use_in_time" validate:"required"`
	TargetItems    string `json:"target_items" validate:"required"`
}

// ArgQuestionBankBinds to do
type ArgQuestionBankBinds struct {
	BandBinds []ArgQuestionBankBind `json:"bind_info" validate:"required"`
}

// ArgQuestionBankBindToDb to do
type ArgQuestionBankBindToDb struct {
	QsBId          int64
	Source         int8
	TargetItemType int8
	UseInTime      int64
	TargetItems    []string
}

// ArgQuestionBankUnbind 关联题库/修改关联
type ArgQuestionBankUnbind struct {
	TargetItems    []int64 `json:"target_items" validate:"required"`
	TargetItemType int8    `json:"target_item_type" validate:"required"`
	Source         int8    `json:"source" validate:"required"`
}

// ArgGetBankBind 查询关联题库信息
type ArgGetBankBind struct {
	TargetItems    []string `json:"target_items" validate:"required"`
	TargetItemType int8     `json:"target_item_type" validate:"required"`
	Source         int8     `json:"source" validate:"required"`
}

// ArgGetQuestion 随机获取一道题
type ArgGetQuestion struct {
	UID            string `form:"uid" json:"uid" validate:"required"`
	TargetItem     string `form:"target_item" json:"target_item" validate:"required"`
	TargetItemType int8   `form:"target_item_type" json:"target_item_type" validate:"required"`
	Source         int8   `form:"source" json:"source" validate:"required"`
	Platform       int8   `form:"platform" json:"platform" validate:"required"`
	ComponentID    int    `form:"component_id" json:"component_id" validate:"required"`
}

// ArgGetBindItems 绑定
type ArgGetBindItems struct {
	ArgPage
	QsBId int64 `form:"qb_id" validate:"required"`
}

// ArgBankSearch 搜索
type ArgBankSearch struct {
	Name string `form:"name" validate:"required"`
}

// ArgCheckAnswer 答案检查
type ArgCheckAnswer struct {
	ArgGetQuestion
	QsID    int64   `json:"qid" validate:"required"`
	Answers []int64 `json:"answers"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
}

// ArgCheckQus 题库检查
type ArgCheckQus struct {
	QusIDs []int64 `json:"ids"`
	Cnt    int     `json:"cnt"`
}
