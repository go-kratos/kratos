package model

//AddGroupArg arg
type AddGroupArg struct {
	Name      string `form:"name" validate:"required"`
	Tag       string `form:"tag" validate:"required"`
	ShortTag  string `form:"short_tag"`
	FontColor string `form:"font_color"`
	BgColor   string `form:"bg_color"`
	Remark    string `from:"remark"`
}

//EditGroupArg arg
type EditGroupArg struct {
	AddArg *AddGroupArg
	ID     int64 `form:"id" validate:"required"`
}

//RemoveGroupArg arg
type RemoveGroupArg struct {
	ID int32 `form:"id" validate:"required"`
}

//GetGroupArg arg
type GetGroupArg struct {
	State int `form:"state"`
}
