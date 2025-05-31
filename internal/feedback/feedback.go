package feedback

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var (
	stdOut      io.Writer
	stdErr      io.Writer
	feedbackOut io.Writer
	feedbackErr io.Writer
	bufferOut   *bytes.Buffer
	bufferErr   *bytes.Buffer
)

func init() {
	reset()
}

// reset resets the feedback package to its initial state, useful for unit testing.
func reset() {
	stdOut = os.Stdout
	stdErr = os.Stderr
	feedbackOut = os.Stdout
	feedbackErr = os.Stderr
	bufferOut = &bytes.Buffer{}
	bufferErr = &bytes.Buffer{}
}

// Result is anything more complex than a sentence that needs to be printed
// for the user.
type Result interface {
	fmt.Stringer
	Data() interface{}
}

// ErrorResult is a result embedding also an error. The error will be printed
// on stderr.
type ErrorResult interface {
	Result
	ErrorString() string
}

// SetOut can be used to change the out writer at runtime.
func SetOut(out io.Writer) {
	stdOut = out
	feedbackOut = io.MultiWriter(bufferOut, stdOut)
}

// SetErr can be used to change the err writer at runtime.
func SetErr(err io.Writer) {
	stdErr = err
	feedbackErr = io.MultiWriter(bufferErr, stdErr)
}

// Printf behaves like fmt.Printf but writes on the out writer and adds a newline.
func Printf(format string, v ...interface{}) {
	Print(fmt.Sprintf(format, v...))
}

// Print behaves like fmt.Print but writes on the out writer and adds a newline.
func Print(v string) {
	fmt.Fprintln(feedbackOut, v) //nolint:errcheck
}

// Warning outputs a warning message.
func Warning(msg string) {
	fmt.Fprintln(feedbackErr, msg) //nolint:errcheck
	logrus.Warning(msg)
}

// FatalError outputs the error and exits with status exitCode.
func FatalError(err error, exitCode ExitCode) {
	Fatal(err.Error(), exitCode)
}

// FatalResult outputs the result and exits with status exitCode.
func FatalResult(res ErrorResult, exitCode ExitCode) {
	PrintResult(res)
	os.Exit(int(exitCode))
}

// Fatal outputs the errorMsg and exits with status exitCode.
func Fatal(errorMsg string, exitCode ExitCode) {
	fmt.Fprintln(stdErr, errorMsg) //nolint:errcheck
	os.Exit(int(exitCode))
}

// PrintResult is a convenient wrapper to provide feedback for complex data.
func PrintResult(res Result) {
	var data string
	var dataErr string
	data = res.String()
	if resErr, ok := res.(ErrorResult); ok {
		dataErr = resErr.ErrorString()
	}
	if data != "" {
		fmt.Fprintln(stdOut, data) //nolint:errcheck
	}
	if dataErr != "" {
		fmt.Fprintln(stdErr, dataErr) //nolint:errcheck
	}
}
