package ragtag

import (
	"reflect"
)

type Func func(val reflect.Value, tag reflect.StructTag) error

func Execute(v interface{}, f Func) error {
	return executePrePost(reflect.ValueOf(v), "", f, nil)
}

func ExecutePrePost(v interface{}, pre, post Func) error {
	return executePrePost(reflect.ValueOf(v), "", pre, post)
}

func executePrePost(val reflect.Value, tag reflect.StructTag, pre, post Func) error {
	if pre != nil {
		if err := pre(val, tag); err != nil {
			return err
		}
	}

	elem := getElement(val)

	switch elem.Kind() {
	case reflect.Struct:
		t := elem.Type()
		for i := 0; i < t.NumField(); i++ {
			// execute recursively
			if err := executePrePost(elem.Field(i), t.Field(i).Tag, pre, post); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < elem.Len(); i++ {
			// execute recursively
			if err := executePrePost(elem.Index(i), "", pre, post); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range elem.MapKeys() {
			// execute recursively
			if err := executePrePost(elem.MapIndex(k), "", pre, post); err != nil {
				return err
			}
		}
	}

	if post != nil {
		if err := post(val, tag); err != nil {
			return err
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
