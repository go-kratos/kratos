package dao

import (
	"context"

	"go-common/app/infra/databus/model"
	"go-common/library/log"
)

const (
	_getAuthSQL = `SELECT auth.group_name,auth.operation,app.app_key,app.app_secret,auth.topic,app.cluster
				FROM auth LEFT JOIN app ON auth.app_id=app.id`

	_getAuth2SQL = `SELECT auth2.group,auth2.operation,app2.app_key,app2.app_secret,auth2.number,topic.topic,topic.cluster
				FROM auth2 LEFT JOIN app2 On auth2.app_id=app2.id LEFT JOIN topic On topic.id=auth2.topic_id WHERE auth2.app_id!=0 AND auth2.is_delete=0`
)

// Auth verify group,topic,key
func (d *Dao) Auth(c context.Context) (auths map[string]*model.Auth, err error) {
	auths = make(map[string]*model.Auth)
	// auth
	rows, err := d.db.Query(c, _getAuthSQL)
	if err != nil {
		log.Error("getAuthStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.Auth{}
		if err = rows.Scan(&a.Group, &a.Operation, &a.Key, &a.Secret, &a.Topic, &a.Cluster); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		auths[a.Group] = a
	}
	// auth2
	rows2, err := d.db.Query(c, _getAuth2SQL)
	if err != nil {
		log.Error("getAuthStmt.Query error(%v)", err)
		return
	}
	defer rows2.Close()
	for rows2.Next() {
		a := &model.Auth{}
		if err = rows2.Scan(&a.Group, &a.Operation, &a.Key, &a.Secret, &a.Batch, &a.Topic, &a.Cluster); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		auths[a.Group] = a
	}
	return
}
