package worker

import (
	"testing"
	"time"
)

func TestIncrease(t *testing.T) {
	var (
		conf = &Conf{
			QueueSize:     10,
			WorkerProcMax: 10,
			WorkerNumber:  1,
		}
		workerPool = New(conf)
	)

	for i := 0; i < 10; i++ {
		workerPool.Add(longtime)
	}
	time.Sleep(6 * time.Second)
	var expect = minInt(conf.WorkerNumber<<1, conf.WorkerProcMax)
	if workerPool.workerNumber != expect {
		t.Logf("worker number=%d, expect=%d", workerPool.workerNumber, expect)
		t.FailNow()
	}

	for i := 0; i < 10; i++ {
		workerPool.Add(longtime)
	}
	time.Sleep(6 * time.Second)
	expect = minInt(conf.WorkerNumber<<2, conf.WorkerProcMax)
	if workerPool.workerNumber != expect {
		t.Logf("worker number=%d, expect=%d", workerPool.workerNumber, expect)
		t.FailNow()
	}

	for i := 0; i < 10; i++ {
		workerPool.Add(longtime)
	}
	time.Sleep(6 * time.Second)
	expect = minInt(conf.WorkerNumber<<3, conf.WorkerProcMax)
	if workerPool.workerNumber != expect {
		t.Logf("worker number=%d, expect=%d", workerPool.workerNumber, expect)
		t.FailNow()
	}

	for i := 0; i < 10; i++ {
		workerPool.Add(longtime)
	}
	time.Sleep(6 * time.Second)
	expect = minInt(conf.WorkerNumber<<4, conf.WorkerProcMax)
	if workerPool.workerNumber != expect {
		t.Logf("worker number=%d, expect=%d", workerPool.workerNumber, expect)
		t.FailNow()
	}
}

func longtime() {
	time.Sleep(20 * time.Second)
}
