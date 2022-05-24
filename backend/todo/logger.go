package todo

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// Logger represents a leveled logger
type Logger interface {
	Debug(msg string)
	Debugf(format string, v ...interface{})
	Info(msg string)
	Infof(format string, v ...interface{})
	Warn(msg string)
	Warnf(format string, v ...interface{})
	Error(msg string)
	Errorf(format string, v ...interface{})
	E(err error)
	Level() LogLevel

	// Standard library logger methods below.
	SetOutput(w io.Writer)
	Output(calldepth int, s string) error
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Flags() int
	SetFlags(flag int)
	Prefix() string
	SetPrefix(prefix string)
	Writer() io.Writer

	// Fatal, Fatalf, Fatalln, Panic, Panicf, and Panicln are not implemented.
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

// LogLevel represents a log level
type LogLevel int

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	default:
		return ""
	}
}

const (
	LogLevelDebug LogLevel = 0
	LogLevelInfo  LogLevel = 3
	LogLevelWarn  LogLevel = 5
	LogLevelError LogLevel = 7
)

// NewLogger creates a new DefaultLogger at LogLevelInfo which writes to ioutil.Discard.
func NewLogger() *DefaultLogger {
	l := &DefaultLogger{
		minLevel: LogLevelInfo,
		Logger:   log.New(ioutil.Discard, "", log.LstdFlags),
	}
	return l
}

// DefaultLogger wraps a log.Logger with a leveled implementation and satisfies the Logger interface.
type DefaultLogger struct {
	minLevel LogLevel
	*log.Logger
}

func (l *DefaultLogger) Debug(msg string) {
	if LogLevelDebug >= l.minLevel {
		l.Logger.Println("[debug]", msg)
	}
}

func (l *DefaultLogger) Debugf(format string, v ...interface{}) { l.Debug(fmt.Sprintf(format, v...)) }

func (l *DefaultLogger) Info(msg string) {
	if LogLevelInfo >= l.minLevel {
		l.Logger.Println("[info]", msg)
	}
}

func (l *DefaultLogger) Infof(format string, v ...interface{}) { l.Info(fmt.Sprintf(format, v...)) }

func (l *DefaultLogger) Warn(msg string) {
	if LogLevelWarn >= l.minLevel {
		l.Logger.Println("[warn]", msg)
	}
}

func (l *DefaultLogger) Warnf(format string, v ...interface{}) { l.Warn(fmt.Sprintf(format, v...)) }

func (l *DefaultLogger) Error(msg string) {
	if LogLevelError >= l.minLevel {
		l.Logger.Println("[error]", msg)
	}
}

func (l *DefaultLogger) Errorf(format string, v ...interface{}) { l.Error(fmt.Sprintf(format, v...)) }

func (l *DefaultLogger) E(err error) { l.Error(err.Error()) }

func (l *DefaultLogger) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		l.minLevel = LogLevelDebug
	case "warn":
		l.minLevel = LogLevelWarn
	case "info":
		l.minLevel = LogLevelInfo
	case "error":
		l.minLevel = LogLevelError
	}
}

func (l *DefaultLogger) Level() LogLevel { return l.minLevel }

func (l *DefaultLogger) SetOutput(w io.Writer) { l.Logger.SetOutput(w) }

func (l *DefaultLogger) Output(calldepth int, s string) error { return l.Logger.Output(calldepth+2, s) }

func toFmtString(v ...interface{}) string {
	fmtstr := strings.Repeat("%v", len(v))
	fmtstr = strings.TrimSpace(fmtstr)
	return "[info] " + fmtstr
}

func (l *DefaultLogger) Printf(format string, v ...interface{}) { l.Print(fmt.Sprintf(format, v...)) }

func (l *DefaultLogger) Print(v ...interface{}) {
	if LogLevelInfo >= l.minLevel {
		l.Logger.Print(fmt.Sprintf(toFmtString(v...), v...))
	}
}

func (l *DefaultLogger) Println(v ...interface{}) {
	if LogLevelInfo >= l.minLevel {
		l.Logger.Println(fmt.Sprintf(toFmtString(v...), v...))
	}
}

func (l *DefaultLogger) Fatal(v ...interface{}) {}

func (l *DefaultLogger) Fatalf(format string, v ...interface{}) {}

func (l *DefaultLogger) Fatalln(v ...interface{}) {}

func (l *DefaultLogger) Panic(v ...interface{}) {}

func (l *DefaultLogger) Panicf(format string, v ...interface{}) {}

func (l *DefaultLogger) Panicln(v ...interface{}) {}

func (l *DefaultLogger) Flags() int { return l.Logger.Flags() }

func (l *DefaultLogger) SetFlags(flag int) { l.Logger.SetFlags(flag) }

func (l *DefaultLogger) Prefix() string { return l.Logger.Prefix() }

func (l *DefaultLogger) SetPrefix(prefix string) { l.Logger.SetPrefix(prefix) }

func (l *DefaultLogger) Writer() io.Writer { return l.Logger.Writer() }
