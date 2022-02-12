package global

import (
	"errors"
	"github.com/go-logfmt/logfmt"
	"io"
	"net/http"
	"net/url"
	"os"
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

type Fields map[string]interface{}

type Logger interface {
	Log(level Level, message string, fields Fields)
	// Error logs an error using Log()
	Error(error)
}

type FileLogger struct {
	Path string
	file *os.File
}

func (f FileLogger) Log(level Level, message string, fields Fields) {
	if f.file == nil {
		file, err := os.OpenFile(f.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if errors.Is(err, os.ErrNotExist) {
			return
		}

		if err != nil {
			panic(err)
		}

		f.file = file
	}

	format(f.file, level, message, fields)
}

func format(w io.Writer, level Level, message string, fields Fields) {
	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	e := logfmt.NewEncoder(w)
	check(e.EncodeKeyval("level", levelMap[level]))
	check(e.EncodeKeyval("message", message))

	for k, v := range fields {
		check(e.EncodeKeyval(k, v))
	}

	check(e.EndRecord())
}

type HTTPLogger struct {
	URL    *url.URL
	Method string
	Client *http.Client
	Body   io.Reader
	Modify func(r *http.Request, level Level, message string, fields Fields)
}

func (h *HTTPLogger) Log(level Level, message string, fields Fields) {
	request, err := http.NewRequest(h.Method, h.URL.String(), h.Body)
	if err != nil {
		panic(err)
	}

	if h.Modify != nil {
		h.Modify(request, level, message, fields)
	}

	if h.Client == nil {
		h.Client = http.DefaultClient
	}

	// send request
	_, err = h.Client.Do(request)
	if err != nil {
		panic(err)
	}
}

type CompositeLogger struct {
	Loggers []Logger
}

func (c CompositeLogger) Log(level Level, message string, fields Fields) {
	for _, logger := range c.Loggers {
		logger.Log(level, message, fields)
	}
}

func init() {
	ProxyLogger = FileLogger{
		Path: ProxyLogFile,
	}

	InternalLogger = FileLogger{
		Path: InternalLogFile,
	}
}

type LogrusCompat struct {
	Logger Logger
}

func (l LogrusCompat) Infof(message string, args ...interface{}) {
	l.Logger.Log(LevelInfo, message, Fields{})
}

func (l LogrusCompat) Errorf(message string, args ...interface{}) {
	l.Logger.Log(LevelError, message, Fields{})
}

func (f FileLogger) Error(err error) {
	f.Log(LevelError, err.Error(), Fields{})
}

func (c CompositeLogger) Error(err error) {
	c.Log(LevelError, err.Error(), Fields{})
}

func (h HTTPLogger) Error(err error) {
	h.Log(LevelError, err.Error(), Fields{})
}
