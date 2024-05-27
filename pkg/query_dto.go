package pkg

import (
	"fmt"
	"reflect"
	"strings"
)

type SortKind string

const (
	SortKind_Desc SortKind = "desc"
	SortKind_Asc  SortKind = "asc"
)

func SortValidate(structParams any) (err error) {
	SortTraversalParams(structParams, func(key string, value SortKind) {
		if err != nil {
			return
		}

		value = SortKind(strings.ToLower(string(value)))
		if value == SortKind_Asc || value == SortKind_Desc {
			return
		}
		err = fmt.Errorf("%v is invalid value, want 'desc' or 'asc'", key)
	})
	return err
}

func SortTraversalParams(structParams any, action func(key string, value SortKind)) {
	v := reflect.ValueOf(structParams)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i)

		if IsZeroValue(vField) {
			continue
		}

		value, ok := vField.Interface().(SortKind)
		if !ok {
			continue
		}

		tField := t.Field(i)
		tag := tField.Tag.Get("sort")
		key := tag
		if key == "" {
			key = tField.Name
		}

		action(key, value)
	}
}
