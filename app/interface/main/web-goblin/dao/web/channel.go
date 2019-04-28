package web

import (
	"context"
	"time"

	"go-common/app/interface/main/web-goblin/model/web"
)

const _cardSQL = "SELECT id,title,tag_id,card_type,card_value,recommand_reason,recommand_state,priority FROM channel_card WHERE stime<? AND etime>? AND `check`=2 AND is_delete=0 ORDER BY priority DESC"

// ChCard channel card.
func (d *Dao) ChCard(ctx context.Context, now time.Time) (res map[int64][]*web.ChCard, err error) {
	res = map[int64][]*web.ChCard{}
	rows, err := d.showDB.Query(ctx, _cardSQL, now, now)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &web.ChCard{}
		if err = rows.Scan(&c.ID, &c.Title, &c.ChannelID, &c.Type, &c.Value, &c.Reason, &c.ReasonType, &c.Pos); err != nil {
			return
		}
		res[c.ChannelID] = append(res[c.ChannelID], c)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
