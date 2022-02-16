package cli

import (
	"fmt"
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/global"
	"os"
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	cmd := NewVersionCommand()
	os.Stderr = c.Tty()
	os.Stdin = c.Tty()
	os.Stdout = c.Tty()

	go func() {
		_, _ = c.ExpectString(fmt.Sprintf("Nest version %s, build %s\n", global.Version, global.Commit))
	}()

	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}
}
