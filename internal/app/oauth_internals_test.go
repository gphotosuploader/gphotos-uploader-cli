package app

import (
	"bytes"
	"testing"
)

func TestAskForAuthCodeInTerminal(t *testing.T)  {
	testCases := []struct{
		name string
		input string
		isErrExpected bool
	} {
		{"Should success", "foo", false},
		{"Should fail if code is empty", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stdin := bytes.NewBufferString(tc.input)
			_, err := askForAuthCodeInTerminal(stdin, "")
			assertExpectedError(t, tc.isErrExpected, err)
		})
	}

}

func assertExpectedError(t *testing.T, errExpected bool, err error, ) {
	if errExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !errExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
