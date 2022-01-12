package util

import (
	"os"
	"strconv"
	"time"
)

func TmpFile() *os.File {
	f, err := os.Create("/tmp/" + TmpName())
	if err != nil {
		panic(err)
	}

	return f
}

func TmpName() string {
	return strconv.Itoa(int(time.Now().UnixNano())) + ".tmp"
}
