// Copyright © 2016 Charles Phillips <charles@doublerebel.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found expanded the LICENSE file.

package bellows

import (
	"reflect"
	"strings"
)

func Expand(value map[string]interface{}) map[string]interface{} {
	return ExpandPrefixed(value, "")
}

func ExpandWithSeparator(value map[string]interface{}, separator string) map[string]interface{} {
	return ExpandPrefixedWithSeparator(value, "", separator)
}

func ExpandPrefixed(value map[string]interface{}, prefix string) map[string]interface{} {
	return ExpandPrefixedWithSeparator(value, prefix, ".")
}

func ExpandPrefixedWithSeparator(value map[string]interface{}, prefix string, separator string) map[string]interface{} {
	m := make(map[string]interface{})
	ExpandPrefixedToResult(value, prefix, m, separator)
	return m
}

func ExpandPrefixedToResult(value map[string]interface{}, prefix string, result map[string]interface{}, separator string) {
	if prefix != "" {
		prefix += separator
	}
	for k, val := range value {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		key := k[len(prefix):]
		idx := strings.Index(key, separator)
		if idx != -1 {
			key = key[:idx]
		}
		if _, ok := result[key]; ok {
			continue
		}
		if idx == -1 {
			result[key] = val
			continue
		}

		// It contains a separator, so it is a more complex structure
		result[key] = ExpandPrefixedWithSeparator(value, k[:len(prefix)+len(key)], separator)
	}
}

func Flatten(value interface{}) map[string]interface{} {
	return FlattenPrefixed(value, "")
}

func FlattenWithSeparator(value map[string]interface{}, separator string) map[string]interface{} {
	return FlattenPrefixedWithSeparator(value, "", separator)
}

func FlattenPrefixed(value interface{}, prefix string) map[string]interface{} {
	return FlattenPrefixedWithSeparator(value, prefix, ".")
}

func FlattenPrefixedWithSeparator(value interface{}, prefix string, separator string) map[string]interface{} {
	m := make(map[string]interface{})
	FlattenPrefixedWithSeparatorToResult(value, prefix, m, separator)
	return m
}

func FlattenPrefixedWithSeparatorToResult(value interface{}, prefix string, m map[string]interface{}, separator string) {
	base := ""
	if prefix != "" {
		base = prefix + separator
	}

	original := reflect.ValueOf(value)
	kind := original.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		original = reflect.Indirect(original)
		kind = original.Kind()
	}
	t := original.Type()

	switch kind {
	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			break
		}
		for _, childKey := range original.MapKeys() {
			childValue := original.MapIndex(childKey)
			FlattenPrefixedWithSeparatorToResult(childValue.Interface(), base+childKey.String(), m, separator)
		}
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			childValue := original.Field(i)
			childKey := t.Field(i).Name
			FlattenPrefixedWithSeparatorToResult(childValue.Interface(), base+childKey, m, separator)
		}
	default:
		if prefix != "" {
			m[prefix] = value
		}
	}
}
