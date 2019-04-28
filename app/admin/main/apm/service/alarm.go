package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/databus"
	"go-common/library/ecode"
	"go-common/library/log"
	"net/http"
)

//Alarm ...
// func (s *Service) Alarm(c context.Context, group, action string) (result *databus.AlarmOpen, err error) {
// 	url := conf.Conf.Alarm.DatabusURL
// 	var jsonBytes []byte
// 	body := &struct {
// 		Action    string `json:"Action"`
// 		PublicKey string `json:"PublicKey"`
// 		Signature int8   `json:"Signature"`
// 		Group     string `json:"Group"`
// 	}{
// 		Action:    action,
// 		PublicKey: conf.Conf.Alarm.DatabusKey,
// 		Signature: 1,
// 		Group:     group,
// 	}
// 	if jsonBytes, err = json.Marshal(body); err != nil {
// 		log.Error("json.Marshal(body) error(%v)", err)
// 		return
// 	}
// 	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBytes)))
// 	if err != nil {
// 		err = ecode.RequestErr
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	// req.Header.Set("Cookie", cookie)
// 	result = &databus.AlarmOpen{}
// 	if err = s.client.Do(c, req, result); err != nil {
// 		fmt.Printf("result=(%v) error=(%v)", result, err)
// 		log.Error("Alarm() error(%v)", err)
// 		err = ecode.RequestErr
// 		return
// 	}
// 	return
// }

//Opsmind ...
func (s *Service) Opsmind(c context.Context, project, group, action, Owners string, percentage, fortime int64, silence bool) (result *databus.Res, err error) {
	var scopes []databus.Scope
	scopes = append(scopes, databus.Scope{Type: 0, Key: "group", Val: []string{group}})
	var owner databus.Owner
	owner.Owner = Owners
	owner.App = project
	url := conf.Conf.Alarm.DatabusURL
	body := &struct {
		Action           string          `json:"Action"`
		PublicKey        string          `json:"PublicKey"`
		Signature        int8            `json:"Signature"`
		PolicyID         string          `json:"PolicyId"`
		CateGory         string          `json:"CateGory"`
		Silence          bool            `json:"Silence"`
		Scope            []databus.Scope `json:"Scope"`
		TriggerName      string          `json:"TriggerName"`
		TriggerOperator  string          `json:"TriggerOperator"`
		TriggerLevel     string          `json:"TriggerLevel"`
		TriggerThreshold int64           `json:"TriggerThreshold"`
		TriggerFor       int64           `json:"TriggerFor"`
		TriggerNoDataFor int64           `json:"TriggerNoDataFor"`
		TriggerNotes     databus.Owner   `json:"TriggerNotes"`
	}{
		Action:           action,
		PublicKey:        conf.Conf.Alarm.DatabusKey,
		Signature:        1,
		PolicyID:         "3mwipx2caggxc",
		CateGory:         "Databus",
		Silence:          silence,
		Scope:            scopes,
		TriggerName:      fmt.Sprintf("%s消费落后告警", group),
		TriggerOperator:  "<",
		TriggerLevel:     "P3",
		TriggerThreshold: percentage,
		TriggerFor:       fortime,
		TriggerNoDataFor: 300,
		TriggerNotes:     owner,
	}
	// json.NewEncoder(os.Stdout).Encode(body)
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(body)
	if err != nil {
		log.Error("json.Marshal(body) error(%v)", err)
		return
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		err = ecode.RequestErr
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Cookie", cookie)
	result = &databus.Res{}
	if err = s.client.Do(c, req, result); err != nil {
		log.Error("Alarm() error(%v)", err)
		err = ecode.RequestErr
		return
	}
	return
}

// OpsmindRemove ...
func (s *Service) OpsmindRemove(c context.Context, adjustid, action string) (result *databus.Res, err error) {
	url := conf.Conf.Alarm.DatabusURL
	body := &struct {
		Action    string `json:"Action"`
		PublicKey string `json:"PublicKey"`
		Signature int8   `json:"Signature"`
		PolicyID  string `json:"PolicyId"`
		AdjustID  string `json:"AdjustId"`
	}{
		Action:    action,
		PublicKey: conf.Conf.Alarm.DatabusKey,
		Signature: 1,
		PolicyID:  "3mwipx2caggxc",
		AdjustID:  adjustid,
	}
	// json.NewEncoder(os.Stdout).Encode(body)
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(body)
	if err != nil {
		log.Error("json.Marshal(body) error(%v)", err)
		return
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		err = ecode.RequestErr
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Cookie", cookie)
	result = &databus.Res{}
	if err = s.client.Do(c, req, result); err != nil {
		log.Error("Alarm() error(%v)", err)
		err = ecode.RequestErr
		return
	}
	return
}

//OpsmindQuery ...
func (s *Service) OpsmindQuery(c context.Context, group, action string) (result *databus.ResQuery, err error) {
	var query []databus.Query
	query = append(query, databus.Query{Key: "group", Val: []string{group}})
	url := conf.Conf.Alarm.DatabusURL
	body := &struct {
		Action    string          `json:"Action"`
		PublicKey string          `json:"PublicKey"`
		Signature int8            `json:"Signature"`
		Query     []databus.Query `json:"Query"`
	}{
		Action:    action,
		PublicKey: conf.Conf.Alarm.DatabusKey,
		Signature: 1,
		Query:     query,
	}
	// json.NewEncoder(os.Stdout).Encode(body)
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(body)
	if err != nil {
		log.Error("json.Marshal(body) error(%v)", err)
		return
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		err = ecode.RequestErr
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Cookie", cookie)
	result = &databus.ResQuery{}
	if err = s.client.Do(c, req, result); err != nil {
		log.Error("Alarm() error(%v)", err)
		err = ecode.RequestErr
	}
	return
}
