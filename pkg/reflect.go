package pkg

import (
	"reflect"
)

func IsZeroValue(v reflect.Value) bool {
	kind := v.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		return v.IsNil()
	}
	return v.IsZero()
}
