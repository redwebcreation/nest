package util

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrompt(t *testing.T) {
	defer func() {
		Stdout = os.Stdout
		Stdin = os.Stdin
	}()
	if Stdout != os.Stdout {
		t.Errorf("util.Stdout should be os.Stdout")
	}
	if Stdin != os.Stdin {
		t.Errorf("util.Stdin should be os.Stdin")
	}

	Stdout = &bytes.Buffer{}
	Stdin = bytes.NewBufferString("felix\n")
	value := Prompt("Name", "", nil)

	output, _ := io.ReadAll(Stdout)
	if string(output) != "Name: " {
		t.Errorf("Expected 'Name: ', got %s", string(output))
	}
	if value != "felix" {
		t.Errorf("Expected 'felix', got '%s'", value)
	}

	Stdout = &bytes.Buffer{}
	Stdin = bytes.NewBufferString("\n")
	value = Prompt("Name?", "", nil)

	output, _ = io.ReadAll(Stdout)
	if string(output) != "Name?: " {
		t.Errorf("Expected 'Name?: ', got %s", string(output))
	}
	if value != "" {
		t.Errorf("Expected '', got '%s'", value)
	}

	Stdout = &bytes.Buffer{}
	Stdin = bytes.NewBufferString("\n16\n24\n")
	value = Prompt("Secret number?", "", func(input string) bool {
		return strings.HasPrefix(input, "2")
	})

	output, _ = io.ReadAll(Stdout)
	if string(output) != "Secret number?: Secret number?: Secret number?: " {
		t.Errorf("Secret number?: Secret number?: Secret number?: ', got %s", string(output))
	}
	if value != "24" {
		t.Errorf("Expected '24', got '%s'", value)
	}

	Stdout = &bytes.Buffer{}
	Stdin = bytes.NewBufferString("\nhello\n")
	value = Prompt("Secret number?", "world", nil)

	output, _ = io.ReadAll(Stdout)
	if string(output) != "Secret number? [world]: " {
		t.Errorf("Expected 'Secret number? [world]: ', got %s", string(output))
	}
	if value != "world" {
		t.Errorf("Expected default value 'world', got '%s'", value)
	}
}
