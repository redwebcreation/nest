package loggy

import (
	"fmt"
	"io"
	"os"
)

type FileLogger struct {
	Path   string
	Writer io.Writer
	Perm   os.FileMode
}

func (f *FileLogger) Write(p []byte) (int, error) {
	if f.Writer == nil && f.Path == "" {
		return 0, fmt.Errorf("no file specified")
	}

	if f.Writer == nil {
		if f.Perm == 0 {
			f.Perm = 0600
		}

		var err error
		f.Writer, err = os.OpenFile(f.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return 0, err
		}

	}

	return f.Writer.Write(p)
}
