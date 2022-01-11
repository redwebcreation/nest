package cli

import (
	"bytes"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/redwebcreation/nest/common"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/util"
)

type Set struct {
	Strategy   string
	Provider   string
	Repository string
	Error      error
}

var dataset = []Set{
	{"remote", "github", "felixdorn/config-test", nil},
	{"remote", "gitlab", "felixdorn/config-test", nil},
	{"remote", "bitbucket", "felixdorn/config-test", nil},
	{"invalidStrategy", "github", "felixdorn/config-test", common.ErrInvalidStrategy},
	{"remote", "invalidProvider", "felixdorn/config-test", common.ErrInvalidProvider},
	{"remote", "github", "invalidRepository", common.ErrInvalidRepository},
}

func Test_ConfigureCommandRuns(t *testing.T) {
	cmd := ConfigCommand()
	cmd.Execute()
}

func Test_ConfigureCommandUsingFlags(t *testing.T) {
	cmd := ConfigureCommand()

	for _, data := range dataset {
		cmd.Flags().Set("strategy", data.Strategy)
		cmd.Flags().Set("provider", data.Provider)
		cmd.Flags().Set("repository", data.Repository)

		global.ConfigLocatorConfigFile = "/tmp/" + strconv.Itoa(int(time.Now().UnixMilli())) + ".json"

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

func Test_ConfigureCommandInteractively(t *testing.T) {
	cmd := ConfigureCommand()

	originalStdin := util.Stdin

	for _, data := range dataset {
		if data.Error != nil {
			continue
		}

		util.Stdin = bytes.NewBufferString(data.Strategy + "\n" + data.Provider + "\n" + data.Repository + "\n")

		global.ConfigLocatorConfigFile = "/tmp/" + strconv.Itoa(int(time.Now().UnixMilli())) + ".json"

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

	util.Stdin = originalStdin
}
