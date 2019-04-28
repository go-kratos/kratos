package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/tool/saga/model"
	"go-common/library/cache/redis"

	"github.com/pkg/errors"
)

func mergeTaskKey(taskType int) string {
	return fmt.Sprintf("saga_task_%d", taskType)
}

func mrIIDKey(mrIID int) string {
	return fmt.Sprintf("saga_mrIID_%d", mrIID)
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// AddMRIID ...
func (d *Dao) AddMRIID(c context.Context, mrIID int, expire int) (err error) {
	var (
		key  = mrIIDKey(mrIID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("SET", key, mrIID); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	if _, err = conn.Do("EXPIRE", key, expire); err != nil {
		return
	}
	return
}

// ExistMRIID ...
func (d *Dao) ExistMRIID(c context.Context, mrIID int) (ok bool, err error) {
	var (
		key  = mrIIDKey(mrIID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = redis.Int(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return false, err
	}
	return true, nil
}

// DeleteMRIID ...
func (d *Dao) DeleteMRIID(c context.Context, mrIID int) (err error) {
	var (
		key  = mrIIDKey(mrIID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// PushMergeTask ...
func (d *Dao) PushMergeTask(c context.Context, taskType int, taskInfo *model.TaskInfo) (err error) {
	var (
		key  = mergeTaskKey(taskType)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(taskInfo); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = conn.Send("LPUSH", key, bs); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// DeleteMergeTask ...
func (d *Dao) DeleteMergeTask(c context.Context, taskType int, taskInfo *model.TaskInfo) (err error) {
	var (
		key  = mergeTaskKey(taskType)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(taskInfo); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = conn.Send("LREM", key, 0, bs); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// MergeTasks ...
func (d *Dao) MergeTasks(c context.Context, taskType int) (count int, taskInfos []*model.TaskInfo, err error) {
	var (
		key    = mergeTaskKey(taskType)
		values [][]byte
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("LLEN", key)); err != nil {
		return
	}
	if values, err = redis.ByteSlices(conn.Do("LRANGE", key, 0, -1)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	taskInfos = make([]*model.TaskInfo, 0, count)
	for _, value := range values {
		taskInfo := &model.TaskInfo{}
		if err = json.Unmarshal(value, &taskInfo); err != nil {
			err = errors.WithStack(err)
			return
		}
		taskInfos = append(taskInfos, taskInfo)
		//taskInfos = append([]*model.TaskInfo{taskInfo}, taskInfos...)
	}
	return
}

func mergeInfoKey(projID int, branch string) string {
	return fmt.Sprintf("saga_mergeInfo_%d_%s", projID, branch)
}

func pathOwnerKey(projID int, branch string, path string) string {
	return fmt.Sprintf("saga_PathOwner_%d_%s_%s", projID, branch, path)
}

func pathReviewerKey(projID int, branch string, path string) string {
	return fmt.Sprintf("saga_PathReviewer_%d_%s_%s", projID, branch, path)
}

func authInfoKey(projID int, mrIID int) string {
	return fmt.Sprintf("saga_auth_%d_%d", projID, mrIID)
}

// SetMergeInfo ...
func (d *Dao) SetMergeInfo(c context.Context, projID int, branch string, mergeInfo *model.MergeInfo) (err error) {
	var (
		key  = mergeInfoKey(projID, branch)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(mergeInfo); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// MergeInfo ...
func (d *Dao) MergeInfo(c context.Context, projID int, branch string) (ok bool, mergeInfo *model.MergeInfo, err error) {
	var (
		key   = mergeInfoKey(projID, branch)
		value []byte
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	mergeInfo = &model.MergeInfo{}
	if err = json.Unmarshal(value, &mergeInfo); err != nil {
		err = errors.WithStack(err)
		return
	}
	ok = true
	return
}

// DeleteMergeInfo ...
func (d *Dao) DeleteMergeInfo(c context.Context, projID int, branch string) (err error) {
	var (
		key  = mergeInfoKey(projID, branch)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// SetPathOwner ...
func (d *Dao) SetPathOwner(c context.Context, projID int, branch string, path string, owners []string) (err error) {
	var (
		key  = pathOwnerKey(projID, branch, path)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(owners); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// PathOwner ...
func (d *Dao) PathOwner(c context.Context, projID int, branch string, path string) (owners []string, err error) {
	var (
		key  = pathOwnerKey(projID, branch, path)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	if err = json.Unmarshal(bs, &owners); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// SetPathReviewer ...
func (d *Dao) SetPathReviewer(c context.Context, projID int, branch string, path string, reviewers []string) (err error) {
	var (
		key  = pathReviewerKey(projID, branch, path)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(reviewers); err != nil {
		return errors.WithStack(err)
	}
	if err = conn.Send("SET", key, bs); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// PathReviewer ...
func (d *Dao) PathReviewer(c context.Context, projID int, branch string, path string) (reviewers []string, err error) {
	var (
		key  = pathReviewerKey(projID, branch, path)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	if err = json.Unmarshal(bs, &reviewers); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// pathAuthKey ...
func pathAuthKey(projID int, branch string, path string) string {
	return fmt.Sprintf("saga_path_auth_%d_%s_%s", projID, branch, path)
}

// PathAuthR ...
func (d *Dao) PathAuthR(c context.Context, projID int, branch string, path string) (authUsers *model.AuthUsers, err error) {
	var (
		key  = pathAuthKey(projID, branch, path)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	authUsers = new(model.AuthUsers)
	if err = json.Unmarshal(bs, authUsers); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// SetPathAuthR ...
func (d *Dao) SetPathAuthR(c context.Context, projID int, branch string, path string, authUsers *model.AuthUsers) (err error) {
	var (
		key  = pathAuthKey(projID, branch, path)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(authUsers); err != nil {
		return errors.WithStack(err)
	}
	if err = conn.Send("SET", key, bs); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// DeletePathAuthR ...
func (d *Dao) DeletePathAuthR(c context.Context, projID int, branch string, path string) (err error) {
	var (
		key  = pathAuthKey(projID, branch, path)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// SetReportStatus ...
func (d *Dao) SetReportStatus(c context.Context, projID int, mrIID int, result bool) (err error) {
	var (
		key  = authInfoKey(projID, mrIID)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(result); err != nil {
		return errors.WithStack(err)
	}
	if err = conn.Send("SET", key, bs); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// ReportStatus ...
func (d *Dao) ReportStatus(c context.Context, projID int, mrIID int) (result bool, err error) {
	var (
		key   = authInfoKey(projID, mrIID)
		value []byte
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	if err = json.Unmarshal(value, &result); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// DeleteReportStatus ...
func (d *Dao) DeleteReportStatus(c context.Context, projID int, mrIID int) (err error) {
	var (
		key  = authInfoKey(projID, mrIID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}
