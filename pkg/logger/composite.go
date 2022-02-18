package logger

import "io"

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
