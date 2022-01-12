package command

import (
	"fmt"
	"os"

	"github.com/redwebcreation/nest/common"
	"github.com/redwebcreation/nest/global"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func WithConfig(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentPreRunE = LoadConf

	return cmd
}

func LoadConf(cmd *cobra.Command, args []string) error {
	commandName := cmd.Name()

	if _, err := os.Stat(global.ConfigLocatorConfigFile); err != nil {
		if commandName == "configure" {
			return nil
		}

		return fmt.Errorf("run `nest configure` to setup nest")
	}

	reader, err := common.LoadConfigReader()
	if err != nil {
		return err
	}

	common.ConfigLocator = reader

	contents, err := reader.Read("nest.yaml")
	if err != nil {
		return err
	}

	var config common.Configuration

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return err
	}

	common.Config = &config

	if commandName == "medic" {
		return nil
	}

	diagnosis := common.DiagnoseConfiguration()

	if len(diagnosis.Errors) == 0 {
		return nil
	}

	return fmt.Errorf("your configuration is invalid, please run `nest medic` to troubleshoot")
}
