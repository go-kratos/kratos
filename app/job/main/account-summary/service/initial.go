package service

// import (
// 	"context"
// 	"sync"
// 	"time"

// 	"go-common/library/log"
// )

// func (s *Service) initialproc(ctx context.Context) {
// 	dataWg := sync.WaitGroup{}
// 	workerWg := sync.WaitGroup{}

// 	jobQueue := make(chan func(), 4096)
// 	initialWorker := func() {
// 		worker := uint64(50)
// 		if s.c.AccountSummary.InitialWriteWorker > 0 {
// 			worker = s.c.AccountSummary.InitialWriteWorker
// 		}
// 		log.Info("Start %d initial write worker", worker)
// 		for i := uint64(0); i < worker; i++ {
// 			workerWg.Add(1)
// 			go func() {
// 				defer workerWg.Done()
// 				for job := range jobQueue {
// 					job()
// 				}
// 			}()
// 		}
// 	}
// 	initialWorker()

// 	initBase := func() {
// 		log.Info("Start to initial member base")
// 		defer dataWg.Done()
// 		baseCh := s.dao.AllMemberBase(ctx)
// 		for chunk := range baseCh {
// 			for _, b := range chunk {
// 				b := b
// 				jobQueue <- func() {
// 					if err := s.SyncToHBase(ctx, b); err != nil {
// 						log.Error("Failed to sync member base in initial process: base: %+v: %+v", b, err)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	initExp := func() {
// 		log.Info("Start to initial member exp")
// 		defer dataWg.Done()
// 		expCh := s.dao.AllMemberExp(ctx)
// 		for chunk := range expCh {
// 			for _, e := range chunk {
// 				e := e
// 				jobQueue <- func() {
// 					if err := s.SyncToHBase(ctx, e); err != nil {
// 						log.Error("Failed to sync member exp in initial process: exp: %+v: %+v", e, err)
// 					}
// 				}
// 			}
// 			time.Sleep(time.Second)
// 		}
// 	}
// 	initOfficial := func() {
// 		log.Info("Start to initial member official")
// 		defer dataWg.Done()
// 		official, err := s.dao.AllOfficial(ctx)
// 		if err != nil {
// 			log.Error("Failed to get all member official: %+v", err)
// 			return
// 		}
// 		for _, o := range official {
// 			o := o
// 			jobQueue <- func() {
// 				if err := s.SyncToHBase(ctx, o); err != nil {
// 					log.Error("Failed to sync member official in initial process: official: %+v: %+v", o, err)
// 				}
// 			}
// 		}
// 	}
// 	initStat := func() {
// 		log.Info("Start to initial relation stat")
// 		defer dataWg.Done()
// 		statsCh := s.dao.AllRelationStat(ctx)
// 		for chunk := range statsCh {
// 			for _, stat := range chunk {
// 				stat := stat
// 				jobQueue <- func() {
// 					if err := s.SyncToHBase(ctx, stat); err != nil {
// 						log.Error("Failed to sync relation stat in initial process: stat: %+v: %+v", stat, err)
// 					}
// 				}
// 			}
// 			time.Sleep(time.Second)
// 		}
// 	}

// 	if s.c.FeatureGate.InitialMemberBase {
// 		dataWg.Add(1)
// 		go initBase()
// 	}

// 	if s.c.FeatureGate.InitialMemberExp {
// 		dataWg.Add(1)
// 		go initExp()
// 	}

// 	if s.c.FeatureGate.InitialMemberOfficial {
// 		dataWg.Add(1)
// 		go initOfficial()
// 	}

// 	if s.c.FeatureGate.InitialRelationStat {
// 		dataWg.Add(1)
// 		go initStat()
// 	}

// 	dataWg.Wait()
// 	close(jobQueue) // all job is enqueued
// 	workerWg.Wait()
// }
