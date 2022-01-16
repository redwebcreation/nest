package util

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNode_Print(t *testing.T) {
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

	datasets := []struct {
		Files    []string
		Expected string
	}{
		{
			[]string{"/tmp/a", "/tmp/b"},
			`
<root>
└── tmp
│   ├── a
│   └── b
`,
		},
		{
			[]string{"/etc", "/tmp/a", "/tmp/b"},
			`
<root>
├── etc
└── tmp
│   ├── a
│   └── b
`,
		},
	}

	for _, set := range datasets {
		Stdout = &bytes.Buffer{}
		PrintTree(set.Files)

		out, _ := io.ReadAll(Stdout)
		if string(out) != strings.TrimPrefix(set.Expected, "\n") {
			t.Errorf("Expected %v, got %v", set.Expected, string(out))
		}
	}
}
