package model

// UpsertParams .
type UpsertParams struct {
	Business   string `form:"business" validate:"required"`
	DataStr    string `form:"data" validate:"required"`
	Insert     bool   `form:"insert" default:"false"`
	UpsertBody []UpsertBody
}

// UpsertBody job的bulk优化参考这个模板 .
type UpsertBody struct {
	IndexName string
	IndexType string
	IndexID   string
	Doc       MapData
}
