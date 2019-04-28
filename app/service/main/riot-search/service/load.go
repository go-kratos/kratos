package service

import (
	"bufio"
	"context"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"

	"go-common/app/service/main/riot-search/model"
	"go-common/library/log"
)

func (s *Service) loadproc(path string) {
	log.Info("loading csv file %s", path)
	file, err := os.Open(path)
	if err != nil {
		log.Error("load csv file failed error:(%v)", err)
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = 2
	reader.LazyQuotes = true
	var arc []*model.Document
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error("read file error: %v line %v", err, line)
			// panic(err)
			continue
		}
		aid, err := strconv.ParseUint(line[0], 10, 64)
		if err != nil {
			log.Error("illegal line %v", line)
			panic(err)
		}
		title := line[1]
		arc = append(arc, &model.Document{
			ID:      aid,
			Content: title,
		})
	}
	log.Info("adding csv data to search engine...")
	s.pool.WaitCount(len(arc))
	for index := range arc {
		i := index
		s.pool.JobQueue <- func() {
			s.dao.Insert(arc[i].ID, arc[i].Content, false)
			s.pool.JobDone()
		}
	}
	s.pool.WaitAll()
	s.dao.Flush()
	/*
	 *********************
	 * load data from db *
	 *********************
	 */
	stime := time.Now().Add(-24 * time.Hour)
	log.Info("sync increment data (mtime>%v) from database", stime)
	for i := 0; i < 12; i++ {
		etime := stime.Add(2 * time.Hour)
		incrData, err := s.dao.IncrementBackup(context.Background(), stime, etime)
		stime = etime
		if err != nil {
			log.Error("database error:(%v)", err)
			s.dao.Flush()
			return
		}
		s.pool.WaitCount(len(incrData))
		for index := range incrData {
			i := index
			s.pool.JobQueue <- func() {
				s.dao.Insert(incrData[i].ID, incrData[i].Content, false)
				s.pool.JobDone()
			}
		}
		s.pool.WaitAll()
	}
	s.dao.Flush()
	log.Info("finish load data")
}
