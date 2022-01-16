package util

import (
	"os"
	"strconv"
	"time"
)

func TmpFile() (*os.File, error) {
	return os.Create("/tmp/" + strconv.Itoa(int(time.Now().UnixNano())) + ".tmp")
}
