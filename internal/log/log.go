package log

import (
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
)

var defaultLog Logger = &stdoutLogger{
	level: logrus.InfoLevel,
}

// Discard is a logger implementation that just discards every log statement
var Discard = &DiscardLogger{}

// Debug prints debug information
func Debug(args ...interface{}) {
	defaultLog.Debug(args...)
}

// Debugf prints formatted debug information
func Debugf(format string, args ...interface{}) {
	defaultLog.Debugf(format, args...)
}

// Info prints info information
func Info(args ...interface{}) {
	defaultLog.Info(args...)
}

// Infof prints formatted information
func Infof(format string, args ...interface{}) {
	defaultLog.Infof(format, args...)
}

// Warn prints warning information
func Warn(args ...interface{}) {
	defaultLog.Warn(args...)
}

// Warnf prints formatted warning information
func Warnf(format string, args ...interface{}) {
	defaultLog.Warnf(format, args...)
}

// Error prints error information
func Error(args ...interface{}) {
	defaultLog.Error(args...)
}

// Errorf prints formatted error information
func Errorf(format string, args ...interface{}) {
	defaultLog.Errorf(format, args...)
}

// Fatal prints fatal error information
func Fatal(args ...interface{}) {
	defaultLog.Fatal(args...)
}

// Fatalf prints formatted fatal error information
func Fatalf(format string, args ...interface{}) {
	defaultLog.Fatalf(format, args...)
}

// Panic prints panic information
func Panic(args ...interface{}) {
	defaultLog.Panic(args...)
}

// Panicf prints formatted panic information
func Panicf(format string, args ...interface{}) {
	defaultLog.Panicf(format, args...)
}

// Done prints done information
func Done(args ...interface{}) {
	defaultLog.Done(args...)
}

// Donef prints formatted info information
func Donef(format string, args ...interface{}) {
	defaultLog.Donef(format, args...)
}

// Fail prints error information
func Fail(args ...interface{}) {
	defaultLog.Fail(args...)
}

// Failf prints formatted error information
func Failf(format string, args ...interface{}) {
	defaultLog.Failf(format, args...)
}

// Print prints information
func Print(level logrus.Level, args ...interface{}) {
	defaultLog.Print(level, args...)
}

// Printf prints formatted information
func Printf(level logrus.Level, format string, args ...interface{}) {
	defaultLog.Printf(level, format, args...)
}

// SetLevel changes the log level of the global logger
func SetLevel(level logrus.Level) {
	defaultLog.SetLevel(level)
}

// StartFileLogging logs the output of the global logger to the file default.log
func StartFileLogging() {
	defaultLogStdout, ok := defaultLog.(*stdoutLogger)
	if ok {
		defaultLogStdout.fileLogger = GetFileLogger("default")
	}
}

// GetInstance returns the Logger instance
func GetInstance() Logger {
	return defaultLog
}

// SetInstance sets the default logger instance
func SetInstance(logger Logger) {
	defaultLog = logger
}

// WriteColored writes a message in color
func WriteColored(message string, color string) {
	_, _ = defaultLog.Write([]byte(ansi.Color(message, color)))
}

// Write writes to the stdout log without formatting the message, but takes care of locking the log and halting a possible wait message
func Write(message []byte) {
	_, _ = defaultLog.Write(message)
}

// WriteString writes to the stdout log without formatting the message, but takes care of locking the log and halting a possible wait message
func WriteString(message string) {
	defaultLog.WriteString(message)
}
