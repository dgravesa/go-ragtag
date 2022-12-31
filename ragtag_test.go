package ragtag_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dgravesa/go-ragtag"
)

func assertEqual(t *testing.T, expected, actual any) {
	if expected != actual {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func Test_Execute_Basic_Success(t *testing.T) {
	// arrange
	type T struct {
		S string `mask:"-"`
	}
	val := T{
		S: "Hello, World!",
	}

	// act
	err := ragtag.Execute("mask", &val, applyMask)

	// assert
	assertEqual(t, nil, err)
	assertEqual(t, "-------------", val.S)
}

func Test_Execute_Basic_Error(t *testing.T) {
	// arrange
	type T struct {
		S string `mask:"-"`
	}
	val := T{
		S: "Hello, World!",
	}
	expectedErr := fmt.Errorf("error occurred")

	// act
	actualErr := ragtag.Execute("mask", &val, func(_ reflect.Value, _ string) error {
		return expectedErr
	})

	// assert
	assertEqual(t, expectedErr, actualErr)
}

func Test_Execute_WithNestedStruct_Success(t *testing.T) {
	// arrange
	type I struct {
		S1 string
		S2 string `mask:"X"`
		S3 string
	}
	type O struct {
		I I
		S string
	}
	val := O{
		I: I{
			S1: "Hello, World!",
			S2: "password",
			S3: "1234",
		},
		S: "outer",
	}

	// act
	err := ragtag.Execute("mask", &val, applyMask)

	// assert
	assertEqual(t, nil, err)
	assertEqual(t, O{
		I: I{
			S1: "Hello, World!",
			S2: "XXXXXXXX",
			S3: "1234",
		},
		S: "outer",
	}, val)
}

func Test_Execute_WithNestedStruct_Error(t *testing.T) {
	// arrange
	type I struct {
		S1 string
		S2 string `mask:"X"`
		S3 string
	}
	type O struct {
		I I
		S string
	}
	val := O{
		I: I{
			S1: "Hello, World!",
			S2: "password",
			S3: "1234",
		},
		S: "outer",
	}
	expectedErr := fmt.Errorf("error occurred")

	// act
	actualErr := ragtag.Execute("mask", &val, func(_ reflect.Value, _ string) error {
		return expectedErr
	})

	// assert
	assertEqual(t, expectedErr, actualErr)
}

func Test_Execute_WithSlice_Success(t *testing.T) {
	// arrange
	type I struct {
		S string `test:"*"`
	}
	val := []I{
		{
			S: "Hello",
		},
		{
			S: "Hi",
		},
	}

	// act
	err := ragtag.Execute("test", val, applyMask)

	// assert
	assertEqual(t, nil, err)
	assertEqual(t, val[0], I{
		S: "*****",
	})
	assertEqual(t, val[1], I{
		S: "**",
	})
}

func Test_Execute_WithSlice_Error(t *testing.T) {
	// arrange
	type I struct {
		S string `test:"*"`
	}
	val := []I{
		{
			S: "Hello",
		},
		{
			S: "Hi",
		},
	}
	expectedErr := fmt.Errorf("error occurred")

	// act
	actualErr := ragtag.Execute("test", &val, func(_ reflect.Value, _ string) error {
		return expectedErr
	})

	// assert
	assertEqual(t, expectedErr, actualErr)
}

func Test_Execute_WithMap_Success(t *testing.T) {
	// arrange
	type I struct {
		S string `test:"x"`
	}
	val := map[string]*I{
		"this one": {
			S: "Test String",
		},
	}

	// act
	err := ragtag.Execute("test", val, applyMask)

	// assert
	assertEqual(t, nil, err)
	assertEqual(t, I{
		S: "xxxxxxxxxxx",
	}, *val["this one"])
}

func Test_Execute_WithMap_Error(t *testing.T) {
	// arrange
	type I struct {
		S string `test:"x"`
	}
	val := map[string]*I{
		"this one": {
			S: "Test String",
		},
	}
	expectedErr := fmt.Errorf("error occurred")

	// act
	actualErr := ragtag.Execute("test", val, func(_ reflect.Value, _ string) error {
		return expectedErr
	})

	// assert
	assertEqual(t, expectedErr, actualErr)
}
