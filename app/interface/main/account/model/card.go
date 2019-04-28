package model

// ArgCardID arg card id.
type ArgCardID struct {
	ID int64 `form:"id" validate:"required,min=1,gte=1"`
}

// ArgGroupID arg group id.
type ArgGroupID struct {
	ID int64 `form:"id" validate:"required,min=1,gte=1"`
}
