package blocker

type Logger interface {
	// Copied from the list of functions provided by https://pkg.go.dev/github.com/coredns/coredns/plugin/pkg/log
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
}

type nopLogger struct{}

func (nopLogger) Error(v ...any)                   {}
func (nopLogger) Errorf(format string, v ...any)   {}
func (nopLogger) Warning(v ...any)                 {}
func (nopLogger) Warningf(format string, v ...any) {}
func (nopLogger) Info(v ...any)                    {}
func (nopLogger) Infof(format string, v ...any)    {}
func (nopLogger) Debug(v ...any)                   {}
func (nopLogger) Debugf(format string, v ...any)   {}
func (nopLogger) Fatal(v ...any)                   {}
func (nopLogger) Fatalf(format string, v ...any)   {}
