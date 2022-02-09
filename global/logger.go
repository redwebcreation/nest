package global

import (
	"encoding/json"
	"fmt"
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
	Path   string
	Stdout bool
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

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	event := make(map[string]interface{}, len(fields)+3)
	event["level"] = level
	event["time"] = time.Now().Format("2006-01-02 15:04:05")
	event["message"] = message

	for _, field := range fields {
		event[field.Key] = field.Val
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	if l.Stdout {
		fmt.Printf("%s\n", bytes)
	}

	_, err = f.Write(bytes)
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte("\n"))
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
		Path:   ProxyLogFile,
		Stdout: true,
	}
	InternalLogger = &Logger{
		Path:   InternalLogFile,
		Stdout: false,
	}
}
