package ragtag_test

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dgravesa/go-ragtag"
)

func Mask(v interface{}) error {
	tagExecutor := ragtag.Executor{
		TagKey:  "mask",
		TagFunc: applyMask,
	}

	return tagExecutor.Execute(v)
}

func applyMask(val reflect.Value, tagVal string) error {
	if len(tagVal) > 1 {
		return fmt.Errorf("mask tag value must be a single character")
	}

	lenS := len(val.String())

	val.SetString(strings.Repeat(tagVal, lenS))

	return nil
}

func ExampleRagTag() {
	type Inner struct {
		OtherStr   string
		MyInnerStr string `mask:"X"`
	}
	type Outer struct {
		MyOuterStr string `mask:"-"`
		Inner      *Inner
	}

	s := Outer{
		MyOuterStr: "Hello, World!",
		Inner: &Inner{
			OtherStr:   "this is another string",
			MyInnerStr: "Good Day",
		},
	}

	err := Mask(&s)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(s.MyOuterStr)
	fmt.Println(s.Inner.OtherStr)
	fmt.Println(s.Inner.MyInnerStr)

	// Output:
	// -------------
	// this is another string
	// XXXXXXXX
}
