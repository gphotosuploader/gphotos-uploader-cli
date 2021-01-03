package cmd_test

import "testing"

func assertExpectedError(t *testing.T, errExpected bool, err error, ) {
	if errExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !errExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
