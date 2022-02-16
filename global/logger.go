package global

import (
	"bytes"
	"fmt"
	"github.com/go-logfmt/logfmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

var ProxyLogger *log.Logger

type CompositeLogger struct {
	Loggers []io.Writer
}

func (c CompositeLogger) Write(p []byte) (int, error) {
	nn := 0
	for _, l := range c.Loggers {
		n, err := l.Write(p)
		nn += n
		if err != nil {
			return nn, err
		}
	}

	return nn, nil
}

type FileLogger struct {
	Path string
	File *os.File
	Perm os.FileMode
}

func (f *FileLogger) Write(p []byte) (int, error) {
	if f.File == nil && f.Path == "" {
		return 0, fmt.Errorf("no file specified")
	}

	if f.File == nil {
		if f.Perm == 0 {
			f.Perm = 0600
		}

		var err error
		f.File, err = os.OpenFile(f.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return 0, err
		}
	}

	return f.File.Write(p)
}

type Fields map[string]interface{}

type Level int

var levelToString = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// LogP writes logs for the reverse proxy.
func LogP(level Level, message string, fields Fields) {
	ProxyLogger.Print(newFields(level, message, fields))
}

// LogI logs internal events.
func LogI(level Level, message string, fields Fields) {
	log.Print(newFields(level, message, fields))
}

func newFields(level Level, message string, fields Fields) Fields {
	if fields == nil {
		fields = make(Fields)
	}

	fields["level"] = levelToString[level]
	fields["message"] = message
	fields["time"] = time.Now().Format("2006-01-02 15:04:05")

	return fields
}

func (f Fields) String() string {
	var buf bytes.Buffer
	enc := logfmt.NewEncoder(&buf)

	var keys []string
	for k := range f {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		err := enc.EncodeKeyval(k, f[k])
		if err != nil {
			panic(err)
		}
	}

	err := enc.EndRecord()
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func init() {
	log.SetPrefix("")
	log.SetFlags(0)
	log.SetOutput(&FileLogger{
		Path: GetInternalLogFile(),
	})

	ProxyLogger = log.New(CompositeLogger{
		Loggers: []io.Writer{
			&FileLogger{
				Path: GetProxyLogFile(),
			},
			&FileLogger{
				File: os.Stdout,
			},
		},
	}, "", 0)
}

type FinisherLogger struct {
	Logger *log.Logger
}

func (l FinisherLogger) Infof(message string, args ...any) {
	l.Logger.Print(newFields(LevelInfo, fmt.Sprintf(message, args...), Fields{
		"tag": "proxy.finisher",
	}))
}

func (l FinisherLogger) Errorf(message string, args ...any) {
	l.Logger.Print(newFields(LevelError, fmt.Sprintf(message, args...), Fields{
		"tag": "proxy.finisher",
	}))
}
