package feature

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

// Feature is the feature name
type Feature string

const (
	flagName = "feature-gates"
)

var (
	// DefaultGate is a shared global Gate.
	DefaultGate = NewGate()
)

// Spec is the spec of the feature
type Spec struct {
	Default bool
}

// Gate parses and stores flag gates for known features from
// a string like feature1=true,feature2=false,...
type Gate interface {
	// AddFlag adds a flag for setting global feature gates to the specified FlagSet.
	AddFlag(fs *flag.FlagSet)
	// Set parses and stores flag gates for known features
	// from a string like feature1=true,feature2=false,...
	Set(value string) error
	// SetFromMap stores flag gates for known features from a map[string]bool or returns an error
	SetFromMap(m map[string]bool) error
	// Enabled returns true if the key is enabled.
	Enabled(key Feature) bool
	// Add adds features to the featureGate.
	Add(features map[Feature]Spec) error
	// KnownFeatures returns a slice of strings describing the Gate's known features.
	KnownFeatures() []string
	// DeepCopy returns a deep copy of the Gate object, such that gates can be
	// set on the copy without mutating the original. This is useful for validating
	// config against potential feature gate changes before committing those changes.
	DeepCopy() Gate
}

// featureGate implements Gate as well as flag.Value for flag parsing.
type featureGate struct {
	// lock guards writes to known, enabled, and reads/writes of closed
	lock sync.Mutex
	// known holds a map[Feature]Spec
	known atomic.Value
	// enabled holds a map[Feature]bool
	enabled atomic.Value
	// closed is set to true when AddFlag is called, and prevents subsequent calls to Add
	closed bool
}

// Set, String, and Type implement flag.Value
var _ flag.Value = &featureGate{}

// NewGate create a feature gate.
func NewGate() *featureGate {
	known := map[Feature]Spec{}
	knownValue := atomic.Value{}
	knownValue.Store(known)

	enabled := map[Feature]bool{}
	enabledValue := atomic.Value{}
	enabledValue.Store(enabled)

	f := &featureGate{
		known:   knownValue,
		enabled: enabledValue,
	}
	return f
}

// Set parses a string of the form "key1=value1,key2=value2,..." into a
// map[string]bool of known keys or returns an error.
func (f *featureGate) Set(value string) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	// Copy existing state
	known := map[Feature]Spec{}
	for k, v := range f.known.Load().(map[Feature]Spec) {
		known[k] = v
	}
	enabled := map[Feature]bool{}
	for k, v := range f.enabled.Load().(map[Feature]bool) {
		enabled[k] = v
	}

	for _, s := range strings.Split(value, ",") {
		if len(s) == 0 {
			continue
		}
		arr := strings.SplitN(s, "=", 2)
		k := Feature(strings.TrimSpace(arr[0]))
		_, ok := known[k]
		if !ok {
			return errors.Errorf("unrecognized key: %s", k)
		}
		if len(arr) != 2 {
			return errors.Errorf("missing bool value for %s", k)
		}
		v := strings.TrimSpace(arr[1])
		boolValue, err := strconv.ParseBool(v)
		if err != nil {
			return errors.Errorf("invalid value of %s: %s, err: %v", k, v, err)
		}
		enabled[k] = boolValue
	}

	// Persist changes
	f.known.Store(known)
	f.enabled.Store(enabled)

	fmt.Fprintf(os.Stderr, "feature gates: %v", enabled)
	return nil
}

// SetFromMap stores flag gates for known features from a map[string]bool or returns an error
func (f *featureGate) SetFromMap(m map[string]bool) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	// Copy existing state
	known := map[Feature]Spec{}
	for k, v := range f.known.Load().(map[Feature]Spec) {
		known[k] = v
	}
	enabled := map[Feature]bool{}
	for k, v := range f.enabled.Load().(map[Feature]bool) {
		enabled[k] = v
	}

	for k, v := range m {
		k := Feature(k)
		_, ok := known[k]
		if !ok {
			return errors.Errorf("unrecognized key: %s", k)
		}
		enabled[k] = v
	}

	// Persist changes
	f.known.Store(known)
	f.enabled.Store(enabled)

	fmt.Fprintf(os.Stderr, "feature gates: %v", f.enabled)
	return nil
}

// String returns a string containing all enabled feature gates, formatted as "key1=value1,key2=value2,...".
func (f *featureGate) String() string {
	pairs := []string{}
	enabled, ok := f.enabled.Load().(map[Feature]bool)
	if !ok {
		return ""
	}
	for k, v := range enabled {
		pairs = append(pairs, fmt.Sprintf("%s=%t", k, v))
	}
	sort.Strings(pairs)
	return strings.Join(pairs, ",")
}

func (f *featureGate) Type() string {
	return "mapStringBool"
}

// Add adds features to the featureGate.
func (f *featureGate) Add(features map[Feature]Spec) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.closed {
		return errors.Errorf("cannot add a feature gate after adding it to the flag set")
	}

	// Copy existing state
	known := map[Feature]Spec{}
	for k, v := range f.known.Load().(map[Feature]Spec) {
		known[k] = v
	}

	for name, spec := range features {
		if existingSpec, found := known[name]; found {
			if existingSpec == spec {
				continue
			}
			return errors.Errorf("feature gate %q with different spec already exists: %v", name, existingSpec)
		}

		known[name] = spec
	}

	// Persist updated state
	f.known.Store(known)

	return nil
}

// Enabled returns true if the key is enabled.
func (f *featureGate) Enabled(key Feature) bool {
	if v, ok := f.enabled.Load().(map[Feature]bool)[key]; ok {
		return v
	}
	return f.known.Load().(map[Feature]Spec)[key].Default
}

// AddFlag adds a flag for setting global feature gates to the specified FlagSet.
func (f *featureGate) AddFlag(fs *flag.FlagSet) {
	f.lock.Lock()
	// TODO(mtaufen): Shouldn't we just close it on the first Set/SetFromMap instead?
	// Not all components expose a feature gates flag using this AddFlag method, and
	// in the future, all components will completely stop exposing a feature gates flag,
	// in favor of componentconfig.
	f.closed = true
	f.lock.Unlock()

	known := f.KnownFeatures()
	fs.Var(f, flagName, ""+
		"A set of key=value pairs that describe feature gates for alpha/experimental features. "+
		"Options are:\n"+strings.Join(known, "\n"))
}

// KnownFeatures returns a slice of strings describing the Gate's known features.
func (f *featureGate) KnownFeatures() []string {
	var known []string
	for k, v := range f.known.Load().(map[Feature]Spec) {
		known = append(known, fmt.Sprintf("%s=true|false (default=%t)", k, v.Default))
	}
	sort.Strings(known)
	return known
}

// DeepCopy returns a deep copy of the Gate object, such that gates can be
// set on the copy without mutating the original. This is useful for validating
// config against potential feature gate changes before committing those changes.
func (f *featureGate) DeepCopy() Gate {
	// Copy existing state.
	known := map[Feature]Spec{}
	for k, v := range f.known.Load().(map[Feature]Spec) {
		known[k] = v
	}
	enabled := map[Feature]bool{}
	for k, v := range f.enabled.Load().(map[Feature]bool) {
		enabled[k] = v
	}

	// Store copied state in new atomics.
	knownValue := atomic.Value{}
	knownValue.Store(known)
	enabledValue := atomic.Value{}
	enabledValue.Store(enabled)

	// Construct a new featureGate around the copied state.
	// We maintain the value of f.closed across the copy.
	return &featureGate{
		known:   knownValue,
		enabled: enabledValue,
		closed:  f.closed,
	}
}
