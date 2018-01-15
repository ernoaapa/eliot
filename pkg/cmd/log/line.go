package log

type Line interface {
	WithProgress(current, total int64) Line
	Infof(format string, args ...interface{}) Line
	Info(a ...interface{}) Line
	Loadingf(format string, args ...interface{}) Line
	Loading(a ...interface{}) Line
	Donef(format string, args ...interface{}) Line
	Done(a ...interface{}) Line
	Warnf(format string, args ...interface{}) Line
	Warn(a ...interface{}) Line
	Errorf(format string, args ...interface{}) Line
	Error(a ...interface{}) Line
	Fatalf(format string, args ...interface{})
	Fatal(a ...interface{})
}
