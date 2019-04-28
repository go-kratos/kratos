package validator

//VerionParams presents version params
type VerionParams struct {
	Group int `form:"group" validate:"min=0,required"`
}

//ListParams presents listAb params
type ListParams struct {
	Pn      int    `form:"pn" validate:"min=1,required"`
	Ps      int    `form:"ps" validate:"min=1,required"`
	Mstatus string `form:"mstatus" validate:"required"`
	Group   int    `form:"group"`
}

//AddAbParams presents addAbtest params
type AddAbParams struct {
	Data  string `form:"data" validate:"required"`
	Group int    `form:"group"`
}

//UpdateAbParams presents updateAbtest params
type UpdateAbParams struct {
	ID    int    `form:"id" validate:"required"`
	Data  string `form:"data" validate:"required"`
	Group int    `form:"group"`
}

//DelAbParams presents deleteAbtest params
type DelAbParams struct {
	ID    int `form:"id" validate:"required"`
	Group int `form:"group"`
}

//UpdateStatusAbParams presents updateStatusAbtest params
type UpdateStatusAbParams struct {
	ID       int    `form:"id" validate:"required"`
	Status   int    `form:"status"`
	Modifier string `form:"modifier" validate:"required"`
	Group    int    `form:"group"`
}
