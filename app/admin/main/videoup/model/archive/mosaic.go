package archive

import (
	"go-common/app/admin/main/videoup/model/utils"
)

//Mosaic 马赛克
type Mosaic struct {
	ID         int64            `json:"id"`
	AID        int64            `json:"aid"`
	CID        int64            `json:"cid"`
	Coordinate string           `json:"coordinate"`
	CTime      utils.FormatTime `json:"ctime"`
}
