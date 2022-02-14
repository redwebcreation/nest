package global

import (
	"errors"
	"github.com/go-logfmt/logfmt"
	"io"
	"os"
	"time"
)

var ProxyLogger Logger
var InternalLogger Logger

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelMap = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

type Fields map[string]any

type Logger interface {
	Log(level Level, message string, fields Fields)
	// Error logs an error using Log()
	Error(error)
}

type FileLogger struct {
	Path string
	File *os.File
}

func (f FileLogger) Log(level Level, message string, fields Fields) {
	if f.File == nil {
		file, err := os.OpenFile(f.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if errors.Is(err, os.ErrNotExist) {
			return
		}

		if err != nil {
			panic(err)
		}

		f.File = file
	}

	write(f.File, level, message, fields)
}

func write(w io.Writer, level Level, message string, fields Fields) {
	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	e := logfmt.NewEncoder(w)
	check(e.EncodeKeyval("level", levelMap[level]))
	check(e.EncodeKeyval("message", message))
	check(e.EncodeKeyval("time", time.Now().Format("2006-01-02 15:04:05")))

	for k, v := range fields {
		check(e.EncodeKeyval(k, v))
	}

	check(e.EndRecord())
}

type CompositeLogger struct {
	Loggers []Logger
}

func (c CompositeLogger) Log(level Level, message string, fields Fields) {
	for _, logger := range c.Loggers {
		logger.Log(level, message, fields)
	}
}

func (c CompositeLogger) Error(err error) {
	for _, logger := range c.Loggers {
		logger.Error(err)
	}
}

func init() {
	ProxyLogger = CompositeLogger{
		Loggers: []Logger{
			FileLogger{
				Path: GetProxyLogFile(),
			},
			FileLogger{
				File: os.Stdout,
			},
		},
	}

	InternalLogger = FileLogger{
		Path: GetInternalLogFile(),
	}
}

type LogrusCompat struct {
	Logger Logger
}

func (l LogrusCompat) Infof(message string, args ...any) {
	l.Logger.Log(LevelInfo, message, Fields{})
}

func (l LogrusCompat) Errorf(message string, args ...any) {
	l.Logger.Log(LevelError, message, Fields{})
}

func (f FileLogger) Error(err error) {
	f.Log(LevelError, err.Error(), Fields{})
}
