package manager

import (
	"context"

	pb "go-common/app/service/main/resource/api/v1"
)

const (
	_specialSQL = "SELECT `id`,`title`,`desc`,`cover`,`scover`,`re_type`,`re_value`,`corner`,`size`,`card` FROM special_card"
)

//Specials get specials cars from DB
func (d *Dao) Specials(c context.Context) (sps map[int64]*pb.SpecialReply, err error) {
	rows, err := d.db.Query(c, _specialSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	sps = make(map[int64]*pb.SpecialReply)
	for rows.Next() {
		sc := &pb.SpecialReply{}
		if err = rows.Scan(&sc.Id, &sc.Title, &sc.Desc, &sc.Cover, &sc.Scover, &sc.ReType, &sc.ReValue, &sc.Corner, &sc.Siz, &sc.Card); err != nil {
			return
		}
		sps[sc.Id] = sc
	}
	err = rows.Err()
	return
}
