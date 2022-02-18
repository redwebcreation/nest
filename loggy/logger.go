// Package loggy contains a simple logger interoperable with the standard library.
// It is name loggy simply to avoid naming conflicts.
package loggy

import (
	"bytes"
	"github.com/go-logfmt/logfmt"
	"io"
	"log"
	"sort"
	"time"
)

type Fields map[string]interface{}

func NewEvent(level Level, message string, fields Fields) Fields {
	if fields == nil {
		fields = make(Fields)
	}

	fields["level"] = level.String()
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

func NewNullLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}
