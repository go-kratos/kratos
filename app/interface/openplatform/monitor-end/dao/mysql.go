package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/interface/openplatform/monitor-end/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	_allGroupsSQL     = "SELECT id, name, receivers, `interval`, ctime, mtime FROM `alert_group` WHERE is_deleted = 0 ORDER BY id DESC"
	_groupSQL         = "SELECT id, name, receivers, `interval`, ctime, mtime FROM `alert_group` WHERE id = ? AND is_deleted = 0"
	_groupsSQL        = "SELECT id, name, receivers, `interval`, ctime, mtime FROM `alert_group` WHERE id in (%s) AND is_deleted = 0"
	_addGroupSQL      = "INSERT INTO `alert_group` (name,receivers,`interval`)VALUES(?,?,?)"
	_groupByNameSQL   = "SELECT id, name, receivers, `interval`, ctime, mtime FROM `alert_group` WHERE name = ? AND is_deleted = 0"
	_updateGroupSQL   = "UPDATE `alert_group` SET name = ?, receivers = ?, `interval` = ? WHERE id = ?"
	_deleteGroupSQL   = "UPDATE `alert_group` SET is_deleted = 1 WHERE id = ?"
	_targetSQL        = "SELECT id, sub_event, event, product, source, group_id, threshold, duration, state, ctime, mtime FROM alert_target WHERE id = ? AND deleted_time = 0"
	_targetQuerySQL   = "SELECT id, sub_event, event, product, source, group_id, threshold, duration, state, ctime, mtime FROM alert_target"
	_allTargetsSQL    = "SELECT id, sub_event, event, product, source, group_id, threshold, duration, state, ctime, mtime FROM alert_target WHERE deleted_time = 0 AND state = ?"
	_countTargetSQL   = "SELECT count(id) as cnt FROM alert_target"
	_addTargetSQL     = "INSERT INTO alert_target (sub_event, event, product, source, group_id, threshold, duration, state)VALUES(?,?,?,?,?,?,?,?)"
	_updateTargetSQL  = "UPDATE alert_target SET sub_event = ?, event =? , product = ?, source = ?, group_id = ?, threshold = ?, duration = ?, state = ? WHERE id = ?"
	_existTargetSQL   = "SELECT id FROM alert_target WHERE sub_event = ? AND event =? AND product = ? AND source = ? AND deleted_time = 0"
	_targetSyncSQL    = "UPDATE alert_target set state = ? WHERE id = ?"
	_deleteTargetSQL  = "UPDATE alert_target set deleted_time = ? WHERE id = ?"
	_productSQL       = "SELECT id, name, group_id, ctime, mtime FROM alert_product WHERE id = ? AND is_deleted = 0"
	_productByNameSQL = "SELECT id, name, group_id, ctime, mtime FROM alert_product WHERE name = ? AND is_deleted = 0"
	_allProductsSQL   = "SELECT id, name, group_id, ctime, mtime FROM alert_product WHERE is_deleted = 0 AND state = 1 ORDER BY ctime desc"
	_addProductSQL    = "INSERT INTO alert_product (name, group_id, state)VALUES(?,?,?)"
	_updateProductSQL = "UPDATE alert_product SET name = ?, group_id = ?, state = ? WHERE id = ?"
	_deleteProductSQL = "UPDATE alert_product SET is_deleted = 1 WHERE id = ?"
)

// Group query group by id.
func (d *Dao) Group(c context.Context, id int64) (res *model.Group, err error) {
	res = &model.Group{}
	if err = d.db.QueryRow(c, _groupSQL, id).Scan(&res.ID, &res.Name, &res.Receivers, &res.Interval, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("d.Group.Scan error(%+v), id(%d)", err, id)
	}
	return
}

// Groups query groups by ids.
func (d *Dao) Groups(c context.Context, ids []int64) (res []*model.Group, err error) {
	if len(ids) == 0 {
		return
	}
	var (
		rows  *xsql.Rows
		query = fmt.Sprintf(_groupsSQL, xstr.JoinInts(ids))
	)
	if rows, err = d.db.Query(c, query); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("d.Group.Scan error(%+v), id(%d)", err, ids)
	}
	for rows.Next() {
		r := &model.Group{}
		if err = rows.Scan(&r.ID, &r.Name, &r.Receivers, &r.Interval, &r.Ctime, &r.Mtime); err != nil {
			log.Error("d.Group.Scan error(%+v), id(%d)", err, ids)
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// GroupByName query group id by name.
func (d *Dao) GroupByName(c context.Context, name string) (res *model.Group, err error) {
	res = &model.Group{}
	if err = d.db.QueryRow(c, _groupByNameSQL, name).Scan(&res.ID, &res.Name, &res.Receivers, &res.Interval, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.GroupByName.Scan error(%+v), name(%s)", err, name)
	}
	return
}

// AddGroup add new group.
func (d *Dao) AddGroup(c context.Context, g *model.Group) (res int64, err error) {
	var r sql.Result
	if r, err = d.db.Exec(c, _addGroupSQL, g.Name, g.Receivers, g.Interval); err != nil {
		log.Error("d.AddGroup.Exec error(%+v), group(%+v)", err, g)
		return
	}
	if res, err = r.LastInsertId(); err != nil {
		log.Error("d.AddGroup.LastInsertId error(%+v), group(%+v)", err, g)
	}
	return
}

// UpdateGroup update group.
func (d *Dao) UpdateGroup(c context.Context, g *model.Group) (res int64, err error) {
	var r sql.Result
	if r, err = d.db.Exec(c, _updateGroupSQL, g.Name, g.Receivers, g.Interval, g.ID); err != nil {
		log.Error("d.UpdateGroup.Exec error(%+v), group(%+v)", err, g)
		return
	}
	if res, err = r.RowsAffected(); err != nil {
		log.Error("d.UpdateGroup.RowsAffected error(%+v), group(%+v)", err, g)
	}
	return
}

// AllGroups return all groups.
func (d *Dao) AllGroups(c context.Context) (res []*model.Group, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _allGroupsSQL); err != nil {
		log.Error("d.AllGroups.Query error(%+v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var g = &model.Group{}
		if err = rows.Scan(&g.ID, &g.Name, &g.Receivers, &g.Interval, &g.Ctime, &g.Mtime); err != nil {
			log.Error("d.AllGroups.Scan error(%+v)]", err)
			return
		}
		res = append(res, g)
	}
	err = rows.Err()
	return
}

// DeleteGroup delete group.
func (d *Dao) DeleteGroup(c context.Context, id int64) (res int64, err error) {
	var r sql.Result
	if r, err = d.db.Exec(c, _deleteGroupSQL, id); err != nil {
		log.Error("d.DeleteGroup.Exec error(%+v), id(%d)", err, id)
		return
	}
	if res, err = r.RowsAffected(); err != nil {
		log.Error("d.DeleteGroup.RowsAffected error(%+v), id(%d)", err, id)
	}
	return
}

// IsExisted return id if target existed by sub_event, event, product, source.
func (d *Dao) IsExisted(c context.Context, t *model.Target) (res int64, err error) {
	if err = d.db.QueryRow(c, _existTargetSQL, t.SubEvent, t.Event, t.Product, t.Source).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.Isexisted.Scan error(%+v), target(%+v)", err, t)
	}
	return
}

// Target get target by id.
func (d *Dao) Target(c context.Context, id int64) (res *model.Target, err error) {
	res = &model.Target{}
	if err = d.db.QueryRow(c, _targetSQL, id).Scan(&res.ID, &res.SubEvent, &res.Event, &res.Product, &res.Source, &res.GroupIDs, &res.Threshold, &res.Duration, &res.State, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("d.Target.Scan error(%+v), id(%d)", err, id)
	}
	if res.GroupIDs != "" {
		var gids []int64
		if gids, err = xstr.SplitInts(res.GroupIDs); err != nil {
			log.Error("d.Product.SplitInts error(%+v), group ids(%s)", err, res.GroupIDs)
			return
		}
		if res.Groups, err = d.Groups(c, gids); err != nil {
			return
		}
	}
	return
}

// AllTargets return all targets by state.
func (d *Dao) AllTargets(c context.Context, state int) (res []*model.Target, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _allTargetsSQL, state); err != nil {
		log.Error("d.AllTargets.Query error(%+v), sql(%s)", err, _allTargetsSQL)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var t = &model.Target{}
		if err = rows.Scan(&t.ID, &t.SubEvent, &t.Event, &t.Product, &t.Source, &t.GroupIDs, &t.Threshold, &t.Duration, &t.State, &t.Ctime, &t.Mtime); err != nil {
			log.Error("d.AllTargets.Scan error(%+v), sql(%s)", err, _allTargetsSQL)
			return
		}
		if t.GroupIDs != "" {
			var gids []int64
			if gids, err = xstr.SplitInts(t.GroupIDs); err != nil {
				log.Error("d.Product.SplitInts error(%+v), group ids(%s)", err, t.GroupIDs)
				return
			}
			if t.Groups, err = d.Groups(c, gids); err != nil {
				return
			}
		}
		res = append(res, t)
	}
	err = rows.Err()
	return
}

// AddTarget add a new target.
func (d *Dao) AddTarget(c context.Context, t *model.Target) (res int64, err error) {
	var r sql.Result
	if r, err = d.db.Exec(c, _addTargetSQL, t.SubEvent, t.Event, t.Product, t.Source, t.GroupIDs, t.Threshold, t.Duration, t.State); err != nil {
		log.Error("d.AddTarget.Exec error(%+v), target(%+v)", err, t)
		return
	}
	if res, err = r.LastInsertId(); err != nil {
		log.Error("d.AddTarget.LastInsertId error(%+v), target(%+v)", err, t)
	}
	return
}

// UpdateTarget uodate target.
func (d *Dao) UpdateTarget(c context.Context, t *model.Target) (res int64, err error) {
	var r sql.Result
	if r, err = d.db.Exec(c, _updateTargetSQL, t.SubEvent, t.Event, t.Product, t.Source, t.GroupIDs, t.Threshold, t.Duration, t.State, t.ID); err != nil {
		log.Error("d.UpdateGroup.Exec error(%+v), target(%+v)", err, t)
		return
	}
	if res, err = r.RowsAffected(); err != nil {
		log.Error("d.UpdateGroup.RowsAffected error(%+v), target(%+v)", err, t)
	}
	return
}

// DeleteTarget delete target by id.
func (d *Dao) DeleteTarget(c context.Context, id int64) (res int64, err error) {
	var (
		r   sql.Result
		now = time.Now()
	)
	if r, err = d.db.Exec(c, _deleteTargetSQL, now, id); err != nil {
		log.Error("d.UpdateGroup.Exec error(%+v), target(%d)", err, id)
		return
	}
	if res, err = r.RowsAffected(); err != nil {
		log.Error("d.UpdateGroup.RowsAffected error(%+v), target(%d)", err, id)
	}
	return
}

// TargetsByQuery query targets by query.
func (d *Dao) TargetsByQuery(c context.Context, where string) (res []*model.Target, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _targetQuerySQL+where); err != nil {
		log.Error("d.TargetsByQuery.Query error(%+v), sql(%s)", err, _targetQuerySQL+where)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var t = &model.Target{}
		if err = rows.Scan(&t.ID, &t.SubEvent, &t.Event, &t.Product, &t.Source, &t.GroupIDs, &t.Threshold, &t.Duration, &t.State, &t.Ctime, &t.Mtime); err != nil {
			log.Error("d.TargetsByQuery.Scan error(%+v), sql(%s)", err, _targetQuerySQL+where)
			return
		}
		if t.GroupIDs != "" {
			var gids []int64
			if gids, err = xstr.SplitInts(t.GroupIDs); err != nil {
				log.Error("d.Product.SplitInts error(%+v), group ids(%s)", err, t.GroupIDs)
				return
			}
			if t.Groups, err = d.Groups(c, gids); err != nil {
				return
			}
		}
		res = append(res, t)
	}
	return
}

// CountTargets .
func (d *Dao) CountTargets(c context.Context, where string) (res int, err error) {
	if err = d.db.QueryRow(c, _countTargetSQL+where).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.CountTargets.Scan error(%+v), sql(%s)", err, _countTargetSQL+where)
	}
	return
}

// TargetSync sync target state by id.
func (d *Dao) TargetSync(c context.Context, id int64, state int) (err error) {
	if _, err = d.db.Exec(c, _targetSyncSQL, state, id); err != nil {
		log.Error("d.TargetSync.Exec error(%+v), id(%d), state(%d)", err, id, state)
	}
	return
}

// Product get product by id.
func (d *Dao) Product(c context.Context, id int64) (res *model.Product, err error) {
	res = &model.Product{}
	if err = d.db.QueryRow(c, _productSQL, id).Scan(&res.ID, &res.Name, &res.GroupIDs, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("d.Product.Scan error(%+v), id(%d)", err, id)
	}
	if res.GroupIDs != "" {
		var gids []int64
		if gids, err = xstr.SplitInts(res.GroupIDs); err != nil {
			log.Error("d.Product.SplitInts error(%+v), group ids(%s)", err, res.GroupIDs)
			return
		}
		if res.Groups, err = d.Groups(c, gids); err != nil {
			return
		}
	}
	return
}

// ProductByName get product bu name.
func (d *Dao) ProductByName(c context.Context, name string) (res *model.Product, err error) {
	res = &model.Product{}
	if err = d.db.QueryRow(c, _productByNameSQL, name).Scan(&res.ID, &res.Name, &res.GroupIDs, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		log.Error("d.ProductByName.Scan error(%+v), name(%s)", err, name)
	}
	if res.GroupIDs != "" {
		var gids []int64
		if gids, err = xstr.SplitInts(res.GroupIDs); err != nil {
			log.Error("d.Product.SplitInts error(%+v), group ids(%s)", err, res.GroupIDs)
			return
		}
		if res.Groups, err = d.Groups(c, gids); err != nil {
			return
		}
	}
	return
}

// AllProducts return all products.
func (d *Dao) AllProducts(c context.Context) (res []*model.Product, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _allProductsSQL); err != nil {
		log.Error("d.AllProducts.Query error(%+v), sql(%s)", err, _allProductsSQL)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var p = &model.Product{}
		if err = rows.Scan(&p.ID, &p.Name, &p.GroupIDs, &p.Ctime, &p.Mtime); err != nil {
			log.Error("d.AllProducts.Scan error(%+v), sql(%s)", err, _allProductsSQL)
			return
		}
		if p.GroupIDs != "" {
			var gids []int64
			if gids, err = xstr.SplitInts(p.GroupIDs); err != nil {
				log.Error("d.Product.SplitInts error(%+v), group ids(%s)", err, p.GroupIDs)
				return
			}
			if p.Groups, err = d.Groups(c, gids); err != nil {
				return
			}
		}
		res = append(res, p)
	}
	err = rows.Err()
	return
}

// AddProduct add a new product.
func (d *Dao) AddProduct(c context.Context, p *model.Product) (res int64, err error) {
	var r sql.Result
	if r, err = d.db.Exec(c, _addProductSQL, p.Name, p.GroupIDs, p.State); err != nil {
		log.Error("d.AddProduct.Exec error(%+v), product(%+v)", err, p)
		return
	}
	if res, err = r.LastInsertId(); err != nil {
		log.Error("d.AddProduct.RowsAffected error(%+v), product(%+v)", err, p)
	}
	return
}

// UpdateProduct update product by id.
func (d *Dao) UpdateProduct(c context.Context, p *model.Product) (res int64, err error) {
	var r sql.Result
	if r, err = d.db.Exec(c, _updateProductSQL, p.Name, p.GroupIDs, p.State, p.ID); err != nil {
		log.Error("d.DeleteProduct.Exec error(%+v), product(%+v)", err, p)
		return
	}
	if res, err = r.RowsAffected(); err != nil {
		log.Error("d.DeleteProduct.RowsAffected error(%+v), product(%+v)", err, p)
	}
	return
}

// DeleteProduct delete a product by id.
func (d *Dao) DeleteProduct(c context.Context, id int64) (res int64, err error) {
	var r sql.Result
	if r, err = d.db.Exec(c, _deleteProductSQL, id); err != nil {
		log.Error("d.DeleteProduct.Exec error(%+v), id(%d)", err, id)
		return
	}
	if res, err = r.RowsAffected(); err != nil {
		log.Error("d.DeleteProduct.RowsAffected error(%+v), id(%d)", err, id)
	}
	return
}
