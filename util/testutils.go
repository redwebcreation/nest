package util

import (
	"os"
	"strconv"
	"time"
)

func TmpFile() *os.File {
	f, err := os.Create("/tmp/" + strconv.Itoa(int(time.Now().UnixNano())) + ".tmp")
	if err != nil {
		panic(err)
	}

	return f
}
