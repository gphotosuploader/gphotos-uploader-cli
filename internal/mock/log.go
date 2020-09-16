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

// Debug marks the function as invoked.
func (l *Logger) Debug(args ...interface{}) {
	l.DebugInvoked = true
}

// Debugf marks the function as invoked.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.DebugfInvoked = true
}

// Info marks the function as invoked.
func (l *Logger) Info(args ...interface{}) {
	l.InfoInvoked = true
}

// Infof marks the function as invoked.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.InfofInvoked = true
}

// Warn marks the function as invoked.
func (l *Logger) Warn(args ...interface{}) {
	l.WarnInvoked = true
}

// Warnf marks the function as invoked.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.WarnfInvoked = true
}

// Error marks the function as invoked.
func (l *Logger) Error(args ...interface{}) {
	l.ErrorInvoked = true
}

// Errorf marks the function as invoked.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.ErrorfInvoked = true
}

// Fatal marks the function as invoked.
func (l *Logger) Fatal(args ...interface{}) {
	l.FatalInvoked = true
}

// Fatalf marks the function as invoked.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.FatalfInvoked = true
}

// Panic marks the function as invoked.
func (l *Logger) Panic(args ...interface{}) {
	l.PanicInvoked = true
}

// Panicf marks the function as invoked.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.PanicfInvoked = true
}

// Done marks the function as invoked.
func (l *Logger) Done(args ...interface{}) {
	l.DoneInvoked = true
}

// Donef marks the function as invoked.
func (l *Logger) Donef(format string, args ...interface{}) {
	l.DonefInvoked = true
}

// Fail marks the function as invoked.
func (l *Logger) Fail(args ...interface{}) {
	l.FailInvoked = true
}

// Failf marks the function as invoked.
func (l *Logger) Failf(format string, args ...interface{}) {
	l.FailfInvoked = true
}

// Print marks the function as invoked.
func (l *Logger) Print(level logrus.Level, args ...interface{}) {
	l.PrintInvoked = true
}

// Printf marks the function as invoked.
func (l *Logger) Printf(level logrus.Level, format string, args ...interface{}) {
	l.PrintfInvoked = true
}

// SetLevel marks the function as invoked.
func (l *Logger) SetLevel(level logrus.Level) {
	l.SetLevelInvoked = true
}

// GetLevel invokes the mock implementation and marks the function as invoked.
func (l *Logger) GetLevel() logrus.Level {
	l.GetLevelInvoked = true
	return l.GetLevelFn()
}

// Write invokes the mock implementation and marks the function as invoked.
func (l *Logger) Write(message []byte) (int, error) {
	l.WriteFnInvoked = true
	return l.WriteFn(message)
}

// WriteString marks the function as invoked.
func (l *Logger) WriteString(message string) {
	l.WriteStringInvoked = true
}
