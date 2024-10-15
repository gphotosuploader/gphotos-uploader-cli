package log

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logdir specifies the relative path to the gphotos-uploader-cli logs
var Logdir = "./.gphotos-uploader-cli/logs/"

var logs = map[string]Logger{}
var logsMutext sync.Mutex

type fileLogger struct {
	logger *logrus.Logger
}

// GetFileLogger returns a logger instance for the specified filename
func GetFileLogger(filename string) Logger {
	logsMutext.Lock()
	defer logsMutext.Unlock()

	log := logs[filename]
	if log == nil {
		newLogger := &fileLogger{
			logger: logrus.New(),
		}
		newLogger.logger.Formatter = &logrus.JSONFormatter{}

		_ = os.MkdirAll(Logdir, os.ModePerm)

		logFile, err := os.OpenFile(Logdir+filename+".log", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			newLogger.Warnf("Unable to open %s log file. Will log to stdout.", filename)
		} else {
			newLogger.logger.SetOutput(logFile)
		}

		logs[filename] = newLogger
	}

	return logs[filename]
}

func (f *fileLogger) Debug(args ...interface{}) {
	f.logger.Debug(args...)
}

func (f *fileLogger) Debugf(format string, args ...interface{}) {
	f.logger.Debugf(format, args...)
}

func (f *fileLogger) Info(args ...interface{}) {
	f.logger.Info(args...)
}

func (f *fileLogger) Infof(format string, args ...interface{}) {
	f.logger.Infof(format, args...)
}

func (f *fileLogger) Warn(args ...interface{}) {
	f.logger.Warn(args...)
}

func (f *fileLogger) Warnf(format string, args ...interface{}) {
	f.logger.Warnf(format, args...)
}

func (f *fileLogger) Error(args ...interface{}) {
	f.logger.Error(args...)
}

func (f *fileLogger) Errorf(format string, args ...interface{}) {
	f.logger.Errorf(format, args...)
}

func (f *fileLogger) Fatal(args ...interface{}) {
	f.logger.Fatal(args...)
}

func (f *fileLogger) Fatalf(format string, args ...interface{}) {
	f.logger.Fatalf(format, args...)
}

func (f *fileLogger) Panic(args ...interface{}) {
	f.logger.Panic(args...)
}

func (f *fileLogger) Panicf(format string, args ...interface{}) {
	f.logger.Panicf(format, args...)
}

func (f *fileLogger) Done(args ...interface{}) {
	f.logger.Info(args...)
}

func (f *fileLogger) Donef(format string, args ...interface{}) {
	f.logger.Infof(format, args...)
}

func (f *fileLogger) Fail(args ...interface{}) {
	f.logger.Error(args...)
}

func (f *fileLogger) Failf(format string, args ...interface{}) {
	f.logger.Errorf(format, args...)
}

func (f *fileLogger) Print(level logrus.Level, args ...interface{}) {
	switch level {
	case logrus.InfoLevel:
		f.Info(args...)
	case logrus.DebugLevel:
		f.Debug(args...)
	case logrus.WarnLevel:
		f.Warn(args...)
	case logrus.ErrorLevel:
		f.Error(args...)
	case logrus.PanicLevel:
		f.Panic(args...)
	case logrus.FatalLevel:
		f.Fatal(args...)
	}
}

func (f *fileLogger) Printf(level logrus.Level, format string, args ...interface{}) {
	switch level {
	case logrus.InfoLevel:
		f.Infof(format, args...)
	case logrus.DebugLevel:
		f.Debugf(format, args...)
	case logrus.WarnLevel:
		f.Warnf(format, args...)
	case logrus.ErrorLevel:
		f.Errorf(format, args...)
	case logrus.PanicLevel:
		f.Panicf(format, args...)
	case logrus.FatalLevel:
		f.Fatalf(format, args...)
	}
}

func (f *fileLogger) SetLevel(level logrus.Level) {
	f.logger.SetLevel(level)
}

func (f *fileLogger) GetLevel() logrus.Level {
	return f.logger.GetLevel()
}

func (f *fileLogger) Write(message []byte) (int, error) {
	return f.logger.Out.Write(message)
}

func (f *fileLogger) WriteString(message string) {
	_, _ = f.logger.Out.Write([]byte(message))
}
