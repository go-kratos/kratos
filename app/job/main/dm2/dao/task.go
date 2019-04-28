package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go-common/app/job/main/dm2/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_taskInfo      = "SELECT id,topic,state,qcount,result,sub,last_index,priority,creator,reviewer,title FROM dm_task WHERE state=? order by priority desc,ctime asc"
	_oneTask       = "SELECT id,topic,state,qcount,result,sub,last_index,priority,creator,reviewer,title FROM dm_task WHERE state in (3,9) order by priority desc,ctime asc limit 1"
	_taskByID      = "SELECT id,topic,state,qcount,result,sub,last_index,priority,creator,reviewer,title FROM dm_task WHERE id=?"
	_uptTask       = "UPDATE dm_task SET state=?,last_index=?,qcount=?,result=? WHERE id=?"
	_uptSubTask    = "UPDATE dm_sub_task SET tcount=tcount+?,end=? WHERE task_id=?"
	_selectSubTask = "SELECT id,operation,rate,tcount,start,end FROM dm_sub_task WHERE task_id=?"

	_delDMIndex     = "UPDATE dm_index_%03d SET state=? WHERE oid=? AND id IN (%s)"
	_uptSubCountSQL = "UPDATE dm_subject_%02d SET count=count-? WHERE type=? AND oid=?"

	_merakMsgURI = "/"
)

// TaskInfos task infos.
func (d *Dao) TaskInfos(c context.Context, state int32) (tasks []*model.TaskInfo, err error) {
	rows, err := d.biliDMWriter.Query(c, _taskInfo, state)
	if err != nil {
		log.Error("d.biliDMWriter.Query(query:%s,state:%d) error(%v)", _taskInfo, state, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		task := &model.TaskInfo{}
		if err = rows.Scan(&task.ID, &task.Topic, &task.State, &task.Count, &task.Result, &task.Sub, &task.LastIndex, &task.Priority, &task.Creator, &task.Reviewer, &task.Title); err != nil {
			log.Error("d.biliDMWriter.Scan(query:%s,state:%d) error(%v)", _taskInfo, state, err)
			return
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.biliDMWriter.rows.Err() error(%v)", err)
	}
	return
}

// OneTask .
func (d *Dao) OneTask(c context.Context) (task *model.TaskInfo, err error) {
	task = &model.TaskInfo{}
	row := d.biliDMWriter.QueryRow(c, _oneTask)
	if err = row.Scan(&task.ID, &task.Topic, &task.State, &task.Count, &task.Result, &task.Sub, &task.LastIndex, &task.Priority, &task.Creator, &task.Reviewer, &task.Title); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			task = nil
		} else {
			log.Error("d.biliDMWriter.Scan(query:%s) error(%v)", _oneTask, err)
		}
	}
	return
}

// TaskInfoByID .
func (d *Dao) TaskInfoByID(c context.Context, id int64) (task *model.TaskInfo, err error) {
	task = &model.TaskInfo{}
	row := d.biliDMWriter.QueryRow(c, _taskByID, id)
	if err = row.Scan(&task.ID, &task.Topic, &task.State, &task.Count, &task.Result, &task.Sub, &task.LastIndex, &task.Priority, &task.Creator, &task.Reviewer, &task.Title); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			task = nil
		} else {
			log.Error("d.biliDMWriter.Scan(query:%s,id:%d) error(%v)", _taskByID, id, err)
		}
	}
	return
}

// UpdateTask update dm task.
func (d *Dao) UpdateTask(c context.Context, task *model.TaskInfo) (affected int64, err error) {
	row, err := d.biliDMWriter.Exec(c, _uptTask, task.State, task.LastIndex, task.Count, task.Result, task.ID)
	if err != nil {
		log.Error("d.biliDMWriter.Exec(query:%s,task:%+v) error(%v)", _uptTask, task, err)
		return
	}
	return row.RowsAffected()
}

// UptSubTask uopdate dm sub task.
func (d *Dao) UptSubTask(c context.Context, taskID, delCount int64, end time.Time) (affected int64, err error) {
	row, err := d.biliDMWriter.Exec(c, _uptSubTask, delCount, end, taskID)
	if err != nil {
		log.Error("d.biliDMWriter.Exec(query:%s) error(%v)", _uptSubTask, err)
		return
	}
	return row.RowsAffected()
}

// SubTask .
func (d *Dao) SubTask(c context.Context, id int64) (subTask *model.SubTask, err error) {
	// TODO: operation time
	subTask = new(model.SubTask)
	row := d.biliDMWriter.QueryRow(c, _selectSubTask, id)
	if err = row.Scan(&subTask.ID, &subTask.Operation, &subTask.Rate, &subTask.Tcount, &subTask.Start, &subTask.End); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			subTask = nil
		}
		log.Error("biliDM.Scan(%s, taskID:%d) error*(%v)", _selectSubTask, id, err)
		return
	}
	return
}

// DelDMs dm task del dms.
func (d *Dao) DelDMs(c context.Context, oid int64, dmids []int64, state int32) (affected int64, err error) {
	rows, err := d.dmWriter.Exec(c, fmt.Sprintf(_delDMIndex, d.hitIndex(oid), xstr.JoinInts(dmids)), state, oid)
	if err != nil {
		log.Error("d.dmWriter.Exec(query:%s,oid:%d,dmids:%v) error(%v)", _delDMIndex, oid, dmids, err)
		return
	}
	return rows.RowsAffected()
}

// UptSubjectCount update count.
func (d *Dao) UptSubjectCount(c context.Context, tp int32, oid, count int64) (affected int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_uptSubCountSQL, d.hitSubject(oid)), count, tp, oid)
	if err != nil {
		log.Error("dmWriter.Exec(query:%s,oid:%d) error(%v)", _uptSubCountSQL, oid, err)
		return
	}
	return res.RowsAffected()
}

// TaskSearchRes get res from BI url
func (d *Dao) TaskSearchRes(c context.Context, task *model.TaskInfo) (count int64, result string, state int32, err error) {

	var (
		resp *http.Response
		res  struct {
			Code     int64    `json:"code"`
			StatusID int32    `json:"statusId"`
			Path     []string `json:"hdfsPath"`
			Count    int64    `json:"count"`
		}
		bs []byte
	)
	// may costing long time use default transport
	if resp, err = http.Get(task.Topic); err != nil {
		log.Error("http.Get(%s) error(%v)", task.Topic, err)
		return
	}
	defer resp.Body.Close()
	if bs, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error("ioutil.ReadAll url:%v error(%v)", task.Topic, err)
		return
	}
	if err = json.Unmarshal(bs, &res); err != nil {
		return
	}
	if res.Code != 200 {
		err = fmt.Errorf("%v", res)
		log.Error("d.httpClient.Get(%s) code(%d)", task.Topic, res.Code)
		return
	}
	if res.StatusID == model.TaskSearchSuc && len(res.Path) > 0 {
		result = res.Path[0]
		count = res.Count
	}
	return count, result, res.StatusID, err
}

// SendWechatWorkMsg send wechat work msg.
func (d *Dao) SendWechatWorkMsg(c context.Context, content, title string, users []string) (err error) {
	userMap := make(map[string]struct{}, len(users))
	unames := make([]string, 0, len(users))
	for _, user := range users {
		if user == "" {
			continue
		}
		if _, ok := userMap[user]; ok {
			continue
		}
		userMap[user] = struct{}{}
		unames = append(unames, user)
	}
	params := url.Values{}
	params.Set("Action", "CreateWechatMessage")
	params.Set("PublicKey", d.conf.TaskConf.MsgPublicKey)
	params.Set("UserName", strings.Join(unames, ","))
	params.Set("Title", title)
	params.Set("Content", content)
	params.Set("Signature", "")
	params.Set("TreeId", "")
	paramStr := params.Encode()
	if strings.IndexByte(paramStr, '+') > -1 {
		paramStr = strings.Replace(paramStr, "+", "%20", -1)
	}
	var (
		buffer bytes.Buffer
		querry string
	)
	buffer.WriteString(paramStr)
	querry = buffer.String()
	res := &struct {
		Code int `json:"RetCode"`
	}{}
	url := d.conf.Host.MerakHost + _merakMsgURI
	req, err := http.NewRequest("POST", url, strings.NewReader(querry))
	if err != nil {
		return
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	if err = d.httpCli.Do(c, req, res); err != nil {
		log.Error("d.SendWechatWorkMsg.client.Do(%v,%v) error(%v)", req, querry, err)
		return
	}
	if res.Code != 0 {
		log.Error("d.SendWechatWorkMsg(%s,%s,%v) res.Code != 0, res(%v)", content, title, users, res)
		err = fmt.Errorf("uri:%s,code:%d", url+querry, res.Code)
	}
	return
}
