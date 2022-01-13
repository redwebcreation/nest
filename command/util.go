package command

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/redwebcreation/nest/common"
	"github.com/redwebcreation/nest/global"
)

func LoadConfigFromCommit(commit string) error {
	reader := common.LocatorConfig{
		Commit: commit,
	}

	contents, err := os.ReadFile(global.ConfigLocatorConfigFile)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(contents, &reader); err != nil && err.Error() == "unknown error: remote: " {
		return fmt.Errorf("the repository %s does not exists", reader.GetRepositoryLocation())
	}

	common.ConfigLocator = &reader

	contents, err = reader.Read("nest.yaml")
	if err != nil {
		return err
	}

	var config common.Configuration

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return err
	}

	common.Config = &config
	return nil
}

func LoadConfig() error {
	return LoadConfigFromCommit("")
}
