package card

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-channel/conf"
	"go-common/app/interface/main/app-channel/model/card"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
)

const (
	_cardSQL = `SELECT c.id,c.title,c.tag_id,c.card_type,c.card_value,c.recommand_reason,c.recommand_state,c.priority FROM channel_card AS c
	WHERE c.stime<? AND c.etime>? AND c.check=2 AND c.is_delete=0 ORDER BY c.priority DESC`
	_cardPlatSQL = `SELECT card_id,plat,conditions,build FROM channel_card_plat WHERE is_delete=0`
	_followSQL   = "SELECT `id`,`type`,`long_title`,`content` FROM `card_follow` WHERE `deleted`=0"
	_cardSetSQL  = `SELECT c.id,c.type,c.value,c.title,c.long_title,c.content FROM card_set AS c WHERE c.deleted=0`
)

// Dao is card dao.
type Dao struct {
	db *sql.DB
	// memcache
	expire int32
	mc     *memcache.Pool
}

// New is card dao new.
func New(c *conf.Config) *Dao {
	d := &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
		// memcache
		expire: int32(time.Duration(c.Memcache.Channels.Expire) / time.Second),
		mc:     memcache.NewPool(c.Memcache.Channels.Config),
	}
	return d
}

// Card channel card
func (d *Dao) Card(ctx context.Context, now time.Time) (res map[int64][]*card.Card, err error) {
	res = map[int64][]*card.Card{}
	rows, err := d.db.Query(ctx, _cardSQL, now, now)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &card.Card{}
		if err = rows.Scan(&c.ID, &c.Title, &c.ChannelID, &c.Type, &c.Value, &c.Reason, &c.ReasonType, &c.Pos); err != nil {
			return
		}
		res[c.ChannelID] = append(res[c.ChannelID], c)
	}
	return
}

// CardPlat channel card  plat
func (d *Dao) CardPlat(ctx context.Context) (res map[string][]*card.CardPlat, err error) {
	res = map[string][]*card.CardPlat{}
	var (
		_initCardPlatKey = "card_platkey_%d_%d"
	)
	rows, err := d.db.Query(ctx, _cardPlatSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &card.CardPlat{}
		if err = rows.Scan(&c.CardID, &c.Plat, &c.Condition, &c.Build); err != nil {
			return
		}
		key := fmt.Sprintf(_initCardPlatKey, c.Plat, c.CardID)
		res[key] = append(res[key], c)
	}
	return
}

// UpCard upper
func (d *Dao) UpCard(ctx context.Context) (res map[int64]*operate.Follow, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _followSQL); err != nil {
		return
	}
	defer rows.Close()
	res = make(map[int64]*operate.Follow)
	for rows.Next() {
		c := &operate.Follow{}
		if err = rows.Scan(&c.ID, &c.Type, &c.Title, &c.Content); err != nil {
			return
		}
		c.Change()
		res[c.ID] = c
	}
	return
}

// CardSet card set
func (d *Dao) CardSet(ctx context.Context) (res map[int64]*operate.CardSet, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _cardSetSQL); err != nil {
		return
	}
	defer rows.Close()
	res = make(map[int64]*operate.CardSet)
	for rows.Next() {
		var (
			c     = &operate.CardSet{}
			value string
		)
		if err = rows.Scan(&c.ID, &c.Type, &value, &c.Title, &c.LongTitle, &c.Content); err != nil {
			return
		}
		c.Value, _ = strconv.ParseInt(value, 10, 64)
		res[c.ID] = c
	}
	return
}

// PingDB ping db
func (d *Dao) PingDB(c context.Context) (err error) {
	return d.db.Ping(c)
}
