package bootstrap

import "flag"

type CommandFlags struct {
	Conf       string
	Env        string
	ConfigHost string
	ConfigType string
}

func NewCommandFlags() *CommandFlags {
	return &CommandFlags{
		Conf:       "",
		Env:        "",
		ConfigHost: "",
		ConfigType: "",
	}
}

func (f *CommandFlags) Init() {
	flag.StringVar(&f.Conf, "conf", "../../configs", "config path, eg: -conf bootstrap.yaml")
	flag.StringVar(&f.Env, "env", "dev", "runtime environment, eg: -env dev")
	flag.StringVar(&f.ConfigHost, "chost", "127.0.0.1:8500", "config server host, eg: -chost 127.0.0.1:8500")
	flag.StringVar(&f.ConfigType, "ctype", "consul", "config server host, eg: -ctype consul")
}
