package election

import (
	"go-common/library/log"
	"testing"
	"time"
)

func TestElection(t *testing.T) {
	//log.AddFilter("all", log.DEBUG, log.NewConsoleLogWriter())
	var hosts = []string{"172.18.33.50:2199"}
	var root = "/microservice/upcredit-service/nodes"

	var elect = New(hosts, root, time.Second*5)
	var err = elect.Init()
	if err != nil {
		log.Error("fail to init elect")
		return
	}

	elect.Elect()
	for {
		isMaster := <-elect.C
		if isMaster {
			log.Info("this is master, node=%s", elect.NodePath)
		} else {
			log.Info("this is follower, node=%s, master=%s", elect.NodePath, elect.MasterPath)
		}
	}
}
