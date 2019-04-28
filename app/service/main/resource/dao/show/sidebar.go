package show

import (
	"context"
	"time"

	"go-common/app/service/main/resource/model"
	"go-common/library/log"
)

const (
	_selSideSQL = `SELECT s.id,s.plat,s.module,s.name,s.logo,s.logo_white,s.param,s.rank,l.build,l.conditions,s.tip,s.need_login,s.white_url,s.logo_selected,s.tab_id,s.red_dot_url,lang.name FROM 
 sidebar AS s,sidebar_limit AS l,language AS lang WHERE s.state=1 AND s.id=l.s_id AND lang.id=s.lang_id AND s.online_time<? ORDER BY s.rank DESC,l.id ASC`
)

// SideBar get side bar.
func (d *Dao) SideBar(ctx context.Context, now time.Time) (ss []*model.SideBar, limits map[int64][]*model.SideBarLimit, err error) {
	rows, err := d.db.Query(ctx, _selSideSQL, now)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	limits = make(map[int64][]*model.SideBarLimit)
	for rows.Next() {
		s := &model.SideBar{}
		if err = rows.Scan(&s.ID, &s.Plat, &s.Module, &s.Name, &s.Logo, &s.LogoWhite, &s.Param, &s.Rank, &s.Build, &s.Conditions, &s.Tip, &s.NeedLogin, &s.WhiteURL, &s.LogoSelected, &s.TabID, &s.Red, &s.Language); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		if _, ok := limits[s.ID]; !ok {
			ss = append(ss, s)
		}
		limit := &model.SideBarLimit{
			ID:        s.ID,
			Build:     s.Build,
			Condition: s.Conditions,
		}
		limits[s.ID] = append(limits[s.ID], limit)
	}
	err = rows.Err()
	return
}
