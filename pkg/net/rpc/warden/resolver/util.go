package resolver

import (
	"flag"
	"fmt"
)

// RegisterTarget will register grpc discovery mock address flag
func RegisterTarget(target *string, discoveryID string) {
	flag.CommandLine.StringVar(
		target,
		fmt.Sprintf("grpc.%s", discoveryID),
		fmt.Sprintf("discovery://default/%s", discoveryID),
		fmt.Sprintf("App's grpc target.\n example: -grpc.%s=\"127.0.0.1:9090\"", discoveryID),
	)
}
