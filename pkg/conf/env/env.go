// Package env get env & app config, all the public field must after init()
// finished and flag.Parse().
package env

import (
	"flag"
	"os"
)

// deploy env.
const (
	DeployEnvDev  = "dev"
	DeployEnvFat  = "fat"
	DeployEnvUat  = "uat"
	DeployEnvPre  = "pre"
	DeployEnvProd = "prod"
)

// env default value.
const (
	// env
	_region    = "region01"
	_zone      = "zone01"
	_deployEnv = "dev"
)

// env configuration.
var (
	// Region avaliable region where app at.
	Region string
	// Zone avaliable zone where app at.
	Zone string
	// Hostname machine hostname.
	Hostname string
	// DeployEnv deploy env where app at.
	DeployEnv string
	// AppID is global unique application id, register by service tree.
	// such as main.arch.disocvery.
	AppID string
	// Color is the identification of different experimental group in one caster cluster.
	Color string
)

func init() {
	var err error
	if Hostname, err = os.Hostname(); err != nil || Hostname == "" {
		Hostname = os.Getenv("HOSTNAME")
	}

	addFlag(flag.CommandLine)
}

func addFlag(fs *flag.FlagSet) {
	// env
	fs.StringVar(&Region, "region", defaultString("REGION", _region), "avaliable region. or use REGION env variable, value: sh etc.")
	fs.StringVar(&Zone, "zone", defaultString("ZONE", _zone), "avaliable zone. or use ZONE env variable, value: sh001/sh002 etc.")
	fs.StringVar(&AppID, "appid", os.Getenv("APP_ID"), "appid is global unique application id, register by service tree. or use APP_ID env variable.")
	fs.StringVar(&DeployEnv, "deploy.env", defaultString("DEPLOY_ENV", _deployEnv), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
	fs.StringVar(&Color, "deploy.color", os.Getenv("DEPLOY_COLOR"), "deploy.color is the identification of different experimental group.")
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}
