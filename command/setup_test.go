package command

import (
	"bytes"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/redwebcreation/nest/util"
	"os"
	"strconv"
	"testing"
	"time"
)

type Set struct {
	Strategy   string
	Provider   string
	Repository string
	Branch     string
	Dir        string
	Error      error
}

var dataset = []Set{
	{"remote", "github", "felixdorn/config-test", "main", "", nil},
	{"remote", "gitlab", "felixdorn/config-test", "main", "", nil},
	{"remote", "bitbucket", "felixdorn/config-test", "main", "", nil},
	{"invalidStrategy", "github", "felixdorn/config-test", "main", "", pkg.ErrInvalidStrategy},
	{"remote", "invalidProvider", "felixdorn/config-test", "main", "", pkg.ErrInvalidProvider},
	{"remote", "github", "invalidRepository", "main", "", pkg.ErrInvalidRepositoryName},
}

func TestNewSetupCommand(t *testing.T) {
	cmd := NewSetupCommand()

	for _, data := range dataset {
		_ = cmd.Flags().Set("strategy", data.Strategy)
		_ = cmd.Flags().Set("provider", data.Provider)
		_ = cmd.Flags().Set("repository", data.Repository)
		_ = cmd.Flags().Set("branch", data.Branch)

		global.LocatorConfigFile = TmpConfig(t).Name()
		defer os.Remove(global.LocatorConfigFile)

		err := cmd.Execute()
		if err != data.Error {
			if data.Error == nil {
				t.Errorf("Expected no error, got %s", err)
			} else {
				t.Errorf("Expected %s, got %s", data.Error, err)
			}
		}

		_ = os.Remove(global.LocatorConfigFile)
	}

}

func TestNewSetupCommand2(t *testing.T) {
	cmd := NewSetupCommand()

	originalStdin := util.Stdin
	originalStdout := util.Stdout

	for _, data := range dataset {
		if data.Error != nil {
			continue
		}

		util.Stdin = bytes.NewBufferString(data.Strategy + "\n" + data.Provider + "\n" + data.Repository + "\n" + data.Branch + "\n")
		util.Stdout = new(bytes.Buffer)

		global.LocatorConfigFile = TmpConfig(t).Name()
		defer os.Remove(global.LocatorConfigFile)
		err := cmd.Execute()
		if err != data.Error {
			if data.Error == nil {
				t.Errorf("Expected no error, got %s, output: %s", err, util.Stdout.(*bytes.Buffer).String())
			} else {
				t.Errorf("Expected %s, got %s, output: %s", data.Error, err, util.Stdout.(*bytes.Buffer).String())
			}
		}

	}

	util.Stdin = originalStdin
	util.Stdout = originalStdout
}

func TmpConfig(t *testing.T) *os.File {
	f, err := os.Create("/tmp/" + strconv.Itoa(int(time.Now().UnixNano())) + ".tmp")
	if err != nil {
		t.Fatalf("Error creating tmp file: %s", err)
	}

	return f
}
