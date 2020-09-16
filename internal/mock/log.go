package mock

import "github.com/sirupsen/logrus"

type Logger struct {
	DebugInvoked  bool
	DebugfInvoked bool

	InfoInvoked  bool
	InfofInvoked bool

	WarnInvoked  bool
	WarnfInvoked bool

	ErrorInvoked  bool
	ErrorfInvoked bool

	FatalInvoked  bool
	FatalfInvoked bool

	PanicInvoked  bool
	PanicfInvoked bool

	DoneInvoked  bool
	DonefInvoked bool

	FailInvoked  bool
	FailfInvoked bool

	PrintInvoked  bool
	PrintfInvoked bool

	WriteFn            func(message []byte) (int, error)
	WriteFnInvoked     bool
	WriteStringInvoked bool

	SetLevelInvoked bool
	GetLevelFn      func() logrus.Level
	GetLevelInvoked bool
}

// Debug implements logger interface
func (l *Logger) Debug(args ...interface{}) {
	l.DebugInvoked = true
}

// Debugf implements logger interface
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.DebugfInvoked = true
}

// Info implements logger interface
func (l *Logger) Info(args ...interface{}) {
	l.InfoInvoked = true
}

// Infof implements logger interface
func (l *Logger) Infof(format string, args ...interface{}) {
	l.InfofInvoked = true
}

// Warn implements logger interface
func (l *Logger) Warn(args ...interface{}) {
	l.WarnInvoked = true
}

// Warnf implements logger interface
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.WarnfInvoked = true
}

// Error implements logger interface
func (l *Logger) Error(args ...interface{}) {
	l.ErrorInvoked = true
}

// Errorf implements logger interface
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.ErrorfInvoked = true
}

// Fatal implements logger interface
func (l *Logger) Fatal(args ...interface{}) {
	l.FatalInvoked = true
}

// Fatalf implements logger interface
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.FatalfInvoked = true
}

// Panic implements logger interface
func (l *Logger) Panic(args ...interface{}) {
	l.PanicInvoked = true
}

// Panicf implements logger interface
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.PanicfInvoked = true
}

// Done implements logger interface
func (l *Logger) Done(args ...interface{}) {
	l.DoneInvoked = true
}

// Donef implements logger interface
func (l *Logger) Donef(format string, args ...interface{}) {
	l.DonefInvoked = true
}

// Fail implements logger interface
func (l *Logger) Fail(args ...interface{}) {
	l.FailInvoked = true
}

// Failf implements logger interface
func (l *Logger) Failf(format string, args ...interface{}) {
	l.FailfInvoked = true
}

// Print implements logger interface
func (l *Logger) Print(level logrus.Level, args ...interface{}) {
	l.PrintInvoked = true
}

// Printf implements logger interface
func (l *Logger) Printf(level logrus.Level, format string, args ...interface{}) {
	l.PrintfInvoked = true
}

// SetLevel implements logger interface
func (l *Logger) SetLevel(level logrus.Level) {
	l.SetLevelInvoked = true
}

// GetLevel implements logger interface
func (l *Logger) GetLevel() logrus.Level {
	l.GetLevelInvoked = true
	return l.GetLevelFn()
}

// Write implements logger interface
func (l *Logger) Write(message []byte) (int, error) {
	l.WriteFnInvoked = true
	return l.WriteFn(message)
}

// WriteString implements logger interface
func (l *Logger) WriteString(message string) {
	l.WriteStringInvoked = true
}
