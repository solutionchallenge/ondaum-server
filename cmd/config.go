package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/spf13/cobra"
)

func NewConfigCommand() *cobra.Command {
	configCommand := &cobra.Command{
		Use:   "config",
		Short: "Show current config",
		Long:  "Show current config",
		Run: func(cmd *cobra.Command, args []string) {
			cfgpath, err := cmd.Flags().GetString("config-path")
			if err != nil {
				panic(err)
			}
			cfgname, err := cmd.Flags().GetString("config-name")
			if err != nil {
				panic(err)
			}
			cfg := utils.LoadConfig(cfgname, cfgpath)
			output, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println("Current Config")
			fmt.Println(string(output))
		},
	}
	configCommand.Flags().StringP("config-path", "p", "./config", "config file path (default is './config')")
	configCommand.Flags().StringP("config-name", "n", "production", "config file name (default is 'production')")
	return configCommand
}
