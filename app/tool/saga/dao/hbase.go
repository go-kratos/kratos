package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/tool/saga/conf"
	"go-common/library/log"

	"github.com/pkg/errors"
	"github.com/tsuna/gohbase/hrpc"
)

const (
	_sagaTable         = "ep:saga"
	_ColFamily         = "saga_auth"
	_cSagaPathOwner    = "path_owner"
	_cSagaPathReviewer = "path_reviewer"
)

// sagaAuthKey ...
func sagaAuthKey(projID int, branch string, path string) string {
	return fmt.Sprintf("saga_auth_%d_%s_%s", projID, branch, path)
}

// SetPathAuthH ...
func (d *Dao) SetPathAuthH(c context.Context, projID int, branch string, path string, owners []string, reviewers []string) (err error) {
	var (
		key       = sagaAuthKey(projID, branch, path)
		auth      = make(map[string][]byte)
		bOwner    []byte
		bReviewer []byte
	)

	if bOwner, err = json.Marshal(owners); err != nil {
		return errors.WithStack(err)
	}
	if bReviewer, err = json.Marshal(reviewers); err != nil {
		return errors.WithStack(err)
	}

	auth[_cSagaPathOwner] = bOwner
	auth[_cSagaPathReviewer] = bReviewer
	values := map[string]map[string][]byte{_ColFamily: auth}

	ctx, cancel := context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	defer cancel()

	if _, err = d.hbase.PutStr(ctx, _sagaTable, key, values); err != nil {
		return errors.Wrapf(err, "hbase PutStr error (key: %s values: %v)", key, values)
	}
	return
}

// PathAuthH ...
func (d *Dao) PathAuthH(ctx context.Context, projID int, branch string, path string) (owners []string, reviewers []string, err error) {
	var (
		key    = sagaAuthKey(projID, branch, path)
		result *hrpc.Result
	)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(conf.Conf.HBase.ReadTimeout))
	defer cancel()

	if result, err = d.hbase.GetStr(ctx, _sagaTable, key); err != nil {
		err = errors.Wrapf(err, "hbase GetStr error (key: %s)", key)
		return
	}

	for _, c := range result.Cells {
		switch string(c.Qualifier) {
		case _cSagaPathOwner:
			if err = json.Unmarshal(c.Value, &owners); err != nil {
				err = errors.WithStack(err)
				return
			}
			log.Info("Get key: (%s), owners Info: (%+v)", key, owners)
		case _cSagaPathReviewer:
			if err = json.Unmarshal(c.Value, &reviewers); err != nil {
				err = errors.WithStack(err)
				return
			}
			log.Info("Get key: (%s), reviewers Info: (%+v)", key, reviewers)
		}
	}
	return
}

// DeletePathAuthH ...
func (d *Dao) DeletePathAuthH(c context.Context, projID int, branch string, path string) (err error) {
	key := sagaAuthKey(projID, branch, path)
	ctx, cancel := context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	defer cancel()

	auth := make(map[string][]byte)
	auth[_cSagaPathOwner] = nil
	auth[_cSagaPathReviewer] = nil
	values := map[string]map[string][]byte{_ColFamily: auth}

	if _, err = d.hbase.Delete(ctx, _sagaTable, key, values); err != nil {
		err = errors.Wrapf(err, "hbase delete error (key: %s)", key)
	}
	return
}
