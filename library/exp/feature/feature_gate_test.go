package feature

import (
	"flag"
	"fmt"
	"strings"
	"testing"
)

func TestFeatureGateOverride(t *testing.T) {
	const testAlphaGate Feature = "TestAlpha"
	const testBetaGate Feature = "TestBeta"

	// Don't parse the flag, assert defaults are used.
	var f Gate = NewGate()
	f.Add(map[Feature]Spec{
		testAlphaGate: {Default: false},
		testBetaGate:  {Default: false},
	})

	f.Set("TestAlpha=true,TestBeta=true")
	if f.Enabled(testAlphaGate) != true {
		t.Errorf("Expected true")
	}
	if f.Enabled(testBetaGate) != true {
		t.Errorf("Expected true")
	}

	f.Set("TestAlpha=false")
	if f.Enabled(testAlphaGate) != false {
		t.Errorf("Expected false")
	}
	if f.Enabled(testBetaGate) != true {
		t.Errorf("Expected true")
	}
}

func TestFeatureGateFlagDefaults(t *testing.T) {
	// gates for testing
	const testAlphaGate Feature = "TestAlpha"
	const testBetaGate Feature = "TestBeta"

	// Don't parse the flag, assert defaults are used.
	var f Gate = NewGate()
	f.Add(map[Feature]Spec{
		testAlphaGate: {Default: false},
		testBetaGate:  {Default: true},
	})

	if f.Enabled(testAlphaGate) != false {
		t.Errorf("Expected false")
	}
	if f.Enabled(testBetaGate) != true {
		t.Errorf("Expected true")
	}
}

func TestFeatureGateSetFromMap(t *testing.T) {
	// gates for testing
	const testAlphaGate Feature = "TestAlpha"
	const testBetaGate Feature = "TestBeta"

	tests := []struct {
		name        string
		setmap      map[string]bool
		expect      map[Feature]bool
		setmapError string
	}{
		{
			name: "set TestAlpha and TestBeta true",
			setmap: map[string]bool{
				"TestAlpha": true,
				"TestBeta":  true,
			},
			expect: map[Feature]bool{
				testAlphaGate: true,
				testBetaGate:  true,
			},
		},
		{
			name: "set TestBeta true",
			setmap: map[string]bool{
				"TestBeta": true,
			},
			expect: map[Feature]bool{
				testAlphaGate: false,
				testBetaGate:  true,
			},
		},
		{
			name: "set TestAlpha false",
			setmap: map[string]bool{
				"TestAlpha": false,
			},
			expect: map[Feature]bool{
				testAlphaGate: false,
				testBetaGate:  false,
			},
		},
		{
			name: "set TestInvaild true",
			setmap: map[string]bool{
				"TestInvaild": true,
			},
			expect: map[Feature]bool{
				testAlphaGate: false,
				testBetaGate:  false,
			},
			setmapError: "unrecognized key:",
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("SetFromMap %s", test.name), func(t *testing.T) {
			f := NewGate()
			f.Add(map[Feature]Spec{
				testAlphaGate: {Default: false},
				testBetaGate:  {Default: false},
			})
			err := f.SetFromMap(test.setmap)
			if test.setmapError != "" {
				if !strings.Contains(err.Error(), test.setmapError) {
					t.Errorf("%d: SetFromMap(%#v) Expected err:%v, Got err:%v", i, test.setmap, test.setmapError, err)
				}
			} else if err != nil {
				t.Errorf("%d: SetFromMap(%#v) Expected success, Got err:%v", i, test.setmap, err)
			}
			for k, v := range test.expect {
				if actual := f.Enabled(k); actual != v {
					t.Errorf("%d: SetFromMap(%#v) Expected %s=%v, Got %s=%v", i, test.setmap, k, v, k, actual)
				}
			}
		})
	}
}

func TestFeatureGateString(t *testing.T) {
	// gates for testing
	const testAlphaGate Feature = "TestAlpha"
	const testBetaGate Feature = "TestBeta"
	const testGAGate Feature = "TestGA"

	featuremap := map[Feature]Spec{
		testGAGate:    {Default: true},
		testAlphaGate: {Default: false},
		testBetaGate:  {Default: true},
	}

	tests := []struct {
		setmap map[string]bool
		expect string
	}{
		{
			setmap: map[string]bool{
				"TestAlpha": false,
			},
			expect: "TestAlpha=false",
		},
		{
			setmap: map[string]bool{
				"TestAlpha": false,
				"TestBeta":  true,
			},
			expect: "TestAlpha=false,TestBeta=true",
		},
		{
			setmap: map[string]bool{
				"TestGA":    true,
				"TestAlpha": false,
				"TestBeta":  true,
			},
			expect: "TestAlpha=false,TestBeta=true,TestGA=true",
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("SetFromMap %s", test.expect), func(t *testing.T) {
			f := NewGate()
			f.Add(featuremap)
			f.SetFromMap(test.setmap)
			result := f.String()
			if result != test.expect {
				t.Errorf("%d: SetFromMap(%#v) Expected %s, Got %s", i, test.setmap, test.expect, result)
			}
		})
	}
}

func TestFeatureGateFlag(t *testing.T) {
	// gates for testing
	const testAlphaGate Feature = "TestAlpha"
	const testBetaGate Feature = "TestBeta"

	tests := []struct {
		arg        string
		expect     map[Feature]bool
		parseError string
	}{
		{
			arg: "",
			expect: map[Feature]bool{
				testAlphaGate: false,
				testBetaGate:  false,
			},
		},
		{
			arg: "fooBarBaz=maybeidk",
			expect: map[Feature]bool{
				testAlphaGate: false,
				testBetaGate:  false,
			},
			parseError: "unrecognized key: fooBarBaz",
		},
		{
			arg: "TestAlpha=true",
			expect: map[Feature]bool{
				testAlphaGate: true,
				testBetaGate:  false,
			},
		},
		{
			arg: "TestAlpha=true",
			expect: map[Feature]bool{
				testAlphaGate: true,
				testBetaGate:  false,
			},
		},
		{
			arg: "TestAlpha=false",
			expect: map[Feature]bool{
				testAlphaGate: false,
				testBetaGate:  false,
			},
		},
		{
			arg: "TestAlpha=false",
			expect: map[Feature]bool{
				testAlphaGate: false,
				testBetaGate:  false,
			},
		},
		{
			arg: "TestBeta=true",
			expect: map[Feature]bool{
				testAlphaGate: false,
				testBetaGate:  true,
			},
		},
	}
	for i, test := range tests {
		fs := flag.NewFlagSet("testfeaturegateflag", flag.ContinueOnError)
		f := NewGate()
		f.Add(map[Feature]Spec{
			testAlphaGate: {Default: false},
			testBetaGate:  {Default: false},
		})
		f.AddFlag(fs)

		err := fs.Parse([]string{fmt.Sprintf("-%s=%s", flagName, test.arg)})
		if test.parseError != "" {
			if !strings.Contains(err.Error(), test.parseError) {
				t.Errorf("%d: Parse() Expected %v, Got %v", i, test.parseError, err)
			}
		} else if err != nil {
			t.Errorf("%d: Parse() Expected nil, Got %v", i, err)
		}
		for k, v := range test.expect {
			if actual := f.enabled.Load().(map[Feature]bool)[k]; actual != v {
				t.Errorf("%d: expected %s=%v, Got %v", i, k, v, actual)
			}
		}
	}
}
