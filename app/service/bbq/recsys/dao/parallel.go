package dao

import (
	"go-common/app/service/bbq/recsys/dao/parallel"

	"github.com/Dai0522/workpool"
)

// parallelTask2 .
func (d *Dao) parallelTask2(tasks map[string]workpool.Task) map[string]workpool.FutureTask {
	ftMap := make(map[string]workpool.FutureTask)
	for name, task := range tasks {
		ft := workpool.NewFutureTask(task)

		retry := 0
		err := d.wp.Submit(ft)
		for err != nil && retry < 3 {
			err = d.wp.Submit(ft)
			retry++
		}
		ftMap[name] = *ft
	}
	return ftMap
}

// parallelTask .
func (d *Dao) parallelTask(tasks []workpool.Task) *[]workpool.FutureTask {
	ftArr := make([]workpool.FutureTask, len(tasks))
	for i := range tasks {
		ft := workpool.NewFutureTask(tasks[i])

		retry := 0
		err := d.wp.Submit(ft)
		for err != nil && retry < 3 {
			err = d.wp.Submit(ft)
			retry++
		}
		ftArr[i] = *ft
	}
	return &ftArr
}

// ParallelRedis run redis cmd parallel
func (d *Dao) ParallelRedis(tasks *[]parallel.RedisTask) *[]workpool.FutureTask {
	ftArr := make([]workpool.FutureTask, len(*tasks))
	for i := range *tasks {
		ft := workpool.NewFutureTask(&(*tasks)[i])

		retry := 0
		err := d.wp.Submit(ft)
		for err != nil && retry < 3 {
			err = d.wp.Submit(ft)
			retry++
		}
		ftArr[i] = *ft
	}
	return &ftArr
}
