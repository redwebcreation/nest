package util

import (
	"fmt"
	"os"
)

func Fatal(message string) {
	_, _ = fmt.Fprintln(os.Stderr, message)

	os.Exit(1)
}
