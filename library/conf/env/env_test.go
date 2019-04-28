package env

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestDefaultString(t *testing.T) {
	v := defaultString("a", "test")
	if v != "test" {
		t.Fatal("v must be test")
	}
	if err := os.Setenv("a", "test1"); err != nil {
		t.Fatal(err)
	}
	v = defaultString("a", "test")
	if v != "test1" {
		t.Fatal("v must be test1")
	}
}

func TestEnv(t *testing.T) {
	tests := []struct {
		flag string
		env  string
		def  string
		val  *string
	}{
		{
			"region",
			"REGION",
			_region,
			&Region,
		},
		{
			"zone",
			"ZONE",
			_zone,
			&Zone,
		},
		{
			"deploy.env",
			"DEPLOY_ENV",
			_deployEnv,
			&DeployEnv,
		},
		{
			"appid",
			"APP_ID",
			"",
			&AppID,
		},
		{
			"http.port",
			"DISCOVERY_HTTP_PORT",
			_httpPort,
			&HTTPPort,
		},
		{
			"gorpc.port",
			"DISCOVERY_GORPC_PORT",
			_gorpcPort,
			&GORPCPort,
		},
		{
			"grpc.port",
			"DISCOVERY_GRPC_PORT",
			_grpcPort,
			&GRPCPort,
		},
		{
			"deploy.color",
			"DEPLOY_COLOR",
			"",
			&Color,
		},
	}
	for _, test := range tests {
		// flag set value
		t.Run(fmt.Sprintf("%s: flag set", test.env), func(t *testing.T) {
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			addFlag(fs)
			err := fs.Parse([]string{fmt.Sprintf("-%s=%s", test.flag, "test")})
			if err != nil {
				t.Fatal(err)
			}
			if *test.val != "test" {
				t.Fatal("val must be test")
			}
		})
		// flag not set, env set
		t.Run(fmt.Sprintf("%s: flag not set, env set", test.env), func(t *testing.T) {
			*test.val = ""
			os.Setenv(test.env, "test2")
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			addFlag(fs)
			err := fs.Parse([]string{})
			if err != nil {
				t.Fatal(err)
			}
			if *test.val != "test2" {
				t.Fatal("val must be test")
			}
		})
		// flag not set, env not set
		t.Run(fmt.Sprintf("%s: flag not set, env not set", test.env), func(t *testing.T) {
			*test.val = ""
			os.Setenv(test.env, "")
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			addFlag(fs)
			err := fs.Parse([]string{})
			if err != nil {
				t.Fatal(err)
			}
			if *test.val != test.def {
				t.Fatal("val must be test")
			}
		})
	}
}
