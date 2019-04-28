package log_test

import (
	"context"

	"go-common/library/log"
)

// This example will logging a text to log file.
func ExampleInfo() {
	fc := &log.Config{
		Family: "test-log",
		Dir:    "/data/log/test",
	}
	log.Init(fc)
	defer log.Close()
	log.Info("test %s", "file log")

	ac := &log.Config{
		Family: "test-log",
		Agent: &log.AgentConfig{
			TaskID: "000003",
			Addr:   "172.16.0.204:514",
			Proto:  "tcp",
			Chan:   1024,
		},
	}
	log.Init(ac)
	defer log.Close()
	log.Info("test %s", "agent log")
}

// This example will logging a structured text to log agent.
func ExampleInfov() {
	ac := &log.Config{
		Family: "test-log",
		Agent: &log.AgentConfig{
			TaskID: "000003",
			Addr:   "172.16.0.204:514",
			Proto:  "tcp",
			Chan:   1024,
		},
	}
	log.Init(ac)
	defer log.Close()
	log.Infov(context.TODO(), log.KV("key1", "val1"), log.KV("key2", "val2"))
}

// This example will set log format
func ExampleSetFormat() {
	log.SetFormat("%L %T %f %M")
	log.Info("hello")
	// log output:
	// INFO 2018-06-28T12:15:48.713784 main.main:8 hello
}
