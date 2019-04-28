package output

import (
	"errors"
	"go-common/app/service/ops/log-agent/conf/configcenter"
)

const (
	lancerConfig = "output.toml"
)

// ReadConfig read config for default output
func ReadConfig() (value string, err error) {
	var ok bool

	// logLevel config
	if value, ok = configcenter.Client.Value(lancerConfig); !ok {
		return "", errors.New("failed to get output.toml")
	} else {
		return value, nil
	}

	return value, nil
}
