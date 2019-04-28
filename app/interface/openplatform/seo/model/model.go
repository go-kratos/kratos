package model

// ArgProID .
type ArgProID struct {
	ID int `form:"id" validate:"required,min=1,gte=1"`
}

// ArgItemID .
type ArgItemID struct {
	ID int `form:"itemsId" validate:"required,min=1,gte=1"`
}
