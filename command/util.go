package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
)

func LoadConfigFromCommit(commit string) error {
	reader := pkg.ConfigLocator{
		ConfigLocatorConfig: pkg.ConfigLocatorConfig{
			Commit: commit,
		},
	}

	contents, err := os.ReadFile(global.ConfigLocatorConfigFile)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(contents, &reader); err != nil && err.Error() == "unknown error: remote: " {
		return fmt.Errorf("the repository %s does not exists", reader.GetRepositoryLocation())
	}

	pkg.Config = &reader

	return nil
}

func LoadConfig() error {
	return LoadConfigFromCommit("")
}
