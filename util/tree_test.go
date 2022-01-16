package util

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewTree(t *testing.T) {
	dataset := []struct {
		Files    []string
		Expected Node
	}{
		{
			[]string{},
			Node{},
		},
		{
			[]string{"/tmp"},
			Node{
				"tmp": Node{},
			},
		},
		{
			[]string{"/tmp", "/tmp/a", "/tmp/b"},
			Node{
				"tmp": Node{
					"a": Node{},
					"b": Node{},
				},
			},
		},
	}

	for _, set := range dataset {
		tree := NewTree(set.Files)

		if reflect.DeepEqual(tree, set.Expected) == false {
			t.Errorf("Expected %v, got %v", set.Expected, tree)
		}
	}
}

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
		Node     Node
		Expected string
	}{
		{
			Node{
				"tmp": Node{
					"a": Node{},
					"b": Node{},
				},
			},
			`
<root>
└── tmp
│   ├── a
│   └── b
`,
		},
		{
			Node{
				"etc": Node{},
				"tmp": Node{
					"a": Node{},
					"b": Node{},
				},
			},
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
		set.Node.Print(0)

		out, _ := io.ReadAll(Stdout)
		if string(out) != strings.TrimPrefix(set.Expected, "\n") {
			t.Errorf("Expected %v, got %v", set.Expected, string(out))
		}
	}
}
