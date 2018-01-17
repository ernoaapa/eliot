package ui

import (
	log "github.com/sirupsen/logrus"
)

// DebugLine is Line implementation which just outputs all the lines to the terminal
type DebugLine struct {
}

// WithProgress display progress bar when line is in loading state
func (r *DebugLine) WithProgress(current, total int64) Line {
	log.Debugf("%d/%d", current, total)
	return r
}

// Infof mark this line to be just blank info line
func (r *DebugLine) Infof(format string, args ...interface{}) Line {
	log.Infof(format, args...)
	return r
}

// Info mark this line to be just blank info line
func (r *DebugLine) Info(a ...interface{}) Line {
	log.Info(a...)
	return r
}

// Loadingf mark this line to be loading (displays loading indicator)
func (r *DebugLine) Loadingf(format string, args ...interface{}) Line {
	log.Debugf(format, args...)
	return r
}

// Loading mark this line to be loading (displays loading indicator)
func (r *DebugLine) Loading(a ...interface{}) Line {
	log.Debug(a...)
	return r
}

// Donef marks this line to be done and updates the text
func (r *DebugLine) Donef(format string, args ...interface{}) Line {
	log.Infof(format, args...)
	return r
}

// Done marks this line to be done and updates the text
func (r *DebugLine) Done(a ...interface{}) Line {
	log.Info(a...)
	return r
}

// Warnf mark this line to be in warning with given format
func (r *DebugLine) Warnf(format string, args ...interface{}) Line {
	log.Warnf(format, args...)
	return r
}

// Warn mark this line to be in warning with given message
func (r *DebugLine) Warn(a ...interface{}) Line {
	log.Warn(a...)
	return r
}

// Errorf mark this line to be in error with given format
func (r *DebugLine) Errorf(format string, args ...interface{}) Line {
	log.Errorf(format, args...)
	return r
}

// Error mark this line to be in error with given message
func (r *DebugLine) Error(a ...interface{}) Line {
	log.Error(a...)
	return r
}

// Fatalf mark this line to be in fatal with given format
// Will exit(1) after rerendering the lines
func (r *DebugLine) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Fatal mark this line to be in error with given message
// Will exit(1) after rerendering the lines
func (r *DebugLine) Fatal(a ...interface{}) {
	log.Fatal(a...)
}
