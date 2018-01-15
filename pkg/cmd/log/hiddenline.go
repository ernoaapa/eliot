package log

// HiddenLine is Line implementation which don't output anything
type HiddenLine struct {
}

// WithProgress display progress bar when line is in loading state
func (r *HiddenLine) WithProgress(current, total int64) Line {
	return r
}

// Infof mark this line to be just blank info line
func (r *HiddenLine) Infof(format string, args ...interface{}) Line {
	return r
}

// Info mark this line to be just blank info line
func (r *HiddenLine) Info(a ...interface{}) Line {
	return r
}

// Loadingf mark this line to be loading (displays loading indicator)
func (r *HiddenLine) Loadingf(format string, args ...interface{}) Line {
	return r
}

// Loading mark this line to be loading (displays loading indicator)
func (r *HiddenLine) Loading(a ...interface{}) Line {
	return r
}

// Donef marks this line to be done and updates the text
func (r *HiddenLine) Donef(format string, args ...interface{}) Line {
	return r
}

// Done marks this line to be done and updates the text
func (r *HiddenLine) Done(a ...interface{}) Line {
	return r
}

// Warnf mark this line to be in warning with given format
func (r *HiddenLine) Warnf(format string, args ...interface{}) Line {
	return r
}

// Warn mark this line to be in warning with given message
func (r *HiddenLine) Warn(a ...interface{}) Line {
	return r
}

// Errorf mark this line to be in error with given format
func (r *HiddenLine) Errorf(format string, args ...interface{}) Line {
	return r
}

// Error mark this line to be in error with given message
func (r *HiddenLine) Error(a ...interface{}) Line {
	return r
}

// Fatalf mark this line to be in fatal with given format
// Will exit(1) after rerendering the lines
func (r *HiddenLine) Fatalf(format string, args ...interface{}) {
}

// Fatal mark this line to be in error with given message
// Will exit(1) after rerendering the lines
func (r *HiddenLine) Fatal(a ...interface{}) {
}
