package util

import (
	"fmt"
	"os"
)

type Color [3]uint8

var Blue Color = [3]uint8{59, 130, 246}

var Red Color = [3]uint8{239, 68, 68}
var Yellow Color = [3]uint8{245, 158, 11}
var Green Color = [3]uint8{16, 185, 129}

var White Color = [3]uint8{255, 255, 255}
var Gray Color = [3]uint8{147, 148, 153}

var Reset = "\x1b[0m"

var AnsiEnabled = true

func init() {
	if os.Getenv("GOOS") == "windows" {
		AnsiEnabled = false
	} else if os.Getenv("TERM") == "dumb" {
		AnsiEnabled = false
	} else {
		for _, arg := range os.Args {
			if arg == "--no-ansi" {
				AnsiEnabled = false
				return
			}
		}
	}

	if !AnsiEnabled {
		Reset = ""
	}
}

func (c Color) String() string {
	return c.Fg()
}

func (c Color) Fg() string {
	if !AnsiEnabled {
		return ""
	}

	return fmt.Sprintf("\x1b[1m\x1b[38;2;%d;%d;%dm", c[0], c[1], c[2])
}

func (c Color) Bg() string {
	if !AnsiEnabled {
		return ""
	}

	return fmt.Sprintf("\x1b[1m\x1b[48;2;%d;%d;%dm", c[0], c[1], c[2])
}
