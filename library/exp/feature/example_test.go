package feature_test

import (
	"flag"
	"fmt"

	"go-common/library/exp/feature"
)

var (
	AStableFeature  feature.Feature = "a-stable-feature"
	AStagingFeature feature.Feature = "a-staging-feature"
)

var exampleFeatures = map[feature.Feature]feature.Spec{
	AStableFeature:  feature.Spec{Default: true},
	AStagingFeature: feature.Spec{Default: false},
}

func init() {
	feature.DefaultGate.Add(exampleFeatures)
	feature.DefaultGate.AddFlag(flag.CommandLine)
}

// This example create an example to using default features.
func Example() {
	knows := feature.DefaultGate.KnownFeatures()
	fmt.Println(knows)

	enabled := feature.DefaultGate.Enabled(AStableFeature)
	fmt.Println(enabled)

	enabled = feature.DefaultGate.Enabled(AStagingFeature)
	fmt.Println(enabled)
	// Output: [a-stable-feature=true|false (default=true) a-staging-feature=true|false (default=false)]
	// true
	// false
}

// This example parsing flag from command line and enable a staging feature.
func ExampleFeature() {
	knows := feature.DefaultGate.KnownFeatures()
	fmt.Println(knows)

	enabled := feature.DefaultGate.Enabled(AStagingFeature)
	fmt.Println(enabled)

	flag.Set("feature-gates", fmt.Sprintf("%s=true", AStagingFeature))
	enabled = feature.DefaultGate.Enabled(AStagingFeature)
	fmt.Println(enabled)
	// Output: [a-stable-feature=true|false (default=true) a-staging-feature=true|false (default=false)]
	// false
	// true
}
