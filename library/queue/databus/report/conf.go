package report

import (
	"time"

	"go-common/library/conf/env"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

var (
	_defaultManagerConfig = map[string]*conf{
		env.DeployEnvFat1: {
			Secret: "971d048a2818e37ae124a0293c300e89",
			Addr:   "172.16.33.158:6205",
		},
		env.DeployEnvUat: {
			Secret: "cde3b480836cc76df3d635470f991caa",
			Addr:   "172.18.33.50:6205",
		},
		env.DeployEnvPre: {
			Secret: "4a933f8170a4711ace8ec363f7e5f23c",
			Addr:   "172.18.21.90:6205",
		},
		env.DeployEnvProd: {
			Secret: "4a933f8170a4711ace8ec363f7e5f23c",
			Addr:   "172.18.21.90:6205",
		},
	}

	_defaultUserConfig = map[string]*conf{
		env.DeployEnvFat1: {
			Secret: "971d048a2818e37ae124a0293c300e89",
			Addr:   "172.16.33.158:6205",
		},
		env.DeployEnvUat: {
			Secret: "cde3b480836cc76df3d635470f991caa",
			Addr:   "172.18.33.50:6205",
		},
		env.DeployEnvPre: {
			Secret: "4a933f8170a4711ace8ec363f7e5f23c",
			Addr:   "172.18.21.90:6205",
		},
		env.DeployEnvProd: {
			Secret: "4a933f8170a4711ace8ec363f7e5f23c",
			Addr:   "172.18.21.90:6205",
		},
	}

	_managerConfig = &databus.Config{
		Key:          "2511663d546f1413",
		Secret:       "971d048a2818e37ae124a0293c300e89",
		Group:        "LogAudit-MainSearch-P",
		Topic:        "LogAudit-T",
		Action:       "pub",
		Buffer:       10240,
		Name:         "log-audit/log-sub",
		Proto:        "tcp",
		Addr:         "172.16.33.158:6205",
		Active:       100,
		Idle:         10,
		DialTimeout:  xtime.Duration(time.Millisecond * 200),
		ReadTimeout:  xtime.Duration(time.Millisecond * 200),
		WriteTimeout: xtime.Duration(time.Millisecond * 200),
		IdleTimeout:  xtime.Duration(time.Second * 80),
	}

	_userConfig = &databus.Config{
		Key:          "2511663d546f1413",
		Secret:       "971d048a2818e37ae124a0293c300e89",
		Group:        "LogUserAction-MainSearch-P",
		Topic:        "LogUserAction-T",
		Action:       "pub",
		Buffer:       10240,
		Name:         "log-user-action/log-sub",
		Proto:        "tcp",
		Addr:         "172.16.33.158:6205",
		Active:       100,
		Idle:         10,
		DialTimeout:  xtime.Duration(time.Millisecond * 200),
		ReadTimeout:  xtime.Duration(time.Millisecond * 200),
		WriteTimeout: xtime.Duration(time.Millisecond * 200),
		IdleTimeout:  xtime.Duration(time.Second * 80),
	}
)
