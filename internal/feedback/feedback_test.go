package feedback

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOutputSelection(t *testing.T) {
	reset()

	myErr := new(bytes.Buffer)
	myOut := new(bytes.Buffer)
	SetOut(myOut)
	SetErr(myErr)

	Print("Hello Foo!")
	require.Equal(t, myOut.String(), "Hello Foo!\n")

	Warning("Hello Bar!")
	require.Equal(t, myErr.String(), "Hello Bar!\n")
}
