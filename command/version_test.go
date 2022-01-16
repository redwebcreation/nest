package command

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/util"
)

func TestNewVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()
	oldVersion := global.Version
	expected := strconv.FormatInt(time.Now().UnixMilli(), 10)
	global.Version = expected
	util.Stdout = new(bytes.Buffer)

	err := cmd.Execute()
	if err != nil {
		t.Error(err)
	}

	output, err := io.ReadAll(util.Stdout)
	if err != nil {
		t.Error(err)
	}

	if string(output) != "nest@"+expected+"\n" {
		t.Errorf("Expected %s, got %s", expected, string(output))
	}

	// cleanup
	global.Version = oldVersion
	util.Stdout = os.Stdout
}
