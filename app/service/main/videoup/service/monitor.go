package service

//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"strings"
//	"time"
//
//	arcmdl "go-common/app/service/videoup/model/archive"
//	"go-common/log"
//)

//monitorConsume
func (s *Service) monitorConsume() {
	// return
	// log.Info("--s.monitorConsume start")
	// if s.c.Monitor.Env != "pro" {
	// 	log.Info("[env check]env (%v)", s.c.Monitor.Env)
	// 	return
	// }
	// log.Info("[before loop]env (%v) and  ", s.c.Monitor.Env)
	// for {
	// 	send := false
	// 	time.Sleep(10 * time.Minute)
	// 	s.locker.Lock()
	// 	bs, _ := json.Marshal(s.monitorMap)
	// 	log.Info("[loop]s.monitorConsume monitorMap json (%s)", bs)
	// 	for _, item := range s.monitorMap {
	// 		if item.Value >= item.Limit {
	// 			send = true
	// 		}
	// 	}
	// 	s.monitorMap = make(map[string]*arcmdl.Alert)
	// 	s.locker.Unlock()
	// 	if send {
	// 		log.Info("s.monitorConsume  retry too many hit")
	// 		if err := s.monitor.Send(context.TODO(), fmt.Sprintf("video-service  retry too many within 10 minute")); err != nil {
	// 			log.Info("[sms]s.monitorConsume  sms error(%v)", err)
	// 		}
	// 	}
	// }
}

//addMonitor
func (s *Service) addMonitor(pre, key string) {
	// return
	// keyJoin := strings.Join([]string{pre, key}, "-")
	// s.locker.Lock()
	// if item, ok := s.monitorMap[keyJoin]; !ok {
	// 	s.monitorMap[keyJoin] = &arcmdl.Alert{keyJoin, 1, s.c.Monitor.Count}
	// } else {
	// 	item.Value++
	// }
	// s.locker.Unlock()
	// log.Info("s.addMonitor hit (%v)", keyJoin)
}
