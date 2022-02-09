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

var tests = []Set{
	{"remote", "github", "felixdorn/config-test", "main", "", nil},
	{"remote", "gitlab", "felixdorn/config-test", "main", "", nil},
	{"remote", "bitbucket", "felixdorn/config-test", "main", "", nil},
	{"invalidStrategy", "github", "felixdorn/config-test", "main", "", pkg.ErrInvalidStrategy},
	{"remote", "invalidProvider", "felixdorn/config-test", "main", "", pkg.ErrInvalidProvider},
	{"remote", "github", "invalidRepository", "main", "", pkg.ErrInvalidRepositoryName},
}

func TestNewSetupCommand(t *testing.T) {
	cmd := NewSetupCommand()

	for _, test := range tests {
		_ = cmd.Flags().Set("strategy", test.Strategy)
		_ = cmd.Flags().Set("provider", test.Provider)
		_ = cmd.Flags().Set("repository", test.Repository)
		_ = cmd.Flags().Set("branch", test.Branch)

		global.LocatorConfigFile = TmpConfig(t).Name()
		defer os.Remove(global.LocatorConfigFile)

		err := cmd.Execute()
		if err != test.Error {
			if test.Error == nil {
				t.Errorf("Expected no error, got %s", err)
			} else {
				t.Errorf("Expected %s, got %s", test.Error, err)
			}
		}

		_ = os.Remove(global.LocatorConfigFile)
	}

}

func TestNewSetupCommand2(t *testing.T) {
	cmd := NewSetupCommand()

	originalStdin := util.Stdin
	originalStdout := util.Stdout

	for _, test := range tests {
		if test.Error != nil {
			continue
		}

		util.Stdin = bytes.NewBufferString(test.Strategy + "\n" + test.Provider + "\n" + test.Repository + "\n" + test.Branch + "\n")
		util.Stdout = new(bytes.Buffer)

		global.LocatorConfigFile = TmpConfig(t).Name()
		defer os.Remove(global.LocatorConfigFile)
		err := cmd.Execute()
		if err != test.Error {
			if test.Error == nil {
				t.Errorf("Expected no error, got %s, output: %s", err, util.Stdout.(*bytes.Buffer).String())
			} else {
				t.Errorf("Expected %s, got %s, output: %s", test.Error, err, util.Stdout.(*bytes.Buffer).String())
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
