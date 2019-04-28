package validator

//AddGroupParams add group params validate
type AddGroupParams struct {
	Name string `form:"name" validate:"required"`
	Desc string `form:"desc" validate:"required"`
}

//UpdateGroupParams update group params validate
type UpdateGroupParams struct {
	ID   int    `form:"id" validate:"required"`
	Name string `form:"name" validate:"required"`
	Desc string `form:"desc" validate:"required"`
}

//DeleteGroupParams delete group params validate
type DeleteGroupParams struct {
	ID int `form:"id" validate:"required"`
}
