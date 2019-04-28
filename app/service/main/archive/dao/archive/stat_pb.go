package archive

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/library/log"
)

// SetStat3 set all stat
func (d *Dao) SetStat3(c context.Context, st *api.Stat) (err error) {
	d.addStatCache3(c, st)
	var clk *api.Click
	if clk, err = d.click3(c, st.Aid); err != nil {
		log.Error("d.stat(%d) error(%v)", st.Aid, err)
		return
	}
	if clk == nil {
		clk = &api.Click{Aid: st.Aid}
	}
	d.addCache(func() {
		d.addClickCache3(context.TODO(), clk)
	})
	return
}

// Stat3 get archive stat.
func (d *Dao) Stat3(c context.Context, aid int64) (st *api.Stat, err error) {
	var cached = true
	if st, err = d.statCache3(c, aid); err != nil {
		log.Error("d.statCache(%d) error(%v)", aid, err)
		cached = false
	}
	if st != nil {
		return
	}
	if st, err = d.stat3(c, aid); err != nil {
		log.Error("d.stat(%d) error(%v)", aid, err)
		return
	}
	if st == nil {
		st = &api.Stat{Aid: aid}
		return
	}
	if cached {
		d.addCache(func() {
			d.addStatCache3(context.TODO(), st)
		})
	}
	return
}

// Stats3 get archives stat.
func (d *Dao) Stats3(c context.Context, aids []int64) (stm map[int64]*api.Stat, err error) {
	if len(aids) == 0 {
		return
	}
	var (
		missed []int64
		missm  map[int64]*api.Stat
		cached = true
	)
	if stm, missed, err = d.statCaches3(c, aids); err != nil {
		log.Error("d.statCaches(%d) error(%v)", aids, err)
		missed = aids
		stm = make(map[int64]*api.Stat, len(aids))
		err = nil // ignore error
		cached = false
	}
	if stm != nil && len(missed) == 0 {
		return
	}
	if missm, err = d.stats3(c, missed); err != nil {
		log.Error("d.stats(%v) error(%v)", missed, err)
		err = nil // ignore error
	}
	for aid, st := range missm {
		stm[aid] = st
		if cached {
			var cst = &api.Stat{}
			*cst = *st
			d.addCache(func() {
				d.addStatCache3(context.TODO(), cst)
			})
		}
	}
	return
}

// Click3 get archive click.
func (d *Dao) Click3(c context.Context, aid int64) (clk *api.Click, err error) {
	var cached = true
	if clk, err = d.clickCache3(c, aid); err != nil {
		log.Error("d.clickCache(%d) error(%v)", aid, err)
		cached = false
	}
	if clk != nil {
		return
	}
	if clk, err = d.click3(c, aid); err != nil {
		log.Error("d.stat(%d) error(%v)", aid, err)
		return
	}
	if clk == nil {
		clk = &api.Click{Aid: aid}
		return
	}
	if cached {
		d.addCache(func() {
			d.addClickCache3(context.TODO(), clk)
		})
	}
	return
}

// InitStatCache3 if db is nil, set nil cache
func (d *Dao) InitStatCache3(c context.Context, aid int64) (err error) {
	var st *api.Stat
	if st, err = d.stat3(c, aid); err != nil {
		log.Error("d.stat(%d) error(%v)", aid, err)
		return
	}
	if st == nil {
		d.addCache(func() {
			d.addStatCache3(context.TODO(), &api.Stat{Aid: aid})
		})
	}
	var clk *api.Click
	if clk, err = d.click3(c, aid); err != nil {
		log.Error("d.stat(%d) error(%v)", aid, err)
		return
	}
	if clk == nil {
		d.addCache(func() {
			d.addClickCache3(context.TODO(), &api.Click{Aid: aid})
		})
	}
	return
}
