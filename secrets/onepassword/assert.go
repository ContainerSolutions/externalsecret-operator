package onepassword

import (
	"fmt"
	"testing"
)

func assertEquals(t *testing.T, expected interface{}, actual interface{}) {
	if actual != expected {
		t.Fail()
		fmt.Printf("expected '%s' got %s'", expected, actual)
	}
}

func assertNotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Fail()
		fmt.Printf("expected '%s' not to be nil", value)
	}
}
