package internal

import (
	"fmt"
	"testing"
)

func AssertEquals(t *testing.T, expected interface{}, actual interface{}) {
	if actual != expected {
		t.Fail()
		fmt.Printf("expected '%s' got %s'", expected, actual)
	}
}

func AssertNotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Fail()
		fmt.Printf("expected '%s' not to be nil", value)
	}
}
