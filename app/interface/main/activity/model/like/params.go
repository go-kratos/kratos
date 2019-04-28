package like

// ParamMsg notify param msg.
type ParamMsg struct {
	Msg string `form:"msg" validate:"required"`
}

// ParamTeams add follow param teams
type ParamTeams struct {
	Teams []string `form:"teams,split" validate:"gt=0,dive,gt=0"`
}

// ParamSid  sid param
type ParamSid struct {
	Sid int64 `form:"sid" validate:"required,min=1"`
}

// ParamAddGuess add guess param
type ParamAddGuess struct {
	ObjID  int64 `form:"obj_id" validate:"required,min=1"`
	Result int64 `form:"result" validate:"required,min=1"`
	Stake  int64 `form:"stake"  validate:"gt=0"`
}

// ParamObject unstart  object param
type ParamObject struct {
	Sid int64 `form:"sid" validate:"required,min=1"`
	Pn  int   `form:"pn" validate:"gt=0"`
	Ps  int   `form:"ps" validate:"gt=0,lte=50"`
}

// ParamAddLikeAct add likeAct param
type ParamAddLikeAct struct {
	Sid   int64 `form:"sid" validate:"required,min=1"`
	Lid   int64 `form:"lid" validate:"required,min=1"`
	Score int64 `form:"score" validate:"min=1,max=5"`
}

// ParamMissionLikeAct add missionAct param
type ParamMissionLikeAct struct {
	Sid int64 `form:"sid" validate:"min=1"`
	Lid int64 `form:"lid" validate:"min=1"`
}

// ParamMissionFriends get mission friends list
type ParamMissionFriends struct {
	Sid  int64 `form:"sid"  validate:"min=1"`
	Lid  int64 `form:"lid"  validate:"min=1"`
	Size int   `form:"size" validate:"min=1,max=50"`
}

// ParamStoryKingAct .
type ParamStoryKingAct struct {
	Sid   int64 `form:"sid" validate:"required,min=1"`
	Lid   int64 `form:"lid" validate:"required,min=1"`
	Score int64 `form:"score" validate:"min=1,max=10"`
}

// ParamList .
type ParamList struct {
	Sid  int64  `form:"sid" validate:"min=1"`
	Type string `form:"type" default:"like"`
	Pn   int    `form:"pn" default:"1" validate:"min=1"`
	Ps   int    `form:"ps" default:"30" validate:"min=1,max=100"`
}
