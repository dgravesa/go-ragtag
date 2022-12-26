package ragtag

import (
	"reflect"
)

type TagFunc func(val reflect.Value, tagVal string) error

type Executor struct {
	TagKey string

	TagFunc TagFunc
}

func (e Executor) Execute(v interface{}) error {
	return e.execute(reflect.ValueOf(v))
}

func getElement(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		return val.Elem()
	}
	return val
}

func (e Executor) execute(val reflect.Value) error {
	val = getElement(val)

	switch val.Kind() {
	case reflect.Struct:
		t := val.Type()
		for i := 0; i < t.NumField(); i++ {
			if tagVal, tagSpecified := t.Field(i).Tag.Lookup(e.TagKey); tagSpecified {
				// apply custom tag operation to field
				if err := e.TagFunc(val.Field(i), tagVal); err != nil {
					return err
				}
			} else {
				// execute recursively
				if err := e.execute(val.Field(i)); err != nil {
					return err
				}
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			// execute recursively
			if err := e.execute(val.Index(i)); err != nil {
				return err
			}
		}
	}

	return nil
}
