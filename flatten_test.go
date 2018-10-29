package bellows_test

import (
	"fmt"
	"github.com/mpecan/bellows"
	"reflect"
	"regexp"
	"testing"
)

type testCase struct {
	name     string
	expanded map[string]interface{}
	flat     map[string]interface{}
}

var basics = []testCase{
	{
		name:     "two levels",
		expanded: map[string]interface{}{"some": map[string]interface{}{"other": "element"}},
		flat:     map[string]interface{}{"some.other": "element"},
	},
	{
		name:     "two levels multi path",
		expanded: map[string]interface{}{"some": map[string]interface{}{"other": "element", "another": "element"}},
		flat:     map[string]interface{}{"some.other": "element", "some.another": "element"},
	},
	{
		name:     "three levels multi path",
		expanded: map[string]interface{}{"some": map[string]interface{}{"other": "element", "another": "element", "anotherScope": map[string]interface{}{"deeper": map[string]interface{}{"we": "go"}}}},
		flat:     map[string]interface{}{"some.other": "element", "some.another": "element", "some.anotherScope.deeper.we": "go"},
	},
}

func prefixMap(in map[string]interface{}, prefix string, separator string) map[string]interface{} {
	out := make(map[string]interface{})
	for key, value := range in {
		out[prefix+separator+key] = value
	}
	return out
}

func replaceSeparator(in map[string]interface{}, separator string, t *testing.T) map[string]interface{} {
	defSepRegexp, err := regexp.Compile("[.]")
	if err != nil {
		t.Error(err)
	}
	out := make(map[string]interface{})
	for key, value := range in {
		out[defSepRegexp.ReplaceAllString(key, separator)] = value
	}
	return out
}

func prefix(caseToPrefix testCase, prefix string) testCase {
	expanded := make(map[string]interface{}, 1)
	expanded[prefix] = caseToPrefix.expanded
	return testCase{
		name:     fmt.Sprintf("%s prefixed with %s", caseToPrefix.name, prefix),
		expanded: expanded,
		flat:     prefixMap(caseToPrefix.flat, prefix, "."),
	}
}

func reSeparator(toReSeparate testCase, separator string, t *testing.T) testCase {
	return testCase{
		name:     fmt.Sprintf("%s with separator:%s", toReSeparate.name, separator),
		expanded: toReSeparate.expanded,
		flat:     replaceSeparator(toReSeparate.flat, separator, t),
	}
}

func TestBasicFlatten(t *testing.T) {
	for _, tt := range basics {
		t.Run(fmt.Sprintf("flatten:%+v", tt.name), func(t *testing.T) {
			flattened := bellows.Flatten(tt.expanded)
			if !reflect.DeepEqual(flattened, tt.flat) {
				t.Errorf("got %q, want %q", flattened, tt.flat)
			}
		})
	}
}

func TestBasicExpand(t *testing.T) {
	for _, tt := range basics {
		t.Run(fmt.Sprintf("expand:%+v", tt.name), func(t *testing.T) {
			expanded := bellows.Expand(tt.flat)
			if !reflect.DeepEqual(expanded, tt.expanded) {
				t.Errorf("got %q, want %q", expanded, tt.flat)
			}
		})
	}
}

func TestPrefixedFlatten(t *testing.T) {
	for _, tt := range basics {
		current := prefix(tt, "somePrefix")
		t.Run(fmt.Sprintf("flatten:%+v", current.name), func(t *testing.T) {
			flattened := bellows.Flatten(current.expanded)
			if !reflect.DeepEqual(flattened, current.flat) {
				t.Errorf("got %q, want %q", flattened, current.flat)
			}
		})
	}
}

func TestPrefixedExpand(t *testing.T) {
	for _, tt := range basics {
		current := prefix(tt, "somePrefix")
		t.Run(fmt.Sprintf("expand:%+v", current.name), func(t *testing.T) {
			expanded := bellows.Expand(current.flat)
			if !reflect.DeepEqual(expanded, current.expanded) {
				t.Errorf("got %q, want %q", expanded, current.flat)
			}
		})
	}
}

func TestCustomSeparatorFlatten(t *testing.T) {
	customSeparator := "-"
	for _, tt := range basics {
		current := reSeparator(tt, customSeparator, t)
		t.Run(fmt.Sprintf("flatten:%+v", current.name), func(t *testing.T) {
			flattened := bellows.FlattenWithSeparator(current.expanded, customSeparator)
			if !reflect.DeepEqual(flattened, current.flat) {
				t.Errorf("got %q, want %q", flattened, current.flat)
			}
		})
	}
}

func TestCustomSeparatorExpand(t *testing.T) {
	customSeparator := "-"
	for _, tt := range basics {
		current := reSeparator(tt, customSeparator, t)
		t.Run(fmt.Sprintf("expand:%+v", current.name), func(t *testing.T) {
			expanded := bellows.ExpandWithSeparator(current.flat, customSeparator)
			if !reflect.DeepEqual(expanded, current.expanded) {
				t.Errorf("got %q, want %q", expanded, current.flat)
			}
		})
	}
}

func TestCustomSeparatorWithPrefixFlatten(t *testing.T) {
	customSeparator := "-"
	customPrefix := "prefix"
	for _, tt := range basics {
		current := reSeparator(prefix(tt, customPrefix), customSeparator, t)
		t.Run(fmt.Sprintf("flatten:%+v", current.name), func(t *testing.T) {
			flattened := bellows.FlattenWithSeparator(current.expanded, customSeparator)
			if !reflect.DeepEqual(flattened, current.flat) {
				t.Errorf("got %q, want %q", flattened, current.flat)
			}
		})
	}
}

func TestCustomSeparatorWithPrefixExpand(t *testing.T) {
	customSeparator := "-"
	customPrefix := "prefix"
	for _, tt := range basics {
		current := reSeparator(prefix(tt, customPrefix), customSeparator, t)
		t.Run(fmt.Sprintf("expand:%+v", current.name), func(t *testing.T) {
			expanded := bellows.ExpandWithSeparator(current.flat, customSeparator)
			if !reflect.DeepEqual(expanded, current.expanded) {
				t.Errorf("got %q, want %q", expanded, current.flat)
			}
		})
	}
}
