package command

import (
	"bytes"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/redwebcreation/nest/util"
	"os"
	"testing"
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

func TestConfigureCommandUsingFlags(t *testing.T) {
	cmd := NewConfigureCommand()

	for _, data := range dataset {
		_ = cmd.Flags().Set("strategy", data.Strategy)
		_ = cmd.Flags().Set("provider", data.Provider)
		_ = cmd.Flags().Set("repository", data.Repository)
		_ = cmd.Flags().Set("branch", data.Branch)

		global.ConfigLocatorConfigFile = util.TmpFile().Name()

		err := cmd.Execute()
		if err != data.Error {
			if data.Error == nil {
				t.Errorf("Expected no error, got %s", err)
			} else {
				t.Errorf("Expected %s, got %s", data.Error, err)
			}
		}

		_ = os.Remove(global.ConfigLocatorConfigFile)
	}

}

func TestConfigureCommandInteractively(t *testing.T) {
	cmd := NewConfigureCommand()

	originalStdin := util.Stdin

	for _, data := range dataset {
		if data.Error != nil {
			continue
		}

		util.Stdin = bytes.NewBufferString(data.Strategy + "\n" + data.Provider + "\n" + data.Repository + "\n" + data.Branch + "\n")

		global.ConfigLocatorConfigFile = util.TmpFile().Name()
		defer os.Remove(global.ConfigLocatorConfigFile)
		err := cmd.Execute()
		if err != data.Error {
			if data.Error == nil {
				t.Errorf("Expected no error, got %s", err)
			} else {
				t.Errorf("Expected %s, got %s", data.Error, err)
			}
		}

	}

	util.Stdin = originalStdin
}
