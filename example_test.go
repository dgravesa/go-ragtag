package ragtag_test

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dgravesa/go-ragtag"
)

func Mask(v interface{}) error {
	return ragtag.Execute(v, applyMask)
}

func applyMask(val reflect.Value, tag reflect.StructTag) error {
	maskTag, hasMaskTag := tag.Lookup("mask")

	if hasMaskTag {
		val.SetString(strings.Repeat(maskTag, val.Len()))
	}

	return nil
}

func ExampleExecute() {
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
