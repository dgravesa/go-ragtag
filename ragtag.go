package ragtag

import (
	"reflect"
)

type Func func(val reflect.Value, tagVal string) error

func Execute(tagKey string, v any, f Func) error {
	return execute(tagKey, reflect.ValueOf(v), f)
}

func execute(tagKey string, val reflect.Value, f Func) error {
	elem := getElement(val)

	switch elem.Kind() {
	case reflect.Struct:
		t := elem.Type()
		for i := 0; i < t.NumField(); i++ {
			// execute tag func if tag is specified
			if tagVal, tagSpecified := t.Field(i).Tag.Lookup(tagKey); tagSpecified {
				if err := f(elem.Field(i), tagVal); err != nil {
					return err
				}
			}

			// execute recursively
			if err := execute(tagKey, elem.Field(i), f); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < elem.Len(); i++ {
			// execute recursively
			if err := execute(tagKey, elem.Index(i), f); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range elem.MapKeys() {
			// execute recursively
			if err := execute(tagKey, elem.MapIndex(k), f); err != nil {
				return err
			}
		}
	}

	return nil
}

func getElement(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		return val.Elem()
	}
	return val
}
