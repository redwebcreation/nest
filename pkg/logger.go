package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"os"
	"time"
)

var (
	ProxyLogger    *Logger
	InternalLogger *Logger
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type Logger struct {
	Path string
}

type Field struct {
	Key string
	Val interface{}
}

func NewField(key string, val interface{}) *Field {
	return &Field{
		Key: key,
		Val: val,
	}
}

func (l Logger) Log(level Level, message string, fields ...*Field) {
	f, err := os.OpenFile(l.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	event := make([]*Field, len(fields)+3)
	event[0] = NewField("level", level)
	event[1] = NewField("time", time.Now().Format("2006-01-02 15:04:05"))
	event[2] = NewField("message", message)

	for i, field := range fields {
		event[i+3] = field
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(bytes)
	if err != nil {
		panic(err)
	}
}

// Infof logs a formatted message at the info level.
// You should not use it directly, this is for compatibility with logrus.
func (l Logger) Infof(message string, a ...interface{}) {
	l.Log(LevelInfo, fmt.Sprintf(message, a...))
}

// Errorf logs a formatted message at the debug level.
// You should not use it directly, this is for compatibility with logrus.
func (l Logger) Errorf(message string, a ...interface{}) {
	l.Log(LevelError, fmt.Sprintf(message, a...))
}

func (l Logger) Error(err error) {
	l.Log(LevelError, err.Error())
}

func init() {
	ProxyLogger = &Logger{
		Path: global.ProxyLogFile,
	}
	InternalLogger = &Logger{
		Path: global.InternalLogFile,
	}
}
