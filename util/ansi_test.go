package util

import (
	"os"
	"testing"
)

func TestCheckAnsi(t *testing.T) {
	testUsingEnv(t, "GOOS", "windows", func() {
		if checkAnsi(nil) {
			t.Error("ANSI should be disabled on windows")
		}
	})

	testUsingEnv(t, "TERM", "dumb", func() {
		if checkAnsi(nil) {
			t.Error("ANSI should be disabled on dumb terminals")
		}
	})

	testUsingEnv(t, "GOOS", "linux", func() {
		testUsingEnv(t, "TERM", "xterm", func() {
			if !checkAnsi(nil) {
				t.Error("ANSI should be enabled on std terminals")
			}
		})
	})

	if checkAnsi([]string{"--no-ansi"}) {
		t.Error("--no-ansi should disable ANSI")
	}

	if !checkAnsi([]string{"--any-arg"}) {
		t.Error("--any-arg should not disable ANSI")
	}
}

func testUsingEnv(t *testing.T, key string, value string, handler func()) {
	previous := os.Getenv(key)
	defer func() {
		err := os.Setenv(key, previous)
		if err != nil {
			t.Fatalf("Environment compromised: %s", err)
		}
	}()

	err := os.Setenv(key, value)
	if err != nil {
		t.Fatal(err)
	}

	handler()
}

func TestColor_String(t *testing.T) {
	color := Color{100, 100, 100}

	if color.String() != color.Fg() {
		t.Error("Color.String() should return the foreground color")
	}
}

func TestColor_Fg(t *testing.T) {
	color := Color{42, 15, 82}

	if color.Fg() != "\x1b[1m\x1b[38;2;42;15;82m" {
		t.Error("Color.Fg() should return the foreground color with bold text")
	}

	AnsiEnabled = false
	if color.Fg() != "" {
		t.Error("Color.Fg() should return an empty string when ANSI is disabled")
	}
	AnsiEnabled = true
}

func TestColor_Bg(t *testing.T) {
	color := Color{1, 10, 0}

	if color.Bg() != "\u001B[1m\x1b[48;2;1;10;0m" {
		t.Error("Color.Bg() should return the background color with bold text")
	}

	AnsiEnabled = false
	if color.Bg() != "" {
		t.Error("Color.Bg() should return an empty string when ANSI is disabled")
	}
	AnsiEnabled = true
}

func TestReset(t *testing.T) {
	if Reset() != "\x1b[0m" {
		t.Error("Reset() should return the reset ANSI code")
	}

	AnsiEnabled = false
	if Reset() != "" {
		t.Error("Reset() should return an empty string when ANSI is disabled")
	}
	AnsiEnabled = true
}
