package command

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/global"
	"os"
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	cmd := NewVersionCommand()
	os.Stderr = c.Tty()
	os.Stdin = c.Tty()
	os.Stdout = c.Tty()

	if err = cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	_, err = c.ExpectString("nest@" + global.Version)
	if err != nil {
		t.Fatal(err)
	}
}
